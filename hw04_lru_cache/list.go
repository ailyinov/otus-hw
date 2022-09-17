package hw04lrucache

import "fmt"

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v any) *ListItem
	PushBack(v any) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value any
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	front *ListItem
	back  *ListItem
	len   int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v any) *ListItem {
	if l.len == 0 {
		return l.pushFirst(v)
	}
	item := ListItem{
		Value: v,
		Next:  l.front,
		Prev:  nil,
	}
	l.front.Prev = &item
	l.front = l.front.Prev
	l.len++
	return l.front
}

func (l *list) PushBack(v any) *ListItem {
	if l.len == 0 {
		return l.pushFirst(v)
	}
	item := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.back,
	}
	l.back.Next = item
	l.back = l.back.Next
	l.len++
	return l.back
}

func (l *list) Remove(i *ListItem) {
	if i, err := l.findItem(i); err == nil {
		if l.Back() == i {
			l.back = i.Prev
		} else {
			i.Next.Prev = i.Prev
		}
		if l.Front() == i {
			l.front = i.Next
		} else {
			i.Prev.Next = i.Next
		}
		l.len--
	}
}

func (l *list) findItem(i *ListItem) (*ListItem, error) {
	for curr := l.front; curr != nil; curr = curr.Next {
		if i == curr {
			return curr, nil
		}
	}
	return nil, fmt.Errorf("item not found: %s", i.Value)
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	i.Next = l.front
	i.Prev = nil
	l.front.Prev = i
	l.front = i
	l.len++
}

func (l *list) pushFirst(v any) *ListItem {
	item := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}
	l.back = item
	l.front = item
	l.len++
	return item
}
