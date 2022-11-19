package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type DbConfig struct {
	ConnectionString string
	DB_HOST          string
	DB_PORT          string
	DB_USERNAME      string
	DB_PASSWORD      string
	DB_NAME          string
}

type AppConfig struct {
	SECRET              string
	JWT_EXPIRATION_TIME string
	MONGO_URL           string
	REDIS_URL           string
	REDIS_PASSWORD      string
}

func GetDbCfg() (DbConfig, error) {
	cfg := DbConfig{}
	err := readEnvironmentVariablesForDb(&cfg)
	if err != nil {
		return DbConfig{}, err
	}
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB_HOST,
		cfg.DB_PORT,
		cfg.DB_USERNAME,
		cfg.DB_PASSWORD,
		cfg.DB_NAME)
	cfg.ConnectionString = connectionString
	return cfg, nil
}
func GetAppCfg() (AppConfig, error) {
	cfg := AppConfig{}
	err := readEnvironmentVariablesForApp(&cfg)
	if err != nil {
		return AppConfig{}, err
	}

	return cfg, nil
}
func readEnvironmentVariablesForDb(cfg *DbConfig) (err error) {
	err = godotenv.Load(".env")
	if err != nil {
		return err
	}
	vars := []string{"DB_HOST",
		"DB_PORT",
		"DB_USERNAME",
		"DB_PASSWORD",
		"DB_NAME",
		"SECRET",
		"JWT_EXPIRATION_TIME"}

	for i := range vars {
		if os.Getenv(vars[i]) == "" {
			return errors.New("empty environment variable " + vars[i])
		}
	}
	cfg.DB_HOST = os.Getenv("DB_HOST")
	cfg.DB_PORT = os.Getenv("DB_PORT")
	cfg.DB_USERNAME = os.Getenv("DB_USERNAME")
	cfg.DB_PASSWORD = os.Getenv("DB_PASSWORD")
	cfg.DB_NAME = os.Getenv("DB_NAME")

	return nil
}

func readEnvironmentVariablesForApp(cfg *AppConfig) (err error) {
	err = godotenv.Load(".env")
	if err != nil {
		return err
	}
	vars := []string{
		"SECRET",
		"MONGO_URL",
		"REDIS_URL",
		"RedisPassword",
		"JWT_EXPIRATION_TIME"}

	for i := range vars {
		if os.Getenv(vars[i]) == "" {
			return errors.New("empty environment variable " + vars[i])
		}
	}

	cfg.SECRET = os.Getenv("SECRET")
	cfg.JWT_EXPIRATION_TIME = os.Getenv("JWT_EXPIRATION_TIME")
	cfg.REDIS_URL = os.Getenv("REDIS_URL")
	cfg.REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")
	cfg.MONGO_URL = os.Getenv("MONGO_URL")
	_, err = time.ParseDuration(cfg.JWT_EXPIRATION_TIME)
	if err != nil {
		return err
	}
	return nil
}
