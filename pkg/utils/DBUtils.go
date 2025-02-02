package utils

import (
	"ddl-server/pkg/database/models"
	"log"

	"gorm.io/gorm"
)

func ListFileKeysFromDB(db *gorm.DB) ([]string, error) {
	// Get the file keys from the database
	var uploads []models.Content
	var fileKeys []string

	if err := db.Model(&models.Content{}).Select("FileName").Find(&uploads).Error; err != nil {
		return nil, err
	}
	for _, upload := range uploads {
		fileKeys = append(fileKeys, upload.FileName)
	}

	return fileKeys, nil
}

func DeleteFromDB(db *gorm.DB, key string) error {
	// Delete the file from the database
	if err := db.Where("FileName = ?", key).Delete(&models.Content{}).Error; err != nil {
		return err
	}
	log.Printf("Deleted %v from the database", key)
	return nil
}

func MarkBrokenDB(db *gorm.DB, key string) error {
	// Mark the file as broken in the database
	if err := db.Model(&models.Content{}).Where("FileName = ?", key).Update("Broken", true).Error; err != nil {
		return err
	}
	log.Printf("Marked %v as broken in the database", key)
	return nil
}