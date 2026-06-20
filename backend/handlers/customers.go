// package handlers

// import (
// 	"encoding/json"
// 	"fmt"
// 	"jwt-auth-backend/database"
// 	"jwt-auth-backend/models"
// 	"net/http"
// 	"os"
// 	"path/filepath"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/gin-gonic/gin"
// )

// // ─────────────────────────────────────────────
// // HELPERS
// // ─────────────────────────────────────────────

// func mustJSON(v []string) string {
// 	b, _ := json.Marshal(v)
// 	return string(b)
// }

// func parseJSON(s string) []string {
// 	var result []string
// 	if err := json.Unmarshal([]byte(s), &result); err != nil {
// 		return []string{}
// 	}
// 	return result
// }

// func toCustomerResponse(c models.Customer) models.CustomerResponse {
// 	persons := make([]models.ContactPersonOutput, 0, len(c.ContactPersons))
// 	for _, p := range c.ContactPersons {
// 		persons = append(persons, models.ContactPersonOutput{
// 			ID:           p.ID,
// 			Name:         p.Name,
// 			Designation:  p.Designation,
// 			ContactType:  p.ContactType,
// 			ContactLevel: p.ContactLevel,
// 			Emails:       parseJSON(p.Emails),
// 			Mobiles:      parseJSON(p.Mobiles),
// 			Whatsapps:    parseJSON(p.Whatsapps),
// 		})
// 	}
// 	return models.CustomerResponse{
// 		ID:                        c.ID,
// 		CreatedAt:                 c.CreatedAt,
// 		UpdatedAt:                 c.UpdatedAt,
// 		CompanyName:               c.CompanyName,
// 		CompanyLogo:               c.CompanyLogo,
// 		Website:                   c.Website,
// 		OfficeAddress:             c.OfficeAddress,
// 		CustomerNote:              c.CustomerNote,
// 		CustomerStatus:            c.CustomerStatus,
// 		AccountManagerName:        c.AccountManagerName,
// 		AccountManagerDesignation: c.AccountManagerDesignation,
// 		AccountManagerEmail:       c.AccountManagerEmail,
// 		AccountManagerContact:     c.AccountManagerContact,
// 		AccountManagerWhatsapp:    c.AccountManagerWhatsapp,
// 		AccountManagerBranch:      c.AccountManagerBranch,
// 		SupportEmails:             parseJSON(c.SupportEmails),
// 		SupportMobiles:            parseJSON(c.SupportMobiles),
// 		SupportWhatsappNumbers:    parseJSON(c.SupportWhatsappNumbers),
// 		SupportWhatsappGroups:     parseJSON(c.SupportWhatsappGroups),
// 		ContactPersons:            persons,
// 	}
// }

// func saveLogoFile(c *gin.Context) (string, error) {
// 	file, err := c.FormFile("company_logo")
// 	if err != nil {
// 		return "", nil // logo is optional
// 	}
// 	ext := strings.ToLower(filepath.Ext(file.Filename))
// 	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
// 	if !allowed[ext] {
// 		return "", fmt.Errorf("unsupported image type: %s", ext)
// 	}
// 	_ = os.MkdirAll("uploads", 0755)
// 	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
// 	dst := filepath.Join("uploads", filename)
// 	if err := c.SaveUploadedFile(file, dst); err != nil {
// 		return "", err
// 	}
// 	return "/uploads/" + filename, nil
// }

// func parseContactPersons(raw string) ([]models.ContactPersonInput, error) {
// 	var persons []models.ContactPersonInput
// 	if err := json.Unmarshal([]byte(raw), &persons); err != nil {
// 		return nil, err
// 	}
// 	return persons, nil
// }

// // ─────────────────────────────────────────────
// // HANDLERS
// // ─────────────────────────────────────────────

// // GET /api/customers
// func ListCustomers(c *gin.Context) {
// 	// FIX: was mistakenly typed as []models.ContactPerson — must be []models.Customer
// 	var customers []models.Customer
// 	query := database.DB.Preload("ContactPersons")

// 	if status := c.Query("status"); status != "" {
// 		query = query.Where("customer_status = ?", status)
// 	}
// 	if search := c.Query("search"); search != "" {
// 		like := "%" + search + "%"
// 		query = query.Where(
// 			"company_name ILIKE ? OR account_manager_name ILIKE ? OR account_manager_email ILIKE ?",
// 			like, like, like,
// 		)
// 	}

// 	if err := query.Find(&customers).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customers"})
// 		return
// 	}

// 	out := make([]models.CustomerResponse, 0, len(customers))
// 	for _, cu := range customers {
// 		out = append(out, toCustomerResponse(cu))
// 	}
// 	c.JSON(http.StatusOK, gin.H{"data": out, "total": len(out)})
// }

// // GET /api/customers/:id
// func GetCustomer(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
// 		return
// 	}
// 	var customer models.Customer
// 	if err := database.DB.Preload("ContactPersons").First(&customer, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, toCustomerResponse(customer))
// }

// // POST /api/customers   (multipart/form-data)
// func CreateCustomer(c *gin.Context) {
// 	logoPath, err := saveLogoFile(c)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	persons, err := parseContactPersons(c.PostForm("contact_persons"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact_persons JSON"})
// 		return
// 	}

