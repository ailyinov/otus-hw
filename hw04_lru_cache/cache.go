package hw04lrucache

import "sync"

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
	mu       sync.Mutex
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

func (l *lruCache) Set(key Key, value interface{}) bool {
	ci := &cacheItem{
		key:   key,
		value: value,
	}

	l.mu.Lock()
	_, ok := l.items[key]
	if ok {
		l.updateItem(key, ci)
	} else {
		l.addItem(key, ci)
	}
	l.mu.Unlock()
	return ok
}

func (l *lruCache) Get(key Key) (value interface{}, ok bool) {
	l.mu.Lock()
	i, ok := l.items[key]
	if ok {
		l.queue.MoveToFront(i)
		value = i.Value.(*cacheItem).value
	}
	l.mu.Unlock()
	return value, ok
}

func (l *lruCache) Clear() {
	l.mu.Lock()
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
	l.mu.Unlock()
}

func (l *lruCache) assureCapacity() {
	if l.capacity == 0 {
		delete(l.items, l.queue.Back().Value.(*cacheItem).key)
		l.queue.Remove(l.queue.Back())
		l.capacity++
	}
}

func (l *lruCache) addItem(key Key, ci *cacheItem) {
	l.assureCapacity()
	l.items[key] = l.queue.PushFront(ci)
	l.capacity--
}

func (l *lruCache) updateItem(key Key, ci *cacheItem) {
	l.items[key].Value = ci
	l.queue.MoveToFront(l.items[key])
}
