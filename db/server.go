package db

import (
	"flag"

	"golang-graphql-crud/utils"
	"log"

	"github.com/joho/godotenv"
)

func Run(){
	
	var dbConfig = DbConfig{}
	err:= godotenv.Load()
	if err != nil{
		log.Fatalf("Error Getting .env File")
	}

	dbConfig.DBHost 		= utils.GetEnv("DB_HOST", "localhost")
	dbConfig.DBUser			= utils.GetEnv("DB_USER", "postgres")
	dbConfig.DBPassword		= utils.GetEnv("DB_PASSWORD", "12345678")
	dbConfig.DBName			= utils.GetEnv("DB_NAME", "go-graphql")
	dbConfig.DBPort			= utils.GetEnv("DB_PORT", "5432")

	flag.Parse()
	arg := flag.Arg(0)
	if arg != "" {
		// initializeDB(dbConfig)
	}else{
		InitializeDb(dbConfig)
	
	}
	

}