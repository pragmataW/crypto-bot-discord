package main

import (
	"log"
	"os"
	"rrs/controller"
	"rrs/service"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var (
	url string
)

func main() {
	src := service.NewRSSService(url)
	ctrl := controller.NewController(src)

	app := fiber.New()
	controller.SetUpRoutes(app, ctrl)
	log.Fatal(app.Listen(":8080"))
}

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(err)
	}

	url = os.Getenv("URL")
}
