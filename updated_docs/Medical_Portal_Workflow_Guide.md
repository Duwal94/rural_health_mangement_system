# Medical Portal Workflow Guide

## Overview

This guide details the workflows for medical personnel (doctors and nurses) using the Medical Portal of the Rural Health Management System. The Medical Portal provides access to patient care functions while maintaining appropriate role-based restrictions.

## User Roles

### Doctor (`doctor`)
- **Portal**: Medical Portal (`/portal/medical/*`)
- **Login Type**: `"medical"`
- **Key Capabilities**: Full medical practice including diagnosis and prescription
- **Restrictions**: Cannot manage staff or modify clinic settings

### Nurse (`nurse`)
- **Portal**: Medical Portal (`/portal/medical/*`)
- **Login Type**: `"medical"`
- **Key Capabilities**: Medical support and patient care
- **Restrictions**: Cannot create diagnoses or prescriptions

## Access Control Matrix

| Action | Doctor | Nurse |
|--------|--------|-------|
| View Patients | ✅ Read-only | ✅ Read-only |
| View Staff | ✅ Read-only | ✅ Read-only |
| Create Visits | ✅ | ✅ |
| Update Visits | ✅ | ✅ |
| View Visits | ✅ | ✅ |
| Create Diagnoses | ✅ | ❌ |
| View Diagnoses | ✅ | ✅ |
| Create Prescriptions | ✅ | ❌ |
| View Prescriptions | ✅ | ✅ |
| Manage Patients | ❌ | ❌ |
| Manage Staff | ❌ | ❌ |

## Doctor Workflow

### Phase 1: Authentication and Setup

#### Step 1.1: Doctor Login
**Endpoint**: `POST /api/v1/auth/clinic-login`

```json
{
  "email": "dr.sarah@mountainclinic.com",
  "password": "DoctorPass123!",
  "login_type": "medical"
}
```

