package handlers

import (
	"fmt"
	"rural_health_management_system/internal/models"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db     *gorm.DB
	jwtKey []byte
}

func NewAuthHandler(db *gorm.DB, jwtKey string) *AuthHandler {
	return &AuthHandler{
		db:     db,
		jwtKey: []byte(jwtKey),
	}
}

// RegisterPatient registers a new patient with authentication
func (h *AuthHandler) RegisterPatient(c *fiber.Ctx) error {
	var req models.RegisterPatientRequest
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

	// Check if clinic exists
	var clinic models.Clinic
	if err := h.db.First(&clinic, req.ClinicID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid clinic ID",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Parse date of birth
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid date format. Use YYYY-MM-DD",
		})
	}

	// Start transaction
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create user
	user := models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		UserType: "patient",
		IsActive: true,
	}
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// Create patient
	patient := models.Patient{
		FullName:    req.FullName,
		Gender:      req.Gender,
		DateOfBirth: dob,
		Address:     req.Address,
		Phone:       req.Phone,
		ClinicID:    req.ClinicID,
		UserID:      &user.ID,
	}
	if err := tx.Create(&patient).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create patient",
		})
	}

	tx.Commit()

	// Generate JWT token
	token, err := h.generateToken(user.ID, user.Email, user.UserType, &patient.ID, nil, nil, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Load patient with clinic
	h.db.Preload("Clinic").First(&patient, patient.ID)

	return c.Status(fiber.StatusCreated).JSON(models.LoginResponse{
		Token:    token,
		UserType: "patient",
		User:     patient,
	})
}

// RegisterClinic registers a new clinic with authentication
func (h *AuthHandler) RegisterClinic(c *fiber.Ctx) error {
	var req models.RegisterClinicRequest
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

	// Start transaction
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create user with clinic_staff type
	user := models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		UserType: "clinic_staff",
		IsActive: true,
	}
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// Create clinic
	clinic := models.Clinic{
		Name:          req.Name,
		Address:       req.Address,
		ContactNumber: req.ContactNumber,
		District:      req.District,
		UserID:        &user.ID,
	}
	if err := tx.Create(&clinic).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create clinic",
		})
	}

	tx.Commit()

	// Generate JWT token
	token, err := h.generateToken(user.ID, user.Email, user.UserType, nil, &clinic.ID, nil, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.LoginResponse{
		Token:    token,
		UserType: "clinic_staff",
		User:     clinic,
	})
}

// Login authenticates a user and returns a JWT token
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Find user
	var user models.User
	if err := h.db.Where("email = ? AND is_active = true", req.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	var patientID, clinicID, staffID *uint
	var staffRole *string
	var userProfile interface{}

	// Load user profile based on type
	switch user.UserType {
	case "patient":
		var patient models.Patient
		if err := h.db.Preload("Clinic").Where("user_id = ?", user.ID).First(&patient).Error; err == nil {
			patientID = &patient.ID
			userProfile = patient
		}
	case "clinic_staff":
		var clinic models.Clinic
		if err := h.db.Where("user_id = ?", user.ID).First(&clinic).Error; err == nil {
			clinicID = &clinic.ID
			userProfile = clinic
		}
	case "doctor", "nurse":
		var staff models.Staff
		if err := h.db.Preload("Clinic").Where("user_id = ?", user.ID).First(&staff).Error; err == nil {
			staffID = &staff.ID
			staffRole = &staff.Role
			clinicID = &staff.ClinicID
			userProfile = staff
		}
	}

	// Generate JWT token
	token, err := h.generateToken(user.ID, user.Email, user.UserType, patientID, clinicID, staffID, staffRole)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(models.LoginResponse{
		Token:    token,
		UserType: user.UserType,
		User:     userProfile,
	})
}

// ChangePassword allows users to change their password
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req models.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get user
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Current password is incorrect",
		})
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Update password
	if err := h.db.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update password",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Password updated successfully",
	})
}

// GetProfile returns the current user's profile
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userType := c.Locals("user_type").(string)
	userID := c.Locals("user_id").(uint)

	switch userType {
	case "patient":
		var patient models.Patient
		if err := h.db.Preload("Clinic").Where("user_id = ?", userID).First(&patient).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Patient profile not found",
			})
		}
		return c.JSON(patient)

	case "clinic_admin":
		var clinic models.Clinic
		if err := h.db.Where("user_id = ?", userID).First(&clinic).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Clinic profile not found",
			})
		}
		return c.JSON(clinic)

	case "doctor", "nurse":
		var staff models.Staff
		if err := h.db.Preload("Clinic").Where("user_id = ?", userID).First(&staff).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Staff profile not found",
			})
		}
		return c.JSON(staff)

	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user type",
		})
	}
}

