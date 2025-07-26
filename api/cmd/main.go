package main

import (
	"devopscontrol-api/internal/handler"
	"devopscontrol-api/internal/service"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	appService := service.NewAppService()
	appHandler := handler.NewAppHandler(appService)

	appHandler.RegisterRoutes(app)

	log.Println("🚀 DevOpsControl API is running on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}
