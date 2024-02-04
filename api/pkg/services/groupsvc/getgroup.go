package groupsvc

import (
	"context"

	"github.com/GDGVIT/configgy-backend/api/pkg/api"
)

func (svc *GroupSvcImpl) GroupGet(ctx context.Context, groupPID string, userPID string) (api.Group, int, error) {
	return api.Group{}, 0, nil
}

func (svc *GroupSvcImpl) GroupsGet(ctx context.Context, userPID string) ([]api.Group, int, error) {
	return []api.Group{}, 0, nil
}
