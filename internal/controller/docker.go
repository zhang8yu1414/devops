package controller

import (
	"context"
	v1 "zhangyudevops.com/api/v1"
	"zhangyudevops.com/internal/service"
)

var Docker = cDocker{}

type cDocker struct{}

func (c *cDocker) PushImagesToHarbor(ctx context.Context, req *v1.DockerPushImagesReq) (res *v1.DockerPushImagesRes, err error) {
	err = service.Docker().LoadImageAndPushToHarbor(ctx)
	if err != nil {
		return nil, err
	}
	return
}
