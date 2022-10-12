package v1

import "github.com/gogf/gf/v2/frame/g"

type UpdateAppReq struct {
	g.Meta   `path:"/app/update" method:"post" tags:"UpdateApp" summary:"Update Application"`
	FileName string `json:"fileName"`
}

type UpdateAppRes struct{}
