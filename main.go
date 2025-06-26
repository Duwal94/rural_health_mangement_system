package main

import (
	"fmt"
	"log"

	"rural_health_management_system/internal/config"
	"rural_health_management_system/internal/database"
	"rural_health_management_system/internal/handlers"

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
	clinicHandler := handlers.NewClinicHandler(db.DB)
	patientHandler := handlers.NewPatientHandler(db.DB)
	staffHandler := handlers.NewStaffHandler(db.DB)
	visitHandler := handlers.NewVisitHandler(db.DB)
	diagnosisHandler := handlers.NewDiagnosisHandler(db.DB)
	prescriptionHandler := handlers.NewPrescriptionHandler(db.DB)

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

	// Clinic routes
	clinics := v1.Group("/clinics")
	clinics.Get("/", clinicHandler.GetClinics)
	clinics.Get("/:id", clinicHandler.GetClinic)
	clinics.Post("/", clinicHandler.CreateClinic)
	clinics.Put("/:id", clinicHandler.UpdateClinic)
	clinics.Delete("/:id", clinicHandler.DeleteClinic)

	// Patient routes
	patients := v1.Group("/patients")
	patients.Get("/", patientHandler.GetPatients)
	patients.Get("/:id", patientHandler.GetPatient)
	patients.Post("/", patientHandler.CreatePatient)
	patients.Put("/:id", patientHandler.UpdatePatient)
	patients.Delete("/:id", patientHandler.DeletePatient)

	// Staff routes
	staff := v1.Group("/staff")
	staff.Get("/", staffHandler.GetStaff)
	staff.Get("/:id", staffHandler.GetStaffByID)
	staff.Post("/", staffHandler.CreateStaff)
	staff.Put("/:id", staffHandler.UpdateStaff)
	staff.Delete("/:id", staffHandler.DeleteStaff) // Visit routes
	visits := v1.Group("/visits")
	visits.Get("/", visitHandler.GetVisits)
	visits.Get("/:id", visitHandler.GetVisitByID)
	visits.Post("/", visitHandler.CreateVisit)
	visits.Put("/:id", visitHandler.UpdateVisit)
	visits.Delete("/:id", visitHandler.DeleteVisit)

	// Diagnosis routes
	diagnoses := v1.Group("/diagnoses")
	diagnoses.Get("/", diagnosisHandler.GetDiagnoses)
	diagnoses.Get("/:id", diagnosisHandler.GetDiagnosisByID)
	diagnoses.Post("/", diagnosisHandler.CreateDiagnosis)
	diagnoses.Put("/:id", diagnosisHandler.UpdateDiagnosis)
	diagnoses.Delete("/:id", diagnosisHandler.DeleteDiagnosis)

	// Prescription routes
	prescriptions := v1.Group("/prescriptions")
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
