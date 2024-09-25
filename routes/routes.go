package routes

import (
	"goSSR/handlers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Get("/", handlers.HandleIndex)
	app.Get("/about", handlers.HandleAbout)
	app.Post("/upload", handlers.HandleUpload)
}
