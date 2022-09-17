package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value any) bool
	Get(key Key) (any, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
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

func (l *lruCache) Set(key Key, value any) bool {
	ci := &cacheItem{
		key:   key,
		value: value,
	}
	_, ok := l.items[key]
	if ok {
		l.items[key].Value = ci
		l.queue.MoveToFront(l.items[key])
	} else {
		l.assureCapacity()
		l.items[key] = l.queue.PushFront(ci)
		l.capacity--
	}
	return ok
}

func (l *lruCache) Get(key Key) (any, bool) {
	if i, ok := l.items[key]; ok {
		l.queue.MoveToFront(i)
		return i.Value.(*cacheItem).value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

func (l *lruCache) assureCapacity() {
	if l.capacity == 0 {
		delete(l.items, l.queue.Back().Value.(*cacheItem).key)
		l.queue.Remove(l.queue.Back())
		l.capacity++
	}
}
