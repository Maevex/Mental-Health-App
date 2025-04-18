package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB
var JWTSecret []byte

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Gagal load .env, lanjut pake environment system...")
	}

	dsn := os.Getenv("DB_URL")
	JWTSecret = []byte(os.Getenv("JWT_SECRET"))

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}

	if err = DB.Ping(); err != nil {
		panic(err.Error())
	}

	fmt.Println("Connected to Database")
}
