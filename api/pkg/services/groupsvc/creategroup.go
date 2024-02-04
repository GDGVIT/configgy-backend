package groupsvc

import (
	"context"

	"github.com/GDGVIT/configgy-backend/api/pkg/api"
	"github.com/GDGVIT/configgy-backend/pkg/tables"
)

func (svc *GroupSvcImpl) GroupCreate(ctx context.Context, req api.GroupCreateRequest, userPID string) (api.GenericMessageResponse, int, error) {
	user, err := svc.DB.GetUserByPID(userPID)
	if err != nil {
		return api.GenericMessageResponse{}, 0, err
	}
	members := map[string]tables.Permission{}
	for _, member := range req.Members {
		members[member.UserPid] = tables.Permission(member.GroupRole)
	}
	err = svc.DB.CreateGroup(user, &tables.Groups{
		Name:        req.Name,
		Description: req.Description,
	}, members)
	if err != nil {
		return api.GenericMessageResponse{}, 0, nil
	}

	return api.GenericMessageResponse{}, 0, nil
}
