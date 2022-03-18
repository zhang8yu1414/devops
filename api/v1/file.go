package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type FileUploadReq struct {
	g.Meta `path:"/file/upload" method:"post" mime:"multipart/form-data" tags:"file" summary:"上传文件"`
	File   *ghttp.UploadFile `json:"file" type:"file" dc:"选择上传文件"`
}

type FileUploadRes struct{}
