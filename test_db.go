package main

import (
	"fmt"
	"log"
	"tasklybe/internal/user"
	"tasklybe/pkg/db"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	db.Connect()

	var count int64
	db.DB.Model(&user.User{}).Count(&count)
	fmt.Printf("Users in database: %d\n", count)

	var users []user.User
	db.DB.Find(&users)
	for _, u := range users {
		fmt.Printf("- %s (%s)\n", u.Name, u.Email)
	}
}
