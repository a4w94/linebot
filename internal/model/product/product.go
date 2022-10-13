package product

import (
	"linebot/internal/config/db"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	CampRoundName string `gorm:"not null"`
	Price         Price  `gorm:"embedded"`
	Size          string
	ImageUri      pq.StringArray `gorm:"type:text[]"`
	Description   string
}

func (product Product) Add() error {
	return db.BeginTransaction(db.DB, func(tx *gorm.DB) error {
		return tx.Create(&product).Error
	})
}

func GetAll() ([]Product, error) {
	var Products []Product
	err := db.DB.Find(&Products).Error

	return Products, err
}

func GetById(Id int64) (Product, error) {
	var GetProduct Product
	err := db.DB.Where("Id=?", Id).Find(&GetProduct).Error

	return GetProduct, err
}
func GetIdByCampRoundName(name string) (Product, error) {
	var GetProduct Product
	err := db.DB.Where("camp_round_name=?", name).Find(&GetProduct).Error

	return GetProduct, err
}

func DeleteById(Id int64) (Product, error) {
	var product Product
	err := db.DB.Where("Id=?", Id).Delete(&product).Error
	return product, err
}
