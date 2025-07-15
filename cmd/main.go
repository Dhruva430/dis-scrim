package main

import (
	"dis-scrim/bot"
	"dis-scrim/configs"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	bot.StartBot(configs.GetToken())
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop
	log.Println("Shutting down bot...")
}