// 	dbPersons := make([]models.ContactPerson, 0, len(persons))
// 	for _, p := range persons {
// 		dbPersons = append(dbPersons, models.ContactPerson{
// 			Name:         p.Name,
// 			Designation:  p.Designation,
// 			ContactType:  p.ContactType,
// 			ContactLevel: p.ContactLevel,
// 			Emails:       mustJSON(p.Emails),
// 			Mobiles:      mustJSON(p.Mobiles),
// 			Whatsapps:    mustJSON(p.Whatsapps),
// 		})
// 	}

// 	customer := models.Customer{
// 		CompanyName:               c.PostForm("company_name"),
// 		CompanyLogo:               logoPath,
// 		Website:                   c.PostForm("website"),
// 		OfficeAddress:             c.PostForm("office_address"),
// 		CustomerNote:              c.PostForm("customer_note"),
// 		CustomerStatus:            c.PostForm("customer_status"),
// 		AccountManagerName:        c.PostForm("account_manager_name"),
// 		AccountManagerDesignation: c.PostForm("account_manager_designation"),
// 		AccountManagerEmail:       c.PostForm("account_manager_email"),
// 		AccountManagerContact:     c.PostForm("account_manager_contact"),
// 		AccountManagerWhatsapp:    c.PostForm("account_manager_whatsapp"),
// 		AccountManagerBranch:      c.PostForm("account_manager_branch"),
// 		SupportEmails:             c.PostForm("support_emails"),
// 		SupportMobiles:            c.PostForm("support_mobiles"),
// 		SupportWhatsappNumbers:    c.PostForm("support_whatsapp_numbers"),
// 		SupportWhatsappGroups:     c.PostForm("support_whatsapp_groups"),
// 		ContactPersons:            dbPersons,
// 	}

// 	if customer.CustomerStatus == "" {
// 		customer.CustomerStatus = "active"
// 	}

// 	if err := database.DB.Create(&customer).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
// 		return
// 	}
// 	c.JSON(http.StatusCreated, toCustomerResponse(customer))
// }

// // PUT /api/customers/:id   (multipart/form-data)
// func UpdateCustomer(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
// 		return
// 	}

// 	// Fetch existing customer with contacts
// 	var existing models.Customer
// 	if err := database.DB.Preload("ContactPersons").First(&existing, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
// 		return
// 	}

// 	// Handle logo update
// 	logoPath := existing.CompanyLogo
// 	newLogo, err := saveLogoFile(c)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if newLogo != "" {
// 		if existing.CompanyLogo != "" {
// 			_ = os.Remove(strings.TrimPrefix(existing.CompanyLogo, "/"))
// 		}
// 		logoPath = newLogo
// 	}

// 	// Parse contact persons JSON
// 	persons, err := parseContactPersons(c.PostForm("contact_persons"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact_persons JSON"})
// 		return
// 	}

// 	tx := database.DB.Begin()
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 			panic(r)
// 		}
// 	}()

// 	// Map to quickly find existing contacts by ID or by (Name+Designation)
// 	existingMap := make(map[uint]models.ContactPerson)
// 	for _, cp := range existing.ContactPersons {
// 		existingMap[cp.ID] = cp
// 	}

// 	var incomingIDs []uint

// 	for _, p := range persons {
// 		var contact models.ContactPerson

// 		// Case 1: Incoming contact has a valid ID -> update existing
// 		if p.ID > 0 {
// 			if _, ok := existingMap[p.ID]; ok {
// 				// Update
// 				contact = models.ContactPerson{
// 					ID:           p.ID,
// 					CustomerID:   uint(id),
// 					Name:         p.Name,
// 					Designation:  p.Designation,
// 					ContactType:  p.ContactType,
// 					ContactLevel: p.ContactLevel,
// 					Emails:       mustJSON(p.Emails),
// 					Mobiles:      mustJSON(p.Mobiles),
// 					Whatsapps:    mustJSON(p.Whatsapps),
// 				}
// 				if err := tx.Model(&models.ContactPerson{}).Where("id = ? AND customer_id = ?", p.ID, id).Updates(contact).Error; err != nil {
// 					tx.Rollback()
// 					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update contact"})
// 					return
// 				}
// 				incomingIDs = append(incomingIDs, p.ID)
// 				continue
// 			}
// 			// ID provided but not found – treat as new (insert)
// 		}

// 		// Case 2: No ID or ID not found – try to find by matching name+designation (optional)
// 		// (Prevents duplicate creation when frontend forgets to send ID)
// 		var found bool
// 		var matchedID uint
// 		for _, cp := range existing.ContactPersons {
// 			if cp.Name == p.Name && cp.Designation == p.Designation {
// 				found = true
// 				matchedID = cp.ID
// 				break
// 			}
// 		}
// 		if found {
// 			// Update the matched contact
// 			contact = models.ContactPerson{
// 				ID:           matchedID,
// 				CustomerID:   uint(id),
// 				Name:         p.Name,
// 				Designation:  p.Designation,
// 				ContactType:  p.ContactType,
// 				ContactLevel: p.ContactLevel,
// 				Emails:       mustJSON(p.Emails),
// 				Mobiles:      mustJSON(p.Mobiles),
// 				Whatsapps:    mustJSON(p.Whatsapps),
// 			}
// 			if err := tx.Model(&models.ContactPerson{}).Where("id = ? AND customer_id = ?", matchedID, id).Updates(contact).Error; err != nil {
// 				tx.Rollback()
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update matched contact"})
// 				return
// 			}
// 			incomingIDs = append(incomingIDs, matchedID)
// 			continue
// 		}

