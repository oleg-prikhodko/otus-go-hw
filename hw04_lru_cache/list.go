package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
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

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Next:  l.front,
		Prev:  nil,
	}
	if l.front == nil {
		l.back = item
	} else {
		l.front.Prev = item
		l.front = item
	}
	l.front = item
	l.len++

	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.back,
	}
	if l.back == nil {
		l.front = item
	} else {
		l.back.Next = item
		l.back = item
	}
	l.back = item
	l.len++

	return item
}

func (l *list) Remove(i *ListItem) {
	prev := i.Prev
	next := i.Next
	l.len--

	if prev == nil {
		l.front = next
	} else {
		prev.Next = next
	}

	if next == nil {
		l.back = prev
	} else {
		next.Prev = prev
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.front {
		return
	}

	next := i.Next
	prev := i.Prev

	prev.Next = next
	if next != nil {
		next.Prev = prev
	} else {
		l.back = prev
	}

	i.Prev = nil
	i.Next = l.front

	l.front.Prev = i
	l.front = i
}

func NewList() List {
	return new(list)
}
