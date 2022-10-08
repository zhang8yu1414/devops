package v1

import "github.com/gogf/gf/v2/frame/g"

type WebsocketReq struct {
	g.Meta `path:"/ws" method:"get" tags:"websocket" summary:"send message"`
}

type WebsocketRes struct {
}
