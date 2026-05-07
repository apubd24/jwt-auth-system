package handlers

import (
	"errors"
	"jwt-auth-backend/database"
	"jwt-auth-backend/models"
	"net/http"
	"strings"

	// ✅ FIX: import added
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// START GET ALL DEVICE
func GetAllDevices(c *gin.Context) {
	var devices []models.Device
	database.DB.Find(&devices)
	c.JSON(http.StatusOK, gin.H{"devices": devices})
}

func CreateDevice(c *gin.Context) {
	// 1. Input DTO
	var input struct {
		Name           string                `json:"name" binding:"required"`
		Serial         string                `json:"serial" binding:"required"`
		Description    string                `json:"description"`
		CustomerID     uint64                `json:"customer_id" binding:"required"`
		CustomerName   string                `json:"customer_name" binding:"required"`
		DeviceName     string                `json:"device_name" binding:"required"`
		DeviceVendor   models.VendorType     `json:"device_vendor" binding:"required"`
		DeviceCategory models.DeviceCategory `json:"device_category" binding:"required"`
		DeviceType     models.PonType        `json:"device_type" binding:"required"`
		IpAddress      string                `json:"ip_address" binding:"required"`
		SnmpCommunity  string                `json:"snmp_community" binding:"required"`
		SnmpVersion    models.SNMPVersion    `json:"snmp_version"`
		IsActive       *bool                 `json:"is_active"` // Keep as pointer
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Logic to handle default value
	// If input.IsActive is nil (not provided), set it to true.
	// If input.IsActive is provided (true or false), keep it.
	finalIsActive := true
	if input.IsActive != nil {
		finalIsActive = *input.IsActive
	}

	// 3. Create Model
	device := models.Device{
		Name:           input.Name,
		Serial:         input.Serial,
		Description:    input.Description,
		CustomerID:     input.CustomerID,
		CustomerName:   input.CustomerName,
		DeviceName:     input.DeviceName,
		DeviceVendor:   input.DeviceVendor,
		DeviceCategory: input.DeviceCategory,
		DeviceType:     input.DeviceType,
		IpAddress:      input.IpAddress,
		SnmpCommunity:  input.SnmpCommunity,
		SnmpVersion:    input.SnmpVersion,
		// Assign the pointer to the variable
		IsActive: &finalIsActive,
	}

	// 4. Insert
	if err := database.DB.Create(&device).Error; err != nil {

		// 🔥 Handle unique constraint errors cleanly
		if strings.Contains(err.Error(), "serial") {
			c.JSON(http.StatusConflict, gin.H{"error": "Serial '" + input.Serial + "' already exists"})
			return
		}
		if strings.Contains(err.Error(), "ip_address") {
			c.JSON(http.StatusConflict, gin.H{"error": "IP address '" + input.IpAddress + "' already exists"})
			return
		}
		// ... (Keep your existing error handling)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create device"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Device created successfully", "data": device})
}

// END Create or Insert new device

// Start Get device by ID

func GetDeviceByID(c *gin.Context) {

	// Get Device ID From URL
	id := c.Param("id")

	// Device Model
	var device models.Device

	// Find Device
	if err := database.DB.Select(
		"device_id",
		"name",
		"serial",
		"description",
		"customer_id",
		"customer_name",
		"device_name",
		"device_vendor",
		"device_category",
		"device_type",
		"ip_address",
		"snmp_community",
		"snmp_version",
		"is_active",
	).First(&device, "device_id = ?", id).Error; err != nil {

		// Device Not Found
		if errors.Is(err, gorm.ErrRecordNotFound) {

			c.JSON(http.StatusNotFound, gin.H{
				"error": "device not found",
			})
			return
		}

		// Database Error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get device",
		})
		return
	}

	// Success Response
	c.JSON(http.StatusOK, gin.H{
		"message": "device fetched successfully",
		"data":    device,
	})
}

// START Device Update func

func UpdateDevice(c *gin.Context) {

	id := c.Param("id")

	var device models.Device

	// Find Device
	if err := database.DB.First(&device, "device_id = ?", id).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "device not found",
		})
		return
	}

	// Input Struct
	var input struct {
		Name           string                 `json:"name"`
		Serial         string                 `json:"serial"`
		Description    string                 `json:"description"`
		CustomerID     *uint64                `json:"customer_id"`
		CustomerName   string                 `json:"customer_name"`
		DeviceName     string                 `json:"device_name"`
		DeviceVendor   *models.VendorType     `json:"device_vendor"`
		DeviceCategory *models.DeviceCategory `json:"device_category"`
		DeviceType     *models.PonType        `json:"device_type"`
		IpAddress      string                 `json:"ip_address"`
		SnmpCommunity  string                 `json:"snmp_community"`
		SnmpVersion    *models.SNMPVersion    `json:"snmp_version"`
		IsActive       *bool                  `json:"is_active"`
	}

	// Bind JSON
	if err := c.ShouldBindJSON(&input); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Update Fields If Provided

	if input.Name != "" {
		device.Name = input.Name
	}

	if input.Serial != "" {
		device.Serial = input.Serial
	}

	if input.Description != "" {
		device.Description = input.Description
	}

	if input.CustomerID != nil {
		device.CustomerID = *input.CustomerID
	}

	if input.CustomerName != "" {
		device.CustomerName = input.CustomerName
	}

	if input.DeviceName != "" {
		device.DeviceName = input.DeviceName
	}

	if input.DeviceVendor != nil {
		device.DeviceVendor = *input.DeviceVendor
	}

	if input.DeviceCategory != nil {
		device.DeviceCategory = *input.DeviceCategory
	}

	if input.DeviceType != nil {
		device.DeviceType = *input.DeviceType
	}

	if input.IpAddress != "" {
		device.IpAddress = input.IpAddress
	}

	if input.SnmpCommunity != "" {
		device.SnmpCommunity = input.SnmpCommunity
	}

	if input.SnmpVersion != nil {
		device.SnmpVersion = *input.SnmpVersion
	}

	if input.IsActive != nil {
		device.IsActive = input.IsActive
	}

	// Save Device
	if err := database.DB.Model(&device).Updates(device).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update device",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "device updated successfully",
		"data":    device,
	})
}

// END device update func

// // END Create or Insert new device

// func UpdateDevice(c *gin.Context) {
// 	id := c.Param("id")
// 	var device models.Device
// 	if err := database.DB.First(&device, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
// 		return
// 	}
// 	var input struct {
// 		Name        string `json:"name"`
// 		Serial      string `json:"serial"`
// 		Description string `json:"description"`
// 	}
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if input.Name != "" {
// 		device.Name = input.Name
// 	}
// 	if input.Serial != "" {
// 		device.Serial = input.Serial
// 	}
// 	if input.Description != "" {
// 		device.Description = input.Description
// 	}
// 	database.DB.Save(&device)
// 	c.JSON(http.StatusOK, device)
// }

func DeleteDevice(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Device{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "device deleted"})
}
