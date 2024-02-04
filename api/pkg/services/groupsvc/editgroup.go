package groupsvc

import (
	"context"
	"errors"

	"github.com/GDGVIT/configgy-backend/api/pkg/api"
	"github.com/GDGVIT/configgy-backend/api/pkg/services/accesscontrolsvc"
	"github.com/GDGVIT/configgy-backend/pkg/tables"
)

func (svc *GroupSvcImpl) GroupEdit(ctx context.Context, groupPID string, req api.GroupEditRequest, userPID string) (api.GenericMessageResponse, int, error) {
	user, err := svc.DB.GetUserByPID(userPID)
	if err != nil {
		return api.GenericMessageResponse{}, 0, nil
	}
	group, err := svc.DB.GetGroupByPID(groupPID)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	hasPermission, err := svc.accesscontrolrsvc.UserHasPermissionToResource(user.ID, group.ID, tables.GroupResource, accesscontrolsvc.UpdateOperation)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	if !hasPermission {
		return api.GenericMessageResponse{}, 0, errors.New("user does not have permission to delete group")
	}
	updatedGroup := tables.Groups{}
	if req.Name != nil {
		updatedGroup.Name = *req.Name
	}
	if req.Description != nil {
		updatedGroup.Name = *req.Name
	}
	err = svc.DB.UpdateGroup(&updatedGroup, group.ID)
	if err != nil {
		return api.GenericMessageResponse{}, 0, nil
	}
	return api.GenericMessageResponse{}, 0, nil
}

func (svc *GroupSvcImpl) GroupAdd(ctx context.Context, groupPID string, req api.GroupAddRequest, userPID string) (api.GenericMessageResponse, int, error) {
	user, err := svc.DB.GetUserByPID(userPID)
	if err != nil {
		return api.GenericMessageResponse{}, 0, nil
	}
	group, err := svc.DB.GetGroupByPID(groupPID)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	hasPermission, err := svc.accesscontrolrsvc.UserHasPermissionToResource(user.ID, group.ID, tables.GroupResource, accesscontrolsvc.UpdateOperation)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	if !hasPermission {
		return api.GenericMessageResponse{}, 0, errors.New("user does not have permission to delete group")
	}

	return api.GenericMessageResponse{}, 0, nil
}

func (svc *GroupSvcImpl) GroupRemove(ctx context.Context, groupPID string, req api.GroupRemoveRequest, userPID string) (api.GenericMessageResponse, int, error) {
	return api.GenericMessageResponse{}, 0, nil
}
