package main

import (
	"fmt"
	"github.com/VolodymyrShabat/Test_ATN/src/api"
	"github.com/VolodymyrShabat/Test_ATN/src/app"
	"github.com/VolodymyrShabat/Test_ATN/src/storage/models"
	"os"
)

func main() {
	services := app.CreateServices(
		[]byte(os.Getenv("JWT_SECRET_KEY")),
		[]byte(os.Getenv("HASH_SALT")))

	server := api.New(
		&models.Config{Env: "info",
			Port:      os.Getenv("SERVER_PORT"),
			Domain:    os.Getenv("DOMAIN"),
			SMTPHost:  os.Getenv("SMTP_HOST"),
			SMTPPass:  os.Getenv("SMTP_PASS"),
			SMTPPort:  os.Getenv("SMTP_PORT"),
			SMTPUser:  os.Getenv("SMTP_USER"),
			EmailFrom: os.Getenv("EMAIL_FROM"),
		}, services)

	err := server.Run()
	if err != nil {
		fmt.Println(err)
	}
}
