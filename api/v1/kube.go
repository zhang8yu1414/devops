package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"zhangyudevops.com/internal/model"
)

type KubeNamespacesReq struct {
	g.Meta `path:"/kube/get-namespaces" method:"get"  tags:"kube" summary:"获取集群namespaces列表"`
}

type KubeNamespacesRes struct {
	Namespaces []*model.Namespace
}

type KubeNodesReq struct {
	g.Meta `path:"/kube/get-nodes" method:"get"  tags:"kube" summary:"获取集群指定node详情"`
}

type KubeNodesRes struct {
	Nodes []*model.Node
}

type KubeNodeInfoReq struct {
	g.Meta   `path:"/kube/get-node-info" method:"post"  tags:"kube" summary:"获取集群指定node详情"`
	NodeName string `json:"nodeName"`
}

type KubeNodeInfoRes struct {
	NodeInfo *model.NodeInfo
}

type KubePodsReq struct {
	g.Meta    `path:"/kube/get-pods" method:"post"  tags:"kube" summary:"获取指定namespace pods列表"`
	Namespace string `json:"namespace"`
	Selector  string `json:"selector,omitempty"`
}

type KubePodsRes struct {
	PodInfo []*model.Pods
}

type KubePodInfoReq struct {
	g.Meta     `path:"/kube/get-pod-info" method:"post"  tags:"kube" summary:"获取指定pod详情"`
	Namespace  string `json:"namespace"`
	Deployment string `json:"deployment"`
}

type KubePodInfoRes struct {
}

type KubeDeploysReq struct {
	g.Meta    `path:"/kube/get-ns-deploy" method:"post"  tags:"kube" summary:"获取指定namespace下的deploy"`
	Namespace string `json:"namespace"`
}

type KubeDeploysRes struct {
	Deploys []*model.Deployments
}
