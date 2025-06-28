# Rural Health Management System - Workflow Overview

## Introduction

The Rural Health Management System provides two primary user workflows: **Clinic Workflow** and **Patient Workflow**. These workflows demonstrate the complete lifecycle of rural healthcare management, from clinic setup to patient care delivery and medical record access.

## Workflow Relationship

```
┌─────────────────┐    ┌─────────────────┐
│ Clinic Workflow │    │ Patient Workflow│
│                 │    │                 │
│ 1. Register     │    │ 1. Register     │
│ 2. Add Staff    │────┤ 2. Update Profile│
│ 3. Add Patients │    │ 3. View Records │
│ 4. Create Visits│────┤ 4. Access Data  │
│ 5. Add Diagnoses│    │                 │
│ 6. Add Rx       │────┤                 │
└─────────────────┘    └─────────────────┘
```

## Postman Collections

### Files Included

1. **`Clinic_Workflow.postman_collection.json`**
   - 19 step-by-step requests
   - Complete clinic onboarding and operations
   - Automated variable management
   - Error handling and validation

2. **`Patient_Workflow.postman_collection.json`**
   - 19 step-by-step requests
   - Patient registration and self-service
   - Medical record access and management
   - Account security and profile updates

### Documentation

1. **`Clinic_Workflow_Guide.md`**
   - Detailed explanation of each clinic workflow step
   - Expected outcomes and error handling
   - Integration points with patient workflow

2. **`Patient_Workflow_Guide.md`**
   - Complete patient journey documentation
   - Security features and data access
   - Dependencies on clinic workflow

## Recommended Execution Order

### Option 1: Complete Healthcare System Demo
```
1. Run Clinic Workflow (Steps 1-19)
   └── Creates: Clinic, Staff, Patients, Visits, Diagnoses, Prescriptions

2. Run Patient Workflow (Steps 1-19)
   └── Demonstrates: Patient access to medical records created above
```

### Option 2: Independent Testing
```
1. Run Patient Workflow (Steps 1-12)
   └── Tests: Registration, profile management, security

2. Run Clinic Workflow (Steps 1-19)
   └── Tests: Complete clinic operations

3. Re-run Patient Workflow (Steps 13-19)
   └── Tests: Access to medical records
```

## Key Features Demonstrated

### Clinic Workflow Features
- **Clinic Registration**: Email/password authentication with JWT tokens
- **Staff Management**: Add doctors, nurses, administrators, pharmacists
- **Patient Registration**: Register patients associated with clinic
- **Visit Documentation**: Create visit records with staff assignments
- **Medical Coding**: ICD-10 diagnosis codes and descriptions
- **Prescription Management**: Detailed medication prescribing
- **Dashboard Analytics**: Real-time statistics and metrics

### Patient Workflow Features
- **Patient Registration**: Self-registration with clinic selection
- **Profile Management**: Update personal information
- **Account Security**: Password changes and authentication
- **Medical Record Access**: View complete health history
- **Visit Details**: Access visit notes and staff information
- **Diagnosis Review**: View all medical diagnoses with codes
- **Prescription Access**: Complete medication history and instructions

## Technical Implementation

### Authentication System
- **JWT Tokens**: Secure authentication with role-based access
- **Automatic Variables**: Postman scripts save tokens for subsequent requests
- **Token Refresh**: Login demonstrates token renewal
- **Role Isolation**: Clinic and patient data access separation

### Data Relationships
```
Clinic
├── Staff (1:many)
├── Patients (1:many)
└── Visits (through patients)
    ├── Diagnoses (1:many)
    └── Prescriptions (1:many)

Patient
├── Visits (1:many)
├── Diagnoses (through visits)
└── Prescriptions (through visits)
```

### API Endpoints Covered

#### Authentication Endpoints
- `POST /auth/register/clinic` - Clinic registration
- `POST /auth/register/patient` - Patient registration
- `POST /auth/login` - User authentication
- `GET /auth/profile` - User profile access
- `POST /auth/change-password` - Password updates

#### Clinic Portal Endpoints
- `GET /portal/clinic/profile` - Clinic profile
- `PUT /portal/clinic/profile` - Update clinic info
- `GET /portal/clinic/dashboard` - Statistics dashboard
- `GET /portal/clinic/patients` - Patient management
- `GET /portal/clinic/patients/{id}` - Patient details
- `POST /portal/clinic/staff` - Staff management
- `GET /portal/clinic/staff` - Staff listings
- `POST /portal/clinic/visits` - Visit creation
- `GET /portal/clinic/visits` - Visit management
- `POST /portal/clinic/diagnoses` - Diagnosis documentation
- `POST /portal/clinic/prescriptions` - Prescription management

#### Patient Portal Endpoints
- `GET /portal/patient/profile` - Patient profile
- `PUT /portal/patient/profile` - Profile updates
- `GET /portal/patient/visits` - Visit history
- `GET /portal/patient/visits/{id}` - Visit details
- `GET /portal/patient/diagnoses` - Diagnosis history
- `GET /portal/patient/prescriptions` - Prescription history

