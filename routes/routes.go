package routes

import (
	"goSSR/handlers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB) {
	h := handlers.NewHandler(db)

	app.Get("/", h.HandleIndex)
	app.Get("/about", h.HandleAbout)
	app.Post("/upload", h.HandleUpload)
}
