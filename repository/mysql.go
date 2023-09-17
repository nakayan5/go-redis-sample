package repository

import (
	// "database/sql"
	"fmt"
	"time"
	// "os"
	
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func Setup() {
	db := Connect()
	defer db.Close()
}

func Connect() *sqlx.DB {
	config := mysql.Config{
		DBName:               "db",
		User:                 "user",
		Passwd:               "password",
		AllowNativePasswords: true,
		Addr:                 "db" + ":" + "3306",
		Net:                  "tcp",
		ParseTime:            true,
		Loc:                  time.Local,
		InterpolateParams:    true,
	}
	db, err := sqlx.Open("mysql", config.FormatDSN())

	if err != nil {
		panic(err)
	}

	fmt.Println("db connected!!")
	return db
}