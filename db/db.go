package db

import (
	"errors"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var password string 
const (
    host     = "john.db.elephantsql.com"
    port     = 5432
    user     = "utpvxqlz"
    dbname   = "utpvxqlz"
)

func getPassword() error {
    if pass,isPresent := os.LookupEnv("POSTELE_PASS");isPresent {
        password = pass
        return nil
    }else {
        return errors.New("Enviorment variable not set")
    }
}

func Connect() (*sqlx.DB,error) {
    err := getPassword()
    if err != nil {
        return nil ,err
    }
    psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
    db, err := sqlx.Open("postgres", psqlconn)
    if err != nil {
        fmt.Println("Connection Failed!")
        fmt.Println(err)
        return nil, err
    }
    fmt.Println("Connection to Database successfull")
    return db, nil
}

