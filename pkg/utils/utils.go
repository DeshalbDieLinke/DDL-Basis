package utils

import (
	content "ddl-server/pkg/database/models"
	"fmt"
	"os"
	"path/filepath"
)

func GetMaterial() ([]content.Content, error) {
	contentItems := []content.Content{}
	var localErr error
	
	materialDir, err := GetMaterialPath()
	if err != nil {
		// handle the error appropriately
		return contentItems, err
	}
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

	// Build the absolute path to the material folder
	materialPath := filepath.Join(wd, "public", "material")
	fmt.Printf("Material path: %s\n", materialPath)
	return materialPath, nil
}