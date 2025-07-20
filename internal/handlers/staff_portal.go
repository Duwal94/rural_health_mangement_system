package handlers

import (
	"rural_health_management_system/internal/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type StaffPortalHandler struct {
	db *gorm.DB
}

func NewStaffPortalHandler(db *gorm.DB) *StaffPortalHandler {
	return &StaffPortalHandler{db: db}
}

// GetMyProfile returns the clinic's own profile
func (h *StaffPortalHandler) GetMyProfile(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)

	var clinic models.Clinic
	if err := h.db.First(&clinic, clinicID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Clinic not found",
		})
	}

	return c.JSON(clinic)
}

// UpdateMyProfile allows clinic staff to update their clinic profile
func (h *StaffPortalHandler) UpdateMyProfile(c *fiber.Ctx) error {
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

// GetDashboardStats returns dashboard statistics for clinic staff
func (h *StaffPortalHandler) GetDashboardStats(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)

	var stats struct {
		TotalPatients int64 `json:"total_patients"`
		TotalStaff    int64 `json:"total_staff"`
		TotalVisits   int64 `json:"total_visits"`
		VisitsToday   int64 `json:"visits_today"`
		ActiveDoctors int64 `json:"active_doctors"`
		ActiveNurses  int64 `json:"active_nurses"`
	}

	// Get total patients
	h.db.Model(&models.Patient{}).Where("clinic_id = ?", clinicID).Count(&stats.TotalPatients)

	// Get total staff
	h.db.Model(&models.Staff{}).Where("clinic_id = ? AND is_active = ?", clinicID, true).Count(&stats.TotalStaff)

	// Get total visits
	h.db.Model(&models.Visit{}).Where("clinic_id = ?", clinicID).Count(&stats.TotalVisits)

	// Get visits today
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	h.db.Model(&models.Visit{}).Where("clinic_id = ? AND visit_date >= ? AND visit_date < ?", clinicID, today, tomorrow).Count(&stats.VisitsToday)

	// Get active doctors
	h.db.Model(&models.Staff{}).Where("clinic_id = ? AND role = ? AND is_active = ?", clinicID, "Doctor", true).Count(&stats.ActiveDoctors)

	// Get active nurses
	h.db.Model(&models.Staff{}).Where("clinic_id = ? AND role = ? AND is_active = ?", clinicID, "Nurse", true).Count(&stats.ActiveNurses)

	return c.JSON(stats)
}

// CreatePatient creates a new patient (staff only)
func (h *StaffPortalHandler) CreatePatient(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)

	var req models.CreatePatientRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Force the clinic ID to be the staff's clinic
	req.ClinicID = clinicID

	// Parse date of birth
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid date format. Use YYYY-MM-DD",
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create patient",
		})
	}

	// Load clinic relationship
	h.db.Preload("Clinic").First(&patient, patient.ID)

	return c.Status(fiber.StatusCreated).JSON(patient)
}

// GetMyPatients returns all patients for the clinic
func (h *StaffPortalHandler) GetMyPatients(c *fiber.Ctx) error {
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
func (h *StaffPortalHandler) GetMyPatient(c *fiber.Ctx) error {
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

// CreateStaff creates a new staff member (staff only)
func (h *StaffPortalHandler) CreateStaff(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)

	var req models.CreateStaffRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Check if email already exists
	var existingUser models.User
	if err := h.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Email already registered",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Determine user type based on role
	var userType string
	switch req.Role {
	case "Doctor":
		userType = "doctor"
	case "Nurse":
		userType = "nurse"
	case "Clinic_Administrator", "Pharmacist":
		// These roles don't get login accounts
		userType = ""
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role",
		})
	}

	// Start transaction
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user *models.User
	// Create user only for Doctor and Nurse roles
	if userType != "" {
		user = &models.User{
			Email:    req.Email,
			Password: string(hashedPassword),
			UserType: userType,
			IsActive: true,
		}
		if err := tx.Create(user).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create user",
			})
		}
	}

	// Create staff
	staff := models.Staff{
		FullName: req.FullName,
		Role:     req.Role,
		Phone:    req.Phone,
		Email:    req.Email,
		ClinicID: clinicID,
		IsActive: true,
	}

	if user != nil {
		staff.UserID = &user.ID
	}

	if err := tx.Create(&staff).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create staff",
		})
	}

	tx.Commit()

	// Load relationships
	h.db.Preload("Clinic").Preload("User").First(&staff, staff.ID)

	return c.Status(fiber.StatusCreated).JSON(staff)
}

// GetMyStaff returns all staff for the clinic
func (h *StaffPortalHandler) GetMyStaff(c *fiber.Ctx) error {
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
		query = query.Where("full_name ILIKE ? OR role ILIKE ? OR phone ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
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

// CreateVisit creates a new visit (staff only)
func (h *StaffPortalHandler) CreateVisit(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)

	var req models.CreateVisitRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Force the clinic ID to be the staff's clinic
	req.ClinicID = clinicID

	// Verify patient belongs to this clinic
	var patient models.Patient
	if err := h.db.Where("id = ? AND clinic_id = ?", req.PatientID, clinicID).First(&patient).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Patient not found in this clinic",
		})
	}

	// Verify staff belongs to this clinic
	var staff models.Staff
	if err := h.db.Where("id = ? AND clinic_id = ?", req.StaffID, clinicID).First(&staff).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Staff not found in this clinic",
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

// GetMyVisits returns all visits for the clinic
func (h *StaffPortalHandler) GetMyVisits(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	if perPage > 100 {
		perPage = 100
	}
	patientID := c.Query("patient_id")
	staffID := c.Query("staff_id")

	offset := (page - 1) * perPage

	var visits []models.Visit
	var total int64

	query := h.db.Model(&models.Visit{}).Where("clinic_id = ?", clinicID)

	if patientID != "" {
		query = query.Where("patient_id = ?", patientID)
	}

	if staffID != "" {
		query = query.Where("staff_id = ?", staffID)
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
func (h *StaffPortalHandler) GetMyVisit(c *fiber.Ctx) error {
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
