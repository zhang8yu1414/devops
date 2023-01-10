package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"zhangyudevops.com/internal/model"
)

// HarborProjectCreateReq 项目创建
type HarborProjectCreateReq struct {
	g.Meta      `path:"/harbor/project/create" method:"post" tags:"Project Create" summary:"创建项目"`
	ProjectName string `json:"projectName"`
}

type HarborProjectCreateRes struct{}

// HarborProjectDeleteReq 项目删除
type HarborProjectDeleteReq struct {
	g.Meta      `path:"/harbor/project/delete" method:"post" tags:"Project Delete" summary:"项目删除"`
	ProjectName string `json:"projectName"`
}

type HarborProjectDeleteRes struct{}

type HarborProjectGetReq struct {
	g.Meta      `path:"/harbor/project/list" method:"post" tags:"Project list" summary:"项目列表"`
	ProjectName string `json:"projectName"`
}

type HarborProjectGetRes struct {
	Info *model.ProjectInfo `json:"info"`
}
