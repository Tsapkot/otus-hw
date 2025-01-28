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

type lruCacheItem struct {
	key   Key
	value interface{}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	if l.capacity == 0 {
		return false
	}
	item, ok := l.items[key]
	if ok {
		item.Value.(*lruCacheItem).value = value
		l.queue.MoveToFront(item)
	} else {
		l.trimCache()
		cachedVal := l.queue.PushFront(&lruCacheItem{key, value})
		l.items[key] = cachedVal
	}
	return ok
}

func (l *lruCache) trimCache() {
	if l.capacity == l.queue.Len() {
		delete(l.items, l.queue.Back().Value.(*lruCacheItem).key)
		l.queue.Remove(l.queue.Back())
	}
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	if value, ok := l.items[key]; ok {
		l.queue.MoveToFront(value)
		return value.Value.(*lruCacheItem).value, true
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
