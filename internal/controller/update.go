package controller

import (
	"context"
	v1 "zhangyudevops.com/api/v1"
	"zhangyudevops.com/internal/service"
)

var Update = cUpdate{}

type cUpdate struct{}

// UncompressedAndPush 解压升级文件，随后推送镜像到harbor
func (c *cUpdate) UncompressedAndPush(ctx context.Context, req *v1.CompressedAndPushReq) (res *v1.CompressedAndPushRes, err error) {
	err = service.Update().UncompressedTarFileAndPushHarbor(ctx, req.FileName)
	if err != nil {
		return nil, err
	}
	return
}
