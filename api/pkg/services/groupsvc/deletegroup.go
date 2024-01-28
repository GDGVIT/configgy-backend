package groupsvc

import (
	"context"
	"errors"

	"github.com/GDGVIT/configgy-backend/api/pkg/api"
	"github.com/GDGVIT/configgy-backend/api/pkg/services/accesscontrolsvc"
	"github.com/GDGVIT/configgy-backend/pkg/tables"
)

func (svc *GroupSvcImpl) GroupDelete(ctx context.Context, groupPID string, userPID string) (api.GenericMessageResponse, int, error) {
	user, err := svc.DB.GetUserByPID(userPID)
	if err != nil {
		return api.GenericMessageResponse{}, 0, nil
	}
	group, err := svc.DB.GetGroupByPID(groupPID)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	hasPermission, err := svc.accesscontrolrsvc.UserHasPermissionToResource(user.ID, group.ID, tables.GroupResource, accesscontrolsvc.DeleteOperation)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	if !hasPermission {
		return api.GenericMessageResponse{}, 0, errors.New("user does not have permission to delete group")
	}
	err = svc.DB.DeleteGroup(group.ID)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	return api.GenericMessageResponse{}, 0, nil
}
