package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wpcodevo/golang-fiber/initializers"
	"github.com/wpcodevo/golang-fiber/models"
	"gorm.io/gorm"
)

func AddEvent(c *fiber.Ctx) error {
	userID := c.Params("userID")

	var user models.User
	if err := initializers.DB.First(&user, "id = ?", userID).Error; err !=nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving user",
		})
	} 

	var event models.Event
	if err := c.BodyParser(&event); err !=nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing event data",
		})
	}

	event.OrganizerID = user.ID

	if err := initializers.DB.Create(&event).Error; err !=nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating event",
		})
	}

	user.Events = append(user.Events, event)

	return c.JSON(fiber.Map{"event": event})
}

func GetUserEvents(c *fiber.Ctx) error {

	userID := c.Params("userID")

	var user models.User
	if err := initializers.DB.Where("id = ?", userID).Preload("Events").First(&user).Error; err !=nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Events not found",
		})

	}

	return c.JSON(user.Events)
}

func GetEventsByID(c *fiber.Ctx) error {
	// userID := c.Params("userID")
	eventID := c.Params("eventID")

	var event models.Event
	res :=initializers.DB.First(&event, "id =?", eventID)
	if err := res.Error; err !=nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": err.Error()})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"event": event}})

}

func GetAllEvents(c *fiber.Ctx) error {
	var events *[]models.Event

	res := initializers.DB.Limit(10).Find(&events)
	if err := res.Error; err !=nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"events": events}})


}