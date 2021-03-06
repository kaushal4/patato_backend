package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var schema string = `
create table users (
    id uuid default uuid_generate_v4 (),
    first_name varchar not null,
    last_name varchar not null,
    email varchar not null unique,
    longitude double precision,
    latitude double precision,
    password varchar not null,
    intrests text[],
    primary key (id)
);

create table listings (
    id uuid default uuid_generate_v4 (),
    item text not null,
    freshness text not null,
    photo text,
    userId uuid not null,
    primary key(id),
    constraint fk_users foreign key(userId) references users(id) on delete cascade
);
`

type User struct{
    Id []byte `db:"id"`
    FirstName string `db:"first_name"`
    LastName string `db:"last_name"`
    Email string `db:"email"`
    Longitude sql.NullFloat64 `db:"longitude"`
    Latitude sql.NullFloat64 `db:"latitude"`
    Intrests pq.StringArray `db:"intrests"`
    Password string `json:"password" form:"password" db:"password"`
}

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
    // db.MustExec(schema)
    fmt.Println("Connection to Database successfull")
    return db, nil
}
