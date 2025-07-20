# Robust Role-Based Access Control System

## Overview

The Rural Health Management System now implements a comprehensive role-based access control (RBAC) system that provides granular permissions and separate login portals for different types of users.

## User Types and Roles

### 1. Patient (`patient`)
- **Login Portal**: Patient Portal (`/portal/patient/*`)
- **Permissions**: Can only access their own data
- **Capabilities**:
  - View own profile
  - Update own profile (limited fields)
  - View own visits
  - View own diagnoses and prescriptions

### 2. Clinic Staff (`clinic_staff`)
- **Login Portal**: Staff Portal (`/portal/staff/*`)
- **Login Type**: `staff` in clinic login
- **Permissions**: Full administrative control over clinic operations
- **Capabilities**:
  - Create/manage patients
  - Create/manage all types of staff (Doctor, Nurse, Clinic_Administrator, Pharmacist)
  - Create/manage visits
  - View diagnoses and prescriptions (cannot create them)
  - Update clinic profile
  - View comprehensive dashboard stats

### 3. Doctor (`doctor`)
- **Login Portal**: Medical Portal (`/portal/medical/*`)
- **Login Type**: `medical` in clinic login
- **Permissions**: Medical practice focused with diagnostic capabilities
- **Capabilities**:
  - View patients (read-only)
  - View staff (read-only)
  - Create/manage visits
  - **Create/manage diagnoses** (exclusive to doctors)
  - **Create/manage prescriptions** (exclusive to doctors)
  - View medical-focused dashboard stats

### 4. Nurse (`nurse`)
- **Login Portal**: Medical Portal (`/portal/medical/*`)
- **Login Type**: `medical` in clinic login
- **Permissions**: Medical support without diagnostic capabilities
- **Capabilities**:
  - View patients (read-only)
  - View staff (read-only)
  - Create/manage visits
  - View diagnoses and prescriptions (cannot create them)
  - View medical-focused dashboard stats

### 5. Non-Login Roles
These roles are managed by clinic staff but don't have login accounts:
- **Clinic_Administrator**: Administrative support staff
- **Pharmacist**: Medication management staff

## Access Control Matrix

| Action | Patient | Clinic Staff | Doctor | Nurse |
|--------|---------|--------------|--------|-------|
| Create Patients | ❌ | ✅ | ❌ | ❌ |
| View Patients | Own Only | ✅ | ✅ | ✅ |
| Update Patients | Own Only | ✅ | ❌ | ❌ |
| Delete Patients | ❌ | ✅ | ❌ | ❌ |
| Create Staff | ❌ | ✅ | ❌ | ❌ |
| View Staff | ❌ | ✅ | ✅ | ✅ |
| Update Staff | ❌ | ✅ | ❌ | ❌ |
| Delete Staff | ❌ | ✅ | ❌ | ❌ |
| Create Visits | ❌ | ✅ | ✅ | ✅ |
| View Visits | Own Only | ✅ | ✅ | ✅ |
| Update Visits | ❌ | ✅ | ✅ | ✅ |
| Create Diagnoses | ❌ | ❌ | ✅ | ❌ |
| View Diagnoses | Own Only | ✅ | ✅ | ✅ |
| Create Prescriptions | ❌ | ❌ | ✅ | ❌ |
| View Prescriptions | Own Only | ✅ | ✅ | ✅ |
| Manage Clinic | ❌ | ✅ | ❌ | ❌ |

## Login System

### Standard Login (`POST /auth/login`)
- Used for patients and backward compatibility
- Accepts email and password

### Clinic Login (`POST /auth/clinic-login`)
- **NEW**: Specialized login for clinic users
- Accepts email, password, and `login_type` field
- **login_type** options:
  - `"staff"`: For clinic administrative staff
  - `"medical"`: For doctors and nurses

#### Example Clinic Login Request:
```json
{
  "email": "doctor@clinic.com",
  "password": "password123",
  "login_type": "medical"
}
```

## Permission System

### Granular Permissions
The system uses a permission-based approach with the following permissions:

#### Patient Management
- `create_patient`
- `update_patient`
- `view_patient`
- `delete_patient`

#### Staff Management
- `create_staff`
- `update_staff`
- `view_staff`
- `delete_staff`

#### Visit Management
- `create_visit`
- `update_visit`
- `view_visit`
- `delete_visit`

#### Medical Permissions
- `create_diagnosis`
- `update_diagnosis`
- `view_diagnosis`
- `delete_diagnosis`
- `create_prescription`
- `update_prescription`
- `view_prescription`
- `delete_prescription`

