package model

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Category struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"not null;unique"`
	SortOrder int       `json:"sort_order" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type Navigation struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	CategoryID  uint      `json:"category_id" gorm:"not null"`
	Title       string    `json:"title" gorm:"not null"`
	URL         string    `json:"url" gorm:"not null"`
	Icon        string    `json:"icon"`
	Description string    `json:"description"`
	SortOrder   int       `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	Category    Category  `json:"category,omitempty" gorm:"foreignKey:CategoryID;references:ID;constraint:OnDelete:CASCADE"`
}

type NavTree struct {
	CategoryID   uint          `json:"category_id"`
	CategoryName string        `json:"category_name"`
	Items        []NavItem     `json:"items"`
}

type NavItem struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Online      bool   `json:"online"`
}

var DB *gorm.DB

func InitDB(dbPath string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}
	return DB.AutoMigrate(&Category{}, &Navigation{})
}
