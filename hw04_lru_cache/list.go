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
	first  *ListItem
	last   *ListItem
	length int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v}

	if l.first == nil {
		l.first = item
		l.last = item
		item.Next = l.last
		item.Prev = l.first
	} else {
		item.Next = l.first
		l.first.Prev = item
		l.first = item
	}

	l.length++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}

	if l.last == nil {
		l.first = item
		l.last = item
		item.Next = l.last
		item.Prev = l.first
	} else {
		l.last.Next = item
		item.Prev = l.last
		l.last = item
	}

	l.length++
	return item
}

func (l *list) Remove(i *ListItem) {
	if i.Prev == nil {
		l.first.Next.Prev = i
		l.first = i.Next
	} else {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	if i.Next == nil {
		l.last.Prev.Next = nil
		l.last = i.Prev
	} else {
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}

	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.first == i {
		return
	}

	if l.last == i {
		l.last.Prev.Next = nil
		i.Next = l.first
		l.first.Prev = i
		l.first = i
		l.first.Prev = nil

		return
	}

	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev

	i.Next = l.first
	l.first.Prev = i
	l.first = i
	l.first.Prev = nil
}
