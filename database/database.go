package database

import (
	"fmt"
	"log"

	"github.com/zollidan/fasadowo/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDatabase() (db *gorm.DB){
	var err error
	db, err = gorm.Open(sqlite.Open("dev.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	err = db.AutoMigrate(
		&models.Category{},
		&models.Subcategory{},
		&models.Collection{},
		&models.Product{},
		&models.User{},
	)
	if err != nil {
		log.Fatal("failed to migrate database")
	}

	category := models.Category{Name: "Electronics"}
	db.Create(&category)

	sub1 := models.Subcategory{Name: "Smartphones", CategoryID: category.ID}
	sub2 := models.Subcategory{Name: "Laptops", CategoryID: category.ID}
	db.Create(&sub1)
	db.Create(&sub2)

	col1 := models.Collection{Name: "Flagship Phones", SubcategoryID: sub1.ID}
	col2 := models.Collection{Name: "Gaming Laptops", SubcategoryID: sub2.ID}
	db.Create(&col1)
	db.Create(&col2)

	products := []models.Product{
		{Name: "iPhone 15 Pro", Price: 1200, CollectionID: col1.ID},
		{Name: "Samsung Galaxy S24", Price: 1100, CollectionID: col1.ID},
		{Name: "Asus ROG Zephyrus", Price: 2000, CollectionID: col2.ID},
		{Name: "MSI Raider GE78", Price: 2200, CollectionID: col2.ID},
	}
	for _, p := range products {
		db.Create(&p)
	}

	fmt.Println("Database initialized with test data ✅")

	return db
}