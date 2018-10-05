package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/labstack/echo"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"
	"github.com/labstack/echo/middleware"
)

type RouteTable struct {
	r   unsafe.Pointer
	cli *client.Client
}

func (route *RouteTable) Get() map[string][]string {
	return *(*map[string][]string)(route.r)
}

func (route *RouteTable) runRebuilder() {
	for {
		time.Sleep(10 * time.Second)
		route.rebuild()
	}
}
func (route *RouteTable) rebuild() {
	newRoutes := make(map[string][]string)
	containers, err := route.cli.ContainerList(
		context.Background(),
		types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {

		funcName, ok := container.Labels["faas.name"]
		if ok {
			containerName := container.Names[0]
			port := container.Labels["faas.port"]
			serveraddr := fmt.Sprint("http:/", containerName, ":", port)
			old, ok := newRoutes[funcName]
			if ok {
				newRoutes[funcName] = append(old, serveraddr)
			} else {
				newRoutes[funcName] = []string{serveraddr}
			}
			fmt.Println("adding container")
			fmt.Println("wih addr: ", serveraddr)
			fmt.Println("%s %s\n", container.ID[:10],
				container.Image)
			fmt.Println("to serve function ", funcName)

		}
	}
	//hotswap routes
	//WARNING: alterations to newRoute past this point are not safe
	atomic.StorePointer(&route.r, unsafe.Pointer(&newRoutes))
}

func makeRouteTable() (res RouteTable) {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.37"))
	if err != nil {
		panic(err.Error())
	}
	res = RouteTable{cli: cli, r: unsafe.Pointer(&map[string][]string{})}
	res.rebuild()
	return
}




func makeGatewayMV() func(echo.HandlerFunc) echo.HandlerFunc {

	// create route table and start update process
	rt := makeRouteTable()
	go rt.runRebuilder()

	getFuncName := func (path string) (funcName string,err error){
		funcName = strings.Trim(path, " /")
		funcNameL := strings.Split(funcName, "/")
		if len(funcNameL) < 1 {
			err = echo.NewHTTPError(http.StatusBadRequest, "No such function")
			return
		}
		funcName = funcNameL[0]
		return
	}

	counter :=map[string]int{} 
	getNext := func(name string, max int) (next int) {
		cur, ok := counter[name]
		if ok {
			next = (cur+1) % max
		}else{
			next =0
		}
		counter[name]= next
		return 
	}
	
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			funcName, err := getFuncName(req.URL.Path)
			if err!=nil {return err}
			req.URL.Path = ""

			// retrive path to server
			urls, ok := rt.Get()[funcName]
			if !ok || len(urls) < 1 {
				return echo.NewHTTPError(http.StatusBadRequest,
					fmt.Sprint("ERROR: ", http.StatusBadRequest,
						", No such function: ", funcName))
			}
			url, err := url.Parse(urls[ getNext(funcName, len(urls))])
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest,
					err.Error())
			}
			fmt.Println("farwording to ", url)

			// fix heder
			if req.Header.Get(echo.HeaderXRealIP) == "" {
				req.Header.Set(echo.HeaderXRealIP, c.RealIP())
			}
			if req.Header.Get(echo.HeaderXForwardedProto) == "" {
				req.Header.Set(echo.HeaderXForwardedProto, c.Scheme())
			}
			if c.IsWebSocket() && req.Header.Get(echo.HeaderXForwardedFor) == "" {
				// For HTTP, it is automatically set by Go HTTP reverse proxy.
				req.Header.Set(echo.HeaderXForwardedFor, c.RealIP())
			}

			//proxy
			httputil.NewSingleHostReverseProxy(url).ServeHTTP(res, req)
			return echo.NewHTTPError(http.StatusBadRequest, "there is nothing here")
		}
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(makeGatewayMV())
	e.Logger.Fatal(e.Start(":80"))

}
