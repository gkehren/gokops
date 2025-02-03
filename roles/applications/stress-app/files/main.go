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
}

func stressHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestCounter.Inc()

	// Prepare response using buffer pool
	buf := bufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()

	buf.WriteString(responsePrefix)
	buf.WriteString(time.Since(start).String())
	buf.WriteString(responseSuffix)

	w.Write(buf.Bytes())
	responseTimeHistogram.Observe(time.Since(start).Seconds())
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
