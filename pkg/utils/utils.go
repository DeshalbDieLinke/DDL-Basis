package utils

import (
	content "ddl-server/pkg/database/models"
	"os"
	"path/filepath"
)

func GetMaterial() ([]content.Content, error) {
	contentItems := []content.Content{}
	var localErr error
	
	materialDir := GetMaterialPath()
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
func GetMaterialPath() string {
	return "/root/material"
}