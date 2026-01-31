package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type ServerConfig struct {
	Port int `validate:"number"`
}

func GetServerConfigFromEnv() ServerConfig {
	portStr := os.Getenv("PORT")
	if portStr == "" {
		log.Fatalln("Missing 'PORT' env")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid 'PORT' env got %s expected int", portStr)
	}

	return ServerConfig{port}
}

type DatabaseConfig struct {
	Database string
	Username string
	Password string
	Host     string
	Port     string
	Schema   string
	SSLMode  string
}

func GetDatabaseConfigFromEnv() DatabaseConfig {
	database := os.Getenv("DB_DATABASE")
	if database == "" {
		log.Fatalln("Missing 'DB_DATABASE' env")
	}

	usr := os.Getenv("DB_USERNAME")
	if usr == "" {
		log.Fatalln("Missing 'DB_USERNAME' env")
	}

	passpwd := os.Getenv("DB_PASSWORD")
	if passpwd == "" {
		log.Fatalln("Missing 'DB_PASSWORD' env")
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		log.Fatalln("Missing 'DB_HOST' env")
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		log.Fatalln("Missing 'DB_PORT' env")
	}

	schema := os.Getenv("DB_SCHEMA")
	if schema == "" {
		schema = "public"
	}

	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	return DatabaseConfig{database, usr, passpwd, host, port, schema, sslmode}
}

func (d DatabaseConfig) ConnStr() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&search_path=%s", d.Username, d.Password, d.Host, d.Port, d.Database, d.SSLMode, d.Schema)
}
