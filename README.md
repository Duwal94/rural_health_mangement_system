# Rural Health Management System

A comprehensive REST API for managing rural health clinics, built with Go, Fiber, and PostgreSQL.

## Features

- **Complete CRUD Operations** for all entities
- **Pagination** and filtering for list endpoints
- **Proper validation** and error handling
- **Database relationships** with GORM
- **Soft deletes** for data integrity
- **RESTful API design** with versioning
- **Health check** endpoint

## Project Structure

```
├── internal/
│   ├── config/         # Configuration management
│   ├── database/       # Database connection and setup
│   ├── handlers/       # HTTP request handlers
│   └── models/         # Data models and DTOs
├── main.go            # Application entry point
├── go.mod            # Go module dependencies
└── .env              # Environment variables
```

## Entities

- **Clinics** - Health facilities
- **Patients** - Registered patients
- **Staff** - Medical and administrative staff
- **Visits** - Patient visits to clinics
- **Diagnoses** - Medical diagnoses for visits
- **Prescriptions** - Prescribed medications

## API Endpoints

### Base URL
```
http://localhost:3000/api/v1
```

### Health Check
```
GET /health
```

### Clinics
```
GET    /clinics          # List all clinics (paginated)
GET    /clinics/:id      # Get clinic by ID
POST   /clinics          # Create new clinic
PUT    /clinics/:id      # Update clinic
DELETE /clinics/:id      # Delete clinic
```

### Patients
```
GET    /patients         # List all patients (paginated)
GET    /patients/:id     # Get patient by ID
POST   /patients         # Create new patient
PUT    /patients/:id     # Update patient
DELETE /patients/:id     # Delete patient
```

### Staff
```
GET    /staff            # List all staff (paginated)
GET    /staff/:id        # Get staff by ID
POST   /staff            # Create new staff
PUT    /staff/:id        # Update staff
DELETE /staff/:id        # Delete staff
```

### Visits
```
GET    /visits           # List all visits (paginated)
GET    /visits/:id       # Get visit by ID
POST   /visits           # Create new visit
```

### Diagnoses
```
GET    /diagnoses        # List all diagnoses (paginated)
GET    /diagnoses/:id    # Get diagnosis by ID
POST   /diagnoses        # Create new diagnosis
PUT    /diagnoses/:id    # Update diagnosis
DELETE /diagnoses/:id    # Delete diagnosis
```

### Prescriptions
```
GET    /prescriptions    # List all prescriptions (paginated)
GET    /prescriptions/:id # Get prescription by ID
POST   /prescriptions    # Create new prescription
PUT    /prescriptions/:id # Update prescription
DELETE /prescriptions/:id # Delete prescription
```

## Query Parameters

### Pagination
- `page` - Page number (default: 1)
- `per_page` - Items per page (default: 10, max: 100)

### Filtering
- **Clinics**: `search`, `district`
- **Patients**: `search`, `clinic_id`
- **Staff**: `clinic_id`, `role`
- **Visits**: `patient_id`, `clinic_id`
- **Diagnoses**: `visit_id`
- **Prescriptions**: `visit_id`

## Sample Requests

### Create a Clinic
```bash
POST /api/v1/clinics
Content-Type: application/json

{
  "name": "Rural Health Center",
  "address": "123 Main St, Village",
  "contact_number": "+1234567890",
  "district": "Central District"
}
```

### Create a Patient
```bash
POST /api/v1/patients
Content-Type: application/json

{
  "full_name": "John Doe",
  "gender": "Male",
  "date_of_birth": "1990-01-15",
  "address": "456 Village Road",
  "phone": "+1234567891",
  "clinic_id": 1
}
```

### Create a Visit
```bash
POST /api/v1/visits
Content-Type: application/json

{
  "patient_id": 1,
  "clinic_id": 1,
  "staff_id": 1,
  "visit_date": "2024-01-15T10:00:00Z",
  "reason": "Regular checkup",
  "notes": "Patient appears healthy"
}
```

## Error Responses

All errors follow a consistent format:

```json
{
  "error": "Error message",
  "details": "Additional error details (optional)"
}
```

## Pagination Response

List endpoints return paginated results:

```json
{
  "data": [...],
  "page": 1,
  "per_page": 10,
  "total": 100,
  "total_pages": 10
}
```

## Setup

### Prerequisites
- Go 1.21+
- PostgreSQL 12+

### Installation

1. Clone the repository
2. Copy `.env.example` to `.env` and update database credentials
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Run the application:
   ```bash
   go run main.go
   ```

### Environment Variables

```env
DATABASE_URL=host=localhost user=postgres password=postgres dbname=rural_health_db port=5432 sslmode=disable
PORT=3000
ENVIRONMENT=development
```

## Database Schema

The application automatically creates the following tables with proper relationships:
- `clinics`
- `patients` (belongs to clinic)
- `staff` (belongs to clinic)
- `visits` (belongs to patient, clinic, and staff)
- `diagnoses` (belongs to visit)
- `prescriptions` (belongs to visit)

## Validation Rules

### Clinic
- Name: required, 2-255 characters
- Address: required, 5-500 characters
- Contact Number: required, 10-20 characters
- District: required, 2-100 characters

### Patient
- Full Name: required, 2-255 characters
- Gender: required, one of: Male, Female, Other
- Date of Birth: required, valid date
- Address: required, 5-500 characters
- Phone: required, 10-20 characters
- Clinic ID: required, must exist

### Staff
- Full Name: required, 2-255 characters
- Role: required, one of: Doctor, Nurse, Administrator, Pharmacist
- Phone: required, 10-20 characters
- Email: required, valid email, unique
- Clinic ID: required, must exist

## Contributing

1. Follow Go best practices
2. Add proper error handling
3. Include validation for all inputs
4. Write tests for new functionality
5. Update documentation

## License

MIT License
