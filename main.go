package main

import (
	"fmt"
	"net/http"
	"potato/backend/db"
	"potato/backend/schema"

	"github.com/gin-gonic/gin"
)

func getPersons(c *gin.Context) {
    con,err := db.Connect()
    if err != nil{
        fmt.Println(err)
        c.JSON(http.StatusInternalServerError,gin.H{"error":"true"})
        return;
    }
    defer con.Close()
    query := fmt.Sprintf(`Select * from persons limit 5`)
    rows,err := con.Queryx(query)
    if err!=nil{
        fmt.Printf("query failed!")
        c.JSON(http.StatusInternalServerError,gin.H{"error":"true"})
        return;
    }
    defer rows.Close()
    results := make([]schema.Persons,0,6)
    var i int = 0
    for rows.Next() {
        var person schema.Persons
        err := rows.StructScan(&person)
        if err != nil {
            fmt.Println("Can't read the row")
            continue
        }
        results = results[:i+1]
        results[i] = person
    }
    c.JSON(http.StatusAccepted,results)
}


func main() {
    port := 8080
    fmt.Println("starting server in port :",port)
    router := gin.Default()
    router.GET("/persons",getPersons)
    router.Run(fmt.Sprintf("localhost:%d",port))
}


