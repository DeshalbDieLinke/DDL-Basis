package database

import (
	content "deshalbdielinke/pkg/database/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func StartDatabase() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&content.Content{})

	// Create
	db.Create(&content.Content{Title: "D42", Description: "The Answer to the Ultimate Question of Life, the Universe, and Everything", ContentType: "text"})

	// Read
	var pieceContent content.Content
	// db.First(&pieceContent, 0) // find product with integer primary key
	db.First(&pieceContent, "title = ?", "D42") // find product with code D42
	

	// Update - update product's price to 200
	// db.Model(&pieceContent).Update("Price", 200)
	// Update - update multiple fields


	// Delete - delete product
	// db.Delete(&pieceContent, 1)
}