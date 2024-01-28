package main

import (
	"github.com/GDGVIT/configgy-backend/cmd"
	"github.com/joho/godotenv"
)

// Message represents the message structure you expect to send to the RabbitMQ queue.

func main() {
	cmd.Execute()
	godotenv.Load()
}
