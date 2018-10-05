package main

import (
	"github.com/pintjuk/faas/function"
)

func main() {
	function.RunFunc([]string{"a", "b"},
		func(a float64, b float64) float64 {
			return a + b
		},
		":8080")
}
