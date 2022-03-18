package v1

import "github.com/gogf/gf/v2/frame/g"

type DockerPushImagesReq struct {
	g.Meta `path:"/docker/push" method:"get" tags:"DockerService" summary:"push the docker images to harbor"`
}

type DockerPushImagesRes struct{}
