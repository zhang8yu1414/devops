package controller

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
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

func (c *cHarbor) DeleteProject(ctx context.Context, req *v1.HarborProjectDeleteReq) (res *v1.HarborProjectDeleteRes, err error) {
	if err = service.Harbor().DeleteProject(ctx, req.ProjectName); err != nil {
		return nil, err
	}

	return
}

func (c *cHarbor) GetProject(ctx context.Context, req *v1.HarborProjectGetReq) (res *v1.HarborProjectGetRes, err error) {
	err, info := service.Harbor().GetProject(ctx, req.ProjectName)
	if err != nil {
		return nil, err
	}
	g.Dump(info)
	res = &v1.HarborProjectGetRes{
		Info: info,
	}

	return
}
