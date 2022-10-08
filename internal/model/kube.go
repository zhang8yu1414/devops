package model

import (
	v1 "k8s.io/api/core/v1"
	"time"
)

type Namespace struct {
	Name              string
	CreationTimestamp time.Time
}

// Node 信息
type Node struct {
	Name              string
	CreationTimestamp time.Time
	NodeStatus        v1.ConditionStatus
}

// NodeInfo 详情
type NodeInfo struct {
	Name      string
	Labels    map[string]string
	Taint     []v1.Taint
	Addresses []v1.NodeAddress
	Capacity
	SystemInfo
}

type Capacity struct {
	Cpu    string
	Memory string
	Pods   string
}

type SystemInfo struct {
	KernelVersion           string
	OSImage                 string
	ContainerRuntimeVersion string
	KubeletVersion          string
	KubeProxyVersion        string
	OperatingSystem         string
	Architecture            string
}

// Pods 返回pod列表
type Pods struct {
	Name         string
	DurationTime string // 创建时间
	PodStatus    v1.PodPhase
	PodIp        string
	HostIp       string
}

type DockerImageInfo struct {
	Name  string
	Image string
}

// Deployments 资源管理
type Deployments struct {
	Name         string
	Namespace    string
	CreationTime string
	Current      int32
	Replicas     int32
	Images       []string
	Status       v1.ConditionStatus
}
