package models

import "time"

//===================

type VendorType string
type HardwareModel string
type DeviceCategory string
type SnmpGroup string
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

	//Hardware Models
	HardwareModelDell  HardwareModel = "Dell PowerEdge R630"
	HardwareModelDell1 HardwareModel = "Dell PowerEdge R640"
	HardwareModelDell2 HardwareModel = "Dell EMC PowerEdge R740"

	// SNMP Group
	PonEPON         SnmpGroup = "EPON"
	PonGPON         SnmpGroup = "GPON"
	PonXPON         SnmpGroup = "XPON"
	PonL2           SnmpGroup = "L2"
	PonL3           SnmpGroup = "L3"
	ServerPOWEREDGE SnmpGroup = "POWEREDGE"

	// SNMP Version
	SNMPv1 SNMPVersion = "v1"
	SNMPv2 SNMPVersion = "v2c"
	SNMPv3 SNMPVersion = "v3"
)

// Device Model

type Device struct {
	DeviceID    uint64 `gorm:"primaryKey;autoIncrement;column:device_id"`
	ItemID      string `gorm:"unique;not null;column:item_id"`
	Name        string `gorm:"not null" json:"name"`
	Serial      string `gorm:"unique;not null" json:"serial"`
	Description string `json:"description"`
	CustomerID  uint64 `gorm:"not null;column:customer_id"`
	// CustomerName      string         `gorm:"size:100;column:customer_name"`
	DeviceName        string         `gorm:"size:100;not null;column:device_name"`
	DeviceVendor      VendorType     `gorm:"type:vendor_type;not null;column:device_vendor"`
	DeviceCategory    DeviceCategory `gorm:"type:device_category;not null;column:device_category"`
	HardwareModel     HardwareModel  `gorm:"type:hardware_model;not null;column:hardware_model"`
	DeviceType        string         `gorm:"not null;column:device_type"`
	HardwareHeight    string         `gorm:"size:150;not null;column:hardware_height"`
	SnmpGroup         SnmpGroup      `gorm:"type:snmp_group;not null;column:snmp_group"`
	AssetStatus       string         `gorm:"size:200;not null;column:asset_status"`
	DatacenterName    string         `gorm:"size:100;column:datacenter_name"`
	RackName          string         `gorm:"size:100;column:rack_name"`
	TeamStoreLocation string         `gorm:"size:100;column:team_store_location"`
	RackPosition      string         `gorm:"size:100;column:rack_position"`
	IpAddress         string         `gorm:"size:45;unique;not null;column:ip_address"`
	SnmpCommunity     string         `gorm:"size:100;not null;column:snmp_community"`
	SnmpVersion       SNMPVersion    `gorm:"type:snmp_ver;default:v1;column:snmp_version"`
	IsActive          *bool          `gorm:"default:true" json:"is_active"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}
