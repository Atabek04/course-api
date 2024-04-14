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
