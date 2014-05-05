package yak

import (
	"testing"
)

func TestNewQueue(t *testing.T) {
	q := NewQueue()
	if q.Empty() != true {
		t.Fatal("New queue length must be 0")
	}
}

func TestEnqueue(t *testing.T) {
	q1 := NewQueue()
	q2 := NewQueue()
	q1.Enqueue("foo")
	if q1.Len() != 1 {
		t.Fatal("Expected 1, got", q1.Len())
	}
	if q2.Empty() != true {
		t.Fatal("Expected true, got", q2.Empty())
	}
	q1.Enqueue("bar")
	if q1.Len() != 2 {
		t.Fatal("Expected 2, got", q1.Len())
	}
	if q2.Empty() != true {
		t.Fatal("Expected true, got", q2.Empty())
	}
	q2.Enqueue("hoge")
	if q1.Len() != 2 {
		t.Fatal("Expected 2, got", q1.Len())
	}
	if q2.Len() != 1 {
		t.Fatal("Expected 1, got", q2.Len())
	}
	if q2.Empty() != false {
		t.Fatal("Expected false, got", q2.Empty())
	}
}

func TestDequeue(t *testing.T) {
	q1 := NewQueue()
	str := []string{"foo", "bar", "qux"}
	for _, i := range str {
		q1.Enqueue(i)
	}
	var v interface{}
	for j, _ := range str {
		v = q1.Dequeue()
		if v != str[j] {
			t.Fatalf("expected %v, got %v", str[j], v)
		}
	}
	v = q1.Dequeue()
	if v != nil {
		t.Fatal("expected nil, got %v", v)
	}
	if q1.Empty() != true {
		t.Fatal("expected true, got", q1.Empty())
	}
}
