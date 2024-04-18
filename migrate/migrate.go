package main

import (
	"course-api/initializers"
	"course-api/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.ModuleInfo{})
}

//goose -dir=./migrations/ postgres "user=postgres password=91926499 dbname=courseapi sslmode=disable" up -force
