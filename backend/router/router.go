package router

import (
	"github.com/KennethTrinh/go-vite-app/controllers"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App) {

	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	items := app.Group("/items")
	items.Get("/", controllers.ListItems)
	items.Post("/", controllers.CreateItem)
	items.Delete("/", controllers.DeleteItems)

}
