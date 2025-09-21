package main

import (
	"go_pz3/internal/database"
	"go_pz3/internal/models"
	"go_pz3/internal/routes"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

	// Отримуємо роутер з налаштованими API-ендпоінтами
	router := routes.SetupRouter(db)

	// Налаштовуємо роздачу статичних файлів (нашого UI)
	router.StaticFS("/ui", http.Dir("./ui"))

	// Робимо редірект з кореня "/" на наш UI
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/ui")
	})

	log.Println("Starting server on port 8080...")
	log.Println("UI available at http://localhost:8080")
	router.Run(":8080")
}
