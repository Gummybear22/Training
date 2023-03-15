package routers

import (
	"Template/pkg/controllers"
	"Template/pkg/controllers/healthchecks"

	"github.com/gofiber/fiber/v2"
)

func SetupPublicRoutes(app *fiber.App) {
	// TEST ROUTES
	testRoutes := app.Group("/test")
	testRoutes.Post("/reg", controllers.RegisterSample)
	testRoutes.Post("/ver", controllers.LoginAuth)
	testRoutes.Post("/update", controllers.UpdateAccount)
	testRoutes.Get("/accounts", controllers.ListAccounts)

	// Endpoints
	apiEndpoint := app.Group("/api")
	publicEndpoint := apiEndpoint.Group("/public")
	v1Endpoint := publicEndpoint.Group("/v1")

	// Service health check
	v1Endpoint.Get("/", healthchecks.CheckServiceHealth)
}

func SetupPublicRoutesB(app *fiber.App) {

	// Endpoints
	apiEndpoint := app.Group("/api")
	publicEndpoint := apiEndpoint.Group("/public")
	v1Endpoint := publicEndpoint.Group("/v1")

	// Service health check
	v1Endpoint.Get("/", healthchecks.CheckServiceHealthB)
}
