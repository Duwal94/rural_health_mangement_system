package main

import (
	"log"
	"time"

	"rural_health_management_system/internal/config"
	"rural_health_management_system/internal/database"
	"rural_health_management_system/internal/models"
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

	log.Println("Starting to seed database...")

	// Create sample clinics
	clinics := []models.Clinic{
		{
			Name:          "Central Rural Health Center",
			Address:       "123 Main Street, Central Village",
			ContactNumber: "+1234567890",
			District:      "Central District",
		},
		{
			Name:          "North Village Clinic",
			Address:       "456 North Road, North Village",
			ContactNumber: "+1234567891",
			District:      "Northern District",
		},
		{
			Name:          "South Community Health",
			Address:       "789 South Avenue, South Village",
			ContactNumber: "+1234567892",
			District:      "Southern District",
		},
	}

	for _, clinic := range clinics {
		if err := db.DB.Create(&clinic).Error; err != nil {
			log.Printf("Error creating clinic %s: %v", clinic.Name, err)
		} else {
			log.Printf("Created clinic: %s", clinic.Name)
		}
	}

	// Create sample staff
	staff := []models.Staff{
		{
			FullName: "Dr. Sarah Johnson",
			Role:     "Doctor",
			Phone:    "+1234567900",
			Email:    "sarah.johnson@clinic.com",
			ClinicID: 1,
		},
		{
			FullName: "Nurse Mary Wilson",
			Role:     "Nurse",
			Phone:    "+1234567901",
			Email:    "mary.wilson@clinic.com",
			ClinicID: 1,
		},
		{
			FullName: "Dr. Michael Brown",
			Role:     "Doctor",
			Phone:    "+1234567902",
			Email:    "michael.brown@clinic.com",
			ClinicID: 2,
		},
		{
			FullName: "Admin Lisa Davis",
			Role:     "Administrator",
			Phone:    "+1234567903",
			Email:    "lisa.davis@clinic.com",
			ClinicID: 2,
		},
		{
			FullName: "Pharmacist John Smith",
			Role:     "Pharmacist",
			Phone:    "+1234567904",
			Email:    "john.smith@clinic.com",
			ClinicID: 3,
		},
	}

	for _, s := range staff {
		if err := db.DB.Create(&s).Error; err != nil {
			log.Printf("Error creating staff %s: %v", s.FullName, err)
		} else {
			log.Printf("Created staff: %s", s.FullName)
		}
	}

	// Create sample patients
	patients := []models.Patient{
		{
			FullName:    "Alice Cooper",
			Gender:      "Female",
			DateOfBirth: time.Date(1985, 3, 15, 0, 0, 0, 0, time.UTC),
			Address:     "101 Village Lane, Central Village",
			Phone:       "+1234567910",
			ClinicID:    1,
		},
		{
			FullName:    "Bob Williams",
			Gender:      "Male",
			DateOfBirth: time.Date(1990, 7, 22, 0, 0, 0, 0, time.UTC),
			Address:     "202 Country Road, Central Village",
			Phone:       "+1234567911",
			ClinicID:    1,
		},
		{
			FullName:    "Carol Martinez",
			Gender:      "Female",
			DateOfBirth: time.Date(1978, 12, 5, 0, 0, 0, 0, time.UTC),
			Address:     "303 Farm Street, North Village",
			Phone:       "+1234567912",
			ClinicID:    2,
		},
		{
			FullName:    "David Lee",
			Gender:      "Male",
			DateOfBirth: time.Date(1982, 9, 18, 0, 0, 0, 0, time.UTC),
			Address:     "404 Rural Route, South Village",
			Phone:       "+1234567913",
			ClinicID:    3,
		},
		{
			FullName:    "Emma Thompson",
			Gender:      "Female",
			DateOfBirth: time.Date(1995, 1, 30, 0, 0, 0, 0, time.UTC),
			Address:     "505 Meadow Drive, South Village",
			Phone:       "+1234567914",
			ClinicID:    3,
		},
	}

	for _, patient := range patients {
		if err := db.DB.Create(&patient).Error; err != nil {
			log.Printf("Error creating patient %s: %v", patient.FullName, err)
		} else {
			log.Printf("Created patient: %s", patient.FullName)
		}
	}

	// Create sample visits
	visits := []models.Visit{
		{
			PatientID: 1,
			ClinicID:  1,
			StaffID:   1,
			VisitDate: time.Now().AddDate(0, 0, -7),
			Reason:    "Annual checkup",
			Notes:     "Patient appears healthy, all vitals normal",
		},
		{
			PatientID: 2,
			ClinicID:  1,
			StaffID:   2,
			VisitDate: time.Now().AddDate(0, 0, -5),
			Reason:    "Cold symptoms",
			Notes:     "Mild fever, recommended rest and fluids",
		},
		{
			PatientID: 3,
			ClinicID:  2,
			StaffID:   3,
			VisitDate: time.Now().AddDate(0, 0, -3),
			Reason:    "Hypertension follow-up",
			Notes:     "Blood pressure improving with medication",
		},
		{
			PatientID: 4,
			ClinicID:  3,
			StaffID:   5,
			VisitDate: time.Now().AddDate(0, 0, -1),
			Reason:    "Medication review",
			Notes:     "Adjusting diabetes medication dosage",
		},
	}

	for i, visit := range visits {
		if err := db.DB.Create(&visit).Error; err != nil {
			log.Printf("Error creating visit %d: %v", i+1, err)
		} else {
			log.Printf("Created visit for patient ID %d", visit.PatientID)
		}
	}

	// Create sample diagnoses
	diagnoses := []models.Diagnosis{
		{
			VisitID:       1,
			DiagnosisCode: "Z00.00",
			Description:   "General health examination",
		},
		{
			VisitID:       2,
			DiagnosisCode: "J00",
			Description:   "Acute upper respiratory infection",
		},
		{
			VisitID:       3,
			DiagnosisCode: "I10",
			Description:   "Essential hypertension",
		},
		{
			VisitID:       4,
			DiagnosisCode: "E11.9",
			Description:   "Type 2 diabetes mellitus without complications",
		},
	}

	for _, diagnosis := range diagnoses {
		if err := db.DB.Create(&diagnosis).Error; err != nil {
			log.Printf("Error creating diagnosis %s: %v", diagnosis.DiagnosisCode, err)
		} else {
			log.Printf("Created diagnosis: %s", diagnosis.DiagnosisCode)
		}
	}

	// Create sample prescriptions
	prescriptions := []models.Prescription{
		{
			VisitID:        2,
			MedicationName: "Acetaminophen",
			Dosage:         "500mg",
			Instructions:   "Take every 6 hours as needed for fever",
			DurationDays:   7,
		},
		{
			VisitID:        3,
			MedicationName: "Lisinopril",
			Dosage:         "10mg",
			Instructions:   "Take once daily in the morning",
			DurationDays:   30,
		},
		{
			VisitID:        4,
			MedicationName: "Metformin",
			Dosage:         "500mg",
			Instructions:   "Take twice daily with meals",
			DurationDays:   30,
		},
	}

	for _, prescription := range prescriptions {
		if err := db.DB.Create(&prescription).Error; err != nil {
			log.Printf("Error creating prescription %s: %v", prescription.MedicationName, err)
		} else {
			log.Printf("Created prescription: %s", prescription.MedicationName)
		}
	}

	log.Println("Database seeding completed successfully!")
}
