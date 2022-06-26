package user

import (
	"fmt"
	"net/http"
	"potato/backend/db"
	"potato/backend/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)


type User struct{
    FirstName string `json:"first_name" form:"first_name" db:"first_name"`
    LastName string `json:"last_name" form:"last_name" db:"last_name"`
    Email string `json:"email" form:"email" db:"email"`
    Password string `json:"password" form:"password" db:"password"`
    Longitude float64 `json:"longitude" form:"longitude" db:"longitude"`
    Latitude float64 `json:"latitude" form:"latitude" db:"latitude"`
}


func SignUp(c *gin.Context){
    var user User
    err := c.Bind(&user)
    if utils.CheckError(err,http.StatusInternalServerError,c){
        return
    }
    con,err := db.Connect()
    if utils.CheckError(err,http.StatusInternalServerError,c){
        return
    }
    hashedPassword,err := bcrypt.GenerateFromPassword([]byte(user.Password),8)
    if utils.CheckError(err,http.StatusInternalServerError,c){
        return
    }
    user.Password = string(hashedPassword)
    _,err = con.NamedExec("insert into users(first_name,last_name,email,password) values(:first_name,:last_name,:email,:password)",&user)
    if utils.CheckError(err,http.StatusInternalServerError,c){
        return
    }
    c.JSON(http.StatusCreated,gin.H{"status":"Insertion Successfull"})
}


func Login(c *gin.Context){
    var user User
    var dbUser db.User
    err := c.Bind(&user)
    if utils.CheckError(err,http.StatusInternalServerError,c){
        return
    }
    con,err := db.Connect()
    if utils.CheckError(err,http.StatusInternalServerError,c){
        return
    }
    fmt.Println(user.Email)
    err = con.Get(&dbUser,`select * from users where email=$1`,user.Email)
    if utils.CheckError(err,http.StatusInternalServerError,c){
        return
    }
    err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password),[]byte(user.Password))
    if utils.CheckError(err,http.StatusUnauthorized,c){
        return
    }
    c.JSON(http.StatusAccepted,gin.H{"status":"Login successfull"})
}
