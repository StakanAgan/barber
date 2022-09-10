package src

import (
	"fmt"
	"github.com/go-redis/redis/v9"
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
		DBName:   os.Getenv("EDGEDB_SERVER_DATABASE"),
		Host:     "db",
		User:     os.Getenv("EDGEDB_SERVER_DATABASE"),
		Password: os.Getenv("EDGEDB_SERVER_PASSWORD"),
		Port:     port,
	}
}

func NewRedisConfig() *redis.Options {
	err := godotenv.Load()
	port, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		log.Fatal(err)
	}
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatal(err)
	}
	redisAddr := fmt.Sprintf("%s:%d", os.Getenv("REDIS_HOST"), port)
	return &redis.Options{
		Addr:     redisAddr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	}
}
