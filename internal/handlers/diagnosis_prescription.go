package handlers

import (
	"math"
	"strconv"

	"rural_health_management_system/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// DiagnosisHandler handles diagnosis-related operations
type DiagnosisHandler struct {
	db *gorm.DB
}

func NewDiagnosisHandler(db *gorm.DB) *DiagnosisHandler {
	return &DiagnosisHandler{db: db}
}

func (h *DiagnosisHandler) GetDiagnoses(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	visitID := c.Query("visit_id", "")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage
	query := h.db.Model(&models.Diagnosis{}).Preload("Visit")

	if visitID != "" {
		query = query.Where("visit_id = ?", visitID)
	}

	var total int64
	query.Count(&total)

	var diagnoses []models.Diagnosis
	if err := query.Offset(offset).Limit(perPage).Find(&diagnoses).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to fetch diagnoses",
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return c.JSON(models.PaginationResponse{
		Data:       diagnoses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

func (h *DiagnosisHandler) GetDiagnosisByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var diagnosis models.Diagnosis
	if err := h.db.Preload("Visit").First(&diagnosis, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "Diagnosis not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to fetch diagnosis",
		})
	}

	return c.JSON(diagnosis)
}

func (h *DiagnosisHandler) CreateDiagnosis(c *fiber.Ctx) error {
	var diagnosis models.Diagnosis

	if err := c.BodyParser(&diagnosis); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Invalid input",
		})
	}

	// Validate required fields
	if diagnosis.VisitID == 0 || diagnosis.DiagnosisCode == "" || diagnosis.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Missing required fields",
		})
	}

	// Check if visit exists
	var visit models.Visit
	if err := h.db.First(&visit, diagnosis.VisitID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Visit not found",
		})
	}

	if err := h.db.Create(&diagnosis).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to create diagnosis",
		})
	}

	h.db.Preload("Visit").First(&diagnosis, diagnosis.ID)
	return c.Status(fiber.StatusCreated).JSON(diagnosis)
}

func (h *DiagnosisHandler) UpdateDiagnosis(c *fiber.Ctx) error {
	id := c.Params("id")

	var diagnosis models.Diagnosis
	if err := h.db.First(&diagnosis, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "Diagnosis not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to fetch diagnosis",
		})
	}

	var updates models.Diagnosis
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Invalid input",
		})
	}

	// Update fields
	if updates.DiagnosisCode != "" {
		diagnosis.DiagnosisCode = updates.DiagnosisCode
	}
	if updates.Description != "" {
		diagnosis.Description = updates.Description
	}

	if err := h.db.Save(&diagnosis).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to update diagnosis",
		})
	}

	h.db.Preload("Visit").First(&diagnosis, diagnosis.ID)
	return c.JSON(diagnosis)
}

func (h *DiagnosisHandler) DeleteDiagnosis(c *fiber.Ctx) error {
	id := c.Params("id")

	var diagnosis models.Diagnosis
	if err := h.db.First(&diagnosis, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "Diagnosis not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to fetch diagnosis",
		})
	}

	if err := h.db.Delete(&diagnosis).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to delete diagnosis",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// PrescriptionHandler handles prescription-related operations
type PrescriptionHandler struct {
	db *gorm.DB
}

func NewPrescriptionHandler(db *gorm.DB) *PrescriptionHandler {
	return &PrescriptionHandler{db: db}
}

func (h *PrescriptionHandler) GetPrescriptions(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	visitID := c.Query("visit_id", "")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage
	query := h.db.Model(&models.Prescription{}).Preload("Visit")

	if visitID != "" {
		query = query.Where("visit_id = ?", visitID)
	}

	var total int64
	query.Count(&total)

	var prescriptions []models.Prescription
	if err := query.Offset(offset).Limit(perPage).Find(&prescriptions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to fetch prescriptions",
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return c.JSON(models.PaginationResponse{
		Data:       prescriptions,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

func (h *PrescriptionHandler) GetPrescriptionByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var prescription models.Prescription
	if err := h.db.Preload("Visit").First(&prescription, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "Prescription not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to fetch prescription",
		})
	}

	return c.JSON(prescription)
}

func (h *PrescriptionHandler) CreatePrescription(c *fiber.Ctx) error {
	var prescription models.Prescription

	if err := c.BodyParser(&prescription); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Invalid input",
		})
	}

	// Validate required fields
	if prescription.VisitID == 0 || prescription.MedicationName == "" ||
		prescription.Dosage == "" || prescription.Instructions == "" || prescription.DurationDays <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Missing required fields",
		})
	}

	// Check if visit exists
	var visit models.Visit
	if err := h.db.First(&visit, prescription.VisitID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Visit not found",
		})
	}

	if err := h.db.Create(&prescription).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to create prescription",
		})
	}

	h.db.Preload("Visit").First(&prescription, prescription.ID)
	return c.Status(fiber.StatusCreated).JSON(prescription)
}

func (h *PrescriptionHandler) UpdatePrescription(c *fiber.Ctx) error {
	id := c.Params("id")

	var prescription models.Prescription
	if err := h.db.First(&prescription, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "Prescription not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to fetch prescription",
		})
	}

	var updates models.Prescription
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Invalid input",
		})
	}

	// Update fields
	if updates.MedicationName != "" {
		prescription.MedicationName = updates.MedicationName
	}
	if updates.Dosage != "" {
		prescription.Dosage = updates.Dosage
	}
	if updates.Instructions != "" {
		prescription.Instructions = updates.Instructions
	}
	if updates.DurationDays > 0 {
		prescription.DurationDays = updates.DurationDays
	}

	if err := h.db.Save(&prescription).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to update prescription",
		})
	}

	h.db.Preload("Visit").First(&prescription, prescription.ID)
	return c.JSON(prescription)
}

func (h *PrescriptionHandler) DeletePrescription(c *fiber.Ctx) error {
	id := c.Params("id")

	var prescription models.Prescription
	if err := h.db.First(&prescription, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "Prescription not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to fetch prescription",
		})
	}

	if err := h.db.Delete(&prescription).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to delete prescription",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