// 		// Case 3: New contact (insert)
// 		contact = models.ContactPerson{
// 			CustomerID:   uint(id),
// 			Name:         p.Name,
// 			Designation:  p.Designation,
// 			ContactType:  p.ContactType,
// 			ContactLevel: p.ContactLevel,
// 			Emails:       mustJSON(p.Emails),
// 			Mobiles:      mustJSON(p.Mobiles),
// 			Whatsapps:    mustJSON(p.Whatsapps),
// 		}
// 		if err := tx.Create(&contact).Error; err != nil {
// 			tx.Rollback()
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create contact"})
// 			return
// 		}
// 		incomingIDs = append(incomingIDs, contact.ID)
// 	}

// 	// Delete contacts that are not in the incoming list
// 	if len(incomingIDs) > 0 {
// 		if err := tx.Where("customer_id = ? AND id NOT IN ?", id, incomingIDs).Delete(&models.ContactPerson{}).Error; err != nil {
// 			tx.Rollback()
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove old contacts"})
// 			return
// 		}
// 	} else {
// 		// No contacts sent → delete all existing contacts for this customer
// 		if err := tx.Where("customer_id = ?", id).Delete(&models.ContactPerson{}).Error; err != nil {
// 			tx.Rollback()
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear all contacts"})
// 			return
// 		}
// 	}

// 	// Update customer main fields
// 	updates := map[string]interface{}{
// 		"company_name":                c.PostForm("company_name"),
// 		"company_logo":                logoPath,
// 		"website":                     c.PostForm("website"),
// 		"office_address":              c.PostForm("office_address"),
// 		"customer_note":               c.PostForm("customer_note"),
// 		"customer_status":             c.PostForm("customer_status"),
// 		"account_manager_name":        c.PostForm("account_manager_name"),
// 		"account_manager_designation": c.PostForm("account_manager_designation"),
// 		"account_manager_email":       c.PostForm("account_manager_email"),
// 		"account_manager_contact":     c.PostForm("account_manager_contact"),
// 		"account_manager_whatsapp":    c.PostForm("account_manager_whatsapp"),
// 		"account_manager_branch":      c.PostForm("account_manager_branch"),
// 		"support_emails":              c.PostForm("support_emails"),
// 		"support_mobiles":             c.PostForm("support_mobiles"),
// 		"support_whatsapp_numbers":    c.PostForm("support_whatsapp_numbers"),
// 		"support_whatsapp_groups":     c.PostForm("support_whatsapp_groups"),
// 	}
// 	if err := tx.Model(&existing).Updates(updates).Error; err != nil {
// 		tx.Rollback()
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer"})
// 		return
// 	}

// 	// Commit transaction
// 	if err := tx.Commit().Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed"})
// 		return
// 	}

// 	// Reload updated customer
// 	database.DB.Preload("ContactPersons").First(&existing, id)
// 	c.JSON(http.StatusOK, toCustomerResponse(existing))
// }

// // DELETE /api/customers/:id
// func DeleteCustomer(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
// 		return
// 	}
// 	var customer models.Customer
// 	if err := database.DB.First(&customer, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
// 		return
// 	}
// 	if customer.CompanyLogo != "" {
// 		_ = os.Remove(strings.TrimPrefix(customer.CompanyLogo, "/"))
// 	}
// 	if err := database.DB.Delete(&customer).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
// }

// package handlers

// import (
// 	"encoding/json"
// 	"fmt"
// 	"jwt-auth-backend/database"
// 	"jwt-auth-backend/models"
// 	"net/http"
// 	"os"
// 	"path/filepath"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm" // Added missing gorm import for gorm.Session
// )

// // ─────────────────────────────────────────────
// // HELPERS
// // ─────────────────────────────────────────────

// func mustJSON(v []string) string {
// 	b, _ := json.Marshal(v)
// 	return string(b)
// }

// func parseJSON(s string) []string {
// 	var result []string
// 	if err := json.Unmarshal([]byte(s), &result); err != nil {
// 		return []string{}
// 	}
// 	return result
// }

