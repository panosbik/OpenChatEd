package initializers

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"

	"OpenChatEd/helpers/security"
)

type Config struct {
	DBHost         string `mapstructure:"MARIADB_HOST"`
	DBUserName     string `mapstructure:"MARIADB_USER"`
	DBUserPassword string `mapstructure:"MARIADB_PASSWORD"`
	DBName         string `mapstructure:"MARIADB_DATABASE"`
	DBPort         string `mapstructure:"MARIADB_PORT"`

	TokenExpiresIn time.Duration `mapstructure:"TOKEN_EXPIRED_IN"`
	TokenMaxAge    int           `mapstructure:"TOKEN_MAXAGE"`
	EmailAddress   string        `mapstructure:"EMAIL_ADDRESS"`
	EmailPassword  string        `mapstructure:"EMAIL_PASSWORD"`
	EmailHost      string        `mapstructure:"EMAIL_HOST"`
	EmailPort      int           `mapstructure:"EMAIL_PORT"`
	ServerUrl      string        `mapstructure:"SERVER_URL"`
	JWT            security.JWT

	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort string `mapstructure:"REDIS_PORT"`
	RedisPass string `mapstructure:"REDIS_PASS"`

	Dialer *gomail.Dialer
}

func LoadConfig(path string) (config Config) {
	viper.AddConfigPath(path)

	// To look specific file name with config: app
	viper.SetConfigName("app")

	// Set config Type that may be env, json and soon other
	viper.SetConfigType("env")

	// To read env values from environment variables
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("Could not load environment variables, \n %s", err))
	}

	// to unmarshal values into target object
	if err = viper.Unmarshal(&config); err != nil {
		panic(fmt.Sprintf("Could not load environment variables, \n %s", err))
	}

	config.Dialer = gomail.NewDialer(config.EmailHost, config.EmailPort, config.EmailAddress, config.EmailPassword)

	prvKey, err := os.ReadFile("./auth.cert/id_rsa")
	if err != nil {
		panic(err)
	}

	pubKey, err := os.ReadFile("./auth.cert/id_rsa.pub")
	if err != nil {
		panic(err)
	}
	config.JWT = security.NewJWT(prvKey, pubKey)

	log.Println("Successfully initialize environment variables")
	return
}
