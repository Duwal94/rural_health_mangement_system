package handlers

import (
	"rural_health_management_system/internal/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PatientPortalHandler struct {
	db *gorm.DB
}

func NewPatientPortalHandler(db *gorm.DB) *PatientPortalHandler {
	return &PatientPortalHandler{db: db}
}

// GetMyProfile returns the patient's own profile
func (h *PatientPortalHandler) GetMyProfile(c *fiber.Ctx) error {
	patientID := c.Locals("patient_id").(uint)

	var patient models.Patient
	if err := h.db.Preload("Clinic").First(&patient, patientID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Patient not found",
		})
	}

	return c.JSON(patient)
}

// UpdateMyProfile allows patients to update their own profile
func (h *PatientPortalHandler) UpdateMyProfile(c *fiber.Ctx) error {
	patientID := c.Locals("patient_id").(uint)

	var req models.UpdatePatientRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var patient models.Patient
	if err := h.db.First(&patient, patientID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Patient not found",
		})
	}

	// Update allowed fields (patients can't change their clinic)
	updates := make(map[string]interface{})
	if req.FullName != nil {
		updates["full_name"] = *req.FullName
	}
	if req.Gender != nil {
		updates["gender"] = *req.Gender
	}
	if req.Address != nil {
		updates["address"] = *req.Address
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.DateOfBirth != nil {
		// Parse date
		// This would need proper date parsing logic
		updates["date_of_birth"] = *req.DateOfBirth
	}

	if err := h.db.Model(&patient).Updates(updates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update patient",
		})
	}

	// Return updated patient with clinic
	h.db.Preload("Clinic").First(&patient, patientID)
	return c.JSON(patient)
}

// GetMyVisits returns all visits for the authenticated patient
func (h *PatientPortalHandler) GetMyVisits(c *fiber.Ctx) error {
	patientID := c.Locals("patient_id").(uint)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	if perPage > 100 {
		perPage = 100
	}

	offset := (page - 1) * perPage

	var visits []models.Visit
	var total int64

	query := h.db.Model(&models.Visit{}).Where("patient_id = ?", patientID)

	query.Count(&total)

	if err := query.Preload("Clinic").Preload("Staff").Preload("Diagnoses").Preload("Prescriptions").
		Offset(offset).Limit(perPage).Order("visit_date DESC").Find(&visits).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch visits",
		})
	}

	totalPages := int(total) / perPage
	if int(total)%perPage != 0 {
		totalPages++
	}

	return c.JSON(models.PaginationResponse{
		Data:       visits,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GetMyVisit returns a specific visit for the authenticated patient
func (h *PatientPortalHandler) GetMyVisit(c *fiber.Ctx) error {
	patientID := c.Locals("patient_id").(uint)
	visitID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid visit ID",
		})
	}

	var visit models.Visit
	if err := h.db.Preload("Clinic").Preload("Staff").Preload("Diagnoses").Preload("Prescriptions").
		Where("id = ? AND patient_id = ?", visitID, patientID).First(&visit).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Visit not found",
		})
	}

	return c.JSON(visit)
}

// GetMyDiagnoses returns all diagnoses for the authenticated patient
func (h *PatientPortalHandler) GetMyDiagnoses(c *fiber.Ctx) error {
	patientID := c.Locals("patient_id").(uint)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	if perPage > 100 {
		perPage = 100
	}

	offset := (page - 1) * perPage

	var diagnoses []models.Diagnosis
	var total int64

	// Join with visits to filter by patient
	query := h.db.Model(&models.Diagnosis{}).
		Joins("JOIN visits ON diagnoses.visit_id = visits.id").
		Where("visits.patient_id = ?", patientID)

	query.Count(&total)

	if err := query.Preload("Visit").Preload("Visit.Clinic").Preload("Visit.Staff").
		Offset(offset).Limit(perPage).Order("created_at DESC").Find(&diagnoses).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch diagnoses",
		})
	}

	totalPages := int(total) / perPage
	if int(total)%perPage != 0 {
		totalPages++
	}

	return c.JSON(models.PaginationResponse{
		Data:       diagnoses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GetMyPrescriptions returns all prescriptions for the authenticated patient
func (h *PatientPortalHandler) GetMyPrescriptions(c *fiber.Ctx) error {
	patientID := c.Locals("patient_id").(uint)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	if perPage > 100 {
		perPage = 100
	}

	offset := (page - 1) * perPage

	var prescriptions []models.Prescription
	var total int64

	// Join with visits to filter by patient
	query := h.db.Model(&models.Prescription{}).
		Joins("JOIN visits ON prescriptions.visit_id = visits.id").
		Where("visits.patient_id = ?", patientID)

	query.Count(&total)

	if err := query.Preload("Visit").Preload("Visit.Clinic").Preload("Visit.Staff").
		Offset(offset).Limit(perPage).Order("created_at DESC").Find(&prescriptions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch prescriptions",
		})
	}

	totalPages := int(total) / perPage
	if int(total)%perPage != 0 {
		totalPages++
	}

	return c.JSON(models.PaginationResponse{
		Data:       prescriptions,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}
