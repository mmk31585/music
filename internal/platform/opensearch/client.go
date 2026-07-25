package opensearch

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"go.uber.org/zap"
)

type Config struct {
	URL      string
	Username string
	Password string
	Index    string
}

type circuitState int32

const (
	circuitClosed   circuitState = 0
	circuitOpen     circuitState = 1
	circuitHalfOpen circuitState = 2
)

type circuitBreaker struct {
	state         circuitState
	failures      int64
	lastFailure   time.Time
	threshold     int64
	resetTimeout  time.Duration
	halfOpenMax   int64
	halfOpenCount int64
	mu            sync.Mutex
}

func newCircuitBreaker() *circuitBreaker {
	return &circuitBreaker{
		state:        circuitClosed,
		threshold:    5,
		resetTimeout: 30 * time.Second,
		halfOpenMax:  3,
	}
}

func (cb *circuitBreaker) allowRequest() bool {
	state := circuitState(atomic.LoadInt32((*int32)(&cb.state)))
	switch state {
	case circuitClosed:
		return true
	case circuitOpen:
		cb.mu.Lock()
		defer cb.mu.Unlock()
		if time.Since(cb.lastFailure) > cb.resetTimeout {
			atomic.StoreInt32((*int32)(&cb.state), int32(circuitHalfOpen))
			atomic.StoreInt64(&cb.halfOpenCount, 0)
			return true
		}
		return false
	case circuitHalfOpen:
		count := atomic.AddInt64(&cb.halfOpenCount, 1)
		return count <= cb.halfOpenMax
	}
	return false
}

func (cb *circuitBreaker) recordSuccess() {
	state := circuitState(atomic.LoadInt32((*int32)(&cb.state)))
	if state == circuitHalfOpen {
		cb.mu.Lock()
		atomic.StoreInt32((*int32)(&cb.state), int32(circuitClosed))
		atomic.StoreInt64(&cb.failures, 0)
		cb.mu.Unlock()
	}
}

func (cb *circuitBreaker) recordFailure() {
	atomic.AddInt64(&cb.failures, 1)
	cb.mu.Lock()
	cb.lastFailure = time.Now()

	failures := atomic.LoadInt64(&cb.failures)
	if failures >= cb.threshold {
		atomic.StoreInt32((*int32)(&cb.state), int32(circuitOpen))
	}
	cb.mu.Unlock()
}

func (cb *circuitBreaker) stateString() string {
	switch circuitState(atomic.LoadInt32((*int32)(&cb.state))) {
	case circuitClosed:
		return "closed"
	case circuitOpen:
		return "open"
	case circuitHalfOpen:
		return "half-open"
	}
	return "unknown"
}

type Client struct {
	client  *opensearch.Client
	logger  *zap.Logger
	healthy bool
	breaker *circuitBreaker
}

func NewClient(cfg Config, logger *zap.Logger) (*Client, error) {
	if cfg.URL == "" {
		return nil, nil
	}

	osCfg := opensearch.Config{
		Addresses: []string{cfg.URL},
		Username:  cfg.Username,
		Password:  cfg.Password,
		Transport: &http.Transport{
			MaxIdleConns:        20,
			MaxIdleConnsPerHost: 20,
			IdleConnTimeout:     90 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		MaxRetries: 2,
		RetryBackoff: func(i int) time.Duration {
			return time.Duration(i*100) * time.Millisecond
		},
		RetryOnStatus: []int{429, 502, 503, 504},
	}

	client, err := opensearch.NewClient(osCfg)
	if err != nil {
		return nil, fmt.Errorf("opensearch client: %w", err)
	}

	c := &Client{
		client:  client,
		logger:  logger,
		breaker: newCircuitBreaker(),
	}

	if err := c.healthCheck(context.Background()); err != nil {
		logger.Warn("opensearch health check failed, continuing degraded", zap.Error(err))
	} else {
		c.healthy = true
	}

	return c, nil
}

func (c *Client) healthCheck(ctx context.Context) error {
	if !c.breaker.allowRequest() {
		return fmt.Errorf("circuit breaker open")
	}

	resp, err := c.client.Info(
		c.client.Info.WithContext(ctx),
	)
	if err != nil {
		c.breaker.recordFailure()
		c.healthy = false
		return err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		c.breaker.recordFailure()
		c.healthy = false
		return fmt.Errorf("opensearch unhealthy: %s", resp.String())
	}

	c.breaker.recordSuccess()
	c.healthy = true
	return nil
}

func (c *Client) IsHealthy() bool {
	return c.healthy && c.breaker.allowRequest()
}

func (c *Client) Search(ctx context.Context, index []string, body any) (*opensearchapi.Response, error) {
	if !c.breaker.allowRequest() {
		return nil, fmt.Errorf("opensearch circuit breaker open")
	}
	resp, err := Search(ctx, c.client, index, body)
	if err != nil {
		c.breaker.recordFailure()
		return nil, err
	}
	c.breaker.recordSuccess()
	return resp, nil
}

func Search(ctx context.Context, client *opensearch.Client, index []string, body any) (*opensearchapi.Response, error) {
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal search body: %w", err)
	}

	return client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(index...),
		client.Search.WithBody(strings.NewReader(string(raw))),
		client.Search.WithTrackTotalHits(false),
	)
}

