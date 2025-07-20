# Role-Based Access Control Examples

## 1. Clinic Registration (Creates Clinic Staff)

```bash
curl -X POST http://localhost:8080/api/v1/auth/register/clinic \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@ruralclinic.com",
    "password": "securepass123",
    "name": "Rural Health Clinic",
    "address": "123 Village Road, Remote Area",
    "contact_number": "+1234567890",
    "district": "Mountain District"
  }'
```

**Response**: Creates a clinic and a clinic_staff user who can manage the clinic.

## 2. Clinic Staff Login (Administrative Portal)

```bash
curl -X POST http://localhost:8080/api/v1/auth/clinic-login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@ruralclinic.com",
    "password": "securepass123",
    "login_type": "staff"
  }'
```

## 3. Clinic Staff Creates a Doctor

```bash
curl -X POST http://localhost:8080/api/v1/portal/staff/staff \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer CLINIC_STAFF_TOKEN" \
  -d '{
    "email": "doctor@ruralclinic.com",
    "password": "doctorpass123",
    "full_name": "Dr. Sarah Johnson",
    "role": "Doctor",
    "phone": "+1234567891"
  }'
```

## 4. Clinic Staff Creates a Nurse

```bash
curl -X POST http://localhost:8080/api/v1/portal/staff/staff \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer CLINIC_STAFF_TOKEN" \
  -d '{
    "email": "nurse@ruralclinic.com",
    "password": "nursepass123",
    "full_name": "Nurse Mary Wilson",
    "role": "Nurse",
    "phone": "+1234567892"
  }'
```

## 5. Clinic Staff Creates Non-Login Staff (Administrator)

```bash
curl -X POST http://localhost:8080/api/v1/portal/staff/staff \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer CLINIC_STAFF_TOKEN" \
  -d '{
    "email": "admin.support@ruralclinic.com",
    "password": "temppass123",
    "full_name": "Admin Assistant Tom",
    "role": "Clinic_Administrator",
    "phone": "+1234567893"
  }'
```

**Note**: Clinic_Administrator role won't get a login account, only the staff record is created.

## 6. Doctor Login (Medical Portal)

```bash
curl -X POST http://localhost:8080/api/v1/auth/clinic-login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "doctor@ruralclinic.com",
    "password": "doctorpass123",
    "login_type": "medical"
  }'
```

## 7. Nurse Login (Medical Portal)

```bash
curl -X POST http://localhost:8080/api/v1/auth/clinic-login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "nurse@ruralclinic.com",
    "password": "nursepass123",
    "login_type": "medical"
  }'
```

## 8. Clinic Staff Creates a Patient

```bash
curl -X POST http://localhost:8080/api/v1/portal/staff/patients \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer CLINIC_STAFF_TOKEN" \
  -d '{
    "full_name": "John Doe",
    "gender": "Male",
    "date_of_birth": "1985-05-15",
    "address": "456 Rural Street, Village",
    "phone": "+1234567894"
  }'
```

## 9. Doctor Creates a Visit

```bash
curl -X POST http://localhost:8080/api/v1/portal/medical/visits \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer DOCTOR_TOKEN" \
  -d '{
    "patient_id": 1,
    "reason": "Regular checkup and flu symptoms",
    "notes": "Patient reports fever and cough for 3 days"
  }'
```

## 10. Doctor Creates a Diagnosis (Only Doctors Can Do This)

```bash
curl -X POST http://localhost:8080/api/v1/portal/medical/diagnoses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer DOCTOR_TOKEN" \
  -d '{
    "visit_id": 1,
    "diagnosis_code": "J11.1",
    "description": "Influenza due to unidentified influenza virus with other respiratory manifestations"
  }'
```

## 11. Doctor Creates a Prescription (Only Doctors Can Do This)

```bash
curl -X POST http://localhost:8080/api/v1/portal/medical/prescriptions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer DOCTOR_TOKEN" \
  -d '{
    "visit_id": 1,
    "medication_name": "Tamiflu",
    "dosage": "75mg",
    "instructions": "Take twice daily with food",
    "duration_days": 5
  }'
```

## 12. Nurse Attempts to Create Diagnosis (Should Fail)

```bash
curl -X POST http://localhost:8080/api/v1/portal/medical/diagnoses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer NURSE_TOKEN" \
  -d '{
    "visit_id": 1,
    "diagnosis_code": "J11.1",
    "description": "Some diagnosis"
  }'
```

