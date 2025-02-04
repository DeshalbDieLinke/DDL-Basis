package endpoints

import (
	"ddl-server/pkg/database/models"
	"ddl-server/pkg/utils"
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
	} else if queryParams.Get("id") != "" { 
		id := queryParams.Get("id")
		if err := db.Where("id = ?", id).Find(&content).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error fetching content"})
		}
		return c.JSON(http.StatusOK, content)
	}
	return c.String(http.StatusInternalServerError, "Not Implemented")
}

// func CreateContentNew(c echo.Context) error {
// 	// Verify the s
// }

func CreateContent(c echo.Context) error {
	
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
	// requiredAccess := 2
	// Check if the official checkbox is checked
	


	// Get the database connection
	db := c.Get("db").(*gorm.DB)

	// user := models.User{}
	// if err := db.Where("email = ?", claims.Email).First(&user).Error; err != nil {
	// 	return c.String(http.StatusUnauthorized, "Invalid token")
	// }
	// if user.AccessLevel != claims.AccessLevel || user.Email != claims.Email {
	// 	log.Printf("INVALID TOKEN USED!!! User: %v, Claims: %v", user, claims)
	// 	return c.String(http.StatusUnauthorized, "Invalid token")
	// }

	ContentIsOfficial := official == "on"

	userMetaData, err := utils.GetUserRoleData(c) 
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	
	isAdmin := *userMetaData.Role == "admin" 
	AuthOfficial := *userMetaData.Official || isAdmin

	canUpload := false
	if ContentIsOfficial {
		canUpload = AuthOfficial
	} else {
		canUpload = isAdmin || *userMetaData.Upload
	}

	if !canUpload { 
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Insufficient permissions"})
	}


	// USER IS AUTHORIZED TO UPLOAD CONTENT
	fileKey := "material/sharepics/" + strings.Replace(formFile.Filename, " ", "_", -1)

	uri := "https://ddl.fra1.cdn.digitaloceanspaces.com/material/sharepics/" + url.PathEscape(strings.Replace(formFile.Filename, " ", "_", -1))
		
	
	usr, err := utils.GetUserFromContext(c) 
	if err != nil { 
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Error getting user"})
	}
	content := models.Content{Title: title, Description: description, Topics: topics, Official: ContentIsOfficial || false, AuthorClerkID: usr.ID , FileName: formFile.Filename, FileKey: fileKey, Uri: &uri, AltText: altText}
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
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error unescaping file key"})
	}
	ErrUploading := utils.UploadToSpace(file, content.FileKey)
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

// type UpdateUserRequest struct {
// 	ID          int     `json:"id"`
// 	Email       *string `json:"email"`
// 	Password    *string `json:"password"`
// 	AccessLevel *int    `json:"accessLevel"`
// 	Username    *string `json:"username"`
// }

// func UpdateUser(c echo.Context) error {
// 	// User who issued the request:
// 	token, err := utils.GetToken(c)
// 	if err != nil {
// 		return c.JSON(401, map[string]string{"error": "No token provided"})
// 	}
// 	claims, err := utils.GetTokenClaims(token)
// 	if err != nil {
// 		return c.String(http.StatusUnauthorized, "Invalid token")
// 	}
// 	db := c.Get("db").(*gorm.DB)

// 	editAuthor := models.User{}
// 	err = db.Where("id = ?", claims.ID).First(&editAuthor).Error
// 	if err != nil {
// 		return c.JSON(http.StatusNotFound, map[string]string{"message": "Error fetching user"})
// 	}

// 	updateRequest := new(UpdateUserRequest)
// 	err = c.Bind(&updateRequest)
// 	if err != nil {
// 		log.Printf("Error binding request: %v", err)
// 		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
// 	}

// 	// Person to be updated
// 	editObject := models.User{}
// 	if err := db.Where("id = ?", updateRequest.ID).First(&editObject).Error; err != nil {
// 		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error fetching user"})
// 	}

// 	editingSelf := editAuthor.ID == editObject.ID
// 	editAllowed := false
// 	editingAccessLevel := false
// 	if updateRequest.AccessLevel != nil {
// 		editingAccessLevel = *updateRequest.AccessLevel != editObject.AccessLevel 
// 	}

// 	denialReason := "Edit request denied"
// 	// Check if the user is editing themselves
	
// 	if editingSelf { 
// 		if editingAccessLevel { 
// 			denialReason = "Edit request denied: Cannot change own access level"
// 			editAllowed = false
// 		} else {
// 			editAllowed = true
// 		}
// 	} else if editAuthor.AccessLevel == 0 {
// 		editAllowed = true
// 	} else {
// 		denialReason = "Edit request denied: Insufficient permissions"
// 	}


// 	if !editAllowed {
// 		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Edit request denied: " + denialReason}) 
// 	}