// RegisterStaff registers a new staff member with authentication (called by clinic admin)
func (h *AuthHandler) RegisterStaff(c *fiber.Ctx) error {
	var req models.RegisterStaffRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Verify requester is clinic staff and owns the clinic
	requesterType := c.Locals("user_type").(string)
	if requesterType != "clinic_staff" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only clinic staff can register staff",
		})
	}

	clinicID := c.Locals("clinic_id").(uint)
	if req.ClinicID != clinicID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only register staff for your own clinic",
		})
	}

	// Check if email already exists
	var existingUser models.User
	if err := h.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Email already registered",
		})
	}

	// Check if clinic exists
	var clinic models.Clinic
	if err := h.db.First(&clinic, req.ClinicID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid clinic ID",
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Clinic_Administrator and Pharmacist roles cannot have login accounts. Only Doctor and Nurse can login.",
		})
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

	// Create user
	user := models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		UserType: userType,
		IsActive: true,
	}
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// Create staff
	staff := models.Staff{
		FullName: req.FullName,
		Role:     req.Role,
		Phone:    req.Phone,
		Email:    req.Email,
		ClinicID: req.ClinicID,
		UserID:   &user.ID,
		IsActive: true,
	}
	if err := tx.Create(&staff).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create staff",
		})
	}

	tx.Commit()

	// Generate JWT token for the new staff member
	token, err := h.generateToken(user.ID, user.Email, user.UserType, nil, &staff.ClinicID, &staff.ID, &staff.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Load staff with clinic
	h.db.Preload("Clinic").First(&staff, staff.ID)

	return c.Status(fiber.StatusCreated).JSON(models.LoginResponse{
		Token:    token,
		UserType: userType,
		User:     staff,
	})
}

// ClinicLogin authenticates clinic users with specific login types
func (h *AuthHandler) ClinicLogin(c *fiber.Ctx) error {
	var req models.ClinicLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Find user
	var user models.User
	if err := h.db.Where("email = ? AND is_active = true", req.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Validate login type against user type
	switch req.LoginType {
	case models.ClinicLoginStaff:
		if user.UserType != "clinic_staff" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Invalid login type for this user",
			})
		}
	case models.ClinicLoginMedical:
		if user.UserType != "doctor" && user.UserType != "nurse" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Invalid login type for this user",
			})
		}
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid login type",
		})
	}

	var patientID, clinicID, staffID *uint
	var staffRole *string
	var userProfile interface{}

	// Load user profile based on type
	switch user.UserType {
	case "clinic_staff":
		var clinic models.Clinic
		if err := h.db.Where("user_id = ?", user.ID).First(&clinic).Error; err == nil {
			clinicID = &clinic.ID
			userProfile = clinic
		}
	case "doctor", "nurse":
		var staff models.Staff
		if err := h.db.Preload("Clinic").Where("user_id = ?", user.ID).First(&staff).Error; err == nil {
			staffID = &staff.ID
			staffRole = &staff.Role
			clinicID = &staff.ClinicID
			userProfile = staff
		}
	}

	// Generate JWT token
	token, err := h.generateToken(user.ID, user.Email, user.UserType, patientID, clinicID, staffID, staffRole)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(models.LoginResponse{
		Token:    token,
		UserType: user.UserType,
		User:     userProfile,
	})
}

// generateToken creates a JWT token
func (h *AuthHandler) generateToken(userID uint, email, userType string, patientID, clinicID, staffID *uint, staffRole *string) (string, error) {
	claims := models.JWTClaims{
		UserID:    userID,
		Email:     email,
		UserType:  userType,
		PatientID: patientID,
		ClinicID:  clinicID,
		StaffID:   staffID,
		StaffRole: staffRole,
		Exp:       time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    claims.UserID,
		"email":      claims.Email,
		"user_type":  claims.UserType,
		"patient_id": claims.PatientID,
		"clinic_id":  claims.ClinicID,
		"staff_id":   claims.StaffID,
		"staff_role": claims.StaffRole,
		"exp":        claims.Exp,
	})

	return token.SignedString(h.jwtKey)
}

// AuthMiddleware validates JWT tokens
func (h *AuthHandler) AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header required",
		})
	}

	// Extract token from "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid authorization header format",
		})
	}

	tokenString := parts[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return h.jwtKey, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}

	// Check token expiration
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token expired",
			})
		}
	}

	// Store user information in context
	c.Locals("user_id", uint(claims["user_id"].(float64)))
	c.Locals("email", claims["email"].(string))
	c.Locals("user_type", claims["user_type"].(string))

	if patientID, ok := claims["patient_id"]; ok && patientID != nil {
		c.Locals("patient_id", uint(patientID.(float64)))
	}

	if clinicID, ok := claims["clinic_id"]; ok && clinicID != nil {
		c.Locals("clinic_id", uint(clinicID.(float64)))
	}

	if staffID, ok := claims["staff_id"]; ok && staffID != nil {
		c.Locals("staff_id", uint(staffID.(float64)))
	}

	if staffRole, ok := claims["staff_role"]; ok && staffRole != nil {
		c.Locals("staff_role", staffRole.(string))
	}

	return c.Next()
}

