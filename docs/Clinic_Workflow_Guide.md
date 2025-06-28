# Clinic Workflow Guide

## Overview
This guide outlines the complete clinic workflow for the Rural Health Management System. It demonstrates how a rural health clinic can onboard to the system and manage their daily operations including staff management, patient registration, visit documentation, diagnoses, and prescriptions.

## Prerequisites
- Rural Health Management System API running on `http://localhost:3000`
- Postman installed
- Import the `Clinic_Workflow.postman_collection.json` file

## Workflow Steps

### Phase 1: Clinic Setup and Onboarding

#### Step 1: Health Check
- **Purpose**: Verify the API is running and accessible
- **Method**: GET `/health`
- **Expected Result**: 200 OK with health status
- **What it does**: Basic connectivity test to ensure the system is operational

#### Step 2: Register Clinic
- **Purpose**: Create a new clinic account with authentication
- **Method**: POST `/auth/register/clinic`
- **Data Required**:
  - Email: `clinic@example.com`
  - Password: `securePassword123`
  - Name: `Rural Health Center`
  - Address: `456 Main Street, Central Village`
  - Contact Number: `+1234567891`
  - District: `Central District`
- **Expected Result**: 201 Created with JWT token and clinic profile
- **Automation**: Token is automatically saved for subsequent requests

#### Step 3: Get Clinic Profile
- **Purpose**: Verify clinic registration and authentication
- **Method**: GET `/portal/clinic/profile`
- **Authentication**: Uses saved clinic token
- **Expected Result**: Clinic profile information
- **What it confirms**: Registration was successful and token works

#### Step 4: Check Dashboard Stats
- **Purpose**: View initial dashboard statistics
- **Method**: GET `/portal/clinic/dashboard`
- **Expected Result**: Zero counts for new clinic (patients, staff, visits, etc.)
- **What it shows**: Baseline metrics before adding data

### Phase 2: Staff Management

#### Step 5: Add Doctor Staff
- **Purpose**: Add a doctor to the clinic staff
- **Method**: POST `/portal/clinic/staff`
- **Data Required**:
  - Full Name: `Dr. Sarah Johnson`
  - Role: `Doctor`
  - Phone: `+1234567900`
  - Email: `sarah.johnson@clinic.com`
- **Expected Result**: 201 Created with doctor's staff record
- **Automation**: Doctor ID is saved for creating visits

#### Step 6: Add Nurse Staff
- **Purpose**: Add a nurse to the clinic staff
- **Method**: POST `/portal/clinic/staff`
- **Data Required**:
  - Full Name: `Nurse Mary Wilson`
  - Role: `Nurse`
  - Phone: `+1234567901`
  - Email: `mary.wilson@clinic.com`
- **Expected Result**: 201 Created with nurse's staff record
- **Automation**: Nurse ID is saved for creating visits

#### Step 7: View All Staff
- **Purpose**: Review all staff members added to the clinic
- **Method**: GET `/portal/clinic/staff`
- **Expected Result**: List containing both doctor and nurse
- **What it confirms**: Staff management is working correctly

### Phase 3: Patient Registration

#### Step 8: Register First Patient
- **Purpose**: Register a patient for the clinic
- **Method**: POST `/auth/register/patient`
- **Data Required**:
  - Email: `john.doe@example.com`
  - Password: `patientPassword123`
  - Full Name: `John Doe`
  - Gender: `Male`
  - Date of Birth: `1990-01-15`
  - Address: `123 Village Road, Rural Area`
  - Phone: `+1234567890`
  - Clinic ID: (automatically uses clinic from Step 2)
- **Expected Result**: 201 Created with patient profile and token
- **Automation**: Patient ID and token are saved

#### Step 9: Register Second Patient
- **Purpose**: Register another patient to demonstrate multiple patient management
- **Method**: POST `/auth/register/patient`
- **Data Required**:
  - Email: `jane.smith@example.com`
  - Password: `patientPassword123`
  - Full Name: `Jane Smith`
  - Gender: `Female`
  - Date of Birth: `1985-05-20`
  - Address: `456 Country Lane, Rural Area`
  - Phone: `+1234567892`
  - Clinic ID: (automatically uses clinic from Step 2)
- **Expected Result**: 201 Created with patient profile and token
- **Automation**: Second patient ID and token are saved

#### Step 10: View All Patients
- **Purpose**: Review all patients registered to the clinic
- **Method**: GET `/portal/clinic/patients`
- **Expected Result**: List containing both registered patients
- **What it confirms**: Patient registration and clinic association work correctly

### Phase 4: Visit Management and Medical Care

#### Step 11: Create Visit for Patient 1
- **Purpose**: Document a patient visit with healthcare provider
- **Method**: POST `/portal/clinic/visits`
- **Data Required**:
  - Patient ID: (from Step 8)
  - Staff ID: (doctor from Step 5)
  - Visit Date: `2024-01-15T10:00:00Z`
  - Reason: `Annual checkup`
  - Notes: `Patient appears healthy, routine examination`
- **Expected Result**: 201 Created with visit record
- **Automation**: Visit ID is saved for adding diagnosis and prescription

