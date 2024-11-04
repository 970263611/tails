package main

import "fmt"

func Aaa() func(int, int) int {
	f := func(a, b int) int {
		return a + b
	}
	return f
}

func main() {
	str := "-a"
	fmt.Println(str[1:])
}
