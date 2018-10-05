# FAAS

This is a simple demo implementing Function as a service (see AWS Lambda). It utilizes docker to manage the scaling of function.

## Installing
### Requirements
1) Docker
2) Docker-compose
3) Govendor

Install go.

Install docker and docker compose.

Install govendor: 

``` bash
go get -u github.com/kardianos/govendor
```

Clone this repo:

``` bash
git clone https://github.com/kardianos/govendor.git
```


## Running FAAS

``` bash
cd faas
docker-compuose up
```

## Creating a custom function
Creating a function is easy using the faas/function library.

``` golang
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
```

Just pas any regular golang function to RunFunc, and the library will spin up a web server running this function.


to make it work with the FAAS gateway you need to run it in a docker container with with labels 

``` 

faas.name={name of your function}
```


and 

``` 
faas.port={the port you used for its web server}
```

With this lables the FAAS gateway will detect your function container automatically and start forwarding function calls to it. 