#### Administrative Permissions
- `manage_clinic`
- `view_reports`
- `manage_inventory`

### Middleware Implementation
The system includes several middleware functions for access control:

1. **`RequirePermission(permission)`**: Checks if user has specific permission
2. **`RequireMultiplePermissions(...permissions)`**: Checks if user has at least one of the permissions
3. **`RequireDoctorAccess()`**: Restricts access to doctors only
4. **`RequireClinicStaffAccess()`**: Restricts access to clinic staff only
5. **`ValidateClinicOwnership()`**: Ensures users only access their own clinic's data

## API Endpoints

### Staff Portal (`/portal/staff/*`)
**Access**: Clinic Staff only
- `GET /profile` - Get clinic profile
- `PUT /profile` - Update clinic profile
- `GET /dashboard` - Staff dashboard stats
- `POST /patients` - Create patient
- `GET /patients` - List patients
- `GET /patients/:id` - Get patient details
- `POST /staff` - Create staff member
- `GET /staff` - List staff
- `POST /visits` - Create visit
- `GET /visits` - List visits
- `GET /visits/:id` - Get visit details

### Medical Portal (`/portal/medical/*`)
**Access**: Doctors and Nurses only
- `GET /profile` - Get staff profile
- `PUT /profile` - Update profile (limited)
- `GET /dashboard` - Medical dashboard stats
- `GET /patients` - List patients (read-only)
- `GET /patients/:id` - Get patient details (read-only)
- `GET /staff` - List staff (read-only)
- `POST /visits` - Create visit
- `GET /visits` - List visits
- `GET /visits/:id` - Get visit details
- `POST /diagnoses` - Create diagnosis (doctors only)
- `POST /prescriptions` - Create prescription (doctors only)

## Security Features

### 1. Separation of Concerns
- **Administrative tasks**: Handled by clinic staff through staff portal
- **Medical tasks**: Handled by medical professionals through medical portal
- **Patient access**: Limited to patient's own data

### 2. Role Validation
- Login type validation ensures users access appropriate portals
- Permission validation at endpoint level
- Clinic ownership validation prevents cross-clinic access

### 3. Staff Creation Control
- Only clinic staff can create new staff members
- Doctor and nurse accounts can only be created by clinic staff
- Non-login roles (Administrator, Pharmacist) are managed without user accounts

### 4. Medical Action Restrictions
- Only doctors can create diagnoses and prescriptions
- Nurses can view but not create medical records
- Clear separation between administrative and medical responsibilities

## Edge Cases Handled

### 1. **Cross-Clinic Access Prevention**
- Users can only access data from their associated clinic
- Clinic ownership validation middleware prevents data leakage

### 2. **Role Escalation Prevention**
- Users cannot access functions beyond their role permissions
- Multiple validation layers prevent privilege escalation

### 3. **Invalid Login Type Handling**
- Clinic login validates user type against requested login type
- Prevents doctors from accessing staff portal and vice versa

### 4. **Orphaned Staff Prevention**
- Staff creation requires clinic association
- User account creation is tied to staff record creation

### 5. **Medical Record Integrity**
- Only qualified medical professionals can create medical records
- Clear audit trail of who created what medical records

### 6. **Token Security**
- JWT tokens include role and clinic information
- Token validation includes permission checks

### 7. **Account Deactivation**
- Staff can be deactivated without deleting historical records
- Deactivated users cannot login but data remains intact

## Migration Notes

### From Old System
- `clinic_admin` user type is now `clinic_staff`
- Old clinic portal routes still work but are deprecated
- New specialized portals provide better separation of concerns

### Backward Compatibility
- Existing login endpoint continues to work
- Old clinic portal endpoints are maintained
- Gradual migration path available

## Best Practices

### 1. **Use Appropriate Portals**
- Clinic staff should use `/portal/staff/*` endpoints
- Medical staff should use `/portal/medical/*` endpoints
- Patients should use `/portal/patient/*` endpoints

### 2. **Role-Specific Operations**
- Use clinic login with appropriate login_type
- Respect permission boundaries
- Don't attempt cross-role operations

### 3. **Data Access Patterns**
- Always include clinic_id validation
- Use pagination for large datasets
- Implement proper error handling

This robust role system ensures that each user type has appropriate access to system functions while maintaining security, data integrity, and clear separation of responsibilities in the healthcare management workflow.
