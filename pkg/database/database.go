package database

import (
	"ddl-server/pkg/database/models"
	content "ddl-server/pkg/database/models"

	// "ddl-server/pkg/utils"
	"log"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GormMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	return func (next echo.HandlerFunc) echo.HandlerFunc { 
		return func(c echo.Context) error { 
			c.Set("db", db)
			return next(c)
		}
	}
}

func StartDatabase() *gorm.DB{
	

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&content.Content{})
	db.AutoMigrate(&models.User{})
	// contentItems, err := utils.GetMaterial()
	// if err != nil { 
	// 	print("Error getting material")
	// }
	// for _, item := range contentItems {
	// 	db.Create(&item)
	// }
	// db.Create(&content.Content{Title: "Test", Description: "Test", ContentType: "Test", Topics: "Test", Official: true})


	// Update - update product's price to 200
	// db.Model(&pieceContent).Update("Price", 200)
	// Update - update multiple fields
	// Delete - delete product
	// db.Delete(&pieceContent, 1)
	log.Printf("Database connection established")
	return db
}