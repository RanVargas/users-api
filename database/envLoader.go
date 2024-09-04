package database

import "os"

/**
  POSTGRES_DB: pgdb
  POSTGRES_PASSWORD: my-password
  POSTGRES_PORT: "5432"
  JWT_SECRET: "]H^IzX.-#V),yMrqdo7-B'YSKo`\z7fC"
DB_HOST=SomeHost
DB_PORT=SomePort
DB_USER=SomeUserName
DB_PASS=SomePassword
DB_NAME=SomeDbName

*/

type EnvVars struct {
	DatabaseUsername string
	DatabasePassword string
	DatabaseHost     string
	DatabaseName     string
	DatabasePort     string
	JwtSecret        string
	XApiKey          string
	DatabaseSslMode  string
	DatabaseTimeZone string
}

func LoadEnvVariables() EnvVars {
	value := EnvVars{
		DatabaseUsername: os.Getenv("POSTGRES_USERNAME"),
		DatabasePassword: os.Getenv("POSTGRES_PASSWORD"),
		DatabaseHost:     os.Getenv("POSTGRES_HOST"),
		DatabaseName:     os.Getenv("POSTGRES_DB"),
		DatabasePort:     os.Getenv("POSTGRES_PORT"),
		JwtSecret:        os.Getenv("JWT_SECRET"),
		XApiKey:          os.Getenv("X-API-KEY"),
		DatabaseSslMode:  os.Getenv("POSTGRES_SSL_MODE"),
		DatabaseTimeZone: os.Getenv("POSTGRES_TIMEZONE"),
	}
	if value.DatabaseUsername == "" || value.DatabasePassword == "" ||
		value.DatabaseHost == "" || value.DatabaseName == "" ||
		value.DatabasePort == "" || value.JwtSecret == "" ||
		value.XApiKey == "" || value.DatabaseSslMode == "" ||
		value.DatabaseTimeZone == "" {
		panic("failed to load essential environment variables")
	}
	return value
}
