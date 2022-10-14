package v1

import "github.com/gogf/gf/v2/frame/g"

type CompressedAndPushReq struct {
	g.Meta   `path:"/app/push" method:"post" tags:"UncompressedAndPushImages" summary:"Uncompressed TARFile, Then push images to harbor"`
	FileName string `json:"fileName"`
}

type CompressedAndPushRes struct{}
