# Dashboard Analytics API Documentation

## Overview

The Dashboard Analytics API provides comprehensive statistical insights for the Rural Health Management System. It offers two main endpoints:

1. **System-wide Analytics** (Public) - Statistics across all clinics
2. **Clinic-specific Analytics** (Protected) - Statistics for a specific clinic

## Routes

### 1. System-wide Dashboard Analytics

**Endpoint**: `GET /api/v1/dashboard/analytics`
**Access**: Public (not protected)
**Description**: Returns comprehensive analytics for all clinics in the system

#### Response Format

```json
{
  "overall_stats": {
    "total_clinics": 25,
    "total_patients": 1250,
    "total_staff": 180,
    "total_visits": 5420,
    "total_diagnoses": 4890,
    "total_prescriptions": 5100,
    "visits_this_month": 312,
    "visits_today": 15
  },
  "top_diagnoses": [
    {
      "diagnosis_code": "J11.1",
      "description": "Influenza due to unidentified influenza virus",
      "count": 340,
      "percentage": 15.2
    },
    {
      "diagnosis_code": "I10",
      "description": "Essential hypertension",
      "count": 285,
      "percentage": 12.8
    }
  ],
  "top_prescriptions": [
    {
      "medication_name": "Acetaminophen",
      "count": 450,
      "percentage": 18.3,
      "avg_duration_days": 7.2,
      "common_dosages": [
        {
          "dosage": "500mg",
          "count": 320
        },
        {
          "dosage": "1000mg",
          "count": 130
        }
      ]
    }
  ],
  "demographics": {
    "age_groups": [
      {
        "age_group": "18-30",
        "count": 325,
        "percentage": 26.0
      },
      {
        "age_group": "31-50",
        "count": 410,
        "percentage": 32.8
      }
    ],
    "gender_distribution": [
      {
        "gender": "Female",
        "count": 680,
        "percentage": 54.4
      },
      {
        "gender": "Male",
        "count": 570,
        "percentage": 45.6
      }
    ]
  },
  "illness_trends": [
    {
      "month": "January",
      "year": 2025,
      "diagnosis_code": "J11.1",
      "count": 45
    }
  ],
  "district_analytics": [
    {
      "district": "Mountain District",
      "total_clinics": 8,
      "total_patients": 420,
      "total_visits": 1850,
      "top_diagnoses": [...],
      "top_prescriptions": [...]
    }
  ],
  "seasonal_trends": [
    {
      "season": "Winter",
      "month": 12,
      "year": 2024,
      "diagnosis_code": "J11.1",
      "count": 85
    }
  ]
}
```

### 2. Clinic-specific Dashboard Analytics

**Endpoint**: `GET /api/v1/portal/staff/dashboard/analytics`
**Access**: Protected (Clinic Staff only)
**Authorization**: `Bearer TOKEN`
**Description**: Returns comprehensive analytics for the authenticated clinic

**Alternative Endpoints** (same functionality):
- `GET /api/v1/portal/medical/dashboard/analytics` (Doctors and Nurses)

#### Headers
```
Authorization: Bearer <clinic_staff_token>
Content-Type: application/json
```

#### Response Format

Same as system-wide analytics but filtered for the specific clinic. The `overall_stats.total_clinics` will always be 1, and `district_analytics` will be empty.

## Analytics Features

### 📊 Overall Statistics
- Total counts for clinics, patients, staff, visits, diagnoses, prescriptions
- Time-based metrics (today, this month)

### 🏥 Top Diagnoses (Top 10)
- Most frequently used diagnosis codes
- Includes ICD-10 codes and descriptions
- Shows count and percentage of total diagnoses

### 💊 Top Prescriptions (Top 10)
- Most frequently prescribed medications
- Average treatment duration
- Common dosage patterns
- Count and percentage statistics

### 👥 Demographics Analysis
- **Age Groups**: Under 18, 18-30, 31-50, 51-70, Over 70
- **Gender Distribution**: Male, Female, Other
- Includes counts and percentages

### 📈 Illness Trends
- Monthly illness patterns over the last 12 months
- Helps identify seasonal health patterns
- Grouped by diagnosis code

### 🌍 District Analytics (System-wide only)
- Health statistics by district
- Compares different geographical areas
- Includes top diagnoses and prescriptions per district

### 🌡️ Seasonal Trends
- Seasonal illness patterns over the last 2 years
- Categorizes by Winter, Spring, Summer, Fall
- Useful for epidemic preparedness

## Use Cases

### For Healthcare Administrators
- **Resource Planning**: Identify most common conditions and medications
- **Seasonal Preparation**: Prepare for seasonal illness outbreaks
- **Staff Allocation**: Understand patient demographics and visit patterns

### For Public Health Officials
- **Disease Surveillance**: Monitor illness trends across districts
- **Policy Making**: Data-driven healthcare policy decisions
- **Resource Distribution**: Allocate resources based on district needs

### For Medical Professionals
- **Clinical Insights**: Understand patient population characteristics
- **Treatment Patterns**: Analyze prescription and treatment trends
- **Research Data**: Access anonymized health statistics

## Example Usage

### cURL Examples

#### Get System-wide Analytics
```bash
curl -X GET "http://localhost:3000/api/v1/dashboard/analytics" \
  -H "Content-Type: application/json"
```

#### Get Clinic-specific Analytics
```bash
curl -X GET "http://localhost:3000/api/v1/portal/staff/dashboard/analytics" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

### JavaScript/Fetch Examples

#### Get System-wide Analytics
```javascript
const response = await fetch('/api/v1/dashboard/analytics');
const analytics = await response.json();
console.log('System Analytics:', analytics);
```

#### Get Clinic-specific Analytics
```javascript
const response = await fetch('/api/v1/portal/staff/dashboard/analytics', {
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
});
const analytics = await response.json();
console.log('Clinic Analytics:', analytics);
```

## Security & Privacy

### Data Protection
- All patient data is anonymized and aggregated
- No personally identifiable information (PII) is exposed
- Statistics are calculated in real-time from the database

### Access Control
- System-wide analytics: Public access (for research and policy making)
- Clinic-specific analytics: Protected by role-based authentication
- Clinic staff can only access their own clinic's data

### Rate Limiting
- Consider implementing rate limiting for public endpoints
- Cache results for better performance

## Performance Considerations

### Database Optimization
- Complex aggregation queries may be slow on large datasets
- Consider implementing caching for frequently accessed data
- Database indexes on date fields recommended

### Recommended Optimizations
1. Add database indexes on `visit_date`, `clinic_id`, `diagnosis_code`
2. Implement Redis caching for analytics results
3. Consider pre-calculating common statistics
4. Use database views for complex aggregations

## Error Handling

### Common Error Responses

#### 401 Unauthorized (for protected endpoints)
```json
{
  "error": "Unauthorized access"
}
```

#### 500 Internal Server Error
```json
{
  "error": "Failed to fetch analytics data"
}
```

#### 404 Not Found
```json
{
  "error": "Clinic not found"
}
```

## Future Enhancements

### Planned Features
- **Real-time Updates**: WebSocket integration for live analytics
- **Custom Date Ranges**: Allow filtering by custom date ranges
- **Export Functionality**: CSV/PDF export of analytics data
- **Comparative Analytics**: Compare clinic performance metrics
- **Predictive Analytics**: Machine learning for health trend prediction

### API Versioning
Current version: `v1`
Future versions will maintain backward compatibility.
