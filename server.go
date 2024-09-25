package main

import (
	"goSSR/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables from
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// template engine
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:        engine,
		ErrorHandler: CustomeErrorHandler,
	})

	// set up routes
	routes.Setup(app)

	log.Fatal(app.Listen(":3000"))
}
