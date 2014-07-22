package qfree

import (
	"fmt"
	"testing"
)

func TestQueue(t *testing.T) {
	q := New()
	q.Enqueue("hi")
	fmt.Println(q.Dequeue())
	q.Enqueue("hi")
	q.Enqueue("hi")
	fmt.Println(q.Dequeue())
	fmt.Println(q.Dequeue())
	fmt.Println(q.Dequeue())
}

func BenchmarkQueue1Enq(b *testing.B) {
	q := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.enqueue1("hi")
	}
}

func BenchmarkQueue2Enq(b *testing.B) {
	q := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.enqueue2("hi")
	}
}

func BenchmarkQueue1Deq(b *testing.B) {
	q := New()
	for i := 0; i < 1000000; i++ {
		q.enqueue1("hi")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.dequeue1()
	}
}

func BenchmarkQueue2Deq(b *testing.B) {
	q := New()
	for i := 0; i < 1000000; i++ {
		q.enqueue2("hi")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.dequeue2()
	}
}
