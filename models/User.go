package models

import "gorm.io/gorm"

// WIP

type User struct {
	gorm.Model

	Username   string
	Passkey    string
	MFAToken   string // 2FA
	EMail      string // For pass recovery & gravatar maybe?
	IsReviewer bool
	IsAdmin    bool
}
