# Clinic Staff Workflow Guide

## Overview

This guide details the complete workflow for clinic administrative staff in the Rural Health Management System. Clinic staff have comprehensive administrative control over clinic operations, including patient management, staff management, and visit coordination.

## User Role: Clinic Staff (`clinic_staff`)

### Access Level
- **Portal**: Staff Portal (`/portal/staff/*`)
- **Login Type**: `"staff"`
- **Permissions**: Full administrative control over clinic operations

### Core Responsibilities
1. **Clinic Management**: Update clinic information and settings
2. **Staff Management**: Create, manage, and oversee all staff types
3. **Patient Management**: Register and manage patient records
4. **Visit Coordination**: Schedule and manage patient visits
5. **Reporting**: Access comprehensive clinic statistics

## Workflow Steps

### Phase 1: Initial Setup

#### Step 1.1: Clinic Registration
**Endpoint**: `POST /api/v1/auth/register/clinic`

Creates the initial clinic and clinic staff account.

```json
{
  "email": "admin@mountainclinic.com",
  "password": "SecurePass123!",
  "name": "Mountain Rural Health Clinic",
  "address": "123 Mountain View Road, Highland Village",
  "contact_number": "+1-555-0123",
  "district": "Highland District"
}
```

**Expected Response**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_type": "clinic_staff",
  "user": {
    "id": 1,
    "name": "Mountain Rural Health Clinic",
    "address": "123 Mountain View Road, Highland Village",
    "contact_number": "+1-555-0123",
    "district": "Highland District",
    "created_at": "2025-07-12T10:00:00Z"
  }
}
```

#### Step 1.2: Staff Portal Login
**Endpoint**: `POST /api/v1/auth/clinic-login`

```json
{
  "email": "admin@mountainclinic.com",
  "password": "SecurePass123!",
  "login_type": "staff"
}
```

**Key Validation**:
- `login_type: "staff"` ensures staff portal access
- User type must be `clinic_staff`
- Returns JWT token with clinic context

### Phase 2: Staff Management

#### Step 2.1: Create Primary Doctor
**Endpoint**: `POST /api/v1/portal/staff/staff`
**Authorization**: `Bearer CLINIC_STAFF_TOKEN`

```json
{
  "email": "dr.sarah@mountainclinic.com",
  "password": "DoctorPass123!",
  "full_name": "Dr. Sarah Johnson",
  "role": "Doctor",
  "phone": "+1-555-0124"
}
```

**Process**:
1. Creates `doctor` user account for medical portal access
2. Creates staff record linked to clinic
3. Enables doctor to login with `login_type: "medical"`

#### Step 2.2: Create Nursing Staff
**Endpoint**: `POST /api/v1/portal/staff/staff`
**Authorization**: `Bearer CLINIC_STAFF_TOKEN`

```json
{
  "email": "nurse.mary@mountainclinic.com",
  "password": "NursePass123!",
  "full_name": "Nurse Mary Wilson",
  "role": "Nurse",
  "phone": "+1-555-0125"
}
```

#### Step 2.3: Create Administrative Support
**Endpoint**: `POST /api/v1/portal/staff/staff`
**Authorization**: `Bearer CLINIC_STAFF_TOKEN`

```json
{
  "email": "admin.tom@mountainclinic.com",
  "password": "AdminPass123!",
  "full_name": "Tom Anderson",
  "role": "Clinic_Administrator",
  "phone": "+1-555-0126"
}
```

**Note**: `Clinic_Administrator` and `Pharmacist` roles create staff records but no login accounts.

#### Step 2.4: Create Pharmacy Staff
**Endpoint**: `POST /api/v1/portal/staff/staff`
**Authorization**: `Bearer CLINIC_STAFF_TOKEN`

```json
{
  "email": "pharmacist.lisa@mountainclinic.com",
  "password": "PharmPass123!",
  "full_name": "Lisa Chen",
  "role": "Pharmacist",
  "phone": "+1-555-0127"
}
```

#### Step 2.5: View All Staff
**Endpoint**: `GET /api/v1/portal/staff/staff`
**Authorization**: `Bearer CLINIC_STAFF_TOKEN`

**Query Parameters**:
- `page=1` (default)
- `per_page=10` (default, max 100)
- `search=Sarah` (optional, searches name, role, phone)

**Response Example**:
```json
{
  "data": [
    {
      "id": 1,
      "full_name": "Dr. Sarah Johnson",
      "role": "Doctor",
      "phone": "+1-555-0124",
      "email": "dr.sarah@mountainclinic.com",
      "is_active": true,
      "user": {
        "id": 2,
        "email": "dr.sarah@mountainclinic.com",
        "user_type": "doctor",
        "is_active": true
      }
    }
  ],
  "page": 1,
  "per_page": 10,
  "total": 4,
  "total_pages": 1
}
```

### Phase 3: Patient Management

#### Step 3.1: Register New Patient
**Endpoint**: `POST /api/v1/portal/staff/patients`
**Authorization**: `Bearer CLINIC_STAFF_TOKEN`

```json
{
  "full_name": "John Doe",
  "gender": "Male",
  "date_of_birth": "1985-05-15",
  "address": "456 Valley Road, Highland Village",
  "phone": "+1-555-0200"
}
```

**Automatic Fields**:
- `clinic_id`: Automatically set from staff's clinic
- `created_at`: Set to current timestamp

#### Step 3.2: Register Multiple Patients
Repeat Step 3.1 for additional patients:

```json
{
  "full_name": "Jane Smith",
  "gender": "Female",
  "date_of_birth": "1992-08-22",
  "address": "789 Pine Street, Highland Village",
  "phone": "+1-555-0201"
}
```

```json
{
  "full_name": "Robert Brown",
  "gender": "Male",
  "date_of_birth": "1978-12-03",
  "address": "321 Oak Avenue, Highland Village",
  "phone": "+1-555-0202"
}
```

#### Step 3.3: View All Patients
**Endpoint**: `GET /api/v1/portal/staff/patients`
**Authorization**: `Bearer CLINIC_STAFF_TOKEN`

**Query Parameters**:
- `page=1`
- `per_page=10`
- `search=John` (searches name and phone)

#### Step 3.4: View Patient Details
**Endpoint**: `GET /api/v1/portal/staff/patients/{id}`
**Authorization**: `Bearer CLINIC_STAFF_TOKEN`

**Example**: `GET /api/v1/portal/staff/patients/1`

**Response includes**:
- Patient information
- Visit history
- Associated diagnoses and prescriptions

### Phase 4: Visit Management

#### Step 4.1: Schedule Patient Visit
**Endpoint**: `POST /api/v1/portal/staff/visits`
**Authorization**: `Bearer CLINIC_STAFF_TOKEN`

```json
{
  "patient_id": 1,
  "staff_id": 1,
  "reason": "Annual checkup and blood pressure monitoring",
  "notes": "Patient reports feeling well, requests routine screening"
}
```

**Automatic Fields**:
- `clinic_id`: Set from staff's clinic
- `visit_date`: Set to current timestamp if not specified

#### Step 4.2: Schedule Follow-up Visit
```json
{
  "patient_id": 1,
  "staff_id": 1,
  "visit_date": "2025-07-20T14:00:00Z",
  "reason": "Follow-up for blood pressure medication",
  "notes": "Two weeks after starting new medication"
}
```

#### Step 4.3: View All Visits
**Endpoint**: `GET /api/v1/portal/staff/visits`
**Authorization**: `Bearer CLINIC_STAFF_TOKEN`

**Query Parameters**:
- `page=1`
- `per_page=20`
- `patient_id=1` (filter by patient)
- `staff_id=1` (filter by staff member)

#### Step 4.4: View Visit Details
**Endpoint**: `GET /api/v1/portal/staff/visits/{id}`
**Authorization**: `Bearer CLINIC_STAFF_TOKEN`

**Response includes**:
- Visit information
- Patient details
- Staff member details
- Associated diagnoses and prescriptions

### Phase 5: Clinic Management

#### Step 5.1: View Clinic Profile
**Endpoint**: `GET /api/v1/portal/staff/profile`
**Authorization**: `Bearer CLINIC_STAFF_TOKEN`

#### Step 5.2: Update Clinic Information
**Endpoint**: `PUT /api/v1/portal/staff/profile`
**Authorization**: `Bearer CLINIC_STAFF_TOKEN`

```json
{
  "name": "Mountain Rural Health Clinic & Wellness Center",
  "address": "123 Mountain View Road, Suite 100, Highland Village",
  "contact_number": "+1-555-0123",
  "district": "Highland District"
}
```

#### Step 5.3: View Dashboard Statistics
**Endpoint**: `GET /api/v1/portal/staff/dashboard`
**Authorization**: `Bearer CLINIC_STAFF_TOKEN`

**Response Example**:
```json
{
  "total_patients": 150,
  "total_staff": 8,
  "total_visits": 1250,
  "visits_today": 12,
  "active_doctors": 3,
  "active_nurses": 4
}
```

## Restricted Actions

Clinic staff **cannot** perform the following actions:

### ❌ Medical Record Creation
- Cannot create diagnoses (doctor-only)
- Cannot create prescriptions (doctor-only)

**Example Failed Request**:
```bash
POST /api/v1/portal/staff/diagnoses
```
**Response**: `403 Forbidden - "Insufficient permissions"`

### ❌ Cross-Clinic Access
- Cannot access other clinics' data
- All queries automatically filtered by clinic ID

## Security Considerations

### 1. Permission Validation
Every request validates:
- JWT token authenticity
- User type matches required role
- Clinic ownership for data access
- Specific permission for action

### 2. Data Isolation
- All data queries include clinic_id filter
- Staff can only see their clinic's data
- Cross-clinic access automatically blocked

### 3. Audit Trail
- All administrative actions logged
- Staff creation/modification tracked
- Patient management activities recorded

## Error Scenarios

### Common Error Cases

#### Invalid Token
```json
{
  "error": "Invalid token"
}
```

#### Insufficient Permissions
```json
{
  "error": "Insufficient permissions for this action"
}
```

#### Invalid Staff Role
```json
{
  "error": "Invalid role. Must be Doctor, Nurse, Clinic_Administrator, or Pharmacist"
}
```

#### Email Already Exists
```json
{
  "error": "Email already registered"
}
```

#### Patient Not Found
```json
{
  "error": "Patient not found in this clinic"
}
```

## Best Practices

### 1. Staff Account Management
- Use strong passwords for all staff accounts
- Provide role-appropriate access only
- Regular review of staff access levels

### 2. Patient Data Management
- Verify patient information accuracy
- Maintain updated contact information
- Respect patient privacy guidelines

### 3. Visit Coordination
- Schedule visits with appropriate staff
- Include detailed visit reasons and notes
- Coordinate between administrative and medical staff

### 4. System Security
- Regular password updates
- Monitor for unusual access patterns
- Report security concerns immediately

## Integration with Medical Portal

### Coordination Points

1. **Staff Creation**: Clinic staff create doctor/nurse accounts
2. **Patient Registration**: Patients registered by staff, visited by medical staff
3. **Visit Scheduling**: Staff schedule visits, medical staff conduct examinations
4. **Record Review**: Staff can view medical records but not create them

### Data Flow
```
Clinic Staff → Creates Staff Accounts → Medical Staff Login
     ↓
Registers Patients → Medical Staff Access → Create Medical Records
     ↓
Schedules Visits → Medical Staff Conduct → Generate Diagnoses/Prescriptions
     ↓
Views Reports ← Administrative Overview ← Medical Data Aggregation
```

This workflow ensures efficient clinic operations while maintaining proper separation of administrative and medical responsibilities.