// 	// Update user in Database
// 	if editAllowed {
// 		if updateRequest.Email != nil {
// 			// Check for email uniqueness if email is being updated
// 			if updateRequest.Email != nil && *updateRequest.Email != editObject.Email {
// 				var count int64
// 				db.Model(&models.User{}).Where("email = ?", *updateRequest.Email).Count(&count)
// 				if count > 0 {
// 					return c.JSON(http.StatusConflict, map[string]string{"message": "Email already in use"})
// 				}
// 			}
// 			editObject.Email = *updateRequest.Email
// 		}
// 		if updateRequest.Password != nil {
// 			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*updateRequest.Password), bcrypt.DefaultCost)
// 			if err != nil {
// 				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error hashing password"})
// 			}
// 			editObject.Password = string(hashedPassword)
// 		}
// 		if updateRequest.AccessLevel != nil {
// 			editObject.AccessLevel = *updateRequest.AccessLevel
// 		}
// 		if updateRequest.Username != nil {
// 			editObject.Username = updateRequest.Username
// 		}
// 		err := db.Save(&editObject)
// 		if err.Error != nil {
// 			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error updating user"})
// 		}
// 		return c.JSON(http.StatusOK, map[string]string{"message": "User updated"})
// 	}
// 	return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error updating user"})

// }

func DeleteContentItem(c echo.Context) error { 
	// Get the content ID
	contentID := c.QueryParam("id")
	if contentID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request: No content ID provided"})
	}
	db := c.Get("db").(*gorm.DB)

	// Get the content from the DB
	content := models.Content{}
	if err := db.Where("id = ?", contentID).First(&content).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Content not found"})
	}

	permitted := utils.VerifyPermissions(2, c, &utils.Target{
		ContentItem: &content,
	})
	if !permitted {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Insufficient permissions"})
	}
	// Delete the content from the database
	if err := db.Delete(&content).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error deleting content"})
	}
	// Delete the file from the bucket
	if err := utils.DeleteFromSpace(content.FileKey); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error deleting file"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Content deleted"})
}

func UpdateContent(c echo.Context) error { 
	// Get the content ID
	contentID := c.QueryParam("id")
	if contentID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request: No content ID provided"})
	}
	db := c.Get("db").(*gorm.DB)
	contentItem := models.Content{}
	if err := db.Where("id = ?", contentID).First(&contentItem).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Content not found"})
	}

	// permitted := utils.VerifyPermissions(2, c, &utils.Target{
	// 	ContentItem: &contentItem,
	// })

	// Permitted if either the user is an admin or the author of the content
	userMetaData, err := utils.GetUserRoleData(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Error getting user"})
	}
	usr, err := utils.GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Error getting user"})
	}
	permitted := *userMetaData.Role == "admin" || contentItem.AuthorClerkID == usr.ID



	if !permitted {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Insufficient permissions"})
	}
	// var RequestData struct { 
	// 	id int `json:"id"`
	// 	Title string `json:"title"`
	// 	Description string `json:"description"`
	// 	Topics string `json:"topics"`
	// 	AltText string `json:"altText"`
	// 	Official bool `json:"official"`
	// 	formFile *multipart.File `json:"file"`
	// 	ContentType string `json:"type"`
	// 	url string `json:"url"`
	// }


	err = c.Request().ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		// Is not mutlipart form. That's okay
	}

	// err := c.Bind(&RequestData)
	title := c.FormValue("title")
	description := c.FormValue("description")
	topics := c.FormValue("topics")
	altText := c.FormValue("altText")
	official := c.FormValue("official")

	// Handle file upload
	formFile, err := c.FormFile("file")
	if err != nil {
		// This is not an error, the file is optional

	}

	// contentType := c.FormValue("type")
	// url := c.FormValue("url")

	

	// Validate the input and the user
	if title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Missing required fields"})
	}

	// Update the content object
	contentItem.Title = title
	contentItem.Description = description
	contentItem.Topics = topics
	contentItem.Official = official == "on"
	contentItem.AltText = altText

	if formFile != nil {
		// Update the URI
		fileKey := "material/sharepics/" + strings.Replace(formFile.Filename, " ", "_", -1)

		uri := "https://ddl.fra1.cdn.digitaloceanspaces.com/material/sharepics/" + url.PathEscape(strings.Replace(formFile.Filename, " ", "_", -1))
		
		contentItem.FileKey = fileKey

		// Remove File Metadata
		src, err := formFile.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error opening file"})
		}
		defer src.Close()

		file, err := utils.CleanFile(src)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error removing metadata"})
		}

		file.Seek(0, io.SeekStart) // Reset file pointer to the start
		log.Printf("File: %v", file)
		// Delete the old file from the bucket
		if err := utils.DeleteFromSpace(contentItem.FileKey); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error deleting file"})
		}

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error unescaping file key"})
		}
		ErrUploading := utils.UploadToSpace(file, contentItem.FileKey)
		if ErrUploading != nil {
			log.Printf("Error uploading file: %v", ErrUploading)
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error uploading file"})
		}

		
		contentItem.Uri = &uri
		contentItem.FileName = formFile.Filename
	}

	if err := db.Save(&contentItem).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error updating content"})
	}
	
	return c.JSON(http.StatusOK, map[string]string{"message": "Content updated"})
}

func Topics(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string][]string{"topics": utils.GetTopics()})
}
