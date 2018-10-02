package function

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

func RunFunction(f func(map[string][]string) (string, error)) {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {

		fmt.Println("#Recived Request:")
		fmt.Println("\t- ", (c.Request()).URL.Host)
		fmt.Println("\t- ", (c.Request()).URL.Path)
		params := c.QueryParams()
		fmt.Println("\t- pram1i: ", params)
		res, err := f(params)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest,
				err.Error())
		}
		return c.String(http.StatusOK, res)
	})
	e.Logger.Fatal(e.Start(":8080"))
}
