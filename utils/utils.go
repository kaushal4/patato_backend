package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func CheckError(err error,code int,c * gin.Context) bool {
    if err != nil{
        fmt.Println(err.Error())
        c.JSON(code,gin.H{"status":err.Error()})
        return true
    }
    return false
}
