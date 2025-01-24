package endpoints

import (
	"ddl-server/pkg/database/models"
	"ddl-server/pkg/utils"
	"log"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"gorm.io/gorm"
)



func HelloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

//TODO Add support for query parameters and search as well as returning the content from the database
func GetContent(c echo.Context) error {
	material, err := utils.GetMaterial()
	if err != nil {
		return c.String(http.StatusInternalServerError, "oopsie")
	}
	return c.JSON(http.StatusOK, material)
}

func SearchContent(c echo.Context) error {
	return c.String(http.StatusInternalServerError, "Not Implemented")
}

func CreateContent(c echo.Context) error {
	// Get the database connection
	db := c.Get("db").(*gorm.DB)

	// Get the FormData
	title := c.FormValue("title")
	description := c.FormValue("description")
	topics := c.FormValue("topics")
	official := c.FormValue("official")
	// Get the file 
	formFile, err  := c.FormFile("file")
	if err != nil { 
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "No file provided"})
	}

	log.Printf("Title: %s, Description: %s, Topics: %s, Official: %s, File: %s", title, description, topics, official, formFile.Filename)

	// Validate the input and the user
	if title == "" || description == "" || topics == "" { 
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Missing required fields"})
	}
	
	// Required Level to upload content
	requiredAccess := 2
	// Check if the official checkbox is checked
	isOfficial := official == "on"
	if isOfficial  { 
		requiredAccess = 1
	}

	// Check if the user has the required access level
	claims, err := VerifyToken("", c)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Invalid token")
	}
	user := models.User{}
	if err := db.Where("email = ?", claims.Email).First(&user).Error; err != nil { 
		return c.String(http.StatusUnauthorized, "Invalid token")
	}
	if user.AccessLevel != claims.AccessLevel || user.Email != claims.Email {
		log.Printf("INVALID TOKEN USED!!! User: %v, Claims: %v", user, claims)
		return c.String(http.StatusUnauthorized, "Invalid token")
	}

	if user.AccessLevel > requiredAccess { 
		return c.String(http.StatusUnauthorized, "Insufficient access level. "+string(requiredAccess)+" Required "+string(claims.AccessLevel)+" Provided")
	}

	// USER IS AUTHORIZED TO UPLOAD CONTENT
	fileName := formFile.Filename
	fileKey := "material/sharepics/" + fileName

	// Create the content object
	content := models.Content{Title: title, Description: description, Topics: topics, Official: isOfficial || false, AuthorID: user.ID, FileName: formFile.Filename}
	var errorCreatingContent error; 

	src, err := formFile.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error opening file"})
	}
	defer src.Close()
	// Upload the file to the bucket server
	ErrUploading := utils.UploadToSpace(src, fileKey)
	if ErrUploading != nil { 
		log.Printf("Error uploading file: %v", ErrUploading)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error uploading file"})
	}
	// Save the content object to the database
	errorCreatingContent = db.Create(&content).Error;

	if errorCreatingContent != nil { 
		log.Printf("Error creating content: %v", errorCreatingContent)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error creating content"})
	}


	return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World!"})
}


