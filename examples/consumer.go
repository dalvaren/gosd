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
