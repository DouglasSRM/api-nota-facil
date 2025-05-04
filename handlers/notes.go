package handlers

import (
	"fmt"
	"api-nota-facil/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) CreateNote(context *fiber.Ctx) error {
	note := models.Notes{}

	if err := context.BodyParser(&note); (err != nil) {
		return context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message":"request failed"})
	}

	if note.ID == "" {
		note.ID = uuid.New().String()
	}

	if err := r.DB.Create(&note).Error; (err != nil) {
		return context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"could not create note"})
	}

	return context.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"note has been added",
		"data":    note,
	})
}

func (r *Repository) GetNotes(context *fiber.Ctx) error{
	noteModels := &[]models.Notes{}

	err := r.DB.Find(noteModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"could not get notes"})
			return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map {
		"message": "notes fetched successfully",
		"data":	noteModels,
	})
	return nil
}

func (r *Repository) DeleteNotes(context *fiber.Ctx) error{
	type RequestBody struct {
		IDs []string `json:"ids"`
	}
	
	var body RequestBody

	if err := context.BodyParser(&body); err != nil {
		return context.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "invalid request body",
		})
	}

	if len(body.IDs) == 0 {
		return context.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "no IDs provided",
		})
	}

	for _, id := range body.IDs {
		if _, err := uuid.Parse(id); err != nil {
			return context.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"message": fmt.Sprintf("invalid UUID: %s", id),
			})
		}
	}

	err := r.DB.Delete(&models.Notes{}, "id IN ?", body.IDs).Error
	if err != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "could not delete notes",
		})
	}

	return context.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "notes deleted successfully",
	})
}

func (r *Repository) UpdateNote(context *fiber.Ctx) error {
	id := context.Params("id")
	if id == "" {
		return context.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
	}

	if _, err := uuid.Parse(id); err != nil {
		return context.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "invalid UUID format",
		})
	}

	var note models.Notes
	if err := r.DB.First(&note, "id = ?", id).Error; err != nil {
		return context.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"message": "note not found",
		})
	}

	var updateData models.Notes
	if err := context.BodyParser(&updateData); err != nil {
		return context.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not parse request body",
		})
	}

	note.Title = updateData.Title
	note.Content = updateData.Content
	note.LastEdited = updateData.LastEdited

	if err := r.DB.Save(&note).Error; err != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "could not update note",
		})
	}

	return context.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "note updated successfully",
		"data":    note,
	})
}