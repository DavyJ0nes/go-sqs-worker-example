package metrics

import (
	"net/http"
	"strconv"

	"github.com/felixge/httpsnoop"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requests *prometheus.CounterVec
	duration *prometheus.HistogramVec
)

// Initialise ensures that the metrics componant is configured correctly
// This should be called from main.go init() function.
//
// For production service this should probably be pulled out as a dependency
// that should be injected at a higher level (maybe within router)
func Initialise() {
	requests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Count of all HTTP requests",
	}, []string{"method", "route", "status_code"})

	duration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "request_duration_seconds",
		Help:    "Time (in secods) spent serving HTTP requests",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "route", "status_code"})

	prometheus.MustRegister(requests)
	prometheus.MustRegister(duration)
}

// Measure is a middleware function for measuring the latency of the application
func Measure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// do not measure requests to /metrics
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		m := httpsnoop.CaptureMetrics(next, w, r)

		// measure request counts
		requests.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(m.Code)).Add(1)
		// measure request duration
		duration.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(m.Code)).Observe(m.Duration.Seconds())
	})
}

// Handler is used to expose prometheus metrics
func Handler() http.Handler {
	return promhttp.Handler()
}
