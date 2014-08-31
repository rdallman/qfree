package qfree

import "testing"

func TestQueue(t *testing.T) {
	q := New()
	q.Enqueue("hi1")
	q.Enqueue("hi2")
	q.Enqueue("hi3")
	d, ok := q.Dequeue()
	if d != "hi1" || !ok {
		t.Error("expected 'hi1' but got", d, ok)
	}
	d, ok = q.Dequeue()
	if d != "hi2" || !ok {
		t.Error("expected 'hi1' but got", d, ok)
	}
	d, ok = q.Dequeue()
	if d != "hi3" || !ok {
		t.Error("expected 'hi1' but got", d, ok)
	}

	d, ok = q.Dequeue()
	if ok || d != "" {
		t.Error("expected no item to dequeue, but got", d, ok)
	}
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
