package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"rural_health_management_system/internal/config"
	"rural_health_management_system/internal/database"
	"rural_health_management_system/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Realistic medical data
var diagnosisCodes = []struct {
	Code        string
	Description string
	Seasonal    bool
	Chronic     bool
}{
	{"Z00.00", "General health examination", false, false},
	{"J11.1", "Influenza with other respiratory manifestations", true, false},
	{"J44.1", "Chronic obstructive pulmonary disease with acute exacerbation", false, true},
	{"I10", "Essential hypertension", false, true},
	{"E11.9", "Type 2 diabetes mellitus without complications", false, true},
	{"K59.00", "Constipation, unspecified", false, false},
	{"J06.9", "Acute upper respiratory infection, unspecified", true, false},
	{"M79.3", "Panniculitis, unspecified", false, false},
	{"R50.9", "Fever, unspecified", false, false},
	{"A09", "Infectious gastroenteritis and colitis", false, false},
	{"J02.9", "Acute pharyngitis, unspecified", true, false},
	{"M25.50", "Pain in unspecified joint", false, false},
	{"R06.02", "Shortness of breath", false, false},
	{"K21.9", "Gastro-esophageal reflux disease without esophagitis", false, true},
	{"L20.9", "Atopic dermatitis, unspecified", false, true},
	{"H10.9", "Unspecified conjunctivitis", false, false},
	{"N39.0", "Urinary tract infection, site not specified", false, false},
	{"J45.9", "Asthma, unspecified", false, true},
	{"F32.9", "Major depressive disorder, single episode, unspecified", false, true},
	{"M54.5", "Low back pain", false, false},
	{"R51", "Headache", false, false},
	{"B34.9", "Viral infection, unspecified", true, false},
	{"J03.90", "Acute tonsillitis, unspecified", true, false},
	{"L30.9", "Dermatitis, unspecified", false, false},
	{"R05", "Cough", true, false},
}

var medications = []struct {
	Name         string
	CommonDosage string
	Duration     []int
}{
	{"Paracetamol", "500mg", []int{3, 5, 7}},
	{"Ibuprofen", "400mg", []int{3, 5, 7}},
	{"Amoxicillin", "500mg", []int{7, 10, 14}},
	{"Metformin", "500mg", []int{30, 60, 90}},
	{"Lisinopril", "10mg", []int{30, 60, 90}},
	{"Amlodipine", "5mg", []int{30, 60, 90}},
	{"Omeprazole", "20mg", []int{14, 30, 60}},
	{"Salbutamol", "100mcg", []int{30, 60, 90}},
	{"Prednisolone", "5mg", []int{5, 7, 14}},
	{"Cetirizine", "10mg", []int{7, 14, 30}},
	{"Azithromycin", "250mg", []int{3, 5}},
	{"Diclofenac", "50mg", []int{5, 7, 14}},
	{"Loratadine", "10mg", []int{7, 14, 30}},
	{"Furosemide", "40mg", []int{30, 60, 90}},
	{"Simvastatin", "20mg", []int{30, 60, 90}},
	{"Aspirin", "75mg", []int{30, 60, 90}},
	{"Insulin", "10 units", []int{30, 60, 90}},
	{"Levothyroxine", "50mcg", []int{30, 60, 90}},
	{"Warfarin", "5mg", []int{30, 60, 90}},
	{"Tramadol", "50mg", []int{3, 5, 7}},
}

var districts = []string{
	"Kathmandu", "Lalitpur", "Bhaktapur", "Kavrepalanchok", "Sindhupalchok",
	"Nuwakot", "Rasuwa", "Dhading", "Makwanpur", "Sindhuli",
	"Ramechhap", "Dolakha", "Solukhumbu", "Okhaldhunga", "Khotang",
	"Udayapur", "Saptari", "Siraha", "Dhanusa", "Mahottari",
}

var clinicNames = []string{
	"Community Health Center", "Primary Health Care", "Rural Health Post",
	"District Hospital", "Medical Clinic", "Health Care Center",
	"Family Medicine Clinic", "General Hospital", "Regional Health Center",
	"Public Health Center",
}

