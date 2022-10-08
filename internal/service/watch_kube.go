package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/informers"
	"time"
	"zhangyudevops.com/internal/model"
)

var inWatch = sWatch{}

type sWatch struct{}

func Watch() *sWatch {
	return &inWatch
}

var (
	Factory = GenerateInformerFactory()
)

// GenerateInformerFactory 初始化factory
func GenerateInformerFactory() informers.SharedInformerFactory {
	informerFactory := informers.NewSharedInformerFactoryWithOptions(
		kubeClient,
		time.Hour*4,
	)
	return informerFactory
}

// RunDeployInformer 执行deployment informer
func (s *sWatch) RunDeployInformer(namespace string) []*v1.Deployment {
	deployInformer := Factory.Apps().V1().Deployments()
	//informer := deployInformer.Informer()
	deployLister := deployInformer.Lister().Deployments(namespace)

	stopper := make(chan struct{})
	defer close(stopper)

	Factory.Start(stopper)
	Factory.WaitForCacheSync(stopper)

	deploy, err := deployLister.List(labels.Everything())
	if err != nil {
		panic(err)
	}

	return deploy
}

func (s *sWatch) GetDeployDataAndTransToWs(ns string) []*model.SourceManager {
	// ns下所有的deployment数量
	deploys := s.RunDeployInformer(ns)
	//
	result := make([]*model.SourceManager, len(deploys))
	for i := 0; i < len(deploys); i++ {
		// 组装images
		images := make([]string, len(deploys[i].Spec.Template.Spec.Containers))
		for j := 0; j < len(deploys[i].Spec.Template.Spec.Containers); j++ {
			images[j] = deploys[i].Spec.Template.Spec.Containers[j].Image
		}

		// 组装deployment列表
		result[i] = &model.SourceManager{
			Name:              deploys[i].Name,
			Namespace:         deploys[i].Namespace,
			CreateTimeStamp:   deploys[i].CreationTimestamp,
			Replicas:          &deploys[i].Status.Replicas,
			Image:             images,
			UpdatedReplicas:   deploys[i].Status.UpdatedReplicas,
			ReadyReplicas:     deploys[i].Status.ReadyReplicas,
			AvailableReplicas: deploys[i].Status.AvailableReplicas,
			Selector:          deploys[i].Spec.Selector.MatchLabels,
		}
	}
	return result
}

func (s *sWatch) TestDeploy(ns string, ch chan []*model.SourceManager) chan []*model.SourceManager {
	for {
		ch <- s.GetDeployDataAndTransToWs(ns)
		fmt.Println("ch", ch)
	}
}

// RunPodInformer watch pod informer
// 各个资源管理器都可以用这个方法获取pod列表
// 如果带selector参数， name为空，表示获取pod列表
// 如果selector为空，name有值就获取pod详情
func (s *sWatch) RunPodInformer(namespace string, sel map[string]string, name string) (ret []*v12.Pod) {

	var (
		selectorKey   string
		selectorValue []string
		//err           error
	)

	podInformer := Factory.Core().V1().Pods()
	//informer := podInformer.Informer()
	podLister := podInformer.Lister().Pods(namespace)

	//informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
	//	AddFunc:    onAdd,
	//	UpdateFunc: onUpdate,
	//	DeleteFunc: onDelete,
	//})

	stopper := make(chan struct{})
	defer close(stopper)

	Factory.Start(stopper)
	//Factory.WaitForCacheSync(stopper)

	if name == "" {
		// 获取selector key value
		for k, v := range sel {
			selectorKey = k
			selectorValue = append(selectorValue, v)
		}

		// 组装label.selector
		se, err := labels.NewRequirement(selectorKey, selection.Equals, selectorValue)
		if err != nil {
			g.Log().Error(context.Background(), err)
		}
		selector := labels.NewSelector()
		selector = selector.Add(*se)

		ret, err = podLister.List(selector)
		if err != nil {
			g.Log().Error(context.Background(), err)
		}
	} else {
		ret = make([]*v12.Pod, 0)
		r, err := podLister.Get(name)
		if err != nil {
			g.Log().Error(context.Background(), err)
		}
		ret = append(ret, r)
	}
	return
}

// GetDeploysPodsDataAndTransToWs 获取指定ns下的指定deploy的pod列表
func (s *sWatch) GetDeploysPodsDataAndTransToWs(ns string, sel map[string]string, name string) []*model.Pod {
	pods := s.RunPodInformer(ns, sel, name)

	podList := make([]*model.Pod, len(pods))
	for i := 0; i < len(pods); i++ {
		podList[i] = &model.Pod{
			Name:            pods[i].Name,
			Namespace:       pods[i].Namespace,
			CreateTimeStamp: pods[i].CreationTimestamp,
			Phase:           pods[i].Status.Phase,
			HostIp:          pods[i].Status.HostIP,
			PodIp:           pods[i].Status.PodIP,
		}
	}
	return podList
}

