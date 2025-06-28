# Rural Health Management System API Documentation

## Base URL
```
http://localhost:3000/api/v1
```

## Overview
The Rural Health Management System API provides comprehensive endpoints for managing rural healthcare facilities, including clinics, patients, staff members, visits, diagnoses, and prescriptions. The system features role-based authentication with separate patient and clinic portals, ensuring data privacy and access control.

## Authentication System
The API uses JWT (JSON Web Token) based authentication with three user types:
- **Patient**: Can only access their own medical data
- **Clinic**: Can access their own clinic data, patients, staff, and visits
- **Admin**: Can access all system data (for system management)

### User Types and Access Levels
- **Patient Portal**: Patients can view their own visits, diagnoses, prescriptions, and update their profile
- **Clinic Portal**: Clinics can manage their patients, staff, visits, diagnoses, and prescriptions
- **Admin Portal**: Full system access for administration purposes

## Response Format
All API responses follow a consistent format:

### Success Response (Single Item)
```json
{
  "id": 1,
  "field1": "value1",
  "field2": "value2",
  "created_at": "2025-06-28T09:37:07.43846Z",
  "updated_at": "2025-06-28T09:37:07.43846Z"
}
```

### Success Response (List with Pagination)
```json
{
  "data": [
    { /* item objects */ }
  ],
  "page": 1,
  "per_page": 10,
  "total": 25,
  "total_pages": 3
}
```

### Error Response
```json
{
  "error": "Error message",
  "details": "Additional error details if available"
}
```

## Authentication
The API uses JWT-based authentication. Include the token in the Authorization header:
```
Authorization: Bearer YOUR_JWT_TOKEN
```

### JWT Token Structure
The JWT tokens contain the following claims:
- `user_id`: Unique user identifier
- `email`: User's email address  
- `user_type`: User role ("patient", "clinic", or "admin")
- `patient_id`: Patient ID (only for patient users)
- `clinic_id`: Clinic ID (only for clinic users)
- `exp`: Token expiration timestamp

### Token Expiration
JWT tokens expire after 24 hours. The frontend should handle token expiration gracefully and redirect users to login when receiving 401 responses.

### Role-Based Access Control
- **Patient tokens**: Can only access patient portal endpoints (`/portal/patient/*`)
- **Clinic tokens**: Can only access clinic portal endpoints (`/portal/clinic/*`) 
- **Admin tokens**: Can access all system management endpoints

---

# API Endpoints

## 1. Authentication

### Register Patient
**POST** `/auth/register/patient`

Register a new patient with authentication credentials.

**Request Body:**
```json
{
  "email": "patient@example.com",
  "password": "securePassword123",
  "full_name": "John Doe",
  "gender": "Male",
  "date_of_birth": "1990-01-15",
  "address": "123 Village Road, Rural Area",
  "phone": "+1234567890",
  "clinic_id": 1
}
```

**Field Validations:**
- `email` (string, required): Valid email address
- `password` (string, required): Minimum 8 characters
- `full_name` (string, required): Patient's full name (2-255 characters)
- `gender` (string, required): Must be "Male", "Female", or "Other"
- `date_of_birth` (string, required): Date in YYYY-MM-DD format
- `address` (string, required): Full address (5-500 characters)
- `phone` (string, required): Phone number (10-20 characters)
- `clinic_id` (integer, required): Valid clinic ID

