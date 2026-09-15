// Package skiplist is a sequential skip list (p = 1/4, maxLevel = 32).
// It is not safe for concurrent use without an external lock.
//
// See https://seamonw.github.io/blog/2026/09/15/golang-skiplist/
package skiplist

import (
	"cmp"
	"iter"
	"math/rand/v2"
)

const (
	maxLevel = 32
	pThresh  = uint32(1) << 30 // p = 1/4
)

type node[K cmp.Ordered, V any] struct {
	key  K
	val  V
	next []*node[K, V]
}

type List[K cmp.Ordered, V any] struct {
	head  *node[K, V]
	level int
	n     int
	rnd   func() uint32
}

func New[K cmp.Ordered, V any]() *List[K, V] {
	return &List[K, V]{
		head:  &node[K, V]{next: make([]*node[K, V], maxLevel)},
		level: 1,
		rnd:   rand.Uint32,
	}
}

func (l *List[K, V]) Len() int { return l.n }

func (l *List[K, V]) randomLevel() int {
	lv := 1
	for lv < maxLevel && l.rnd() < pThresh {
		lv++
	}
	return lv
}

func (l *List[K, V]) search(target K, prev *[maxLevel]*node[K, V]) {
	x := l.head
	for i := l.level - 1; i >= 0; i-- {
		for nxt := x.next[i]; nxt != nil && nxt.key < target; nxt = x.next[i] {
			x = nxt
		}
		prev[i] = x
	}
}

func (l *List[K, V]) Get(key K) (V, bool) {
	x := l.head
	for i := l.level - 1; i >= 0; i-- {
		for nxt := x.next[i]; nxt != nil && nxt.key < key; nxt = x.next[i] {
			x = nxt
		}
	}
	if nxt := x.next[0]; nxt != nil && nxt.key == key {
		return nxt.val, true
	}
	var zero V
	return zero, false
}

func (l *List[K, V]) Set(key K, val V) {
	var prev [maxLevel]*node[K, V]
	l.search(key, &prev)

	if nxt := prev[0].next[0]; nxt != nil && nxt.key == key {
		nxt.val = val
		return
	}

	lv := l.randomLevel()
	if lv > l.level {
		for i := l.level; i < lv; i++ {
			prev[i] = l.head
		}
		l.level = lv
	}

	x := &node[K, V]{key: key, val: val, next: make([]*node[K, V], lv)}
	for i := 0; i < lv; i++ {
		x.next[i] = prev[i].next[i]
		prev[i].next[i] = x
	}
	l.n++
}

func (l *List[K, V]) Delete(key K) bool {
	var prev [maxLevel]*node[K, V]
	l.search(key, &prev)

	x := prev[0].next[0]
	if x == nil || x.key != key {
		return false
	}
	for i := 0; i < len(x.next); i++ {
		prev[i].next[i] = x.next[i]
	}
	for l.level > 1 && l.head.next[l.level-1] == nil {
		l.level--
	}
	l.n--
	return true
}

func (l *List[K, V]) Range(start K, f func(K, V) bool) {
	x := l.head
	for i := l.level - 1; i >= 0; i-- {
		for nxt := x.next[i]; nxt != nil && nxt.key < start; nxt = x.next[i] {
			x = nxt
		}
	}
	for cur := x.next[0]; cur != nil; cur = cur.next[0] {
		if !f(cur.key, cur.val) {
			return
		}
	}
}

func (l *List[K, V]) From(start K) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) { l.Range(start, yield) }
}

func (l *List[K, V]) LevelHist() []int {
	h := make([]int, l.level)
	for i := range h {
		for cur := l.head.next[i]; cur != nil; cur = cur.next[i] {
			h[i]++
		}
	}
	return h
}
