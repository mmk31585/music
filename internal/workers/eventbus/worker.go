package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Worker struct {
	db     *sqlx.DB
	rdb    *redis.Client
	logger *zap.Logger
}

func NewWorker(db *sqlx.DB, rdb *redis.Client, logger *zap.Logger) *Worker {
	return &Worker{db: db, rdb: rdb, logger: logger}
}

func (w *Worker) Name() string { return "eventbus" }

func (w *Worker) Run(ctx context.Context) error {
	defer func() {
		if r := recover(); r != nil {
			w.logger.Error("eventbus worker panicked", zap.Any("recover", r))
		}
	}()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := w.processBatch(ctx); err != nil {
				w.logger.Error("eventbus worker failed", zap.Error(err))
			}
		}
	}
}

type eventRow struct {
	ID        string          `db:"id"`
	EventType string          `db:"event_type"`
	Payload   json.RawMessage `db:"payload"`
	TraceID   *string         `db:"trace_id"`
}

func (w *Worker) processBatch(ctx context.Context) error {
	rows, err := w.db.QueryContext(ctx, `
		SELECT id, event_type, payload
		FROM event_outbox
		WHERE status = 'pending'
		ORDER BY created_at ASC
		LIMIT 100
		FOR UPDATE SKIP LOCKED
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var e eventRow
		if err := rows.Scan(&e.ID, &e.EventType, &e.Payload); err != nil {
			continue
		}

		if err := w.publishToStream(ctx, e); err != nil {
			w.logger.Error("stream publish failed, marking failed",
				zap.String("id", e.ID), zap.Error(err))
			w.markFailed(ctx, e.ID, err.Error())
			continue
		}

		if _, err := w.db.ExecContext(ctx,
			`UPDATE event_outbox SET status = 'processed', processed_at = NOW() WHERE id = $1`,
			e.ID,
		); err != nil {
			w.logger.Error("failed to mark event processed",
				zap.String("id", e.ID), zap.Error(err))
		}
	}

	return rows.Err()
}

func (w *Worker) publishToStream(ctx context.Context, e eventRow) error {
	if w.rdb == nil {
		return nil
	}

	streamName := fmt.Sprintf("stream:%s", e.EventType)
	values := map[string]string{
		"id":         e.ID,
		"event_type": e.EventType,
		"payload":    string(e.Payload),
		"timestamp":  fmt.Sprintf("%d", time.Now().UnixMilli()),
	}

	if err := w.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: values,
		MaxLen: 10000,
		Approx: true,
	}).Err(); err != nil {
		return fmt.Errorf("xadd %s: %w", streamName, err)
	}

	return nil
}

func (w *Worker) markFailed(ctx context.Context, id string, errMsg string) {
	w.db.ExecContext(ctx,
		`UPDATE event_outbox SET status = 'failed', retry_count = retry_count + 1, last_error = $2
		 WHERE id = $1`, id, errMsg)
}

type ConsumerGroupWorker struct {
	rdb      *redis.Client
	logger   *zap.Logger
	group    string
	consumer string
}

func NewConsumerGroupWorker(rdb *redis.Client, logger *zap.Logger, group, consumer string) *ConsumerGroupWorker {
	return &ConsumerGroupWorker{
		rdb:      rdb,
		logger:   logger,
		group:    group,
		consumer: consumer,
	}
}

func (cgw *ConsumerGroupWorker) Run(ctx context.Context, eventTypes []string, handler func(ctx context.Context, eventType string, payload json.RawMessage) error) error {
	for _, et := range eventTypes {
		streamName := fmt.Sprintf("stream:%s", et)
		if err := cgw.rdb.XGroupCreateMkStream(ctx, streamName, cgw.group, "0").Err(); err != nil &&
			err.Error() != "BUSYGROUP Consumer Group name already exists" {
			cgw.logger.Warn("xgroup create", zap.String("stream", streamName), zap.Error(err))
		}
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			for _, et := range eventTypes {
				if err := cgw.consumeStream(ctx, et, handler); err != nil {
					cgw.logger.Error("consume stream failed",
						zap.String("event_type", et), zap.Error(err))
				}
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func (cgw *ConsumerGroupWorker) consumeStream(ctx context.Context, eventType string, handler func(ctx context.Context, eventType string, payload json.RawMessage) error) error {
	streamName := fmt.Sprintf("stream:%s", eventType)

	result, err := cgw.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    cgw.group,
		Consumer: cgw.consumer,
		Streams:  []string{streamName, ">"},
		Count:    10,
		Block:    0,
	}).Result()
	if err != nil {
		return err
	}

	for _, stream := range result {
		for _, message := range stream.Messages {
			id := message.ID
			payloadStr, ok := message.Values["payload"].(string)
			if !ok {
				cgw.rdb.XAck(ctx, streamName, cgw.group, id)
				continue
			}

			if err := handler(ctx, eventType, json.RawMessage(payloadStr)); err != nil {
				cgw.logger.Error("handler failed",
					zap.String("event_type", eventType),
					zap.String("message_id", id),
					zap.Error(err))
				continue
			}

			cgw.rdb.XAck(ctx, streamName, cgw.group, id)
		}
	}

	return nil
}
