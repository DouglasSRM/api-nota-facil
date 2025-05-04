package models

import (
	"gorm.io/gorm"
	"time"
)

type Notes struct {
	ID 			string		`gorm:"type:uuid;primary_key" json:"id"`
	Title		string 		`json:"title"`
	Content		string 		`json:"content"`
	LastEdited	time.Time 	`json:"lastEdited"`
}

func MigrateNotes(db *gorm.DB) error {
	err := db.AutoMigrate(&Notes{})
	return err
}