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

func (l list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v, Next: l.front, Prev: nil}
	if l.front == nil && l.back == nil {
		l.back = item
	} else {
		l.front.Prev = item
	}
	l.front = item
	l.len++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v, Next: nil, Prev: l.back}
	if l.back == nil && l.front == nil {
		l.front = item
	} else {
		l.back.Next = item
	}
	l.back = item
	l.len++
	return item
}

func (l *list) Remove(i *ListItem) {
	tempBack := l.back
	tempFront := l.front
	if i == tempBack {
		if i.Prev == nil {
			l.back = nil
		} else {
			i.Prev.Next = nil
			l.back = i.Prev
		}
	}

	if i == tempFront {
		if i.Next == nil {
			l.front = nil
		} else {
			i.Next.Prev = nil
			l.front = i.Next
		}
	}

	if i != tempBack && i != tempFront {
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.front == i {
		return
	}
	i.Prev.Next = i.Next
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}
	i.Next = l.front
	l.front.Prev = i
	i.Prev = nil
	l.front = i
}

func NewList() List {
	return new(list)
}
