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
	token, err := h.generateToken(user.ID, user.Email, user.UserType, &patient.ID, nil)
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

	// Create user
	user := models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		UserType: "clinic",
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
	token, err := h.generateToken(user.ID, user.Email, user.UserType, nil, &clinic.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.LoginResponse{
		Token:    token,
		UserType: "clinic",
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

	var patientID, clinicID *uint
	var userProfile interface{}

	// Load user profile based on type
	switch user.UserType {
	case "patient":
		var patient models.Patient
		if err := h.db.Preload("Clinic").Where("user_id = ?", user.ID).First(&patient).Error; err == nil {
			patientID = &patient.ID
			userProfile = patient
		}
	case "clinic":
		var clinic models.Clinic
		if err := h.db.Where("user_id = ?", user.ID).First(&clinic).Error; err == nil {
			clinicID = &clinic.ID
			userProfile = clinic
		}
	}

	// Generate JWT token
	token, err := h.generateToken(user.ID, user.Email, user.UserType, patientID, clinicID)
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

	case "clinic":
		var clinic models.Clinic
		if err := h.db.Where("user_id = ?", userID).First(&clinic).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Clinic profile not found",
			})
		}
		return c.JSON(clinic)

	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user type",
		})
	}
}

// generateToken creates a JWT token
func (h *AuthHandler) generateToken(userID uint, email, userType string, patientID, clinicID *uint) (string, error) {
	claims := models.JWTClaims{
		UserID:    userID,
		Email:     email,
		UserType:  userType,
		PatientID: patientID,
		ClinicID:  clinicID,
		Exp:       time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    claims.UserID,
		"email":      claims.Email,
		"user_type":  claims.UserType,
		"patient_id": claims.PatientID,
		"clinic_id":  claims.ClinicID,
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
