package controllers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"OpenChatEd/initializers"
)

// Controller interface defines the common behavior for all controllers
type Controller interface {
	Index(c *fiber.Ctx) error
	Show(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

// BaseController implements common behavior for all controllers
type BaseController struct {
	DB     *gorm.DB
	Config *initializers.Config
}

func NewBaseController(db *gorm.DB, config *initializers.Config) *BaseController {
	return &BaseController{
		DB:     db,
		Config: config,
	}
}

func (bc *BaseController) Index(c *fiber.Ctx) error {
	// Implementation here
	return nil
}

func (bc *BaseController) Show(c *fiber.Ctx) error {
	// Implementation here
	return nil
}

func (bc *BaseController) Create(c *fiber.Ctx) error {
	// Implementation here
	return nil
}

func (bc *BaseController) Update(c *fiber.Ctx) error {
	// Implementation here
	return nil
}

func (bc *BaseController) Delete(c *fiber.Ctx) error {
	// Implementation here
	return nil
}
