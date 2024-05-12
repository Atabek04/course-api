package main

import (
	"course-api/cmd/routers"
	"course-api/initializers"
	"course-api/mailer"
	"course-api/models"
	"fmt"
	"github.com/pressly/goose"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
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

	//go CheckActivationAndToken(DB)

	router := routers.NewRouter(DB)
	log.Println("Server is started in port: 3000")
	http.ListenAndServe(":3000", router)
}

func CheckActivationAndToken(db *gorm.DB) {
	for {
		var unactivatedUsers []models.User
		db.Where("activated = ?", false).Find(&unactivatedUsers)

		if len(unactivatedUsers) == 0 {
			fmt.Println("Unactivated users list is empty")
		} else {
			fmt.Println("User list isn't empty")
		}

		for _, user := range unactivatedUsers {
			fmt.Println(user)

			var tokens []models.Token
			db.Where("user_id = ?", user.ID).Find(&tokens)

			for _, token := range tokens {
				if token.Expiry.After(time.Now()) {
					continue
				}

				newToken := generateNewToken(user)
				err := sendEmailWithToken(user, newToken)
				if err != nil {
					fmt.Printf("Error sending email to %s: %v\n", user.Email, err)
					continue
				}
				fmt.Printf("New token sent to user %s\n", user.Name)
			}
		}

		time.Sleep(30 * time.Second)
	}
}

func generateNewToken(user models.User) string {
	newToken, err := models.GenerateToken(user.ID, 3*24*time.Hour, models.ScopeActivation)
	if err != nil {
		log.Printf("Error generating new token for user %d: %v\n", user.ID, err)
	}
	return newToken.Plaintext
}

func sendEmailWithToken(user models.User, token string) error {
	data := map[string]any{
		"activationToken": token,
		"userID":          user.ID}

	mailService := mailer.New("sandbox.smtp.mailtrap.io", 587, "3155d87cf6e478", "11ce409c255576", "otabek.shadimatov@gmail.com")
	err := mailService.Send(user.Email, "user_welcome.tmpl", data)
	if err != nil {
		log.Printf("Failed to send email")
	}
	return err
}
