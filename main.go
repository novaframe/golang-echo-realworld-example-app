package main

import (
	"os"
	"github.com/xesina/golang-echo-realworld-example-app/db"
	"github.com/xesina/golang-echo-realworld-example-app/handler"
	"github.com/xesina/golang-echo-realworld-example-app/router"
	"github.com/xesina/golang-echo-realworld-example-app/store"
	"github.com/jinzhu/gorm"
	echoSwagger "github.com/swaggo/echo-swagger"                 // echo-swagger middleware
	_ "github.com/xesina/golang-echo-realworld-example-app/docs" // docs is generated by Swag CLI, you have to import it.
)

// @title Swagger Example API
// @version 1.0
// @description Conduit API
// @title Conduit API

// @host 127.0.0.1:8585
// @BasePath /api

// @schemes http https
// @produce	application/json
// @consumes application/json

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	r := router.New()
	r.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := r.Group("/api")
	var d	*gorm.DB
	switch os.Getenv("BACKEND_ENV") {
	case "development":
		d = db.DevDB()
		r.Logger.Info("ENV is Development:SQLite")
	case "production":
		mysqlUser := os.Getenv("MYSQL_USER")
		mysqlPassword := os.Getenv("MYSQL_PASSWORD")
		mysqlDatabase := os.Getenv("MYSQL_DATABASE")
		
		var dbDSN string
		if mysqlUser == "" || mysqlPassword == "" || mysqlDatabase == "" {
			dbDSN = "realworld:realworld@tcp(db:3306)/realworld?charset=utf8&parseTime=True"
		} else {
			dbDSN = mysqlUser + ":" + mysqlPassword + "@tcp(db:3306)/" + mysqlDatabase +"?charset=utf8&parseTime=True"
		}
		d = db.PrdDB(dbDSN)
		r.Logger.Info("ENV is Production:MySQL")
	default:
		d = db.DevDB()
		r.Logger.Info("ENV is Not Set:Default-Sqlite")
	}
	
	if debug := os.Getenv("DEBUG"); debug =="true" {
		d.LogMode(true)
	}  else if debug == "false" {
		d.LogMode(false)
	}

	db.AutoMigrate(d)

	us := store.NewUserStore(d)
	as := store.NewArticleStore(d)
	h := handler.NewHandler(us, as)
	h.Register(v1)
	r.Logger.Fatal(r.Start("0.0.0.0:8080"))
}
