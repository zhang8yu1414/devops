package controller

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	v1 "zhangyudevops.com/api/v1"
	"zhangyudevops.com/internal/service"
)

var Kube = cKube{}

type cKube struct{}

// GetNamespaces 获取集群namespace
func (c *cKube) GetNamespaces(ctx context.Context, req *v1.KubeNamespacesReq) (res *v1.KubeNamespacesRes, err error) {
	namespaceInfo, err := service.Kube().GetNamespaces(ctx)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}
	res = &v1.KubeNamespacesRes{
		Namespaces: namespaceInfo,
	}
	return
}

// GetNodes 获取node列表
func (c *cKube) GetNodes(ctx context.Context, req *v1.KubeNodesReq) (res *v1.KubeNodesRes, err error) {
	nodeList, err := service.Kube().GetNodesList(ctx)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}
	res = &v1.KubeNodesRes{
		Nodes: nodeList,
	}
	return
}

// GetNodeInfo 获取node信息列表
func (c cKube) GetNodeInfo(ctx context.Context, req *v1.KubeNodeInfoReq) (res *v1.KubeNodeInfoRes, err error) {
	nodeInfo, err := service.Kube().GetNodeInfo(ctx, req.NodeName)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	res = &v1.KubeNodeInfoRes{
		NodeInfo: nodeInfo,
	}
	return
}

// GetPods 获取指定namespace pods列表
func (c cKube) GetPods(ctx context.Context, req *v1.KubePodsReq) (res *v1.KubePodsRes, err error) {
	podInfo, err := service.Kube().GetPods(ctx, req.Namespace, req.Selector)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	res = &v1.KubePodsRes{PodInfo: podInfo}
	return
}

// GetPodInfo 获取pod详情
func (c *cKube) GetPodInfo(ctx context.Context, req *v1.KubePodInfoReq) (res *v1.KubePodInfoRes, err error) {
	//podInfo, err := service.Kube().GetDeployPodInfo(ctx, req.Namespace, req.Deployment)
	err = service.Kube().GetDeployPodInfo(ctx, req.Namespace, req.Deployment)
	if err != nil {
		return
	}

	return
}

// GetNamespaceDeploy 获取指定namespace下的deployment资源
func (c *cKube) GetNamespaceDeploy(ctx context.Context, req *v1.KubeDeploysReq) (res *v1.KubeDeploysRes, err error) {
	deploys, err := service.Kube().GetNamespaceDeploys(ctx, req.Namespace)
	if err != nil {
		return
	}

	res = &v1.KubeDeploysRes{Deploys: deploys}
	return
}
