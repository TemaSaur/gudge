package main

import "fmt"

func sum(a, b int) int {
	return a + b
}

func main() {
	var t int
	fmt.Scan(&t)
	for range t {
		var a, b int
		fmt.Scan(&a, &b)
		fmt.Println(sum(a, b))
	}
}
