package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.Lock()
	defer c.Unlock()

	item, ok := c.items[key]
	if ok {
		item.Value = cacheItem{
			key:   key,
			value: value,
		}
		c.queue.MoveToFront(item)
		return true
	}

	c.queue.PushFront(cacheItem{
		key:   key,
		value: value,
	})
	c.items[key] = c.queue.Front()

	if c.queue.Len() > c.capacity {
		removingItem := c.queue.Back()
		c.queue.Remove(removingItem)
		delete(c.items, removingItem.Value.(cacheItem).key)
	}

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.Lock()
	defer c.Unlock()

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	c.queue.PushFront(item)
	return item.Value.(cacheItem).value, true
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
