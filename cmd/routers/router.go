package routers

import (
	"course-api/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", createModuleInfo(db)).Methods("POST")
	router.HandleFunc("/{id}", getModuleInfo(db)).Methods("GET")
	router.HandleFunc("/{id}", updateModuleInfo(db)).Methods("PUT")
	router.HandleFunc("/{id}", deleteModuleInfo(db)).Methods("DELETE")

	router.HandleFunc("/departments", createDepartmentInfo(db)).Methods("POST")
	router.HandleFunc("/departments/{id}", getDepartmentInfo(db)).Methods("GET")

	return router
}

func createDepartmentInfo(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var deparment models.DepartmentInfo
		err := json.NewDecoder(r.Body).Decode(&deparment)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		result := db.Create(&deparment)
		if result.Error != nil {
			http.Error(w, "Failed to create module info", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(deparment)
	}
}

func getDepartmentInfo(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var department models.DepartmentInfo
		result := db.First(&department, id)
		if result.Error != nil {
			http.Error(w, "Module info not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(department)
	}
}

func createModuleInfo(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var module models.ModuleInfo
		err := json.NewDecoder(r.Body).Decode(&module)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		module.CreatedAt = time.Now()
		module.UpdatedAt = time.Now()

		result := db.Create(&module)
		if result.Error != nil {
			http.Error(w, "Failed to create module info", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(module)
	}
}

func getModuleInfo(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var module models.ModuleInfo
		result := db.First(&module, id)
		if result.Error != nil {
			http.Error(w, "Module info not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(module)
	}
}

func updateModuleInfo(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var module models.ModuleInfo
		result := db.First(&module, id)
		if result.Error != nil {
			http.Error(w, "Module info not found", http.StatusNotFound)
			return
		}

		var updatedModule models.ModuleInfo
		err := json.NewDecoder(r.Body).Decode(&updatedModule)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		module.ModuleName = updatedModule.ModuleName
		module.ModuleDuration = updatedModule.ModuleDuration
		module.ExamType = updatedModule.ExamType
		module.Version = updatedModule.Version
		module.UpdatedAt = time.Now()

		db.Save(&module)

		json.NewEncoder(w).Encode(module)
	}
}

func deleteModuleInfo(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var module models.ModuleInfo
		result := db.First(&module, id)
		if result.Error != nil {
			http.Error(w, "Module info not found", http.StatusNotFound)
			return
		}

		db.Delete(&module)

		w.WriteHeader(http.StatusNoContent)
	}
}
