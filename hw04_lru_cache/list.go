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
		if i.Next != nil {
			i.Next.Prev = i.Prev
		} else {
			l.back = i.Prev
		}
		if i.Prev != nil {
			i.Prev.Next = i.Next
		} else {
			l.front = i.Next
		}
		l.len--
		return
	}
}

func (l *list) findItem(i *ListItem) (*ListItem, error) {
	curr := l.front
	for curr != nil {
		if i == curr {
			return curr, nil
		}
		curr = curr.Next
	}
	return nil, fmt.Errorf("item not found: %s", i.Value)
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	i.Next = l.front
	i.Prev = nil
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
