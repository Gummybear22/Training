package controllers

import (
	"Template/pkg/models"
	"Template/pkg/utils/go-utils/database"

	"github.com/gofiber/fiber/v2"
)

func ListAccounts(c *fiber.Ctx) error {
	accounts := []models.Accounts{}

	database.DBConn.Table("accounts").Find(&accounts)

	return c.Render("index", fiber.Map{
		"Title":    "Register",
		"Subtitle": "",
		"Accounts": accounts,
	})
}

func NewAccountView(c *fiber.Ctx) error {
	return c.Render("new", fiber.Map{
		"Title":    "New Fact",
		"Subtitle": "Add something new",
	})
}

func CreateAccount(c *fiber.Ctx) error {
	accounts := new(models.Accounts)
	if err := c.BodyParser(accounts); err != nil {
		return NewAccountView(c)
	}

	result := database.DBConn.Create(&accounts)
	if result.Error != nil {
		return NewAccountView(c)
	}

	return ListAccounts(c)
}
