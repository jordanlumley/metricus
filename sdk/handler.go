package metricus

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ExposeMetricsHandler returns an HTTP handler that exposes metrics
func (ms *MetricStore) ExposeMetricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics := ms.GetMetrics()

		// convert key value to json object
		data := make(map[string]int)
		for key, metric := range metrics {
			data[key] = metric.Value
		}

		// Convert the map to a JSON object (marshal it)
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	}
}
