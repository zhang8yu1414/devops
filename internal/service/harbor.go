package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"zhangyudevops.com/internal/model"
)

type sHarbor struct{}

var (
	inHarbor     = sHarbor{}
	HarborUrl    = getHarborUrl
	HarborClient = initialHttpClient()
)

func Harbor() *sHarbor {
	return &inHarbor
}

// 组装harbor api url
func getHarborUrl(suffix string) (url string) {

	versionVar, _ := g.Config().Get(context.Background(), "harbor.version")
	version := versionVar.String()

	ipVar, _ := g.Config().Get(context.Background(), "harbor.ip")
	ip := ipVar.String()

	url = fmt.Sprintf("http://%s/api/%s/%s", ip, version, suffix)
	return
}

// 初始化http client
func initialHttpClient() *gclient.Client {
	userVar, _ := g.Config().Get(context.Background(), "harbor.username")
	user := userVar.String()

	passVar, _ := g.Config().Get(context.Background(), "harbor.password")
	password := passVar.String()
	c := g.Client().ContentJson()
	c.SetBasicAuth(user, password)
	return c
}

// @todo: 需要添加在返回错误数据的情况下，反正错误消息
func (s *sHarbor) CreateProject(ctx context.Context, name string) (err error) {
	url := HarborUrl("projects")

	if r, err := HarborClient.Post(ctx, url, g.Map{
		"project_name": name,
	}); err != nil {
		return err
	} else {
		var msg model.ApiResMessage
		err = gjson.Unmarshal(r.ReadAll(), &msg)
		if err != nil {
			return err
		}

		g.Dump(msg)
	}

	return
}
