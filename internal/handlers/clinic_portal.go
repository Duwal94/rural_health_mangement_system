package handlers

import (
	"rural_health_management_system/internal/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ClinicPortalHandler struct {
	db *gorm.DB
}

func NewClinicPortalHandler(db *gorm.DB) *ClinicPortalHandler {
	return &ClinicPortalHandler{db: db}
}

// GetMyProfile returns the clinic's own profile
func (h *ClinicPortalHandler) GetMyProfile(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)

	var clinic models.Clinic
	if err := h.db.First(&clinic, clinicID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Clinic not found",
		})
	}

	return c.JSON(clinic)
}

// UpdateMyProfile allows clinics to update their own profile
func (h *ClinicPortalHandler) UpdateMyProfile(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)

	var req map[string]interface{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var clinic models.Clinic
	if err := h.db.First(&clinic, clinicID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Clinic not found",
		})
	}

	// Update allowed fields
	allowedFields := []string{"name", "address", "contact_number", "district"}
	updates := make(map[string]interface{})

	for _, field := range allowedFields {
		if value, exists := req[field]; exists {
			updates[field] = value
		}
	}

	if len(updates) > 0 {
		if err := h.db.Model(&clinic).Updates(updates).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update clinic",
			})
		}
	}

	// Return updated clinic
	h.db.First(&clinic, clinicID)
	return c.JSON(clinic)
}

// GetMyPatients returns all patients for the clinic
func (h *ClinicPortalHandler) GetMyPatients(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	if perPage > 100 {
		perPage = 100
	}
	search := c.Query("search")

	offset := (page - 1) * perPage

	var patients []models.Patient
	var total int64

	query := h.db.Model(&models.Patient{}).Where("clinic_id = ?", clinicID)

	if search != "" {
		query = query.Where("full_name ILIKE ? OR phone ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)

	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&patients).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch patients",
		})
	}

	totalPages := int(total) / perPage
	if int(total)%perPage != 0 {
		totalPages++
	}

	return c.JSON(models.PaginationResponse{
		Data:       patients,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GetMyPatient returns a specific patient for the clinic
func (h *ClinicPortalHandler) GetMyPatient(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)
	patientID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid patient ID",
		})
	}

	var patient models.Patient
	if err := h.db.Preload("Visits").Preload("Visits.Diagnoses").Preload("Visits.Prescriptions").
		Where("id = ? AND clinic_id = ?", patientID, clinicID).First(&patient).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Patient not found",
		})
	}

	return c.JSON(patient)
}

// GetMyStaff returns all staff for the clinic
func (h *ClinicPortalHandler) GetMyStaff(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	if perPage > 100 {
		perPage = 100
	}
	role := c.Query("role")

	offset := (page - 1) * perPage

	var staff []models.Staff
	var total int64

	query := h.db.Model(&models.Staff{}).Where("clinic_id = ?", clinicID)

	if role != "" {
		query = query.Where("role = ?", role)
	}

	query.Count(&total)

	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&staff).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch staff",
		})
	}

	totalPages := int(total) / perPage
	if int(total)%perPage != 0 {
		totalPages++
	}

	return c.JSON(models.PaginationResponse{
		Data:       staff,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

// CreateStaff allows clinics to add new staff
func (h *ClinicPortalHandler) CreateStaff(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)

	var req map[string]interface{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Set clinic ID
	req["clinic_id"] = clinicID

	staff := models.Staff{
		FullName: req["full_name"].(string),
		Role:     req["role"].(string),
		Phone:    req["phone"].(string),
		Email:    req["email"].(string),
		ClinicID: clinicID,
	}

	if err := h.db.Create(&staff).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create staff",
		})
	}

	// Load with clinic relationship
	h.db.Preload("Clinic").First(&staff, staff.ID)
	return c.Status(fiber.StatusCreated).JSON(staff)
}

// GetMyVisits returns all visits for the clinic
func (h *ClinicPortalHandler) GetMyVisits(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	if perPage > 100 {
		perPage = 100
	}
	patientIDStr := c.Query("patient_id")

	offset := (page - 1) * perPage

	var visits []models.Visit
	var total int64

	query := h.db.Model(&models.Visit{}).Where("clinic_id = ?", clinicID)

	if patientIDStr != "" {
		if patientID, err := strconv.ParseUint(patientIDStr, 10, 32); err == nil {
			query = query.Where("patient_id = ?", patientID)
		}
	}

	query.Count(&total)

	if err := query.Preload("Patient").Preload("Staff").Preload("Diagnoses").Preload("Prescriptions").
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

// CreateVisit allows clinics to create new visits
func (h *ClinicPortalHandler) CreateVisit(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)

	var req models.CreateVisitRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Verify patient belongs to this clinic
	var patient models.Patient
	if err := h.db.Where("id = ? AND clinic_id = ?", req.PatientID, clinicID).First(&patient).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Patient not found or doesn't belong to this clinic",
		})
	}

	// Verify staff belongs to this clinic
	var staff models.Staff
	if err := h.db.Where("id = ? AND clinic_id = ?", req.StaffID, clinicID).First(&staff).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Staff not found or doesn't belong to this clinic",
		})
	}

	// Set clinic ID
	req.ClinicID = clinicID

	visit := models.Visit{
		PatientID: req.PatientID,
		ClinicID:  clinicID,
		StaffID:   req.StaffID,
		VisitDate: req.VisitDate,
		Reason:    req.Reason,
		Notes:     req.Notes,
	}

	if err := h.db.Create(&visit).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create visit",
		})
	}

	// Load with relationships
	h.db.Preload("Patient").Preload("Clinic").Preload("Staff").First(&visit, visit.ID)
	return c.Status(fiber.StatusCreated).JSON(visit)
}

