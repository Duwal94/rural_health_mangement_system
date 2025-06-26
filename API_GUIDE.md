# Rural Health Management System API

## Quick Start

1. **Setup Database:**
   ```bash
   # Create PostgreSQL database
   createdb rural_health_db
   ```

2. **Run the application:**
   ```bash
   go run .
   ```

3. **Seed sample data:**
   ```bash
   go run cmd/seed/main.go
   ```

## API Testing Examples

### 1. Create a Clinic
```bash
curl -X POST http://localhost:3000/api/v1/clinics \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Rural Health Center",
    "address": "123 Main St, Village",
    "contact_number": "+1234567890",
    "district": "Central District"
  }'
```

### 2. Create a Patient
```bash
curl -X POST http://localhost:3000/api/v1/patients \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "gender": "Male",
    "date_of_birth": "1990-01-15",
    "address": "456 Village Road",
    "phone": "+1234567891",
    "clinic_id": 1
  }'
```

### 3. Create a Visit
```bash
curl -X POST http://localhost:3000/api/v1/visits \
  -H "Content-Type: application/json" \
  -d '{
    "patient_id": 1,
    "clinic_id": 1,
    "staff_id": 1,
    "visit_date": "2024-01-15T10:00:00Z",
    "reason": "Regular checkup",
    "notes": "Patient appears healthy"
  }'
```

### 4. Get All Patients (with pagination)
```bash
curl "http://localhost:3000/api/v1/patients?page=1&per_page=10&search=John"
```

### 5. Health Check
```bash
curl http://localhost:3000/health
```

## Environment Variables

Copy `.env.example` to `.env` and update:

```
DATABASE_URL=host=localhost user=postgres password=postgres dbname=rural_health_db port=5432 sslmode=disable
PORT=3000
ENVIRONMENT=development
```
