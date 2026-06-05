package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "music_http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "music_http_request_duration_seconds",
			Help:    "HTTP request latency in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "music_http_requests_in_flight",
			Help: "Number of HTTP requests currently being processed.",
		},
	)

	WorkerRunsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "music_worker_runs_total",
			Help: "Total number of worker job runs.",
		},
		[]string{"worker", "status"},
	)

	WorkerRunDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "music_worker_run_duration_seconds",
			Help:    "Worker job duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"worker"},
	)

	EventsPublishedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "music_events_published_total",
			Help: "Total number of domain events published.",
		},
		[]string{"event_type"},
	)

	DatabaseQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "music_database_queries_total",
			Help: "Total number of database queries.",
		},
		[]string{"operation", "status"},
	)

	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "music_database_query_duration_seconds",
			Help:    "Database query latency in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)
)

func HTTPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		HTTPRequestsInFlight.Inc()
		defer HTTPRequestsInFlight.Dec()

		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		duration := time.Since(start).Seconds()

		HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
		HTTPRequestDuration.WithLabelValues(method, path, status).Observe(duration)
	}
}

func ObserveWorkerRun(workerName string, start time.Time, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}

	WorkerRunsTotal.WithLabelValues(workerName, status).Inc()
	WorkerRunDuration.WithLabelValues(workerName).Observe(time.Since(start).Seconds())
}

func IncEventPublished(eventType string) {
	EventsPublishedTotal.WithLabelValues(eventType).Inc()
}

func ObserveDatabaseQuery(operation string, start time.Time, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}

	DatabaseQueriesTotal.WithLabelValues(operation, status).Inc()
	DatabaseQueryDuration.WithLabelValues(operation).Observe(time.Since(start).Seconds())
}
