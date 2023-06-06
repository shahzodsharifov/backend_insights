package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wpcodevo/golang-fiber/initializers"
	"github.com/wpcodevo/golang-fiber/models"
	"gorm.io/gorm"
)

func AddVaccancy(c *fiber.Ctx) error {
	userID := c.Params("userID")

	var user models.User
	if err := initializers.DB.First(&user, "id = ?", userID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving user",
		})
	}

	var vaccancy models.Vaccancy
	if err := c.BodyParser(&vaccancy); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing vaccancy data",
		})
	}

	vaccancy.EmployerID = user.ID

	if err := initializers.DB.Create(&vaccancy).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating vaccancy",
		})
	}

	user.Vaccancies = append(user.Vaccancies, vaccancy)

	return c.JSON(fiber.Map{"vaccancy": vaccancy})

}

func GetUserVaccancies(c *fiber.Ctx) error {

	userID := c.Params("userID")

	var user models.User

	if err := initializers.DB.Where("id = ?", userID).Preload("Vaccancies").First(&user).Error ; err !=nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Vaccancy not found",
		})
	}

	return c.JSON(user.Vaccancies)
}

func GetAllVacancies(c *fiber.Ctx) error {
	var vaccancies *[]models.Vaccancy

	res:= initializers.DB.Limit(20).Find(&vaccancies)
	if err := res.Error; err !=nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status":"success", "data": fiber.Map{"vaccancies": vaccancies}})
}

func GetVaccancyByID(c *fiber.Ctx) error {
	// userID := c.Params("userID")
	vaccancyID := c.Params("vaccancyID")

	var vaccancy models.Vaccancy
	res := initializers.DB.Find(&vaccancy, "id =?", vaccancyID)
	if err := res.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": err.Error()})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})

	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"vaccancy": vaccancy}})

}