**Response (201 Created):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_type": "patient",
  "user": {
    "id": 1,
    "full_name": "John Doe",
    "gender": "Male",
    "date_of_birth": "1990-01-15T00:00:00Z",
    "address": "123 Village Road, Rural Area",
    "phone": "+1234567890",
    "clinic_id": 1,
    "created_at": "2025-06-28T09:38:48.047385Z",
    "updated_at": "2025-06-28T09:38:48.047385Z",
    "clinic": {
      "id": 1,
      "name": "Rural Health Center",
      "address": "456 Main Street",
      "contact_number": "+1234567891",
      "district": "Central District",
      "created_at": "2025-06-28T09:37:07.43846Z",
      "updated_at": "2025-06-28T09:37:07.43846Z"
    }
  }
}
```

### Register Clinic
**POST** `/auth/register/clinic`

Register a new clinic with authentication credentials.

**Request Body:**
```json
{
  "email": "clinic@example.com",
  "password": "securePassword123",
  "name": "Rural Health Center",
  "address": "456 Main Street, Central Village",
  "contact_number": "+1234567891",
  "district": "Central District"
}
```

**Field Validations:**
- `email` (string, required): Valid email address
- `password` (string, required): Minimum 8 characters
- `name` (string, required): Clinic name (2-255 characters)
- `address` (string, required): Full address (5-500 characters)
- `contact_number` (string, required): Phone number (10-20 characters)
- `district` (string, required): District name (2-100 characters)

**Response (201 Created):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_type": "clinic",
  "user": {
    "id": 1,
    "name": "Rural Health Center",
    "address": "456 Main Street, Central Village",
    "contact_number": "+1234567891",
    "district": "Central District",
    "created_at": "2025-06-28T09:37:07.43846Z",
    "updated_at": "2025-06-28T09:37:07.43846Z"
  }
}
```

### Login
**POST** `/auth/login`

Authenticate a user and receive a JWT token.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "userPassword123"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_type": "patient", // or "clinic" or "admin"
  "user": {
    // User profile data based on user type
  }
}
```

### Get Profile
**GET** `/auth/profile`

Get the current authenticated user's profile.

**Headers:**
```
Authorization: Bearer YOUR_JWT_TOKEN
```

**Response (200 OK):**
```json
{
  // Patient profile if user_type is "patient"
  "id": 1,
  "full_name": "John Doe",
  "gender": "Male",
  "date_of_birth": "1990-01-15T00:00:00Z",
  "address": "123 Village Road",
  "phone": "+1234567890",
  "clinic_id": 1,
  "created_at": "2025-06-28T09:38:48.047385Z",
  "updated_at": "2025-06-28T09:38:48.047385Z",
  "clinic": {
    // Associated clinic data
  }
}
```

### Change Password
**POST** `/auth/change-password`

Change the current user's password.

**Headers:**
```
Authorization: Bearer YOUR_JWT_TOKEN
```

**Request Body:**
```json
{
  "current_password": "currentPassword123",
  "new_password": "newPassword123"
}
```

**Response (200 OK):**
```json
{
  "message": "Password changed successfully"
}
```

---

## 2. Patient Portal

All patient portal endpoints require authentication and patient role.

**Headers:**
```
Authorization: Bearer PATIENT_JWT_TOKEN
```

### Get My Profile
**GET** `/portal/patient/profile`

Get the authenticated patient's profile information.

**Response (200 OK):**
```json
{
  "id": 1,
  "full_name": "John Doe",
  "gender": "Male",
  "date_of_birth": "1990-01-15T00:00:00Z",
  "address": "123 Village Road",
  "phone": "+1234567890",
  "clinic_id": 1,
  "created_at": "2025-06-28T09:38:48.047385Z",
  "updated_at": "2025-06-28T09:38:48.047385Z",
  "clinic": {
    "id": 1,
    "name": "Rural Health Center",
    "address": "456 Main Street",
    "contact_number": "+1234567891",
    "district": "Central District"
  }
}
```

### Update My Profile
**PUT** `/portal/patient/profile`

Update the authenticated patient's profile information.

**Request Body:**
```json
{
  "full_name": "John Doe Updated",
  "gender": "Male",
  "address": "123 Updated Village Road",
  "phone": "+1234567890",
  "date_of_birth": "1990-01-15"
}
```

**Note:** Patients cannot change their associated clinic. All fields are optional.

**Response (200 OK):** Updated patient profile in same format as Get My Profile.

### Get My Visits
**GET** `/portal/patient/visits`

Get all visits for the authenticated patient.

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `per_page` (integer, optional): Items per page (default: 10, max: 100)

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": 1,
      "patient_id": 1,
      "clinic_id": 1,
      "staff_id": 1,
      "visit_date": "2024-01-15T10:00:00Z",
      "reason": "Annual checkup",
      "notes": "Patient appears healthy",
      "created_at": "2025-06-28T09:42:32.645963Z",
      "updated_at": "2025-06-28T09:42:32.645963Z",
      "clinic": {
        "id": 1,
        "name": "Rural Health Center",
        "address": "456 Main Street",
        "contact_number": "+1234567891",
        "district": "Central District"
      },
      "staff": {
        "id": 1,
        "full_name": "Dr. Sarah Johnson",
        "role": "Doctor",
        "phone": "+1234567900",
        "email": "sarah.johnson@clinic.com"
      },
      "diagnoses": [
        {
          "id": 1,
          "diagnosis_code": "Z00.00",
          "description": "General health examination",
          "created_at": "2025-06-28T09:43:32.111086Z"
        }
      ],
      "prescriptions": [
        {
          "id": 1,
          "medication_name": "Acetaminophen",
          "dosage": "500mg",
          "instructions": "Take every 6 hours as needed",
          "duration_days": 7,
          "created_at": "2025-06-28T09:45:33.372822Z"
        }
      ]
    }
  ],
  "page": 1,
  "per_page": 10,
  "total": 1,
  "total_pages": 1
}
```

