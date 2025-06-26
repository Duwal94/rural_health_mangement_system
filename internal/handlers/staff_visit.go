package handlers

import (
	"math"
	"strconv"
	"time"

	"rural_health_management_system/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// StaffHandler handles staff-related operations
type StaffHandler struct {
	db *gorm.DB
}

func NewStaffHandler(db *gorm.DB) *StaffHandler {
	return &StaffHandler{db: db}
}

func (h *StaffHandler) GetStaff(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	clinicID := c.Query("clinic_id", "")
	role := c.Query("role", "")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage
	query := h.db.Model(&models.Staff{}).Preload("Clinic")

	if clinicID != "" {
		query = query.Where("clinic_id = ?", clinicID)
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}

	var total int64
	query.Count(&total)

	var staff []models.Staff
	if err := query.Offset(offset).Limit(perPage).Find(&staff).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to fetch staff",
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return c.JSON(models.PaginationResponse{
		Data:       staff,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

func (h *StaffHandler) GetStaffByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var staff models.Staff
	if err := h.db.Preload("Clinic").First(&staff, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "Staff member not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to fetch staff member",
		})
	}

	return c.JSON(staff)
}

func (h *StaffHandler) CreateStaff(c *fiber.Ctx) error {
	var staff models.Staff

	if err := c.BodyParser(&staff); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Invalid input",
		})
	}

	// Validate required fields
	if staff.FullName == "" || staff.Role == "" || staff.Email == "" || staff.ClinicID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Missing required fields",
		})
	}

	// Validate role
	validRoles := []string{"Doctor", "Nurse", "Administrator", "Pharmacist"}
	isValidRole := false
	for _, role := range validRoles {
		if staff.Role == role {
			isValidRole = true
			break
		}
	}
	if !isValidRole {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Invalid role. Must be Doctor, Nurse, Administrator, or Pharmacist",
		})
	}

	// Check if clinic exists
	var clinic models.Clinic
	if err := h.db.First(&clinic, staff.ClinicID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Clinic not found",
		})
	}

	if err := h.db.Create(&staff).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to create staff member",
		})
	}

	h.db.Preload("Clinic").First(&staff, staff.ID)
	return c.Status(fiber.StatusCreated).JSON(staff)
}

func (h *StaffHandler) UpdateStaff(c *fiber.Ctx) error {
	id := c.Params("id")

	var staff models.Staff
	if err := h.db.First(&staff, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "Staff member not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to fetch staff member",
		})
	}

	var updates models.Staff
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Invalid input",
		})
	}

	// Update fields
	if updates.FullName != "" {
		staff.FullName = updates.FullName
	}
	if updates.Role != "" {
		validRoles := []string{"Doctor", "Nurse", "Administrator", "Pharmacist"}
		isValidRole := false
		for _, role := range validRoles {
			if updates.Role == role {
				isValidRole = true
				break
			}
		}
		if !isValidRole {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error: "Invalid role",
			})
		}
		staff.Role = updates.Role
	}
	if updates.Phone != "" {
		staff.Phone = updates.Phone
	}
	if updates.Email != "" {
		staff.Email = updates.Email
	}
	if updates.ClinicID != 0 {
		var clinic models.Clinic
		if err := h.db.First(&clinic, updates.ClinicID).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error: "Clinic not found",
			})
		}
		staff.ClinicID = updates.ClinicID
	}

	staff.UpdatedAt = time.Now()

	if err := h.db.Save(&staff).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to update staff member",
		})
	}

	h.db.Preload("Clinic").First(&staff, staff.ID)
	return c.JSON(staff)
}

func (h *StaffHandler) DeleteStaff(c *fiber.Ctx) error {
	id := c.Params("id")

	var staff models.Staff
	if err := h.db.First(&staff, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "Staff member not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to fetch staff member",
		})
	}

	// Check if staff has visits
	var visitCount int64
	h.db.Model(&models.Visit{}).Where("staff_id = ?", id).Count(&visitCount)

	if visitCount > 0 {
		return c.Status(fiber.StatusConflict).JSON(models.ErrorResponse{
			Error: "Cannot delete staff member with existing visits",
		})
	}

	if err := h.db.Delete(&staff).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to delete staff member",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// VisitHandler handles visit-related operations
type VisitHandler struct {
	db *gorm.DB
}

func NewVisitHandler(db *gorm.DB) *VisitHandler {
	return &VisitHandler{db: db}
}

