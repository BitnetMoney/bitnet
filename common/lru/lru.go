// Copyright 2023 Bitnet
// This file is part of the Bitnet library.
//
// This software is provided "as is", without warranty of any kind,
// express or implied, including but not limited to the warranties
// of merchantability, fitness for a particular purpose and
// noninfringement. In no even shall the authors or copyright
// holders be liable for any claim, damages, or other liability,
// whether in an action of contract, tort or otherwise, arising
// from, out of or in connection with the software or the use or
// other dealings in the software.

package lru

import "sync"

// Cache is a LRU cache.
// This type is safe for concurrent use.
type Cache[K comparable, V any] struct {
	cache BasicLRU[K, V]
	mu    sync.Mutex
}

// NewCache creates an LRU cache.
func NewCache[K comparable, V any](capacity int) *Cache[K, V] {
	return &Cache[K, V]{cache: NewBasicLRU[K, V](capacity)}
}

// Add adds a value to the cache. Returns true if an item was evicted to store the new item.
func (c *Cache[K, V]) Add(key K, value V) (evicted bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cache.Add(key, value)
}

// Contains reports whether the given key exists in the cache.
func (c *Cache[K, V]) Contains(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cache.Contains(key)
}

// Get retrieves a value from the cache. This marks the key as recently used.
func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cache.Get(key)
}

// Len returns the current number of items in the cache.
func (c *Cache[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cache.Len()
}

// Peek retrieves a value from the cache, but does not mark the key as recently used.
func (c *Cache[K, V]) Peek(key K) (value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cache.Peek(key)
}

// Purge empties the cache.
func (c *Cache[K, V]) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache.Purge()
}

// Remove drops an item from the cache. Returns true if the key was present in cache.
func (c *Cache[K, V]) Remove(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cache.Remove(key)
}

// Keys returns all keys of items currently in the LRU.
func (c *Cache[K, V]) Keys() []K {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cache.Keys()
}
