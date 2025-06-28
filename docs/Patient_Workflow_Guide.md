# Patient Workflow Guide

## Overview
This guide outlines the complete patient workflow for the Rural Health Management System. It demonstrates how patients can register for healthcare services, manage their profiles, and access their medical information including visits, diagnoses, and prescriptions.

## Prerequisites
- Rural Health Management System API running on `http://localhost:3000`
- At least one clinic registered in the system (use Clinic Workflow first)
- Postman installed
- Import the `Patient_Workflow.postman_collection.json` file

## Workflow Steps

### Phase 1: Discovery and Registration

#### Step 1: Health Check
- **Purpose**: Verify the API is running and accessible
- **Method**: GET `/health`
- **Expected Result**: 200 OK with health status
- **What it does**: Ensures the system is operational before patient registration

#### Step 2: Check Available Clinics
- **Purpose**: View available clinics for registration
- **Method**: GET `/clinics?page=1&per_page=10`
- **Authentication**: None required (public information)
- **Expected Result**: List of available clinics with locations and contact info
- **What patients see**: Available healthcare facilities in their area

#### Step 3: Register as Patient
- **Purpose**: Create a patient account with clinic association
- **Method**: POST `/auth/register/patient`
- **Data Required**:
  - Email: `patient.demo@example.com`
  - Password: `patientPassword123`
  - Full Name: `Alice Johnson`
  - Gender: `Female`
  - Date of Birth: `1992-03-10`
  - Address: `789 Rural Street, Countryside`
  - Phone: `+1234567893`
  - Clinic ID: `1` (associate with existing clinic)
- **Expected Result**: 201 Created with JWT token and patient profile
- **Automation**: Patient token, ID, and clinic ID are automatically saved

### Phase 2: Profile Management

#### Step 4: Get My Profile
- **Purpose**: View personal profile and clinic association
- **Method**: GET `/portal/patient/profile`
- **Authentication**: Uses saved patient token
- **Expected Result**: Complete patient profile with associated clinic information
- **What it shows**: Personal details and which clinic provides care

#### Step 5: Update My Profile
- **Purpose**: Modify personal information
- **Method**: PUT `/portal/patient/profile`
- **Data Updated**:
  - Full Name: `Alice Johnson (Updated)`
  - Address: `789 Updated Rural Street, Countryside`
  - Phone: `+1234567893`
- **Expected Result**: 200 OK with updated profile
- **What it demonstrates**: Patients can maintain their own information
- **Note**: Patients cannot change their associated clinic

### Phase 3: Initial Medical Record Check

#### Step 6: Check Initial Visits
- **Purpose**: View visit history (should be empty for new patient)
- **Method**: GET `/portal/patient/visits?page=1&per_page=10`
- **Expected Result**: Empty list or 200 OK with no data
- **What it establishes**: Baseline for medical history tracking

#### Step 7: Check Initial Diagnoses
- **Purpose**: View diagnosis history (should be empty for new patient)
- **Method**: GET `/portal/patient/diagnoses?page=1&per_page=10`
- **Expected Result**: Empty list or 200 OK with no data
- **What it establishes**: No previous medical diagnoses on record

#### Step 8: Check Initial Prescriptions
- **Purpose**: View prescription history (should be empty for new patient)
- **Method**: GET `/portal/patient/prescriptions?page=1&per_page=10`
- **Expected Result**: Empty list or 200 OK with no data
- **What it establishes**: No current or past medications on record

### Phase 4: Account Security

#### Step 9: Login Again (Simulate Return Visit)
- **Purpose**: Demonstrate login functionality for returning patients
- **Method**: POST `/auth/login`
- **Data Required**:
  - Email: `patient.demo@example.com`
  - Password: `patientPassword123`
- **Expected Result**: 200 OK with fresh JWT token
- **Automation**: New token is automatically saved
- **What it simulates**: Patient returning to access their records

#### Step 10: Get Profile via Auth
- **Purpose**: Alternative method to access profile information
- **Method**: GET `/auth/profile`
- **Authentication**: Uses current patient token
- **Expected Result**: Patient profile (same as Step 4 but different endpoint)
- **What it demonstrates**: Multiple ways to access profile data

#### Step 11: Change Password
- **Purpose**: Update account password for security
- **Method**: POST `/auth/change-password`
- **Data Required**:
  - Current Password: `patientPassword123`
  - New Password: `newPatientPassword123`
- **Expected Result**: 200 OK with success message
- **What it demonstrates**: Account security management

#### Step 12: Login with New Password
- **Purpose**: Verify password change was successful
- **Method**: POST `/auth/login`
- **Data Required**:
  - Email: `patient.demo@example.com`
  - Password: `newPatientPassword123` (updated password)
- **Expected Result**: 200 OK with JWT token
- **Automation**: Updated token is saved
- **What it confirms**: Password change functionality works correctly

### Phase 5: Accessing Medical Records

#### Step 13: View My Visits After Medical Care
- **Purpose**: Check visits after clinic has created medical records
- **Method**: GET `/portal/patient/visits?page=1&per_page=10`
- **Expected Result**: List of visits if clinic workflow has been run
- **What patients see**: Complete visit history with clinic staff, dates, and reasons
- **Note**: Run clinic workflow first to create sample data

