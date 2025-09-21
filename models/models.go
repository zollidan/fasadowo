package models

import "time"

// Category 1 -> many Subcategories
type Category struct {
	ID          uint          `gorm:"primaryKey"`
	Name        string        `gorm:"type:text;not null"`
	CreatedAt   time.Time
	Subcategories []Subcategory `gorm:"foreignKey:CategoryID"`
}

// Subcategory 1 -> many Collections
type Subcategory struct {
	ID          uint         `gorm:"primaryKey"`
	Name        string       `gorm:"type:text;not null"`
	CreatedAt   time.Time
	CategoryID  uint         `gorm:"not null"`
	Collections []Collection `gorm:"foreignKey:SubcategoryID"`
}

// Collection 1 -> many Products
type Collection struct {
	ID           uint      `gorm:"primaryKey"`
	Name         string    `gorm:"type:text;not null"`
	CreatedAt    time.Time
	SubcategoryID uint     `gorm:"not null"`
	Products     []Product `gorm:"foreignKey:CollectionID"`
}

// Product belongs to Collection
type Product struct {
	ID           uint      `gorm:"primaryKey"`
	Name         string    `gorm:"type:text;not null"`
	Price        float64   `gorm:"type:numeric(10,2)"`
	CreatedAt    time.Time
	CollectionID uint      `gorm:"not null"`
}
