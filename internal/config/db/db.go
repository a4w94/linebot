package db

import (
	"fmt"
	"linebot/internal/config"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

//之後修改成自動init
func Init() {
	fmt.Println("初始化數據庫")
	// DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	// if err != nil {
	// 	log.Fatalf("Error opening database: %q", err)
	// }

	dbConfig := config.Config.DB
	if config.Config.DB.Adapter == "mysql" {
		dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)
		fmt.Println("dsn:", dsn)

		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	} else if config.Config.DB.Adapter == "postgres" {

		dsn := fmt.Sprintf("postgres://%v:%v@%v/%v", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Name)
		fmt.Println("dsn:", dsn)

		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		fmt.Println("connet database fail,pleaes check parametre", err)

		//退出程序
		os.Exit(1)
	}

	// DB.Migrator().DropTable(&product.Camp{}, &model.Account{}, &model.Tag{}, &model.TagMap{})

	// //migrate table
	// _ = DB.AutoMigrate(&model.Camp{}, &model.Account{}, &model.Tag{}, &model.TagMap{})
}

func BeginTransaction(DB *gorm.DB, process func(tx *gorm.DB) error) error {
	tx := DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	err := process(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error

	return err
}