### Get My Visit by ID
**GET** `/portal/patient/visits/{id}`

Get details of a specific visit for the authenticated patient.

**Path Parameters:**
- `id` (integer, required): Visit ID

**Response (200 OK):** Same format as individual visit in Get My Visits.

### Get My Diagnoses
**GET** `/portal/patient/diagnoses`

Get all diagnoses for the authenticated patient.

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `per_page` (integer, optional): Items per page (default: 10, max: 100)

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": 1,
      "visit_id": 1,
      "diagnosis_code": "Z00.00",
      "description": "General health examination",
      "created_at": "2025-06-28T09:43:32.111086Z",
      "updated_at": "2025-06-28T09:43:32.111086Z",
      "visit": {
        "id": 1,
        "visit_date": "2024-01-15T10:00:00Z",
        "reason": "Annual checkup",
        "clinic": {
          "name": "Rural Health Center"
        },
        "staff": {
          "full_name": "Dr. Sarah Johnson",
          "role": "Doctor"
        }
      }
    }
  ],
  "page": 1,
  "per_page": 10,
  "total": 1,
  "total_pages": 1
}
```

### Get My Prescriptions
**GET** `/portal/patient/prescriptions`

Get all prescriptions for the authenticated patient.

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `per_page` (integer, optional): Items per page (default: 10, max: 100)

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": 1,
      "visit_id": 1,
      "medication_name": "Acetaminophen",
      "dosage": "500mg",
      "instructions": "Take every 6 hours as needed for fever",
      "duration_days": 7,
      "created_at": "2025-06-28T09:45:33.372822Z",
      "updated_at": "2025-06-28T09:45:33.372822Z",
      "visit": {
        "id": 1,
        "visit_date": "2024-01-15T10:00:00Z",
        "reason": "Annual checkup",
        "clinic": {
          "name": "Rural Health Center"
        },
        "staff": {
          "full_name": "Dr. Sarah Johnson",
          "role": "Doctor"
        }
      }
    }
  ],
  "page": 1,
  "per_page": 10,
  "total": 1,
  "total_pages": 1
}
```

---

## 3. Clinic Portal

All clinic portal endpoints require authentication and clinic role.

**Headers:**
```
Authorization: Bearer CLINIC_JWT_TOKEN
```

### Get My Profile
**GET** `/portal/clinic/profile`

