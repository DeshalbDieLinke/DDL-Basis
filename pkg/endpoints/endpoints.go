package endpoints

import (
	"ddl-server/pkg/database/models"
	"ddl-server/pkg/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	echo "github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func HelloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// TODO Add support for query parameters and search as well as returning the content from the database
func GetContent(c echo.Context) error {
	queryParams := c.QueryParams()
	db := c.Get("db").(*gorm.DB)

	if len(queryParams) > 0 {
		// Search for content
		return SearchContent(c, queryParams)
	}
	// Get the content ID
	content := []models.Content{}
	if err := db.Find(&content).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error fetching content"})
	}
	log.Printf("Content: %v", content)

	return c.JSON(http.StatusOK, content)
}

func SearchContent(c echo.Context, queryParams url.Values) error {
	db := c.Get("db").(*gorm.DB)
	content := []models.Content{}

	log.Printf("Query Params: %v", queryParams)
	if queryParams.Get("author") != "" {
		author := queryParams.Get("author")
		if err := db.Where("author_id = ?", author).Find(&content).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error fetching content"})
		}
		log.Printf("Content: %v", content)
		return c.JSON(http.StatusOK, content)
	}
	return c.String(http.StatusInternalServerError, "Not Implemented")
}

func CreateContent(c echo.Context) error {
	// Get the database connection
	db := c.Get("db").(*gorm.DB)

	// Get the FormData
	title := c.FormValue("title")
	description := c.FormValue("description")
	topics := c.FormValue("topics")
	altText := c.FormValue("altText")

	official := c.FormValue("official")
	// Get the file
	formFile, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "No file provided"})
	}

	// Check if the file is too large. 10MB Limit.
	if formFile.Size > 10000000 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "File too large"})
	}
	// Validate the input and the user
	if title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Missing required fields"})
	}

	// Required Level to upload content
	requiredAccess := 2
	// Check if the official checkbox is checked
	isOfficial := official == "on"
	if isOfficial {
		requiredAccess = 1
	}

	// Check if the user has the required access level
	token, err := GetToken(c)
	if err != nil {
		return c.JSON(401, map[string]string{"error": "No token provided"})
	}
	claims, err := GetTokenClaims(token)
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
		return c.String(http.StatusUnauthorized, "Insufficient access level. "+fmt.Sprint(requiredAccess)+" Required "+fmt.Sprint(claims.AccessLevel)+" Provided")
	}

	// USER IS AUTHORIZED TO UPLOAD CONTENT
	fileName := formFile.Filename
	fileKey := "material/sharepics/" + fileName

	// Create the content object
	uri := "https://ddl.fra1.cdn.digitaloceanspaces.com/" + fileKey
	uri = strings.Replace(uri, " ", "%20", -1)

	content := models.Content{Title: title, Description: description, Topics: topics, Official: isOfficial || false, AuthorID: user.ID, FileName: formFile.Filename, Uri: &uri, AltText: altText}
	var errorCreatingContent error

	src, err := formFile.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error opening file"})
	}
	defer src.Close()

	// Remove File Metadata
	file, err := utils.CleanFile(src)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error removing metadata"})
	}

	file.Seek(0, io.SeekStart) // Reset file pointer to the start
	log.Printf("File: %v", file)

	// Upload the file to the bucket server
	ErrUploading := utils.UploadToSpace(file, fileKey)
	if ErrUploading != nil {
		log.Printf("Error uploading file: %v", ErrUploading)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error uploading file"})
	}
	// Save the content object to the database
	errorCreatingContent = db.Create(&content).Error

	if errorCreatingContent != nil {
		log.Printf("Error creating content: %v", errorCreatingContent)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error creating content AHHHHHHHHHH"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World!"})
}

type UpdateUserRequest struct {
	Email       *string `json:"email"`
	Password    *string `json:"password"`
	AccessLevel *int    `json:"accessLevel"`
	Username    *string `json:"username"`
}

func UpdateUser(c echo.Context) error {
	// User who issued the request:
	token, err := GetToken(c)
	if err != nil {
		return c.JSON(401, map[string]string{"error": "No token provided"})
	}
	claims, err := GetTokenClaims(token)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Invalid token")
	}
	db := c.Get("db").(*gorm.DB)

	editAuthor := models.User{}
	err = db.Where("email = ?", claims.Email).First(&editAuthor).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error fetching user"})
	}

	updateRequest := new(UpdateUserRequest)
	err = c.Bind(&updateRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}
	// Person to be updated
	editObject := models.User{}
	if err := db.Where("email = ?", updateRequest.Email).First(&editObject).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error fetching user"})
	}

	editAllowed := false
	if editAuthor.AccessLevel == 0 {
		editAllowed = true
	} else if editAuthor == editObject && updateRequest.AccessLevel == nil {
		editAllowed = true
	}
	if !editAllowed {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Edit request denied"})
	}

	// Update user in Database
	if editAllowed {
		if updateRequest.Email != nil {
			editObject.Email = *updateRequest.Email
		}
		if updateRequest.Password != nil {
			editObject.Password = *updateRequest.Password
		}
		if updateRequest.AccessLevel != nil {
			editObject.AccessLevel = *updateRequest.AccessLevel
		}
		if updateRequest.Username != nil {
			editObject.Username = updateRequest.Username
		}
		db.Save(&editObject)
		return c.JSON(http.StatusOK, map[string]string{"message": "User updated"})
	}
	return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error updating user"})

}

func Topics(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string][]string{"topics": utils.GetTopics()})
}
