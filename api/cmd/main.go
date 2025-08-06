package main

import (
	"api/internal/handler"
	"api/internal/service"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	appService := service.NewAppService()
	appHandler := handler.NewAppHandler(appService)

	appHandler.RegisterRoutes(app)

	log.Println("🚀 DevOpsControl API is running on http://localhost:3000")
	app.Static("/", "./web")
	log.Fatal(app.Listen(":3000"))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./web/index.html")
	})

	app.Get("/apps", appHandler.ListApps)
	app.Delete("/:name", appHandler.DeleteApp)
}