// RequireUserType middleware ensures the user has the required user type
func (h *AuthHandler) RequireUserType(userTypes ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userType := c.Locals("user_type").(string)

		for _, allowedType := range userTypes {
			if userType == allowedType {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}
}

// RequireOwnership middleware ensures users can only access their own data
func (h *AuthHandler) RequireOwnership(c *fiber.Ctx) error {
	userType := c.Locals("user_type").(string)

	switch userType {
	case "patient":
		patientID := c.Locals("patient_id").(uint)

		// Check if accessing patient data
		if c.Route().Path == "/api/v1/patients/:id" ||
			c.Route().Path == "/api/v1/portal/patient/visits" ||
			c.Route().Path == "/api/v1/portal/patient/visits/:id" {

			if c.Route().Path == "/api/v1/patients/:id" {
				id, err := strconv.ParseUint(c.Params("id"), 10, 32)
				if err != nil || uint(id) != patientID {
					return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
						"error": "You can only access your own data",
					})
				}
			}
		}

	case "clinic":
		clinicID := c.Locals("clinic_id").(uint)

		// For clinic users, they can access their own clinic data and related patients/staff/visits
		// This will be handled in individual handlers
		c.Locals("owner_clinic_id", clinicID)
	}

	return c.Next()
}

// RequireRole middleware ensures the user has one of the required roles
func (h *AuthHandler) RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userType := c.Locals("user_type").(string)
		staffRole, hasStaffRole := c.Locals("staff_role").(string)

		for _, allowedRole := range roles {
			if userType == allowedRole {
				return c.Next()
			}
			// For staff, also check their specific role
			if hasStaffRole && staffRole == allowedRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions for this action",
		})
	}
}

// RequireClinicAccess middleware ensures the user has access to clinic operations
func (h *AuthHandler) RequireClinicAccess() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userType := c.Locals("user_type").(string)

		// Allow clinic staff, doctors, and nurses
		allowedTypes := []string{"clinic_staff", "doctor", "nurse"}
		for _, allowedType := range allowedTypes {
			if userType == allowedType {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied. Only clinic staff can access this resource",
		})
	}
}

// RequireDoctorAccess middleware ensures only doctors can perform medical actions
func (h *AuthHandler) RequireDoctorAccess() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userType := c.Locals("user_type").(string)

		if userType != "doctor" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied. Only doctors can perform this action",
			})
		}

		return c.Next()
	}
}

// RequireClinicStaffAccess middleware ensures only clinic staff can perform administrative actions
func (h *AuthHandler) RequireClinicStaffAccess() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userType := c.Locals("user_type").(string)

		if userType != "clinic_staff" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied. Only clinic staff can perform this action",
			})
		}

		return c.Next()
	}
}

// ValidateClinicOwnership middleware ensures users can only access data from their own clinic
func (h *AuthHandler) ValidateClinicOwnership() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userType := c.Locals("user_type").(string)
		userClinicID, hasClinicID := c.Locals("clinic_id").(uint)

		if !hasClinicID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "No clinic association found",
			})
		}

		// For clinic-related operations, ensure user belongs to the clinic
		switch userType {
		case "clinic_staff", "doctor", "nurse":
			// These user types should have clinic_id in their context
			c.Locals("validated_clinic_id", userClinicID)
			return c.Next()
		default:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied to clinic resources",
			})
		}
	}
}

// RequirePermission middleware ensures the user has the specified permission
func (h *AuthHandler) RequirePermission(permission models.Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userType := c.Locals("user_type").(string)
		var staffRole *string

		if role, hasRole := c.Locals("staff_role").(string); hasRole {
			staffRole = &role
		}

		if !models.HasPermission(userType, staffRole, permission) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions for this action",
			})
		}

		return c.Next()
	}
}

// RequireMultiplePermissions middleware ensures the user has at least one of the specified permissions
func (h *AuthHandler) RequireMultiplePermissions(permissions ...models.Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userType := c.Locals("user_type").(string)
		var staffRole *string

		if role, hasRole := c.Locals("staff_role").(string); hasRole {
			staffRole = &role
		}

		for _, permission := range permissions {
			if models.HasPermission(userType, staffRole, permission) {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions for this action",
		})
	}
}
