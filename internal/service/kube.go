package service

import (
	"context"
	"flag"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"path/filepath"
	"zhangyudevops.com/internal/model"
	"zhangyudevops.com/internal/utils"
)

type sKube struct{}

var (
	inKube     = sKube{}
	kubeClient = GenerateK8sClient()
)

func Kube() *sKube {
	return &inKube
}

// GenerateK8sClient 生成k8s客户端
func GenerateK8sClient() *kubernetes.Clientset {
	var kubeConfig *string
	configPath, _ := g.Cfg().Get(context.Background(), "kube.configPath")
	path := configPath.String()
	kubeConfig = flag.String("kubeConfig", filepath.Join(path, "kubeConfig"), "absolute path to kube config")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		panic(err)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return clientSet
}

// GetNamespaces 获取namespaces信息
func (s *sKube) GetNamespaces(ctx context.Context) (namespace []*model.Namespace, err error) {
	ns, err := kubeClient.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	namespace = make([]*model.Namespace, len(ns.Items))
	for i, n := range ns.Items {
		namespace[i] = &model.Namespace{
			Name:              n.Name,
			CreationTimestamp: n.CreationTimestamp.Time,
		}
	}

	return
}

// GetNodesList 查询node列表信息
func (s *sKube) GetNodesList(ctx context.Context) (nodeList []*model.Node, err error) {
	nodes, err := kubeClient.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var nodeStatus v1.ConditionStatus
	nodeList = make([]*model.Node, len(nodes.Items))
	for i, node := range nodes.Items {
		for _, condition := range node.Status.Conditions {
			if condition.Type == "Ready" {
				nodeStatus = condition.Status
			}
		}
		nodeList[i] = &model.Node{
			Name:              node.Name,
			CreationTimestamp: node.CreationTimestamp.Time,
			NodeStatus:        nodeStatus,
		}
	}
	return
}

// GetNodeInfo 查询node信息详情
func (s *sKube) GetNodeInfo(ctx context.Context, nodeName string) (nodeInfo *model.NodeInfo, err error) {
	node, err := kubeClient.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	nodeInfo = &model.NodeInfo{
		Name:      node.Name,
		Labels:    node.Labels,
		Taint:     node.Spec.Taints,
		Addresses: node.Status.Addresses,
		Capacity: model.Capacity{
			Cpu:    node.Status.Capacity.Cpu().String(),
			Memory: node.Status.Capacity.Memory().String(),
			Pods:   node.Status.Capacity.Pods().String(),
		},
		SystemInfo: model.SystemInfo{
			KernelVersion:           node.Status.NodeInfo.KernelVersion,
			OSImage:                 node.Status.NodeInfo.OSImage,
			ContainerRuntimeVersion: node.Status.NodeInfo.ContainerRuntimeVersion,
			KubeletVersion:          node.Status.NodeInfo.KubeletVersion,
			KubeProxyVersion:        node.Status.NodeInfo.KubeProxyVersion,
			OperatingSystem:         node.Status.NodeInfo.OperatingSystem,
			Architecture:            node.Status.NodeInfo.Architecture,
		},
	}

	return
}

// GetPods 获取k8s 指定namespace pod列表
func (s *sKube) GetPods(ctx context.Context, namespace string, selector string) (podsList []*model.Pods, err error) {
	podList, err := kubeClient.CoreV1().Pods(namespace).Watch(ctx, metav1.ListOptions{
		LabelSelector: selector,
	})
	if err != nil {
		return
	}

	podsList = make([]*model.Pods, len(podList.ResultChan()))
	go func() {

		for event := range podList.ResultChan() {
			p, ok := event.Object.(*v1.Pod)
			if !ok {
				g.Log().Print(ctx, "unexpected type")
			}
			podsList = append(podsList, &model.Pods{
				Name:         p.Name,
				DurationTime: utils.TransformTimestamp(p.CreationTimestamp.Time),
				PodStatus:    p.Status.Phase,
				PodIp:        p.Status.PodIP,
				HostIp:       p.Status.HostIP,
			})
		}
	}()
	return
}

// GetDeployPodInfo 获取指定pod的详细信息
func (s *sKube) GetDeployPodInfo(ctx context.Context, namespace string, deploy string) (err error) {
	var (
		key   string
		value string
	)

	deployment, err := kubeClient.AppsV1().Deployments(namespace).Get(ctx, deploy, metav1.GetOptions{})
	if err != nil {
		return
	}

	label := deployment.Spec.Selector.MatchLabels
	for k, v := range label {
		key, value = k, v
	}

	// 组装当前deploy label selector
	selector := fmt.Sprintf("%s=%s", key, value)
	podList, err := s.GetPods(ctx, namespace, selector)
	if err != nil {
		return
	}
	g.Dump("podList", podList)
	return
}

// GetNamespaceDeploys 获取指定namespace下的deployment资源
func (s *sKube) GetNamespaceDeploys(ctx context.Context, namespace string) (deploys []*model.Deployments, err error) {
	deploysList, err := kubeClient.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return
	}

	deploys = make([]*model.Deployments, len(deploysList.Items))
	var status v1.ConditionStatus

	for i, deploy := range deploysList.Items {
		// deploy的状态
		for _, condition := range deploy.Status.Conditions {
			if condition.Type == "Available" {
				status = condition.Status
			}
		}

		// 组装container images
		images := make([]string, len(deploy.Spec.Template.Spec.Containers))
		for j, c := range deploy.Spec.Template.Spec.Containers {
			images[j] = c.Image
		}

		durationTime := utils.TransformTimestamp(deploy.CreationTimestamp.Time)
		deploys[i] = &model.Deployments{
			Name:         deploy.Name,
			Namespace:    deploy.Namespace,
			CreationTime: durationTime,
			Current:      deploy.Status.ReadyReplicas,
			Replicas:     deploy.Status.Replicas,
			Images:       images,
			Status:       status,
		}
	}

	return
}

func (s *sKube) EditDeploy(ctx context.Context, namespace string) {
	//kubeClient.AppsV1().Deployments(namespace).Apply(ctx, )
}

func (s *sKube) ApplyK8sYaml() {

}
