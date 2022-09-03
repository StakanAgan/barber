package src

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type DBConfig struct {
	DBName   string
	Host     string
	User     string
	Password string
	Port     int
}

func NewDBConfig() *DBConfig {
	err := godotenv.Load()

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal(err)
	}
	return &DBConfig{
		DBName:   "benny",
		Host:     "db",
		User:     os.Getenv("EDGEDB_SERVER_DATABASE"),
		Password: os.Getenv("EDGEDB_SERVER_PASSWORD"),
		Port:     port,
	}
}
