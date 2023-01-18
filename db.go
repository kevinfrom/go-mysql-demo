package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	ID        uint         `json:"id" gorm:"primarykey"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at" gorm:"index"`
	Name      string       `json:"name" gorm:"not null"`
	Price     int64        `json:"price" gorm:"not null"`
}

type DbNotFoundError struct {
	Message string
}

func (e *DbNotFoundError) Error() string {
	return e.Message
}

func GetDb() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	CheckError(err)

	err = db.AutoMigrate(&Product{})
	CheckError(err)

	return db
}

func GetProducts() []Product {
	var products []Product
	db := GetDb()
	db.Find(&products)

	return products
}

func CreateProduct(name string, price int64) Product {
	product := Product{Name: name, Price: price}

	db := GetDb()
	db.Select("Name", "Price").Create(&product)

	return product
}

func GetProduct(id uint64) (Product, error) {
	var product Product
	var err error

	db := GetDb()
	result := db.First(&product, id)

	if result.RowsAffected < 1 {
		err = gorm.ErrRecordNotFound
	}

	return product, err
}

func SaveProduct(product Product) (bool, error) {
	var success bool
	var err error

	db := GetDb()

	result := db.Save(&product)

	if result.RowsAffected > 0 {
		success = true
	} else {
		err = result.Error
	}

	return success, err
}

func DeleteProduct(id uint64) (bool, error) {
	var success bool
	var err error

	db := GetDb()
	result := db.Delete(&Product{}, id)

	if result.RowsAffected < 1 {
		err = gorm.ErrRecordNotFound
	} else {
		success = true
	}

	return success, err
}
