package main

import (
	"fmt"
	"potato/backend/listings"
	"potato/backend/user"

	"github.com/gin-gonic/gin"
)

func main() {
    port := 8080
    fmt.Println("starting server in port :",port)
    router := gin.New()
    router.MaxMultipartMemory = 8<<20 //8mb
    fmt.Println(8<<20)
    userRoutes := router.Group("/user")
    {
        userRoutes.POST("/",user.SignUp)
        userRoutes.GET("/",user.CheckToken,user.GetUser)
        userRoutes.POST("/login",user.Login)
        userRoutes.POST("/location",user.CheckToken,user.Location)
        userRoutes.POST("/intrest",user.CheckToken,user.Intrests)
        userRoutes.GET("/test",user.CheckToken,user.MiddlewareTest)
        userRoutes.GET("/refresh",user.CheckToken,user.RefreshToken)

    }
    listRoutes := router.Group("/listings")
    {
        listRoutes.POST("",user.CheckToken,listings.Listings)
    }
    router.Run(fmt.Sprintf("localhost:%d",port))
}


