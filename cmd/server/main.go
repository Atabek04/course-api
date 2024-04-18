package main

import (
	"course-api/cmd/routers"
	"course-api/initializers"
	"github.com/pressly/goose"
	"gorm.io/gorm"
	"log"
	"net/http"
)

var (
	DB *gorm.DB
)

func main() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	DB = initializers.DB

	db, _ := DB.DB()
	err := goose.Up(db, "./migrations")
	if err != nil {
		log.Fatalf("Error while migrating: %v", err)
	}

	router := routers.NewRouter(DB)
	log.Println("Server is started in port: 3000")
	http.ListenAndServe(":3000", router)
}