#### System Endpoints
- `GET /health` - System health check
- `GET /clinics` - Available clinics (for patient registration)

## Automation Features

### Postman Automation
- **Variable Management**: Automatic saving of IDs and tokens
- **Test Scripts**: Validation and data extraction
- **Error Handling**: Proper response validation
- **Console Logging**: Progress tracking and debugging

### Example Automation
```javascript
// Automatic token saving
if (pm.response.code === 201) {
    const responseJson = pm.response.json();
    pm.collectionVariables.set('clinicToken', responseJson.token);
    pm.collectionVariables.set('clinicId', responseJson.user.id);
    console.log('Clinic registered successfully!');
}
```

## Use Cases Demonstrated

### Rural Healthcare Scenarios

#### Clinic Setup
1. **New Rural Clinic**: Establishment of healthcare facility
2. **Staff Onboarding**: Adding medical professionals
3. **Patient Base Building**: Community registration
4. **Operational Readiness**: System preparation for patient care

#### Patient Care Delivery
1. **Routine Checkups**: Annual health examinations
2. **Acute Care**: Flu and illness treatment
3. **Multi-Provider Care**: Doctor and nurse collaboration
4. **Complete Documentation**: Comprehensive medical records

#### Patient Self-Service
1. **Health Record Access**: Personal medical history
2. **Prescription Tracking**: Medication management
3. **Provider Information**: Healthcare team details
4. **Secure Access**: Protected health information

## Security Features

### Access Control
- **Role-Based Authentication**: Separate clinic and patient access
- **Data Isolation**: Users only see their own data
- **Token-Based Security**: JWT authentication throughout
- **Password Management**: Secure password change functionality

### Privacy Protection
- **Patient Data Protection**: HIPAA-style access controls
- **Clinic Data Isolation**: Multi-tenant security
- **Audit Trails**: Complete request logging
- **Secure Endpoints**: All medical data requires authentication

## Performance Features

### Pagination
- **Large Dataset Handling**: Paginated responses for visits, diagnoses, prescriptions
- **Configurable Page Sizes**: Adjustable per_page parameters
- **Total Count Information**: Complete dataset statistics

### Efficient Queries
- **Filtered Results**: Search and filter capabilities
- **Related Data Loading**: Visits include diagnoses and prescriptions
- **Optimized Responses**: Only necessary data returned

## Error Handling

### Common Error Scenarios
- **Authentication Failures**: Invalid credentials or expired tokens
- **Authorization Errors**: Accessing unauthorized data
- **Validation Errors**: Invalid data formats or missing fields
- **Resource Not Found**: Invalid IDs or non-existent records

### Testing Error Conditions
Both workflows include scenarios that can test:
- Token expiration and renewal
- Invalid data submissions
- Unauthorized access attempts
- Missing resource references

## Getting Started

### Prerequisites
1. **API Server**: Rural Health Management System running on localhost:3000
2. **Postman**: Latest version with collection import capability
3. **Test Environment**: Clean database for demonstration

### Quick Start
1. **Import Collections**: Load both JSON files into Postman
2. **Set Base URL**: Ensure `{{baseUrl}}` points to your API server
3. **Run Clinic Workflow**: Execute all 19 steps in sequence
4. **Run Patient Workflow**: Execute all 19 steps to see patient perspective
5. **Review Results**: Check console logs and response data

### Verification
After running both workflows, you should have:
- 1 registered clinic with authentication
- 2 staff members (doctor and nurse)
- 3 patients (2 from clinic workflow, 1 from patient workflow)
- 2 documented visits with complete medical records
- 2 diagnoses and 2 prescriptions
- Updated dashboard statistics

## Troubleshooting

### Common Issues
1. **API Not Running**: Ensure server is started on correct port
2. **Database Issues**: Check database connectivity and schema
3. **Token Errors**: Verify JWT configuration and secret keys
4. **Variable Issues**: Check Postman variable scope and values

### Debugging Tips
- Check Postman console for automatic logging
- Verify response status codes match expectations
- Review API server logs for detailed error information
- Test individual requests before running full workflows

## Extension Opportunities

### Additional Workflows
- **Admin Workflow**: System-wide management and oversight
- **Multi-Clinic Scenarios**: Multiple clinics with patient transfers
- **Advanced Medical Scenarios**: Complex diagnoses and treatment plans

### Integration Testing
- **Load Testing**: Multiple concurrent users
- **Performance Testing**: Large datasets and pagination
- **Security Testing**: Penetration testing and vulnerability assessment

## Conclusion

These workflows provide a comprehensive demonstration of the Rural Health Management System's capabilities, showing both the clinic operations perspective and the patient experience. They serve as functional tests, integration demonstrations, and documentation of the system's healthcare management features.
