package main

import (
	"context"
	"flag"
	"path/filepath"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// 1. 导入 client-go
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// 2. 创建 Kubernetes 客户端
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// 3. 编写业务逻辑：这里删除名为 `demo-deployment` 的 Deployment 资源
	if err := clientset.AppsV1().Deployments(apiv1.NamespaceDefault).Delete(context.Background(), "demo-deployment", metav1.DeleteOptions{}); err != nil {
		panic(err)
	}
}