// func toCustomerResponse(c models.Customer) models.CustomerResponse {
// 	persons := make([]models.ContactPersonOutput, 0, len(c.ContactPersons))
// 	for _, p := range c.ContactPersons {
// 		persons = append(persons, models.ContactPersonOutput{
// 			ID:           p.ID,
// 			Name:         p.Name,
// 			Designation:  p.Designation,
// 			ContactType:  p.ContactType,
// 			ContactLevel: p.ContactLevel,
// 			Emails:       parseJSON(p.Emails),
// 			Mobiles:      parseJSON(p.Mobiles),
// 			Whatsapps:    parseJSON(p.Whatsapps),
// 		})
// 	}
// 	return models.CustomerResponse{
// 		ID:                        c.ID,
// 		CreatedAt:                 c.CreatedAt,
// 		UpdatedAt:                 c.UpdatedAt,
// 		CompanyName:               c.CompanyName,
// 		CompanyLogo:               c.CompanyLogo,
// 		Website:                   c.Website,
// 		OfficeAddress:             c.OfficeAddress,
// 		CustomerNote:              c.CustomerNote,
// 		CustomerStatus:            c.CustomerStatus,
// 		AccountManagerName:        c.AccountManagerName,
// 		AccountManagerDesignation: c.AccountManagerDesignation,
// 		AccountManagerEmail:       c.AccountManagerEmail,
// 		AccountManagerContact:     c.AccountManagerContact,
// 		AccountManagerWhatsapp:    c.AccountManagerWhatsapp,
// 		AccountManagerBranch:      c.AccountManagerBranch,
// 		SupportEmails:             parseJSON(c.SupportEmails),
// 		SupportMobiles:            parseJSON(c.SupportMobiles),
// 		SupportWhatsappNumbers:    parseJSON(c.SupportWhatsappNumbers),
// 		SupportWhatsappGroups:     parseJSON(c.SupportWhatsappGroups),
// 		ContactPersons:            persons,
// 	}
// }

// func saveLogoFile(c *gin.Context) (string, error) {
// 	file, err := c.FormFile("company_logo")
// 	if err != nil {
// 		return "", nil // logo is optional
// 	}
// 	ext := strings.ToLower(filepath.Ext(file.Filename))
// 	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
// 	if !allowed[ext] {
// 		return "", fmt.Errorf("unsupported image type: %s", ext)
// 	}
// 	_ = os.MkdirAll("uploads", 0755)
// 	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
// 	dst := filepath.Join("uploads", filename)
// 	if err := c.SaveUploadedFile(file, dst); err != nil {
// 		return "", err
// 	}
// 	return "/uploads/" + filename, nil
// }

// func parseContactPersons(raw string) ([]models.ContactPersonInput, error) {
// 	var persons []models.ContactPersonInput
// 	if raw == "" || raw == "[]" {
// 		return persons, nil
// 	}
// 	if err := json.Unmarshal([]byte(raw), &persons); err != nil {
// 		return nil, err
// 	}
// 	return persons, nil
// }

// // ─────────────────────────────────────────────
// // HANDLERS
// // ─────────────────────────────────────────────

// // GET /api/customers
// func ListCustomers(c *gin.Context) {
// 	var customers []models.Customer
// 	query := database.DB.Preload("ContactPersons")

// 	if status := c.Query("status"); status != "" {
// 		query = query.Where("customer_status = ?", status)
// 	}
// 	if search := c.Query("search"); search != "" {
// 		like := "%" + search + "%"
// 		query = query.Where(
// 			"company_name ILIKE ? OR account_manager_name ILIKE ? OR account_manager_email ILIKE ?",
// 			like, like, like,
// 		)
// 	}

// 	if err := query.Find(&customers).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customers"})
// 		return
// 	}

// 	out := make([]models.CustomerResponse, 0, len(customers))
// 	for _, cu := range customers {
// 		out = append(out, toCustomerResponse(cu))
// 	}
// 	c.JSON(http.StatusOK, gin.H{"data": out, "total": len(out)})
// }

// // GET /api/customers/:id
// func GetCustomer(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
// 		return
// 	}
// 	var customer models.Customer
// 	if err := database.DB.Preload("ContactPersons").First(&customer, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, toCustomerResponse(customer))
// }

// // POST /api/customers   (multipart/form-data)
// func CreateCustomer(c *gin.Context) {
// 	logoPath, err := saveLogoFile(c)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	persons, err := parseContactPersons(c.PostForm("contact_persons"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact_persons JSON"})
// 		return
// 	}

// 	dbPersons := make([]models.ContactPerson, 0, len(persons))
// 	for _, p := range persons {
// 		dbPersons = append(dbPersons, models.ContactPerson{
// 			Name:         p.Name,
// 			Designation:  p.Designation,
// 			ContactType:  p.ContactType,
// 			ContactLevel: p.ContactLevel,
// 			Emails:       mustJSON(p.Emails),
// 			Mobiles:      mustJSON(p.Mobiles),
// 			Whatsapps:    mustJSON(p.Whatsapps),
// 		})
// 	}

// 	customer := models.Customer{
// 		CompanyName:               c.PostForm("company_name"),
// 		CompanyLogo:               logoPath,
// 		Website:                   c.PostForm("website"),
// 		OfficeAddress:             c.PostForm("office_address"),
// 		CustomerNote:              c.PostForm("customer_note"),
// 		CustomerStatus:            c.PostForm("customer_status"),
// 		AccountManagerName:        c.PostForm("account_manager_name"),
// 		AccountManagerDesignation: c.PostForm("account_manager_designation"),
// 		AccountManagerEmail:       c.PostForm("account_manager_email"),
// 		AccountManagerContact:     c.PostForm("account_manager_contact"),
// 		AccountManagerWhatsapp:    c.PostForm("account_manager_whatsapp"),
// 		AccountManagerBranch:      c.PostForm("account_manager_branch"),
// 		SupportEmails:             c.PostForm("support_emails"),
// 		SupportMobiles:            c.PostForm("support_mobiles"),
// 		SupportWhatsappNumbers:    c.PostForm("support_whatsapp_numbers"),
// 		SupportWhatsappGroups:     c.PostForm("support_whatsapp_groups"),
// 		ContactPersons:            dbPersons,
// 	}

