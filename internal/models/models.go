package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents the authentication entity
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"not null;size:255;uniqueIndex" validate:"required,email"`
	Password  string         `json:"-" gorm:"not null;size:255" validate:"required,min=8"`
	UserType  string         `json:"user_type" gorm:"not null;size:20" validate:"required,oneof=patient clinic admin"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	PatientProfile *Patient `json:"patient_profile,omitempty" gorm:"foreignKey:UserID"`
	ClinicProfile  *Clinic  `json:"clinic_profile,omitempty" gorm:"foreignKey:UserID"`
}

type Clinic struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Name          string         `json:"name" gorm:"not null;size:255" validate:"required,min=2,max=255"`
	Address       string         `json:"address" gorm:"not null;size:500" validate:"required,min=5,max=500"`
	ContactNumber string         `json:"contact_number" gorm:"not null;size:20" validate:"required,min=10,max=20"`
	District      string         `json:"district" gorm:"not null;size:100" validate:"required,min=2,max=100"`
	UserID        *uint          `json:"user_id,omitempty" gorm:"index"` // Link to User for authentication
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User     *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Patients []Patient `json:"patients,omitempty" gorm:"foreignKey:ClinicID"`
	Staff    []Staff   `json:"staff,omitempty" gorm:"foreignKey:ClinicID"`
	Visits   []Visit   `json:"visits,omitempty" gorm:"foreignKey:ClinicID"`
}

type Patient struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	FullName    string         `json:"full_name" gorm:"not null;size:255" validate:"required,min=2,max=255"`
	Gender      string         `json:"gender" gorm:"not null;size:10" validate:"required,oneof=Male Female Other"`
	DateOfBirth time.Time      `json:"date_of_birth" gorm:"not null" validate:"required"`
	Address     string         `json:"address" gorm:"not null;size:500" validate:"required,min=5,max=500"`
	Phone       string         `json:"phone" gorm:"not null;size:20" validate:"required,min=10,max=20"`
	ClinicID    uint           `json:"clinic_id" gorm:"not null" validate:"required"`
	UserID      *uint          `json:"user_id,omitempty" gorm:"index"` // Link to User for authentication
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User   *User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Clinic *Clinic `json:"clinic,omitempty" gorm:"foreignKey:ClinicID;references:ID"`
	Visits []Visit `json:"visits,omitempty" gorm:"foreignKey:PatientID"`
}

type Staff struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	FullName  string         `json:"full_name" gorm:"not null;size:255" validate:"required,min=2,max=255"`
	Role      string         `json:"role" gorm:"not null;size:100" validate:"required,oneof=Doctor Nurse Administrator Pharmacist"`
	Phone     string         `json:"phone" gorm:"not null;size:20" validate:"required,min=10,max=20"`
	Email     string         `json:"email" gorm:"not null;size:255;uniqueIndex" validate:"required,email"`
	ClinicID  uint           `json:"clinic_id" gorm:"not null" validate:"required"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Clinic *Clinic `json:"clinic,omitempty" gorm:"foreignKey:ClinicID;references:ID"`
	Visits []Visit `json:"visits,omitempty" gorm:"foreignKey:StaffID"`
}

type Visit struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	PatientID uint           `json:"patient_id" gorm:"not null" validate:"required"`
	ClinicID  uint           `json:"clinic_id" gorm:"not null" validate:"required"`
	StaffID   uint           `json:"staff_id" gorm:"not null" validate:"required"`
	VisitDate time.Time      `json:"visit_date" gorm:"not null" validate:"required"`
	Reason    string         `json:"reason" gorm:"not null;size:500" validate:"required,min=5,max=500"`
	Notes     string         `json:"notes" gorm:"size:1000" validate:"max=1000"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Patient       *Patient       `json:"patient,omitempty" gorm:"foreignKey:PatientID;references:ID"`
	Clinic        *Clinic        `json:"clinic,omitempty" gorm:"foreignKey:ClinicID;references:ID"`
	Staff         *Staff         `json:"staff,omitempty" gorm:"foreignKey:StaffID;references:ID"`
	Diagnoses     []Diagnosis    `json:"diagnoses,omitempty" gorm:"foreignKey:VisitID"`
	Prescriptions []Prescription `json:"prescriptions,omitempty" gorm:"foreignKey:VisitID"`
}

