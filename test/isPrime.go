package main

import "fmt"

func main() {
	ans := isPrime(13)
	fmt.Println(ans)
}

func isPrime(x int) bool {
	for i := 2; i < x; i++ {
		if x%i == 0 {
			return false
		}
	}
	return true
}
