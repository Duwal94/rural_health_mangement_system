# Rural Health Management System API Documentation

## Base URL
```
http://localhost:3000/api/v1
```

## Overview
The Rural Health Management System API provides comprehensive endpoints for managing rural healthcare facilities, including clinics, patients, staff members, visits, diagnoses, and prescriptions. All endpoints support pagination, filtering, and search capabilities where applicable.

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
Currently, the API does not require authentication. This may change in future versions.

---

# API Endpoints

## 1. Health Check

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

## 2. Clinics Management

### Get All Clinics
**GET** `/clinics`

Retrieve a paginated list of all clinics with optional filtering.

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

### Get Clinic by ID
**GET** `/clinics/{id}`

Retrieve detailed information about a specific clinic.

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

### Create Clinic
**POST** `/clinics`

Create a new clinic.

**Request Body:**
```json
{
  "name": "Central Rural Health Center",
  "address": "123 Main Street, Central Village",
  "contact_number": "+1234567890",
  "district": "Central District"
}
```

**Field Validations:**
- `name` (string, required): Clinic name (max 255 characters)
- `address` (string, required): Full address (max 500 characters)
- `contact_number` (string, required): Phone number with country code
- `district` (string, required): District name (max 255 characters)

**Response (201 Created):**
```json
{
  "id": 1,
  "name": "Central Rural Health Center",
  "address": "123 Main Street, Central Village",
  "contact_number": "+1234567890",
  "district": "Central District",
  "created_at": "2025-06-28T09:37:07.438460838Z",
  "updated_at": "2025-06-28T09:37:07.438460838Z"
}
```

### Update Clinic
**PUT** `/clinics/{id}`

Update an existing clinic.

**Path Parameters:**
- `id` (integer, required): Clinic ID

