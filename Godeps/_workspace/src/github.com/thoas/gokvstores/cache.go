package gokvstores

import (
	"container/list"
	"fmt"
	"sync"
)

type CacheKVStore struct {
	Cache *CacheKVStoreConnection
}

func NewCacheKVStore(maxEntries int) KVStore {
	return &CacheKVStore{Cache: NewCacheKVStoreConnection(maxEntries)}
}

// Cache is an LRU cache, safe for concurrent access.
type CacheKVStoreConnection struct {
	maxEntries int

	mu    sync.Mutex
	ll    *list.List
	cache map[string]*list.Element
}

// *entry is the type stored in each *list.Element.
type entry struct {
	key   string
	value interface{}
}

// New returns a new cache with the provided maximum items.
func NewCacheKVStoreConnection(maxEntries int) *CacheKVStoreConnection {
	return &CacheKVStoreConnection{
		maxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[string]*list.Element),
	}
}

func (c *CacheKVStore) Connection() KVStoreConnection {
	return c.Cache
}

func (c *CacheKVStore) Close() error {
	return nil
}

func (c *CacheKVStoreConnection) Flush() error {
	c.ll = list.New()
	c.cache = make(map[string]*list.Element)

	return nil
}

func (c *CacheKVStoreConnection) Close() error {
	return nil
}

func (c *CacheKVStoreConnection) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.cache[key]; !ok {
		return fmt.Errorf("Key %s does not exist", key)
	}

	delete(c.cache, key)

	return nil
}

func (c *CacheKVStoreConnection) Exists(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.cache[key]

	return ok
}

// Add adds the provided key and value to the cache, evicting
// an old item if necessary.
func (c *CacheKVStoreConnection) Set(key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Already in cache?
	if ee, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ee)
		ee.Value.(*entry).value = value
		return nil
	}

	// Add to cache if not present
	ele := c.ll.PushFront(&entry{key, value})
	c.cache[key] = ele

	if c.ll.Len() > c.maxEntries && c.maxEntries != -1 {
		c.removeOldest()
	}

	return nil
}

// Get fetches the key's value from the cache.
// The ok result will be true if the item was found.
func (c *CacheKVStoreConnection) Get(key string) interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ele, hit := c.cache[key]; hit {
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).value
	}

	return nil
}

// RemoveOldest removes the oldest item in the cache and returns its key and value.
// If the cache is empty, the empty string and nil are returned.
func (c *CacheKVStoreConnection) RemoveOldest() (key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.removeOldest()
}

// note: must hold c.mu
func (c *CacheKVStoreConnection) removeOldest() (key string, value interface{}) {
	ele := c.ll.Back()
	if ele == nil {
		return
	}
	c.ll.Remove(ele)
	ent := ele.Value.(*entry)
	delete(c.cache, ent.key)
	return ent.key, ent.value

}

// Len returns the number of items in the cache.
func (c *CacheKVStoreConnection) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ll.Len()
}
