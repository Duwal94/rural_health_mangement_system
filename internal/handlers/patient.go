package handlers

import (
	"math"
	"strconv"
	"time"

	"rural_health_management_system/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PatientHandler struct {
	db *gorm.DB
}

func NewPatientHandler(db *gorm.DB) *PatientHandler {
	return &PatientHandler{db: db}
}

// GetPatients - GET /patients with pagination and search
func (h *PatientHandler) GetPatients(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	search := c.Query("search", "")
	clinicID := c.Query("clinic_id", "")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	query := h.db.Model(&models.Patient{}).Preload("Clinic")

	// Apply filters
	if search != "" {
		query = query.Where("full_name ILIKE ? OR phone ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if clinicID != "" {
		query = query.Where("clinic_id = ?", clinicID)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to count patients",
		})
	}

	// Get patients with pagination
	var patients []models.Patient
	if err := query.Offset(offset).Limit(perPage).Find(&patients).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to fetch patients",
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return c.JSON(models.PaginationResponse{
		Data:       patients,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GetPatient - GET /patients/:id
func (h *PatientHandler) GetPatient(c *fiber.Ctx) error {
	id := c.Params("id")

	var patient models.Patient
	if err := h.db.Preload("Clinic").Preload("Visits").First(&patient, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "Patient not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to fetch patient",
		})
	}

	return c.JSON(patient)
}

// CreatePatient - POST /patients
func (h *PatientHandler) CreatePatient(c *fiber.Ctx) error {
	var req models.CreatePatientRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid input",
			Details: err.Error(),
		})
	}

	// Validate required fields
	if req.FullName == "" || req.Gender == "" || req.DateOfBirth == "" ||
		req.Address == "" || req.Phone == "" || req.ClinicID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Missing required fields",
		})
	}

	// Validate gender
	if req.Gender != "Male" && req.Gender != "Female" && req.Gender != "Other" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Invalid gender. Must be Male, Female, or Other",
		})
	}

	// Parse date of birth
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Invalid date format. Use YYYY-MM-DD",
		})
	}

	// Check if clinic exists
	var clinic models.Clinic
	if err := h.db.First(&clinic, req.ClinicID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error: "Clinic not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to validate clinic",
		})
	}

	patient := models.Patient{
		FullName:    req.FullName,
		Gender:      req.Gender,
		DateOfBirth: dob,
		Address:     req.Address,
		Phone:       req.Phone,
		ClinicID:    req.ClinicID,
	}

	if err := h.db.Create(&patient).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to create patient",
			Details: err.Error(),
		})
	}

	// Preload clinic for response
	h.db.Preload("Clinic").First(&patient, patient.ID)

	return c.Status(fiber.StatusCreated).JSON(patient)
}

// UpdatePatient - PUT /patients/:id
func (h *PatientHandler) UpdatePatient(c *fiber.Ctx) error {
	id := c.Params("id")

	var patient models.Patient
	if err := h.db.First(&patient, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "Patient not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to fetch patient",
		})
	}

	var req models.UpdatePatientRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid input",
			Details: err.Error(),
		})
	}

	// Update fields if provided
	if req.FullName != nil {
		patient.FullName = *req.FullName
	}
	if req.Gender != nil {
		if *req.Gender != "Male" && *req.Gender != "Female" && *req.Gender != "Other" {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error: "Invalid gender. Must be Male, Female, or Other",
			})
		}
		patient.Gender = *req.Gender
	}
	if req.DateOfBirth != nil {
		dob, err := time.Parse("2006-01-02", *req.DateOfBirth)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error: "Invalid date format. Use YYYY-MM-DD",
			})
		}
		patient.DateOfBirth = dob
	}
	if req.Address != nil {
		patient.Address = *req.Address
	}
	if req.Phone != nil {
		patient.Phone = *req.Phone
	}
	if req.ClinicID != nil {
		// Check if clinic exists
		var clinic models.Clinic
		if err := h.db.First(&clinic, *req.ClinicID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
					Error: "Clinic not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
				Error: "Failed to validate clinic",
			})
		}
		patient.ClinicID = *req.ClinicID
	}

	patient.UpdatedAt = time.Now()

	if err := h.db.Save(&patient).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to update patient",
			Details: err.Error(),
		})
	}

	// Preload clinic for response
	h.db.Preload("Clinic").First(&patient, patient.ID)

	return c.JSON(patient)
}

// DeletePatient - DELETE /patients/:id
func (h *PatientHandler) DeletePatient(c *fiber.Ctx) error {
	id := c.Params("id")

	// Check if patient exists
	var patient models.Patient
	if err := h.db.First(&patient, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "Patient not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to fetch patient",
		})
	}

	// Check if patient has visits
	var visitCount int64
	h.db.Model(&models.Visit{}).Where("patient_id = ?", id).Count(&visitCount)

	if visitCount > 0 {
		return c.Status(fiber.StatusConflict).JSON(models.ErrorResponse{
			Error: "Cannot delete patient with existing visits",
		})
	}

	// Soft delete
	if err := h.db.Delete(&patient).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to delete patient",
			Details: err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
