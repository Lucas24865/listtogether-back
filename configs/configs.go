package configs

import (
	"github.com/joho/godotenv"
)

func SetConfigs() {
	err := godotenv.Load(".env")

	if err != nil {
		panic("Error loading .env file")
	}
}
