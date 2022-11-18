package syncslice_test

import (
	"fmt"

	"github.com/nbd-wtf/syncslice"
)

func Example() {
	s := syncslice.Make[string](2, 3)
	fmt.Println(s.Len())

	s.Set(0, "hello")
	s.Set(1, "world")
	fmt.Println(s.Len())

	s.Append("!")
	fmt.Println(s.Len())

	for item := range s.Iter() {
		fmt.Println(item.Index, item.Value)
	}

	next := s.Slice(0, 1)
	next.Append("cruel")
	next.Append(s.Get(1))
	next.Append("bogus")
	next.Append("bogus")
	next.Append("bogus")
	next.Append("bogus")
	next.Append("bogus")

	next.Range(func(index int, value string) bool {
		if value == "bogus" {
			return false
		}
		fmt.Println(index, value)
		return true
	})

	// Output:
	// 2
	// 2
	// 3
	// 0 hello
	// 1 world
	// 2 !
	// 0 hello
	// 1 cruel
	// 2 world
}
