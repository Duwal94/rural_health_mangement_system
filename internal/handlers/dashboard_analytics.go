package handlers

import (
	"rural_health_management_system/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type DashboardAnalyticsHandler struct {
	db *gorm.DB
}

func NewDashboardAnalyticsHandler(db *gorm.DB) *DashboardAnalyticsHandler {
	return &DashboardAnalyticsHandler{db: db}
}

// DiagnosisAnalytics represents diagnosis statistics
type DiagnosisAnalytics struct {
	DiagnosisCode string  `json:"diagnosis_code"`
	Description   string  `json:"description"`
	Count         int64   `json:"count"`
	Percentage    float64 `json:"percentage"`
}

// PrescriptionAnalytics represents prescription statistics
type PrescriptionAnalytics struct {
	MedicationName  string       `json:"medication_name"`
	Count           int64        `json:"count"`
	Percentage      float64      `json:"percentage"`
	AvgDurationDays float64      `json:"avg_duration_days"`
	CommonDosages   []DosageInfo `json:"common_dosages"`
}

// DosageInfo represents dosage patterns
type DosageInfo struct {
	Dosage string `json:"dosage"`
	Count  int64  `json:"count"`
}

// DemographicAnalytics represents patient demographics
type DemographicAnalytics struct {
	AgeGroups          []AgeGroupInfo `json:"age_groups"`
	GenderDistribution []GenderInfo   `json:"gender_distribution"`
}

