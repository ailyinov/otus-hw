package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity     int
	currCapacity int
	queue        List
	items        map[Key]*ListItem
	mu           sync.Mutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity:     capacity,
		currCapacity: capacity,
		queue:        NewList(),
		items:        make(map[Key]*ListItem, capacity),
	}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

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

func (l *lruCache) Get(key Key) (value interface{}, ok bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	i, ok := l.items[key]
	if ok {
		l.queue.MoveToFront(i)
		value = i.Value.(*cacheItem).value
	}
	return value, ok
}

func (l *lruCache) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
	l.currCapacity = l.capacity
}

func (l *lruCache) assureCapacity() {
	if l.currCapacity == 0 {
		delete(l.items, l.queue.Back().Value.(*cacheItem).key)
		l.queue.Remove(l.queue.Back())
		l.currCapacity++
	}
}

func (l *lruCache) addItem(key Key, ci *cacheItem) {
	l.assureCapacity()
	l.items[key] = l.queue.PushFront(ci)
	l.currCapacity--
}

func (l *lruCache) updateItem(key Key, ci *cacheItem) {
	l.items[key].Value = ci
	l.queue.MoveToFront(l.items[key])
}
