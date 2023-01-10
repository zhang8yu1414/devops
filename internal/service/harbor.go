package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/util/gconv"
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
	c := g.Client()
	c.SetBasicAuth(user, password)
	return c
}

// CreateProject Harbor创建项目
func (s *sHarbor) CreateProject(ctx context.Context, name string) (err error) {
	url := HarborUrl("projects")

	if r, err := HarborClient.ContentJson().Post(ctx, url, g.Map{"project_name": name}); err != nil {
		return err
	} else {
		// 如果返回的内容为空，则直接返回空，如果不为空，则需要解析
		b := r.ReadAll()
		if gconv.String(b) == "" {
			g.Log().Infof(ctx, "Harbor项目名称%s创建成功", name)
			return nil
		}
		var (
			msg     *model.ApiResMessage
			message string
		)
		if err = gjson.Unmarshal(b, &msg); err != nil {
			return err
		}

		// 取出返回信息里的message信息
		for i, info := range msg.Errors {
			if i == 0 {
				message = info.Message
			}
		}

		// 把获取的message组成错误信息返回
		err = errors.New(message)
		g.Log().Infof(ctx, "Harbor项目名称%s已创建，本次创建失败", name)
		return err
	}

}

// GetProject 查询当前项目名称
func (s *sHarbor) GetProject(ctx context.Context, name string) (err error, info *model.ProjectInfo) {
	url := HarborUrl("projects")
	if r, err := HarborClient.Get(ctx, url, g.Map{"name": name, "with_detail": false}); err != nil {
		return err, nil
	} else {
		b := r.ReadAll()
		if gconv.String(b) == "[]\n" {
			err = errors.New("查询的项目名称不存在")
			return err, nil
		}
		var infos []*model.ProjectInfo
		if err = gjson.Unmarshal(b, &infos); err != nil {
			return err, nil
		}

		for i, in := range infos {
			if i == 0 {
				info = in
			}
		}
		return nil, info
	}

}

// DeleteProject Harbor删除项目
func (s *sHarbor) DeleteProject(ctx context.Context, name string) (err error) {
	suffix := fmt.Sprintf("projects/%s", name)
	url := HarborUrl(suffix)

	if _, err = HarborClient.ContentJson().Delete(ctx, url); err != nil {
		return err
	}
	return nil
}
