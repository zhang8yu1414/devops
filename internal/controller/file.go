package controller

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	v1 "zhangyudevops.com/api/v1"
	"zhangyudevops.com/internal/service"
)

var File = cFile{}

type cFile struct{}

// Upload 文件上传
func (c *cFile) Upload(ctx context.Context, req *v1.FileUploadReq) (res *v1.FileUploadRes, err error) {
	if req.File == nil {
		return nil, gerror.NewCode(gcode.CodeMissingParameter, "请选择需要上传的文件")
	}
	err = service.File().UploadFile(ctx, req.File)
	if err != nil {
		return nil, err
	}
	return
}