// 	if customer.CustomerStatus == "" {
// 		customer.CustomerStatus = "active"
// 	}

// 	if err := database.DB.Create(&customer).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
// 		return
// 	}
// 	c.JSON(http.StatusCreated, toCustomerResponse(customer))
// }

// // PUT /api/customers/:id   (multipart/form-data)
// func UpdateCustomer(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
// 		return
// 	}

// 	var existing models.Customer
// 	if err := database.DB.Preload("ContactPersons").First(&existing, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
// 		return
// 	}

// 	logoPath := existing.CompanyLogo
// 	newLogo, err := saveLogoFile(c)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if newLogo != "" {
// 		if existing.CompanyLogo != "" {
// 			_ = os.Remove(strings.TrimPrefix(existing.CompanyLogo, "/"))
// 		}
// 		logoPath = newLogo
// 	}

// 	persons, err := parseContactPersons(c.PostForm("contact_persons"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact_persons JSON"})
// 		return
// 	}

// 	tx := database.DB.Begin()
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 		}
// 	}()

// 	existingMap := make(map[uint]models.ContactPerson)
// 	for _, cp := range existing.ContactPersons {
// 		existingMap[cp.ID] = cp
// 	}

// 	var incomingIDs []uint

// 	for _, p := range persons {
// 		if p.ID > 0 {
// 			if _, ok := existingMap[p.ID]; ok {
// 				err := tx.Session(&gorm.Session{}).
// 					Model(&models.ContactPerson{}).
// 					Where("id = ? AND customer_id = ?", p.ID, id).
// 					Updates(map[string]interface{}{
// 						"name":          p.Name,
// 						"designation":   p.Designation,
// 						"contact_type":  p.ContactType,
// 						"contact_level": p.ContactLevel,
// 						"emails":        mustJSON(p.Emails),
// 						"mobiles":       mustJSON(p.Mobiles),
// 						"whatsapps":     mustJSON(p.Whatsapps),
// 					}).Error

// 				if err != nil {
// 					tx.Rollback()
// 					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update contact"})
// 					return
// 				}
// 				incomingIDs = append(incomingIDs, p.ID)
// 				continue
// 			}
// 		}

// 		var found bool
// 		var matchedID uint
// 		for _, cp := range existing.ContactPersons {
// 			if cp.Name == p.Name && cp.Designation == p.Designation {
// 				found = true
// 				matchedID = cp.ID
// 				break
// 			}
// 		}

// 		if found {
// 			err := tx.Session(&gorm.Session{}).
// 				Model(&models.ContactPerson{}).
// 				Where("id = ? AND customer_id = ?", matchedID, id).
// 				Updates(map[string]interface{}{
// 					"name":          p.Name,
// 					"designation":   p.Designation,
// 					"contact_type":  p.ContactType,
// 					"contact_level": p.ContactLevel,
// 					"emails":        mustJSON(p.Emails),
// 					"mobiles":       mustJSON(p.Mobiles),
// 					"whatsapps":     mustJSON(p.Whatsapps),
// 				}).Error

// 			if err != nil {
// 				tx.Rollback()
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update matched contact"})
// 				return
// 			}
// 			incomingIDs = append(incomingIDs, matchedID)
// 			continue
// 		}

// 		contact := models.ContactPerson{
// 			CustomerID:   uint(id),
// 			Name:         p.Name,
// 			Designation:  p.Designation,
// 			ContactType:  p.ContactType,
// 			ContactLevel: p.ContactLevel,
// 			Emails:       mustJSON(p.Emails),
// 			Mobiles:      mustJSON(p.Mobiles),
// 			Whatsapps:    mustJSON(p.Whatsapps),
// 		}
// 		if err := tx.Session(&gorm.Session{}).Create(&contact).Error; err != nil {
// 			tx.Rollback()
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create contact"})
// 			return
// 		}
// 		incomingIDs = append(incomingIDs, contact.ID)
// 	}

// 	cleanTx := tx.Session(&gorm.Session{})

// 	if len(incomingIDs) > 0 {
// 		if err := cleanTx.Exec(
// 			"DELETE FROM contact_people WHERE customer_id = ? AND id NOT IN (?)",
// 			id, incomingIDs,
// 		).Error; err != nil {
// 			tx.Rollback()
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove stale contacts"})
// 			return
// 		}
// 	} else {
// 		if err := cleanTx.Exec(
// 			"DELETE FROM contact_people WHERE customer_id = ?",
// 			id,
// 		).Error; err != nil {
// 			tx.Rollback()
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear all contacts"})
// 			return
// 		}
// 	}

