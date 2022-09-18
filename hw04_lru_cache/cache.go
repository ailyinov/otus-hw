package hw04lrucache

import "sync"

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
	mu       sync.RWMutex
}

type cacheItem struct {
	key   Key
	value any
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
		l.updateItem(key, ci)
	} else {
		l.addItem(key, ci)
	}
	return ok
}

func (l *lruCache) Get(key Key) (item any, ok bool) {
	l.mu.RLock()
	i, ok := l.items[key]
	if ok {
		l.queue.MoveToFront(i)
		item = i.Value.(*cacheItem).value
	}
	l.mu.RUnlock()
	return item, ok
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
	l.mu.Lock()
	l.assureCapacity()
	l.items[key] = l.queue.PushFront(ci)
	l.capacity--
	l.mu.Unlock()
}

func (l *lruCache) updateItem(key Key, ci *cacheItem) {
	l.mu.Lock()
	l.items[key].Value = ci
	l.queue.MoveToFront(l.items[key])
	l.mu.Unlock()
}
