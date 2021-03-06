/*
This is a library for running a go function as a HTTP micro service

Just call RunFunc with your function, and it will create and run a web-server
        Runfunc([]string{"arg1", "arg2"}, myFunction, ":8080")
*/
package function

import (
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"reflect"
	"strconv"
)

func validKind(k reflect.Kind) bool {
	switch k {
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return true

	case reflect.String:
		return true
	case reflect.Bool:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}

}

func parseValue(s string, k reflect.Kind) (r reflect.Value, err error) {
	switch k {
	case reflect.Uint:
		res, err := strconv.ParseUint(s, 10, 64)
		return reflect.ValueOf(uint(res)), err
	case reflect.Uint8:
		res, err := strconv.ParseUint(s, 10, 64)
		return reflect.ValueOf(uint8(res)), err
	case reflect.Uint16:
		res, err := strconv.ParseUint(s, 10, 64)
		return reflect.ValueOf(uint16(res)), err
	case reflect.Uint32:
		res, err := strconv.ParseUint(s, 10, 64)
		return reflect.ValueOf(uint32(res)), err
	case reflect.Uint64:
		res, err := strconv.ParseUint(s, 10, 64)
		return reflect.ValueOf(uint64(res)), err
	case reflect.Int:
		res, err := strconv.ParseInt(s, 10, 64)
		return reflect.ValueOf(int(res)), err
	case reflect.Int8:
		res, err := strconv.ParseInt(s, 10, 64)
		return reflect.ValueOf(int8(res)), err
	case reflect.Int16:
		res, err := strconv.ParseInt(s, 10, 64)
		return reflect.ValueOf(int16(res)), err
	case reflect.Int32:
		res, err := strconv.ParseInt(s, 10, 64)
		return reflect.ValueOf(int32(res)), err
	case reflect.Int64:
		res, err := strconv.ParseInt(s, 10, 64)
		return reflect.ValueOf(int64(res)), err
	case reflect.String:
		return reflect.ValueOf(s), nil
	case reflect.Bool:
		res, err := strconv.ParseBool(s)
		return reflect.ValueOf(res), err
	case reflect.Float32:
		res, err := strconv.ParseFloat(s, 32)
		return reflect.ValueOf(float32(res)), err
	case reflect.Float64:
		res, err := strconv.ParseFloat(s, 64)
		return reflect.ValueOf(float64(res)), err
	default:
		panic("Function has bad parameter type")
	}
}

func callFunc(f interface{}, argNames []string, c echo.Context) (res string, err error) {
	fValue := reflect.ValueOf(f)
	fType := fValue.Type()

	args := make([]reflect.Value, len(argNames))

	for i, argName := range argNames {
		args[i] = reflect.Zero(fType.In(i))
		queryValue := c.QueryParam(argName)
		fmt.Println(argName)
		fmt.Println(argName, ": ", queryValue)

		args[i], err = parseValue(queryValue, fType.In(i).Kind())
		if err != nil {
			err = errors.New(fmt.Sprint("Bad parameter: ", queryValue))
		}
	}
	ress := fValue.Call(args)
	res = fmt.Sprint(ress[0].Interface())
	if fType.NumOut() > 1 {
		switch v := ress[1].Interface().(type) {
		case error:
			err = v
		}
	}
	return
}

/*
Starts a webs server to run this function f  as a micro service.

ArgNames

List of query parameter names to use as input arguments to your function.

- have to match arguments of function f.

Function

Function used for processing requests.

- f's arguments have a primitive or string type.

- f first return value have a perimeter or a string type.

- f may have a second return value, witch has to be an error.
RunFunc panics if any of above conditions are violated.

Path

path and port number for the http server.
*/
func RunFunc(argNames []string, function interface{}, addr string) {
	/* Validate input */
	fValue := reflect.ValueOf(function)
	fType := fValue.Type()
	if fValue.Kind() != reflect.Func {
		panic("First argument to RunFunc has to be a function")
	}
	if fType.NumOut() > 2 {
		panic(fmt.Sprint("Function has to have 1 or 2 return values, it has ", fType.NumOut()))
	}

	// TODO: this comparison does not really work for some reason
	/*
		if fType.NumOut()==2 && (fType.Out(1) != reflect.TypeOf(errors.New(""))) {
			panic(fmt.Sprint("Second return value has to be of type error, it is: ", fType.Out(1).Name(), ",",reflect.TypeOf().Name()))
		}*/

	if fType.NumIn() != len(argNames) {
		panic("Argument Names do not match the function arguments")
	}
	if !validKind(fType.Out(0).Kind()) {
		panic("First output has invalid type, use only primitive types and strings")
	}
	for i, _ := range argNames {
		if !validKind(fType.In(i).Kind()) {
			panic(fmt.Sprint("Argument number ", i, "has invalid type, use only primitive types and strings"))
		}
	}
	/* Start web-server */
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.GET("/", func(c echo.Context) error {

		res, err := callFunc(function, argNames, c)
		if err != nil {
			fmt.Println(err)
			return echo.NewHTTPError(http.StatusBadRequest,
				err.Error())
		}
		return c.HTML(http.StatusOK, res)
	})
	e.Logger.Fatal(e.Start(addr))
}
