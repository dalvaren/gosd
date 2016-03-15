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
