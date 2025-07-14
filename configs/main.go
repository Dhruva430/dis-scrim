package configs

import "os"

func GetToken() string {
	return os.Getenv("TOKEN")
}