Get the authenticated clinic's profile information.

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Rural Health Center",
  "address": "456 Main Street, Central Village",
  "contact_number": "+1234567891",
  "district": "Central District",
  "created_at": "2025-06-28T09:37:07.43846Z",
  "updated_at": "2025-06-28T09:37:07.43846Z"
}
```

### Update My Profile
**PUT** `/portal/clinic/profile`

Update the authenticated clinic's profile information.

**Request Body:**
```json
{
  "name": "Updated Rural Health Center",
  "address": "456 Updated Main Street, Central Village",
  "contact_number": "+1234567891",
  "district": "Central District"
}
```

**Response (200 OK):** Updated clinic profile in same format as Get My Profile.

### Get Dashboard Stats
**GET** `/portal/clinic/dashboard`

Get dashboard statistics for the authenticated clinic.

**Response (200 OK):**
```json
{
  "total_patients": 150,
  "total_staff": 8,
  "total_visits": 1250,
  "visits_this_month": 85,
  "total_diagnoses": 1100,
  "total_prescriptions": 950
}
```

### Get My Patients
**GET** `/portal/clinic/patients`

Get all patients for the authenticated clinic.

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `per_page` (integer, optional): Items per page (default: 10, max: 100)
- `search` (string, optional): Search by patient name or phone

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": 1,
      "full_name": "John Doe",
      "gender": "Male",
      "date_of_birth": "1990-01-15T00:00:00Z",
      "address": "123 Village Road",
      "phone": "+1234567890",
      "clinic_id": 1,
      "created_at": "2025-06-28T09:38:48.047385Z",
      "updated_at": "2025-06-28T09:38:48.047385Z"
    }
  ],
  "page": 1,
  "per_page": 10,
  "total": 150,
  "total_pages": 15
}
```

### Get My Patient by ID
**GET** `/portal/clinic/patients/{id}`

Get details of a specific patient for the authenticated clinic.

**Path Parameters:**
- `id` (integer, required): Patient ID

**Response (200 OK):**
```json
{
  "id": 1,
  "full_name": "John Doe",
  "gender": "Male",
  "date_of_birth": "1990-01-15T00:00:00Z",
  "address": "123 Village Road",
  "phone": "+1234567890",
  "clinic_id": 1,
  "created_at": "2025-06-28T09:38:48.047385Z",
  "updated_at": "2025-06-28T09:38:48.047385Z",
  "visits": [
    {
      "id": 1,
      "visit_date": "2024-01-15T10:00:00Z",
      "reason": "Annual checkup",
      "diagnoses": [...],
      "prescriptions": [...]
    }
  ]
}
```

### Get My Staff
**GET** `/portal/clinic/staff`

Get all staff for the authenticated clinic.

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `per_page` (integer, optional): Items per page (default: 10, max: 100)
- `role` (string, optional): Filter by role (Doctor, Nurse, Administrator, Pharmacist)

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": 1,
      "full_name": "Dr. Sarah Johnson",
      "role": "Doctor",
      "phone": "+1234567900",
      "email": "sarah.johnson@clinic.com",
      "clinic_id": 1,
      "created_at": "2025-06-28T09:40:04.003575Z",
      "updated_at": "2025-06-28T09:40:04.003575Z"
    }
  ],
  "page": 1,
  "per_page": 10,
  "total": 8,
  "total_pages": 1
}
```

### Create Staff
**POST** `/portal/clinic/staff`

Add a new staff member to the authenticated clinic.

**Request Body:**
```json
{
  "full_name": "Dr. Sarah Johnson",
  "role": "Doctor",
  "phone": "+1234567900",
  "email": "sarah.johnson@clinic.com"
}
```

**Field Validations:**
- `full_name` (string, required): Staff member's full name
- `role` (string, required): Must be "Doctor", "Nurse", "Administrator", or "Pharmacist"
- `phone` (string, required): Phone number
- `email` (string, required): Valid email address

**Response (201 Created):**
```json
{
  "id": 1,
  "full_name": "Dr. Sarah Johnson",
  "role": "Doctor",
  "phone": "+1234567900",
  "email": "sarah.johnson@clinic.com",
  "clinic_id": 1,
  "created_at": "2025-06-28T09:40:04.003575Z",
  "updated_at": "2025-06-28T09:40:04.003575Z",
  "clinic": {
    "id": 1,
    "name": "Rural Health Center"
  }
}
```

### Get My Visits
**GET** `/portal/clinic/visits`

Get all visits for the authenticated clinic.

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `per_page` (integer, optional): Items per page (default: 10, max: 100)
- `patient_id` (integer, optional): Filter by patient ID

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": 1,
      "patient_id": 1,
      "clinic_id": 1,
      "staff_id": 1,
      "visit_date": "2024-01-15T10:00:00Z",
      "reason": "Annual checkup",
      "notes": "Patient appears healthy",
      "created_at": "2025-06-28T09:42:32.645963Z",
      "updated_at": "2025-06-28T09:42:32.645963Z",
      "patient": {
        "id": 1,
        "full_name": "John Doe",
        "gender": "Male",
        "phone": "+1234567890"
      },
      "staff": {
        "id": 1,
        "full_name": "Dr. Sarah Johnson",
        "role": "Doctor"
      },
      "diagnoses": [...],
      "prescriptions": [...]
    }
  ],
  "page": 1,
  "per_page": 10,
  "total": 1250,
  "total_pages": 125
}
```

