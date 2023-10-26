package config

import (
	"strconv"

	"fmt"
	"github.com/joho/godotenv"
	"os"
	"time"
)

type Config struct {
	Env        string
	DBType     string
	PgSQL      PgSQL
	Redis      Redis
	HttpServer HttpServer
}

type PgSQL struct {
	Port     string
	Host     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type Redis struct {
	Port     string
	Host     string
	Password string
	DB       int
	Timeout  time.Duration
}

type HttpServer struct {
	Port         string
	Host         string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

func LoadConfig(configPath string) (Config, error) {
	if configPath == "" {
		return Config{}, fmt.Errorf("config path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return Config{}, fmt.Errorf("config file does not exist: %s", configPath)

	}

	err := godotenv.Load(configPath)
	if err != nil {
		return Config{}, err
	}
	var config Config

	config.Env = os.Getenv("ENV")
	config.DBType = os.Getenv("DBType")

	config.PgSQL = PgSQL{
		Port:     os.Getenv("POSTGRES_PORT"),
		Host:     os.Getenv("POSTGRES_HOST"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
		SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}

	config.Redis = Redis{
		Port:     os.Getenv("REDIS_PORT"),
		Host:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       getInt("REDIS_DB"),
		Timeout:  time.Second * getDuration("REDIS_TIMEOUT"),
	}

	config.HttpServer = HttpServer{
		Port:         os.Getenv("SERVER_PORT"),
		Host:         os.Getenv("SERVER_HOST"),
		WriteTimeout: time.Second * getDuration("SERVER_WRITE_TIMEOUT"),
		ReadTimeout:  time.Second * getDuration("SERVER_READ_TIMEOUT"),
	}

	return config, nil
}

func getInt(s string) int {
	str := os.Getenv(s)
	i, _ := strconv.Atoi(str)
	return i
}

func getDuration(s string) time.Duration {
	str := os.Getenv(s)
	d, _ := time.ParseDuration(str)
	return d
}