// 	updates := map[string]interface{}{
// 		"company_name":                c.PostForm("company_name"),
// 		"company_logo":                logoPath,
// 		"website":                     c.PostForm("website"),
// 		"office_address":              c.PostForm("office_address"),
// 		"customer_note":               c.PostForm("customer_note"),
// 		"customer_status":             c.PostForm("customer_status"),
// 		"account_manager_name":        c.PostForm("account_manager_name"),
// 		"account_manager_designation": c.PostForm("account_manager_designation"),
// 		"account_manager_email":       c.PostForm("account_manager_email"),
// 		"account_manager_contact":     c.PostForm("account_manager_contact"),
// 		"account_manager_whatsapp":    c.PostForm("account_manager_whatsapp"),
// 		"account_manager_branch":      c.PostForm("account_manager_branch"),
// 		"support_emails":              c.PostForm("support_emails"),
// 		"support_mobiles":             c.PostForm("support_mobiles"),
// 		"support_whatsapp_numbers":    c.PostForm("support_whatsapp_numbers"),
// 		"support_whatsapp_groups":     c.PostForm("support_whatsapp_groups"),
// 	}

// 	if err := tx.Session(&gorm.Session{}).Model(&existing).Updates(updates).Error; err != nil {
// 		tx.Rollback()
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer"})
// 		return
// 	}

// 	if err := tx.Commit().Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed"})
// 		return
// 	}

// 	database.DB.Preload("ContactPersons").First(&existing, id)
// 	c.JSON(http.StatusOK, toCustomerResponse(existing))
// }

// // DELETE /api/customers/:id
// func DeleteCustomer(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
// 		return
// 	}
// 	var customer models.Customer
// 	if err := database.DB.First(&customer, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
// 		return
// 	}
// 	if customer.CompanyLogo != "" {
// 		_ = os.Remove(strings.TrimPrefix(customer.CompanyLogo, "/"))
// 	}
// 	if err := database.DB.Delete(&customer).Error; err != nil {
// 		// FIXED HERE: Added missing c.JSON(http.StatusInternalServerError, ...) wrapping method
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
// }

package handlers

import (
	"encoding/json"
	"fmt"
	"jwt-auth-backend/database"
	"jwt-auth-backend/models"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ─────────────────────────────────────────────
// HELPERS
// ─────────────────────────────────────────────

func mustJSON(v []string) string {
	if v == nil {
		return "[]"
	}
	b, _ := json.Marshal(v)
	return string(b)
}

func parseJSON(s string) []string {
	var result []string
	if s == "" {
		return []string{}
	}
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		return []string{}
	}
	return result
}

func toCustomerResponse(c models.Customer) models.CustomerResponse {
	persons := make([]models.ContactPersonOutput, 0, len(c.ContactPersons))
	for _, p := range c.ContactPersons {
		persons = append(persons, models.ContactPersonOutput{
			ID:           p.ID,
			Name:         p.Name,
			Designation:  p.Designation,
			ContactType:  p.ContactType,
			ContactLevel: p.ContactLevel,
			Emails:       parseJSON(p.Emails),
			Mobiles:      parseJSON(p.Mobiles),
			Whatsapps:    parseJSON(p.Whatsapps),
		})
	}
	return models.CustomerResponse{
		ID:                        c.ID,
		CreatedAt:                 c.CreatedAt,
		UpdatedAt:                 c.UpdatedAt,
		CompanyName:               c.CompanyName,
		CompanyLogo:               c.CompanyLogo,
		Website:                   c.Website,
		OfficeAddress:             c.OfficeAddress,
		CustomerNote:              c.CustomerNote,
		CustomerStatus:            c.CustomerStatus,
		AccountManagerName:        c.AccountManagerName,
		AccountManagerDesignation: c.AccountManagerDesignation,
		AccountManagerEmail:       c.AccountManagerEmail,
		AccountManagerContact:     c.AccountManagerContact,
		AccountManagerWhatsapp:    c.AccountManagerWhatsapp,
		AccountManagerBranch:      c.AccountManagerBranch,
		SupportEmails:             parseJSON(c.SupportEmails),
		SupportMobiles:            parseJSON(c.SupportMobiles),
		SupportWhatsappNumbers:    parseJSON(c.SupportWhatsappNumbers),
		SupportWhatsappGroups:     parseJSON(c.SupportWhatsappGroups),
		ContactPersons:            persons,
	}
}

