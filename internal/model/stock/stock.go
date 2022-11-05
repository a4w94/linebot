package stock

import (
	"errors"
	"fmt"
	"linebot/internal/config/db"
	"linebot/internal/model/product"
	"log"
	"time"

	"gorm.io/gorm"
)

//庫存
type Stock struct {
	gorm.Model
	Date      time.Time
	ProductId uint
	RemainNum int
}

func (stock Stock) Add() error {
	product, _ := product.GetById(int64(stock.ProductId))
	if stock.RemainNum > product.TotlaNum {

		err := errors.New("remain number can't bigger than total num")
		log.Fatal(err)
	}
	return db.BeginTransaction(db.DB, func(tx *gorm.DB) error {
		return tx.Create(&stock).Error
	})
}
func GetAll() ([]Stock, error) {
	var Stocks []Stock
	err := db.DB.Find(&Stocks).Error
	return Stocks, err
}

func GetStocks_By_ID_and_DateRange(pid uint, start, end time.Time) ([]Stock, error) {

	var stocks []Stock
	var tmp_time time.Time
	tmp_time = start
	var err error
	for tmp_time != end {
		fmt.Println("tmp time", tmp_time)
		var s []Stock

		tmp_time_last := tmp_time.Add(time.Hour*23 + time.Minute*59 + time.Second*59)
		err := db.DB.Where("product_id=? AND date BETWEEN ? AND ?", pid, tmp_time, tmp_time_last).Find(&s).Error
		fmt.Println("s", s)
		if err != nil {
			log.Println("GetStocks_By_ID_and_DateRange failed")
		}
		if len(s) != 1 {
			err = errors.New("get stock date unexist")
			return []Stock{}, err

		}
		tmp_time = tmp_time.AddDate(0, 0, 1)
	}
	fmt.Println("stocks", stocks, err)
	return stocks, err
}

func (stock *Stock) UpdateStock() error {
	return db.BeginTransaction(db.DB, func(tx *gorm.DB) error {
		return tx.Save(&stock).Error
	})
}
