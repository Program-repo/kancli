package llist

import (
	"fmt"
)

type node[T any] struct {
	value      T
	next, prev *node[T]
}

func (n node[T]) Value() T {
	return n.value
}

func (n node[T]) Next() *node[T] {
	return n.next
}

func (n node[T]) Prev() *node[T] {
	return n.prev
}

type LinkedList[T any] struct {
	head, tail *node[T]
	length     int
}

func New[T any]() *LinkedList[T] {
	return &LinkedList[T]{}
}

func (l *LinkedList[T]) PushFront(value T) {

	newNode := &node[T]{value: value}
	if l.head == nil {
		l.head = newNode
		l.tail = newNode
	} else {
		oldNode := l.head
		oldNode.prev = newNode
		newNode.next = oldNode
		l.head = newNode
	}
	l.length++
}

func (l *LinkedList[T]) PushBack(value T) {

	newNode := &node[T]{value: value}
	if l.tail == nil {
		l.head = newNode
		l.tail = newNode
	} else {
		oldNode := l.tail
		newNode.prev = oldNode
		oldNode.next = newNode
		l.tail = newNode
	}
	l.length++
}

func (l *LinkedList[T]) DeleteAt(index int) error {
	// size := l.length
	// special case: input index is out of range, return error
	if index < 0 || index > l.length-1 { //size-1 {
		return fmt.Errorf("index out of range")
	}
	// special case: index is 0, delete at head by use DeleteAtHead method
	if index == 0 {
		if l.head == nil {
			return fmt.Errorf("list is empty")
		}
		l.head = l.head.next
		// node := l.head
		// node.prev = nil
		l.head.prev = nil
		l.length--
		return nil
	}

	// general case: find the node before the index, then delete its' next node
	current := l.head
	for i := 0; i < index-1; i++ {
		current = current.next
	}

	// node := current.next.next
	if current.next.next == nil {
		current.next = nil
		l.tail = current
	} else {
		current.next.next.prev = current.next.prev
		// node.prev = current.next.prev
		current.next = current.next.next
	}

	l.length--
	return nil

}

func (l *LinkedList[T]) Head() *node[T] {
	return l.head
}

func (l *LinkedList[T]) Tail() *node[T] {
	return l.tail
}
func (l *LinkedList[T]) Length() int {
	return l.length
}
