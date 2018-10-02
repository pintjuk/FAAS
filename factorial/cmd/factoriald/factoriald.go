package main

import (
	"errors"
	"github.com/pintjuk/faas/function"
	"strconv"
	"strings"
)

var memory map[int]int

func factorial(x int) (res int, err error) {
	if x < 0 {
		err = errors.New("ERROR:\n\tFactorial undefined for negative values!")
		return
	}
	err = nil
	res, ok := memory[x]
	if !ok {
		switch x {
		case 0:
			res = 1
		case 1:
			res = 1
		default:
			res, err = factorial(x - 1)
			res *= x
		}
		if res <= 0 {
			err = errors.New("Integer overflow!")
		}
		memory[x] = res
	}
	return
}

func factorialWrap(param map[string][]string) (out string, err error) {
	param1s, ok := param["param-1"]
	if !ok {
		err = errors.New("ERROR:\n\t Mising param-1")
		return
	}
	param1i, err := strconv.Atoi(strings.Join(param1s, ""))
	if err != nil {
		err = errors.New("Parameters are:\n\tparam-1 (integer)")
		return
	}
	res, err := factorial(param1i)
	out = strconv.Itoa(res)
	return
}

func main() {
	memory = make(map[int]int)
	function.RunFunction(factorialWrap)
}
