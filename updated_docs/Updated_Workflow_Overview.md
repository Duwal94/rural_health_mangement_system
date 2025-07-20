# Updated Workflow Overview - RBAC System

## Introduction

This document provides a comprehensive overview of the updated Rural Health Management System workflows, incorporating the new Role-Based Access Control (RBAC) system. The system now features separate login portals and granular permissions for different user types.

## System Architecture

### User Types and Access Levels

The system now supports five distinct user types with specific access levels:

1. **Clinic Staff** (`clinic_staff`) - Administrative control
2. **Doctors** (`doctor`) - Medical practice with diagnostic capabilities
3. **Nurses** (`nurse`) - Medical support without diagnostic capabilities
4. **Patients** (`patient`) - Personal health data access
5. **System Admin** (`admin`) - System-wide management

### Portal Structure

```
├── Patient Portal (/portal/patient/*)
│   └── Personal health data access only
├── Staff Portal (/portal/staff/*)
│   └── Clinic administrative functions
├── Medical Portal (/portal/medical/*)
│   └── Medical practice functions
└── Admin Portal (/*)
    └── System-wide management
```

## Workflow 1: Clinic Setup and Staff Management

### Step 1: Clinic Registration

**Endpoint**: `POST /api/v1/auth/register/clinic`

**Purpose**: Creates a new clinic and the primary clinic staff account.

**Request Body**:
```json
{
  "email": "admin@ruralclinic.com",
  "password": "securepass123",
  "name": "Rural Health Clinic",
  "address": "123 Village Road, Remote Area",
  "contact_number": "+1234567890",
  "district": "Mountain District"
}
```

**Response**: 
- Creates a `clinic_staff` user account
- Links user to the clinic
- Returns JWT token for immediate access

### Step 2: Clinic Staff Login

**Endpoint**: `POST /api/v1/auth/clinic-login`

**Purpose**: Authenticates clinic staff for administrative portal access.

**Request Body**:
```json
{
  "email": "admin@ruralclinic.com",
  "password": "securepass123",
  "login_type": "staff"
}
```

**Key Features**:
- `login_type: "staff"` ensures access to staff portal
- Validates user type matches login type
- Returns JWT with clinic context

### Step 3: Create Medical Staff (Doctors)

**Endpoint**: `POST /api/v1/portal/staff/staff`

**Authorization**: Clinic Staff only

**Purpose**: Creates doctor accounts with login capabilities.

**Request Body**:
```json
{
  "email": "doctor@ruralclinic.com",
  "password": "doctorpass123",
  "full_name": "Dr. Sarah Johnson",
  "role": "Doctor",
  "phone": "+1234567891"
}
```

**Process**:
1. Creates `doctor` user account
2. Creates staff record linked to clinic
3. Enables medical portal access

### Step 4: Create Medical Staff (Nurses)

**Endpoint**: `POST /api/v1/portal/staff/staff`

**Authorization**: Clinic Staff only

**Purpose**: Creates nurse accounts with login capabilities.

**Request Body**:
```json
{
  "email": "nurse@ruralclinic.com",
  "password": "nursepass123",
  "full_name": "Nurse Mary Wilson",
  "role": "Nurse",
  "phone": "+1234567892"
}
```

### Step 5: Create Non-Login Staff

**Endpoint**: `POST /api/v1/portal/staff/staff`

**Authorization**: Clinic Staff only

**Purpose**: Creates administrative support staff without login capabilities.

**Request Body**:
```json
{
  "email": "admin.support@ruralclinic.com",
  "password": "temppass123",
  "full_name": "Admin Assistant Tom",
  "role": "Clinic_Administrator",
  "phone": "+1234567893"
}
```

**Note**: `Clinic_Administrator` and `Pharmacist` roles create staff records but no user accounts.

## Workflow 2: Patient Management

### Step 1: Patient Registration by Clinic Staff

**Endpoint**: `POST /api/v1/portal/staff/patients`

**Authorization**: Clinic Staff only

**Purpose**: Registers new patients in the clinic system.

