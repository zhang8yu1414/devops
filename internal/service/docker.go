package service

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"zhangyudevops.com/internal/utils"
)

type sDocker struct{}

var insDocker = sDocker{}

func Docker() *sDocker {
	return &insDocker
}

// GenerateDockerClient 生成docker客户端
func GenerateDockerClient() *client.Client {
	var (
		err          error
		once         sync.Once
		dockerClient *client.Client
	)

	once.Do(func() {
		dockerClient, err = client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			glog.Error(context.Background(), err)
		}
	})
	return dockerClient
}

// LoginHarbor 登录harbor仓库
func LoginHarbor() error {
	var (
		once         sync.Once
		ctx          context.Context
		harborConfig *types.AuthConfig
		status       string
	)
	// 获取harbor配置
	addressVar, _ := g.Cfg().Get(ctx, "harbor.ip")
	usernameVar, _ := g.Cfg().Get(ctx, "harbor.username")
	passwordVar, _ := g.Cfg().Get(ctx, "harbor.password")
	harborConfig.Username = usernameVar.String()
	harborConfig.Password = passwordVar.String()
	harborConfig.ServerAddress = addressVar.String()

	// 登录仓库
	// TODO: 这里需要确认 status 返回的值是什么
	once.Do(func() {
		res, err := GenerateDockerClient().RegistryLogin(ctx, *harborConfig)
		if err != nil {
			glog.Error(ctx, err)
		}
		status = res.Status
		g.Dump("status", status)
	})
	return nil
}

// LoadImageAndPushToHarbor 导入docker镜像并推送到harbor仓库
func (s *sDocker) LoadImageAndPushToHarbor(ctx context.Context) (err error) {
	path, err := g.Config().Get(ctx, "docker.filePath")
	if err != nil {
		return err
	}
	// 获取当前目录下的文件
	files := utils.GetDirectoryDockerFileList(fmt.Sprintf("%s/%s", path, "docker"))

	for _, file := range files {
		_, filename := filepath.Split(file)

		// 查询tar包的状态是否为已经推送过仓库
		value, err := g.Model("file").Fields("import").Where("name=", filename).Value()
		if err != nil {
			return err
		}

		//	如果状态为0则表示则表示未处理过，需要写后续逻辑
		if value.Int() == 0 {
			f, err := os.Open(fmt.Sprintf("%s/%s/%s", path, "docker", filename))
			if err != nil {
				return err
			}

			res, err := GenerateDockerClient().ImageLoad(ctx, f, true)
			if err != nil {
				return err
			}
			body, err := ioutil.ReadAll(res.Body)
			strBody := gconv.String(body)
			image := strings.Split(strBody, " ")[2]
			oldImage := strings.Split(image, "\\")[0]

			// 开始处理tag逻辑
			// 拼凑harbor image tag
			harborVar, _ := g.Cfg().Get(ctx, "harbor.ip")
			harborIpAddress := harborVar.String()
			newImage := ""
			if ip := utils.FindIpAddress(image); ip != "" {
				//g.Dump("ip", ip)
				newImage = strings.Replace(oldImage, ip, harborIpAddress, 1)
			} else {
				newImage = fmt.Sprintf("%s/%s", harborIpAddress, oldImage)
			}

			// 生成新的tag image
			err = GenerateDockerClient().ImageTag(ctx, oldImage, newImage)
			if err != nil {
				return err
			}

			// 推送harbor仓库
			_, err = GenerateDockerClient().ImagePush(ctx, newImage, types.ImagePushOptions{})
			if err != nil {
				return err
			}

			// 操作数据库
			err = g.DB().Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
				//	插入image数据
				queryId, err := tx.Ctx(ctx).Model("file").
					Fields("id").
					Where("name", filename).Value()
				id := queryId.Int64()
				_, err = tx.Ctx(ctx).Model("image").
					Data(g.Map{"New": newImage, "Old": oldImage, "FileId": id}).Insert()
				if err != nil {
					return err
				}
				_, err = tx.Ctx(ctx).Model("file").
					Data(g.Map{"Import": 1}).
					Where("name", filename).Update()
				if err != nil {
					return err
				}

				return nil
			})
			if err != nil {
				return err
			}
		}
	}

	return
}
