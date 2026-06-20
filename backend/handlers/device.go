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

// Custom struct for dropdown response (only ID and DeviceName)
type DeviceDropdown struct {
	ID            uint   `json:"id"`
	ItemID        string `json:"item_id"`
	DeviceName    string `json:"device_name"`
	Serial        string `json:"serial"`
	DeviceVendor  string `json:"device_vendor"`
	HardwareModel string `json:"hardware_model"`
}

func GetDeviceDropdown(c *gin.Context) {

	var devices []DeviceDropdown

	err := database.DB.
		Model(&models.Device{}).
		Select("device_id as id, item_id, device_name, serial, device_vendor, hardware_model").
		Order("device_name ASC").
		Find(&devices).Error

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data": devices,
	})
}

func CreateDevice(c *gin.Context) {
	// 1. Input DTO
	var input struct {
		ItemID      string `json:"item_id" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Serial      string `json:"serial" binding:"required"`
		Description string `json:"description"`
		CustomerID  uint64 `json:"customer_id" binding:"required"`
		// CustomerName      string                `json:"customer_name"`
		DeviceName        string                `json:"device_name" binding:"required"`
		DeviceVendor      models.VendorType     `json:"device_vendor" binding:"required"`
		DeviceCategory    models.DeviceCategory `json:"device_category" binding:"required"`
		HardwareModel     models.HardwareModel  `json:"hardware_model" binding:"required"`
		DeviceType        string                `json:"device_type" binding:"required"`
		HardwareHeight    string                `json:"hardware_height"`
		SnmpGroup         models.SnmpGroup      `json:"snmp_group" binding:"required"`
		AssetStatus       string                `json:"asset_status" binding:"required"`
		DatacenterName    string                `json:"datacenter_name"`
		RackName          string                `json:"rack_name"`
		RackPosition      string                `json:"rack_position"`
		TeamStoreLocation string                `json:"team_store_location"`
		IpAddress         string                `json:"ip_address" binding:"required"`
		SnmpCommunity     string                `json:"snmp_community" binding:"required"`
		SnmpVersion       models.SNMPVersion    `json:"snmp_version"`
		IsActive          *bool                 `json:"is_active"` // Keep as pointer
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	finalIsActive := true
	if input.IsActive != nil {
		finalIsActive = *input.IsActive
	}

	// 3. Create Model
	device := models.Device{
		ItemID:      input.ItemID,
		Name:        input.Name,
		Serial:      input.Serial,
		Description: input.Description,
		CustomerID:  input.CustomerID,
		// CustomerName:      input.CustomerName,
		DeviceName:        input.DeviceName,
		DeviceVendor:      input.DeviceVendor,
		DeviceCategory:    input.DeviceCategory,
		HardwareModel:     input.HardwareModel,
		DeviceType:        input.DeviceType,
		HardwareHeight:    input.HardwareHeight,
		SnmpGroup:         input.SnmpGroup,
		AssetStatus:       input.AssetStatus,
		DatacenterName:    input.DatacenterName,
		RackName:          input.RackName,
		RackPosition:      input.RackPosition,
		TeamStoreLocation: input.TeamStoreLocation,
		IpAddress:         input.IpAddress,
		SnmpCommunity:     input.SnmpCommunity,
		SnmpVersion:       input.SnmpVersion,
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
		"item_id",
		"name",
		"serial",
		"description",
		"customer_id",
		// "customer_name",
		"device_name",
		"device_vendor",
		"device_category",
		"hardware_model",
		"device_type",
		"hardware_height",
		"snmp_group",
		"asset_status",
		"datacenter_name",
		"rack_name",
		"rack_position",
		"team_store_location",
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
		Name        string  `json:"name"`
		Serial      string  `json:"serial"`
		Description string  `json:"description"`
		CustomerID  *uint64 `json:"customer_id"`
		// CustomerName      string                 `json:"customer_name"`
		DeviceName        string                 `json:"device_name"`
		DeviceVendor      *models.VendorType     `json:"device_vendor"`
		DeviceCategory    *models.DeviceCategory `json:"device_category"`
		HardwareModel     *models.HardwareModel  `json:"hardware_model"`
		DeviceType        string                 `json:"device_type"`
		HardwareHeight    string                 `json:"hardware_height"`
		SnmpGroup         *models.SnmpGroup      `json:"snmp_group"`
		AssetStatus       string                 `json:"asset_status"`
		DatacenterName    string                 `json:"datacenter_name"`
		RackName          string                 `json:"rack_name"`
		RackPosition      string                 `json:"rack_position"`
		TeamStoreLocation string                 `json:"team_store_location"`
		IpAddress         string                 `json:"ip_address"`
		SnmpCommunity     string                 `json:"snmp_community"`
		SnmpVersion       *models.SNMPVersion    `json:"snmp_version"`
		IsActive          *bool                  `json:"is_active"`
	}

	// Bind JSON
	if err := c.ShouldBindJSON(&input); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Update Fields If Provided

	if input.Serial != "" {
		device.Serial = input.Serial
	}

	if input.Description != "" {
		device.Description = input.Description
	}

	if input.CustomerID != nil {
		device.CustomerID = *input.CustomerID
	}

	// if input.CustomerName != "" {
	// 	device.CustomerName = input.CustomerName
	// }

	if input.DeviceName != "" {
		device.DeviceName = input.DeviceName
	}

	if input.DeviceVendor != nil {
		device.DeviceVendor = *input.DeviceVendor
	}

	if input.DeviceCategory != nil {
		device.DeviceCategory = *input.DeviceCategory
	}

	if input.HardwareModel != nil {
		device.HardwareModel = *input.HardwareModel
	}

	if input.DeviceType != "" {
		device.DeviceType = input.DeviceType
	}

	if input.HardwareHeight != "" {
		device.HardwareHeight = input.HardwareHeight
	}

	if input.SnmpGroup != nil {
		device.SnmpGroup = *input.SnmpGroup
	}

	if input.AssetStatus != "" {

		device.AssetStatus = input.AssetStatus

		switch input.AssetStatus {

		case "Live":

			device.DatacenterName = input.DatacenterName
			device.RackName = input.RackName
			device.RackPosition = input.RackPosition

			device.TeamStoreLocation = ""

		case "Available":

			device.DatacenterName = input.DatacenterName
			device.RackName = input.RackName
			device.RackPosition = input.RackPosition

			device.TeamStoreLocation = ""

		case "Team Store":

			device.DatacenterName = ""
			device.RackName = ""
			device.RackPosition = ""

			device.TeamStoreLocation = input.TeamStoreLocation

		default:

			device.DatacenterName = ""
			device.RackName = ""
			device.RackPosition = ""

			device.TeamStoreLocation = ""
		}
	}

	// if input.DatacenterName != "" {
	// 	device.DatacenterName = input.DatacenterName
	// }

	// if input.RackName != "" {
	// 	device.RackName = input.RackName
	// }

	// if input.RackPosition != "" {
	// 	device.RackPosition = input.RackPosition
	// }

	// if input.TeamStoreLocation != "" {
	// 	device.TeamStoreLocation = input.TeamStoreLocation
	// }

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
	updateData := map[string]interface{}{
		"name":        device.Name,
		"serial":      device.Serial,
		"description": device.Description,
		"customer_id": device.CustomerID,
		// "customer_name":       device.CustomerName,
		"device_name":         device.DeviceName,
		"device_vendor":       device.DeviceVendor,
		"device_category":     device.DeviceCategory,
		"hardware_model":      device.HardwareModel,
		"device_type":         device.DeviceType,
		"hardware_height":     device.HardwareHeight,
		"snmp_group":          device.SnmpGroup,
		"asset_status":        device.AssetStatus,
		"datacenter_name":     device.DatacenterName,
		"rack_name":           device.RackName,
		"rack_position":       device.RackPosition,
		"team_store_location": device.TeamStoreLocation,
		"ip_address":          device.IpAddress,
		"snmp_community":      device.SnmpCommunity,
		"snmp_version":        device.SnmpVersion,
		"is_active":           device.IsActive,
	}

	if err := database.DB.
		Model(&device).
		Updates(updateData).Error; err != nil {

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

func DeleteDevice(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Device{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "device deleted"})
}