**Request Body:**
```json
{
  "name": "Updated Central Rural Health Center",
  "address": "123 Updated Main Street, Central Village",
  "contact_number": "+1234567890",
  "district": "Central District"
}
```

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Updated Central Rural Health Center",
  "address": "123 Updated Main Street, Central Village",
  "contact_number": "+1234567890",
  "district": "Central District",
  "created_at": "2025-06-28T09:37:07.43846Z",
  "updated_at": "2025-06-28T09:38:16.560056728Z"
}
```

### Delete Clinic
**DELETE** `/clinics/{id}`

Soft delete a clinic. Cannot delete if clinic has associated patients or staff.

**Path Parameters:**
- `id` (integer, required): Clinic ID

**Response (204 No Content)**

---

## 3. Patients Management

### Get All Patients
**GET** `/patients`

Retrieve a paginated list of all patients with optional filtering.

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `per_page` (integer, optional): Items per page (default: 10, max: 100)
- `search` (string, optional): Search by full name or phone number
- `clinic_id` (integer, optional): Filter by clinic ID

**Example Request:**
```
GET /api/v1/patients?page=1&per_page=10&search=Alice&clinic_id=1
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": 1,
      "full_name": "Alice Cooper",
      "gender": "Female",
      "date_of_birth": "1985-03-15T00:00:00Z",
      "address": "101 Village Lane, Central Village",
      "phone": "+1234567910",
      "clinic_id": 1,
      "created_at": "2025-06-28T09:38:48.047385Z",
      "updated_at": "2025-06-28T09:38:48.047385Z",
      "clinic": {
        "id": 1,
        "name": "Updated Central Rural Health Center",
        "address": "123 Updated Main Street, Central Village",
        "contact_number": "+1234567890",
        "district": "Central District",
        "created_at": "2025-06-28T09:37:07.43846Z",
        "updated_at": "2025-06-28T09:38:16.560056Z"
      }
    }
  ],
  "page": 1,
  "per_page": 10,
  "total": 1,
  "total_pages": 1
}
```

### Get Patient by ID
**GET** `/patients/{id}`

Retrieve detailed information about a specific patient with clinic information.

**Path Parameters:**
- `id` (integer, required): Patient ID

**Response (200 OK):**
```json
{
  "id": 1,
  "full_name": "Alice Cooper",
  "gender": "Female",
  "date_of_birth": "1985-03-15T00:00:00Z",
  "address": "101 Village Lane, Central Village",
  "phone": "+1234567910",
  "clinic_id": 1,
  "created_at": "2025-06-28T09:38:48.047385Z",
  "updated_at": "2025-06-28T09:38:48.047385Z",
  "clinic": {
    "id": 1,
    "name": "Updated Central Rural Health Center",
    "address": "123 Updated Main Street, Central Village",
    "contact_number": "+1234567890",
    "district": "Central District",
    "created_at": "2025-06-28T09:37:07.43846Z",
    "updated_at": "2025-06-28T09:38:16.560056Z"
  }
}
```

### Create Patient
**POST** `/patients`

Create a new patient.

**Request Body:**
```json
{
  "full_name": "Alice Cooper",
  "gender": "Female",
  "date_of_birth": "1985-03-15",
  "address": "101 Village Lane, Central Village",
  "phone": "+1234567910",
  "clinic_id": 1
}
```

**Field Validations:**
- `full_name` (string, required): Patient's full name (max 255 characters)
- `gender` (string, required): Must be "Male", "Female", or "Other"
- `date_of_birth` (string, required): Date in YYYY-MM-DD format
- `address` (string, required): Full address (max 500 characters)
- `phone` (string, required): Phone number with country code
- `clinic_id` (integer, required): Valid clinic ID

**Response (201 Created):**
```json
{
  "id": 1,
  "full_name": "Alice Cooper",
  "gender": "Female",
  "date_of_birth": "1985-03-15T00:00:00Z",
  "address": "101 Village Lane, Central Village",
  "phone": "+1234567910",
  "clinic_id": 1,
  "created_at": "2025-06-28T09:38:48.047385Z",
  "updated_at": "2025-06-28T09:38:48.047385Z",
  "clinic": {
    "id": 1,
    "name": "Updated Central Rural Health Center",
    "address": "123 Updated Main Street, Central Village",
    "contact_number": "+1234567890",
    "district": "Central District",
    "created_at": "2025-06-28T09:37:07.43846Z",
    "updated_at": "2025-06-28T09:38:16.560056Z"
  }
}
```

### Update Patient
**PUT** `/patients/{id}`

Update an existing patient. All fields are optional.

**Path Parameters:**
- `id` (integer, required): Patient ID

**Request Body:**
```json
{
  "full_name": "Alice Cooper Updated",
  "gender": "Female",
  "date_of_birth": "1985-03-15",
  "address": "101 Updated Village Lane, Central Village",
  "phone": "+1234567910",
  "clinic_id": 1
}
```

**Response (200 OK):** Same format as Create Patient response with updated data.

### Delete Patient
**DELETE** `/patients/{id}`

Soft delete a patient. Cannot delete if patient has existing visits.

**Path Parameters:**
- `id` (integer, required): Patient ID

**Response (204 No Content)**

---

## 4. Staff Management

### Get All Staff
**GET** `/staff`

Retrieve a paginated list of all staff members with optional filtering.

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `per_page` (integer, optional): Items per page (default: 10, max: 100)
- `clinic_id` (integer, optional): Filter by clinic ID
- `role` (string, optional): Filter by role (Doctor, Nurse, Administrator, Pharmacist)

**Example Request:**
```
GET /api/v1/staff?page=1&per_page=10&clinic_id=1&role=Doctor
```

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
      "updated_at": "2025-06-28T09:40:04.003575Z",
      "clinic": {
        "id": 1,
        "name": "Updated Central Rural Health Center",
        "address": "123 Updated Main Street, Central Village",
        "contact_number": "+1234567890",
        "district": "Central District",
        "created_at": "2025-06-28T09:37:07.43846Z",
        "updated_at": "2025-06-28T09:38:16.560056Z"
      }
    }
  ],
  "page": 1,
  "per_page": 10,
  "total": 1,
  "total_pages": 1
}
```

### Get Staff by ID
**GET** `/staff/{id}`

Retrieve detailed information about a specific staff member with clinic information.

**Path Parameters:**
- `id` (integer, required): Staff ID