type Diagnosis struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	VisitID       uint           `json:"visit_id" gorm:"not null" validate:"required"`
	DiagnosisCode string         `json:"diagnosis_code" gorm:"not null;size:20" validate:"required,min=2,max=20"`
	Description   string         `json:"description" gorm:"not null;size:1000" validate:"required,min=5,max=1000"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Visit *Visit `json:"visit,omitempty" gorm:"foreignKey:VisitID;references:ID"`
}

type Prescription struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	VisitID        uint           `json:"visit_id" gorm:"not null" validate:"required"`
	MedicationName string         `json:"medication_name" gorm:"not null;size:255" validate:"required,min=2,max=255"`
	Dosage         string         `json:"dosage" gorm:"not null;size:100" validate:"required,min=2,max=100"`
	Instructions   string         `json:"instructions" gorm:"not null;size:500" validate:"required,min=5,max=500"`
	DurationDays   int            `json:"duration_days" gorm:"not null" validate:"required,min=1,max=365"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Visit *Visit `json:"visit,omitempty" gorm:"foreignKey:VisitID;references:ID"`
}

// Request/Response DTOs
type CreatePatientRequest struct {
	FullName    string `json:"full_name" validate:"required,min=2,max=255"`
	Gender      string `json:"gender" validate:"required,oneof=Male Female Other"`
	DateOfBirth string `json:"date_of_birth" validate:"required"` // Will be parsed to time.Time
	Address     string `json:"address" validate:"required,min=5,max=500"`
	Phone       string `json:"phone" validate:"required,min=10,max=20"`
	ClinicID    uint   `json:"clinic_id" validate:"required"`
}

type UpdatePatientRequest struct {
	FullName    *string `json:"full_name,omitempty" validate:"omitempty,min=2,max=255"`
	Gender      *string `json:"gender,omitempty" validate:"omitempty,oneof=Male Female Other"`
	DateOfBirth *string `json:"date_of_birth,omitempty"`
	Address     *string `json:"address,omitempty" validate:"omitempty,min=5,max=500"`
	Phone       *string `json:"phone,omitempty" validate:"omitempty,min=10,max=20"`
	ClinicID    *uint   `json:"clinic_id,omitempty"`
}

type CreateVisitRequest struct {
	PatientID uint      `json:"patient_id" validate:"required"`
	ClinicID  uint      `json:"clinic_id" validate:"required"`
	StaffID   uint      `json:"staff_id" validate:"required"`
	VisitDate time.Time `json:"visit_date,omitempty"`
	Reason    string    `json:"reason" validate:"required,min=5,max=500"`
	Notes     string    `json:"notes,omitempty" validate:"max=1000"`
}

type UpdateVisitRequest struct {
	PatientID *uint      `json:"patient_id,omitempty"`
	ClinicID  *uint      `json:"clinic_id,omitempty"`
	StaffID   *uint      `json:"staff_id,omitempty"`
	VisitDate *time.Time `json:"visit_date,omitempty"`
	Reason    *string    `json:"reason,omitempty" validate:"omitempty,min=5,max=500"`
	Notes     *string    `json:"notes,omitempty" validate:"omitempty,max=1000"`
}

type CreateDiagnosisRequest struct {
	VisitID       uint   `json:"visit_id" validate:"required"`
	DiagnosisCode string `json:"diagnosis_code" validate:"required,min=2,max=20"`
	Description   string `json:"description" validate:"required,min=5,max=1000"`
}

type CreatePrescriptionRequest struct {
	VisitID        uint   `json:"visit_id" validate:"required"`
	MedicationName string `json:"medication_name" validate:"required,min=2,max=255"`
	Dosage         string `json:"dosage" validate:"required,min=2,max=100"`
	Instructions   string `json:"instructions" validate:"required,min=5,max=500"`
	DurationDays   int    `json:"duration_days" validate:"required,min=1,max=365"`
}

type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

type ErrorResponse struct {
	Error   string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}

// Authentication DTOs
type RegisterPatientRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	FullName    string `json:"full_name" validate:"required,min=2,max=255"`
	Gender      string `json:"gender" validate:"required,oneof=Male Female Other"`
	DateOfBirth string `json:"date_of_birth" validate:"required"` // Will be parsed to time.Time
	Address     string `json:"address" validate:"required,min=5,max=500"`
	Phone       string `json:"phone" validate:"required,min=10,max=20"`
	ClinicID    uint   `json:"clinic_id" validate:"required"`
}

type RegisterClinicRequest struct {
	Email         string `json:"email" validate:"required,email"`
	Password      string `json:"password" validate:"required,min=8"`
	Name          string `json:"name" validate:"required,min=2,max=255"`
	Address       string `json:"address" validate:"required,min=5,max=500"`
	ContactNumber string `json:"contact_number" validate:"required,min=10,max=20"`
	District      string `json:"district" validate:"required,min=2,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token    string      `json:"token"`
	UserType string      `json:"user_type"`
	User     interface{} `json:"user"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

type JWTClaims struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	UserType  string `json:"user_type"`
	PatientID *uint  `json:"patient_id,omitempty"`
	ClinicID  *uint  `json:"clinic_id,omitempty"`
	Exp       int64  `json:"exp"`
}
