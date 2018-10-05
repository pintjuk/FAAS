package main

import (
	"errors"
	"github.com/pintjuk/faas/function"
)

var memory map[int]int

func factorial(x int) (res int, err error) {
	if x < 0 {
		err = errors.New("Factorial undefined for negative values!")
		return
	}
	err = nil
	res, ok := memory[x]
	if !ok {
		switch x {
		case 0:
			res = 1
		default:
			res, err = factorial(x - 1)
			res *= x
		}
		if res <= 0 {
			err = errors.New("Integer overflow!")
			return
		}
		memory[x] = res
	}
	return
}

func main() {
	memory = make(map[int]int)
	function.RunFunc([]string{"value"}, factorial, ":8080")
}
