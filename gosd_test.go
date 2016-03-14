// To run:
// go get github.com/githubnemo/CompileDaemon
// CompileDaemon -command="./gervice"

package main

import "encoding/json"
import "github.com/gin-gonic/gin"
import "gopkg.in/validator.v2"
import "github.com/parnurzeal/gorequest"
import "github.com/jinzhu/gorm"
import _ "github.com/mattn/go-sqlite3"

type User struct {
  gorm.Model
  Name  string `binding:"required" validate:"min=3,max=40,regexp=^[a-zA-Z]*$"`
}

func main() {
    // open connection
    DB,_ := gorm.Open("sqlite3", "./service.db")
    DB.DB()

    r := gin.Default()
    r.GET("/param/:id", func(c *gin.Context) {
      c.JSON(200, c.Param("id"))
    })
    r.GET("/me", func(c *gin.Context) {
      request := gorequest.New()
      resp, body, errs := request.Post("http://localhost:3333/pong").
        Send(`{"name":"backy"}`).
        End()

      if errs != nil {
        c.JSON(resp.StatusCode, errs)
        return
      }
      if resp.StatusCode >= 400 {
        c.JSON(resp.StatusCode, body)
        return
      }

      var resCli map[string]interface{}
      json.Unmarshal([]byte(body), &resCli)
      c.JSON(resp.StatusCode, resCli)
    })
    r.GET("/me2", func(c *gin.Context) {
      user := User{}
      user.Name = "Lauren"
      request := gorequest.New()
      resp, body, errs := request.Post("http://localhost:3333/pong").
        Send(user).
        End()
      if errs != nil {
        c.JSON(resp.StatusCode, errs)
        return
      }
      if resp.StatusCode >= 400 {
        c.JSON(resp.StatusCode, body)
        return
      }

      var resCli map[string]interface{}
      json.Unmarshal([]byte(body), &resCli)
      c.JSON(resp.StatusCode, resCli)
    })
    r.GET("/me3", func(c *gin.Context) {
      dat := make(map[string]interface{})
      dat["name"] = "Frank"
      jsonString, _ := json.Marshal(dat)
      request := gorequest.New()
      resp, body, errs := request.Post("http://localhost:3333/pong").
        // Send("{\"name\":\"Rick\"}").
        Send(string(jsonString)).
        End()
      if errs != nil {
        c.JSON(resp.StatusCode, errs)
        return
      }
      if resp.StatusCode >= 400 {
        c.JSON(resp.StatusCode, body)
        return
      }

      var resCli map[string]interface{}
      json.Unmarshal([]byte(body), &resCli)
      c.JSON(resp.StatusCode, resCli)
    })
    r.GET("/ping", func(c *gin.Context) {
      user := User{}
      user.Name = "John"
      address := make(map[string]interface{})
      address["city"] = "São Paulo"
      address["state"] = "São Paulo"
      address["country"] = "Brazil"
      dat := make(map[string]interface{})
      dat["a"] = 1
      dat["b"] = "my"
      dat["user"] = user
      dat["address"] = address
        c.JSON(200, dat)
    })
    r.POST("/pong", func(c *gin.Context) {
      newUser := User{}
      c.BindJSON(&newUser)
      err := validator.Validate(newUser)
      if err != nil {
        c.JSON(400, err.Error())
        return
      }
      c.JSON(200, newUser)
    })
    r.GET("/migrate", func(c *gin.Context) {
      // open connection
      db,_ := gorm.Open("sqlite3", "./service.db")
      db.DB()

      // migration
      db.CreateTable(&User{})
      c.JSON(200, "success")
    })
    r.GET("/adduser", func(c *gin.Context) {
      // open connection
      db,_ := gorm.Open("sqlite3", "./service.db")
      db.DB()

      // save to DB
      user := User{Name: "Jinzhu"}
      db.Create(&user)

      users := db.Find(&[]User{})
      c.JSON(200, users)
    })
    r.GET("/where", func(c *gin.Context) {
      users := DB.Find(&[]User{})
      c.JSON(200, users)
    })
    r.Run(":3333") // listen and server on 0.0.0.0:3333
}
