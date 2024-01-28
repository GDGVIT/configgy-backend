package tables

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Permission string

const (
	OwnerPermission Permission = "owner"
	AdminPermission Permission = "admin"
	EditPermission  Permission = "edit"
	ViewPermission  Permission = "view"
)

type IdentityType string

const (
	UserIdentity  IdentityType = "user"
	GroupIdentity IdentityType = "group"
)

type ResourceTypes string

const (
	VaultResource      ResourceTypes = "vault"
	CredentialResource ResourceTypes = "credential"
	GroupResource      ResourceTypes = "group"
)

type PermissionAssignments struct {
	ID             int           `gorm:"column:permission_assignment_id;primaryKey;autoIncrement"`
	PID            string        `gorm:"column:permission_assignment_pid;unique;type:varchar(40)"`
	PermissionName Permission    `gorm:"column:permission_name;not null"`
	VaultID        int           `gorm:"column:vault_id"`
	ResourcePID    string        `gorm:"column:resource_pid;type:varchar(40)"`
	ResourceType   ResourceTypes `gorm:"column:resource_type"`
	IdentityPID    string        `gorm:"column:identity_pid;type:varchar(40)"`
	IdentityType   IdentityType  `gorm:"column:identity_type"`
	CredentialID   int           `gorm:"column:credential_id"`
	GroupID        int           `gorm:"column:group_id"`
	UserID         int           `gorm:"column:user_id;"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (t *PermissionAssignments) TableName() string {
	return "permission_assignments"
}

// Create a new vault
func (db *DB) CreatePermissionAssignment(permission_assignment *PermissionAssignments) *gorm.DB {
	return db.gormDB.Create(permission_assignment)
}

func (db *DB) DeletePermissionAssignment(permission_assignment *PermissionAssignments) *gorm.DB {
	return db.gormDB.Delete(&permission_assignment)
}

// Get a permission assignment by id
func (db *DB) GetPermissionAssignmentByID(id int) (*PermissionAssignments, error) {
	var permission_assignment PermissionAssignments
	err := db.gormDB.Where("permission_assignment_id = ?", id).First(&permission_assignment).Error
	return &permission_assignment, err
}

// Get a permission assignment by pid
func (db *DB) GetPermissionAssignmentByPID(pid string) (*PermissionAssignments, error) {
	var permission_assignment PermissionAssignments
	err := db.gormDB.Where("permission_assignment_pid = ?", pid).First(&permission_assignment).Error
	return &permission_assignment, err
}

func (db *DB) GetPermissionsForUserOnResource(userID int, resourceID int, resourceType ResourceTypes) (*PermissionAssignments, error) {
	var permission_assignment PermissionAssignments
	err := db.gormDB.Where(fmt.Sprintf("user_id = ? AND %s_id = ?", string(resourceType)), userID, resourceID).First(&permission_assignment).Error
	return &permission_assignment, err
}

func (db *DB) GetPermissionsForGroupOnResource(groupID int, resourceID int, resourceType ResourceTypes) (*PermissionAssignments, error) {
	var permission_assignment PermissionAssignments
	err := db.gormDB.Where(fmt.Sprintf("group_id = ? AND %s_id = ?", string(resourceType)), groupID, resourceID).First(&permission_assignment).Error
	return &permission_assignment, err
}

func (db *DB) GetPermissionsForResource(resourceID int, resourceType ResourceTypes) ([]*PermissionAssignments, error) {
	var permission_assignments []*PermissionAssignments
	err := db.gormDB.Where(fmt.Sprintf("%s_id = ?", string(resourceType)), resourceID).Find(&permission_assignments).Error
	return permission_assignments, err
}

func (db *DB) GetPermissionsForUser(userID int) ([]*PermissionAssignments, error) {
	var permission_assignments []*PermissionAssignments
	err := db.gormDB.Where("user_id = ?", userID).Find(&permission_assignments).Error
	return permission_assignments, err
}

func (db *DB) GetPermissionsForGroup(groupID int) ([]*PermissionAssignments, error) {
	var permission_assignments []*PermissionAssignments
	err := db.gormDB.Where("group_id = ?", groupID).Find(&permission_assignments).Error
	return permission_assignments, err
}
