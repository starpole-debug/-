package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username     string    `gorm:"uniqueIndex"`
	Email        string    `gorm:"uniqueIndex"`
	PasswordHash string
	Nickname     string
	IsAdmin      bool
}

func main() {
	dsn := "host=localhost user=aiapp password=aiapp dbname=aiapp port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&User{})

	password := "admin123"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	u := User{
		Username:     "admin",
		Email:        "admin@example.com",
		PasswordHash: string(hash),
		Nickname:     "Admin",
		IsAdmin:      true,
	}

	var existing User
	if err := db.Where("email = ?", u.Email).First(&existing).Error; err == nil {
		fmt.Println("admin already exists:", existing.Email)
		return
	}

	if err := db.Create(&u).Error; err != nil {
		log.Fatal(err)
	}

	fmt.Println("admin seeded:", u.Email)
}
