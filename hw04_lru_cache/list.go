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
	itemsCount  int
	front, back *ListItem
}

func NewList() List {
	return new(list)
}

func (l list) Len() int {
	return l.itemsCount
}

func (l list) Front() *ListItem {
	return l.front
}

func (l list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	newListItem := ListItem{Value: v, Prev: nil, Next: l.front}

	// обновление связей первого элемента
	if l.front != nil {
		l.front.Prev = &newListItem
		newListItem.Next = l.front
	}

	// обновление последнего элемента, если это 1й итем
	if l.itemsCount == 0 {
		l.back = &newListItem
	}

	l.itemsCount++
	l.front = &newListItem

	return &newListItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newListItem := ListItem{Value: v, Prev: l.back, Next: nil}

	// обновление связей последнегр элемента
	if l.back != nil {
		l.back.Next = &newListItem
		newListItem.Prev = l.back
	}

	// обновление первого элемента, если это 1й итем
	if l.itemsCount == 0 {
		l.front = &newListItem
	}

	l.itemsCount++
	l.back = &newListItem

	return &newListItem
}

func (l *list) Remove(i *ListItem) {
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	l.itemsCount--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.itemsCount < 2 || l.front == i {
		return
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	l.front.Prev = i
	i.Next = l.front
	l.front = i
}
