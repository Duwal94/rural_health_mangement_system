package main

import (
	"fmt"
	"log"

	"rural_health_management_system/internal/config"
	"rural_health_management_system/internal/database"
	"rural_health_management_system/internal/handlers"
	"rural_health_management_system/internal/models"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db.DB, cfg.JWTSecret)
	clinicHandler := handlers.NewClinicHandler(db.DB)
	patientHandler := handlers.NewPatientHandler(db.DB)
	staffHandler := handlers.NewStaffHandler(db.DB)
	visitHandler := handlers.NewVisitHandler(db.DB)
	diagnosisHandler := handlers.NewDiagnosisHandler(db.DB)
	prescriptionHandler := handlers.NewPrescriptionHandler(db.DB)
	patientPortalHandler := handlers.NewPatientPortalHandler(db.DB)
	// Old clinic portal handler - deprecated
	clinicPortalHandler := handlers.NewClinicPortalHandler(db.DB)
	// New separate portal handlers
	staffPortalHandler := handlers.NewStaffPortalHandler(db.DB)
	medicalPortalHandler := handlers.NewMedicalPortalHandler(db.DB)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Setup middleware
	handlers.SetupMiddleware(app)

	// Health check endpoint
	app.Get("/health", handlers.HealthCheck)

	// API version 1 routes
	v1 := app.Group("/api/v1")

	// Authentication routes (public)
	auth := v1.Group("/auth")
	auth.Post("/register/patient", authHandler.RegisterPatient)
	auth.Post("/register/clinic", authHandler.RegisterClinic)
	auth.Post("/login", authHandler.Login)
	auth.Post("/clinic-login", authHandler.ClinicLogin) // New clinic-specific login
	auth.Get("/clinics", clinicHandler.GetClinics)
	// Protected routes - require authentication
	protected := v1.Group("/", authHandler.AuthMiddleware)
	protected.Post("/auth/change-password", authHandler.ChangePassword)
	protected.Get("/auth/profile", authHandler.GetProfile)
	protected.Post("/auth/register/staff", authHandler.RequireUserType("clinic_staff"), authHandler.RegisterStaff) // Only clinic staff can register staff
	// Patient Portal routes (patient access only)
	patientPortal := v1.Group("/portal/patient", authHandler.AuthMiddleware, authHandler.RequireUserType("patient"))
	patientPortal.Get("/profile", patientPortalHandler.GetMyProfile)
	patientPortal.Put("/profile", patientPortalHandler.UpdateMyProfile)
	patientPortal.Get("/visits", patientPortalHandler.GetMyVisits)
	patientPortal.Get("/visits/:id", patientPortalHandler.GetMyVisit)
	patientPortal.Get("/diagnoses", patientPortalHandler.GetMyDiagnoses)
	patientPortal.Get("/prescriptions", patientPortalHandler.GetMyPrescriptions)

	// Clinic Portal routes (clinic access only) - DEPRECATED, use staff or medical portals
	clinicPortal := v1.Group("/portal/clinic", authHandler.AuthMiddleware, authHandler.RequireUserType("clinic_staff"))
	clinicPortal.Get("/profile", clinicPortalHandler.GetMyProfile)
	clinicPortal.Put("/profile", clinicPortalHandler.UpdateMyProfile)
	clinicPortal.Get("/dashboard", clinicPortalHandler.GetDashboardStats)
	clinicPortal.Get("/patients", clinicPortalHandler.GetMyPatients)
	clinicPortal.Get("/patients/:id", clinicPortalHandler.GetMyPatient)
	clinicPortal.Get("/staff", clinicPortalHandler.GetMyStaff)
	clinicPortal.Post("/staff", clinicPortalHandler.CreateStaff)
	clinicPortal.Get("/visits", clinicPortalHandler.GetMyVisits)
	clinicPortal.Get("/visits/:id", clinicPortalHandler.GetMyVisit)
	clinicPortal.Post("/visits", clinicPortalHandler.CreateVisit)
	clinicPortal.Post("/diagnoses", clinicPortalHandler.CreateDiagnosis)
	clinicPortal.Post("/prescriptions", clinicPortalHandler.CreatePrescription)

	// Staff Portal routes (clinic staff only) - NEW
	staffPortal := v1.Group("/portal/staff", authHandler.AuthMiddleware, authHandler.RequireUserType("clinic_staff"), authHandler.ValidateClinicOwnership())
	staffPortal.Get("/profile", staffPortalHandler.GetMyProfile)
	staffPortal.Put("/profile", staffPortalHandler.UpdateMyProfile)
	staffPortal.Get("/dashboard", staffPortalHandler.GetDashboardStats)

	// Patient management (staff only)
	staffPortal.Post("/patients", authHandler.RequirePermission(models.PermissionCreatePatient), staffPortalHandler.CreatePatient)
	staffPortal.Get("/patients", authHandler.RequirePermission(models.PermissionViewPatient), staffPortalHandler.GetMyPatients)
	staffPortal.Get("/patients/:id", authHandler.RequirePermission(models.PermissionViewPatient), staffPortalHandler.GetMyPatient)

	// Staff management (staff only)
	staffPortal.Post("/staff", authHandler.RequirePermission(models.PermissionCreateStaff), staffPortalHandler.CreateStaff)
	staffPortal.Get("/staff", authHandler.RequirePermission(models.PermissionViewStaff), staffPortalHandler.GetMyStaff)

	// Visit management (staff can create, view all)
	staffPortal.Post("/visits", authHandler.RequirePermission(models.PermissionCreateVisit), staffPortalHandler.CreateVisit)
	staffPortal.Get("/visits", authHandler.RequirePermission(models.PermissionViewVisit), staffPortalHandler.GetMyVisits)
	staffPortal.Get("/visits/:id", authHandler.RequirePermission(models.PermissionViewVisit), staffPortalHandler.GetMyVisit)

	// Medical Portal routes (doctors and nurses only) - NEW
	medicalPortal := v1.Group("/portal/medical", authHandler.AuthMiddleware, authHandler.RequireUserType("doctor", "nurse"), authHandler.ValidateClinicOwnership())
	medicalPortal.Get("/profile", medicalPortalHandler.GetMyProfile)
	medicalPortal.Put("/profile", medicalPortalHandler.UpdateMyProfile)
	medicalPortal.Get("/dashboard", medicalPortalHandler.GetDashboardStats)

	// Patient access (medical staff can view)
	medicalPortal.Get("/patients", authHandler.RequirePermission(models.PermissionViewPatient), medicalPortalHandler.GetMyPatients)
	medicalPortal.Get("/patients/:id", authHandler.RequirePermission(models.PermissionViewPatient), medicalPortalHandler.GetMyPatient)

	// Staff access (medical staff can view)
	medicalPortal.Get("/staff", authHandler.RequirePermission(models.PermissionViewStaff), medicalPortalHandler.GetStaff)

	// Visit management (medical staff can create, view)
	medicalPortal.Post("/visits", authHandler.RequirePermission(models.PermissionCreateVisit), medicalPortalHandler.CreateVisit)
	medicalPortal.Get("/visits", authHandler.RequirePermission(models.PermissionViewVisit), medicalPortalHandler.GetMyVisits)
	medicalPortal.Get("/visits/:id", authHandler.RequirePermission(models.PermissionViewVisit), medicalPortalHandler.GetMyVisit)

	// Medical actions (doctors only)
	medicalPortal.Post("/diagnoses", authHandler.RequireDoctorAccess(), authHandler.RequirePermission(models.PermissionCreateDiagnosis), medicalPortalHandler.CreateDiagnosis)
	medicalPortal.Post("/prescriptions", authHandler.RequireDoctorAccess(), authHandler.RequirePermission(models.PermissionCreatePrescription), medicalPortalHandler.CreatePrescription)

	// Medical data access (doctors and nurses can view their own patients' data)
	medicalPortal.Get("/diagnoses", authHandler.RequirePermission(models.PermissionViewDiagnosis), medicalPortalHandler.GetMyDiagnoses)
	medicalPortal.Get("/diagnoses/:id", authHandler.RequirePermission(models.PermissionViewDiagnosis), medicalPortalHandler.GetMyDiagnosis)
	medicalPortal.Get("/prescriptions", authHandler.RequirePermission(models.PermissionViewPrescription), medicalPortalHandler.GetMyPrescriptions)
	medicalPortal.Get("/prescriptions/:id", authHandler.RequirePermission(models.PermissionViewPrescription), medicalPortalHandler.GetMyPrescription)

	// Admin/System routes - require authentication and admin permissions
	admin := v1.Group("/", authHandler.AuthMiddleware, authHandler.RequireUserType("admin"))

	// Clinic routes (admin only for system management)
	clinics := admin.Group("/clinics")
	clinics.Get("/:id", clinicHandler.GetClinic)
	clinics.Post("/", clinicHandler.CreateClinic)
	clinics.Put("/:id", clinicHandler.UpdateClinic)
	clinics.Delete("/:id", clinicHandler.DeleteClinic)

	// Patient routes (admin only for system management)
	patients := admin.Group("/patients")
	patients.Get("/", patientHandler.GetPatients)
	patients.Get("/:id", patientHandler.GetPatient)
	patients.Post("/", patientHandler.CreatePatient)
	patients.Put("/:id", patientHandler.UpdatePatient)
	patients.Delete("/:id", patientHandler.DeletePatient)

	// Staff routes (admin only for system management)
	staff := admin.Group("/staff")
	staff.Get("/", staffHandler.GetStaff)
	staff.Get("/:id", staffHandler.GetStaffByID)
	staff.Post("/", staffHandler.CreateStaff)
	staff.Put("/:id", staffHandler.UpdateStaff)
	staff.Delete("/:id", staffHandler.DeleteStaff)

	// Visit routes (admin only for system management)
	visits := admin.Group("/visits")
	visits.Get("/", visitHandler.GetVisits)
	visits.Get("/:id", visitHandler.GetVisitByID)
	visits.Post("/", visitHandler.CreateVisit)
	visits.Put("/:id", visitHandler.UpdateVisit)
	visits.Delete("/:id", visitHandler.DeleteVisit)

	// Diagnosis routes (admin only for system management)
	diagnoses := admin.Group("/diagnoses")
	diagnoses.Get("/", diagnosisHandler.GetDiagnoses)
	diagnoses.Get("/:id", diagnosisHandler.GetDiagnosisByID)
	diagnoses.Post("/", diagnosisHandler.CreateDiagnosis)
	diagnoses.Put("/:id", diagnosisHandler.UpdateDiagnosis)
	diagnoses.Delete("/:id", diagnosisHandler.DeleteDiagnosis)

	// Prescription routes (admin only for system management)
	prescriptions := admin.Group("/prescriptions")
	prescriptions.Get("/", prescriptionHandler.GetPrescriptions)
	prescriptions.Get("/:id", prescriptionHandler.GetPrescriptionByID)
	prescriptions.Post("/", prescriptionHandler.CreatePrescription)
	prescriptions.Put("/:id", prescriptionHandler.UpdatePrescription)
	prescriptions.Delete("/:id", prescriptionHandler.DeletePrescription)

	// 404 handler
	app.Use(handlers.NotFound)

	// Start server
	port := ":" + cfg.Port
	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(app.Listen(port))
}
