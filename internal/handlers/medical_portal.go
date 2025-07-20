package handlers

import (
	"rural_health_management_system/internal/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type MedicalPortalHandler struct {
	db *gorm.DB
}

func NewMedicalPortalHandler(db *gorm.DB) *MedicalPortalHandler {
	return &MedicalPortalHandler{db: db}
}

// GetMyProfile returns the staff member's own profile
func (h *MedicalPortalHandler) GetMyProfile(c *fiber.Ctx) error {
	staffID := c.Locals("staff_id").(uint)

	var staff models.Staff
	if err := h.db.Preload("Clinic").Preload("User").First(&staff, staffID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Staff profile not found",
		})
	}

	return c.JSON(staff)
}

// UpdateMyProfile allows medical staff to update their own profile (limited fields)
func (h *MedicalPortalHandler) UpdateMyProfile(c *fiber.Ctx) error {
	staffID := c.Locals("staff_id").(uint)

	var req map[string]interface{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var staff models.Staff
	if err := h.db.First(&staff, staffID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Staff profile not found",
		})
	}

	// Medical staff can only update limited fields
	allowedFields := []string{"phone"}
	updates := make(map[string]interface{})

	for _, field := range allowedFields {
		if value, exists := req[field]; exists {
			updates[field] = value
		}
	}

	if len(updates) > 0 {
		if err := h.db.Model(&staff).Updates(updates).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update profile",
			})
		}
	}

	// Return updated staff
	h.db.Preload("Clinic").Preload("User").First(&staff, staffID)
	return c.JSON(staff)
}

// GetMyPatients returns patients that the medical staff can access
func (h *MedicalPortalHandler) GetMyPatients(c *fiber.Ctx) error {
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

	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	return c.JSON(models.PaginationResponse{
		Data:       patients,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GetMyPatient returns a specific patient
func (h *MedicalPortalHandler) GetMyPatient(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)
	patientID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid patient ID",
		})
	}

	var patient models.Patient
	if err := h.db.Preload("Clinic").Preload("Visits.Staff").Preload("Visits.Diagnoses").Preload("Visits.Prescriptions").Where("id = ? AND clinic_id = ?", patientID, clinicID).First(&patient).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Patient not found",
		})
	}

	return c.JSON(patient)
}

// CreateVisit creates a new visit (medical staff)
func (h *MedicalPortalHandler) CreateVisit(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)
	staffID := c.Locals("staff_id").(uint)

	var req models.CreateVisitRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Force the clinic ID and staff ID
	req.ClinicID = clinicID
	req.StaffID = staffID

	// Verify patient belongs to this clinic
	var patient models.Patient
	if err := h.db.Where("id = ? AND clinic_id = ?", req.PatientID, clinicID).First(&patient).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Patient not found in this clinic",
		})
	}

	// Set visit date to now if not provided
	if req.VisitDate.IsZero() {
		req.VisitDate = time.Now()
	}

	visit := models.Visit{
		PatientID: req.PatientID,
		ClinicID:  req.ClinicID,
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

	// Load relationships
	h.db.Preload("Patient").Preload("Clinic").Preload("Staff").First(&visit, visit.ID)

	return c.Status(fiber.StatusCreated).JSON(visit)
}

// GetMyVisits returns visits where the medical staff was involved
func (h *MedicalPortalHandler) GetMyVisits(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)
	staffID := c.Locals("staff_id").(uint)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	if perPage > 100 {
		perPage = 100
	}
	patientID := c.Query("patient_id")
	showAll := c.Query("show_all") // If "true", show all clinic visits

	offset := (page - 1) * perPage

	var visits []models.Visit
	var total int64

	query := h.db.Model(&models.Visit{}).Where("clinic_id = ?", clinicID)

	// If not showing all, filter by staff ID
	if showAll != "true" {
		query = query.Where("staff_id = ?", staffID)
	}

	if patientID != "" {
		query = query.Where("patient_id = ?", patientID)
	}

	query.Count(&total)

	if err := query.Preload("Patient").Preload("Staff").Preload("Diagnoses").Preload("Prescriptions").Offset(offset).Limit(perPage).Order("visit_date DESC").Find(&visits).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch visits",
		})
	}

	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	return c.JSON(models.PaginationResponse{
		Data:       visits,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GetMyVisit returns a specific visit
func (h *MedicalPortalHandler) GetMyVisit(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)
	visitID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid visit ID",
		})
	}

	var visit models.Visit
	if err := h.db.Preload("Patient").Preload("Clinic").Preload("Staff").Preload("Diagnoses").Preload("Prescriptions").Where("id = ? AND clinic_id = ?", visitID, clinicID).First(&visit).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Visit not found",
		})
	}

	return c.JSON(visit)
}

