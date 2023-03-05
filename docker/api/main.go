package main

import (
	"fmt"
	"log"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	goRedis "github.com/redis/go-redis/v9"

	"OpenChatEd/controllers"
	"OpenChatEd/helpers/redis"
	"OpenChatEd/initializers"
	"OpenChatEd/routes"
)

var BaseController *controllers.BaseController

func init() {
	config := initializers.LoadConfig(".")
	db := initializers.ConnectDB(&config)
	BaseController = controllers.NewBaseController(db, &config)
	redis.Client = goRedis.NewClient(&goRedis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		Password: config.RedisPass,
		DB:       0, // use default DB
	})
}

func main() {
	app := fiber.New(fiber.Config{
		Prefork:      true,
		ServerHeader: "",
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
	})
	routes.APIRoutes(app, *BaseController)
	log.Fatal(app.Listen(":3000"))
}
