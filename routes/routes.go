package routes

import (
	"api-nota-facil/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, h *handlers.Repository) {
	api := app.Group("/api")
	api.Post("/create_note", h.CreateNote)
	api.Delete("/delete_notes", h.DeleteNotes)
	api.Get("/notes", h.GetNotes)
	api.Put("/update_note/:id", h.UpdateNote)
}
