package model

import (
	"github.com/jinzhu/gorm"
)

// Project represents a project
type Project struct {
	gorm.Model
	Name        string `json:"name" gorm:"unique_index"`
	Description string `json:"description"`
}