var nepaliFemaleNames = []string{
	"Sita Sharma", "Gita Poudel", "Maya Tamang", "Kamala Thapa", "Laxmi Gurung",
	"Parvati Shrestha", "Radha Rai", "Sushila Magar", "Bindu Limbu", "Sunita Chhetri",
	"Mira Karki", "Devi Pandey", "Sarita Dhital", "Kiran Bajracharya", "Nirmala Ghimire",
	"Anita Joshi", "Rina Bhattarai", "Sangita Maharjan", "Kumari Dahal", "Pushpa Adhikari",
}

var nepaliMaleNames = []string{
	"Ram Bahadur Thapa", "Shyam Prasad Sharma", "Hari Krishna Poudel", "Krishna Tamang", "Gopal Gurung",
	"Bishal Shrestha", "Prakash Rai", "Dipak Magar", "Sanjay Limbu", "Rajesh Chhetri",
	"Nabin Karki", "Suresh Pandey", "Ramesh Dhital", "Arjun Bajracharya", "Mahesh Ghimire",
	"Deepak Joshi", "Bikash Bhattarai", "Kamal Maharjan", "Umesh Dahal", "Santosh Adhikari",
}

var reasons = []string{
	"Routine checkup", "Fever and headache", "Cough and cold", "Stomach pain",
	"Back pain", "Chest pain", "Joint pain", "Skin rash", "Breathing difficulty",
	"Diabetes follow-up", "Hypertension monitoring", "Vaccination", "Wound dressing",
	"Health screening", "Medication refill", "Pregnancy checkup", "Eye examination",
	"Dental consultation", "Mental health consultation", "Injury treatment",
}

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Seed the database
	log.Println("Starting database seeding...")

	if err := seedDatabase(db.DB); err != nil {
		log.Fatal("Failed to seed database:", err)
	}

	log.Println("Database seeding completed successfully!")
}

