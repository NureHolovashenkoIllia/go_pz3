package main

import (
	"go_pz3/internal/database"
	"go_pz3/internal/models"
	"go_pz3/internal/routes"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// 1. Завантаження .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. Підключення до БД
	database.Connect()
	db := database.DB

	// 3. Міграція
	log.Println("Running migrations...")
	db.AutoMigrate(&models.Actor{}, &models.Performance{}, &models.Role{}, &models.Rehearsal{})
	log.Println("Migrations completed successfully.")

	// 4. Налаштування та запуск роутера
	router := routes.SetupRouter(db) // <-- Викликаємо нашу функцію
	log.Println("Starting server on port 8080...")
	router.Run(":8080")
}
