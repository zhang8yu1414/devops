package model

type ApiResMessage struct {
	Errors []info `json:"errors"`
}

type info struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ProjectInfo 项目详情
type ProjectInfo struct {
	Name       string `json:"name"`
	RepoCount  int    `json:"repo_count"`
	CreateTime string `json:"creation_time"`
	UpdateTime string `json:"update_time"`
}
