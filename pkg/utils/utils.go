package utils

import (
	"bytes"
	"ddl-server/pkg/types"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"

	"github.com/golang-jwt/jwt/v5"
)

func CleanFile(file io.ReadSeeker) (io.ReadSeeker, error) {
    // Decode EXIF data
    _, err := exif.Decode(file)
    if err != nil {
        // If there's no EXIF data, return the original file
        if !exif.IsExifError(err) {
            return file, nil
        }
        return nil, err // Return other errors
    }
    // Create a buffer to hold the cleaned file content
    var buf bytes.Buffer
    // Reset the file reader to the beginning
    if seeker, ok := file.(io.Seeker); ok {
        _, err = seeker.Seek(0, 0)
        if err != nil {
            return nil, err
        }
    }
    // Copy the file content to the buffer (without EXIF data)
    _, err = io.Copy(&buf, file)
    if err != nil {
        return nil, err
    }

	returnFile := bytes.NewReader(buf.Bytes())

    return returnFile, nil
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
	log.Printf("Material path: %s \n", materialPath)
	return materialPath, nil
}

// Function to generate a token
func GenerateToken(email string, accessLevel int) (string, error) {
	// Set custom claims
	claims := &types.JWTClaims{
		Email:       email,
		AccessLevel: accessLevel,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)), // Token valid for 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),

		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}

func GetTopics() []string {
	topics := []string{
		"Klima",
		"Frieden",
		"Demokratie",
		"Antifaschismus",
		"Antidiskriminierung",
		"Antikapitalismus",
		"Feminismus",
		"Queer",
		"Dort und Agrar",
		"Vielfalt",
		"integration",
		"Mieten",
		"Mieten",
		"Soziale Gerechtigkeit",
		"Bahn",
		"Lebensmittel",
		"Infrastruktur",
		"Arbeit und Inflation",
		"Klassenkampf",
	}
	return topics
}