### Get My Visit by ID
**GET** `/portal/clinic/visits/{id}`

Get details of a specific visit for the authenticated clinic.

**Path Parameters:**
- `id` (integer, required): Visit ID

**Response (200 OK):** Same format as individual visit in Get My Visits.

### Create Visit
**POST** `/portal/clinic/visits`

Create a new visit for a patient in the authenticated clinic.

**Request Body:**
```json
{
  "patient_id": 1,
  "staff_id": 1,
  "visit_date": "2024-01-15T10:00:00Z",
  "reason": "Annual checkup",
  "notes": "Patient appears healthy"
}
```

**Field Validations:**
- `patient_id` (integer, required): Valid patient ID belonging to this clinic
- `staff_id` (integer, required): Valid staff ID belonging to this clinic
- `visit_date` (string, optional): ISO 8601 datetime. If not provided, current time is used
- `reason` (string, required): Reason for visit (5-500 characters)
- `notes` (string, optional): Additional notes (max 1000 characters)

**Response (201 Created):** Visit object with patient, clinic, and staff relationships.

### Create Diagnosis
**POST** `/portal/clinic/diagnoses`

Add a diagnosis to a visit for the authenticated clinic.

**Request Body:**
```json
{
  "visit_id": 1,
  "diagnosis_code": "Z00.00",
  "description": "General health examination"
}
```

**Field Validations:**
- `visit_id` (integer, required): Valid visit ID belonging to this clinic
- `diagnosis_code` (string, required): ICD-10 or similar code (2-20 characters)
- `description` (string, required): Diagnosis description (5-1000 characters)

**Response (201 Created):**
```json
{
  "id": 1,
  "visit_id": 1,
  "diagnosis_code": "Z00.00",
  "description": "General health examination",
  "created_at": "2025-06-28T09:43:32.111086Z",
  "updated_at": "2025-06-28T09:43:32.111086Z",
  "visit": {
    "id": 1,
    "visit_date": "2024-01-15T10:00:00Z",
    "reason": "Annual checkup"
  }
}
```

### Create Prescription
**POST** `/portal/clinic/prescriptions`

Add a prescription to a visit for the authenticated clinic.

**Request Body:**
```json
{
  "visit_id": 1,
  "medication_name": "Acetaminophen",
  "dosage": "500mg",
  "instructions": "Take every 6 hours as needed for fever",
  "duration_days": 7
}
```

**Field Validations:**
- `visit_id` (integer, required): Valid visit ID belonging to this clinic
- `medication_name` (string, required): Name of medication (2-255 characters)
- `dosage` (string, required): Dosage information (2-100 characters)
- `instructions` (string, required): Usage instructions (5-500 characters)
- `duration_days` (integer, required): Duration in days (1-365)

**Response (201 Created):**
```json
{
  "id": 1,
  "visit_id": 1,
  "medication_name": "Acetaminophen",
  "dosage": "500mg",
  "instructions": "Take every 6 hours as needed for fever",
  "duration_days": 7,
  "created_at": "2025-06-28T09:45:33.372822Z",
  "updated_at": "2025-06-28T09:45:33.372822Z",
  "visit": {
    "id": 1,
    "visit_date": "2024-01-15T10:00:00Z",
    "reason": "Annual checkup"
  }
}
```

---

## 4. Health Check

### Health Check
**GET** `/health`

Check the health status of the API.

**Response (200 OK):**
```json
{
  "status": "healthy",
  "timestamp": "2025-06-28T09:37:07.43846Z"
}
```

---

## 5. Admin Portal (System Management)

