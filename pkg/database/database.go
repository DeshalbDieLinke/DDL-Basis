package database

import (
	content "deshalbdielinke/pkg/database/models"
	"deshalbdielinke/pkg/utils"
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func StartDatabase() *gorm.DB{
	

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&content.Content{})
	contentItems, err := utils.GetMaterial()
	if err != nil { 
		print("Error getting material")
	}
	for _, item := range contentItems {
		db.Create(&item)
	}
	// db.Create(&content.Content{Title: "Test", Description: "Test", ContentType: "Test", Topics: "Test", Official: true})

	// Read
	var pieceContent content.Content
	// db.First(&pieceContent, 0) // find product with integer primary key
	result := db.First(&pieceContent)
	if result.Error != nil {
		log.Fatalf("Error fetching record: %v", result.Error)
	}
	fmt.Printf("Test: %v", pieceContent.Title )
	

	// Update - update product's price to 200
	// db.Model(&pieceContent).Update("Price", 200)
	// Update - update multiple fields


	// Delete - delete product
	// db.Delete(&pieceContent, 1)
	return db
}