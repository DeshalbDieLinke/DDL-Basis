package models

import "gorm.io/gorm"


type Content struct {
	gorm.Model
	Description string
	Title   string
	Content *string
	ContentType string
	Uri *string
	Author *string
	Topics string
	Official bool
}