All admin endpoints require authentication and admin role. These endpoints are used for system administration and have access to all data across all clinics.

**Headers:**
```
Authorization: Bearer ADMIN_JWT_TOKEN
```

### Clinics Management (Admin)

### Get All Clinics (Admin)
**GET** `/clinics`

Retrieve a paginated list of all clinics (admin access only).

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `per_page` (integer, optional): Items per page (default: 10, max: 100)
- `search` (string, optional): Search by clinic name or address
- `district` (string, optional): Filter by district

**Example Request:**
```
GET /api/v1/clinics?page=1&per_page=10&search=central&district=Central District
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Central Rural Health Center",
      "address": "123 Main Street, Central Village",
      "contact_number": "+1234567890",
      "district": "Central District",
      "created_at": "2025-06-28T09:37:07.43846Z",
      "updated_at": "2025-06-28T09:37:07.43846Z"
    }
  ],
  "page": 1,
  "per_page": 10,
  "total": 1,
  "total_pages": 1
}
```

### Get Clinic by ID (Admin)
**GET** `/clinics/{id}`

Retrieve detailed information about a specific clinic (admin access only).

**Path Parameters:**
- `id` (integer, required): Clinic ID

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Central Rural Health Center",
  "address": "123 Main Street, Central Village",
  "contact_number": "+1234567890",
  "district": "Central District",
  "created_at": "2025-06-28T09:37:07.43846Z",
  "updated_at": "2025-06-28T09:37:07.43846Z"
}
```

### Create Clinic (Admin)
**POST** `/clinics`

Create a new clinic (admin access only).

**Request Body:**
```json
{
  "name": "Central Rural Health Center",
  "address": "123 Main Street, Central Village",
  "contact_number": "+1234567890",
  "district": "Central District"
}
```

**Note:** This creates a clinic without authentication. For clinics with authentication, use `/auth/register/clinic`.

### Update Clinic (Admin)
**PUT** `/clinics/{id}`

Update an existing clinic (admin access only).

### Delete Clinic (Admin)
**DELETE** `/clinics/{id}`

Soft delete a clinic (admin access only). Cannot delete if clinic has associated patients or staff.

---

## 6. Patients Management (Admin)

### Get All Patients (Admin)
**GET** `/patients`

Retrieve a paginated list of all patients (admin access only).

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `per_page` (integer, optional): Items per page (default: 10, max: 100)
- `search` (string, optional): Search by full name or phone number
- `clinic_id` (integer, optional): Filter by clinic ID

**Note:** For clinic-specific patient access, use `/portal/clinic/patients`.

### Get Patient by ID (Admin)
**GET** `/patients/{id}`

Retrieve detailed information about a specific patient (admin access only).

### Create Patient (Admin)
**POST** `/patients`

Create a new patient (admin access only).

**Note:** This creates a patient without authentication. For patients with authentication, use `/auth/register/patient`.

### Update Patient (Admin)
**PUT** `/patients/{id}`

Update an existing patient (admin access only).

### Delete Patient (Admin)
**DELETE** `/patients/{id}`

Soft delete a patient (admin access only).

---

## 7. Staff Management (Admin)

### Get All Staff (Admin)
**GET** `/staff`

Retrieve a paginated list of all staff members (admin access only).

**Note:** For clinic-specific staff access, use `/portal/clinic/staff`.

### Get Staff by ID (Admin)
**GET** `/staff/{id}`

Retrieve detailed information about a specific staff member (admin access only).

### Create Staff (Admin)
**POST** `/staff`

Create a new staff member (admin access only).

**Note:** For clinic-managed staff creation, use `/portal/clinic/staff`.

### Update Staff (Admin)
**PUT** `/staff/{id}`

Update an existing staff member (admin access only).

### Delete Staff (Admin)
**DELETE** `/staff/{id}`

Soft delete a staff member (admin access only).

---

## 8. Visits Management (Admin)

### Get All Visits (Admin)
**GET** `/visits`

Retrieve a paginated list of all visits (admin access only).

**Note:** For patient-specific visits, use `/portal/patient/visits`. For clinic-specific visits, use `/portal/clinic/visits`.

### Get Visit by ID (Admin)
**GET** `/visits/{id}`

Retrieve detailed information about a specific visit (admin access only).

### Create Visit (Admin)
**POST** `/visits`

Create a new visit (admin access only).

**Note:** For clinic-managed visit creation, use `/portal/clinic/visits`.

### Update Visit (Admin)
**PUT** `/visits/{id}`

Update an existing visit (admin access only).

### Delete Visit (Admin)
**DELETE** `/visits/{id}`

Soft delete a visit (admin access only).

---

## 9. Diagnoses Management (Admin)

### Get All Diagnoses (Admin)
**GET** `/diagnoses`

Retrieve a paginated list of all diagnoses (admin access only).

**Note:** For patient-specific diagnoses, use `/portal/patient/diagnoses`. For clinic-managed diagnoses, use `/portal/clinic/diagnoses`.

### Get Diagnosis by ID (Admin)
**GET** `/diagnoses/{id}`

Retrieve detailed information about a specific diagnosis (admin access only).

### Create Diagnosis (Admin)
**POST** `/diagnoses`

Create a new diagnosis (admin access only).

**Note:** For clinic-managed diagnosis creation, use `/portal/clinic/diagnoses`.

### Update Diagnosis (Admin)
**PUT** `/diagnoses/{id}`

Update an existing diagnosis (admin access only).

### Delete Diagnosis (Admin)
**DELETE** `/diagnoses/{id}`

Soft delete a diagnosis (admin access only).

---

## 10. Prescriptions Management (Admin)

### Get All Prescriptions (Admin)
**GET** `/prescriptions`

Retrieve a paginated list of all prescriptions (admin access only).

**Note:** For patient-specific prescriptions, use `/portal/patient/prescriptions`. For clinic-managed prescriptions, use `/portal/clinic/prescriptions`.

### Get Prescription by ID (Admin)
**GET** `/prescriptions/{id}`

Retrieve detailed information about a specific prescription (admin access only).

### Create Prescription (Admin)
**POST** `/prescriptions`

Create a new prescription (admin access only).

**Note:** For clinic-managed prescription creation, use `/portal/clinic/prescriptions`.

### Update Prescription (Admin)
**PUT** `/prescriptions/{id}`

Update an existing prescription (admin access only).

### Delete Prescription (Admin)
**DELETE** `/prescriptions/{id}`

Soft delete a prescription (admin access only).

---

## Data Model Relationships

### Entity Relationships
1. **Clinic** → has many → **Patients**, **Staff**
2. **Patient** → belongs to → **Clinic**
3. **Patient** → has many → **Visits**
4. **Staff** → belongs to → **Clinic**
5. **Staff** → has many → **Visits** (as attending staff)
6. **Visit** → belongs to → **Patient**, **Clinic**, **Staff**
7. **Visit** → has many → **Diagnoses**, **Prescriptions**
8. **Diagnosis** → belongs to → **Visit**
9. **Prescription** → belongs to → **Visit**

### Required Fields Summary
- **Clinic**: name, address, contact_number, district
- **Patient**: full_name, gender, date_of_birth, address, phone, clinic_id
- **Staff**: full_name, role, phone, email, clinic_id
- **Visit**: patient_id, clinic_id, staff_id, reason
- **Diagnosis**: visit_id, diagnosis_code, description
- **Prescription**: visit_id, medication_name, dosage, instructions, duration_days

### Enumerated Values
- **Patient Gender**: "Male", "Female", "Other"
- **Staff Role**: "Doctor", "Nurse", "Administrator", "Pharmacist"

### Date/Time Fields
- All `created_at` and `updated_at` fields are in ISO 8601 format
- `date_of_birth` is stored as date (YYYY-MM-DD format for input)
- `visit_date` is stored as datetime (ISO 8601 format)

### Pagination
- Default page size: 10 items
- Maximum page size: 100 items
- All list endpoints support pagination
- Response includes: `page`, `per_page`, `total`, `total_pages`

### Soft Deletes
- All delete operations are soft deletes (records are marked as deleted but not removed)
- Relationships prevent deletion when dependent records exist
- Cannot delete clinic with patients/staff
- Cannot delete patient with visits
- Cannot delete staff with visits
- Cannot delete visit with diagnoses/prescriptions

---

## Common HTTP Status Codes

- **200 OK**: Successful GET, PUT requests
- **201 Created**: Successful POST requests  
- **204 No Content**: Successful DELETE requests
- **400 Bad Request**: Invalid request data or malformed JSON
- **401 Unauthorized**: Missing, invalid, or expired JWT token
- **403 Forbidden**: Insufficient permissions for the requested resource
- **404 Not Found**: Resource not found
- **409 Conflict**: Resource already exists (e.g., email already registered)
- **422 Unprocessable Entity**: Validation errors in request data
- **500 Internal Server Error**: Server error

## Authentication Error Responses

### 401 Unauthorized
Returned when authentication is required but missing or invalid:

```json
{
  "error": "Authorization header required"
}
```

```json
{
  "error": "Invalid authorization header format"
}
```

```json
{
  "error": "Invalid token"
}
```

```json
{
  "error": "Token expired"
}
```

### 403 Forbidden
Returned when user doesn't have permission to access the resource:

```json
{
  "error": "Insufficient permissions"
}
```

```json
{
  "error": "You can only access your own data"
}
```

### 409 Conflict
Returned during registration when email already exists:

```json
{
  "error": "Email already registered"
}
```

## Data Validation Rules

### Common Validation Constraints
- **Email fields**: Must be valid email format
- **Password fields**: Minimum 8 characters
- **Phone fields**: 10-20 characters
- **Name fields**: 2-255 characters
- **Address fields**: 5-500 characters
- **Gender field**: Must be "Male", "Female", or "Other"
- **User Type**: Must be "patient", "clinic", or "admin"
- **Staff Role**: Must be "Doctor", "Nurse", "Administrator", or "Pharmacist"
- **Date fields**: Use ISO 8601 format (YYYY-MM-DD or YYYY-MM-DDTHH:MM:SSZ)

### Validation Error Response
When validation fails, the API returns a 422 status with details:

```json
{
  "error": "Validation failed",
  "details": {
    "field_name": "error message for specific field"
  }
}
```

## Sample Workflow for Frontend Development

### For Patient Portal Frontend:
1. **Patient Registration** → POST `/auth/register/patient`
2. **Patient Login** → POST `/auth/login`
3. **View Profile** → GET `/portal/patient/profile`
4. **Update Profile** → PUT `/portal/patient/profile`
5. **View Medical History** → GET `/portal/patient/visits`
6. **View Specific Visit** → GET `/portal/patient/visits/{id}`
7. **View Diagnoses** → GET `/portal/patient/diagnoses`
8. **View Prescriptions** → GET `/portal/patient/prescriptions`

### For Clinic Portal Frontend:
1. **Clinic Registration** → POST `/auth/register/clinic`
2. **Clinic Login** → POST `/auth/login`
3. **View Dashboard** → GET `/portal/clinic/dashboard`
4. **Manage Staff** → GET/POST `/portal/clinic/staff`
5. **Manage Patients** → GET `/portal/clinic/patients`
6. **Record Visits** → POST `/portal/clinic/visits`
7. **Add Diagnoses** → POST `/portal/clinic/diagnoses`
8. **Add Prescriptions** → POST `/portal/clinic/prescriptions`

### For Admin System Frontend:
1. **Admin Login** → POST `/auth/login`
2. **System Management** → Use admin endpoints (`/clinics`, `/patients`, `/staff`, etc.)
3. **Full CRUD Operations** → All admin endpoints support GET, POST, PUT, DELETE

### Frontend Architecture Recommendations:
- **Separate Portals**: Build separate interfaces for patient, clinic, and admin users
- **Route Protection**: Use JWT user_type to protect routes based on user roles
- **Data Isolation**: Patients see only their data, clinics see only their data
- **Token Management**: Store JWT tokens securely and handle expiration gracefully
- **Error Handling**: Implement proper error handling for all HTTP status codes

This API provides a complete foundation for building a comprehensive rural health management system frontend with features for clinic management, patient registration, visit tracking, medical records, and prescription management.
