package conf

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
	"strconv"
)

type Config struct {
	AppPort    int
	DbHost     string
	DbPort     int
	DbUser     string
	DbPassword string
	DbName     string
	JwtSecret  string
}

var conf *Config

func LoadConfig() *Config {
	logger := InitLogger()
	defer logger.Sync()

	if err := godotenv.Load("../.env"); err != nil {
		logger.Fatal("Ошибка загрузки .env файла", zap.Error(err))
	}

	appPort, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		logger.Fatal("Ошибка конвертации APP_PORT", zap.Error(err))
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		logger.Fatal("Ошибка конвертации DB_PORT", zap.Error(err))
	}

	return &Config{
		AppPort:    appPort,
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     dbPort,
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbName:     os.Getenv("DB_NAME"),
		JwtSecret:  os.Getenv("JWT_SECRET"),
	}

}
