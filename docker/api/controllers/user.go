package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"path"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"OpenChatEd/constants"
	"OpenChatEd/helpers"
	"OpenChatEd/helpers/pagination"
	"OpenChatEd/helpers/redis"
	"OpenChatEd/helpers/security"
	"OpenChatEd/helpers/validation"
	"OpenChatEd/models"
	"OpenChatEd/schemas"
)

// UserController
type UserController struct {
	BaseController
}

func (uc *UserController) Index(c *fiber.Ctx) error {
	// Implementation here
	return nil
}

func (uc *UserController) Show(c *fiber.Ctx) error {
	// Implementation here
	return nil
}

func (uc *UserController) Create(c *fiber.Ctx) error {
	var payload schemas.SignUpInput

	if err := c.BodyParser(&payload); err != nil {
		// Return an error response if the request body could not be parsed
		return c.Status(fiber.StatusUnprocessableEntity).JSON(schemas.NewResponse(nil, constants.ErrBodyParse.Error()))
	}

	// Validate the request body against the defined schema
	if errs := validation.ValidateStruct(payload); errs != nil {
		// Return an error response if the request body fails validation
		return c.Status(fiber.StatusUnprocessableEntity).JSON(schemas.NewResponse(nil, errs))
	}

	newUser := models.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	}

	// Create the new user in database
	if err := uc.DB.
		Create(&newUser).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(schemas.NewResponse(nil, err.Error()))
	}

	// Send Email verification
	token, _, _ := uc.Config.JWT.EncodeJWToken(time.Duration(30*time.Minute), newUser.ID)
	go helpers.SendEmail(
		uc.Config.Dialer,
		"Email Verification",
		"./templates/email-verification.html",
		struct{ URL string }{
			URL: path.Join(
				uc.Config.ServerUrl,
				fmt.Sprintf("confirm-email?token=%s", token),
			),
		},
		[]string{newUser.Email},
	)
	// Return the new user struct
	return c.Status(fiber.StatusCreated).JSON(schemas.NewResponse(newUser, nil))
}

func (uc *UserController) Update(c *fiber.Ctx) error {
	// Implementation here
	return nil
}

func (uc *UserController) Delete(c *fiber.Ctx) error {
	// Implementation here
	return nil
}

// TODO: Render template for visual feedback
func (uc *UserController) ConfirmEmail(c *fiber.Ctx) error {
	jwtString := c.Query("token")

	if uuid, err := uc.Config.JWT.DecodeJWToken(jwtString); err == nil {

		var user models.User
		if err := uc.DB.First(&user, *uuid).Updates(models.User{IsActive: true}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusForbidden).JSON(schemas.NewResponse(nil, constants.ErrAccessDenied.Error()))
			} else {
				return c.Status(fiber.StatusInternalServerError).JSON(schemas.NewResponse(nil, constants.ErrServer.Error()))
			}
		}
		return nil
	}
	return c.Status(fiber.StatusForbidden).JSON(schemas.NewResponse(nil, constants.ErrInvalidToken.Error()))
}

// Login user
func (uc *UserController) Login(c *fiber.Ctx) error {
	// Validate the input body for create the new user
	payload := new(schemas.SignInInput)

	if err := c.BodyParser(&payload); err != nil {
		// Return an error response if the request body could not be parsed
		return c.Status(fiber.StatusUnprocessableEntity).JSON(schemas.NewResponse(nil, constants.ErrBodyParse.Error()))
	}

	// Validate the request body against the defined schema
	if errs := validation.ValidateStruct(payload); errs != nil {
		// Return an error response if the request body fails validation
		return c.Status(fiber.StatusUnprocessableEntity).JSON(schemas.NewResponse(nil, errs))
	}

	var (
		userID uint
		err    error
	)

	if payload.GrantType == "password" {
		user := new(models.User)
		// Retrieve user from the database by email address
		if err = uc.DB.First(&user, "email = ? AND is_active = ?", payload.Email, true).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusBadRequest).JSON(schemas.NewResponse(nil, constants.ErrInvalidLoginCredentials.Error()))
			}
			return c.Status(fiber.StatusInternalServerError).JSON(schemas.NewResponse(nil, constants.ErrServer.Error()))
		}

		// Compare send in pass with saved user pass hash
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(schemas.NewResponse(nil, constants.ErrInvalidLoginCredentials.Error()))
		}

		// Save the user ID save
		userID = user.ID

	} else {
		// Get the user ID for redis cache by refresh token and save it
		if userID, err = redis.GetUserIDByRefreshToken(payload.RefreshToken); err != nil {
			return c.Status(fiber.StatusForbidden).JSON(schemas.NewResponse(nil, constants.ErrInvalidRefreshToken.Error()))
		}
	}

	// Generate a jwt token
	token, exp, err := uc.Config.JWT.EncodeJWToken(uc.Config.TokenExpiresIn, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(schemas.NewResponse(nil, constants.ErrServer.Error()))
	}

	// Generate a refresh Token
	refreshToken, err := security.GenerateRefreshToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(schemas.NewResponse(nil, constants.ErrServer.Error()))
	}

	// Store the refresh Token to the redis cache
	if err = redis.SaveRefreshToken(refreshToken, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(schemas.NewResponse(nil, constants.ErrServer.Error()))
	}

	// Return the jwt token
	return c.Status(fiber.StatusOK).JSON(schemas.NewResponse(schemas.NewToken(token, refreshToken, exp), nil))
}

// Return the struct of logon user
func (uc *UserController) Me(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(schemas.NewResponse(c.Locals(constants.CurrentUser), nil))
}

// Searching for a user
func (uc *UserController) Search(c *fiber.Ctx) error {
	// Validate the required query parameter
	term := c.Query("term")
	if term == "" {
		return c.Status(fiber.StatusBadRequest).JSON(schemas.NewResponse(nil, constants.ErrInvalidTermSearchUser.Error()))
	}

	users := new([]models.User)

	// Get pagination parameters from query string
	page := c.QueryInt("page")
	pageSize := c.QueryInt("page_size")

	// Search for a users by term
	q := uc.DB.
		Model(&models.User{}).
		Omit("password", "is_active", "icon_path", "created", "updated", "deleted").
		Where("is_active = ?", true).
		Where("username LIKE @term OR email LIKE @term", sql.Named("term", term+"%"))

	// Paginate the results using the provided page and page size
	paginationResults, err := pagination.NewPagingResult(q, users, page, pageSize)

	if err != nil {
		// Return a server error response if an error occurred during pagination
		return c.Status(fiber.StatusInternalServerError).JSON(schemas.NewResponse(nil, constants.ErrServer.Error()))
	}

	// Return list of users struct
	return c.Status(fiber.StatusOK).JSON(schemas.NewResponse(paginationResults, nil))
}
