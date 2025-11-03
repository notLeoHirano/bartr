package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/notLeoHirano/bartr/middleware"
	"github.com/notLeoHirano/bartr/models"
)

func (h *Handler) GetItems(c *gin.Context) {
	userID := middleware.GetUserID(c)
	excludeOwn := c.Query("exclude_own") == "true"

	items, err := h.service.GetItems(userID, excludeOwn)
	if err != nil {
		log.Printf("Error fetching items: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *Handler) CreateItem(c *gin.Context) {
	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
    var verr validator.ValidationErrors
    if errors.As(err, &verr) {
        for _, fe := range verr {
            switch fe.Field() {
            case "Title":
                c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
                return
            }
        }
    }
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
    return
}

	item.UserID = middleware.GetUserID(c)

	if err := h.service.CreateItem(&item); err != nil {
		if err.Error() == "title is required" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Error creating item: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}

	log.Printf("Created item: %d - %s", item.ID, item.Title)
	c.JSON(http.StatusCreated, item)
}

func (h *Handler) DeleteItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	userID := middleware.GetUserID(c)

	if err := h.service.DeleteItem(id, userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Deleted item: %d", id)
	c.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
}