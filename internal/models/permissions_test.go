package models

import (
	"testing"
)

func TestHasPermission(t *testing.T) {
	tests := []struct {
		name       string
		userType   string
		staffRole  *string
		permission Permission
		expected   bool
	}{
		{
			name:       "Clinic staff can create patients",
			userType:   "clinic_staff",
			staffRole:  nil,
			permission: PermissionCreatePatient,
			expected:   true,
		},
		{
			name:       "Doctor can create diagnosis",
			userType:   "doctor",
			staffRole:  nil,
			permission: PermissionCreateDiagnosis,
			expected:   true,
		},
		{
			name:       "Nurse cannot create diagnosis",
			userType:   "nurse",
			staffRole:  nil,
			permission: PermissionCreateDiagnosis,
			expected:   false,
		},
		{
			name:       "Doctor cannot create staff",
			userType:   "doctor",
			staffRole:  nil,
			permission: PermissionCreateStaff,
			expected:   false,
		},
		{
			name:       "Nurse can view patients",
			userType:   "nurse",
			staffRole:  nil,
			permission: PermissionViewPatient,
			expected:   true,
		},
		{
			name:       "Doctor can create prescriptions",
			userType:   "doctor",
			staffRole:  nil,
			permission: PermissionCreatePrescription,
			expected:   true,
		},
		{
			name:       "Nurse cannot create prescriptions",
			userType:   "nurse",
			staffRole:  nil,
			permission: PermissionCreatePrescription,
			expected:   false,
		},
		{
			name:       "Clinic staff can manage clinic",
			userType:   "clinic_staff",
			staffRole:  nil,
			permission: PermissionManageClinic,
			expected:   true,
		},
		{
			name:       "Doctor cannot manage clinic",
			userType:   "doctor",
			staffRole:  nil,
			permission: PermissionManageClinic,
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasPermission(tt.userType, tt.staffRole, tt.permission)
			if result != tt.expected {
				t.Errorf("HasPermission(%s, %v, %s) = %v, want %v",
					tt.userType, tt.staffRole, tt.permission, result, tt.expected)
			}
		})
	}
}

func TestRolePermissions(t *testing.T) {
	// Test that all roles have their expected permissions

	// Clinic staff should have administrative permissions
	clinicStaffPerms := RolePermissions["clinic_staff"]
	if !contains(clinicStaffPerms, PermissionCreatePatient) {
		t.Error("Clinic staff should have create patient permission")
	}
	if !contains(clinicStaffPerms, PermissionCreateStaff) {
		t.Error("Clinic staff should have create staff permission")
	}
	if contains(clinicStaffPerms, PermissionCreateDiagnosis) {
		t.Error("Clinic staff should NOT have create diagnosis permission")
	}

	// Doctor should have medical permissions
	doctorPerms := RolePermissions["doctor"]
	if !contains(doctorPerms, PermissionCreateDiagnosis) {
		t.Error("Doctor should have create diagnosis permission")
	}
	if !contains(doctorPerms, PermissionCreatePrescription) {
		t.Error("Doctor should have create prescription permission")
	}
	if contains(doctorPerms, PermissionCreateStaff) {
		t.Error("Doctor should NOT have create staff permission")
	}

	// Nurse should have limited permissions
	nursePerms := RolePermissions["nurse"]
	if !contains(nursePerms, PermissionViewPatient) {
		t.Error("Nurse should have view patient permission")
	}
	if contains(nursePerms, PermissionCreateDiagnosis) {
		t.Error("Nurse should NOT have create diagnosis permission")
	}
	if contains(nursePerms, PermissionCreatePrescription) {
		t.Error("Nurse should NOT have create prescription permission")
	}
}

func contains(perms []Permission, target Permission) bool {
	for _, p := range perms {
		if p == target {
			return true
		}
	}
	return false
}
