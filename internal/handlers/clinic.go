package handlers

import (
	"math"
	"strconv"
	"time"

	"rural_health_management_system/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ClinicHandler struct {
	db *gorm.DB
}

func NewClinicHandler(db *gorm.DB) *ClinicHandler {
	return &ClinicHandler{db: db}
}

// GetClinics - GET /clinics with pagination and search
func (h *ClinicHandler) GetClinics(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	search := c.Query("search", "")
	district := c.Query("district", "")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	query := h.db.Model(&models.Clinic{})

	// Apply filters
	if search != "" {
		query = query.Where("name ILIKE ? OR address ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if district != "" {
		query = query.Where("district ILIKE ?", "%"+district+"%")
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to count clinics",
		})
	}

	// Get clinics with pagination
	var clinics []models.Clinic
	if err := query.Offset(offset).Limit(perPage).Find(&clinics).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to fetch clinics",
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return c.JSON(models.PaginationResponse{
		Data:       clinics,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GetClinic - GET /clinics/:id
func (h *ClinicHandler) GetClinic(c *fiber.Ctx) error {
	id := c.Params("id")

	var clinic models.Clinic
	if err := h.db.Preload("Patients").Preload("Staff").First(&clinic, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "Clinic not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to fetch clinic",
		})
	}

	return c.JSON(clinic)
}

// CreateClinic - POST /clinics
func (h *ClinicHandler) CreateClinic(c *fiber.Ctx) error {
	var clinic models.Clinic

	if err := c.BodyParser(&clinic); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid input",
			Details: err.Error(),
		})
	}

	// Validate required fields
	if clinic.Name == "" || clinic.Address == "" || clinic.ContactNumber == "" || clinic.District == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Missing required fields: name, address, contact_number, district",
		})
	}

	if err := h.db.Create(&clinic).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to create clinic",
			Details: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(clinic)
}

// UpdateClinic - PUT /clinics/:id
func (h *ClinicHandler) UpdateClinic(c *fiber.Ctx) error {
	id := c.Params("id")

	var clinic models.Clinic
	if err := h.db.First(&clinic, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "Clinic not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to fetch clinic",
		})
	}

	var updates models.Clinic
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid input",
			Details: err.Error(),
		})
	}

	// Update fields
	if updates.Name != "" {
		clinic.Name = updates.Name
	}
	if updates.Address != "" {
		clinic.Address = updates.Address
	}
	if updates.ContactNumber != "" {
		clinic.ContactNumber = updates.ContactNumber
	}
	if updates.District != "" {
		clinic.District = updates.District
	}

	clinic.UpdatedAt = time.Now()

	if err := h.db.Save(&clinic).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to update clinic",
			Details: err.Error(),
		})
	}

	return c.JSON(clinic)
}

// DeleteClinic - DELETE /clinics/:id
func (h *ClinicHandler) DeleteClinic(c *fiber.Ctx) error {
	id := c.Params("id")

	// Check if clinic exists
	var clinic models.Clinic
	if err := h.db.First(&clinic, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "Clinic not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to fetch clinic",
		})
	}

	// Check if clinic has associated patients or staff
	var patientCount, staffCount int64
	h.db.Model(&models.Patient{}).Where("clinic_id = ?", id).Count(&patientCount)
	h.db.Model(&models.Staff{}).Where("clinic_id = ?", id).Count(&staffCount)

	if patientCount > 0 || staffCount > 0 {
		return c.Status(fiber.StatusConflict).JSON(models.ErrorResponse{
			Error: "Cannot delete clinic with existing patients or staff",
		})
	}

	// Soft delete
	if err := h.db.Delete(&clinic).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to delete clinic",
			Details: err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
