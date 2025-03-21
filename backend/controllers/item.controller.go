package controllers

import (
	"github.com/KennethTrinh/go-vite-app/initializers"
	"github.com/KennethTrinh/go-vite-app/models"
	"github.com/KennethTrinh/go-vite-app/utils"
	"github.com/gofiber/fiber/v3"
)

// curl -X POST http://localhost:8000/items -H "Content-Type: application/json" -d '{"name":"test","description":"test","icon":"❤️","color":"test","time":1}'
func CreateItem(c fiber.Ctx) error {
	payload := &models.Item{}
	if err := c.Bind().Body(payload); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, err, "Invalid request payload")
	}

	validateErr := models.ValidateStruct(payload)
	if validateErr != nil {
		return utils.Error(c, fiber.StatusBadRequest, validateErr, "Error validating data")
	}

	// initializers.DB.Create(&payload)
	if err := initializers.DB.Create(&payload).Error; err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, err, "Failed to create item")
	}

	return c.SendStatus(fiber.StatusCreated)

}

// curl http://localhost:8000/items
func ListItems(c fiber.Ctx) error {

	var items []models.Item
	initializers.DB.Find(&items)

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"items": items,
		},
	)

}

// curl -X DELETE http://localhost:8000/items
func DeleteItems(c fiber.Ctx) error {
	// Delete all items from the database
	result := initializers.DB.Exec("DELETE FROM items")

	// Check if there was an error during deletion
	if result.Error != nil {
		return utils.Error(c, fiber.StatusInternalServerError,
			result.Error,
			"Failed to delete all items")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
