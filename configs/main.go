package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}

func GetToken() string {
	token := os.Getenv("BOT_TOKEN")
	fmt.Println("Token:", token)
	return token
}