#### Step 12: Add Diagnosis to Visit 1
- **Purpose**: Add medical diagnosis to the visit
- **Method**: POST `/portal/clinic/diagnoses`
- **Data Required**:
  - Visit ID: (from Step 11)
  - Diagnosis Code: `Z00.00`
  - Description: `General health examination - patient in good health`
- **Expected Result**: 201 Created with diagnosis record
- **What it demonstrates**: ICD-10 coding and diagnosis documentation

#### Step 13: Add Prescription to Visit 1
- **Purpose**: Add medication prescription to the visit
- **Method**: POST `/portal/clinic/prescriptions`
- **Data Required**:
  - Visit ID: (from Step 11)
  - Medication Name: `Multivitamin`
  - Dosage: `1 tablet`
  - Instructions: `Take one tablet daily with breakfast`
  - Duration: `30 days`
- **Expected Result**: 201 Created with prescription record
- **What it demonstrates**: Prescription management and medication tracking

#### Step 14: Create Visit for Patient 2
- **Purpose**: Document a second patient visit (different scenario)
- **Method**: POST `/portal/clinic/visits`
- **Data Required**:
  - Patient ID: (from Step 9)
  - Staff ID: (nurse from Step 6)
  - Visit Date: `2024-01-16T14:30:00Z`
  - Reason: `Flu symptoms consultation`
  - Notes: `Patient reports fever, cough, and fatigue for 2 days`
- **Expected Result**: 201 Created with visit record
- **What it demonstrates**: Different staff member handling visits

#### Step 15: Add Diagnosis to Visit 2
- **Purpose**: Add flu diagnosis to the second visit
- **Method**: POST `/portal/clinic/diagnoses`
- **Data Required**:
  - Visit ID: (from Step 14)
  - Diagnosis Code: `J11.1`
  - Description: `Influenza with respiratory manifestations`
- **Expected Result**: 201 Created with diagnosis record
- **What it demonstrates**: Illness diagnosis documentation

#### Step 16: Add Prescription to Visit 2
- **Purpose**: Add fever medication prescription
- **Method**: POST `/portal/clinic/prescriptions`
- **Data Required**:
  - Visit ID: (from Step 14)
  - Medication Name: `Acetaminophen`
  - Dosage: `500mg`
  - Instructions: `Take every 6 hours as needed for fever`
  - Duration: `7 days`
- **Expected Result**: 201 Created with prescription record
- **What it demonstrates**: Treatment prescription for acute illness

### Phase 5: Review and Monitoring

#### Step 17: View All Visits
- **Purpose**: Review all visits documented in the clinic
- **Method**: GET `/portal/clinic/visits`
- **Expected Result**: List of all visits with patient, staff, diagnosis, and prescription details
- **What it shows**: Complete visit history with related medical information

#### Step 18: View Specific Patient Details
- **Purpose**: Get comprehensive patient information including visit history
- **Method**: GET `/portal/clinic/patients/{patient1Id}`
- **Expected Result**: Patient profile with complete visit history
- **What it demonstrates**: Patient-centric view of medical history

#### Step 19: Final Dashboard Check
- **Purpose**: Review updated dashboard statistics after completing workflow
- **Method**: GET `/portal/clinic/dashboard`
- **Expected Result**: Updated counts (2 patients, 2 staff, 2 visits, 2 diagnoses, 2 prescriptions)
- **What it confirms**: All data has been properly recorded and is reflected in dashboard metrics

## Key Features Demonstrated

### 1. Role-Based Access Control
- Clinic authentication with JWT tokens
- Clinic-specific data access (only see own patients, staff, visits)
- Secure API endpoints requiring proper authentication

### 2. Comprehensive Healthcare Management
- Staff management (multiple roles: Doctor, Nurse)
- Patient registration and profile management
- Visit documentation with date, time, reason, and notes
- Medical diagnosis tracking with ICD-10 codes
- Prescription management with detailed medication information

### 3. Data Relationships
- Clinics have multiple patients and staff
- Visits link patients with staff members
- Diagnoses and prescriptions are tied to specific visits
- Complete audit trail of medical interactions

### 4. Automation and Workflow
- Automatic token management in Postman
- Variable passing between requests
- Step-by-step progression that builds complete medical records

## Expected Outcomes

After completing this workflow, you will have:

1. **Functional Clinic**: A registered clinic with authentication
2. **Medical Staff**: Doctor and nurse staff members
3. **Patient Base**: Two registered patients with different demographics
4. **Medical Records**: Complete visit records with diagnoses and prescriptions
5. **Dashboard Metrics**: Updated statistics showing clinic activity

## Error Handling

Common issues and solutions:

- **401 Unauthorized**: Check if API is running and tokens are valid
- **400 Bad Request**: Verify all required fields are provided with correct data types
- **404 Not Found**: Ensure IDs from previous steps are correctly saved in variables
- **500 Internal Server Error**: Check API logs for database or server issues

## Next Steps

After completing this workflow:
1. Run the Patient Workflow to see the patient perspective
2. Experiment with different data scenarios
3. Test error conditions (invalid data, unauthorized access)
4. Explore the admin endpoints for system-wide management
