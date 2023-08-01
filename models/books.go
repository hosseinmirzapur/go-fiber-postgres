package models

import "gorm.io/gorm"

type Books struct {
	ID        uint    `gorm:"primary_key;auto_increment" json:"id"`
	Author    *string `json:"author"`
	Title     *string `json:"title"`
	Publisher *string `json:"publisher"`
}

func MigrateBooks(db *gorm.DB) error {
	return db.AutoMigrate(&Books{})
}