func saveLogoFile(c *gin.Context) (string, error) {
	file, err := c.FormFile("company_logo")
	if err != nil {
		return "", nil
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowed[ext] {
		return "", fmt.Errorf("unsupported image type: %s", ext)
	}
	_ = os.MkdirAll("uploads", 0755)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join("uploads", filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		return "", err
	}
	return "/uploads/" + filename, nil
}

func parseContactPersons(raw string) ([]models.ContactPersonInput, error) {
	var persons []models.ContactPersonInput
	if raw == "" || raw == "[]" {
		return persons, nil
	}
	if err := json.Unmarshal([]byte(raw), &persons); err != nil {
		return nil, err
	}
	return persons, nil
}

// ─────────────────────────────────────────────
// HANDLERS
// ─────────────────────────────────────────────

// GET /api/customers
func ListCustomers(c *gin.Context) {
	var customers []models.Customer
	query := database.DB.Preload("ContactPersons")

	if status := c.Query("status"); status != "" {
		query = query.Where("customer_status = ?", status)
	}
	if search := c.Query("search"); search != "" {
		like := "%" + search + "%"
		query = query.Where(
			"company_name ILIKE ? OR account_manager_name ILIKE ? OR account_manager_email ILIKE ?",
			like, like, like,
		)
	}

	if err := query.Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customers"})
		return
	}

	out := make([]models.CustomerResponse, 0, len(customers))
	for _, cu := range customers {
		out = append(out, toCustomerResponse(cu))
	}
	c.JSON(http.StatusOK, gin.H{"data": out, "total": len(out)})
}

// Custom endpoint to get only id and company_name for dropdowns
func GetCustomerDropdown(c *gin.Context) {

	var customers []struct {
		ID          uint   `json:"id"`
		CompanyName string `json:"company_name"`
	}

	err := database.DB.
		Model(&models.Customer{}).
		Select("id, company_name").
		Order("company_name ASC").
		Find(&customers).Error

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data": customers,
	})
}

// GET /api/customers/:id
func GetCustomer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var customer models.Customer
	if err := database.DB.Preload("ContactPersons").First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	c.JSON(http.StatusOK, toCustomerResponse(customer))
}

// POST /api/customers
func CreateCustomer(c *gin.Context) {
	logoPath, err := saveLogoFile(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	persons, err := parseContactPersons(c.PostForm("contact_persons"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact_persons JSON"})
		return
	}

	dbPersons := make([]models.ContactPerson, 0, len(persons))
	for _, p := range persons {
		dbPersons = append(dbPersons, models.ContactPerson{
			Name:         p.Name,
			Designation:  p.Designation,
			ContactType:  p.ContactType,
			ContactLevel: p.ContactLevel,
			Emails:       mustJSON(p.Emails),
			Mobiles:      mustJSON(p.Mobiles),
			Whatsapps:    mustJSON(p.Whatsapps),
		})
	}

	customer := models.Customer{
		CompanyName:               c.PostForm("company_name"),
		CompanyLogo:               logoPath,
		Website:                   c.PostForm("website"),
		OfficeAddress:             c.PostForm("office_address"),
		CustomerNote:              c.PostForm("customer_note"),
		CustomerStatus:            c.PostForm("customer_status"),
		AccountManagerName:        c.PostForm("account_manager_name"),
		AccountManagerDesignation: c.PostForm("account_manager_designation"),
		AccountManagerEmail:       c.PostForm("account_manager_email"),
		AccountManagerContact:     c.PostForm("account_manager_contact"),
		AccountManagerWhatsapp:    c.PostForm("account_manager_whatsapp"),
		AccountManagerBranch:      c.PostForm("account_manager_branch"),
		SupportEmails:             c.PostForm("support_emails"),
		SupportMobiles:            c.PostForm("support_mobiles"),
		SupportWhatsappNumbers:    c.PostForm("support_whatsapp_numbers"),
		SupportWhatsappGroups:     c.PostForm("support_whatsapp_groups"),
		ContactPersons:            dbPersons,
	}

	if customer.CustomerStatus == "" {
		customer.CustomerStatus = "active"
	}

	if err := database.DB.Create(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}
	c.JSON(http.StatusCreated, toCustomerResponse(customer))
}

// PUT /api/customers/:id
func UpdateCustomer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	var existing models.Customer
	if err := database.DB.Preload("ContactPersons").First(&existing, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	logoPath := existing.CompanyLogo
	newLogo, err := saveLogoFile(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if newLogo != "" {
		if existing.CompanyLogo != "" {
			_ = os.Remove(strings.TrimPrefix(existing.CompanyLogo, "/"))
		}
		logoPath = newLogo
	}

	persons, err := parseContactPersons(c.PostForm("contact_persons"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact_persons JSON"})
		return
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	existingMap := make(map[uint]models.ContactPerson)
	for _, cp := range existing.ContactPersons {
		existingMap[cp.ID] = cp
	}

	var incomingIDs []uint

	for _, p := range persons {
		// Case 1: Has explicit entry ID -> Execute structural column updates
		if p.ID > 0 {
			if _, ok := existingMap[p.ID]; ok {
				err := tx.Session(&gorm.Session{}).
					Model(&models.ContactPerson{}).
					Where("id = ? AND customer_id = ?", p.ID, id).
					Updates(map[string]interface{}{
						"name":          p.Name,
						"designation":   p.Designation,
						"contact_type":  p.ContactType,
						"contact_level": p.ContactLevel,
						"emails":        mustJSON(p.Emails),
						"mobiles":       mustJSON(p.Mobiles),
						"whatsapps":     mustJSON(p.Whatsapps),
					}).Error

				if err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update contact"})
					return
				}
				incomingIDs = append(incomingIDs, p.ID)
				continue
			}
		}

		// Case 2: Matching validation fallback criteria to capture accidental omissions
		var found bool
		var matchedID uint
		for _, cp := range existing.ContactPersons {
			if cp.Name == p.Name && cp.Designation == p.Designation {
				found = true
				matchedID = cp.ID
				break
			}
		}

		if found {
			err := tx.Session(&gorm.Session{}).
				Model(&models.ContactPerson{}).
				Where("id = ? AND customer_id = ?", matchedID, id).
				Updates(map[string]interface{}{
					"name":          p.Name,
					"designation":   p.Designation,
					"contact_type":  p.ContactType,
					"contact_level": p.ContactLevel,
					"emails":        mustJSON(p.Emails),
					"mobiles":       mustJSON(p.Mobiles),
					"whatsapps":     mustJSON(p.Whatsapps),
				}).Error

			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update matched contact"})
				return
			}
			incomingIDs = append(incomingIDs, matchedID)
			continue
		}

		// Case 3: Completely new row entry processing
		contact := models.ContactPerson{
			CustomerID:   uint(id),
			Name:         p.Name,
			Designation:  p.Designation,
			ContactType:  p.ContactType,
			ContactLevel: p.ContactLevel,
			Emails:       mustJSON(p.Emails),
			Mobiles:      mustJSON(p.Mobiles),
			Whatsapps:    mustJSON(p.Whatsapps),
		}
		if err := tx.Session(&gorm.Session{}).Create(&contact).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create contact"})
			return
		}
		incomingIDs = append(incomingIDs, contact.ID)
	}

	// ─────────────────────────────────────────────────────────────────────────
	// BULLETPROOF RAW SQL DELETIONS (Targeting exactly through mapped instances)
	// ─────────────────────────────────────────────────────────────────────────
	cleanTx := tx.Session(&gorm.Session{})

	if len(incomingIDs) > 0 {
		if err := cleanTx.Exec(
			"DELETE FROM contact_people WHERE customer_id = ? AND id NOT IN (?)",
			id, incomingIDs,
		).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove stale contacts"})
			return
		}
	} else {
		if err := cleanTx.Exec(
			"DELETE FROM contact_people WHERE customer_id = ?",
			id,
		).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear all contacts"})
			return
		}
	}

	// Update parent customer layout
	updates := map[string]interface{}{
		"company_name":                c.PostForm("company_name"),
		"company_logo":                logoPath,
		"website":                     c.PostForm("website"),
		"office_address":              c.PostForm("office_address"),
		"customer_note":               c.PostForm("customer_note"),
		"customer_status":             c.PostForm("customer_status"),
		"account_manager_name":        c.PostForm("account_manager_name"),
		"account_manager_designation": c.PostForm("account_manager_designation"),
		"account_manager_email":       c.PostForm("account_manager_email"),
		"account_manager_contact":     c.PostForm("account_manager_contact"),
		"account_manager_whatsapp":    c.PostForm("account_manager_whatsapp"),
		"account_manager_branch":      c.PostForm("account_manager_branch"),
		"support_emails":              c.PostForm("support_emails"),
		"support_mobiles":             c.PostForm("support_mobiles"),
		"support_whatsapp_numbers":    c.PostForm("support_whatsapp_numbers"),
		"support_whatsapp_groups":     c.PostForm("support_whatsapp_groups"),
	}

	if err := tx.Session(&gorm.Session{}).Model(&existing).Updates(updates).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed"})
		return
	}

	database.DB.Preload("ContactPersons").First(&existing, id)
	c.JSON(http.StatusOK, toCustomerResponse(existing))
}

