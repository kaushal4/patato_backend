package db

import (
    "fmt"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

const (
    host     = "john.db.elephantsql.com"
    port     = 5432
    user     = "utpvxqlz"
    password = "NXd_Nu8q9i1kA1KYkCp5pTZvgE8i0KqI"
    dbname   = "utpvxqlz"
)

func Connect() *sqlx.DB {
    psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
    db, err := sqlx.Open("postgres", psqlconn)
    if err != nil {
        fmt.Println("Connection Failed!")
        fmt.Println(err)
        return nil
    }
    fmt.Println("Connection to Database successfull")
    return db
}

