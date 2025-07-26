package handler

import (
	"devopscontrol-api/internal/model"
	"devopscontrol-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AppHandler struct {
	AppService *service.AppService
}

func NewAppHandler(s *service.AppService) *AppHandler {
	return &AppHandler{AppService: s}
}

func (h *AppHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api/apps")
	api.Post("/", h.CreateApp)
	api.Get("/", h.ListApps)
	api.Get("/:id", h.GetApp)
}

func (h *AppHandler) CreateApp(c *fiber.Ctx) error {
	var req model.App
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}
	app := h.AppService.CreateApp(req)
	return c.Status(fiber.StatusCreated).JSON(app)
}

func (h *AppHandler) GetApp(c *fiber.Ctx) error {
	id := c.Params("id")
	app, exists := h.AppService.GetApp(id)
	if !exists {
		return fiber.NewError(fiber.StatusNotFound, "App not found")
	}
	return c.JSON(app)
}

func (h *AppHandler) ListApps(c *fiber.Ctx) error {
	return c.JSON(h.AppService.ListApps())
}
