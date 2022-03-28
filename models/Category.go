package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model

	Name        string
	Description string
	Key         string
}

type Categories []Category