func (h *VisitHandler) GetVisits(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	patientID := c.Query("patient_id", "")
	clinicID := c.Query("clinic_id", "")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage
	query := h.db.Model(&models.Visit{}).Preload("Patient").Preload("Clinic").Preload("Staff")

	if patientID != "" {
		query = query.Where("patient_id = ?", patientID)
	}
	if clinicID != "" {
		query = query.Where("clinic_id = ?", clinicID)
	}

	var total int64
	query.Count(&total)

	var visits []models.Visit
	if err := query.Order("visit_date DESC").Offset(offset).Limit(perPage).Find(&visits).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to fetch visits",
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return c.JSON(models.PaginationResponse{
		Data:       visits,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

func (h *VisitHandler) GetVisitByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var visit models.Visit
	if err := h.db.Preload("Patient").Preload("Clinic").Preload("Staff").
		Preload("Diagnoses").Preload("Prescriptions").First(&visit, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "Visit not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to fetch visit",
		})
	}

	return c.JSON(visit)
}

func (h *VisitHandler) CreateVisit(c *fiber.Ctx) error {
	var visit models.Visit

	if err := c.BodyParser(&visit); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Invalid input",
		})
	}

	// Validate required fields
	if visit.PatientID == 0 || visit.ClinicID == 0 || visit.StaffID == 0 || visit.Reason == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Missing required fields",
		})
	}

	// Validate foreign keys exist
	var patient models.Patient
	if err := h.db.First(&patient, visit.PatientID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Patient not found",
		})
	}

	var clinic models.Clinic
	if err := h.db.First(&clinic, visit.ClinicID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Clinic not found",
		})
	}

	var staff models.Staff
	if err := h.db.First(&staff, visit.StaffID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Staff member not found",
		})
	}

	// Set visit date to now if not provided
	if visit.VisitDate.IsZero() {
		visit.VisitDate = time.Now()
	}

	if err := h.db.Create(&visit).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to create visit",
		})
	}

	h.db.Preload("Patient").Preload("Clinic").Preload("Staff").First(&visit, visit.ID)
	return c.Status(fiber.StatusCreated).JSON(visit)
}

func (h *VisitHandler) UpdateVisit(c *fiber.Ctx) error {
	id := c.Params("id")

	var visit models.Visit
	if err := h.db.First(&visit, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "Visit not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to fetch visit",
		})
	}

	var updates models.Visit
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Invalid input",
		})
	}

	// Update fields if provided
	if updates.PatientID != 0 {
		var patient models.Patient
		if err := h.db.First(&patient, updates.PatientID).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error: "Patient not found",
			})
		}
		visit.PatientID = updates.PatientID
	}
	if updates.ClinicID != 0 {
		var clinic models.Clinic
		if err := h.db.First(&clinic, updates.ClinicID).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error: "Clinic not found",
			})
		}
		visit.ClinicID = updates.ClinicID
	}
	if updates.StaffID != 0 {
		var staff models.Staff
		if err := h.db.First(&staff, updates.StaffID).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error: "Staff member not found",
			})
		}
		visit.StaffID = updates.StaffID
	}
	if !updates.VisitDate.IsZero() {
		visit.VisitDate = updates.VisitDate
	}
	if updates.Reason != "" {
		visit.Reason = updates.Reason
	}
	if updates.Notes != "" {
		visit.Notes = updates.Notes
	}

	visit.UpdatedAt = time.Now()

	if err := h.db.Save(&visit).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to update visit",
		})
	}

	h.db.Preload("Patient").Preload("Clinic").Preload("Staff").First(&visit, visit.ID)
	return c.JSON(visit)
}

func (h *VisitHandler) DeleteVisit(c *fiber.Ctx) error {
	id := c.Params("id")

	var visit models.Visit
	if err := h.db.First(&visit, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "Visit not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to fetch visit",
		})
	}

	// Check if visit has diagnoses or prescriptions
	var diagnosisCount, prescriptionCount int64
	h.db.Model(&models.Diagnosis{}).Where("visit_id = ?", id).Count(&diagnosisCount)
	h.db.Model(&models.Prescription{}).Where("visit_id = ?", id).Count(&prescriptionCount)

	if diagnosisCount > 0 || prescriptionCount > 0 {
		return c.Status(fiber.StatusConflict).JSON(models.ErrorResponse{
			Error: "Cannot delete visit with existing diagnoses or prescriptions",
		})
	}

	if err := h.db.Delete(&visit).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to delete visit",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
