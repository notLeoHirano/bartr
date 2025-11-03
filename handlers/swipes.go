package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/notLeoHirano/bartr/middleware"
	"github.com/notLeoHirano/bartr/models"
)

func (h *Handler) CreateSwipe(c *gin.Context) {
	print()
	var swipe models.Swipe
	if err := c.ShouldBindJSON(&swipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	swipe.UserID = middleware.GetUserID(c)

	if err := h.service.CreateSwipe(&swipe); err != nil {
		if err.Error() == "direction must be 'left' or 'right'" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Error creating swipe: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create swipe"})
		return
	}

	c.JSON(http.StatusCreated, swipe)
}



func (h *Handler) GetMatches(c *gin.Context) {
	userID := middleware.GetUserID(c)

	matches, err := h.service.GetMatches(userID)
	if err != nil {
		log.Printf("Error fetching matches: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch matches"})
		return
	}

	c.JSON(http.StatusOK, matches)
}
