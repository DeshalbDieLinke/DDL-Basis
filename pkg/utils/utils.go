package utils

import (
	content "ddl-server/pkg/database/models"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GetMaterial() ([]content.Content, error) {
	contentItems := []content.Content{}
	var localErr error
	
	materialDir, err := GetMaterialPath()
	if err != nil {
		// handle the error appropriately
		return contentItems, err
	}
	fmt.Printf("Material path Print: %s", materialDir)
	filepath.WalkDir(materialDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			localErr = err
		}

		// Apply the function only to files
		if !d.IsDir() {
			contentItems = append(contentItems, content.Content{Title: d.Name(), Description: path, ContentType: "image"})
		}
		return nil
	})

	return contentItems, localErr
}


// GetMaterialPath returns the absolute path to the public/material/ folder.
func GetMaterialPath() (string, error) {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	wd = strings.ReplaceAll(wd, "build/", "")
	// Build the absolute path to the material folder
	materialPath := filepath.Join(wd, "public", "material")
	log.Printf("Material path: %s", materialPath)
	fmt.Printf("Material path Print: %s", materialPath)
	return materialPath, nil
}