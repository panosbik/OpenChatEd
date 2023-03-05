package routes

import (
	"github.com/gofiber/fiber/v2"

	"OpenChatEd/controllers"
	"OpenChatEd/middleware"
)

func userRoutes(app *fiber.App, base controllers.BaseController) {
	userController := controllers.UserController{BaseController: base}
	app.Get("/confirm-email", userController.ConfirmEmail)
	app.Post("/auth.register", userController.Create)
	app.Post("/auth.login", userController.Login)
	app.Get("/users/me", middleware.JWTAuthorization(base.Config.JWT, base.DB), userController.Me)
	app.Get("/users/search", middleware.JWTAuthorization(base.Config.JWT, base.DB), userController.Search)
}