#### Step 14: View Specific Visit Details
- **Purpose**: Get detailed information about a particular visit
- **Method**: GET `/portal/patient/visits/{visitId}`
- **Path Parameter**: Visit ID (update manually in Postman variable)
- **Expected Result**: Detailed visit information including diagnoses and prescriptions
- **What it provides**: Comprehensive view of a single medical encounter

#### Step 15: View My Diagnoses
- **Purpose**: Access all medical diagnoses from clinic visits
- **Method**: GET `/portal/patient/diagnoses?page=1&per_page=10`
- **Expected Result**: List of diagnoses with ICD codes and descriptions
- **What patients see**: 
  - Diagnosis codes (e.g., Z00.00, J11.1)
  - Medical descriptions
  - Associated visit information
  - Attending healthcare provider

#### Step 16: View My Prescriptions
- **Purpose**: Access all medication prescriptions
- **Method**: GET `/portal/patient/prescriptions?page=1&per_page=10`
- **Expected Result**: List of all prescribed medications
- **What patients see**:
  - Medication names
  - Dosage information
  - Instructions for use
  - Duration of treatment
  - Prescribing visit details

### Phase 6: Data Exploration and Navigation

#### Step 17: Search My Prescriptions
- **Purpose**: Demonstrate pagination and filtering capabilities
- **Method**: GET `/portal/patient/prescriptions?page=1&per_page=5`
- **Expected Result**: Smaller page size showing pagination works
- **What it demonstrates**: Patients can navigate through large amounts of medical data

#### Step 18: View Recent Visits
- **Purpose**: Focus on most recent medical encounters
- **Method**: GET `/portal/patient/visits?page=1&per_page=5`
- **Expected Result**: Recent visits with reduced page size
- **What it provides**: Quick access to latest medical care

#### Step 19: Final Profile Check
- **Purpose**: Confirm all profile updates are maintained
- **Method**: GET `/portal/patient/profile`
- **Expected Result**: Profile with all updates from Step 5
- **What it confirms**: Data persistence and profile management work correctly

## Key Features Demonstrated

### 1. Patient Self-Service
- Independent registration and account creation
- Profile management and updates
- Password security management
- Access to complete medical history

### 2. Medical Record Access
- View all visits with healthcare providers
- Access diagnosis information with medical codes
- Review prescription history with detailed medication information
- See associated clinic and staff member details

### 3. Data Privacy and Security
- JWT-based authentication
- Patient can only access their own medical records
- Secure password change functionality
- Token-based session management

### 4. User Experience Features
- Pagination for large datasets
- Detailed visit information including related diagnoses and prescriptions
- Multiple endpoints for accessing the same information
- Clear data relationships (visits → diagnoses/prescriptions)

### 5. Healthcare Integration
- Association with specific clinic
- Complete medical history tracking
- Prescription and medication management
- Healthcare provider identification

## Data Relationships Shown

### Patient → Clinic
- Each patient is associated with one clinic
- Patients can see their clinic's information
- Clinic association cannot be changed by patient

### Patient → Visits
- Patients can view all their visits
- Each visit shows the attending staff member
- Visit details include date, time, reason, and clinical notes

### Visits → Diagnoses
- Each visit can have multiple diagnoses
- Diagnoses include ICD codes and descriptions
- Patients see complete diagnostic history

### Visits → Prescriptions
- Each visit can have multiple prescriptions
- Prescriptions include medication details, dosage, and instructions
- Complete medication history is accessible

## Expected Outcomes

After completing this workflow, the patient will have:

1. **Active Account**: Registered patient profile with authentication
2. **Updated Profile**: Modified personal information
3. **Security Setup**: Changed password and verified login
4. **Medical Record Access**: Ability to view visits, diagnoses, and prescriptions
5. **System Familiarity**: Understanding of navigation and data access

## Data Dependencies

To see medical records (Steps 13-18), you need:
1. **Existing Clinic**: Run the Clinic Workflow first
2. **Visit Records**: Clinic must create visits for this patient
3. **Medical Data**: Diagnoses and prescriptions added by clinic staff

## Integration Testing

This workflow is designed to work with:
- **Clinic Workflow**: Provides the clinic and medical data
- **Admin Workflow**: System-wide management and oversight
- **Real Clinic Usage**: Actual patient registration at existing clinics

## Error Scenarios and Testing

### Common Issues:
- **Empty Medical Records**: Normal for new patients before clinic visits
- **401 Unauthorized**: Token expiration or invalid credentials
- **404 Not Found**: Incorrect visit IDs or non-existent records
- **403 Forbidden**: Attempting to access other patients' data

### Security Testing:
- Try accessing another patient's visit ID (should fail)
- Use expired or invalid tokens (should fail)
- Attempt to change clinic association (should fail)

## Real-World Usage

This workflow simulates:
1. **Patient Registration**: New patient signing up for healthcare services
2. **Profile Management**: Keeping personal information current
3. **Return Visits**: Existing patients accessing their records
4. **Medical History Review**: Patients checking their health records
5. **Prescription Tracking**: Monitoring current and past medications

## Next Steps

After completing this workflow:
1. Run the Clinic Workflow to create medical data for this patient
2. Test the Admin Workflow for system management perspective
3. Experiment with multiple patients and clinic scenarios
4. Test error conditions and security boundaries
