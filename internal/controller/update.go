package controller

import (
	"context"
	v1 "zhangyudevops.com/api/v1"
	"zhangyudevops.com/internal/service"
)

var Update = cUpdate{}

type cUpdate struct{}

func (c *cUpdate) UpdateApp(ctx context.Context, req *v1.UpdateAppReq) (res *v1.UpdateAppRes, err error) {
	err = service.Update().UpdateApp(ctx, req.FileName)
	if err != nil {
		return nil, err
	}
	return
}
