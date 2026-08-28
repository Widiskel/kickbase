package middleware

import (
	"strconv"
	"time"

	"kickbase/internal/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests processed",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency distributions",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	teamsTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "teams_total",
			Help: "Total number of teams",
		},
	)

	playersTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "players_total",
			Help: "Total number of players",
		},
	)

	matchesTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "matches_total",
			Help: "Total number of matches",
		},
	)

	matchesCompletedTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "matches_completed_total",
			Help: "Total number of completed matches",
		},
	)
)

func init() {
	prometheus.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
		teamsTotal,
		playersTotal,
		matchesTotal,
		matchesCompletedTotal,
	)
}

// PrometheusMetrics middleware collects request count and latency
func PrometheusMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		httpRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)
	}
}

// PrometheusHandler returns a Gin handler that serves Prometheus metrics and updates DB gauges
func PrometheusHandler(
	teamRepo interfaces.TeamRepository,
	playerRepo interfaces.PlayerRepository,
	matchRepo interfaces.MatchRepository,
) gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		if teamRepo != nil {
			if count, err := teamRepo.CountTotal(); err == nil {
				teamsTotal.Set(float64(count))
			}
		}
		if playerRepo != nil {
			if count, err := playerRepo.CountTotal(); err == nil {
				playersTotal.Set(float64(count))
			}
		}
		if matchRepo != nil {
			if count, err := matchRepo.CountTotal(); err == nil {
				matchesTotal.Set(float64(count))
			}
			if count, err := matchRepo.CountCompleted(); err == nil {
				matchesCompletedTotal.Set(float64(count))
			}
		}

		h.ServeHTTP(c.Writer, c.Request)
	}
}

// GetMetrics returns string format metrics (for tests)
func GetMetrics() string {
	return "# Prometheus standard metrics enabled at /metrics\n"
}
