package database

import (
	"fmt"
	"foodbuddy/model"
	"foodbuddy/helper"

	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	databaseCredentials := helper.GetEnvVariables()

	dsn := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v?charset=utf8mb4&parseTime=True&loc=Local", databaseCredentials.DBUser, databaseCredentials.DBPassword, databaseCredentials.DBName)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("unable to connect to database, ", databaseCredentials.DBName)
	} else {
		fmt.Println("connection to database :OK")
	}

}

func AutoMigrate() {
	DB.AutoMigrate(&model.User{})
	DB.AutoMigrate(&model.Restaurant{})
	DB.AutoMigrate(&model.Category{})
	DB.AutoMigrate(&model.Product{})
	DB.AutoMigrate(&model.FavouriteProduct{})
	DB.AutoMigrate(&model.Address{})
	DB.AutoMigrate(&model.Admin{})
	DB.AutoMigrate(&model.VerificationTable{})
	DB.AutoMigrate(&model.CartItems{})
	DB.AutoMigrate(&model.Order{})
	DB.AutoMigrate(&model.OrderItem{})
	DB.AutoMigrate(&model.Payment{})
	DB.AutoMigrate(&model.PasswordReset{})
	DB.AutoMigrate(&model.CouponInventory{})
	DB.AutoMigrate(&model.CouponUsage{})
	DB.AutoMigrate(&model.UserWalletHistory{})
	DB.AutoMigrate(&model.RestaurantWalletHistory{})
	DB.AutoMigrate(&model.UserReferralHistory{})

}
