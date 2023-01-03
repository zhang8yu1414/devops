package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type HarborProjectCreateReq struct {
	g.Meta      `path:"/harbor/project/create" method:"post" tags:"Project Create" summary:"创建项目"`
	ProjectName string `json:"projectName"`
}

type HarborProjectCreateRes struct {
}
