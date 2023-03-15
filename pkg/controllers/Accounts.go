package controllers

import (
	"Template/pkg/models"
	"Template/pkg/models/errors"
	"Template/pkg/utils/go-utils/database"
	"Template/pkg/utils/go-utils/passwordHashing"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/go-playground/validator.v9"
)

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

var validate = validator.New()

func ValidateStruct(user models.Accounts) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func CheckIfExixst(accounts models.Accounts, contacts models.Contacts) []*errors.ErrorModel {
	// declare exists error
	var exisErr []*errors.ErrorModel

	// query to check if username already exist
	var checker bool
	err := database.DBConn.Raw("SELECT EXISTS(SELECT 1 FROM accounts WHERE username = $1 OR email = $1 OR contact = $1)", accounts.Username).Row().Scan(&checker)
	if err != nil {
		return exisErr
	}

	// check if username exist
	if checker {
		var message errors.ErrorModel
		message.Message = "Username Already Exist"
		message.IsSuccess = false
		message.Error = err
		exisErr = append(exisErr, &message)
	}

	// query to check if eamil already exist
	err = database.DBConn.Raw("SELECT EXISTS(SELECT 1 FROM accounts WHERE email = $1)", contacts.Email).Row().Scan(&checker)
	if err != nil {
		var message errors.ErrorModel
		message.Message = "Query Error"
		message.IsSuccess = false
		message.Error = err
		exisErr = append(exisErr, &message)
	}

	// check if email exist
	if checker {
		var message errors.ErrorModel
		message.Message = "Email Already Exist"
		message.IsSuccess = false
		message.Error = err
		exisErr = append(exisErr, &message)
	}

	// query to check if contact already exist
	err = database.DBConn.Raw("SELECT EXISTS(SELECT 1 FROM accounts WHERE contact = $1)", contacts.Contact).Row().Scan(&checker)
	if err != nil {
		var message errors.ErrorModel
		message.Message = "Query Error"
		message.IsSuccess = false
		message.Error = err
		exisErr = append(exisErr, &message)
	}

	// check if contact exist
	if checker {
		var message errors.ErrorModel
		message.Message = "Contact Already Exist"
		message.IsSuccess = false
		message.Error = err
		exisErr = append(exisErr, &message)
	}

	return exisErr
}

// controller for registration
func RegisterSample(c *fiber.Ctx) error {
	// point to models
	accounts := &models.Accounts{}
	contacts := &models.Contacts{}

	// body parser, parses data submitted
	if parsErr := c.BodyParser(accounts); parsErr != nil {
		return c.JSON(fiber.Map{
			"Error in parsing:": parsErr.Error(),
		})
	}

	// body parser, parses data submitted
	if parsErr := c.BodyParser(contacts); parsErr != nil {
		return c.JSON(fiber.Map{
			"Error in parsing:": parsErr.Error(),
		})
	}

	// validate data
	error := ValidateStruct(*accounts)
	if error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(error)

	}

	// check if exist already
	errorMessage := CheckIfExixst(*accounts, *contacts)
	if errorMessage != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorMessage)
	}

	// hash incoming password
	hashedPassword, err := passwordHashing.HashPassword(accounts.Password)
	if err != nil {
		return err
	}

	// Insert new user with hashed password
	var lastId int
	err = database.DBConn.Raw("INSERT INTO accounts (first_name, last_name, username, password) VALUES (?, ?, ?, ?) RETURNING account_id",
		accounts.First_name, accounts.Last_name, accounts.Username, hashedPassword).Scan(&lastId).Error
	if err != nil {
		return err
	}

	// Insert new email and contact
	err = database.DBConn.Exec("INSERT INTO contacts (account_id, email, contact) VALUES (?, ?, ?)",
		lastId, contacts.Email, contacts.Contact).Error
	if err != nil {
		return err
	}

	// return message
	return c.JSON(fiber.Map{
		"result": "registration successful",
	})
}

func LoginAuth(c *fiber.Ctx) error {
	userModel := &models.Accounts{}

	if parsErr := c.BodyParser(userModel); parsErr != nil {
		return c.JSON(fiber.Map{
			"Error in parsing:": parsErr.Error(),
		})
	}

	var dbpass string

	err := database.DBConn.Raw("SELECT password FROM accounts WHERE username=$1 OR email=$1 OR contact=$1", userModel.Username).Row().Scan(&dbpass)
	if err != nil {
		fmt.Print(err)
		return c.JSON(fiber.Map{
			"query error": err,
		})
	}

	if !passwordHashing.CheckPasswordHash(userModel.Password, dbpass) {
		fmt.Print(userModel.Password, dbpass)
		return c.JSON(fiber.Map{
			"Error": "Invalid Password",
		})
	}

	return c.JSON(fiber.Map{
		"Result": "login successful",
	})
}

func UpdateAccount(c *fiber.Ctx) error {
	userModel := &models.Accounts{}

	if parsErr := c.BodyParser(userModel); parsErr != nil {
		return c.JSON(fiber.Map{
			"Error in parsing:": parsErr.Error(),
		})
	}

	// validate data
	error := ValidateStruct(*userModel)
	if error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(error)

	}

	// check if exist already
	// errorMessage := CheckIfExixst(*userModel)
	// if errorMessage != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(errorMessage)
	// }

	// hashedPassword, err := passwordHashing.HashPassword(userModel.Password)
	// if err != nil {
	// 	return err
	// }

	// err = database.DBConn.Exec("UPDATE accounts	SET first_name = ?, last_name = ?, username = ?, password = ?, email = ?, contact = ?, house_no = ?, street = ?, subdivision = ?, barangay = ?, city = ?, province = ?, country = ?, zip_code = ? WHERE id = ?",
	// 	userModel.First_name, userModel.Last_name, userModel.Username, hashedPassword, userModel.Email, userModel.Contact, userModel.House_no, userModel.Street, userModel.Subdivision, userModel.Barangay, userModel.City, userModel.Province, userModel.Country, userModel.Zip_code, userModel.Id).Error
	// if err != nil {
	// 	fmt.Print(err)
	// 	return c.JSON(fiber.Map{
	// 		"query error": err.Error(),
	// 	})
	// }

	return c.JSON(fiber.Map{
		"Result": "update successful",
	})
}

func ListAccounts(c *fiber.Ctx) error {
	accounts := []models.Accounts{}

	err := database.DBConn.Raw("SELECT * FROM accounts ORDER BY ID ASC").Find(&accounts).Error
	if err != nil {
		return c.JSON(fiber.Map{
			"Result": err,
		})
	}

	return c.JSON(fiber.Map{
		"Result": accounts,
	})
}
