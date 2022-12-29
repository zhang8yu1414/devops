package controller

import (
	"context"
	"errors"
	v1 "zhangyudevops.com/api/v1"
	"zhangyudevops.com/internal/service"
)

var Update = cUpdate{}

type cUpdate struct{}

// UncompressedAndPush 解压升级文件，随后推送镜像到harbor
func (c *cUpdate) UncompressedAndPush(ctx context.Context, req *v1.CompressedAndPushReq) (res *v1.CompressedAndPushRes, err error) {
	// post请求，filename必填，如果为空，则报错
	if req.FileName == "" {
		err = errors.New("filename can't be null")
		return
	}
	err = service.Update().UncompressedTarFileAndPushHarbor(ctx, req.FileName)
	if err != nil {
		return nil, err
	}
	return
}
