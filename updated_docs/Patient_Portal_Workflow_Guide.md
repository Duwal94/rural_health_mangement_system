# Patient Portal Workflow Guide

## Overview

This guide provides a comprehensive walkthrough for patients using the Patient Portal of the Rural Health Management System. The Patient Portal offers secure access to personal health information and healthcare services while maintaining appropriate privacy controls.

## Patient Role

### Patient (`patient`)
- **Portal**: Patient Portal (`/portal/patient/*`)
- **Login Type**: `"patient"`
- **Key Capabilities**: View personal medical records, manage profile, schedule visits
- **Restrictions**: Cannot access other patients' records or administrative functions

## Access Control Matrix

| Action | Patient |
|--------|---------|
| View Own Profile | ✅ |
| Update Own Profile | ✅ |
| View Own Visits | ✅ |
| View Own Diagnoses | ✅ |
| View Own Prescriptions | ✅ |
| Request New Visit | ✅ |
| View Clinic Information | ✅ |
| Access Other Patients | ❌ |
| View Staff Details | ❌ |
| Administrative Functions | ❌ |

## Patient Workflow

### Phase 1: Authentication and Setup

#### Step 1.1: Patient Registration
**Note**: Patients are registered by clinic staff through the Staff Portal

**Staff Registration Process** (for reference):
```json
{
  "full_name": "Alice Johnson",
  "gender": "Female",
  "date_of_birth": "1990-08-15",
  "phone": "+1-555-0300",
  "email": "alice.johnson@email.com",
  "address": "789 Oak Street, Rural Town, ST 12345",
  "emergency_contact_name": "Bob Johnson",
  "emergency_contact_phone": "+1-555-0301",
  "password": "PatientPass123!"
}
```

#### Step 1.2: Patient First-Time Login
**Endpoint**: `POST /api/v1/auth/patient-login`

```json
{
  "email": "alice.johnson@email.com",
  "password": "PatientPass123!"
}
```