**Response (200 OK):**
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
    "name": "Updated Central Rural Health Center",
    "address": "123 Updated Main Street, Central Village",
    "contact_number": "+1234567890",
    "district": "Central District",
    "created_at": "2025-06-28T09:37:07.43846Z",
    "updated_at": "2025-06-28T09:38:16.560056Z"
  }
}
```

### Create Staff
**POST** `/staff`

Create a new staff member.

**Request Body:**
```json
{
  "full_name": "Dr. Sarah Johnson",
  "role": "Doctor",
  "phone": "+1234567900",
  "email": "sarah.johnson@clinic.com",
  "clinic_id": 1
}
```

**Field Validations:**
- `full_name` (string, required): Staff member's full name (max 255 characters)
- `role` (string, required): Must be "Doctor", "Nurse", "Administrator", or "Pharmacist"
- `phone` (string, required): Phone number with country code
- `email` (string, required): Valid email address
- `clinic_id` (integer, required): Valid clinic ID

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
    "name": "Updated Central Rural Health Center",
    "address": "123 Updated Main Street, Central Village",
    "contact_number": "+1234567890",
    "district": "Central District",
    "created_at": "2025-06-28T09:37:07.43846Z",
    "updated_at": "2025-06-28T09:38:16.560056Z"
  }
}
```

### Update Staff
**PUT** `/staff/{id}`

Update an existing staff member.

**Path Parameters:**
- `id` (integer, required): Staff ID

**Request Body:**
```json
{
  "full_name": "Dr. Sarah Johnson Updated",
  "role": "Doctor",
  "phone": "+1234567900",
  "email": "sarah.johnson.updated@clinic.com",
  "clinic_id": 1
}
```

**Response (200 OK):**
```json
{
  "id": 1,
  "full_name": "Dr. Sarah Johnson Updated",
  "role": "Doctor",
  "phone": "+1234567900",
  "email": "sarah.johnson.updated@clinic.com",
  "clinic_id": 1,
  "created_at": "2025-06-28T09:40:04.003575Z",
  "updated_at": "2025-06-28T09:41:14.173253Z",
  "clinic": {
    "id": 1,
    "name": "Updated Central Rural Health Center",
    "address": "123 Updated Main Street, Central Village",
    "contact_number": "+1234567890",
    "district": "Central District",
    "created_at": "2025-06-28T09:37:07.43846Z",
    "updated_at": "2025-06-28T09:38:16.560056Z"
  }
}
```

### Delete Staff
**DELETE** `/staff/{id}`

Soft delete a staff member. Cannot delete if staff has existing visits.

**Path Parameters:**
- `id` (integer, required): Staff ID

**Response (204 No Content)**

---

## 5. Visits Management

### Get All Visits
**GET** `/visits`

Retrieve a paginated list of all visits with patient, clinic, and staff information.

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `per_page` (integer, optional): Items per page (default: 10, max: 100)
- `patient_id` (integer, optional): Filter by patient ID
- `clinic_id` (integer, optional): Filter by clinic ID

