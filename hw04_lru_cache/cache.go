package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	_, exists := c.items[key]
	if exists {
		c.items[key].Value = cacheItem{value: value, key: key}
		c.queue.MoveToFront(c.items[key])
	} else {
		if c.queue.Len() == c.capacity {
			value := c.queue.Back()
			cacheItem, ok := value.Value.(cacheItem)

			if ok {
				delete(c.items, cacheItem.key)
			}
			c.queue.Remove(value)
		}
		c.items[key] = c.queue.PushFront(cacheItem{value: value, key: key})
	}
	return exists
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	value, exists := c.items[key]

	if exists {
		cacheItem, ok := value.Value.(cacheItem)

		if ok {
			return cacheItem.value, true
		}
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue = NewList()
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
