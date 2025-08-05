package handler

import (
	"api/internal/model"
	"api/internal/service"
	"api/internal/utils"

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
	api.Get("/apps", h.ListApps)
	api.Get("/:id", h.GetApp)
}

func (h *AppHandler) CreateApp(c *fiber.Ctx) error {
	var req model.App

	// convert json format to input struct
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}

	// create application in memory
	app := h.AppService.CreateApp(req)

	err := utils.GenerateHelmFiles(req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate helm files"+err.Error())
	}

	err = utils.CopyTemplates(req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to copy template files"+err.Error())
	}

	// generate manifest using helm
	err = service.GenerateManifest(req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate manifest: "+err.Error())
	}

	err = service.InstallHelmRelease(req)
	if err != nil {
		return c.Status(500).SendString("Helm install failed: " + err.Error())
	}

	// Return manifest and app together
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"app": app,
	})
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
	apps, err := service.ListDeployedApps()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(apps)
}
