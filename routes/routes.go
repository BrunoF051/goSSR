package routes

import (
	"goSSR/auth"
	"goSSR/handlers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB) {
	h := handlers.NewHandler(db)
	authHandler := auth.NewHandler(db)

	app.Get("/", h.HandleIndex)
	app.Get("/about", h.HandleAbout)
	app.Post("/upload", auth.RequireAuth, h.HandleUpload)
	// Google OAuth routes
	app.Get("/auth/google", authHandler.GoogleLoginHandler)
	app.Get("/auth/google/callback", authHandler.GoogleCallbackHandler)

}
