#!/bin/bash

# Test script for new dashboard routes
# Run this after starting the server

echo "Testing Dashboard Analytics Routes"
echo "================================="

# Base URL
BASE_URL="http://localhost:8080/api/v1"

echo ""
echo "1. Testing Public Dashboard Analytics (All Clinics)"
echo "Route: GET $BASE_URL/dashboard/analytics"
curl -s -w "\nStatus: %{http_code}\n" "$BASE_URL/dashboard/analytics" | head -20

echo ""
echo "2. Testing Public Dashboard Content (Alternative route)"
echo "Route: GET $BASE_URL/dashboard/content"
curl -s -w "\nStatus: %{http_code}\n" "$BASE_URL/dashboard/content" | head -20

echo ""
echo "3. Testing Protected Clinic Dashboard (Requires Authentication)"
echo "Route: GET $BASE_URL/portal/clinic/dashboard/content"
echo "Note: This will return 401 without valid JWT token"
curl -s -w "\nStatus: %{http_code}\n" "$BASE_URL/portal/clinic/dashboard/content"

echo ""
echo "4. Testing Protected Staff Dashboard (Requires Authentication)"
echo "Route: GET $BASE_URL/portal/staff/dashboard/content"
echo "Note: This will return 401 without valid JWT token"
curl -s -w "\nStatus: %{http_code}\n" "$BASE_URL/portal/staff/dashboard/content"

echo ""
echo "5. Testing Protected Medical Dashboard (Requires Authentication)"
echo "Route: GET $BASE_URL/portal/medical/dashboard/content"
echo "Note: This will return 401 without valid JWT token"
curl -s -w "\nStatus: %{http_code}\n" "$BASE_URL/portal/medical/dashboard/content"

echo ""
echo "Testing completed!"
echo ""
echo "To test protected routes with authentication:"
echo "1. First login to get a JWT token:"
echo "   curl -X POST $BASE_URL/auth/login -H 'Content-Type: application/json' -d '{\"email\":\"staff@clinic.com\", \"password\":\"password\"}'"
echo ""
echo "2. Then use the token in subsequent requests:"
echo "   curl -H 'Authorization: Bearer YOUR_TOKEN' $BASE_URL/portal/staff/dashboard/content"
