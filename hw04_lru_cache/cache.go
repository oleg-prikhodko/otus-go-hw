package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	if l.capacity == 0 {
		return false
	}

	if i, ok := l.items[key]; ok {
		l.items[key].Value = cacheItem{key, value}
		l.queue.MoveToFront(i)

		return true
	}

	l.items[key] = l.queue.PushFront(cacheItem{key, value})
	if l.queue.Len() > l.capacity {
		last := l.queue.Back()
		l.queue.Remove(last)
		delete(l.items, last.Value.(cacheItem).key)
	}

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	if i, ok := l.items[key]; ok {
		l.queue.MoveToFront(i)
		return i.Value.(cacheItem).value, true
	}

	return nil, false
}

func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