// GetPodInfoDataAndTransToWs 获取指定pod的详情
func (s *sWatch) GetPodInfoDataAndTransToWs(ns string, sel map[string]string, name string) (ret *model.PodInfo) {
	var (
		info []*model.ContainerInfo     // container info
		port []*model.ContainerPortInfo // container port
	)

	podInfo := s.RunPodInformer(ns, sel, name)

	// 组装container信息
	info = make([]*model.ContainerInfo, len(podInfo[0].Spec.Containers))
	for i := 0; i < len(podInfo[0].Spec.Containers); i++ {
		// 组装pod详情
		port = make([]*model.ContainerPortInfo, len(podInfo[0].Spec.Containers[i].Ports))
		for j := 0; j < len(podInfo[0].Spec.Containers[i].Ports); j++ {
			port[j] = &model.ContainerPortInfo{
				PortName:      podInfo[0].Spec.Containers[i].Ports[j].Name,
				ContainerPort: podInfo[0].Spec.Containers[i].Ports[j].ContainerPort,
				Protocol:      podInfo[0].Spec.Containers[i].Ports[j].Protocol,
			}
		}
		info[i] = &model.ContainerInfo{
			ContainerName:     podInfo[0].Spec.Containers[i].Name,
			ImagePullPolicy:   podInfo[0].Spec.Containers[i].ImagePullPolicy,
			ImageName:         podInfo[0].Spec.Containers[i].Image,
			ContainerPortInfo: port,
		}
	}

	// 组装pod状态
	condition := make([]*model.Condition, len(podInfo[0].Status.Conditions))
	for i, c := range podInfo[0].Status.Conditions {
		condition[i] = &model.Condition{
			Type:   c.Type,
			Status: c.Status,
		}
	}

	// 返回数据
	ret = &model.PodInfo{
		Name:            podInfo[0].Name,
		Namespace:       podInfo[0].Namespace,
		CreateTimeStamp: podInfo[0].CreationTimestamp,
		Phase:           podInfo[0].Status.Phase,
		HostIp:          podInfo[0].Status.HostIP,
		PodIp:           podInfo[0].Status.PodIP,
		ContainerInfo:   info,
		Conditions:      condition,
	}

	return
}

// RunStatefulSetInformer 生成sts informer
func (s *sWatch) RunStatefulSetInformer(namespace string) (ret []*model.SourceManager) {
	stsInformer := Factory.Apps().V1().StatefulSets()
	//informer := stsInformer.Informer()
	stsLister := stsInformer.Lister().StatefulSets(namespace)

	//informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
	//	AddFunc:    onAdd,
	//	UpdateFunc: onUpdate,
	//	DeleteFunc: onDelete,
	//})

	stopper := make(chan struct{})
	defer close(stopper)

	Factory.Start(stopper)
	//Factory.WaitForCacheSync(stopper)

	stsList, err := stsLister.List(labels.Everything())
	if err != nil {
		g.Log().Error(context.Background(), err)
	}

	// 把返回对象列表组装成自定义结构体
	ret = make([]*model.SourceManager, len(stsList))
	for i, set := range stsList {
		//组转pod images 列表
		images := make([]string, len(set.Spec.Template.Spec.Containers))
		for i2, container := range set.Spec.Template.Spec.Containers {
			images[i2] = container.Image
		}

		ret[i] = &model.SourceManager{
			Name:              set.Name,
			Namespace:         set.Namespace,
			CreateTimeStamp:   set.CreationTimestamp,
			Replicas:          set.Spec.Replicas,
			UpdatedReplicas:   set.Status.UpdatedReplicas,
			ReadyReplicas:     set.Status.ReadyReplicas,
			AvailableReplicas: set.Status.AvailableReplicas,
			Selector:          set.Spec.Selector.MatchLabels,
			Image:             images,
		}
	}
	return
}

