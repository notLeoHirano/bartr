package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/notLeoHirano/bartr/database"
	"github.com/notLeoHirano/bartr/handlers"
	"github.com/notLeoHirano/bartr/middleware"
	"github.com/notLeoHirano/bartr/service"
	"github.com/notLeoHirano/bartr/store"
)

func main() {
	// Initialize database
	db, err := database.New("./bartr.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	if err := db.Init(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize layers
	st := store.New(db.DB)
	svc := service.New(st)
	handler := handlers.New(svc)

	// Setup router
	r := gin.Default()
	
	// CORS must be first
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * 3600,
	}))

	// Public routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
	}

	// Protected routes
	api := r.Group("/")
	api.Use(middleware.AuthRequired())
	{
		// User
		api.GET("/me", handler.GetMe)

		// Items
		api.GET("/items", handler.GetItems)
		api.POST("/items", handler.CreateItem)
		api.DELETE("/items/:id", handler.DeleteItem)

		// Swipes
		api.POST("/swipes", handler.CreateSwipe)

		// Matches
		api.GET("/matches", handler.GetMatches)

		// Comments
		api.POST("/comments", handler.CreateComment)
		api.GET("/matches/:match_id/comments", handler.GetComments)
	}

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}