package main

import (
	"goSSR/database"
	"goSSR/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	postgresStorage "github.com/gofiber/storage/postgres/v3"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	// Load environment variables from
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize storage config
	storage := postgresStorage.New(database.ConfigStorage)

	// Close the storage when the program terminates
	defer storage.Close()

	// Initialize database
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// migrate the schemas
	db.AutoMigrate(database.GetModels()...)

	// template engine
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:        engine,
		ErrorHandler: CustomeErrorHandler,
	})

	// set up routes
	routes.Setup(app, db)

	log.Fatal(app.Listen(":3000"))
}
