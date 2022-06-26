package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
    "potato/backend/user"
)

func main() {
    port := 8080
    fmt.Println("starting server in port :",port)
    router := gin.Default()
    router.POST("/user",user.SignUp)
    router.POST("/user/login",user.Login)
    router.Run(fmt.Sprintf("localhost:%d",port))
}


