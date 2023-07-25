package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

var DB *sql.DB

func InitializeDb(dbConfig DbConfig){
	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	dbConfig.DBHost,
	dbConfig.DBPort,
	dbConfig.DBUser,
	dbConfig.DBPassword,
	dbConfig.DBName,
)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)

	fmt.Println("konek ke pg")
	// defer db.Close()

	DB = db
}