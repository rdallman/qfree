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

func (self *Queue) Enqueue(x string) {
	self.enqueue2(x)
}

func (self *Queue) Dequeue() (string, bool) {
	return self.dequeue2()
}

func (self *Queue) enqueue2(x string) {
	q := unsafe.Pointer(&Node{body: x, next: nil})
	p := atomic.LoadPointer(&self.tail)
	oldp := p
	for {
		for ((*Node)(p)).next != nil {
			p = ((*Node)(p).next)
		}
		if atomic.CompareAndSwapPointer(&(((*Node)(p)).next), nil, q) {
			break
		}
	}
	atomic.CompareAndSwapPointer(&(self.tail), oldp, q)
}

func (self *Queue) dequeue2() (string, bool) {
	for {
		p := self.head
		next := ((*Node)(p)).next
		if next == nil {
			return "", false
		}
		if atomic.CompareAndSwapPointer(&(self.head), p, next) {
			return ((*Node)(next)).body, true
		}
	}
}

func (self *Queue) enqueue1(x string) {
	newValue := unsafe.Pointer(&Node{body: x, next: nil})
	var tail, next unsafe.Pointer
	for {
		tail = self.tail
		next = ((*Node)(tail)).next
		if next != nil {
			atomic.CompareAndSwapPointer(&(self.tail), tail, next)
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

func (self *Queue) dequeue1() (string, bool) {
	var head, tail, next unsafe.Pointer
	for {
		head = self.head
		tail = self.tail
		next = ((*Node)(head)).next
		if head == tail {
			if next == nil {
				return "", false
			} else {
				atomic.CompareAndSwapPointer(&(self.tail), tail, next)
			}
		} else {
			val := ((*Node)(next)).body
			if atomic.CompareAndSwapPointer(&(self.head), head, next) {
				return val, true
			}
		}
		runtime.Gosched()
	}
}
