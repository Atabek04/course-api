package routers

import (
	"course-api/mailer"
	"course-api/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/modules", createModuleInfo(db)).Methods("POST")
	router.HandleFunc("/modules/{id}", getModuleInfo(db)).Methods("GET")
	router.HandleFunc("/modules/{id}", updateModuleInfo(db)).Methods("PUT")
	router.HandleFunc("/modules/{id}", deleteModuleInfo(db)).Methods("DELETE")

	router.HandleFunc("/departments", createDepartmentInfo(db)).Methods("POST")
	router.HandleFunc("/departments/{id}", getDepartmentInfo(db)).Methods("GET")

	router.HandleFunc("/users", registerUserHandler(db)).Methods("POST")
	router.HandleFunc("/users", getAllUserInfoHandler(db)).Methods("GET") // Correct mapping here
	router.HandleFunc("/users/{id}", getUserInfoHandler(db)).Methods("GET")
	router.HandleFunc("/users/{id}", editUserInfoHandler(db)).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteUserInfoHandler(db)).Methods("DELETE")

	return router
}

func registerUserHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		user := &models.User{
			Name:      input.Name,
			Email:     input.Email,
			Activated: false,
		}

		err = user.SetPassword(input.Password)
		if err != nil {
			http.Error(w, "Failed to set password", http.StatusInternalServerError)
			return
		}

		result := db.Create(user)
		if result.Error != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		go func() {
			mailService := mailer.New("sandbox.smtp.mailtrap.io", 587, "3155d87cf6e478", "11ce409c255576", "otabek.shadimatov@gmail.com")
			err = mailService.Send(user.Email, "user_welcome.tmpl", user)
			if err != nil {
				http.Error(w, "Failed to send email", http.StatusInternalServerError)
				return
			}
		}()

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(user)
	}
}

func getUserInfoHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var userInfo models.User
		result := db.First(&userInfo, id)
		if result.Error != nil {
			http.Error(w, "User info not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(userInfo)
	}
}

func getAllUserInfoHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var users []models.User
		query := db.Model(&models.User{})

		// Handle filtering by name if the "name" query parameter is provided.
		if filterParam := r.URL.Query().Get("name"); filterParam != "" {
			query = query.Where("name = ?", filterParam)
		}

		// Handle sorting if the "sort" query parameter is provided.
		if sortBy := r.URL.Query().Get("sort"); sortBy != "" {
			query = query.Order(sortBy)
		}

		result := db.Find(&users)
		if result.Error != nil {
			http.Error(w, "Failed to fetch user information", http.StatusInternalServerError)
			return
		}

		// Return the response in JSON format.
		json.NewEncoder(w).Encode(users)
	}
}

func editUserInfoHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var userInfo models.User
		result := db.First(&userInfo, id)
		if result.Error != nil {
			http.Error(w, "User info not found", http.StatusNotFound)
			return
		}

		var updatedUserInfo models.User
		err := json.NewDecoder(r.Body).Decode(&updatedUserInfo)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		userInfo.Name = updatedUserInfo.Name
		userInfo.Email = updatedUserInfo.Email

		db.Save(&userInfo)

		json.NewEncoder(w).Encode(userInfo)
	}
}

func deleteUserInfoHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var userInfo models.User
		result := db.First(&userInfo, id)
		if result.Error != nil {
			http.Error(w, "User info not found", http.StatusNotFound)
			return
		}

		db.Delete(&userInfo)

		w.WriteHeader(http.StatusNoContent)
	}
}

func createDepartmentInfo(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var department models.DepartmentInfo
		err := json.NewDecoder(r.Body).Decode(&department)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		result := db.Create(&department)
		if result.Error != nil {
			http.Error(w, "Failed to create module info", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(department)
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
