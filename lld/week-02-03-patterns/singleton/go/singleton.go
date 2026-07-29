// Package singleton implements a thread-safe AppConfig singleton using
// sync.Once so first-access under concurrent goroutines only ever
// constructs one instance.
package singleton

import "sync"

type AppConfig struct {
	id       int
	settings map[string]string
	mu       sync.RWMutex
}

var (
	instance     *AppConfig
	once         sync.Once
	instanceSeq  int
	instanceSeqM sync.Mutex
)

// GetInstance returns the single shared AppConfig, constructing it exactly
// once regardless of how many goroutines call this concurrently.
func GetInstance() *AppConfig {
	once.Do(func() {
		instanceSeqM.Lock()
		instanceSeq++
		id := instanceSeq
		instanceSeqM.Unlock()

		instance = &AppConfig{
			id:       id,
			settings: make(map[string]string),
		}
	})
	return instance
}

// ID uniquely identifies which construction produced this instance; used
// in tests to prove concurrent callers all observe the same instance.
func (c *AppConfig) ID() int {
	return c.id
}

func (c *AppConfig) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.settings[key] = value
}

func (c *AppConfig) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.settings[key]
	return v, ok
}

// resetForTest clears the singleton so each test starts from a clean state.
// Only exported test code in this package should call it.
func resetForTest() {
	instance = nil
	once = sync.Once{}
}
