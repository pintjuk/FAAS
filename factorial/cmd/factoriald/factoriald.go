package main

import (
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
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
		memory[x] = res
	}
	return
}

func main() {
	memory = make(map[int]int)
	e := echo.New()
	e.GET("/", func(c echo.Context) error {

		fmt.Println("#Recived Request:")
		fmt.Println("\t- ", (c.Request()).URL.Host)
		fmt.Println("\t- ", (c.Request()).URL.Path)
		param1s := c.QueryParam("param-1")
		fmt.Println("\t- pram1: ", param1s)
		param1i, err := strconv.Atoi(param1s)
		if err != nil {
			fmt.Println("ERROR:\n\t invalid parameter in request")
			return c.String(http.StatusBadRequest,
				"Parameters are:\n\tparam-1 (integer)")
		}

		fmt.Println("\t- pram1i: ", param1i)
		fmt.Println("\t- err: ", err)
		res, err := factorial(param1i)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest,
				err.Error())
		}
		return c.String(http.StatusOK, strconv.Itoa(res))
	})
	e.Logger.Fatal(e.Start(":8080"))
}