**Example Request:**
```
GET /api/v1/visits?page=1&per_page=10&patient_id=1&clinic_id=1
```

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
      "notes": "Patient appears healthy, all vitals normal",
      "created_at": "2025-06-28T09:42:32.645963Z",
      "updated_at": "2025-06-28T09:42:32.645963Z",
      "patient": {
        "id": 1,
        "full_name": "Alice Cooper Updated",
        "gender": "Female",
        "date_of_birth": "1985-03-15T00:00:00Z",
        "address": "101 Updated Village Lane, Central Village",
        "phone": "+1234567910",
        "clinic_id": 1,
        "created_at": "2025-06-28T09:38:48.047385Z",
        "updated_at": "2025-06-28T09:39:44.982823Z"
      },
      "clinic": {
        "id": 1,
        "name": "Updated Central Rural Health Center",
        "address": "123 Updated Main Street, Central Village",
        "contact_number": "+1234567890",
        "district": "Central District",
        "created_at": "2025-06-28T09:37:07.43846Z",
        "updated_at": "2025-06-28T09:38:16.560056Z"
      },
      "staff": {
        "id": 1,
        "full_name": "Dr. Sarah Johnson Updated",
        "role": "Doctor",
        "phone": "+1234567900",
        "email": "sarah.johnson.updated@clinic.com",
        "clinic_id": 1,
        "created_at": "2025-06-28T09:40:04.003575Z",
        "updated_at": "2025-06-28T09:41:14.173253Z"
      }
    }
  ],
  "page": 1,
  "per_page": 10,
  "total": 1,
  "total_pages": 1
}
```

### Get Visit by ID
**GET** `/visits/{id}`

Retrieve detailed information about a specific visit with patient, clinic, staff, diagnoses, and prescriptions.

**Path Parameters:**
- `id` (integer, required): Visit ID

**Response (200 OK):**
```json
{
  "id": 1,
  "patient_id": 1,
  "clinic_id": 1,
  "staff_id": 1,
  "visit_date": "2024-01-15T10:00:00Z",
  "reason": "Annual checkup",
  "notes": "Patient appears healthy, all vitals normal",
  "created_at": "2025-06-28T09:42:32.645963Z",
  "updated_at": "2025-06-28T09:42:32.645963Z",
  "patient": {
    "id": 1,
    "full_name": "Alice Cooper Updated",
    "gender": "Female",
    "date_of_birth": "1985-03-15T00:00:00Z",
    "address": "101 Updated Village Lane, Central Village",
    "phone": "+1234567910",
    "clinic_id": 1,
    "created_at": "2025-06-28T09:38:48.047385Z",
    "updated_at": "2025-06-28T09:39:44.982823Z"
  },
  "clinic": {
    "id": 1,
    "name": "Updated Central Rural Health Center",
    "address": "123 Updated Main Street, Central Village",
    "contact_number": "+1234567890",
    "district": "Central District",
    "created_at": "2025-06-28T09:37:07.43846Z",
    "updated_at": "2025-06-28T09:38:16.560056Z"
  },
  "staff": {
    "id": 1,
    "full_name": "Dr. Sarah Johnson Updated",
    "role": "Doctor",
    "phone": "+1234567900",
    "email": "sarah.johnson.updated@clinic.com",
    "clinic_id": 1,
    "created_at": "2025-06-28T09:40:04.003575Z",
    "updated_at": "2025-06-28T09:41:14.173253Z"
  }
}
```

### Create Visit
**POST** `/visits`

Create a new visit. If visit_date is not provided, current time will be used.

**Request Body:**
```json
{
  "patient_id": 1,
  "clinic_id": 1,
  "staff_id": 1,
  "visit_date": "2024-01-15T10:00:00Z",
  "reason": "Annual checkup",
  "notes": "Patient appears healthy, all vitals normal"
}
```

**Field Validations:**
- `patient_id` (integer, required): Valid patient ID
- `clinic_id` (integer, required): Valid clinic ID
- `staff_id` (integer, required): Valid staff ID
- `visit_date` (string, optional): ISO 8601 datetime format. If not provided, current time is used
- `reason` (string, required): Reason for visit (max 500 characters)
- `notes` (string, optional): Additional notes (max 1000 characters)

**Response (201 Created):** Same format as Get Visit by ID response.

### Update Visit
**PUT** `/visits/{id}`

Update an existing visit.

**Path Parameters:**
- `id` (integer, required): Visit ID

**Request Body:**
```json
{
  "patient_id": 1,
  "clinic_id": 1,
  "staff_id": 1,
  "visit_date": "2024-01-15T10:00:00Z",
  "reason": "Updated annual checkup",
  "notes": "Patient appears healthy, all vitals normal - updated notes"
}
```

**Response (200 OK):** Same format as Get Visit by ID response with updated data.

### Delete Visit
**DELETE** `/visits/{id}`

Soft delete a visit. Cannot delete if visit has existing diagnoses or prescriptions.

**Path Parameters:**
- `id` (integer, required): Visit ID

**Response (204 No Content)**

---

## 6. Diagnoses Management

### Get All Diagnoses
**GET** `/diagnoses`

Retrieve a paginated list of all diagnoses with visit information.

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `per_page` (integer, optional): Items per page (default: 10, max: 100)
- `visit_id` (integer, optional): Filter by visit ID

**Example Request:**
```
GET /api/v1/diagnoses?page=1&per_page=10&visit_id=1
```

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
        "patient_id": 1,
        "clinic_id": 1,
        "staff_id": 1,
        "visit_date": "2024-01-15T10:00:00Z",
        "reason": "Updated annual checkup",
        "notes": "Patient appears healthy, all vitals normal - updated notes",
        "created_at": "2025-06-28T09:42:32.645963Z",
        "updated_at": "2025-06-28T09:42:56.175682Z"
      }
    }
  ],
  "page": 1,
  "per_page": 10,
  "total": 1,
  "total_pages": 1
}
```

