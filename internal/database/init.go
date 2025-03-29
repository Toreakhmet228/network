package database

import (
	"chat/pkg/conf"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	log := conf.InitLogger()

	config := conf.LoadConfig() // Загружаем конфиг

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DbHost, config.DbPort, config.DbUser, config.DbPassword, config.DbName,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("Ошибка подключения к базе данных: ")
	}

	fmt.Println("✅ Подключение к базе данных успешно!")
}
