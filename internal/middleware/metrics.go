package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	requestCount    int64
	requestDuration float64
)

func PrometheusMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		requestCount++
		requestDuration += duration

		// Store metrics for the /metrics endpoint
		c.Set("metric_duration", duration)
	}
}

func GetMetrics() string {
	return `# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total ` + strconv.FormatInt(requestCount, 10) + `

# HELP http_request_duration_seconds Total HTTP request duration in seconds
# TYPE http_request_duration_seconds counter
http_request_duration_seconds ` + strconv.FormatFloat(requestDuration, 'f', 6, 64) + `

# HELP teams_total Total number of teams
# TYPE teams_total gauge
teams_total 0

# HELP players_total Total number of players
# TYPE players_total gauge
players_total 0

# HELP matches_total Total number of matches
# TYPE matches_total gauge
matches_total 0

# HELP matches_completed_total Total number of completed matches
# TYPE matches_completed_total gauge
matches_completed_total 0
`
}
