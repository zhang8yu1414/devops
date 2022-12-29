package controller

import (
	"context"
	v1 "zhangyudevops.com/api/v1"
)

var Harbor = cHarbor{}

type cHarbor struct{}

func (c *cHarbor) CreateProject(ctx context.Context, req *v1.HarborProjectCreateReq) (res *v1.HarborProjectCreateRes, err error) {

}
