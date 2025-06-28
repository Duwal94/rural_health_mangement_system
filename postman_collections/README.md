# Postman Collections for Rural Health Management System

This directory contains workflow-specific Postman collections that demonstrate the complete functionality of the Rural Health Management System.

## Collections

### 🏥 Clinic_Workflow.postman_collection.json
**Complete clinic operations workflow**
- 19 step-by-step requests
- Clinic registration and setup
- Staff management (doctors, nurses)
- Patient registration and management
- Visit documentation and medical records
- Diagnosis and prescription management
- Dashboard analytics

### 👤 Patient_Workflow.postman_collection.json
**Complete patient experience workflow**
- 19 step-by-step requests
- Patient registration and authentication
- Profile management and security
- Medical record access (visits, diagnoses, prescriptions)
- Account management and password changes

## Quick Start

1. **Import Collections**
   ```
   - Open Postman
   - Click Import
   - Drag and drop both JSON files
   ```

2. **Set Environment**
   ```
   Base URL: http://localhost:3000/api/v1
   ```

3. **Run Workflows**
   ```
   1. Execute Clinic Workflow (all 19 steps)
   2. Execute Patient Workflow (all 19 steps)
   ```

## Features

✅ **Automated Variable Management** - IDs and tokens saved automatically  
✅ **Complete Error Handling** - Proper validation and logging  
✅ **Step-by-Step Documentation** - Each request fully explained  
✅ **Real Healthcare Scenarios** - Authentic medical workflows  
✅ **Security Testing** - Authentication and authorization  

## Documentation

- [`../docs/Clinic_Workflow_Guide.md`](../docs/Clinic_Workflow_Guide.md) - Detailed clinic workflow guide
- [`../docs/Patient_Workflow_Guide.md`](../docs/Patient_Workflow_Guide.md) - Complete patient workflow guide  
- [`../docs/Workflow_Overview.md`](../docs/Workflow_Overview.md) - System overview and integration guide

## Prerequisites

- Rural Health Management System API running on `localhost:3000`
- Postman (latest version)
- Clean database for testing

## Expected Results

After running both workflows:
- 1 registered clinic with 2 staff members
- 3 registered patients 
- 2 documented medical visits
- 2 diagnoses and 2 prescriptions
- Complete audit trail and dashboard metrics

## Support

For issues or questions:
1. Check the detailed workflow guides in `/docs`
2. Verify API server is running
3. Review Postman console logs
4. Check API server logs for errors
