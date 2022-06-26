package user

import (
	"errors"
	"fmt"
	"net/http"
	"potato/backend/db"
	"potato/backend/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
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
    err := c.BindJSON(&user)
    if utils.CheckError(err,http.StatusInternalServerError,c){
        return
    }
    con,err := db.Connect()
    if utils.CheckError(err,http.StatusInternalServerError,c){
        return
    }
    defer con.Close()
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

var jwtKey = []byte("my_secret_key")


type Claims struct {
	Id string `json:"id"`
	jwt.StandardClaims
}


func Login(c *gin.Context){
    var user User
    var dbUser db.User
    err := c.BindJSON(&user)
    if utils.CheckError(err,http.StatusInternalServerError,c){
        return
    }
    con,err := db.Connect()
    if utils.CheckError(err,http.StatusInternalServerError,c){
        return
    }
    defer con.Close()
    fmt.Println(user.Email)
    err = con.Get(&dbUser,`select * from users where email=$1`,user.Email)
    if utils.CheckError(err,http.StatusInternalServerError,c){
        return
    }
    err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password),[]byte(user.Password))
    if utils.CheckError(err,http.StatusUnauthorized,c){
        return
    }
    expirationTime := time.Now().Add(5 * time.Hour)
    claims := &Claims{
        Id: string(dbUser.Id),
        StandardClaims: jwt.StandardClaims{
            // In JWT, the expiry time is expressed as unix milliseconds
            ExpiresAt: expirationTime.Unix(),
        },
    }
    // Declare the token with the algorithm used for signing, and the claims
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    // Create the JWT string
    tokenString, err := token.SignedString(jwtKey)
    if utils.CheckError(err,http.StatusUnauthorized,c){
        return
    }
    c.SetCookie("token",tokenString,int(5*60*60),"/","localhost",false,false)
    c.JSON(http.StatusAccepted,gin.H{"status":"Login successfull","token":tokenString})
}

func CheckToken(c *gin.Context){
    token,err := c.Cookie("token")
    if err != nil {
        if err == http.ErrNoCookie {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }
        c.AbortWithStatus(http.StatusBadRequest)
        return
    }
    var claims Claims
    tkn, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if err != nil {
        fmt.Println(err)
        if err == jwt.ErrSignatureInvalid {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }
        c.AbortWithStatus(http.StatusBadRequest)
        return
    }
    if !tkn.Valid {
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }
    c.Set("claims",claims)
    c.Next()
}

func checkClaims(c* gin.Context) (Claims,error) {
    var claims any
    var isPresent bool
    if claims,isPresent = c.Get("claims");!isPresent{
        c.Status(http.StatusUnauthorized)
        return Claims{},errors.New("No Claims, did you forget to add the middleware?")
    }
    _,ok := claims.(Claims)
    if !ok{
        c.Status(http.StatusUnauthorized)
        return Claims{},errors.New("No Claims")
    }
    return claims.(Claims),nil
}

func MiddlewareTest(c * gin.Context){
    claims,err := checkClaims(c)
    if utils.CheckError(err,http.StatusInternalServerError,c){
        return
    }
    c.JSON(http.StatusOK,gin.H{"status":"works!","id":claims.Id})
}

func RefreshToken(c * gin.Context){
    claims,err := checkClaims(c)
    if utils.CheckError(err,http.StatusInternalServerError,c){
        return
    }
    expirationTime := time.Now().Add(5 * time.Minute)
    claims.ExpiresAt = expirationTime.Unix()
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        c.Status(http.StatusInternalServerError)
        return
    }
    c.SetCookie("token",tokenString,int(5*60*60),"/","localhost",false,false)
    c.JSON(http.StatusAccepted,gin.H{"status":"Refreshed","token":tokenString})
}

type location struct{
    Latitude float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
}

func Location(c * gin.Context){
    claims,err := checkClaims(c)
    if utils.CheckError(err,http.StatusBadRequest,c){
        return
    }
    var loc location
    if err := c.ShouldBindJSON(&loc);err != nil{
        c.Status(http.StatusBadRequest)
        return
    }
    con,err := db.Connect()
    if utils.CheckError(err,http.StatusInternalServerError,c){
        return
    }
    defer con.Close()
    _ ,err = con.Exec(`update users set longitude=$1,latitude=$2 where id=$3`,loc.Longitude,loc.Latitude,claims.Id)
    if utils.CheckError(err,http.StatusInternalServerError,c) {
        return
    }
    c.JSON(http.StatusAccepted,gin.H{"update":"successfull"})
}
