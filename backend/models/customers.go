package models

import "time"

// Helper function to convert ContactPerson DB model to ContactPersonOutput DTO
// Force GORM to use exactly "contact_people" as your table name
func (ContactPerson) TableName() string {
	return "contact_people"
}

// ─────────────────────────────────────────────
// DB MODELS
// ─────────────────────────────────────────────

type ContactPerson struct {
	ID           uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	CustomerID   uint   `json:"customer_id"`
	Name         string `json:"name"`
	Designation  string `json:"designation"`
	ContactType  string `json:"contact_type"`
	ContactLevel string `json:"contact_level"`
	Emails       string `json:"emails" gorm:"type:text"` // JSON array string
	Mobiles      string `json:"mobiles" gorm:"type:text"`
	Whatsapps    string `json:"whatsapps" gorm:"type:text"`
}

type Customer struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Company Details
	CompanyName    string `json:"company_name"`
	CompanyLogo    string `json:"company_logo"`
	Website        string `json:"website"`
	OfficeAddress  string `json:"office_address"`
	CustomerNote   string `json:"customer_note"`
	CustomerStatus string `json:"customer_status" gorm:"default:active"`

	// Account Manager
	AccountManagerName        string `json:"account_manager_name"`
	AccountManagerDesignation string `json:"account_manager_designation"`
	AccountManagerEmail       string `json:"account_manager_email"`
	AccountManagerContact     string `json:"account_manager_contact"`
	AccountManagerWhatsapp    string `json:"account_manager_whatsapp"`
	AccountManagerBranch      string `json:"account_manager_branch"`

	// Support contacts (JSON arrays stored as strings)
	SupportEmails          string `json:"support_emails" gorm:"type:text"`
	SupportMobiles         string `json:"support_mobiles" gorm:"type:text"`
	SupportWhatsappNumbers string `json:"support_whatsapp_numbers" gorm:"type:text"`
	SupportWhatsappGroups  string `json:"support_whatsapp_groups" gorm:"type:text"`

	// Relations
	ContactPersons []ContactPerson `json:"contact_persons" gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"`
}

// ─────────────────────────────────────────────
// REQUEST DTOs
// ─────────────────────────────────────────────

// FIX: ID must be included here.
// When React sends contact_persons JSON, each person carries its DB id.
// id=0  → Go will INSERT a new record
// id>0  → Go will UPDATE the existing record
// Any DB record whose id is NOT in the list → Go will DELETE it
type ContactPersonInput struct {
	ID           uint     `json:"id"` // ← THE KEY FIX: was missing, causing all persons to be treated as new
	Name         string   `json:"name"`
	Designation  string   `json:"designation"`
	ContactType  string   `json:"contact_type"`
	ContactLevel string   `json:"contact_level"`
	Emails       []string `json:"emails"`
	Mobiles      []string `json:"mobiles"`
	Whatsapps    []string `json:"whatsapps"`
}

// ─────────────────────────────────────────────
// RESPONSE DTOs
// ─────────────────────────────────────────────

type ContactPersonOutput struct {
	ID           uint     `json:"id"`
	Name         string   `json:"name"`
	Designation  string   `json:"designation"`
	ContactType  string   `json:"contact_type"`
	ContactLevel string   `json:"contact_level"`
	Emails       []string `json:"emails"`
	Mobiles      []string `json:"mobiles"`
	Whatsapps    []string `json:"whatsapps"`
}

type CustomerResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	CompanyName    string `json:"company_name"`
	CompanyLogo    string `json:"company_logo"`
	Website        string `json:"website"`
	OfficeAddress  string `json:"office_address"`
	CustomerNote   string `json:"customer_note"`
	CustomerStatus string `json:"customer_status"`

	AccountManagerName        string `json:"account_manager_name"`
	AccountManagerDesignation string `json:"account_manager_designation"`
	AccountManagerEmail       string `json:"account_manager_email"`
	AccountManagerContact     string `json:"account_manager_contact"`
	AccountManagerWhatsapp    string `json:"account_manager_whatsapp"`
	AccountManagerBranch      string `json:"account_manager_branch"`

	SupportEmails          []string              `json:"support_emails"`
	SupportMobiles         []string              `json:"support_mobiles"`
	SupportWhatsappNumbers []string              `json:"support_whatsapp_numbers"`
	SupportWhatsappGroups  []string              `json:"support_whatsapp_groups"`
	ContactPersons         []ContactPersonOutput `json:"contact_persons"`
}
