package agent

import (
	"fmt"
	"net/http"
)

func UpgradeSSE(res http.ResponseWriter) {
	res.Header().Add("Content-Type", "text/event-stream")
	res.Header().Add("Cache-Control", "no-cache")
	res.Header().Add("Connection", "keep-alive")
}

func SendSSE(res http.ResponseWriter, data []byte) error {
	_, err := fmt.Fprintf(res, "data: %v\n\n", string(data))
	if err != nil {
		return fmt.Errorf("failed to send SSE: %w", err)
	}
	res.(http.Flusher).Flush()

	return nil
}