**Request Body**:
```json
{
  "full_name": "John Doe",
  "gender": "Male",
  "date_of_birth": "1985-05-15",
  "address": "456 Rural Street, Village",
  "phone": "+1234567894"
}
```

**Process**:
1. Validates clinic staff permissions
2. Creates patient record linked to clinic
3. No user account created (patients register separately if needed)

### Step 2: Patient Self-Registration (Optional)

**Endpoint**: `POST /api/v1/auth/register/patient`

**Purpose**: Allows patients to create their own portal access.

**Request Body**:
```json
{
  "email": "patient@example.com",
  "password": "patientpass123",
  "full_name": "Jane Patient",
  "gender": "Female",
  "date_of_birth": "1990-03-20",
  "address": "789 Patient Lane",
  "phone": "+1234567896",
  "clinic_id": 1
}
```

### Step 3: Patient Login

**Endpoint**: `POST /api/v1/auth/login`

**Purpose**: Authenticates patients for portal access.

**Request Body**:
```json
{
  "email": "patient@example.com",
  "password": "patientpass123"
}
```

## Workflow 3: Medical Practice (Doctor)

### Step 1: Doctor Login

**Endpoint**: `POST /api/v1/auth/clinic-login`

**Purpose**: Authenticates doctors for medical portal access.

**Request Body**:
```json
{
  "email": "doctor@ruralclinic.com",
  "password": "doctorpass123",
  "login_type": "medical"
}
```

**Key Features**:
- `login_type: "medical"` ensures access to medical portal
- Provides read-only access to patients and staff
- Enables medical record creation

### Step 2: View Patients

**Endpoint**: `GET /api/v1/portal/medical/patients`

**Authorization**: Doctor/Nurse

**Purpose**: Lists all patients in the clinic (read-only).

**Query Parameters**:
- `page`: Page number (default: 1)
- `per_page`: Items per page (default: 10, max: 100)
- `search`: Search term for patient name or phone

### Step 3: Create Patient Visit

**Endpoint**: `POST /api/v1/portal/medical/visits`

**Authorization**: Doctor/Nurse

**Purpose**: Records a patient visit.

**Request Body**:
```json
{
  "patient_id": 1,
  "reason": "Regular checkup and flu symptoms",
  "notes": "Patient reports fever and cough for 3 days"
}
```

**Automatic Fields**:
- `clinic_id`: Set from user's clinic context
- `staff_id`: Set from user's staff ID
- `visit_date`: Set to current timestamp if not provided

### Step 4: Create Diagnosis (Doctors Only)

**Endpoint**: `POST /api/v1/portal/medical/diagnoses`

**Authorization**: Doctors only

**Purpose**: Records medical diagnosis for a visit.

**Request Body**:
```json
{
  "visit_id": 1,
  "diagnosis_code": "J11.1",
  "description": "Influenza due to unidentified influenza virus with other respiratory manifestations"
}
```

**Validation**:
- Only doctors can create diagnoses
- Nurses attempting this receive 403 Forbidden
- Visit must belong to doctor's clinic

### Step 5: Create Prescription (Doctors Only)

**Endpoint**: `POST /api/v1/portal/medical/prescriptions`

**Authorization**: Doctors only

**Purpose**: Prescribes medication for a patient visit.

**Request Body**:
```json
{
  "visit_id": 1,
  "medication_name": "Tamiflu",
  "dosage": "75mg",
  "instructions": "Take twice daily with food",
  "duration_days": 5
}
```

## Workflow 4: Medical Support (Nurse)

### Step 1: Nurse Login

**Endpoint**: `POST /api/v1/auth/clinic-login`

**Purpose**: Authenticates nurses for medical portal access.

**Request Body**:
```json
{
  "email": "nurse@ruralclinic.com",
  "password": "nursepass123",
  "login_type": "medical"
}
```

### Step 2: Support Patient Care

