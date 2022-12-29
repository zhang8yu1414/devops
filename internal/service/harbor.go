package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"zhangyudevops.com/internal/model"
)

type sHarbor struct{}

var (
	inHarbor   = sHarbor{}
	HarborUser = getHarborUserInfo()
)

func Harbor() *sHarbor {
	return &inHarbor
}

// @todo: 这里需要考虑是否使用g.Client去调用接口，把配置写入全局调用
func getHarborUserInfo() (userInfo *model.LoginInfo) {
	userVar, _ := g.Config().Get(context.Background(), "harbor.username")
	user := userVar.String()

	passVar, _ := g.Config().Get(context.Background(), "harbor.password")
	password := passVar.String()

	info := fmt.Sprintf("%s:%s", user, password)

	versionVar, _ := g.Config().Get(context.Background(), "harbor.version")
	version := versionVar.String()

	ipVar, _ := g.Config().Get(context.Background(), "harbor.ip")
	ip := ipVar.String()

	url := fmt.Sprintf("http://%s/api/%s", ip, version)

	userInfo = &model.LoginInfo{
		Url:      url,
		UserInfo: info,
	}

	return
}

func (s *sHarbor) CreateProject(ctx context.Context, name string) error {

	g.Client().Post(ctx, HarborUser.Url)
}
