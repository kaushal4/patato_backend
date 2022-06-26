package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
    "potato/backend/user"
)

func main() {
    port := 8080
    fmt.Println("starting server in port :",port)
    router := gin.New()
    userRoutes := router.Group("/user")
    {
        userRoutes.POST("/",user.SignUp)
        userRoutes.POST("/login",user.Login)
        userRoutes.POST("/location",user.CheckToken,user.Location)
        userRoutes.GET("/test",user.CheckToken,user.MiddlewareTest)
        userRoutes.GET("/refresh",user.CheckToken,user.RefreshToken)

    }
    router.Run(fmt.Sprintf("localhost:%d",port))
}


