package utils

import (
	"ddl-server/pkg/database/models"
	"log"

	"gorm.io/gorm"
)

func ListFileKeysfromDB(db *gorm.DB) ([]string, error) {
	// Get the file keys from the database
	var uploads []models.Content
	var FileKeys []string

	if err := db.Model(&models.Content{}).Select("FileKey").Find(&uploads).Error; err != nil {
		return nil, err
	}
	for _, upload := range uploads {
		FileKeys = append(FileKeys, upload.FileKey)
	}

	return FileKeys, nil
}

func DeleteFromDB(db *gorm.DB, Url string) error {
	// Delete the file from the database
	if err := db.Where("Url = ?",Url).Delete(&models.Content{}).Error; err != nil {
		return err
	}
	log.Printf("Deleted %v from the database", Url)
	return nil
}

func MarkBrokenDB(db *gorm.DB, key string) error {
	// Mark the file as broken in the database
	if err := db.Model(&models.Content{}).Where("FileKey = ?", key).Update("Broken", true).Error; err != nil {
		return err
	}
	log.Printf("Marked %v as broken in the database", key)
	return nil
}