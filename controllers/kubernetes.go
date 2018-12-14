package controllers

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"strconv"
	"strings"
)


//默认首页登录入口控制器




type KubernetesController struct {
	//beego.Controller
	BaseController
}



func (c *KubernetesController) GetNodePage() {
	//如何获取用户名，除了cookie?
	//fmt.Println("username is ",c.GetString("username"))
	c.GetNodeData()

	c.TplName = "opsmanager/k8s/node.html"
	c.Layout = "common/layout_home.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["treeviewhtml"] = "role/ops/treeview.html"
}

func (c *KubernetesController) GetNamespace() {
	//如何获取用户名，除了cookie?
	//fmt.Println("username is ",c.GetString("username"))
	c.GetNamespaces()

	c.Data["json"] =c.Data["namespaces"]
	c.ServeJSON()
}

func (c *KubernetesController) GetNamespaceInformation() {
	//如何获取用户名，除了cookie?
	//fmt.Println("username is ",c.GetString("username"))
	nsName := c.Ctx.Input.Param(":ns")
	c.Data["nsName"] = nsName
	c.GetDeploymentData()
	c.GetPodsData()
	c.GetServiceData()
	c.TplName = "opsmanager/k8s/namespaces.html"
	c.Layout = "common/layout_home.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["treeviewhtml"] = "role/ops/treeview.html"
}

func (c *KubernetesController) GetNamespaces() {
	namespace, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		fmt.Println("k8s get node err is ",err)
	}
	var namespaces  []string

	for _,ns := range namespace.Items {
		namespaces = append(namespaces,ns.Name)
	}
	c.Data["namespaces"] = namespaces
}

func (c *KubernetesController) GetNodeData() {
	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		fmt.Println("k8s get node err is ",err)
	}
	nodeInformations :=  []k8sNode{}
	for _,node := range nodes.Items{
		cpu,_ := node.Status.Allocatable.Cpu().AsInt64()
		memory,_ := node.Status.Allocatable.Memory().AsInt64()
		memory = memory / 1024 / 1024 / 1024 + 1
		nodeInformation	:= k8sNode{
			Name: node.Name,
			Ip: node.Status.Addresses[0].Address,
			Status: string(node.Status.Conditions[len(node.Status.Conditions)-1].Type),
			Cpu: cpu,
			Memory: memory,
			JoinDate: node.CreationTimestamp.Format("2006-01-02"),
			System: node.Status.NodeInfo.OSImage,
			DockerVersion: node.Status.NodeInfo.ContainerRuntimeVersion,
		}
		nodeInformations = append(nodeInformations,nodeInformation)

	}
	c.Data["k8snode"] = nodeInformations

}


type deployMent struct {
	Name string
	Label []string
	Replic int32
	ReadyReplic int32
	CreateDate string
	Image string
}

type nsPod struct {
	Name string
	Status string
	NodeIp string
	PodIp string
	RestartCount int32
	CreateDate string
}

type nsService struct {
	Name string
	Label []string
	ClusterIp string
	DomainName string
	EndPoints []string
	CreateDate string
}


func (c *KubernetesController) GetDeploymentData() {

	deploymentList, err := clientset.AppsV1beta1().Deployments(c.Data["nsName"].(string)).List(metav1.ListOptions{})
	if err != nil {
		log.Println(err)
	}
	deploymentInformation := []deployMent{}
	for _,dp := range deploymentList.Items {
		var labelstr []string
		for k,v := range dp.Labels {
			tmpstr := k + ": " + v
			labelstr = append(labelstr,tmpstr)
		}

		tmpdp := deployMent{
			Name: dp.Name,
			Replic: *dp.Spec.Replicas,
			ReadyReplic: dp.Status.ReadyReplicas,
			CreateDate: dp.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Image: dp.Spec.Template.Spec.Containers[0].Image,
			Label: labelstr,
		}
		deploymentInformation = append(deploymentInformation,tmpdp)
	}


	c.Data["k8sdeploy"] = deploymentInformation


}