**Expected Response**: 403 Forbidden - "Access denied. Only doctors can perform this action"

## 13. Doctor Attempts to Create Staff (Should Fail)

```bash
curl -X POST http://localhost:8080/api/v1/portal/staff/staff \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer DOCTOR_TOKEN" \
  -d '{
    "email": "newstaff@ruralclinic.com",
    "password": "password123",
    "full_name": "New Staff Member",
    "role": "Nurse",
    "phone": "+1234567895"
  }'
```

**Expected Response**: 403 Forbidden - "Insufficient permissions"

## 14. Clinic Staff Views Dashboard

```bash
curl -X GET http://localhost:8080/api/v1/portal/staff/dashboard \
  -H "Authorization: Bearer CLINIC_STAFF_TOKEN"
```

**Response**: Comprehensive clinic statistics including all patients, staff, visits, etc.

## 15. Doctor Views Medical Dashboard

```bash
curl -X GET http://localhost:8080/api/v1/portal/medical/dashboard \
  -H "Authorization: Bearer DOCTOR_TOKEN"
```

**Response**: Medical-focused statistics including doctor's visits, diagnoses, prescriptions, etc.

## 16. Patient Registration and Login

```bash
# Register patient
curl -X POST http://localhost:8080/api/v1/auth/register/patient \
  -H "Content-Type: application/json" \
  -d '{
    "email": "patient@example.com",
    "password": "patientpass123",
    "full_name": "Jane Patient",
    "gender": "Female",
    "date_of_birth": "1990-03-20",
    "address": "789 Patient Lane",
    "phone": "+1234567896",
    "clinic_id": 1
  }'

# Patient login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "patient@example.com",
    "password": "patientpass123"
  }'
```

## 17. Patient Views Own Visits

```bash
curl -X GET http://localhost:8080/api/v1/portal/patient/visits \
  -H "Authorization: Bearer PATIENT_TOKEN"
```

**Response**: Only visits for this specific patient.

## Access Control Validation Examples

### Example 1: Cross-Clinic Access Prevention
If Doctor A from Clinic 1 tries to access Patient B from Clinic 2:

```bash
curl -X GET http://localhost:8080/api/v1/portal/medical/patients/5 \
  -H "Authorization: Bearer DOCTOR_CLINIC_1_TOKEN"
```

**Expected Response**: 404 Not Found (patient not found in doctor's clinic)

### Example 2: Role Validation
If a doctor tries to use staff portal:

```bash
curl -X GET http://localhost:8080/api/v1/portal/staff/dashboard \
  -H "Authorization: Bearer DOCTOR_TOKEN"
```

**Expected Response**: 403 Forbidden - "Insufficient permissions"

### Example 3: Medical Action Restrictions
If a nurse tries to create a prescription:

```bash
curl -X POST http://localhost:8080/api/v1/portal/medical/prescriptions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer NURSE_TOKEN" \
  -d '{
    "visit_id": 1,
    "medication_name": "Aspirin",
    "dosage": "500mg",
    "instructions": "Take as needed",
    "duration_days": 7
  }'
```

**Expected Response**: 403 Forbidden - "Access denied. Only doctors can perform this action"

## Summary of Role Capabilities

### Clinic Staff Portal (`clinic_staff`)
- ✅ Manage clinic profile
- ✅ Create/manage patients  
- ✅ Create/manage all staff types
- ✅ Create/manage visits
- ✅ View diagnoses and prescriptions
- ❌ Create diagnoses or prescriptions

### Medical Portal - Doctor (`doctor`)
- ✅ View patients (read-only)
- ✅ View staff (read-only)
- ✅ Create/manage visits
- ✅ Create/manage diagnoses
- ✅ Create/manage prescriptions
- ❌ Manage staff or patients
- ❌ Modify clinic settings

### Medical Portal - Nurse (`nurse`)
- ✅ View patients (read-only)
- ✅ View staff (read-only)
- ✅ Create/manage visits
- ✅ View diagnoses and prescriptions
- ❌ Create diagnoses or prescriptions
- ❌ Manage staff or patients
- ❌ Modify clinic settings

### Patient Portal (`patient`)
- ✅ View own profile
- ✅ Update own profile (limited)
- ✅ View own visits
- ✅ View own diagnoses and prescriptions
- ❌ Access any other patient's data
- ❌ Access clinic management functions
