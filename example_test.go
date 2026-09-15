package skiplist_test

import (
	"fmt"

	"github.com/seamonw/skiplist"
)

func Example() {
	l := skiplist.New[int, string]()
	l.Set(3, "c")
	l.Set(1, "a")
	l.Set(2, "b")
	l.Range(1, func(k int, v string) bool {
		fmt.Println(k, v)
		return true
	})
	// Output:
	// 1 a
	// 2 b
	// 3 c
}