func seedDatabase(db *gorm.DB) error {
	rand.Seed(time.Now().UnixNano())

	log.Println("Creating admin user...")
	if err := createAdminUser(db); err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	log.Println("Creating clinics...")
	clinics, err := createClinics(db, 15) // Create 15 clinics
	if err != nil {
		return fmt.Errorf("failed to create clinics: %w", err)
	}

	log.Println("Creating staff...")
	allStaff, err := createStaff(db, clinics)
	if err != nil {
		return fmt.Errorf("failed to create staff: %w", err)
	}

	log.Println("Creating patients...")
	allPatients, err := createPatients(db, clinics, 200) // Create 200 patients total
	if err != nil {
		return fmt.Errorf("failed to create patients: %w", err)
	}

	log.Println("Creating visits...")
	visits, err := createVisits(db, clinics, allStaff, allPatients, 1500) // Create 1500 visits
	if err != nil {
		return fmt.Errorf("failed to create visits: %w", err)
	}

	log.Println("Creating diagnoses...")
	if err := createDiagnoses(db, visits); err != nil {
		return fmt.Errorf("failed to create diagnoses: %w", err)
	}

	log.Println("Creating prescriptions...")
	if err := createPrescriptions(db, visits); err != nil {
		return fmt.Errorf("failed to create prescriptions: %w", err)
	}

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func createAdminUser(db *gorm.DB) error {
	hashedPassword, err := hashPassword("admin123")
	if err != nil {
		return err
	}

	admin := models.User{
		Email:    "admin@ruralhealth.com",
		Password: hashedPassword,
		UserType: "admin",
		IsActive: true,
	}

	return db.FirstOrCreate(&admin, models.User{Email: admin.Email}).Error
}

func createClinics(db *gorm.DB, count int) ([]models.Clinic, error) {
	var clinics []models.Clinic

	for i := 0; i < count; i++ {
		// Create user for clinic
		hashedPassword, err := hashPassword("clinic123")
		if err != nil {
			return nil, err
		}

		user := models.User{
			Email:    fmt.Sprintf("clinic%d@ruralhealth.com", i+1),
			Password: hashedPassword,
			UserType: "clinic_staff",
			IsActive: true,
		}

		if err := db.Create(&user).Error; err != nil {
			return nil, err
		}

		// Create clinic
		district := districts[rand.Intn(len(districts))]
		clinicName := fmt.Sprintf("%s %s", district, clinicNames[rand.Intn(len(clinicNames))])

		clinic := models.Clinic{
			Name:          clinicName,
			Address:       fmt.Sprintf("%s District, Ward No. %d", district, rand.Intn(15)+1),
			ContactNumber: fmt.Sprintf("01-%d", 4000000+rand.Intn(999999)),
			District:      district,
			UserID:        &user.ID,
		}

		if err := db.Create(&clinic).Error; err != nil {
			return nil, err
		}

		clinics = append(clinics, clinic)
		log.Printf("Created clinic: %s", clinic.Name)
	}

	return clinics, nil
}

func createStaff(db *gorm.DB, clinics []models.Clinic) ([]models.Staff, error) {
	var allStaff []models.Staff
	roles := []string{"Doctor", "Nurse", "Clinic_Administrator", "Pharmacist"}

	for _, clinic := range clinics {
		// Create 3-5 staff members per clinic
		staffCount := rand.Intn(3) + 3
		for i := 0; i < staffCount; i++ {
			role := roles[rand.Intn(len(roles))]
			var name string
			if rand.Float32() < 0.6 { // 60% chance of female
				name = nepaliFemaleNames[rand.Intn(len(nepaliFemaleNames))]
			} else {
				name = nepaliMaleNames[rand.Intn(len(nepaliMaleNames))]
			}

			// Create user for staff
			hashedPassword, err := hashPassword("staff123")
			if err != nil {
				return nil, err
			}

			var userType string
			switch role {
			case "Doctor":
				userType = "doctor"
			case "Nurse":
				userType = "nurse"
			default:
				userType = "clinic_staff"
			}

			user := models.User{
				Email:    fmt.Sprintf("staff%d_%d@ruralhealth.com", clinic.ID, i+1),
				Password: hashedPassword,
				UserType: userType,
				IsActive: true,
			}

			if err := db.Create(&user).Error; err != nil {
				return nil, err
			}

			staff := models.Staff{
				FullName: name,
				Role:     role,
				Phone:    fmt.Sprintf("98%08d", rand.Intn(99999999)),
				Email:    user.Email,
				ClinicID: clinic.ID,
				UserID:   &user.ID,
				IsActive: true,
			}

			if err := db.Create(&staff).Error; err != nil {
				return nil, err
			}

			allStaff = append(allStaff, staff)
		}
	}

	log.Printf("Created %d staff members", len(allStaff))
	return allStaff, nil
}

func createPatients(db *gorm.DB, clinics []models.Clinic, totalCount int) ([]models.Patient, error) {
	var allPatients []models.Patient

	for i := 0; i < totalCount; i++ {
		clinic := clinics[rand.Intn(len(clinics))]

		// Generate realistic age distribution
		age := generateRealisticAge()
		dateOfBirth := time.Now().AddDate(-age, -rand.Intn(12), -rand.Intn(28))

		var name, gender string
		if rand.Float32() < 0.55 { // 55% chance of female (realistic for healthcare)
			name = nepaliFemaleNames[rand.Intn(len(nepaliFemaleNames))]
			gender = "Female"
		} else {
			name = nepaliMaleNames[rand.Intn(len(nepaliMaleNames))]
			gender = "Male"
		}

		// Create user for some patients (not all patients need user accounts)
		var userID *uint
		if rand.Float32() < 0.3 { // 30% of patients have user accounts
			hashedPassword, err := hashPassword("patient123")
			if err != nil {
				return nil, err
			}

			user := models.User{
				Email:    fmt.Sprintf("patient%d@ruralhealth.com", i+1),
				Password: hashedPassword,
				UserType: "patient",
				IsActive: true,
			}

			if err := db.Create(&user).Error; err != nil {
				return nil, err
			}
			userID = &user.ID
		}

		patient := models.Patient{
			FullName:    name,
			Gender:      gender,
			DateOfBirth: dateOfBirth,
			Address:     fmt.Sprintf("%s District, Ward No. %d", clinic.District, rand.Intn(15)+1),
			Phone:       fmt.Sprintf("98%08d", rand.Intn(99999999)),
			ClinicID:    clinic.ID,
			UserID:      userID,
		}

		if err := db.Create(&patient).Error; err != nil {
			return nil, err
		}

		allPatients = append(allPatients, patient)
	}

	log.Printf("Created %d patients", len(allPatients))
	return allPatients, nil
}

func generateRealisticAge() int {
	// Generate realistic age distribution for healthcare
	r := rand.Float32()
	switch {
	case r < 0.15: // 15% children (0-17)
		return rand.Intn(18)
	case r < 0.35: // 20% young adults (18-30)
		return rand.Intn(13) + 18
	case r < 0.55: // 20% adults (31-50)
		return rand.Intn(20) + 31
	case r < 0.80: // 25% middle-aged (51-70)
		return rand.Intn(20) + 51
	default: // 20% elderly (71+)
		return rand.Intn(20) + 71
	}
}

func createVisits(db *gorm.DB, clinics []models.Clinic, allStaff []models.Staff, allPatients []models.Patient, totalCount int) ([]models.Visit, error) {
	var visits []models.Visit

	// Create a map of clinic ID to staff for quick lookup
	clinicStaff := make(map[uint][]models.Staff)
	for _, staff := range allStaff {
		clinicStaff[staff.ClinicID] = append(clinicStaff[staff.ClinicID], staff)
	}

	for i := 0; i < totalCount; i++ {
		patient := allPatients[rand.Intn(len(allPatients))]
		clinic := clinics[0] // Find the patient's clinic
		for _, c := range clinics {
			if c.ID == patient.ClinicID {
				clinic = c
				break
			}
		}

		// Get staff from the same clinic
		staffList := clinicStaff[clinic.ID]
		if len(staffList) == 0 {
			continue // Skip if no staff in this clinic
		}
		staff := staffList[rand.Intn(len(staffList))]

		// Generate visit date within the last 12 months
		daysAgo := rand.Intn(365)
		visitDate := time.Now().AddDate(0, 0, -daysAgo)

		// Add some seasonal bias for certain conditions
		month := visitDate.Month()
		isWinterMonth := month == 12 || month == 1 || month == 2
		isSummerMonth := month >= 6 && month <= 8

		var reason string
		if isWinterMonth && rand.Float32() < 0.4 {
			winterReasons := []string{"Cough and cold", "Fever and headache", "Breathing difficulty"}
			reason = winterReasons[rand.Intn(len(winterReasons))]
		} else if isSummerMonth && rand.Float32() < 0.3 {
			summerReasons := []string{"Stomach pain", "Skin rash", "Vaccination"}
			reason = summerReasons[rand.Intn(len(summerReasons))]
		} else {
			reason = reasons[rand.Intn(len(reasons))]
		}

		visit := models.Visit{
			PatientID: patient.ID,
			ClinicID:  clinic.ID,
			StaffID:   staff.ID,
			VisitDate: visitDate,
			Reason:    reason,
			Notes:     fmt.Sprintf("Patient examined by %s. General condition stable.", staff.FullName),
			CreatedAt: visitDate,
			UpdatedAt: visitDate,
		}

		if err := db.Create(&visit).Error; err != nil {
			return nil, err
		}

		visits = append(visits, visit)
	}

	log.Printf("Created %d visits", len(visits))
	return visits, nil
}

func createDiagnoses(db *gorm.DB, visits []models.Visit) error {
	diagnosisCount := 0
	for _, visit := range visits {
		// 85% of visits get a diagnosis
		if rand.Float32() < 0.85 {
			// Some visits get multiple diagnoses (15% chance)
			diagnosesPerVisit := 1
			if rand.Float32() < 0.15 {
				diagnosesPerVisit = 2
			}

			for i := 0; i < diagnosesPerVisit; i++ {
				// Bias diagnosis selection based on season and patient age
				diagnosis := selectDiagnosis(visit.VisitDate, visit.PatientID)

				diagnosisRecord := models.Diagnosis{
					VisitID:       visit.ID,
					DiagnosisCode: diagnosis.Code,
					Description:   diagnosis.Description,
					CreatedAt:     visit.VisitDate,
					UpdatedAt:     visit.VisitDate,
				}

				if err := db.Create(&diagnosisRecord).Error; err != nil {
					return err
				}
				diagnosisCount++
			}
		}
	}

	log.Printf("Created %d diagnoses", diagnosisCount)
	return nil
}

func selectDiagnosis(visitDate time.Time, patientID uint) struct {
	Code        string
	Description string
	Seasonal    bool
	Chronic     bool
} {
	month := visitDate.Month()
	isWinterMonth := month == 12 || month == 1 || month == 2

	// Winter months favor seasonal illnesses
	if isWinterMonth && rand.Float32() < 0.6 {
		seasonalDiagnoses := []struct {
			Code        string
			Description string
			Seasonal    bool
			Chronic     bool
		}{
			{"J11.1", "Influenza with other respiratory manifestations", true, false},
			{"J06.9", "Acute upper respiratory infection, unspecified", true, false},
			{"J02.9", "Acute pharyngitis, unspecified", true, false},
			{"J03.90", "Acute tonsillitis, unspecified", true, false},
			{"R05", "Cough", true, false},
			{"B34.9", "Viral infection, unspecified", true, false},
		}
		return seasonalDiagnoses[rand.Intn(len(seasonalDiagnoses))]
	}

	// Otherwise select from all diagnoses with some bias toward common conditions
	if rand.Float32() < 0.4 {
		// 40% chance of common conditions
		commonDiagnoses := []struct {
			Code        string
			Description string
			Seasonal    bool
			Chronic     bool
		}{
			{"Z00.00", "General health examination", false, false},
			{"I10", "Essential hypertension", false, true},
			{"E11.9", "Type 2 diabetes mellitus without complications", false, true},
			{"M54.5", "Low back pain", false, false},
			{"R51", "Headache", false, false},
		}
		return commonDiagnoses[rand.Intn(len(commonDiagnoses))]
	}

	return diagnosisCodes[rand.Intn(len(diagnosisCodes))]
}

func createPrescriptions(db *gorm.DB, visits []models.Visit) error {
	prescriptionCount := 0
	for _, visit := range visits {
		// 75% of visits get prescriptions
		if rand.Float32() < 0.75 {
			// Number of medications (1-3)
			medicationCount := rand.Intn(3) + 1

			for i := 0; i < medicationCount; i++ {
				medication := medications[rand.Intn(len(medications))]
				duration := medication.Duration[rand.Intn(len(medication.Duration))]

				// Generate realistic instructions
				instructions := generateInstructions(medication.Name, medication.CommonDosage)

				prescription := models.Prescription{
					VisitID:        visit.ID,
					MedicationName: medication.Name,
					Dosage:         medication.CommonDosage,
					Instructions:   instructions,
					DurationDays:   duration,
					CreatedAt:      visit.VisitDate,
					UpdatedAt:      visit.VisitDate,
				}

				if err := db.Create(&prescription).Error; err != nil {
					return err
				}
				prescriptionCount++
			}
		}
	}

	log.Printf("Created %d prescriptions", prescriptionCount)
	return nil
}

func generateInstructions(medication, dosage string) string {
	instructions := []string{
		fmt.Sprintf("Take %s twice daily after meals", dosage),
		fmt.Sprintf("Take %s three times daily before meals", dosage),
		fmt.Sprintf("Take %s once daily in the morning", dosage),
		fmt.Sprintf("Take %s as needed for pain", dosage),
		fmt.Sprintf("Take %s with plenty of water", dosage),
		fmt.Sprintf("Take %s on empty stomach", dosage),
		fmt.Sprintf("Apply %s to affected area twice daily", dosage),
		fmt.Sprintf("Take %s at bedtime", dosage),
	}

	// Add specific instructions for certain medications
	switch medication {
	case "Insulin":
		return "Inject subcutaneously 30 minutes before meals"
	case "Salbutamol":
		return "Inhale 2 puffs as needed for breathing difficulty"
	case "Prednisolone":
		return "Take with food to avoid stomach upset. Do not stop suddenly."
	default:
		return instructions[rand.Intn(len(instructions))]
	}
}