// AgeGroupInfo represents age group distribution
type AgeGroupInfo struct {
	AgeGroup   string  `json:"age_group"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

// GenderInfo represents gender distribution
type GenderInfo struct {
	Gender     string  `json:"gender"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

// IllnessTrend represents illness trends over time
type IllnessTrend struct {
	Month         string `json:"month"`
	Year          int    `json:"year"`
	DiagnosisCode string `json:"diagnosis_code"`
	Count         int64  `json:"count"`
}

// DistrictAnalytics represents district-level analytics
type DistrictAnalytics struct {
	District         string                  `json:"district"`
	TotalClinics     int64                   `json:"total_clinics"`
	TotalPatients    int64                   `json:"total_patients"`
	TotalVisits      int64                   `json:"total_visits"`
	TopDiagnoses     []DiagnosisAnalytics    `json:"top_diagnoses"`
	TopPrescriptions []PrescriptionAnalytics `json:"top_prescriptions"`
}

// ComprehensiveDashboard represents the complete dashboard analytics
type ComprehensiveDashboard struct {
	OverallStats      OverallStats            `json:"overall_stats"`
	TopDiagnoses      []DiagnosisAnalytics    `json:"top_diagnoses"`
	TopPrescriptions  []PrescriptionAnalytics `json:"top_prescriptions"`
	Demographics      DemographicAnalytics    `json:"demographics"`
	IllnessTrends     []IllnessTrend          `json:"illness_trends"`
	DistrictAnalytics []DistrictAnalytics     `json:"district_analytics"`
	SeasonalTrends    []SeasonalTrend         `json:"seasonal_trends"`
}

// OverallStats represents overall system statistics
type OverallStats struct {
	TotalClinics       int64 `json:"total_clinics"`
	TotalPatients      int64 `json:"total_patients"`
	TotalStaff         int64 `json:"total_staff"`
	TotalVisits        int64 `json:"total_visits"`
	TotalDiagnoses     int64 `json:"total_diagnoses"`
	TotalPrescriptions int64 `json:"total_prescriptions"`
	VisitsThisMonth    int64 `json:"visits_this_month"`
	VisitsToday        int64 `json:"visits_today"`
}

// SeasonalTrend represents seasonal illness patterns
type SeasonalTrend struct {
	Season        string `json:"season"`
	Month         int    `json:"month"`
	Year          int    `json:"year"`
	DiagnosisCode string `json:"diagnosis_code"`
	Count         int64  `json:"count"`
}

// GetSystemDashboard returns comprehensive analytics for all clinics (public)
func (h *DashboardAnalyticsHandler) GetSystemDashboard(c *fiber.Ctx) error {
	dashboard := ComprehensiveDashboard{}

	// Get overall stats
	dashboard.OverallStats = h.getOverallStats()

	// Get top diagnoses across all clinics
	dashboard.TopDiagnoses = h.getTopDiagnoses(nil, 10)

	// Get top prescriptions across all clinics
	dashboard.TopPrescriptions = h.getTopPrescriptions(nil, 10)

	// Get demographics across all clinics
	dashboard.Demographics = h.getDemographics(nil)

	// Get illness trends (last 12 months)
	dashboard.IllnessTrends = h.getIllnessTrends(nil, 12)

	// Get district analytics
	dashboard.DistrictAnalytics = h.getDistrictAnalytics()

	// Get seasonal trends
	dashboard.SeasonalTrends = h.getSeasonalTrends(nil)

	return c.JSON(dashboard)
}

// GetClinicDashboard returns comprehensive analytics for a specific clinic
func (h *DashboardAnalyticsHandler) GetClinicDashboard(c *fiber.Ctx) error {
	clinicID := c.Locals("clinic_id").(uint)

	dashboard := ComprehensiveDashboard{}

	// Get overall stats for this clinic
	dashboard.OverallStats = h.getClinicStats(clinicID)

	// Get top diagnoses for this clinic
	dashboard.TopDiagnoses = h.getTopDiagnoses(&clinicID, 10)

	// Get top prescriptions for this clinic
	dashboard.TopPrescriptions = h.getTopPrescriptions(&clinicID, 10)

	// Get demographics for this clinic
	dashboard.Demographics = h.getDemographics(&clinicID)

	// Get illness trends for this clinic (last 12 months)
	dashboard.IllnessTrends = h.getIllnessTrends(&clinicID, 12)

	// Get seasonal trends for this clinic
	dashboard.SeasonalTrends = h.getSeasonalTrends(&clinicID)

	return c.JSON(dashboard)
}

// Helper functions

func (h *DashboardAnalyticsHandler) getOverallStats() OverallStats {
	var stats OverallStats

	h.db.Model(&models.Clinic{}).Count(&stats.TotalClinics)
	h.db.Model(&models.Patient{}).Count(&stats.TotalPatients)
	h.db.Model(&models.Staff{}).Where("is_active = ?", true).Count(&stats.TotalStaff)
	h.db.Model(&models.Visit{}).Count(&stats.TotalVisits)
	h.db.Model(&models.Diagnosis{}).Count(&stats.TotalDiagnoses)
	h.db.Model(&models.Prescription{}).Count(&stats.TotalPrescriptions)

	// Visits this month
	now := time.Now()
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	h.db.Model(&models.Visit{}).Where("visit_date >= ?", firstDayOfMonth).Count(&stats.VisitsThisMonth)

	// Visits today
	today := now.Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	h.db.Model(&models.Visit{}).Where("visit_date >= ? AND visit_date < ?", today, tomorrow).Count(&stats.VisitsToday)

	return stats
}

func (h *DashboardAnalyticsHandler) getClinicStats(clinicID uint) OverallStats {
	var stats OverallStats

	stats.TotalClinics = 1 // This clinic
	h.db.Model(&models.Patient{}).Where("clinic_id = ?", clinicID).Count(&stats.TotalPatients)
	h.db.Model(&models.Staff{}).Where("clinic_id = ? AND is_active = ?", clinicID, true).Count(&stats.TotalStaff)
	h.db.Model(&models.Visit{}).Where("clinic_id = ?", clinicID).Count(&stats.TotalVisits)

	// Count diagnoses and prescriptions for this clinic's visits
	h.db.Table("diagnoses").
		Joins("JOIN visits ON diagnoses.visit_id = visits.id").
		Where("visits.clinic_id = ?", clinicID).
		Count(&stats.TotalDiagnoses)

	h.db.Table("prescriptions").
		Joins("JOIN visits ON prescriptions.visit_id = visits.id").
		Where("visits.clinic_id = ?", clinicID).
		Count(&stats.TotalPrescriptions)

	// Visits this month for this clinic
	now := time.Now()
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	h.db.Model(&models.Visit{}).Where("clinic_id = ? AND visit_date >= ?", clinicID, firstDayOfMonth).Count(&stats.VisitsThisMonth)

	// Visits today for this clinic
	today := now.Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	h.db.Model(&models.Visit{}).Where("clinic_id = ? AND visit_date >= ? AND visit_date < ?", clinicID, today, tomorrow).Count(&stats.VisitsToday)

	return stats
}

func (h *DashboardAnalyticsHandler) getTopDiagnoses(clinicID *uint, limit int) []DiagnosisAnalytics {
	var results []DiagnosisAnalytics

	query := h.db.Table("diagnoses").
		Select("diagnosis_code, description, COUNT(*) as count").
		Group("diagnosis_code, description").
		Order("count DESC").
		Limit(limit)

	if clinicID != nil {
		query = query.Joins("JOIN visits ON diagnoses.visit_id = visits.id").
			Where("visits.clinic_id = ?", *clinicID)
	}

	var rawResults []struct {
		DiagnosisCode string `json:"diagnosis_code"`
		Description   string `json:"description"`
		Count         int64  `json:"count"`
	}

	query.Scan(&rawResults)

	// Calculate total for percentages
	var total int64
	for _, result := range rawResults {
		total += result.Count
	}

	// Convert to response format with percentages
	for _, result := range rawResults {
		percentage := float64(0)
		if total > 0 {
			percentage = float64(result.Count) / float64(total) * 100
		}

		results = append(results, DiagnosisAnalytics{
			DiagnosisCode: result.DiagnosisCode,
			Description:   result.Description,
			Count:         result.Count,
			Percentage:    percentage,
		})
	}

	return results
}

func (h *DashboardAnalyticsHandler) getTopPrescriptions(clinicID *uint, limit int) []PrescriptionAnalytics {
	var results []PrescriptionAnalytics

	query := h.db.Table("prescriptions").
		Select("medication_name, COUNT(*) as count, AVG(duration_days) as avg_duration_days").
		Group("medication_name").
		Order("count DESC").
		Limit(limit)

	if clinicID != nil {
		query = query.Joins("JOIN visits ON prescriptions.visit_id = visits.id").
			Where("visits.clinic_id = ?", *clinicID)
	}

	var rawResults []struct {
		MedicationName  string  `json:"medication_name"`
		Count           int64   `json:"count"`
		AvgDurationDays float64 `json:"avg_duration_days"`
	}

	query.Scan(&rawResults)

	// Calculate total for percentages
	var total int64
	for _, result := range rawResults {
		total += result.Count
	}

	// Convert to response format with percentages and dosage info
	for _, result := range rawResults {
		percentage := float64(0)
		if total > 0 {
			percentage = float64(result.Count) / float64(total) * 100
		}

		// Get common dosages for this medication
		dosages := h.getCommonDosages(result.MedicationName, clinicID)

		results = append(results, PrescriptionAnalytics{
			MedicationName:  result.MedicationName,
			Count:           result.Count,
			Percentage:      percentage,
			AvgDurationDays: result.AvgDurationDays,
			CommonDosages:   dosages,
		})
	}

	return results
}

func (h *DashboardAnalyticsHandler) getCommonDosages(medicationName string, clinicID *uint) []DosageInfo {
	var dosages []DosageInfo

	query := h.db.Table("prescriptions").
		Select("dosage, COUNT(*) as count").
		Where("medication_name = ?", medicationName).
		Group("dosage").
		Order("count DESC").
		Limit(5)

	if clinicID != nil {
		query = query.Joins("JOIN visits ON prescriptions.visit_id = visits.id").
			Where("visits.clinic_id = ?", *clinicID)
	}

	query.Scan(&dosages)
	return dosages
}

func (h *DashboardAnalyticsHandler) getDemographics(clinicID *uint) DemographicAnalytics {
	var demographics DemographicAnalytics

	// Age groups
	demographics.AgeGroups = h.getAgeGroups(clinicID)

	// Gender distribution
	demographics.GenderDistribution = h.getGenderDistribution(clinicID)

	return demographics
}

func (h *DashboardAnalyticsHandler) getAgeGroups(clinicID *uint) []AgeGroupInfo {
	var results []AgeGroupInfo

	// Calculate age groups based on date of birth
	query := h.db.Table("patients").
		Select(`
			CASE 
				WHEN DATE_PART('year', AGE(CURRENT_DATE, date_of_birth)) < 18 THEN 'Under 18'
				WHEN DATE_PART('year', AGE(CURRENT_DATE, date_of_birth)) BETWEEN 18 AND 30 THEN '18-30'
				WHEN DATE_PART('year', AGE(CURRENT_DATE, date_of_birth)) BETWEEN 31 AND 50 THEN '31-50'
				WHEN DATE_PART('year', AGE(CURRENT_DATE, date_of_birth)) BETWEEN 51 AND 70 THEN '51-70'
				ELSE 'Over 70'
			END as age_group,
			COUNT(*) as count
		`).
		Group("age_group").
		Order("age_group")

	if clinicID != nil {
		query = query.Where("clinic_id = ?", *clinicID)
	}

	var rawResults []struct {
		AgeGroup string `json:"age_group"`
		Count    int64  `json:"count"`
	}

	query.Scan(&rawResults)

	// Calculate total for percentages
	var total int64
	for _, result := range rawResults {
		total += result.Count
	}

	// Convert to response format with percentages
	for _, result := range rawResults {
		percentage := float64(0)
		if total > 0 {
			percentage = float64(result.Count) / float64(total) * 100
		}

		results = append(results, AgeGroupInfo{
			AgeGroup:   result.AgeGroup,
			Count:      result.Count,
			Percentage: percentage,
		})
	}

	return results
}

func (h *DashboardAnalyticsHandler) getGenderDistribution(clinicID *uint) []GenderInfo {
	var results []GenderInfo

	query := h.db.Table("patients").
		Select("gender, COUNT(*) as count").
		Group("gender").
		Order("gender")

	if clinicID != nil {
		query = query.Where("clinic_id = ?", *clinicID)
	}

	var rawResults []struct {
		Gender string `json:"gender"`
		Count  int64  `json:"count"`
	}

	query.Scan(&rawResults)

	// Calculate total for percentages
	var total int64
	for _, result := range rawResults {
		total += result.Count
	}

	// Convert to response format with percentages
	for _, result := range rawResults {
		percentage := float64(0)
		if total > 0 {
			percentage = float64(result.Count) / float64(total) * 100
		}

		results = append(results, GenderInfo{
			Gender:     result.Gender,
			Count:      result.Count,
			Percentage: percentage,
		})
	}

	return results
}

func (h *DashboardAnalyticsHandler) getIllnessTrends(clinicID *uint, months int) []IllnessTrend {
	var trends []IllnessTrend

	// Get illness trends for the last N months
	startDate := time.Now().AddDate(0, -months, 0)

	query := h.db.Table("diagnoses").
		Select(`
			TO_CHAR(visits.visit_date, 'Month') as month,
			DATE_PART('year', visits.visit_date) as year,
			diagnosis_code,
			COUNT(*) as count
		`).
		Joins("JOIN visits ON diagnoses.visit_id = visits.id").
		Where("visits.visit_date >= ?", startDate).
		Group("month, year, diagnosis_code").
		Order("year, DATE_PART('month', visits.visit_date), count DESC")

	if clinicID != nil {
		query = query.Where("visits.clinic_id = ?", *clinicID)
	}

	query.Scan(&trends)
	return trends
}

func (h *DashboardAnalyticsHandler) getSeasonalTrends(clinicID *uint) []SeasonalTrend {
	var trends []SeasonalTrend

	// Get seasonal patterns for the last 2 years
	startDate := time.Now().AddDate(-2, 0, 0)

	query := h.db.Table("diagnoses").
		Select(`
			CASE 
				WHEN DATE_PART('month', visits.visit_date) IN (12, 1, 2) THEN 'Winter'
				WHEN DATE_PART('month', visits.visit_date) IN (3, 4, 5) THEN 'Spring'
				WHEN DATE_PART('month', visits.visit_date) IN (6, 7, 8) THEN 'Summer'
				ELSE 'Fall'
			END as season,
			DATE_PART('month', visits.visit_date) as month,
			DATE_PART('year', visits.visit_date) as year,
			diagnosis_code,
			COUNT(*) as count
		`).
		Joins("JOIN visits ON diagnoses.visit_id = visits.id").
		Where("visits.visit_date >= ?", startDate).
		Group("season, month, year, diagnosis_code").
		Order("year, month, count DESC")

	if clinicID != nil {
		query = query.Where("visits.clinic_id = ?", *clinicID)
	}

	query.Scan(&trends)
	return trends
}

func (h *DashboardAnalyticsHandler) getDistrictAnalytics() []DistrictAnalytics {
	var districts []DistrictAnalytics

	// Get basic district stats
	var districtStats []struct {
		District      string `json:"district"`
		TotalClinics  int64  `json:"total_clinics"`
		TotalPatients int64  `json:"total_patients"`
		TotalVisits   int64  `json:"total_visits"`
	}

	h.db.Table("clinics").
		Select(`
			district,
			COUNT(*) as total_clinics,
			SUM((SELECT COUNT(*) FROM patients WHERE patients.clinic_id = clinics.id)) as total_patients,
			SUM((SELECT COUNT(*) FROM visits WHERE visits.clinic_id = clinics.id)) as total_visits
		`).
		Group("district").
		Order("total_visits DESC").
		Scan(&districtStats)

	// For each district, get top diagnoses and prescriptions
	for _, stat := range districtStats {
		district := DistrictAnalytics{
			District:      stat.District,
			TotalClinics:  stat.TotalClinics,
			TotalPatients: stat.TotalPatients,
			TotalVisits:   stat.TotalVisits,
		}

		// Get top diagnoses for this district
		district.TopDiagnoses = h.getTopDiagnosesForDistrict(stat.District, 5)

		// Get top prescriptions for this district
		district.TopPrescriptions = h.getTopPrescriptionsForDistrict(stat.District, 5)

		districts = append(districts, district)
	}

	return districts
}

func (h *DashboardAnalyticsHandler) getTopDiagnosesForDistrict(district string, limit int) []DiagnosisAnalytics {
	var results []DiagnosisAnalytics

	var rawResults []struct {
		DiagnosisCode string `json:"diagnosis_code"`
		Description   string `json:"description"`
		Count         int64  `json:"count"`
	}

	h.db.Table("diagnoses").
		Select("diagnosis_code, description, COUNT(*) as count").
		Joins("JOIN visits ON diagnoses.visit_id = visits.id").
		Joins("JOIN clinics ON visits.clinic_id = clinics.id").
		Where("clinics.district = ?", district).
		Group("diagnosis_code, description").
		Order("count DESC").
		Limit(limit).
		Scan(&rawResults)

	// Calculate total for percentages
	var total int64
	for _, result := range rawResults {
		total += result.Count
	}

	// Convert to response format with percentages
	for _, result := range rawResults {
		percentage := float64(0)
		if total > 0 {
			percentage = float64(result.Count) / float64(total) * 100
		}

		results = append(results, DiagnosisAnalytics{
			DiagnosisCode: result.DiagnosisCode,
			Description:   result.Description,
			Count:         result.Count,
			Percentage:    percentage,
		})
	}

	return results
}

func (h *DashboardAnalyticsHandler) getTopPrescriptionsForDistrict(district string, limit int) []PrescriptionAnalytics {
	var results []PrescriptionAnalytics

	var rawResults []struct {
		MedicationName  string  `json:"medication_name"`
		Count           int64   `json:"count"`
		AvgDurationDays float64 `json:"avg_duration_days"`
	}

	h.db.Table("prescriptions").
		Select("medication_name, COUNT(*) as count, AVG(duration_days) as avg_duration_days").
		Joins("JOIN visits ON prescriptions.visit_id = visits.id").
		Joins("JOIN clinics ON visits.clinic_id = clinics.id").
		Where("clinics.district = ?", district).
		Group("medication_name").
		Order("count DESC").
		Limit(limit).
		Scan(&rawResults)

	// Calculate total for percentages
	var total int64
	for _, result := range rawResults {
		total += result.Count
	}

	// Convert to response format with percentages
	for _, result := range rawResults {
		percentage := float64(0)
		if total > 0 {
			percentage = float64(result.Count) / float64(total) * 100
		}

		results = append(results, PrescriptionAnalytics{
			MedicationName:  result.MedicationName,
			Count:           result.Count,
			Percentage:      percentage,
			AvgDurationDays: result.AvgDurationDays,
		})
	}

	return results
}