// RunDaemonSetInformer 生成ds informer
func (s *sWatch) RunDaemonSetInformer(namespace string) (ret []*model.DaemonSet) {
	dsInformer := Factory.Apps().V1().DaemonSets()
	dsLister := dsInformer.Lister().DaemonSets(namespace)

	stopper := make(chan struct{})
	defer close(stopper)

	Factory.Start(stopper)
	//Factory.WaitForCacheSync(stopper)

	dsList, err := dsLister.List(labels.Everything())
	if err != nil {
		g.Log().Error(context.Background(), err)
	}

	ret = make([]*model.DaemonSet, len(dsList))
	for i, set := range dsList {
		//组转pod images 列表
		images := make([]string, len(set.Spec.Template.Spec.Containers))
		for i2, container := range set.Spec.Template.Spec.Containers {
			images[i2] = container.Image
		}

		ret[i] = &model.DaemonSet{
			Name:                   set.Name,
			Namespace:              set.Namespace,
			Image:                  images,
			CreateTimeStamp:        set.CreationTimestamp,
			CurrentNumberScheduled: set.Status.CurrentNumberScheduled,
			DesiredNumberScheduled: set.Status.DesiredNumberScheduled,
			NumberReady:            set.Status.NumberReady,
			Selector:               set.Spec.Selector.MatchLabels,
		}
	}
	return
}

func (s *sWatch) RunNodeInformer() (ret []*model.Nodes) {
	var (
		internalIP string
		status     v12.ConditionStatus
	)

	nodeInformer := Factory.Core().V1().Nodes()
	nodeLister := nodeInformer.Lister()

	stopper := make(chan struct{})
	defer close(stopper)

	Factory.Start(stopper)
	//Factory.WaitForCacheSync(stopper)

	nodesList, err := nodeLister.List(labels.Everything())
	if err != nil {
		g.Log().Error(context.Background(), err)
	}

	// 把返回对象列表组装成自定义结构体
	ret = make([]*model.Nodes, len(nodesList))
	for i, n := range nodesList {
		// 组装ip
		for _, addr := range n.Status.Addresses {
			if addr.Type == "InternalIP" {
				internalIP = addr.Address
			}
		}

		// 组装status
		for _, c := range n.Status.Conditions {
			if c.Type == "Ready" {
				status = c.Status
			}
		}

		ret[i] = &model.Nodes{
			Name:                    n.Name,
			InternalIP:              internalIP,
			CreateTimeStamp:         n.CreationTimestamp,
			KubeletVersion:          n.Status.NodeInfo.KubeletVersion,
			OsImage:                 n.Status.NodeInfo.OSImage,
			ContainerRuntimeVersion: n.Status.NodeInfo.ContainerRuntimeVersion,
			KernelVersion:           n.Status.NodeInfo.KernelVersion,
			Status:                  status,
		}
	}

	return
}

// GetNodeInfoDataAndTransToWs 获取node详情并传递给ws
func (s *sWatch) GetNodeInfoDataAndTransToWs(name string) (ret *model.NodeInfos) {
	nodeInformer := Factory.Core().V1().Nodes()
	nodeLister := nodeInformer.Lister()

	stopper := make(chan struct{})
	defer close(stopper)

	Factory.Start(stopper)
	//Factory.WaitForCacheSync(stopper)

	info, err := nodeLister.Get(name)
	if err != nil {
		g.Log().Error(context.Background(), err)
	}

	// 组装node source capacity
	capacity := &model.NodeSource{
		Cpu:              info.Status.Capacity.Cpu(),
		EphemeralStorage: info.Status.Capacity.StorageEphemeral(),
		Memory:           info.Status.Capacity.Memory(),
		Pods:             info.Status.Capacity.Pods(),
	}

	// 组装node source allocatable
	allocatable := &model.NodeSource{
		Cpu:              info.Status.Allocatable.Cpu(),
		EphemeralStorage: info.Status.Allocatable.StorageEphemeral(),
		Memory:           info.Status.Allocatable.Memory(),
		Pods:             info.Status.Allocatable.Pods(),
	}

	ret = &model.NodeInfos{
		Name:                    info.Name,
		Address:                 info.Status.Addresses,
		CreateTimeStamp:         info.CreationTimestamp,
		Conditions:              info.Status.Conditions,
		Capacity:                capacity,
		Allocatable:             allocatable,
		Taints:                  info.Spec.Taints,
		PodCIDR:                 info.Spec.PodCIDR,
		KubeletVersion:          info.Status.NodeInfo.KubeletVersion,
		OsImage:                 info.Status.NodeInfo.OSImage,
		ContainerRuntimeVersion: info.Status.NodeInfo.ContainerRuntimeVersion,
		KernelVersion:           info.Status.NodeInfo.KernelVersion,
		OperatingSystem:         info.Status.NodeInfo.OperatingSystem,
		Architecture:            info.Status.NodeInfo.Architecture,
		KubeProxyVersion:        info.Status.NodeInfo.KubeProxyVersion,
		Images:                  info.Status.Images,
		Labels:                  info.Labels,
		Annotations:             info.Annotations,
	}

	return
}
