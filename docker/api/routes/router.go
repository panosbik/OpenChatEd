package routes

import (
	"github.com/gofiber/fiber/v2"

	"OpenChatEd/controllers"
)

func APIRoutes(app *fiber.App, baseController controllers.BaseController) {
	// If you want to forward with a specific domain. You have to use proxy.DomainForward.
	userRoutes(app, baseController)
}
