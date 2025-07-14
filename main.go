package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Dhruva430/dis-scrim.git/bot"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env: %v", err)
	}

	bot.StartBot(os.Getenv("TOKEN"))
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop
	log.Println("Shutting down bot...")
}
