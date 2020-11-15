package db

import (
	"fmt"

	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/xesina/golang-echo-realworld-example-app/model"
)

func DevDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./realworld.db")
	if err != nil {
		fmt.Println("storage err: ", err)
	}
	db.DB().SetMaxIdleConns(3)
	return db
}

func PrdDB(dbDSN string) *gorm.DB {
  //dsn := "realworld:realworld@tcp(db:3306)/realworld?charset=utf8&parseTime=True"
	db, err := gorm.Open("mysql", dbDSN)
	if err != nil {
		fmt.Println("storage err: ", err)
	}
	db.DB().SetMaxIdleConns(10)
	return db
}

func TestDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./../realworld_test.db")
	if err != nil {
		fmt.Println("storage err: ", err)
	}
	db.DB().SetMaxIdleConns(3)
	db.LogMode(true)
	return db
}

func DropTestDB() error {
	if err := os.Remove("./../realworld_test.db"); err != nil {
		return err
	}
	return nil
}

//TODO: err check
func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&model.User{},
		&model.Follow{},
		&model.Article{},
		&model.Comment{},
		&model.Tag{},
	)
}
