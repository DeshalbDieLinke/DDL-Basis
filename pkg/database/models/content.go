package models

import "gorm.io/gorm"


type Content struct {
	gorm.Model
	Description string
	Title   string
	Content *string
	ContentType string
	AltText string
	Uri *string
	AuthorID uint 
	Topics string
	Official bool
	FileName string
	Broken bool
}