package routers

import (
	"bytes"
	"course-api/models"
	"fmt"
	"github.com/pressly/goose"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var (
	DB *gorm.DB
)

func init() {
	var err error
	//dsn := os.Getenv("DB_URL")
	dsn := os.Getenv("DB_TEST_URL")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Failed to connect ot database")
	}

	db, _ := DB.DB()
	if err := goose.Up(db, "../../migrations"); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	fmt.Println("Database migration completed successfully")
}

func TestRegisterUser(t *testing.T) {
	payload := []byte(`{
		"name": "Lol Doe",
		"email": "lol222@gmail.com",
		"password": "password123"
	}`)

	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(registerUserHandler(DB))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %v, got %v", http.StatusCreated, rr.Code)
	}
}

func TestGetAllUsers(t *testing.T) {
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getAllUserInfoHandler(DB))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}

	t.Logf("Response Body: %s", rr.Body.String())
}

func TestGetUserInfo(t *testing.T) {
	req, err := http.NewRequest("GET", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(getUserInfoHandler(DB))
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}

	t.Logf("Response Body: %s", rr.Body.String())
}

func TestUpdateUser(t *testing.T) {
	payload := []byte(`{
		"name": "Updated Name",
		"email": "updatedemail@example.com"
	}`)

	req, err := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(editUserInfoHandler(DB))
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}
}

func TestDeleteUserInfo(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/users/7", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(deleteUserInfoHandler(DB))
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status %v, got %v", http.StatusNoContent, rr.Code)
	}
}

func TestGenerateToken(t *testing.T) {
	token, err := models.TokenModel{}.New(DB, 99, 3*24*time.Hour, models.ScopeActivation)
	if err != nil {
		t.Fatal(err)
	}

	if token.Plaintext != "" {
		t.Errorf("token is empty '%s'", token.Plaintext)
	}
}
