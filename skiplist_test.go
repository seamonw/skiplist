package skiplist

import (
	"math"
	"slices"
	"testing"
)

func TestEmpty(t *testing.T) {
	l := New[int, string]()
	if l.Len() != 0 {
		t.Fatal(l.Len())
	}
	if _, ok := l.Get(1); ok {
		t.Fatal("empty get")
	}
	if l.Delete(1) {
		t.Fatal("empty delete")
	}
}

func TestSetGetUpsert(t *testing.T) {
	l := New[int, string]()
	l.Set(1, "a")
	l.Set(2, "b")
	l.Set(1, "A")
	if l.Len() != 2 {
		t.Fatalf("len %d", l.Len())
	}
	v, ok := l.Get(1)
	if !ok || v != "A" {
		t.Fatalf("got %q %v", v, ok)
	}
}

func TestDeleteAndLevel(t *testing.T) {
	l := New[int, int]()
	for i := 0; i < 100; i++ {
		l.Set(i, i)
	}
	if !l.Delete(42) {
		t.Fatal("missing 42")
	}
	if l.Delete(42) {
		t.Fatal("double delete")
	}
	if _, ok := l.Get(42); ok {
		t.Fatal("42 still there")
	}
	for i := 0; i < 100; i++ {
		l.Delete(i)
	}
	if l.Len() != 0 {
		t.Fatal(l.Len())
	}
	if l.level != 1 {
		t.Fatalf("level %d want 1", l.level)
	}
}

func TestRangeOrder(t *testing.T) {
	l := New[int, int]()
	for _, k := range []int{5, 1, 9, 3, 7} {
		l.Set(k, k*10)
	}
	var keys []int
	l.Range(math.MinInt, func(k, v int) bool {
		if v != k*10 {
			t.Fatalf("%d -> %d", k, v)
		}
		keys = append(keys, k)
		return true
	})
	if !slices.Equal(keys, []int{1, 3, 5, 7, 9}) {
		t.Fatalf("%v", keys)
	}

	keys = keys[:0]
	l.Range(5, func(k, _ int) bool {
		keys = append(keys, k)
		return true
	})
	if !slices.Equal(keys, []int{5, 7, 9}) {
		t.Fatalf("from 5: %v", keys)
	}
}

func TestFromBreak(t *testing.T) {
	l := New[int, int]()
	for i := 0; i < 10; i++ {
		l.Set(i, i)
	}
	var n int
	for range l.From(0) {
		n++
		if n == 3 {
			break
		}
	}
	if n != 3 {
		t.Fatal(n)
	}
}

func TestLevelHistShape(t *testing.T) {
	l := New[int, struct{}]()
	const n = 10000
	for i := 0; i < n; i++ {
		l.Set(i, struct{}{})
	}
	h := l.LevelHist()
	if h[0] != n {
		t.Fatalf("L0 %d", h[0])
	}
	if len(h) > 1 {
		ratio := float64(h[1]) / float64(n)
		if ratio < 0.15 || ratio > 0.40 {
			t.Fatalf("L1 ratio %v hist %v", ratio, h)
		}
	}
}