### Get Diagnosis by ID
**GET** `/diagnoses/{id}`

Retrieve detailed information about a specific diagnosis with visit information.

**Path Parameters:**
- `id` (integer, required): Diagnosis ID

**Response (200 OK):**
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
    "patient_id": 1,
    "clinic_id": 1,
    "staff_id": 1,
    "visit_date": "2024-01-15T10:00:00Z",
    "reason": "Updated annual checkup",
    "notes": "Patient appears healthy, all vitals normal - updated notes",
    "created_at": "2025-06-28T09:42:32.645963Z",
    "updated_at": "2025-06-28T09:42:56.175682Z"
  }
}
```

### Create Diagnosis
**POST** `/diagnoses`

Create a new diagnosis for a visit.

**Request Body:**
```json
{
  "visit_id": 1,
  "diagnosis_code": "Z00.00",
  "description": "General health examination"
}
```

**Field Validations:**
- `visit_id` (integer, required): Valid visit ID
- `diagnosis_code` (string, required): ICD-10 or similar diagnosis code (max 50 characters)
- `description` (string, required): Diagnosis description (max 500 characters)

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
    "patient_id": 1,
    "clinic_id": 1,
    "staff_id": 1,
    "visit_date": "2024-01-15T10:00:00Z",
    "reason": "Updated annual checkup",
    "notes": "Patient appears healthy, all vitals normal - updated notes",
    "created_at": "2025-06-28T09:42:32.645963Z",
    "updated_at": "2025-06-28T09:42:56.175682Z"
  }
}
```

### Update Diagnosis
**PUT** `/diagnoses/{id}`

Update an existing diagnosis.

**Path Parameters:**
- `id` (integer, required): Diagnosis ID

**Request Body:**
```json
{
  "diagnosis_code": "Z00.01",
  "description": "Updated general health examination"
}
```

**Response (200 OK):**
```json
{
  "id": 1,
  "visit_id": 1,
  "diagnosis_code": "Z00.01",
  "description": "Updated general health examination",
  "created_at": "2025-06-28T09:43:32.111086Z",
  "updated_at": "2025-06-28T09:44:14.86456Z",
  "visit": {
    "id": 1,
    "patient_id": 1,
    "clinic_id": 1,
    "staff_id": 1,
    "visit_date": "2024-01-15T10:00:00Z",
    "reason": "Updated annual checkup",
    "notes": "Patient appears healthy, all vitals normal - updated notes",
    "created_at": "2025-06-28T09:42:32.645963Z",
    "updated_at": "2025-06-28T09:42:56.175682Z"
  }
}
```

### Delete Diagnosis
**DELETE** `/diagnoses/{id}`

Soft delete a diagnosis.

**Path Parameters:**
- `id` (integer, required): Diagnosis ID

**Response (204 No Content)**

---

## 7. Prescriptions Management

### Get All Prescriptions
**GET** `/prescriptions`

Retrieve a paginated list of all prescriptions with visit information.

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `per_page` (integer, optional): Items per page (default: 10, max: 100)
- `visit_id` (integer, optional): Filter by visit ID

