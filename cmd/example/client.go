package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	metricus "github.com/jordanlumley/metricus/sdk"
)

func main() {
	// Initialize the metric store with auto-flushing every 10 seconds
	// BadgerDB will be managed internally by the SDK in ./data directory
	store, err := metricus.NewMetricStore(metricus.DefaultMetricStoreOpts)
	if err != nil {
		log.Fatalf("Failed to initialize metric store: %v", err)
	}

	// Define a new counter
	opsProcessed := store.NewCounter("ops_processed_total")
	specificOp1Processed := store.NewCounter("specific_op1_processed_total")
	specificOp2Processed := store.NewCounter("specific_op2_processed_total")

	// Set up the router for the user's REST API
	r := mux.NewRouter()
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			opsProcessed.Inc()
			h.ServeHTTP(w, r)
		})
	})
	// Integrate the metrics endpoint into the user's API
	r.Handle("/metrics", store.ExposeMetricsHandler())
	r.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	r.Handle("/op1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		specificOp1Processed.Inc()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	r.Handle("/op2", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		specificOp2Processed.Inc()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
