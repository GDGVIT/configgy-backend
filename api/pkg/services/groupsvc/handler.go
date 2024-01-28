package groupsvc

import (
	"github.com/GDGVIT/configgy-backend/api/pkg/services/accesscontrolsvc"
	"github.com/GDGVIT/configgy-backend/pkg/logger"
	"github.com/GDGVIT/configgy-backend/pkg/tables"
)

type GroupSvcImpl struct {
	DB                DB
	logger            logger.Logger
	accesscontrolrsvc AccessControllrSvc
	AuthnSvc          AuthnSvc
}

type DB interface {
	GetUserByPID(pid string) (*tables.Users, error)
	GetGroupByPID(pid string) (*tables.Groups, error)

	CreateGroup(user *tables.Users, group *tables.Groups, groupMembers map[string]tables.Permission) error

	UpdateGroup(group *tables.Groups, groupID int) error

	DeleteGroup(groupID int) error
}
type AccessControllrSvc interface {
	UserHasPermissionToResource(userID int, resourceID int, resourceType tables.ResourceTypes, operation accesscontrolsvc.CRUDOperation) (bool, error)
}
type AuthnSvc interface{}

func Handler(db DB, logger logger.Logger, accesscontrolrsvc AccessControllrSvc, AuthnSvc AuthnSvc) *GroupSvcImpl {
	return &GroupSvcImpl{DB: db, logger: logger, accesscontrolrsvc: accesscontrolrsvc, AuthnSvc: AuthnSvc}
}
