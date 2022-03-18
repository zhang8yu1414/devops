package controller

import (
	"context"
	v1 "zhangyudevops.com/api/v1"
)

var Docker = cDocker{}

type cDocker struct{}

func (c *cDocker) PushImagesToHarbor(ctx context.Context, req *v1.DockerPushImagesReq) (res *v1.DockerPushImagesRes, err error) {

	return
}