**Response**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_type": "doctor",
  "user": {
    "id": 1,
    "full_name": "Dr. Sarah Johnson",
    "role": "Doctor",
    "phone": "+1-555-0124",
    "email": "dr.sarah@mountainclinic.com",
    "clinic": {
      "id": 1,
      "name": "Mountain Rural Health Clinic"
    }
  }
}
```

#### Step 1.2: View Medical Dashboard
**Endpoint**: `GET /api/v1/portal/medical/dashboard`
**Authorization**: `Bearer DOCTOR_TOKEN`

**Response Example**:
```json
{
  "my_visits_today": 5,
  "my_total_visits": 287,
  "total_patients": 150,
  "total_staff": 8,
  "my_diagnoses": 245,
  "my_prescriptions": 312
}
```

### Phase 2: Patient Care

#### Step 2.1: View Patient List
**Endpoint**: `GET /api/v1/portal/medical/patients`
**Authorization**: `Bearer DOCTOR_TOKEN`

**Query Parameters**:
- `page=1`
- `per_page=20`
- `search=John` (search by name or phone)

**Response Example**:
```json
{
  "data": [
    {
      "id": 1,
      "full_name": "John Doe",
      "gender": "Male",
      "date_of_birth": "1985-05-15",
      "phone": "+1-555-0200",
      "created_at": "2025-07-01T10:00:00Z"
    }
  ],
  "page": 1,
  "per_page": 20,
  "total": 150,
  "total_pages": 8
}
```

#### Step 2.2: View Patient Details
**Endpoint**: `GET /api/v1/portal/medical/patients/{id}`
**Authorization**: `Bearer DOCTOR_TOKEN`

**Example**: `GET /api/v1/portal/medical/patients/1`

**Response includes**:
- Complete patient information
- Medical history (visits, diagnoses, prescriptions)
- Recent visits and outcomes

### Phase 3: Visit Management

#### Step 3.1: Create Patient Visit
**Endpoint**: `POST /api/v1/portal/medical/visits`
**Authorization**: `Bearer DOCTOR_TOKEN`

```json
{
  "patient_id": 1,
  "reason": "Patient presents with chest pain and shortness of breath",
  "notes": "Symptoms started 2 hours ago, mild exertion triggers discomfort"
}
```

**Automatic Fields**:
- `clinic_id`: Set from doctor's clinic
- `staff_id`: Set from doctor's staff ID
- `visit_date`: Set to current timestamp

#### Step 3.2: View Own Visits
**Endpoint**: `GET /api/v1/portal/medical/visits`
**Authorization**: `Bearer DOCTOR_TOKEN`

**Query Parameters**:
- `page=1`
- `per_page=10`
- `patient_id=1` (filter by specific patient)
- `show_all=false` (default: shows only doctor's visits)
- `show_all=true` (shows all clinic visits)

#### Step 3.3: View Specific Visit
**Endpoint**: `GET /api/v1/portal/medical/visits/{id}`
**Authorization**: `Bearer DOCTOR_TOKEN`

**Response includes**:
- Visit details
- Patient information
- Existing diagnoses and prescriptions

### Phase 4: Medical Diagnosis (Doctors Only)

#### Step 4.1: Create Diagnosis
**Endpoint**: `POST /api/v1/portal/medical/diagnoses`
**Authorization**: `Bearer DOCTOR_TOKEN`

```json
{
  "visit_id": 1,
  "diagnosis_code": "I20.9",
  "description": "Angina pectoris, unspecified. Chest pain likely cardiac in origin, requires further evaluation with ECG and cardiac enzymes."
}
```

**Validation**:
- Only doctors can create diagnoses
- Visit must exist and belong to doctor's clinic
- Diagnosis code should follow standard coding (ICD-10)

#### Step 4.2: Create Additional Diagnosis
```json
{
  "visit_id": 1,
  "diagnosis_code": "Z51.11",
  "description": "Encounter for antineoplastic chemotherapy - follow-up care plan discussed"
}
```

### Phase 5: Prescription Management (Doctors Only)

#### Step 5.1: Create Prescription
**Endpoint**: `POST /api/v1/portal/medical/prescriptions`
**Authorization**: `Bearer DOCTOR_TOKEN`

```json
{
  "visit_id": 1,
  "medication_name": "Nitroglycerin",
  "dosage": "0.4mg sublingual",
  "instructions": "Place under tongue as needed for chest pain, may repeat every 5 minutes up to 3 doses. If pain persists, seek emergency care.",
  "duration_days": 30
}
```

#### Step 5.2: Create Additional Prescriptions
```json
{
  "visit_id": 1,
  "medication_name": "Aspirin",
  "dosage": "81mg",
  "instructions": "Take once daily with food for cardiovascular protection",
  "duration_days": 90
}
```

```json
{
  "visit_id": 1,
  "medication_name": "Metoprolol",
  "dosage": "25mg",
  "instructions": "Take twice daily, monitor blood pressure and heart rate",
  "duration_days": 30
}
```

## Nurse Workflow

### Phase 1: Authentication and Setup

#### Step 1.1: Nurse Login
**Endpoint**: `POST /api/v1/auth/clinic-login`

```json
{
  "email": "nurse.mary@mountainclinic.com",
  "password": "NursePass123!",
  "login_type": "medical"
}
```

#### Step 1.2: View Nursing Dashboard
**Endpoint**: `GET /api/v1/portal/medical/dashboard`
**Authorization**: `Bearer NURSE_TOKEN`

**Response focuses on nursing metrics**:
```json
{
  "my_visits_today": 8,
  "my_total_visits": 456,
  "total_patients": 150,
  "total_staff": 8,
  "my_diagnoses": 0,
  "my_prescriptions": 0
}
```

### Phase 2: Patient Support

#### Step 2.1: View Patients (Read-Only)
**Endpoint**: `GET /api/v1/portal/medical/patients`
**Authorization**: `Bearer NURSE_TOKEN`

Same endpoint as doctors, but nurses have read-only access.

#### Step 2.2: Support Patient Care
**Endpoint**: `POST /api/v1/portal/medical/visits`
**Authorization**: `Bearer NURSE_TOKEN`

```json
{
  "patient_id": 1,
  "reason": "Vital signs monitoring and medication administration",
  "notes": "Blood pressure: 145/90, Heart rate: 78 bpm, Temperature: 98.6°F, administered prescribed medications"
}
```

### Phase 3: Visit Documentation

#### Step 3.1: Document Nursing Care
```json
{
  "patient_id": 2,
  "reason": "Post-procedure monitoring and patient education",
  "notes": "Patient stable post-procedure, vital signs within normal limits, provided discharge instructions and follow-up care guidelines"
}
```

#### Step 3.2: View Visit History
Nurses can view all visits to understand patient care continuity:

**Endpoint**: `GET /api/v1/portal/medical/visits`
**Authorization**: `Bearer NURSE_TOKEN`

### Phase 4: Restricted Actions

#### Step 4.1: Attempted Diagnosis Creation (Should Fail)
**Endpoint**: `POST /api/v1/portal/medical/diagnoses`
**Authorization**: `Bearer NURSE_TOKEN`

```json
{
  "visit_id": 1,
  "diagnosis_code": "I20.9",
  "description": "Some diagnosis"
}
```

**Expected Response**: `403 Forbidden`
```json
{
  "error": "Access denied. Only doctors can perform this action"
}
```

#### Step 4.2: Attempted Prescription Creation (Should Fail)
**Endpoint**: `POST /api/v1/portal/medical/prescriptions`
**Authorization**: `Bearer NURSE_TOKEN`

**Expected Response**: `403 Forbidden`
```json
{
  "error": "Access denied. Only doctors can perform this action"
}
```

## Collaborative Workflows

### Scenario 1: Doctor-Nurse Team Care

#### Step 1: Nurse Initial Assessment
```json
{
  "patient_id": 3,
  "reason": "Initial nursing assessment for new admission",
  "notes": "Patient admitted with complaints of fatigue and weight loss. Vital signs stable, baseline measurements recorded."
}
```

#### Step 2: Doctor Examination
```json
{
  "patient_id": 3,
  "reason": "Medical evaluation following nursing assessment",
  "notes": "Physical examination reveals enlarged lymph nodes, recommending lab work and imaging studies."
}
```

#### Step 3: Doctor Diagnosis
```json
{
  "visit_id": 5,
  "diagnosis_code": "R59.9",
  "description": "Enlarged lymph nodes, unspecified. Requires further evaluation to rule out malignancy."
}
```

#### Step 4: Doctor Prescription
```json
{
  "visit_id": 5,
  "medication_name": "Multivitamin",
  "dosage": "1 tablet",
  "instructions": "Take once daily with breakfast while awaiting test results",
  "duration_days": 14
}
```

#### Step 5: Nurse Follow-up
```json
{
  "patient_id": 3,
  "reason": "Patient education and follow-up care coordination",
  "notes": "Discussed test procedures with patient, scheduled lab appointments, provided contact information for questions."
}
```

### Scenario 2: Emergency Response

#### Step 1: Nurse Triage
```json
{
  "patient_id": 4,
  "reason": "Emergency triage - chest pain presentation",
  "notes": "Patient presents with severe chest pain 8/10, diaphoretic, vital signs: BP 160/95, HR 110, O2 sat 95%"
}
```

#### Step 2: Doctor Emergency Assessment
```json
{
  "patient_id": 4,
  "reason": "Emergency medical assessment - acute chest pain",
  "notes": "EKG shows ST elevation in leads II, III, aVF. Immediate cardiac intervention required."
}
```

#### Step 3: Doctor Emergency Diagnosis
```json
{
  "visit_id": 7,
  "diagnosis_code": "I21.19",
  "description": "ST elevation myocardial infarction involving inferior wall. Immediate transfer to cardiac center required."
}
```

#### Step 4: Doctor Emergency Prescriptions
```json
{
  "visit_id": 7,
  "medication_name": "Aspirin",
  "dosage": "325mg",
  "instructions": "Chew immediately, given for acute MI",
  "duration_days": 1
}
```

## Staff Information Access

### View Clinic Staff
**Endpoint**: `GET /api/v1/portal/medical/staff`
**Authorization**: `Bearer DOCTOR_TOKEN` or `Bearer NURSE_TOKEN`

**Purpose**: View colleague information for collaboration

**Response Example**:
```json
{
  "data": [
    {
      "id": 1,
      "full_name": "Dr. Sarah Johnson",
      "role": "Doctor",
      "phone": "+1-555-0124",
      "is_active": true
    },
    {
      "id": 2,
      "full_name": "Nurse Mary Wilson",
      "role": "Nurse",
      "phone": "+1-555-0125",
      "is_active": true
    }
  ]
}
```

## Error Handling

### Medical Portal Specific Errors

#### Insufficient Medical Permissions
```json
{
  "error": "Access denied. Only doctors can perform this action"
}
```

#### Visit Not Found
```json
{
  "error": "Visit not found in this clinic"
}
```

#### Invalid Patient Access
```json
{
  "error": "Patient not found in this clinic"
}
```

#### Cross-Portal Access Attempt
```json
{
  "error": "Insufficient permissions"
}
```

## Security Features

### 1. Role-Based Restrictions
- Nurses cannot create medical records
- Doctors cannot manage staff
- All users limited to clinic data

### 2. Medical Record Integrity
- Only qualified doctors create diagnoses
- Prescription creation restricted to doctors
- Clear audit trail for medical decisions

### 3. Data Protection
- Patient data only accessible within clinic
- Medical records protected by role validation
- Automatic clinic ownership verification

## Best Practices

### For Doctors
1. **Complete Documentation**: Include detailed diagnosis descriptions
2. **Accurate Coding**: Use proper medical coding standards
3. **Clear Prescriptions**: Provide detailed medication instructions
4. **Team Communication**: Coordinate with nursing staff

### For Nurses
1. **Detailed Observations**: Document comprehensive patient assessments
2. **Timely Updates**: Record patient status changes promptly
3. **Care Coordination**: Collaborate effectively with doctors
4. **Patient Education**: Document patient teaching activities

### For Both
1. **Security Awareness**: Protect login credentials
2. **Patient Privacy**: Respect patient confidentiality
3. **System Integrity**: Report technical issues promptly
4. **Professional Standards**: Maintain appropriate medical documentation

This Medical Portal workflow ensures effective patient care while maintaining proper role-based access control and medical record integrity.
