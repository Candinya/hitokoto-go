package models

import "gorm.io/gorm"

type Version struct {
	gorm.Model

	BundleVersion string
}
