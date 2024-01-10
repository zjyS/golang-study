package main

import (
	"fmt"
)

type SS []int

func main() {
	a := []int{1, 2, 3, 4, 5, 6}
	s1 := a[:3]
	s2 := a[4:]
	s3 := s2[:1]
	fmt.Println("print a")
	forEach(a)
	fmt.Println("print s1")
	forEach(s1)
	fmt.Println("print s2")
	forEach(s2)
	fmt.Println("print s3")
	forEach(s3)

	b := SS{1, 2, 3, 4, 5}
	fmt.Println("before forEach")
	fmt.Println(b)
	b.forEach(func(a *int) {
		*a = *a + 1
	})
	fmt.Println("after forEach")
	fmt.Println(b)

}

func forEach(i []int) {
	for n := 0; n < len(i); n = n + 1 {
		fmt.Println(i[n])
	}
}

func (s SS) forEach(op func(*int)) {
	for n := 0; n < len(s); n = n + 1 {
		op(&s[n])
	}
}
