package handlers

import (
	"jwt-auth-backend/database"
	"jwt-auth-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetAllUsers – only admin
func GetAllUsers(c *gin.Context) {
	var users []models.User
	database.DB.Select("id, fullname, username, role, is_active, created_at, updated_at").Find(&users)
	c.JSON(http.StatusOK, gin.H{"users": users})
}

// GetUserByID – fetch single user for edit form
func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	var user models.User

	// Select only required fields (good practice)
	if err := database.DB.
		Select("id, fullname, username, role, is_active, created_at, updated_at").
		First(&user, id).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// CreateUser – admin creates any user (including admin)
func CreateUser(c *gin.Context) {
	var input struct {
		Fullname string `json:"fullname" binding:"required"` // NEW Modifyied"
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role"` // "admin" or "readonly"
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	role := input.Role
	if role != "admin" {
		role = "readonly"
	}
	user := models.User{
		Fullname: input.Fullname,
		Username: input.Username,
		Password: string(hashed),
		Role:     role,
		IsActive: true,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username exists"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "user created", "user": user})
}

// UpdateUser – admin can update role and active status
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	var input struct {
		Fullname string  `json:"fullname" binding:"required"` // NEW Modifyied"
		Role     *string `json:"role"`
		IsActive *bool   `json:"is_active"`
		Password string  `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Fullname != "" {
		user.Fullname = input.Fullname
	}

	if input.Role != nil {
		if *input.Role == "admin" || *input.Role == "readonly" {
			user.Role = *input.Role
		}
	}
	if input.IsActive != nil {
		user.IsActive = *input.IsActive
	}
	if input.Password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		user.Password = string(hashed)
	}
	database.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"message": "user updated", "user": user})
}

// User Password change API
func ChangePassword(c *gin.Context) {
	id := c.Param("id")

	// Get logged-in user from JWT
	currentUserID, _ := c.Get("userID")
	role, _ := c.Get("role")

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// ONLY ADMIN OR OWNER
	if role != "admin" && currentUserID != id {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed"})
		return
	}

	var input struct {
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "password hashing failed"})
		return
	}

	user.Password = string(hashed)
	database.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"message": "password updated successfully",
	})
}

// DeleteUser – hard delete (or soft delete if you prefer)
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
