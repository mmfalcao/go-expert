package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Category struct {
	ID       int `gorm:"primaryKey"`
	Name     string
	Products []Product
}

type Product struct {
	ID           int `gorm:"primaryKey"`
	Name         string
	Price        float64
	CategoryID   int
	Category     Category
	SerialNumber SerialNumber
	gorm.Model
}

type SerialNumber struct {
	ID        int `gorm:"primaryKey"`
	Number    string
	ProductID int
}

func main() {
	dsn := "root:root@tcp(localhost:3306)/goexpert?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Product{}, &Category{}, &SerialNumber{})

	// create categories
	electronics := []Category{
		{Name: "Accessories & Supplies"},
		{Name: "Camera & Photo"},
		{Name: "Car & Vehicle Eletronics"},
		{Name: "Cell Phones & Accessories"},
		{Name: "Computers & Accessories"},
		{Name: "GPS & Navigation"},
		{Name: "Headphones"},
		{Name: "Home Audio"},
		{Name: "Office Eletronics"},
		{Name: "Portable Audio & Video"},
		{Name: "Security & Surveillance"},
		{Name: "Service Plans"},
		{Name: "Television & Video"},
		{Name: "Video Game Consoles & Accessories"},
		{Name: "Video Projectors"},
		{Name: "Wearable Technology"},
		{Name: "eBooks Readers & Accessories"},
	}
	db.Create(&electronics)

	// create
	// db.Create(&Product{
	// 	Name:  "Notebook",
	// 	Price: 1000.00,
	// })

	// create batch
	products := []Product{
		{Name: "iPhone Charger [Apple MFi Certified]", Price: 8.48, CategoryID: electronics[0].ID},
		{Name: "SanDisk 128GB Ultra microSDXC UHS-I Memory Card with Adapter", Price: 14.79, CategoryID: electronics[1].ID},
		{Name: "Magnetic Wireless Car Charger", Price: 19.90, CategoryID: electronics[2].ID},
		{Name: "Newmowa 60 LED High Power Rechargeable Clip", Price: 30.59, CategoryID: electronics[3].ID},
		{Name: "Lenovo 2022 Newest Ideapad 3 Laptop", Price: 383.75, CategoryID: electronics[4].ID},
	}
	db.Create(&products)

	serials := []SerialNumber{
		{Number: "HCgWnEiZ", ProductID: products[0].ID},
		{Number: "MyGm473b", ProductID: products[1].ID},
		{Number: "uukJGYkf", ProductID: products[2].ID},
		{Number: "xqisDaTX", ProductID: products[3].ID},
		{Number: "hytPixrK", ProductID: products[4].ID},
	}
	db.Create(&serials)

	// HasMany
	var categories []Category
	err = db.Model(&Category{}).Preload("Products.SerialNumber").Find(&categories).Error
	if err != nil {
		panic(err)
	}
	for _, category := range categories {
		fmt.Println(category.Name, ":")
		for _, product := range category.Products {
			fmt.Println("- Product: ", product.Name, " - SerialNumber: ", product.SerialNumber.Number)
		}
	}

	// select one
	// var product Product
	// db.First(&product, 1)
	// db.First(&product, "name = ?", "Mouse")
	// fmt.Println(product)

	// select all
	/*var pdts []Product
	db.Preload("Category").Preload("SerialNumber").Find(&pdts)
	for _, product := range pdts {
		fmt.Println("Category: ", product.Category.Name, " | Product: ", product.Name, " | Serial Number: ", product.SerialNumber.Number)
	}*/

	// var products []Product
	// db.Limit(2).Offset(2).Find(&products)
	// for _, product := range products {
	// 	fmt.Println(product)
	// }

	// where
	// var products []Product
	// db.Where("price > ?", 260).Find(&products)
	// for _, product := range products {
	// 	fmt.Println(product)
	// }

	// like
	// var products []Product
	// db.Where("name LIKE ?", "%k%").Find(&products)
	// for _, product := range products {
	// 	fmt.Println(product)
	// }

	// var p Product
	// db.First(&p, 1)
	// p.Name = "Mouse"
	// db.Save(&p)

	// var p2 Product
	// db.First(&p2, 1)
	// fmt.Println(p2.Name)
	// db.Delete(&p2)
}
