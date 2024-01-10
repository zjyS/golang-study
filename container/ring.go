package main

import (
	"container/ring"
	"fmt"
)

func forEach() {
	r := ring.New(3)
	r.Value = 1
	for p := r.Next(); p != r; p = p.Next() {
		p.Value = 1
	}
	r.Do(Println)
}

func Println(v any) {
	fmt.Println(v)
}