// DELETE /api/customers/:id
func DeleteCustomer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var customer models.Customer
	if err := database.DB.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	if customer.CompanyLogo != "" {
		_ = os.Remove(strings.TrimPrefix(customer.CompanyLogo, "/"))
	}
	if err := database.DB.Delete(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}

// DELETE /api/customers/:id/contacts
// Query Params: ?contact_id=1328 (Optional. If omitted, clears ALL contacts for this customer)
func DeleteCustomerContacts(c *gin.Context) {
	// 1. Parse and validate the Customer ID from the URL path
	customerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	// 2. Check if the customer actually exists first
	var count int64
	if err := database.DB.Model(&models.Customer{}).Where("id = ?", customerID).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error checking customer"})
		return
	}
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	// 3. Handle deletion logic based on whether 'contact_id' query param is provided
	contactIDStr := c.Query("contact_id")

	if contactIDStr != "" {
		// Scenario A: Delete one specific contact person under this customer
		contactID, err := strconv.Atoi(contactIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact_id parameter"})
			return
		}

		// Execute deletion ensuring BOTH contact_id and customer_id match
		result := database.DB.Exec(
			"DELETE FROM contact_people WHERE id = ? AND customer_id = ?",
			contactID, customerID,
		)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete contact person"})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contact person not found for this customer"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Contact person deleted successfully",
			"details": gin.H{"contact_id": contactID, "customer_id": customerID},
		})
		return

	} else {
		// Scenario B: Clear ALL contact persons under this customer
		result := database.DB.Exec("DELETE FROM contact_people WHERE customer_id = ?", customerID)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear customer contacts"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":       "All contact persons cleared successfully for this customer",
			"rows_affected": result.RowsAffected,
		})
		return
	}
}
