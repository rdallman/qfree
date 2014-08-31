package qfree

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

type Queue struct {
	head, tail unsafe.Pointer
}

type Node struct {
	body string
	next unsafe.Pointer
}

func New() *Queue {
	q := new(Queue)
	q.head = unsafe.Pointer(new(Node))
	q.tail = q.head
	return q
}

func (q *Queue) Enqueue(x string) {
	q.enqueue2(x)
}

func (q *Queue) Dequeue() (string, bool) {
	return q.dequeue2()
}

func (q *Queue) enqueue2(x string) {
	node := unsafe.Pointer(&Node{body: x, next: nil})
	p := atomic.LoadPointer(&q.tail)
	oldp := p
	for {
		for ((*Node)(p)).next != nil {
			p = ((*Node)(p).next)
		}
		if atomic.CompareAndSwapPointer(&(((*Node)(p)).next), nil, node) {
			break
		}
	}
	atomic.CompareAndSwapPointer(&(q.tail), oldp, node)
}

func (q *Queue) dequeue2() (string, bool) {
	for {
		p := q.head
		next := ((*Node)(p)).next
		if next == nil {
			return "", false
		}
		if atomic.CompareAndSwapPointer(&(q.head), p, next) {
			return ((*Node)(next)).body, true
		}
	}
}

func (q *Queue) enqueue1(x string) {
	newValue := unsafe.Pointer(&Node{body: x, next: nil})
	var tail, next unsafe.Pointer
	for {
		tail = q.tail
		next = ((*Node)(tail)).next
		if next != nil {
			atomic.CompareAndSwapPointer(&(q.tail), tail, next)
		} else if atomic.CompareAndSwapPointer(&((*Node)(tail).next), nil, newValue) {
			break
		}
		runtime.Gosched()
	}
}

type ErrQueueEmpty struct{}

func (e ErrQueueEmpty) Error() string {
	return "queue is empty"
}

func (q *Queue) dequeue1() (string, bool) {
	var head, tail, next unsafe.Pointer
	for {
		head = q.head
		tail = q.tail
		next = ((*Node)(head)).next
		if head == tail {
			if next == nil {
				return "", false
			} else {
				atomic.CompareAndSwapPointer(&(q.tail), tail, next)
			}
		} else {
			val := ((*Node)(next)).body
			if atomic.CompareAndSwapPointer(&(q.head), head, next) {
				return val, true
			}
		}
		runtime.Gosched()
	}
}
