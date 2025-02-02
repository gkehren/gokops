// main.go
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define a counter metric to track the number of requests.
var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_requests_total",
			Help: "Total number of requests received",
		},
		[]string{"endpoint"},
	)
	// A histogram to record response times.
	responseTimeHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "app_response_time_seconds",
			Help:    "Response time in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)
)

func init() {
	// Register Prometheus metrics.
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(responseTimeHistogram)
}

func stressHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	requestCounter.WithLabelValues(r.URL.Path).Inc()

	// Simulate load: random sleep between 10 and 100ms.
	sleepDuration := time.Duration(rand.Intn(90)+10) * time.Millisecond
	time.Sleep(sleepDuration)

	// Log the request (this log can be picked up by Loki).
	log.Printf("Processed request for %s in %v", r.URL.Path, time.Since(startTime))

	// Record the response time.
	responseTimeHistogram.WithLabelValues(r.URL.Path).Observe(time.Since(startTime).Seconds())

	fmt.Fprintf(w, "Hello, world! Processed in %v\n", time.Since(startTime))
}

func main() {
	// Seed the random generator.
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", stressHandler)
	// Expose Prometheus metrics.
	http.Handle("/metrics", promhttp.Handler())

	port := "4242"
	log.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
