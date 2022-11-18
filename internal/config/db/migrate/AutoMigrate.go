package migrate

import (
	"linebot/internal/config/db"
)

func init() {
	// DropTable(&order.Order{}, &product.Product{}, &product.Price{}, &account.Account{}, &stock.Stock{})
	// // // // AutoMigrate(&product.Product{})
	// // // DB.Migrator().CreateConstraint(&order.Order{}, "OrderItem")
	// AutoMigrate(&order.Order{})
	// AutoMigrate(&product.Product{}, &product.Price{})
	// AutoMigrate(&account.Account{})
	// AutoMigrate(&stock.Stock{})
}

// AutoMigrate run auto migration
func AutoMigrate(values ...interface{}) {
	for _, value := range values {

		db.DB.AutoMigrate(value)
	}
}

func DropTable(values ...interface{}) {
	for _, value := range values {
		db.DB.Migrator().DropTable(value)
	}
}