**Example Request:**
```
GET /api/v1/prescriptions?page=1&per_page=10&visit_id=1
```

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
        "patient_id": 1,
        "clinic_id": 1,
        "staff_id": 1,
        "visit_date": "2024-01-15T10:00:00Z",
        "reason": "Updated annual checkup",
        "notes": "Patient appears healthy, all vitals normal - updated notes",
        "created_at": "2025-06-28T09:42:32.645963Z",
        "updated_at": "2025-06-28T09:42:56.175682Z"
      }
    }
  ],
  "page": 1,
  "per_page": 10,
  "total": 1,
  "total_pages": 1
}
```

### Get Prescription by ID
**GET** `/prescriptions/{id}`

Retrieve detailed information about a specific prescription with visit information.

**Path Parameters:**
- `id` (integer, required): Prescription ID

**Response (200 OK):**
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
    "patient_id": 1,
    "clinic_id": 1,
    "staff_id": 1,
    "visit_date": "2024-01-15T10:00:00Z",
    "reason": "Updated annual checkup",
    "notes": "Patient appears healthy, all vitals normal - updated notes",
    "created_at": "2025-06-28T09:42:32.645963Z",
    "updated_at": "2025-06-28T09:42:56.175682Z"
  }
}
```

### Create Prescription
**POST** `/prescriptions`

Create a new prescription for a visit.

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
- `visit_id` (integer, required): Valid visit ID
- `medication_name` (string, required): Name of the medication (max 255 characters)
- `dosage` (string, required): Dosage information (max 100 characters)
- `instructions` (string, required): Instructions for taking the medication (max 1000 characters)
- `duration_days` (integer, required): Number of days the prescription is valid (min: 1, max: 365)

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
    "patient_id": 1,
    "clinic_id": 1,
    "staff_id": 1,
    "visit_date": "2024-01-15T10:00:00Z",
    "reason": "Updated annual checkup",
    "notes": "Patient appears healthy, all vitals normal - updated notes",
    "created_at": "2025-06-28T09:42:32.645963Z",
    "updated_at": "2025-06-28T09:42:56.175682Z"
  }
}
```

### Update Prescription
**PUT** `/prescriptions/{id}`

Update an existing prescription.

**Path Parameters:**
- `id` (integer, required): Prescription ID

**Request Body:**
```json
{
  "medication_name": "Acetaminophen Updated",
  "dosage": "750mg",
  "instructions": "Take every 8 hours as needed for fever",
  "duration_days": 10
}
```

**Response (200 OK):**
```json
{
  "id": 1,
  "visit_id": 1,
  "medication_name": "Acetaminophen Updated",
  "dosage": "750mg",
  "instructions": "Take every 8 hours as needed for fever",
  "duration_days": 10,
  "created_at": "2025-06-28T09:45:33.372822Z",
  "updated_at": "2025-06-28T09:46:17.024381Z",
  "visit": {
    "id": 1,
    "patient_id": 1,
    "clinic_id": 1,
    "staff_id": 1,
    "visit_date": "2024-01-15T10:00:00Z",
    "reason": "Updated annual checkup",
    "notes": "Patient appears healthy, all vitals normal - updated notes",
    "created_at": "2025-06-28T09:42:32.645963Z",
    "updated_at": "2025-06-28T09:42:56.175682Z"
  }
}
```

### Delete Prescription
**DELETE** `/prescriptions/{id}`

Soft delete a prescription.

**Path Parameters:**
- `id` (integer, required): Prescription ID

**Response (204 No Content)**

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
- **400 Bad Request**: Invalid request data
- **404 Not Found**: Resource not found
- **422 Unprocessable Entity**: Validation errors
- **500 Internal Server Error**: Server error

## Sample Workflow for Frontend Development

1. **Create a new clinic** → POST `/clinics`
2. **Add staff to clinic** → POST `/staff`
3. **Register patients** → POST `/patients`
4. **Record patient visits** → POST `/visits`
5. **Add diagnoses to visits** → POST `/diagnoses`
6. **Add prescriptions to visits** → POST `/prescriptions`
7. **View patient history** → GET `/visits?patient_id={id}`
8. **Search and filter** → Use query parameters on list endpoints

This API provides a complete foundation for building a comprehensive rural health management system frontend with features for clinic management, patient registration, visit tracking, medical records, and prescription management.
