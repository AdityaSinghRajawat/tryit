package utils

import (
	"container/list"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
)

// Cache is an LRU with optional disk-backed persistence (IMPL §9.5).
// Construct via NewCache; inject into services that need it.
type Cache struct {
	mu      sync.Mutex
	ll      *list.List
	items   map[string]*list.Element
	cap     int
	ttl     time.Duration
	diskDir string
	enabled bool
}

type cacheEntry struct {
	key       string
	value     []byte
	expiresAt time.Time
}

func NewCache() *Cache {
	c := &Cache{
		ll:      list.New(),
		items:   make(map[string]*list.Element),
		cap:     config.GetCacheLRUCapacity(),
		ttl:     config.GetCacheTTL(),
		diskDir: config.GetCacheDiskDir(),
		enabled: config.GetCacheEnabled(),
	}
	if c.enabled && c.diskDir != "" {
		_ = os.MkdirAll(c.diskDir, 0o700)
	}
	return c
}

// Get returns the cached payload, or (nil, false) on miss / expiry / disabled.
// In-memory misses fall through to disk when a disk dir is configured.
func (c *Cache) Get(key string) ([]byte, bool) {
	if !c.enabled {
		return nil, false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		e := elem.Value.(*cacheEntry)
		if time.Now().After(e.expiresAt) {
			c.removeElement(elem)
			return nil, false
		}
		c.ll.MoveToFront(elem)
		return clone(e.value), true
	}

	if c.diskDir == "" {
		return nil, false
	}

	b, exp, ok := c.readDisk(key)
	if !ok {
		return nil, false
	}

	if time.Now().After(exp) {
		c.deleteDisk(key)
		return nil, false
	}

	c.setLocked(key, b, exp)
	return clone(b), true
}

// Set stores value at key with the configured TTL, evicting LRU when over cap.
func (c *Cache) Set(key string, value []byte) {
	if !c.enabled {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	expiresAt := time.Now().Add(c.ttl)
	c.setLocked(key, value, expiresAt)

	if c.diskDir != "" {
		c.writeDisk(key, value, expiresAt)
	}
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, ok := c.items[key]; ok {
		c.removeElement(elem)
	}

	if c.diskDir != "" {
		c.deleteDisk(key)
	}
}

func (c *Cache) setLocked(key string, value []byte, expiresAt time.Time) {
	if elem, ok := c.items[key]; ok {
		e := elem.Value.(*cacheEntry)
		e.value = append(e.value[:0], value...)
		e.expiresAt = expiresAt
		c.ll.MoveToFront(elem)
		return
	}
	stored := clone(value)
	elem := c.ll.PushFront(&cacheEntry{key: key, value: stored, expiresAt: expiresAt})
	c.items[key] = elem
	for c.ll.Len() > c.cap {
		c.removeElement(c.ll.Back())
	}
}

func (c *Cache) removeElement(elem *list.Element) {
	if elem == nil {
		return
	}
	c.ll.Remove(elem)
	delete(c.items, elem.Value.(*cacheEntry).key)
}

type cacheDiskEnvelope struct {
	ExpiresAt time.Time `json:"expiresAt"`
	Value     []byte    `json:"value"`
}

func (c *Cache) diskPath(key string) string {
	return filepath.Join(c.diskDir, key+".json")
}

func (c *Cache) writeDisk(key string, value []byte, expiresAt time.Time) {
	b, err := json.Marshal(cacheDiskEnvelope{ExpiresAt: expiresAt, Value: value})
	if err != nil {
		return
	}
	_ = os.WriteFile(c.diskPath(key), b, 0o600)
}

func (c *Cache) readDisk(key string) ([]byte, time.Time, bool) {
	b, err := os.ReadFile(c.diskPath(key))
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, time.Time{}, false
		}
		return nil, time.Time{}, false
	}
	var env cacheDiskEnvelope
	if jerr := json.Unmarshal(b, &env); jerr != nil {
		return nil, time.Time{}, false
	}
	return env.Value, env.ExpiresAt, true
}

func (c *Cache) deleteDisk(key string) {
	_ = os.Remove(c.diskPath(key))
}

func clone(b []byte) []byte {
	out := make([]byte, len(b))
	copy(out, b)
	return out
}