func (c *KubernetesController) GetPodsData() {

	pods,err := clientset.CoreV1().Pods(c.Data["nsName"].(string)).List(metav1.ListOptions{})
	if err != nil {
		log.Println(err.Error())
	}
	podSlce := []nsPod{}
	//获取POD容器状态，默认pod中只有一个应用容器


	for _,pod := range pods.Items {

		var status string
		if pod.Status.ContainerStatuses[0].State.Terminated != nil {
			status = "Terminated " +  pod.Status.ContainerStatuses[0].State.Terminated.Reason
		} else if pod.Status.ContainerStatuses[0].State.Running != nil {
			status =  "Running"
		} else if pod.Status.ContainerStatuses[0].State.Waiting != nil {
			status = "Waiting " + pod.Status.ContainerStatuses[0].State.Waiting.Reason
		} else {
			status = "Unknown State ! "
		}
		//get pod nodename
		var nodename string
		nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
		if err != nil {
			fmt.Println("k8s get node err is ",err)
		}
		for _,node := range nodes.Items {
			if node.Status.Addresses[0].Address == pod.Status.HostIP {
				nodename = node.Name
			}
		}

		tmpod := nsPod{
			Name: pod.Name,
			Status: status,
			NodeIp: nodename,
			PodIp: pod.Status.PodIP,
			RestartCount: pod.Status.ContainerStatuses[0].RestartCount,
			CreateDate: pod.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}

		podSlce = append(podSlce,tmpod)
	}
	c.Data["nsPod"] = podSlce


}

func (c *KubernetesController) GetServiceData() {

	svcs,err := clientset.CoreV1().Services(c.Data["nsName"].(string)).List(metav1.ListOptions{})
	if err != nil {
		log.Println(err.Error())
	}
	svcSlice := []nsService{}

	for _,svc := range svcs.Items {
		var epSlice []string
		var labelstr []string
		for k,v := range svc.Spec.Selector {
			tmpstr := k + ": " + v
			labelstr = append(labelstr,tmpstr)
		}

		var podip string
		endpoint,err := clientset.CoreV1().Endpoints(c.Data["nsName"].(string)).List(metav1.ListOptions{})
		if err != nil {
			log.Println(err.Error())
		}
		nodename := ""
		for _,ep := range  endpoint.Items {


				if strings.TrimSpace(svc.Name) == strings.TrimSpace(ep.Name) && strings.TrimSpace(ep.Name) != "kubernetes" && len(ep.Subsets) > 0 {
					//后端podIP只取一个，如果多个replics，怎么处理？
					podip = ep.Subsets[0].Addresses[0].IP //获取service后端podIp

					nodename = *ep.Subsets[0].Addresses[0].NodeName

					break    //bug here
				} else {
					podip = "null"
				}
		}

		for _,p := range svc.Spec.Ports {
			if nodename != "" {
				epSlice = append(epSlice, svc.Spec.ClusterIP+":"+strconv.Itoa(int(p.Port))+"→"+podip+":"+p.TargetPort.String()+" ( "+nodename+" ) ")
			} else {
				//epSlice = append(epSlice, svc.Spec.ClusterIP+":"+strconv.Itoa(int(p.Port))+"→"+podip+":"+p.TargetPort.String())
				epSlice = append(epSlice, "null")
			}
		}
		tmsvc := nsService{
			Name: svc.Name,
			Label: labelstr,
			ClusterIp: svc.Spec.ClusterIP,
			DomainName:svc.Name + "." + svc.Namespace,
			EndPoints: epSlice,
			CreateDate: svc.CreationTimestamp.Format("2006-01-02 15:04:05"),

		}

		svcSlice = append(svcSlice,tmsvc)


	}

	c.Data["nsSvc"] = svcSlice


}





