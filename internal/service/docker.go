package service

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gogf/gf/v2/os/glog"
	"sync"
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

// GetContainersList 获取docker容器列表
func (s *sDocker) GetContainersList(ctx context.Context) ([]types.Container, error) {
	containers, err := GenerateDockerClient().ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		glog.Error(ctx, err)
	}
	return containers, err
}