func (c *Client) IndexDocument(ctx context.Context, index, documentID string, body any) error {
	if !c.breaker.allowRequest() {
		return fmt.Errorf("opensearch circuit breaker open")
	}
	err := IndexDocument(ctx, c.client, index, documentID, body)
	if err != nil {
		c.breaker.recordFailure()
		return err
	}
	c.breaker.recordSuccess()
	return nil
}

func IndexDocument(ctx context.Context, client *opensearch.Client, index, documentID string, body any) error {
	raw, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal index body: %w", err)
	}

	req := opensearchapi.IndexRequest{
		Index:      index,
		DocumentID: documentID,
		Body:       strings.NewReader(string(raw)),
	}

	resp, err := req.Do(ctx, client)
	if err != nil {
		return fmt.Errorf("index request: %w", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("index error: %s", string(b))
	}

	return nil
}

func (c *Client) Delete(ctx context.Context, index, documentID string) error {
	if !c.breaker.allowRequest() {
		return fmt.Errorf("opensearch circuit breaker open")
	}
	req := opensearchapi.DeleteRequest{
		Index:      index,
		DocumentID: documentID,
	}

	resp, err := req.Do(ctx, c.client)
	if err != nil {
		c.breaker.recordFailure()
		return err
	}
	defer resp.Body.Close()
	c.breaker.recordSuccess()
	return nil
}

func (c *Client) Bulk(ctx context.Context, body *strings.Reader) error {
	if !c.breaker.allowRequest() {
		return fmt.Errorf("opensearch circuit breaker open")
	}
	req := opensearchapi.BulkRequest{
		Body: body,
	}

	resp, err := req.Do(ctx, c.client)
	if err != nil {
		c.breaker.recordFailure()
		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		c.breaker.recordFailure()
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bulk error: %s", string(b))
	}

	c.breaker.recordSuccess()
	return nil
}

func (c *Client) IndexExists(ctx context.Context, name string) (bool, error) {
	if !c.breaker.allowRequest() {
		return false, fmt.Errorf("opensearch circuit breaker open")
	}
	resp, err := c.client.Indices.Exists([]string{name})
	if err != nil {
		c.breaker.recordFailure()
		return false, err
	}
	defer resp.Body.Close()
	c.breaker.recordSuccess()
	return resp.StatusCode != 404, nil
}

func (c *Client) CreateIndex(ctx context.Context, name, mapping string) error {
	if !c.breaker.allowRequest() {
		return fmt.Errorf("opensearch circuit breaker open")
	}
	exists, err := c.IndexExists(ctx, name)
	if err != nil {
		c.breaker.recordFailure()
		return err
	}
	if exists {
		c.breaker.recordSuccess()
		return nil
	}

	req := opensearchapi.IndicesCreateRequest{
		Index: name,
		Body:  strings.NewReader(mapping),
	}

	resp, err := req.Do(ctx, c.client)
	if err != nil {
		c.breaker.recordFailure()
		return fmt.Errorf("create index %s: %w", name, err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		c.breaker.recordFailure()
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("create index %s error: %s", name, string(b))
	}

	c.breaker.recordSuccess()
	return nil
}

func (c *Client) Client() *opensearch.Client {
	return c.client
}

func (c *Client) Close() {
	c.healthy = false
}
