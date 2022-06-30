package messages

import (
	"net/http"
	"potato/backend/db"
	"potato/backend/user"
	"potato/backend/utils"

	"github.com/gin-gonic/gin"
)

type itemId struct{
    ItemId []byte `json:"transactionid"`
}

func CreateTransaction(c *gin.Context){
    claims,err := user.CheckClaims(c)
    if utils.CheckError(err,http.StatusBadRequest,c){
        return
    }
    var itemid itemId
    err = c.BindQuery(&itemid)
    if utils.CheckError(err,http.StatusBadRequest,c){
        return
    }
    con,err := db.Connect()
    if utils.CheckError(err,http.StatusInternalServerError,c){
        return
    }
    defer con.Close()
    _,err = con.Exec("insert into transactions(itemId,userId,status) values($1,$2,$3)",itemid,claims.Id,"unaccepted")
    if utils.CheckError(err,http.StatusInternalServerError,c){
        return
    }
    c.JSON(http.StatusCreated,gin.H{"status":"successfully created"})
}
