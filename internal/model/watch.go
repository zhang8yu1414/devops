package model

import (
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SourceManager 详情
// deploy. ds 资源管理器，通用结构
type SourceManager struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	CreateTimeStamp   v1.Time           `json:"createTimeStamp"`
	Replicas          *int32            `json:"replicas"`
	Image             []string          `json:"image"`
	UpdatedReplicas   int32             `json:"updatedReplicas"`
	ReadyReplicas     int32             `json:"readyReplicas"`
	AvailableReplicas int32             `json:"availableReplicas"`
	Selector          map[string]string `json:"selector"`
}

// Pod 基本信息
type Pod struct {
	Name            string       `json:"name"`
	Namespace       string       `json:"namespace"`
	CreateTimeStamp v1.Time      `json:"createTimeStamp"`
	Phase           v12.PodPhase `json:"phase"`
	HostIp          string       `json:"hostIp"`
	PodIp           string       `json:"podIp"`
}

// PodInfo pod的详情
type PodInfo struct {
	Name            string           `json:"name"`
	Namespace       string           `json:"namespace"`
	CreateTimeStamp v1.Time          `json:"createTimeStamp"`
	Phase           v12.PodPhase     `json:"phase"`
	HostIp          string           `json:"hostIp"`
	PodIp           string           `json:"podIp"`
	Conditions      []*Condition     `json:"conditions"`
	ContainerInfo   []*ContainerInfo `json:"containerInfo"`
}

// Condition pod 状态
type Condition struct {
	Type   v12.PodConditionType `json:"type"`
	Status v12.ConditionStatus  `json:"status"`
}

// ContainerInfo pod里面的容器详情
type ContainerInfo struct {
	ContainerName     string               `json:"name"`
	ImagePullPolicy   v12.PullPolicy       `json:"imagePullPolicy"`
	ImageName         string               `json:"imageName"`
	ContainerPortInfo []*ContainerPortInfo `json:"containerPortInfo"`
}

// ContainerPortInfo 容器端口详情
type ContainerPortInfo struct {
	PortName      string       `json:"portName"`
	ContainerPort int32        `json:"containerPort"`
	Protocol      v12.Protocol `json:"protocol"`
}

// DaemonSet 结构
type DaemonSet struct {
	Name                   string            `json:"name"`
	Namespace              string            `json:"namespace"`
	CreateTimeStamp        v1.Time           `json:"createTimeStamp"`
	Image                  []string          `json:"image"`
	CurrentNumberScheduled int32             `json:"currentNumberScheduled"`
	DesiredNumberScheduled int32             `json:"desiredNumberScheduled"`
	NumberReady            int32             `json:"numberReady"`
	Selector               map[string]string `json:"selector"`
}

// Nodes 获取node基本信息
type Nodes struct {
	Name                    string              `json:"name"`
	InternalIP              string              `json:"internalIP"`
	CreateTimeStamp         v1.Time             `json:"createTimeStamp"`
	KubeletVersion          string              `json:"kubeletVersion"`
	OsImage                 string              `json:"osImage"`
	ContainerRuntimeVersion string              `json:"containerRuntimeVersion"`
	KernelVersion           string              `json:"kernelVersion"`
	Status                  v12.ConditionStatus `json:"status"`
}

// NodeInfos 获取node详细信息
type NodeInfos struct {
	Name                    string               `json:"name"`
	Address                 []v12.NodeAddress    `json:"address"` // 存internalIP与hostname
	CreateTimeStamp         v1.Time              `json:"createTimeStamp"`
	Conditions              []v12.NodeCondition  `json:"conditions"`
	Capacity                *NodeSource          `json:"capacity"`
	Allocatable             *NodeSource          `json:"allocatable"`
	Taints                  []v12.Taint          `json:"taints"`
	PodCIDR                 string               `json:"podCIDR"`
	KubeletVersion          string               `json:"kubeletVersion"`
	OsImage                 string               `json:"osImage"`
	ContainerRuntimeVersion string               `json:"containerRuntimeVersion"`
	KernelVersion           string               `json:"kernelVersion"`
	OperatingSystem         string               `json:"operatingSystem"`
	Architecture            string               `json:"architecture"`
	KubeProxyVersion        string               `json:"kubeProxyVersion"`
	Images                  []v12.ContainerImage `json:"images"`
	Labels                  map[string]string    `json:"labels"`
	Annotations             map[string]string    `json:"annotations"`
}

// NodeSource node资源情况
type NodeSource struct {
	Cpu              *resource.Quantity `json:"cpu"`
	EphemeralStorage *resource.Quantity `json:"ephemeralStorage"`
	Memory           *resource.Quantity `json:"memory"`
	Pods             *resource.Quantity `json:"pods"`
}
