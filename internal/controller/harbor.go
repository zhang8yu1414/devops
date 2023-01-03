package controller

import (
	"context"
	v1 "zhangyudevops.com/api/v1"
	"zhangyudevops.com/internal/service"
)

var Harbor = cHarbor{}

type cHarbor struct{}

func (c *cHarbor) CreateProject(ctx context.Context, req *v1.HarborProjectCreateReq) (res *v1.HarborProjectCreateRes, err error) {
	if err = service.Harbor().CreateProject(ctx, req.ProjectName); err != nil {
		return nil, err
	}

	return
}
