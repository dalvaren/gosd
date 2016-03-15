# gosd

Scale your Go microservices in 5 minutes with this simple service discovery client for small architectures.

It's ideal for starting projects because you don't need a separated network to run your service discovery. So it's very important you set your service discovery Database in a high available system, like Amazon ElastiCache.

Today the only driver implemented works for Redis, but other are planned (etcd, SQLite, custom API...).

## Installation

Just add the library to your project.

```
$ go get github.com/dalvaren/gosd
```

## Configuration

Following the [12 factor app rules](http://12factor.net/), the configurations for this library shall be done using environment variables. Below are the list with the needed configuration

```
# For the Redis Driver
export gosdRedisDB=0                  # Redis database
export gosdRedisAddr="localhost:6379" # Redis host address
export gosdRedisPassword=""           # Redis password

# For the Library
export gosdTryRefreshAmount=3         # Number of times the library tries to reach the Service Discovery database for updating.
export gosdTryFindServiceAmount=5     # Number of times the service tries to find the service URL in its cache. It's important to wait for other services to start.
export gosdTryFindServiceDelay=3      # Delay in seconds between attempts of gosdTryFindServiceAmount
```

## Usage

First, you need to import it.

```
import gosd "github.com/dalvaren/gosd"
```

Start the Service Discovery client and choose the desired driver. Do it at your application startup. Basically, this command makes the microservice register itself in service discovery when it starts.

```
currentName := gosd.Start("my-service-name", "http://localhost:8885", gosd.DriverRedis{})
```

***currentName*** is the unique name of your service, used to close only it when this microservice closes. ***my-service-name*** is the desired service name. ***"http://localhost:8885"*** is an example of a reachable host for this service. You can implement some way to get it automatically. The service will unregister automatically with some problem occurs with the application.

Now every time you want to get the most updated version of the registered services list you run:

```
gosd.Get()
```

And for the last, when you need some specific microservice URL (for an API request for example) you just call:

```
gosd.IterateServiceRoute("my-other-service-name")
```

This command iterates on local service list, using a Round Robin algorithm to give you the next service URL.

And that's it.
If you set right the configurations, you can start microservices on demand and they will be discovered by who is using it, after they call `gosd.Get()` or `gosd.UpdateByCron()` .

You can see other (and important) features in next section.

## Advanced option

1. Update using GOSD cron. You can call this command on all microservice endpoints, put it in a middleware or something similar.

  ```
  gosd.UpdateByCron()
  ```

1. Delete some service manually, by its URL. It shall be done after you make some request and the service returns no response (you can perform 3 to 5 attempts before removing some server). Take care with this command, since you can unregister the own service:

  ```
  gosd.DeleteServiceWithURL("http://localhost:8881")
  ```

1. Sometimes it's interesting to register some services locally, for that you can use the command below. But remember to register again after each `gosd.Get()` or `gosd.UpdateByCron()` :

  ```
  gosd.AddServiceManually("service-name", "http://localhost:8886")
  ```


## Example

As example let's run 2 identical services called "provider", where they responds with their unique IDs. I'm using Gin as framework.

```
// provider.go
package main

import "os"
import "github.com/gin-gonic/gin"
import "github.com/dalvaren/gosd"

func main() {
  gosd.Start("provider", "http://localhost" + os.Args[1], gosd.DriverRedis{})
  gosd.Get()

  router := gin.Default()
  router.GET("/ping", func(ginContext *gin.Context) {
    gosd.UpdateByCron()
    ginContext.JSON(200, "Provider: " + os.Args[2])
  })
  router.Run(os.Args[1])
}

```

  - Open a terninal and run that with `go run provider.go :3333 1`
  - Open another terninal and run that with `go run provider.go :3334 2`

And the "consumer", who makes 10 requests for the 2 services. Note that it does not need to know them individually and you also can start them some seconds later.

```
// consumer.go
package main

import "fmt"
import "github.com/dalvaren/gosd"
import "github.com/parnurzeal/gorequest"

func main() {
  gosd.Start("consumer", "http://localhost:fake", gosd.DriverRedis{})
  gosd.Get()

  for n:=0; n < 10; n++ {
    providerURL := gosd.IterateServiceRoute("provider")
    request := gorequest.New()
    _, body, _ := request.Get(providerURL + "/ping").End()
    fmt.Print(body)
  }
}

```

## Contribute creating drivers for it!

Right now, I've only implemented a driver for Redis.
Fork it and feel free to develop other drivers (etcd, memcache, custom apis...).
The driver documentation can be found [here](https://github.com/dalvaren/gosd/blob/master/DRIVERS.md)!

## Author

Daniel Campos
