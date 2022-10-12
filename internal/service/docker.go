package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"zhangyudevops.com/internal/utils"
)

type sDocker struct{}

var (
	insDocker    = sDocker{}
	DockerClient = GenerateDockerClient()
)

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

// GenerateHarborAuthConfig 生成harbor登录密匙
func GenerateHarborAuthConfig(ctx context.Context) types.ImagePushOptions {
	var (
		harborConfig types.AuthConfig
	)
	// 获取harbor配置
	addressVar, _ := g.Cfg().Get(ctx, "harbor.ip")
	usernameVar, _ := g.Cfg().Get(ctx, "harbor.username")
	passwordVar, _ := g.Cfg().Get(ctx, "harbor.password")
	harborConfig.Username = usernameVar.String()
	harborConfig.Password = passwordVar.String()
	harborConfig.ServerAddress = addressVar.String()

	authConfigBytes, _ := json.Marshal(harborConfig)
	authConfigEncoded := base64.URLEncoding.EncodeToString(authConfigBytes)
	opts := types.ImagePushOptions{RegistryAuth: authConfigEncoded}
	return opts
}

// LoadImageAndPushToHarbor 导入docker镜像并推送到harbor仓库
func (s *sDocker) LoadImageAndPushToHarbor(ctx context.Context, tarFileName string, uncompressedPath string) (err error) {
	// 获取当前目录下的文件
	files := utils.GetDirectoryDockerFileList(fmt.Sprintf("%s/%s", uncompressedPath, "images"))

	for _, file := range files {
		_, filename := filepath.Split(file)

		f, err := os.Open(fmt.Sprintf("%s/%s/%s", uncompressedPath, "images", filename))
		if err != nil {
			return err
		}

		res, err := DockerClient.ImageLoad(ctx, f, true)
		if err != nil {
			return err
		}
		body, err := io.ReadAll(res.Body)
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
		err = DockerClient.ImageTag(ctx, oldImage, newImage)
		if err != nil {
			return err
		}

		// 推送harbor仓库
		opts := GenerateHarborAuthConfig(ctx)
		_, _ = DockerClient.ImagePush(ctx, newImage, opts)
		g.Log().Infof(ctx, "%s , image push success", newImage)

		// 操作数据库
		err = g.DB().Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
			//	插入image数据
			queryId, err := tx.Ctx(ctx).Model("file").
				Fields("id").
				Where("name", tarFileName).Value()
			id := queryId.Int64()
			// 考虑到重复推送镜像，删除同文件ID下的导入记录
			_, _ = tx.Ctx(ctx).Model("image").Where(g.Map{
				"file_id": id,
				"new":     newImage,
				"old":     oldImage,
			}).Delete()

			_, err = tx.Ctx(ctx).Model("image").
				Data(g.Map{"New": newImage, "Old": oldImage, "FileId": id}).Insert()
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return err
		}
	}

	return
}
