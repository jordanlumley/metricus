package metricus

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/dgraph-io/badger/v3"
)

var (
	DefaultMetricStoreOpts = MetricStoreOpts{
		FlushIntervalSeconds: 10,
		Path:                 "/tmp/data",
	}
)

type MetricStoreOpts struct {
	FlushIntervalSeconds int
	Path                 string
}

// MetricStore manages all metrics in memory and persistently stores them using BadgerDB
type MetricStore struct {
	metrics map[string]*Counter
	mu      sync.Mutex
	db      *badger.DB
	stopCh  chan struct{}
}

// NewMetricStore initializes BadgerDB internally and sets up metric storage
func NewMetricStore(opts MetricStoreOpts) (*MetricStore, error) {
	dbOpts := badger.DefaultOptions(opts.Path)
	db, err := badger.Open(dbOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to open BadgerDB: %v", err)
	}

	// Create the metric store and start automatic flushing
	store := &MetricStore{
		metrics: make(map[string]*Counter),
		db:      db,
		stopCh:  make(chan struct{}),
	}

	// Load saved metrics from BadgerDB
	if err := store.load(); err != nil {
		return nil, err
	}

	// Start the background auto-flush
	go store.startAutoFlush(opts.FlushIntervalSeconds)

	return store, nil
}

// NewCounter creates a new counter
func (ms *MetricStore) NewCounter(name string) *Counter {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	// Check if the counter already exists
	if counter, exists := ms.metrics[name]; exists {
		return counter // Return the existing counter
	}

	// Otherwise, create a new counter
	counter := &Counter{}
	ms.metrics[name] = counter
	return counter
}

// GetMetrics retrieves all metrics for exposure
func (ms *MetricStore) GetMetrics() map[string]*Counter {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	// Create a copy of the metrics
	copyMetrics := make(map[string]*Counter, len(ms.metrics))
	for k, v := range ms.metrics {
		copyMetrics[k] = &Counter{Value: v.Value}
	}

	// Sort the metrics by key
	sortedMetrics := make(map[string]*Counter)
	var keys []string
	for k := range copyMetrics {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fmt.Println(keys)
	for _, k := range keys {
		sortedMetrics[k] = copyMetrics[k]
	}

	fmt.Println(sortedMetrics)

	return sortedMetrics
}

// startAutoFlush automatically flushes metrics to BadgerDB at regular intervals
func (ms *MetricStore) startAutoFlush(flushIntervalSeconds int) {
	var interval = time.Duration(flushIntervalSeconds) * time.Second
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ticker.C:
			ms.flush()
			fmt.Println("Flushed metrics to DB")
		case <-ms.stopCh:
			ticker.Stop()
			return
		}
	}
}

// load reads saved metrics from BadgerDB and populates the in-memory store
func (ms *MetricStore) load() error {
	fmt.Println("Loading metrics from BadgerDB...")

	err := ms.db.View(func(txn *badger.Txn) error {
		itr := txn.NewIterator(badger.DefaultIteratorOptions)
		defer itr.Close()

		for itr.Rewind(); itr.Valid(); itr.Next() {
			item := itr.Item()
			name := string(item.Key())
			fmt.Printf("Loading metric: %s\n", name)

			// Retrieve the value (counter) from BadgerDB
			err := item.Value(func(val []byte) error {
				value, err := strconv.Atoi(string(val))
				if err != nil {
					return err
				}

				// Restore the counter in memory
				ms.metrics[name] = &Counter{Value: value}
				fmt.Printf("Loaded metric: %s with value: %d\n", name, value)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to load metrics from DB: %v", err)
	}

	return nil
}

// flush writes metrics from memory to BadgerDB
func (ms *MetricStore) flush() {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	fmt.Println("Flushing metrics to BadgerDB...")

	for name, counter := range ms.metrics {

		err := ms.db.Update(func(txn *badger.Txn) error {
			return txn.Set([]byte(name), []byte(strconv.Itoa(counter.Value)))
		})
		if err != nil {
			fmt.Printf("Failed to write metric %s to BadgerDB: %v", name, err)
		} else {
			fmt.Printf("Successfully flushed metric: %s\n", name)
		}
	}
}

// StopAutoFlush stops the background auto-flush mechanism
func (ms *MetricStore) StopAutoFlush() {
	ms.stopCh <- struct{}{}
}

// Close closes the BadgerDB connection when the app shuts down
func (ms *MetricStore) Close() {
	ms.flush()

	ms.db.Close()
}
