package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost string
	Port       string

	DBUser     string
	DBPassword string
	DBHost string
	DBPort string
	DBName     string

	JWTExpirationInSeconds int64
	JWTSecret string
}

// create a singleton
var Envs = initConfig()

func initConfig() Config {
	godotenv.Load() //load environment variables from .env file
	return Config{
		PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
		Port: getEnv("PORT", "8080"),
		DBUser: getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "1234"),
		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "5433"),
		DBName: getEnv("DB_NAME", "personal_blog"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXP", 3600*24*7),
		JWTSecret: getEnv("JWT_SECRET", "secret-not-secret-anymore?"), // for signature and checking JWT
	}
}

func getEnv(key, fallback string) string {
	//look for environment variable by key
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback // default
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback //default
		}
		return i // converted value
	}
	return fallback //default
}