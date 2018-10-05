package main

import (
//	"context"
	//	"errors"
	//	"fmt"
	//	"github.com/docker/docker/api/types"
	//	"github.com/docker/docker/api/types/filters"
	//	"github.com/docker/docker/client"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/url"
)

func main() {
	/*	cli, err := client.NewClientWithOpts(client.WithVersion("1.37"))
		if err != nil {
			panic(err.Error())
		}*/

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// Setup proxy
	url1, err := url.Parse("http://faas_addition_1:8080")
	if err != nil {
		e.Logger.Fatal(err)
	}
	url2, err := url.Parse("http://faas_factorial_1:8080")
	if err != nil {
		e.Logger.Fatal(err)
	}
	targets := []*middleware.ProxyTarget{
			{
				URL: url1,
			},
			{
				URL: url2,
			},
		}
	balancer := middleware.NewRoundRobinBalancer(targets)
     //   g := e.Group("/e")
	e.Use(middleware.Proxy(balancer))
	
	e.Logger.Fatal(e.Start(":80"))

}

/*
func main() {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.37"))
	if err != nil {
		panic(err.Error())
	}

	e := echo.New()
	e.GET("/:name", func(c echo.Context) error {
		name := c.Param("name")
		filters := filters.NewArgs()
		filters.Add("label", fmt.Sprint("faas.name=", name))

		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{Filters: filters})

		if len(containers) == 0 {
			err = errors.New(fmt.Sprint("No container with label faas.name=", name))
		}
		/*
			for _, container := range containers {
				fmt.Println("filtered %s %s\n", container.ID[:10], container.Image)
				fmt.Println()
			}*/
/*container_name := containers[0].Names[0]
		port := containers[0].Labels["faas.port"]
		uri := fmt.Sprint("https:/", container_name, ":", port, "/?a=1&b=2")
		fmt.Println("GET: ", uri)
		resp, err := http.Get(uri)
		if err != nil {
			fmt.Println(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		respCode := resp.StatusCode
		fmt.Println(respCode, " | ", body, " | ", err)

		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())

		} else {
			return c.String(http.StatusOK, "Hello, World!")
		}
	})
	e.Logger.Fatal(e.Start(":80"))
	for {
		fmt.Println("Changed sorce")
	}
}*/