// GetMyVisit returns a specific visit for the clinic
func (h *ClinicPortalHandler) GetMyVisit(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)
	visitID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid visit ID",
		})
	}

	var visit models.Visit
	if err := h.db.Preload("Patient").Preload("Staff").Preload("Diagnoses").Preload("Prescriptions").
		Where("id = ? AND clinic_id = ?", visitID, clinicID).First(&visit).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Visit not found",
		})
	}

	return c.JSON(visit)
}

// CreateDiagnosis allows clinics to add diagnoses to visits
func (h *ClinicPortalHandler) CreateDiagnosis(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)

	var req models.CreateDiagnosisRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Verify visit belongs to this clinic
	var visit models.Visit
	if err := h.db.Where("id = ? AND clinic_id = ?", req.VisitID, clinicID).First(&visit).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Visit not found or doesn't belong to this clinic",
		})
	}

	diagnosis := models.Diagnosis{
		VisitID:       req.VisitID,
		DiagnosisCode: req.DiagnosisCode,
		Description:   req.Description,
	}

	if err := h.db.Create(&diagnosis).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create diagnosis",
		})
	}

	// Load with visit relationship
	h.db.Preload("Visit").First(&diagnosis, diagnosis.ID)
	return c.Status(fiber.StatusCreated).JSON(diagnosis)
}

// CreatePrescription allows clinics to add prescriptions to visits
func (h *ClinicPortalHandler) CreatePrescription(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)

	var req models.CreatePrescriptionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Verify visit belongs to this clinic
	var visit models.Visit
	if err := h.db.Where("id = ? AND clinic_id = ?", req.VisitID, clinicID).First(&visit).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Visit not found or doesn't belong to this clinic",
		})
	}

	prescription := models.Prescription{
		VisitID:        req.VisitID,
		MedicationName: req.MedicationName,
		Dosage:         req.Dosage,
		Instructions:   req.Instructions,
		DurationDays:   req.DurationDays,
	}

	if err := h.db.Create(&prescription).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create prescription",
		})
	}

	// Load with visit relationship
	h.db.Preload("Visit").First(&prescription, prescription.ID)
	return c.Status(fiber.StatusCreated).JSON(prescription)
}

// GetDashboardStats returns dashboard statistics for the clinic
func (h *ClinicPortalHandler) GetDashboardStats(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)

	var stats struct {
		TotalPatients      int64 `json:"total_patients"`
		TotalStaff         int64 `json:"total_staff"`
		TotalVisits        int64 `json:"total_visits"`
		VisitsThisMonth    int64 `json:"visits_this_month"`
		TotalDiagnoses     int64 `json:"total_diagnoses"`
		TotalPrescriptions int64 `json:"total_prescriptions"`
	}

	// Count patients
	h.db.Model(&models.Patient{}).Where("clinic_id = ?", clinicID).Count(&stats.TotalPatients)

	// Count staff
	h.db.Model(&models.Staff{}).Where("clinic_id = ?", clinicID).Count(&stats.TotalStaff)

	// Count total visits
	h.db.Model(&models.Visit{}).Where("clinic_id = ?", clinicID).Count(&stats.TotalVisits)

	// Count visits this month
	h.db.Model(&models.Visit{}).Where("clinic_id = ? AND DATE_TRUNC('month', visit_date) = DATE_TRUNC('month', CURRENT_DATE)", clinicID).Count(&stats.VisitsThisMonth)

	// Count diagnoses
	h.db.Model(&models.Diagnosis{}).
		Joins("JOIN visits ON diagnoses.visit_id = visits.id").
		Where("visits.clinic_id = ?", clinicID).
		Count(&stats.TotalDiagnoses)

	// Count prescriptions
	h.db.Model(&models.Prescription{}).
		Joins("JOIN visits ON prescriptions.visit_id = visits.id").
		Where("visits.clinic_id = ?", clinicID).
		Count(&stats.TotalPrescriptions)

	return c.JSON(stats)
}
