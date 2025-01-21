package content

import "gorm.io/gorm"


type Content struct {
	gorm.Model
	Title   string
	Content string
	ContentType string
	Uri string
	Author string
	
}