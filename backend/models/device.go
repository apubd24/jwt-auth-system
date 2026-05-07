package models

import "time"

//===================

type VendorType string
type DeviceCategory string
type PonType string
type SNMPVersion string

const (
	// Vendor
	VendorBDCOM     VendorType = "BDCOM"
	VendorVSOL      VendorType = "VSOL"
	VendorHUAWEI    VendorType = "HUAWEI"
	VendorZTE       VendorType = "ZTE"
	VendorFIBERHOME VendorType = "FIBERHOME"
	VendorDELL      VendorType = "DELL"

	// Device Category
	CategoryRouter   DeviceCategory = "ROUTER"
	CategorySwitch   DeviceCategory = "SWITCH"
	CategoryFirewall DeviceCategory = "FIREWALL"
	CategoryOLT      DeviceCategory = "OLT"
	CategoryServer   DeviceCategory = "SERVER"

	// PON Type
	PonEPON         PonType = "EPON"
	PonGPON         PonType = "GPON"
	PonXPON         PonType = "XPON"
	PonL2           PonType = "L2"
	PonL3           PonType = "L3"
	ServerPOWEREDGE PonType = "POWEREDGE"

	// SNMP Version
	SNMPv1 SNMPVersion = "v1"
	SNMPv2 SNMPVersion = "v2"
	SNMPv3 SNMPVersion = "V3"
)

// Device Model

type Device struct {
	DeviceID       uint64         `gorm:"primaryKey;autoIncrement;column:device_id"`
	Name           string         `gorm:"not null" json:"name"`
	Serial         string         `gorm:"unique;not null" json:"serial"`
	Description    string         `json:"description"`
	CustomerID     uint64         `gorm:"not null;column:customer_id"`
	CustomerName   string         `gorm:"size:100;not null;column:customer_name"`
	DeviceName     string         `gorm:"size:100;not null;column:device_name"`
	DeviceVendor   VendorType     `gorm:"type:vendor_type;not null;column:device_vendor"`
	DeviceCategory DeviceCategory `gorm:"type:device_category;not null;column:device_category"`
	DeviceType     PonType        `gorm:"type:pon_type;not null;column:device_type"`
	IpAddress      string         `gorm:"size:45;unique;not null;column:ip_address"`
	SnmpCommunity  string         `gorm:"size:100;not null;column:snmp_community"`
	SnmpVersion    SNMPVersion    `gorm:"type:snmp_ver;default:v1;column:snmp_version"`
	IsActive       *bool          `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}