**Response**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_type": "patient",
  "patient": {
    "id": 1,
    "full_name": "Alice Johnson",
    "gender": "Female",
    "date_of_birth": "1990-08-15",
    "phone": "+1-555-0300",
    "email": "alice.johnson@email.com",
    "clinic": {
      "id": 1,
      "name": "Mountain Rural Health Clinic"
    }
  }
}
```

#### Step 1.3: View Patient Dashboard
**Endpoint**: `GET /api/v1/portal/patient/dashboard`
**Authorization**: `Bearer PATIENT_TOKEN`

**Response Example**:
```json
{
  "welcome_message": "Welcome back, Alice!",
  "total_visits": 12,
  "recent_visits": 3,
  "upcoming_appointments": 1,
  "active_prescriptions": 2,
  "last_visit_date": "2024-12-15T14:30:00Z",
  "clinic_info": {
    "name": "Mountain Rural Health Clinic",
    "phone": "+1-555-0100",
    "address": "123 Main Street, Rural Town, ST 12345"
  }
}
```

### Phase 2: Profile Management

#### Step 2.1: View Personal Profile
**Endpoint**: `GET /api/v1/portal/patient/profile`
**Authorization**: `Bearer PATIENT_TOKEN`

**Response includes**:
```json
{
  "id": 1,
  "full_name": "Alice Johnson",
  "gender": "Female",
  "date_of_birth": "1990-08-15",
  "phone": "+1-555-0300",
  "email": "alice.johnson@email.com",
  "address": "789 Oak Street, Rural Town, ST 12345",
  "emergency_contact_name": "Bob Johnson",
  "emergency_contact_phone": "+1-555-0301",
  "is_active": true,
  "created_at": "2024-07-01T10:00:00Z",
  "updated_at": "2024-12-01T15:30:00Z"
}
```

#### Step 2.2: Update Personal Information
**Endpoint**: `PUT /api/v1/portal/patient/profile`
**Authorization**: `Bearer PATIENT_TOKEN`

```json
{
  "phone": "+1-555-0302",
  "email": "alice.j.updated@email.com",
  "address": "789 Oak Street, Apt 2B, Rural Town, ST 12345",
  "emergency_contact_name": "Bob Johnson",
  "emergency_contact_phone": "+1-555-0301"
}
```

**Response**:
```json
{
  "message": "Profile updated successfully",
  "patient": {
    "id": 1,
    "full_name": "Alice Johnson",
    "phone": "+1-555-0302",
    "email": "alice.j.updated@email.com",
    "updated_at": "2024-12-20T10:15:00Z"
  }
}
```

### Phase 3: Visit History

#### Step 3.1: View Visit History
**Endpoint**: `GET /api/v1/portal/patient/visits`
**Authorization**: `Bearer PATIENT_TOKEN`

**Query Parameters**:
- `page=1`
- `per_page=10`
- `year=2024` (filter by year)
- `from_date=2024-01-01`
- `to_date=2024-12-31`

**Response Example**:
```json
{
  "data": [
    {
      "id": 1,
      "visit_date": "2024-12-15T14:30:00Z",
      "reason": "Annual physical examination",
      "notes": "Routine check-up, all vitals normal",
      "staff_name": "Dr. Sarah Johnson",
      "staff_role": "Doctor",
      "diagnoses": [
        {
          "id": 1,
          "diagnosis_code": "Z00.00",
          "description": "Encounter for general adult medical examination without abnormal findings"
        }
      ],
      "prescriptions": [
        {
          "id": 1,
          "medication_name": "Multivitamin",
          "dosage": "1 tablet",
          "instructions": "Take once daily with breakfast"
        }
      ]
    },
    {
      "id": 2,
      "visit_date": "2024-11-20T09:15:00Z",
      "reason": "Follow-up for hypertension",
      "notes": "Blood pressure controlled, continue current medication",
      "staff_name": "Dr. Michael Brown",
      "staff_role": "Doctor",
      "diagnoses": [
        {
          "id": 2,
          "diagnosis_code": "I10",
          "description": "Essential hypertension, well controlled"
        }
      ],
      "prescriptions": [
        {
          "id": 2,
          "medication_name": "Lisinopril",
          "dosage": "10mg",
          "instructions": "Take once daily in the morning"
        }
      ]
    }
  ],
  "page": 1,
  "per_page": 10,
  "total": 12,
  "total_pages": 2
}
```

#### Step 3.2: View Specific Visit Details
**Endpoint**: `GET /api/v1/portal/patient/visits/{id}`
**Authorization**: `Bearer PATIENT_TOKEN`

**Example**: `GET /api/v1/portal/patient/visits/1`

**Response includes**:
- Complete visit information
- All diagnoses for that visit
- All prescriptions from that visit
- Staff member who conducted the visit

### Phase 4: Medical Records

#### Step 4.1: View All Diagnoses
**Endpoint**: `GET /api/v1/portal/patient/diagnoses`
**Authorization**: `Bearer PATIENT_TOKEN`

**Query Parameters**:
- `page=1`
- `per_page=20`
- `visit_id=1` (filter by specific visit)

**Response Example**:
```json
{
  "data": [
    {
      "id": 1,
      "diagnosis_code": "Z00.00",
      "description": "Encounter for general adult medical examination without abnormal findings",
      "visit_date": "2024-12-15T14:30:00Z",
      "visit_reason": "Annual physical examination",
      "doctor_name": "Dr. Sarah Johnson"
    },
    {
      "id": 2,
      "diagnosis_code": "I10",
      "description": "Essential hypertension, well controlled",
      "visit_date": "2024-11-20T09:15:00Z",
      "visit_reason": "Follow-up for hypertension",
      "doctor_name": "Dr. Michael Brown"
    }
  ],
  "page": 1,
  "per_page": 20,
  "total": 8,
  "total_pages": 1
}
```

#### Step 4.2: View All Prescriptions
**Endpoint**: `GET /api/v1/portal/patient/prescriptions`
**Authorization**: `Bearer PATIENT_TOKEN`

**Query Parameters**:
- `page=1`
- `per_page=20`
- `active_only=true` (show only current prescriptions)

**Response Example**:
```json
{
  "data": [
    {
      "id": 1,
      "medication_name": "Multivitamin",
      "dosage": "1 tablet",
      "instructions": "Take once daily with breakfast",
      "duration_days": 365,
      "prescribed_date": "2024-12-15T14:30:00Z",
      "expires_date": "2025-12-15T14:30:00Z",
      "is_active": true,
      "visit_date": "2024-12-15T14:30:00Z",
      "doctor_name": "Dr. Sarah Johnson"
    },
    {
      "id": 2,
      "medication_name": "Lisinopril",
      "dosage": "10mg",
      "instructions": "Take once daily in the morning",
      "duration_days": 90,
      "prescribed_date": "2024-11-20T09:15:00Z",
      "expires_date": "2025-02-18T09:15:00Z",
      "is_active": true,
      "visit_date": "2024-11-20T09:15:00Z",
      "doctor_name": "Dr. Michael Brown"
    }
  ],
  "page": 1,
  "per_page": 20,
  "total": 2,
  "total_pages": 1
}
```

### Phase 5: Visit Requests

#### Step 5.1: Request New Visit
**Endpoint**: `POST /api/v1/portal/patient/visit-requests`
**Authorization**: `Bearer PATIENT_TOKEN`

```json
{
  "reason": "Experiencing persistent headaches for the past week",
  "preferred_date": "2024-12-25",
  "preferred_time": "morning",
  "symptoms": "Severe headaches, sensitivity to light, some nausea",
  "urgency_level": "moderate"
}
```

**Response**:
```json
{
  "message": "Visit request submitted successfully",
  "request": {
    "id": 1,
    "reason": "Experiencing persistent headaches for the past week",
    "preferred_date": "2024-12-25",
    "preferred_time": "morning",
    "status": "pending",
    "submitted_at": "2024-12-20T10:30:00Z"
  }
}
```

#### Step 5.2: View Visit Request Status
**Endpoint**: `GET /api/v1/portal/patient/visit-requests`
**Authorization**: `Bearer PATIENT_TOKEN`

**Response Example**:
```json
{
  "data": [
    {
      "id": 1,
      "reason": "Experiencing persistent headaches for the past week",
      "preferred_date": "2024-12-25",
      "preferred_time": "morning",
      "status": "scheduled",
      "scheduled_date": "2024-12-26T09:00:00Z",
      "submitted_at": "2024-12-20T10:30:00Z",
      "updated_at": "2024-12-20T14:00:00Z"
    },
    {
      "id": 2,
      "reason": "Routine blood pressure check",
      "preferred_date": "2024-12-30",
      "preferred_time": "afternoon",
      "status": "pending",
      "submitted_at": "2024-12-20T11:00:00Z"
    }
  ]
}
```

### Phase 6: Health Summaries

#### Step 6.1: Get Health Summary
**Endpoint**: `GET /api/v1/portal/patient/health-summary`
**Authorization**: `Bearer PATIENT_TOKEN`

**Response Example**:
```json
{
  "patient_info": {
    "full_name": "Alice Johnson",
    "age": 34,
    "gender": "Female"
  },
  "recent_vitals": {
    "last_visit": "2024-12-15T14:30:00Z",
    "blood_pressure": "120/80",
    "heart_rate": "72 bpm",
    "weight": "145 lbs",
    "height": "5'6\""
  },
  "current_conditions": [
    {
      "condition": "Essential hypertension",
      "diagnosed_date": "2024-11-20",
      "status": "well controlled"
    }
  ],
  "active_medications": [
    {
      "medication": "Lisinopril 10mg",
      "instructions": "Take once daily in the morning",
      "prescribed_date": "2024-11-20",
      "expires_date": "2025-02-18"
    },
    {
      "medication": "Multivitamin",
      "instructions": "Take once daily with breakfast",
      "prescribed_date": "2024-12-15",
      "expires_date": "2025-12-15"
    }
  ],
  "upcoming_appointments": [
    {
      "date": "2024-12-26T09:00:00Z",
      "reason": "Follow-up for headaches",
      "doctor": "Dr. Sarah Johnson"
    }
  ]
}
```

### Phase 7: Communication

#### Step 7.1: Contact Clinic
**Endpoint**: `GET /api/v1/portal/patient/clinic-info`
**Authorization**: `Bearer PATIENT_TOKEN`

**Response Example**:
```json
{
  "clinic": {
    "name": "Mountain Rural Health Clinic",
    "phone": "+1-555-0100",
    "email": "info@mountainclinic.com",
    "address": "123 Main Street, Rural Town, ST 12345",
    "hours": {
      "monday": "8:00 AM - 6:00 PM",
      "tuesday": "8:00 AM - 6:00 PM",
      "wednesday": "8:00 AM - 6:00 PM",
      "thursday": "8:00 AM - 6:00 PM",
      "friday": "8:00 AM - 5:00 PM",
      "saturday": "9:00 AM - 1:00 PM",
      "sunday": "Closed"
    },
    "emergency_contact": "+1-555-0199"
  }
}
```

#### Step 7.2: Send Message to Clinic
**Endpoint**: `POST /api/v1/portal/patient/messages`
**Authorization**: `Bearer PATIENT_TOKEN`

```json
{
  "subject": "Question about medication",
  "message": "I'm experiencing some side effects from my Lisinopril. Can I schedule a consultation to discuss alternatives?",
  "priority": "normal"
}
```

**Response**:
```json
{
  "message": "Message sent successfully",
  "message_id": 1,
  "sent_at": "2024-12-20T11:45:00Z"
}
```

## Common Patient Scenarios

### Scenario 1: New Patient First Visit

#### Step 1: Staff Registration (Completed by clinic)
```json
{
  "full_name": "Robert Smith",
  "gender": "Male",
  "date_of_birth": "1975-03-10",
  "phone": "+1-555-0400",
  "email": "robert.smith@email.com",
  "address": "456 Pine Avenue, Rural Town, ST 12345",
  "emergency_contact_name": "Susan Smith",
  "emergency_contact_phone": "+1-555-0401",
  "password": "RobertPass123!"
}
```

#### Step 2: Patient First Login
```json
{
  "email": "robert.smith@email.com",
  "password": "RobertPass123!"
}
```

#### Step 3: Profile Review and Update
```json
{
  "phone": "+1-555-0402",
  "address": "456 Pine Avenue, Unit 3, Rural Town, ST 12345"
}
```

#### Step 4: First Visit Request
```json
{
  "reason": "New patient initial consultation",
  "preferred_date": "2024-12-28",
  "preferred_time": "morning",
  "symptoms": "General health assessment needed",
  "urgency_level": "routine"
}
```

### Scenario 2: Chronic Condition Management

#### Step 1: Check Current Medications
**Endpoint**: `GET /api/v1/portal/patient/prescriptions?active_only=true`

#### Step 2: Review Recent Visits
**Endpoint**: `GET /api/v1/portal/patient/visits?from_date=2024-10-01`

#### Step 3: Request Follow-up
```json
{
  "reason": "Hypertension follow-up and medication review",
  "preferred_date": "2025-01-15",
  "preferred_time": "afternoon",
  "symptoms": "Blood pressure seems higher lately",
  "urgency_level": "moderate"
}
```

### Scenario 3: Emergency Health Concern

#### Step 1: Urgent Visit Request
```json
{
  "reason": "Severe chest pain and shortness of breath",
  "preferred_date": "2024-12-20",
  "preferred_time": "immediate",
  "symptoms": "Sharp chest pain 8/10, difficulty breathing, sweating",
  "urgency_level": "urgent"
}
```

#### Step 2: Contact Clinic Directly
**Use emergency contact**: `+1-555-0199`

#### Step 3: Follow-up on Emergency Visit
**Endpoint**: `GET /api/v1/portal/patient/visits/{emergency_visit_id}`

## Error Handling

### Patient Portal Specific Errors

#### Invalid Patient Access
```json
{
  "error": "Patient record not found"
}
```

#### Cross-Patient Data Access
```json
{
  "error": "Access denied. You can only view your own medical records"
}
```

#### Invalid Visit Request
```json
{
  "error": "Invalid preferred date. Date must be in the future"
}
```

#### Profile Update Validation Error
```json
{
  "error": "Invalid phone number format"
}
```

## Security Features

### 1. Patient Data Protection
- Patients can only access their own records
- No access to other patients' information
- No administrative functions available

### 2. Medical Record Security
- Read-only access to diagnoses and prescriptions
- Complete visit history with doctor information
- Secure authentication for all access

### 3. Privacy Controls
- Personal information protection
- Secure communication with clinic
- Audit trail for all patient access

## Mobile-Friendly Features

### 1. Responsive Design
- Works on all device sizes
- Touch-friendly interface
- Optimized for mobile viewing

### 2. Quick Actions
- Easy profile updates
- Simple visit requests
- Fast prescription viewing

### 3. Emergency Access
- Quick clinic contact information
- Emergency number prominently displayed
- Urgent visit request option

## Best Practices for Patients

### 1. Account Security
- Use strong passwords
- Log out when finished
- Don't share login credentials
- Report suspicious activity

### 2. Health Information Management
- Keep profile information current
- Review visit summaries regularly
- Track medication compliance
- Note any medication side effects

### 3. Communication
- Use clear, specific language in visit requests
- Provide accurate symptom descriptions
- Update emergency contact information
- Respond promptly to clinic communications

### 4. Appointment Management
- Schedule routine visits in advance
- Provide accurate availability
- Cancel if unable to attend
- Arrive prepared with questions

## Technical Support

### Common Issues and Solutions

#### 1. Login Problems
- Verify email and password
- Check internet connection
- Contact clinic if locked out

#### 2. Missing Information
- Recent visits may take time to appear
- Contact clinic for clarification
- Check all tabs and sections

#### 3. Technical Errors
- Refresh the page
- Clear browser cache
- Try different browser
- Contact clinic IT support

This Patient Portal workflow ensures patients have secure, convenient access to their healthcare information while maintaining appropriate privacy and security controls.
