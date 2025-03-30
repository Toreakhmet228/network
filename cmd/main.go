package main

import (
	"chat/internal/database"
	"chat/internal/midllwere"
	"chat/internal/routes"
	"chat/pkg/conf"
	"github.com/gofiber/fiber/v2"
)

func main() {
	database.InitDB()
	logger := conf.InitLogger()
	app := fiber.New()
	app.Use(middleware.AuthMiddleware())
	api := app.Group("/api")
	routes.ProtectedRoutes(api)
	logger.Info("started ")
	app.Listen(":8888")

}
