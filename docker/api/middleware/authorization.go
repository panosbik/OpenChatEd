package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"OpenChatEd/constants"
	"OpenChatEd/helpers/security"
	"OpenChatEd/models"
	"OpenChatEd/schemas"
)

// Authorization middleware which verify the JWT on Headers
// and saves the logon user in Local pointer
func JWTAuthorization(j security.JWT, db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		jwtString := c.Get("Authorization", "")
		if jwtString == "" || !strings.HasPrefix(jwtString, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(schemas.NewResponse(nil, constants.ErrAccessDenied.Error()))
		}

		// Remove the schema type from token
		jwtString = strings.TrimPrefix(jwtString, "Bearer ")

		// Parse the jwt string
		if uuid, err := j.DecodeJWToken(jwtString); err == nil {
			var user models.User
			// Find the user form sub claim
			if err := user.First(db, *uuid); err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return c.Status(fiber.StatusForbidden).JSON(schemas.NewResponse(nil, constants.ErrAccessDenied.Error()))
				} else {
					return c.Status(fiber.StatusInternalServerError).JSON(schemas.NewResponse(nil, constants.ErrServer.Error()))
				}
			}

			// Check if user is active
			if user.IsActive {
				return c.Status(fiber.StatusUnauthorized).JSON(
					schemas.NewResponse(nil, constants.ErrEmailVerificationError.Error()))
			}

			// Save the user to Locals pointer
			c.Locals(constants.CurrentUser, user)

			return c.Next()
		} else {
			return c.Status(fiber.StatusForbidden).JSON(schemas.NewResponse(nil, err.Error()))
		}
	}
}
