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

	// Import search metrics
	ImportSearchTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "music_import_search_total",
			Help: "Total number of import search requests.",
		},
		[]string{"source", "cached"},
	)

	ImportSearchDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "music_import_search_duration_seconds",
			Help:    "Import search provider latency in seconds.",
			Buckets: []float64{.05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"source"},
	)

	ImportSearchErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "music_import_search_errors_total",
			Help: "Total number of import search provider errors.",
		},
		[]string{"source"},
	)

	ImportCacheOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "music_import_cache_operations_total",
			Help: "Total import search cache operations.",
		},
		[]string{"operation"},
	)

	// Acquisition metrics
	ImportAcquisitionTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "music_import_acquisition_total",
			Help: "Total number of acquisition resolve attempts.",
		},
		[]string{"source", "status"},
	)

	ImportAcquisitionDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "music_import_acquisition_duration_seconds",
			Help:    "Acquisition resolver latency in seconds.",
			Buckets: []float64{1, 2.5, 5, 10, 15, 30},
		},
		[]string{"source"},
	)

	// Import job metrics
	ImportJobsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "music_import_jobs_total",
			Help: "Total number of import job outcomes.",
		},
		[]string{"status"},
	)

	ImportJobDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "music_import_job_duration_seconds",
			Help:    "Import job duration by stage in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"stage"},
	)

	ImportJobsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "music_import_jobs_in_flight",
			Help: "Number of import jobs currently being processed.",
		},
	)

	// Circuit breaker metrics
	ImportCircuitBreakerState = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "music_import_circuit_breaker_state",
			Help: "Circuit breaker state per provider (0=closed, 1=half-open, 2=open).",
		},
		[]string{"provider"},
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

func ObserveImportSearch(source string, start time.Time) {
	ImportSearchDuration.WithLabelValues(source).Observe(time.Since(start).Seconds())
}

func IncImportSearchTotal(source string, cached bool) {
	val := "false"
	if cached {
		val = "true"
	}
	ImportSearchTotal.WithLabelValues(source, val).Inc()
}

func IncImportSearchError(source string) {
	ImportSearchErrorsTotal.WithLabelValues(source).Inc()
}

func IncImportCacheOp(op string) {
	ImportCacheOperationsTotal.WithLabelValues(op).Inc()
}

func IncImportAcquisition(source string, status string, start time.Time) {
	ImportAcquisitionTotal.WithLabelValues(source, status).Inc()
	ImportAcquisitionDuration.WithLabelValues(source).Observe(time.Since(start).Seconds())
}

func IncImportJob(status string) {
	ImportJobsTotal.WithLabelValues(status).Inc()
}

func ObserveImportJobStage(stage string, start time.Time) {
	ImportJobDuration.WithLabelValues(stage).Observe(time.Since(start).Seconds())
}

func SetImportJobsInFlight(n int) {
	ImportJobsInFlight.Set(float64(n))
}

func SetCircuitBreakerState(provider string, state int) {
	ImportCircuitBreakerState.WithLabelValues(provider).Set(float64(state))
}
