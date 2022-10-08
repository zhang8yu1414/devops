package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	v1 "zhangyudevops.com/api/v1"
	"zhangyudevops.com/internal/consts"
	"zhangyudevops.com/internal/model"
	"zhangyudevops.com/internal/service"
)

type cWs struct {
}

var Ws = cWs{}

func (c *cWs) Websocket(ctx context.Context, req *v1.WebsocketReq) (res *v1.WebsocketRes, err error) {
	var (
		r       = g.RequestFromCtx(ctx)
		ws      *ghttp.WebSocket
		msg     model.WsMessage
		msgByte []byte
	)

	if ws, err = r.WebSocket(); err != nil {
		g.Log().Error(ctx, err)
		return
	}

	for {
		_, msgByte, err = ws.ReadMessage()
		if err != nil {
			return
		}

		// message decode
		if err = gjson.DecodeTo(msgByte, &msg); err != nil {
			_ = c.write(ws, model.WsMessage{
				Type: consts.WsTypeError,
				Msg:  fmt.Sprintf("invalid message: %s", err.Error()),
			})
			continue
		}

		g.Log().Print(ctx, msg)

		// 把msg.Data的数据转换为struct
		msgData := msg.Data
		b, _ := json.Marshal(msgData)
		reqData := model.RequestData{}
		_ = json.Unmarshal(b, &reqData)
		g.Log().Print(ctx, reqData)

		// 根据资源类型来进行ws数据传输
		switch msg.Type {
		case consts.WsTypeDeploy:
			// todo: 这里ws获取不到数据，需要改
			if msg.Data != nil {
				ch := make(chan []*model.SourceManager)
				g.Dump("deploy", <-service.Watch().TestDeploy(reqData.Namespace, ch))
				if err = c.write(ws, model.WsMessage{
					Type: consts.WsTypeDeploy,
					//Data: service.Watch().GetDeployDataAndTransToWs(reqData.Namespace),
					Data: <-service.Watch().TestDeploy(reqData.Namespace, ch),
				}); err != nil {
					g.Log().Error(ctx, err)
				}
			}
		case consts.WsTypePod:
			if msg.Data != nil {
				if err = c.write(ws, model.WsMessage{
					Type: consts.WsTypePod,
					Data: service.Watch().GetDeploysPodsDataAndTransToWs(
						reqData.Namespace,
						reqData.Selector,
						reqData.PodName,
					),
				}); err != nil {
					g.Log().Error(ctx, err)
				}
			}
		case consts.WsTypePodInfo:
			if msg.Data != nil {
				if err = c.write(ws, model.WsMessage{
					Type: consts.WsTypePodInfo,
					Data: service.Watch().GetPodInfoDataAndTransToWs(
						reqData.Namespace,
						reqData.Selector,
						reqData.PodName,
					),
				}); err != nil {
					g.Log().Error(ctx, err)
				}
			}
		case consts.WsTypeStatefulSet:
			if msg.Data != nil {
				if err = c.write(ws, model.WsMessage{
					Type: consts.WsTypeStatefulSet,
					Data: service.Watch().RunStatefulSetInformer(reqData.Namespace),
				}); err != nil {
					g.Log().Error(ctx, err)
				}
			}
		case consts.WsTypeDaemonSet:
			if msg.Data != nil {
				if err = c.write(ws, model.WsMessage{
					Type: consts.WsTypeDaemonSet,
					Data: service.Watch().RunDaemonSetInformer(reqData.Namespace),
				}); err != nil {
					g.Log().Error(ctx, err)
				}
			}
		case consts.WsTypeNode:
			if msg.Data != nil {
				if err = c.write(ws, model.WsMessage{
					Type: consts.WsTypeNode,
					Data: service.Watch().RunNodeInformer(),
				}); err != nil {
					g.Log().Error(ctx, err)
				}
			}
		case consts.WsTypeNodeInfo:
			if msg.Data != nil {
				if err = c.write(ws, model.WsMessage{
					Type: consts.WsTypeNodeInfo,
					Data: service.Watch().GetNodeInfoDataAndTransToWs(reqData.NodeName),
				}); err != nil {
					g.Log().Error(ctx, err)
				}
			}
		}
	}
}

// write sends message to current client.
func (c *cWs) write(ws *ghttp.WebSocket, msg model.WsMessage) error {
	msgBytes, err := gjson.Encode(msg)
	if err != nil {
		return err
	}

	return ws.WriteMessage(ghttp.WsMsgText, msgBytes)
}
