package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

func CheckError(err error,code int,c * gin.Context) bool {
    if err != nil{
        fmt.Println(err.Error())
        c.JSON(code,gin.H{"status":err.Error()})
        return true
    }
    return false
}

func UploadImage(path string,file *multipart.FileHeader,c * gin.Context) (string,error) {
    fileType := file.Header["Content-Type"][0]
    if !strings.Contains(fileType,"image") {
        c.JSON(http.StatusBadRequest,gin.H{"status":"upload a image dummy"})
        return "",errors.New("The file is not a image")
    }
    u2 := uuid.NewV4().String()
    filename := u2+file.Filename[strings.LastIndex(file.Filename,"."):]
    err := c.SaveUploadedFile(file, path+filename)
    if CheckError(err,http.StatusInternalServerError,c){
        return "",errors.New("Could'nt save file to filesystem")
    }
    return filename,nil
}
