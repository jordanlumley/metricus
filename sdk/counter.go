package metricus

import (
	"sync"
)

// Counter represents a single counter metric
type Counter struct {
	Value int
	mu    sync.Mutex
}

// Inc increments the counter
func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Value++
}