// CreateDiagnosis creates a new diagnosis (doctors only)
func (h *MedicalPortalHandler) CreateDiagnosis(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)

	var req models.CreateDiagnosisRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Verify visit exists and belongs to this clinic
	var visit models.Visit
	if err := h.db.Where("id = ? AND clinic_id = ?", req.VisitID, clinicID).First(&visit).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Visit not found in this clinic",
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

	// Load relationships
	h.db.Preload("Visit").First(&diagnosis, diagnosis.ID)

	return c.Status(fiber.StatusCreated).JSON(diagnosis)
}

// CreatePrescription creates a new prescription (doctors only)
func (h *MedicalPortalHandler) CreatePrescription(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)

	var req models.CreatePrescriptionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Verify visit exists and belongs to this clinic
	var visit models.Visit
	if err := h.db.Where("id = ? AND clinic_id = ?", req.VisitID, clinicID).First(&visit).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Visit not found in this clinic",
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

	// Load relationships
	h.db.Preload("Visit").First(&prescription, prescription.ID)

	return c.Status(fiber.StatusCreated).JSON(prescription)
}

// GetStaff returns staff information (read-only for medical staff)
func (h *MedicalPortalHandler) GetStaff(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	if perPage > 100 {
		perPage = 100
	}
	search := c.Query("search")

	offset := (page - 1) * perPage

	var staff []models.Staff
	var total int64

	query := h.db.Model(&models.Staff{}).Where("clinic_id = ?", clinicID)

	if search != "" {
		query = query.Where("full_name ILIKE ? OR role ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)

	if err := query.Preload("User").Offset(offset).Limit(perPage).Order("created_at DESC").Find(&staff).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch staff",
		})
	}

	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	return c.JSON(models.PaginationResponse{
		Data:       staff,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GetDashboardStats returns dashboard statistics for medical staff
func (h *MedicalPortalHandler) GetDashboardStats(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)
	staffID := c.Locals("staff_id").(uint)

	var stats struct {
		MyVisitsToday   int64 `json:"my_visits_today"`
		MyTotalVisits   int64 `json:"my_total_visits"`
		TotalPatients   int64 `json:"total_patients"`
		TotalStaff      int64 `json:"total_staff"`
		MyDiagnoses     int64 `json:"my_diagnoses"`
		MyPrescriptions int64 `json:"my_prescriptions"`
	}

	// Get my visits today
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	h.db.Model(&models.Visit{}).Where("clinic_id = ? AND staff_id = ? AND visit_date >= ? AND visit_date < ?", clinicID, staffID, today, tomorrow).Count(&stats.MyVisitsToday)

	// Get my total visits
	h.db.Model(&models.Visit{}).Where("clinic_id = ? AND staff_id = ?", clinicID, staffID).Count(&stats.MyTotalVisits)

	// Get total patients in clinic
	h.db.Model(&models.Patient{}).Where("clinic_id = ?", clinicID).Count(&stats.TotalPatients)

	// Get total staff in clinic
	h.db.Model(&models.Staff{}).Where("clinic_id = ? AND is_active = ?", clinicID, true).Count(&stats.TotalStaff)

	// Get my diagnoses (through visits)
	h.db.Table("diagnoses").
		Joins("JOIN visits ON diagnoses.visit_id = visits.id").
		Where("visits.clinic_id = ? AND visits.staff_id = ?", clinicID, staffID).
		Count(&stats.MyDiagnoses)

	// Get my prescriptions (through visits)
	h.db.Table("prescriptions").
		Joins("JOIN visits ON prescriptions.visit_id = visits.id").
		Where("visits.clinic_id = ? AND visits.staff_id = ?", clinicID, staffID).
		Count(&stats.MyPrescriptions)

	return c.JSON(stats)
}

// GetMyDiagnoses returns diagnoses for patients that this doctor has treated
func (h *MedicalPortalHandler) GetMyDiagnoses(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)
	staffID := c.Locals("staff_id").(uint)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	if perPage > 100 {
		perPage = 100
	}

	visitID := c.Query("visit_id")
	patientID := c.Query("patient_id")

	offset := (page - 1) * perPage

	var diagnoses []models.Diagnosis
	var total int64

	// Join with visits to ensure doctor can only see diagnoses from their visits
	query := h.db.Model(&models.Diagnosis{}).
		Joins("JOIN visits ON diagnoses.visit_id = visits.id").
		Where("visits.clinic_id = ? AND visits.staff_id = ?", clinicID, staffID).
		Preload("Visit").
		Preload("Visit.Patient").
		Preload("Visit.Staff")

	// Additional filters
	if visitID != "" {
		query = query.Where("diagnoses.visit_id = ?", visitID)
	}

	if patientID != "" {
		query = query.Where("visits.patient_id = ?", patientID)
	}

	query.Count(&total)

	if err := query.Offset(offset).Limit(perPage).Order("diagnoses.created_at DESC").Find(&diagnoses).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch diagnoses",
		})
	}

	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	return c.JSON(models.PaginationResponse{
		Data:       diagnoses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GetMyPrescriptions returns prescriptions for patients that this doctor has treated
func (h *MedicalPortalHandler) GetMyPrescriptions(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)
	staffID := c.Locals("staff_id").(uint)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	if perPage > 100 {
		perPage = 100
	}

	visitID := c.Query("visit_id")
	patientID := c.Query("patient_id")
	activeOnly := c.Query("active_only") == "true"

	offset := (page - 1) * perPage

	var prescriptions []models.Prescription
	var total int64

	// Join with visits to ensure doctor can only see prescriptions from their visits
	query := h.db.Model(&models.Prescription{}).
		Joins("JOIN visits ON prescriptions.visit_id = visits.id").
		Where("visits.clinic_id = ? AND visits.staff_id = ?", clinicID, staffID).
		Preload("Visit").
		Preload("Visit.Patient").
		Preload("Visit.Staff")

	// Additional filters
	if visitID != "" {
		query = query.Where("prescriptions.visit_id = ?", visitID)
	}

	if patientID != "" {
		query = query.Where("visits.patient_id = ?", patientID)
	}

	// Filter for active prescriptions (not expired)
	if activeOnly {
		// Calculate expiry date: created_at + duration_days (PostgreSQL syntax)
		query = query.Where("prescriptions.created_at + INTERVAL '1 day' * prescriptions.duration_days > ?", time.Now())
	}

	query.Count(&total)

	if err := query.Offset(offset).Limit(perPage).Order("prescriptions.created_at DESC").Find(&prescriptions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch prescriptions",
		})
	}

	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	return c.JSON(models.PaginationResponse{
		Data:       prescriptions,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GetMyDiagnosis returns a specific diagnosis if it belongs to this doctor's patient
func (h *MedicalPortalHandler) GetMyDiagnosis(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)
	staffID := c.Locals("staff_id").(uint)
	diagnosisID := c.Params("id")

	var diagnosis models.Diagnosis
	if err := h.db.Model(&models.Diagnosis{}).
		Joins("JOIN visits ON diagnoses.visit_id = visits.id").
		Where("diagnoses.id = ? AND visits.clinic_id = ? AND visits.staff_id = ?", diagnosisID, clinicID, staffID).
		Preload("Visit").
		Preload("Visit.Patient").
		Preload("Visit.Staff").
		First(&diagnosis).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Diagnosis not found or access denied",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch diagnosis",
		})
	}

	return c.JSON(diagnosis)
}

// GetMyPrescription returns a specific prescription if it belongs to this doctor's patient
func (h *MedicalPortalHandler) GetMyPrescription(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)
	staffID := c.Locals("staff_id").(uint)
	prescriptionID := c.Params("id")

	var prescription models.Prescription
	if err := h.db.Model(&models.Prescription{}).
		Joins("JOIN visits ON prescriptions.visit_id = visits.id").
		Where("prescriptions.id = ? AND visits.clinic_id = ? AND visits.staff_id = ?", prescriptionID, clinicID, staffID).
		Preload("Visit").
		Preload("Visit.Patient").
		Preload("Visit.Staff").
		First(&prescription).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Prescription not found or access denied",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch prescription",
		})
	}

	return c.JSON(prescription)
}
