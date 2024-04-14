package main

import (
	"course-api/initializers"
	"course-api/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pressly/goose"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

var (
	DB *gorm.DB
)

func main() {
	DB = initializers.DB

	db, _ := DB.DB()
	err := goose.Up(db, "./migrations")
	if err != nil {
		log.Fatalf("Error while migrating: %v", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/moduleinfo", createModuleInfo).Methods("POST")
	router.HandleFunc("/moduleinfo/{id}", getModuleInfo).Methods("GET")
	router.HandleFunc("/moduleinfo/{id}", updateModuleInfo).Methods("PUT")
	router.HandleFunc("/moduleinfo/{id}", deleteModuleInfo).Methods("DELETE")

	log.Println("Server is started in port: 8080")
	http.ListenAndServe(":8080", router)
}

func createModuleInfo(w http.ResponseWriter, router *http.Request) {
	var module models.ModuleInfo
	err := json.NewDecoder(router.Body).Decode(&module)
	if err != nil {
		http.Error(w, "Ошибка при чтении запроса", http.StatusBadRequest)
		return
	}

	module.CreatedAt = time.Now()
	module.UpdatedAt = time.Now()

	result := DB.Create(&module)
	if result.Error != nil {
		http.Error(w, "Ошибка при создании записи", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(module)
}

func getModuleInfo(w http.ResponseWriter, router *http.Request) {
	params := mux.Vars(router)
	id := params["id"]

	var module models.ModuleInfo
	result := DB.First(&module, id)
	if result.Error != nil {
		http.Error(w, "Запись не найдена", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(module)
}

func updateModuleInfo(w http.ResponseWriter, router *http.Request) {
	params := mux.Vars(router)
	id := params["id"]

	var module models.ModuleInfo
	result := DB.First(&module, id)
	if result.Error != nil {
		http.Error(w, "Запись не найдена", http.StatusNotFound)
		return
	}

	var updatedModule models.ModuleInfo
	err := json.NewDecoder(router.Body).Decode(&updatedModule)
	if err != nil {
		http.Error(w, "Ошибка при чтении запроса", http.StatusBadRequest)
		return
	}

	module.ModuleName = updatedModule.ModuleName
	module.ModuleDuration = updatedModule.ModuleDuration
	module.ExamType = updatedModule.ExamType
	module.Version = updatedModule.Version
	module.UpdatedAt = time.Now()

	DB.Save(&module)

	json.NewEncoder(w).Encode(module)
}

func deleteModuleInfo(w http.ResponseWriter, router *http.Request) {
	params := mux.Vars(router)
	id := params["id"]

	var module models.ModuleInfo
	result := DB.First(&module, id)
	if result.Error != nil {
		http.Error(w, "Запись не найдена", http.StatusNotFound)
		return
	}

	DB.Delete(&module)

	w.WriteHeader(http.StatusNoContent)
}
