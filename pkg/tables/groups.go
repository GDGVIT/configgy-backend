package tables

import "gorm.io/gorm"

type Groups struct {
	ID          int    `gorm:"column:group_id;primaryKey;autoIncrement"`
	PID         string `gorm:"column:group_pid;unique;type:varchar(100)"`
	Name        string `gorm:"column:group_name;not null"`
	Description string `gorm:"column:group_description;not null"`
}

func (t *Groups) TableName() string {
	return "groups"
}

// Create a new group
func (db *DB) CreateGroup(user *Users, group *Groups, groupMembers map[string]Permission) error {
	txns := []*gorm.DB{}
	if group.PID == "" {
		group.PID = UUIDWithPrefix("group")
	}
	tx := db.gormDB.Create(&group)
	if tx.Error != nil {
		return tx.Error
	}
	txns = append(txns, tx)

	tx = db.gormDB.Create(&PermissionAssignments{
		GroupID:        group.ID,
		UserID:         user.ID,
		PermissionName: OwnerPermission,
		ResourcePID:    group.PID,
		ResourceType:   GroupResource,
		IdentityPID:    user.PID,
		IdentityType:   UserIdentity,
	})
	if tx.Error != nil {
		db.RollbackTxns(txns)
		return tx.Error
	}
	txns = append(txns, tx)

	for memberPID, memberPermission := range groupMembers {
		member, err := db.GetUserByPID(memberPID)
		if err != nil {
			db.RollbackTxns(txns)
			return err
		}
		tx = db.gormDB.Create(&PermissionAssignments{
			GroupID:        group.ID,
			UserID:         member.ID,
			PermissionName: memberPermission,
			ResourcePID:    group.PID,
			ResourceType:   GroupResource,
			IdentityPID:    memberPID,
			IdentityType:   UserIdentity,
		})
		if tx.Error != nil {
			db.RollbackTxns(txns)
			return tx.Error
		}
		txns = append(txns, tx)
	}

	return nil
}

// Get a group by id
func (db *DB) GetGroupByID(id int) (*Groups, error) {
	var group Groups
	err := db.gormDB.Where("group_id = ?", id).First(&group).Error
	return &group, err
}

// Get a group by pid
func (db *DB) GetGroupByPID(pid string) (*Groups, error) {
	var group Groups
	err := db.gormDB.Where("group_pid = ?", pid).First(&group).Error
	return &group, err
}

// Create a new group member
func (db *DB) CreateGroupMember(groupMember *PermissionAssignments) *gorm.DB {
	return db.gormDB.Create(groupMember)
}

// Get a group member by id
func (db *DB) GetGroupMemberByID(id int) (*PermissionAssignments, error) {
	var groupMember PermissionAssignments
	err := db.gormDB.Where("group_member_id = ?", id).First(&groupMember).Error
	return &groupMember, err
}

// Get a group member by pid
func (db *DB) GetGroupMemberByPID(pid string) (*PermissionAssignments, error) {
	var groupMember PermissionAssignments
	err := db.gormDB.Where("group_member_pid = ?", pid).First(&groupMember).Error
	return &groupMember, err
}

// Get all group members for a group
func (db *DB) GetGroupMembersByGroupID(groupID int) ([]*PermissionAssignments, error) {
	var groupMembers []*PermissionAssignments
	err := db.gormDB.Where("group_id = ?", groupID).Find(&groupMembers).Error
	return groupMembers, err
}

// Get all group memberships for a user
func (db *DB) GetGroupsByUserID(userID int) ([]*PermissionAssignments, error) {
	var groups []*PermissionAssignments
	err := db.gormDB.Where("user_id = ?", userID).Find(&groups).Error
	return groups, err
}

// Delete a group
func (db *DB) DeleteGroup(groupID int) error {
	groupMemberDeleteTx := db.gormDB.Where("group_id = ?", groupID).Delete(&PermissionAssignments{})
	if groupMemberDeleteTx.Error != nil {
		groupMemberDeleteTx.Rollback()
		return groupMemberDeleteTx.Error
	}

	groupDeleteTx := db.gormDB.Where("group_id = ?", groupID).Delete(&Groups{})
	if groupDeleteTx.Error != nil {
		groupDeleteTx.Rollback()
		return groupDeleteTx.Error
	}
	return nil
}

// Delete a group member
func (db *DB) DeleteGroupMemberByGroupIDAndUserID(groupID int, userID int) error {
	groupMemberDeleteTx := db.gormDB.Where("group_id = ? AND user_id = ?", groupID, userID).Delete(&PermissionAssignments{})
	if groupMemberDeleteTx.Error != nil {
		groupMemberDeleteTx.Rollback()
		return groupMemberDeleteTx.Error
	}
	return nil
}

func (db *DB) IsUserInGroup(userID int, groupID int) (bool, error) {
	var groupMember PermissionAssignments
	err := db.gormDB.Where("user_id = ? AND group_id = ?", userID, groupID).First(&groupMember).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *DB) UpdateGroup(group *Groups, groupID int) error {
	err := db.gormDB.Model(&Groups{}).Where("group_id = ?", groupID).Updates(group).Error
	if err != nil {
		return err
	}
	return nil
}
