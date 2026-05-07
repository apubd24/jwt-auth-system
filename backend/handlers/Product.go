package handlers

import (
	"jwt-auth-backend/database"
	"jwt-auth-backend/models"
	"net/http"

	// ✅ FIX: import added
	"github.com/gin-gonic/gin"
)

// START GET ALL DEVICE
func GetAllProducts(c *gin.Context) {
	var Products []models.Product
	database.DB.Find(&Products)
	c.JSON(http.StatusOK, gin.H{"Products": Products})
}

// ===========OLD Configuration==============

func CreateProduct(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Serial      string `json:"serial" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product := models.Product{
		Name:        input.Name,
		Serial:      input.Serial,
		Description: input.Description,
	}
	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "serial already exists"})
		return
	}
	c.JSON(http.StatusCreated, product)
}
