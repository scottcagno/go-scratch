package queue

import (
	"container/list"
)

type Item[T any] struct {
	*list.Element
}

func (it *Item[T]) Next() T {
	return it.Element.Next().Value.(T)
}

func (it *Item[T]) Prev() T {
	return it.Element.Prev().Value.(T)
}

// Queue is a generic queue type that uses a doubly linked list
type Queue[T any] struct {
	list *list.List
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		list: list.New(),
	}
}

// Back returns the last item of list q or nil if the list is empty.
func (q *Queue[T]) Back() T {
	return q.list.Back().Value.(T)
}

// Front returns the first item of list q or nil if the list is empty.
func (q *Queue[T]) Front() T {
	return q.list.Front().Value.(T)
}

// InsertAfter inserts a new item it with value v immediately after mark
// and returns it. If mark is not an item of q, the list is not modified.
// The mark must not be nil.
func (q *Queue[T]) InsertAfter(v T, it *Item[T]) *Item[T] {
	return q.list.InsertAfter(v, it.Element).Value.(*Item[T])
}

// InsertBefore inserts a new item it with value v immediately before mark
// and returns it. If mark is not an item of q, the list is not modified.
// The mark must not be nil.
func (q *Queue[T]) InsertBefore(v T, it *Item[T]) *Item[T] {
	return q.list.InsertBefore(v, it.Element).Value.(*Item[T])
}

// Len returns the length which is the total number of items in the queue
func (q *Queue[T]) Len() int {
	return q.list.Len()
}

// MoveAfter moves item it to its new position after mark. If e or mark
// is not an element of q, or it == mark, the list is not modified. The
// item and mark must not be nil.
func (q *Queue[T]) MoveAfter(it, mark *Item[T]) {
	q.list.MoveAfter(it.Element, mark.Element)
}

// MoveBefore moves item it to its new position before mark. If e or mark
// is not an element of q, or it == mark, the list is not modified. The
// item and mark must not be nil.
func (q *Queue[T]) MoveBefore(it, mark *Item[T]) {
	q.list.MoveBefore(it.Element, mark.Element)
}

// MoveToBack moves item it to the back of list q. If it is not an item of
// q, the list is not modified. The item must not be nil.
func (q *Queue[T]) MoveToBack(it *Item[T]) {
	q.list.MoveToBack(it.Element)
}

// MoveToFront moves item it to the front of list q. If it is not an item of
// q, the list is not modified. The item must not be nil.
func (q *Queue[T]) MoveToFront(it *Item[T]) {
	q.list.MoveToBack(it.Element)
}

// PushBack inserts a new item it with value v at the back of the list q and
// returns item it.
func (q *Queue[T]) PushBack(it *Item[T]) *Item[T] {
	return q.list.PushBack(it.Element).Value.(*Item[T])
}

// PushFront inserts a new item it with value v at the front of the list q and
// returns item it.
func (q *Queue[T]) PushFront(it *Item[T]) *Item[T] {
	return q.list.PushFront(it.Element).Value.(*Item[T])
}

// Remove removes item it from list q if it is an item of list q. It returns
// the item value. The item must not be nil.
func (q *Queue[T]) Remove(it *Item[T]) T {
	return q.list.Remove(it.Element).(T)
}