**Available Actions**:
- ✅ View patients (read-only)
- ✅ View staff (read-only)
- ✅ Create and manage visits
- ✅ View existing diagnoses and prescriptions
- ❌ Create diagnoses (doctor-only)
- ❌ Create prescriptions (doctor-only)

**Key Endpoints**:
```bash
GET /api/v1/portal/medical/patients          # View patients
GET /api/v1/portal/medical/patients/{id}     # View patient details
POST /api/v1/portal/medical/visits           # Create visit
GET /api/v1/portal/medical/visits            # View visits
GET /api/v1/portal/medical/dashboard         # View medical dashboard
```

## Workflow 5: Patient Portal Access

### Step 1: Patient Login

**Endpoint**: `POST /api/v1/auth/login`

**Purpose**: Standard login for patient portal access.

### Step 2: View Personal Health Data

**Available Endpoints**:
```bash
GET /api/v1/portal/patient/profile           # View own profile
PUT /api/v1/portal/patient/profile           # Update own profile
GET /api/v1/portal/patient/visits            # View own visits
GET /api/v1/portal/patient/visits/{id}       # View specific visit
GET /api/v1/portal/patient/diagnoses         # View own diagnoses
GET /api/v1/portal/patient/prescriptions     # View own prescriptions
```

**Security**: Patients can only access their own data, enforced by user ID validation.

## Access Control Matrix

| Action | Clinic Staff | Doctor | Nurse | Patient |
|--------|--------------|--------|-------|---------|
| Create Patients | ✅ | ❌ | ❌ | Self Only |
| View Patients | ✅ All | ✅ Clinic Only | ✅ Clinic Only | ✅ Self Only |
| Manage Staff | ✅ | ❌ | ❌ | ❌ |
| Create Visits | ✅ | ✅ | ✅ | ❌ |
| Create Diagnosis | ❌ | ✅ | ❌ | ❌ |
| Create Prescription | ❌ | ✅ | ❌ | ❌ |
| View Medical Records | ✅ Clinic | ✅ Clinic | ✅ Clinic | ✅ Self Only |
| Manage Clinic | ✅ | ❌ | ❌ | ❌ |

## Error Handling

### Authentication Errors

**401 Unauthorized**:
- Missing or invalid JWT token
- Token expired
- Invalid credentials

**Example Response**:
```json
{
  "error": "Invalid token"
}
```

### Authorization Errors

**403 Forbidden**:
- Insufficient permissions for action
- Wrong portal access attempt
- Cross-clinic access attempt

**Example Responses**:
```json
{
  "error": "Access denied. Only doctors can perform this action"
}
```

```json
{
  "error": "Insufficient permissions"
}
```

### Validation Errors

**400 Bad Request**:
- Invalid request body
- Missing required fields
- Invalid data format

**Example Response**:
```json
{
  "error": "Invalid request body",
  "details": {
    "full_name": "required field missing"
  }
}
```

## Security Features

### 1. Multi-Level Validation
- JWT token validation
- User type validation
- Permission-level validation
- Clinic ownership validation

### 2. Portal Separation
- Staff portal: Administrative functions
- Medical portal: Medical functions
- Patient portal: Personal data access

### 3. Cross-Clinic Protection
- Users can only access their clinic's data
- Automatic clinic ID injection in requests
- Validation at database query level

### 4. Role-Specific Restrictions
- Doctors: Cannot manage staff or modify clinic
- Nurses: Cannot create medical records
- Staff: Cannot create medical records
- Patients: Cannot access administrative functions

## Best Practices

### 1. Token Management
- Store JWT tokens securely
- Handle token expiration gracefully
- Use appropriate token for each portal

### 2. Error Handling
- Implement proper error handling for all scenarios
- Display user-friendly error messages
- Log security violations for audit

### 3. Data Validation
- Validate all input data
- Use proper data types and formats
- Implement client-side and server-side validation

### 4. Permission Checks
- Always check permissions before actions
- Use appropriate middleware for route protection
- Implement fallback security measures

This workflow system ensures secure, role-appropriate access to clinic management functions while maintaining clear separation of responsibilities between administrative and medical staff.
