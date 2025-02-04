package main

import (
	"bytes"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Pre-allocate metrics with fixed labels
var (
	requestCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name:        "app_requests_total",
			Help:        "Total number of requests received",
			ConstLabels: prometheus.Labels{"endpoint": "root"},
		},
	)

	responseTimeHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:        "app_response_time_seconds",
			Help:        "Response time in seconds",
			Buckets:     prometheus.ExponentialBuckets(0.0001, 2, 16),
			ConstLabels: prometheus.Labels{"endpoint": "root"},
		},
	)

	activeRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name:        "app_active_requests",
			Help:        "Current number of active requests",
			ConstLabels: prometheus.Labels{"endpoint": "root"},
		},
	)

	responseSizeHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:        "app_response_size_bytes",
			Help:        "Size of responses in bytes",
			Buckets:     prometheus.ExponentialBuckets(128, 2, 10),
			ConstLabels: prometheus.Labels{"endpoint": "root"},
		},
	)
)

// Buffer pool for response building
var bufPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 128))
	},
}

const (
	responsePrefix = "Hello, world! Processed in "
	responseSuffix = "\n"
)

func init() {
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(responseTimeHistogram)
	prometheus.MustRegister(activeRequests)
	prometheus.MustRegister(responseSizeHistogram)
}

func stressHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestCounter.Inc()
	activeRequests.Inc()
	defer activeRequests.Dec()

	buf := bufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()

	buf.WriteString(responsePrefix)
	buf.WriteString(time.Since(start).String())
	buf.WriteString(responseSuffix)

	respBytes := buf.Bytes()
	w.Write(respBytes)
	responseTimeHistogram.Observe(time.Since(start).Seconds())
	responseSizeHistogram.Observe(float64(len(respBytes)))
}

func main() {
	// Create a single router that handles both API and metrics paths.
	mux := http.NewServeMux()
	mux.HandleFunc("/", stressHandler)
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    ":4242",
		Handler: mux,

		// Timeouts to protect against slow clients
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,

		// Optimize for high traffic
		MaxHeaderBytes: 1 << 18, // 256KB
	}

	server.ListenAndServe()
}
