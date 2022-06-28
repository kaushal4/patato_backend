package listings

import (
    "potato/backend/db"
    "potato/backend/user"
	"fmt"
	"net/http"
    "potato/backend/utils"
	"github.com/gin-gonic/gin"
)

type Listings struct {

}

func Listings(c *gin.Context){
    claims,err := user.CheckClaims(c)
    if utils.CheckError(err,http.StatusBadRequest,c){
        return
    }
    file,err := c.FormFile("photo")
    if err != nil{
        fmt.Println(err)
        c.Status(http.StatusBadRequest)
        return
    }
    filename,err := utils.UploadImage("listings/photos/",file,c)
    if err != nil {
        fmt.Println(err)
        return
    }
    con,err := db.Connect()
    if utils.CheckError(err,http.StatusInternalServerError,c){
        return
    }
    defer con.Close()
    c.JSON(http.StatusOK, gin.H{"status":"inserted"})
}
