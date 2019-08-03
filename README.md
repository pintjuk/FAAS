# Toy FaaS implementation in Golang
This is just a litle project i did to learn docker, docker compose and golang.

So this docker thing seems to be getting some traction, and so does function as a service (FaaS). It would be nice if we could have our own FaaS infrastructure in place instead of paying Amazon for it. There are some alternatives out there e.g. OpenFaaS. But we have a bad case of “Not Invented Here” syndrome, so we built one ourselves.

![system design](https://pintjuk.github.io/img/2019-06/faas/scetch.png)

This is a small GoLang infrastructure app that enables one to deploy simple REST servers, acting as single functions. The goal was to deploy and scale the functions horizontally with ease. To achieve this, we leveraged the scaling and service discovery that comes for free in a docker-compose and swarm.

The core component is the FaaS Gateway. The image above illustrates the idea. The gateway uses docker to discover function dynamically and updating an internal table, a table of functions, and the containers providing them.  When the gateway receives the request: `GET lambda/function-3?param-1=123&param-2=456` redirect the request containers with the tag `faas:name=function-3` in a round-robin fashion.

## Creating your function

A function service is a docker container running a micro-service implemented as a single simple golang function.

Creating a function service is straight forward using the faas/function library.

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

Just pass any GoLang function to RunFunc, and the library will spin up a web server running this function. It will automatically detect the number of parameters and map them to query parameters `?param-1&param-2& ...&param-n` of the incoming requests.

To make it work with the FAAS-Gateway, you need to deploy it in a docker container with the following labels:

``` 
faas.name={name of your function}
```

and

``` 
faas.port={the port you used for its web server}
```

With these labels, the FAAS gateway will detect your function container automatically and start forwarding function calls to it as soon as your container is up.

This is a docker file for the example function addition function above:

FROM golang:1.9

```Dockerfile
# go stuff
RUN go get -u github.com/kardianos/govendor && \
go get -u -v github.com/codeskyblue/fswatch
RUN mkdir -p /go/src/github.com/pintjuk/faas/addition
WORKDIR /go/src/github.com/pintjuk/faas/addition
EXPOSE 8080
CMD govendor fetch +m ; go run cmd/*/*.go
```

And this is its docker compose entry:

```ymal
addition:
        build:
            context: ./addition
        labels:
            - "faas.name=addition"
            - "faas.port=8080"
        volumes:
            - ./addition:/go/src/github.com/pintjuk/faas/addition
```

## Installation

### Requirements

1. Docker
2. Docker-compose
3. Govendor

### Installation steps

1. Install GoLang.
2. Install docker and docker-compose.
3. Install Govendor:
``` bash
go get -u github.com/kardianos/govendor
```
1. Clone this repo:

``` bash
git clone https://github.com/kardianos/govendor.git
```

## Running FAAS

1. up faas-gateway and function services

``` bash
cd faas
docker-compose up
```
