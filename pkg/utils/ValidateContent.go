package utils

import (
	"ddl-server/pkg/database/models"
	"log"

	"gorm.io/gorm"
)

// Cross Reference the Database with the Space server

func SyncFileContent(db *gorm.DB) (bool, error) {
	// Fetch file keys from both Spaces and Database
	spaceFileKeys, err := ListFilesFromSpace()
	if err != nil {
		return false, err
	}

	dbFileKeys, err := ListFileKeysFromDB(db)
	if err != nil {
		return false, err
	}

	// Log the keys for debugging
	log.Printf("Space Keys: %v", spaceFileKeys)
	log.Printf("DB Keys: %v", dbFileKeys)

	// Create sets for fast lookup
	spaceFileSet := make(map[string]struct{})
	dbFileSet := make(map[string]struct{})
	for _, key := range spaceFileKeys {
		spaceFileSet[key] = struct{}{}
	}
	for _, key := range dbFileKeys {
		dbFileSet[key] = struct{}{}
	}

	// Track deletions and changes
	var deleteFromSpace, deleteFromDB []string
	isEqual := true

	// 1. Check files in Spaces but not in DB, delete from Space
	for _, key := range spaceFileKeys {
		if _, exists := dbFileSet[key]; !exists {
			DeleteFromSpace(key)
			deleteFromSpace = append(deleteFromSpace, key)
			isEqual = false
		}
	}

	// 2. Check files in DB but not in Spaces, delete from DB
	for _, key := range dbFileKeys {
		if _, exists := spaceFileSet[key]; !exists {
			MarkBrokenDB(db, key)
			deleteFromDB = append(deleteFromDB, key)
			isEqual = false
		} else {
			// Mark the file as not broken
			if err := db.Model(&models.Content{}).Where("FileName = ?", key).Update("Broken", false).Error; err != nil {
				return false, err
			}
	}
	}
	// Log which files were deleted
	if len(deleteFromSpace) > 0 {
		log.Printf("Deleted from Space: %v", deleteFromSpace)
	}
	if len(deleteFromDB) > 0 {
		log.Printf("Deleted from DB: %v", deleteFromDB)
	}

	// Final validation: After deletions, check again if everything is in sync
	if !isEqual {
		// Re-fetch the file keys after deletion
		spaceFileKeys, err = ListFilesFromSpace()
		if err != nil {
			return false, err
		}

		dbFileKeys, err = ListFileKeysFromDB(db)
		if err != nil {
			return false, err
		} 

		// Perform final check: Are the files in sync now?
		spaceFileSet = make(map[string]struct{})
		dbFileSet = make(map[string]struct{})
		for _, key := range spaceFileKeys {
			spaceFileSet[key] = struct{}{}
		}
		for _, key := range dbFileKeys {
			dbFileSet[key] = struct{}{}
		}

		// Check if everything is now in sync
		for _, key := range spaceFileKeys {
			if _, exists := dbFileSet[key]; !exists {
				isEqual = false
				break
			}
		}

		for _, key := range dbFileKeys {
			if _, exists := spaceFileSet[key]; !exists {
				MarkBrokenDB(db, key)
				isEqual = false
				break
			} else {
				// Mark the file as not broken
				if err := db.Model(&models.Content{}).Where("FileName = ?", key).Update("Broken", false).Error; err != nil {
					return false, err
				}
			}
		}
	} 


	// Log the final result
	if isEqual {
		log.Println("Content is fully synchronized.")
	} else {
		log.Println("There are still mismatches between DB and Spaces.")
	}

	return isEqual, nil

}