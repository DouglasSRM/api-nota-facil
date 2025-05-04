package main

import (
	"api-nota-facil/handlers"
	"api-nota-facil/models"
	"api-nota-facil/routes"
	"api-nota-facil/storage"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host: 		os.Getenv("DB_HOST"),
		Port: 		os.Getenv("DB_PORT"),
		User: 		os.Getenv("DB_USER"),
		Password: 	os.Getenv("DB_PASS"),
		SSLMode: 	os.Getenv("DB_SSLMODE"),
		DBName: 	os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)
	if (err != nil) {
		log.Fatal("could not load the database")
	}

	err = models.MigrateNotes(db)
	if err != nil {
		log.Fatal("could not migrade db")
	}

	app := fiber.New()
	app.Use(cors.New())

	r := handlers.NewRepository(db)
	routes.SetupRoutes(app, r)

	log.Fatal(app.Listen(":8080"))
}