package main

import (
	"path/filepath"

	"github.com/spf13/pflag"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	defaultKubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")

	kubeconfig = pflag.StringP("kubeconfig", "c", defaultKubeconfig, "Absolute path to the kubeconfig file.")
)

// 获取kube-system命名空间下的Pod列表
func main() {
	pflag.Parse()

	// 创建 *restclient.Config 类型的 REST 客户端配置
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	// 创建 DiscoveryClient
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err)
	}
	_ = discoveryClient

	// 创建 DynamicClient
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	_ = dynamicClient

	// 创建 ClientSet
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	_ = clientset

	// 创建 RESTClient
	// 注意：创建 RESTClient时，要对 *restclient.Config 进行必要的设置
	// 配置API路径, 请求的HTTP路径
	config.APIPath = "api" // pods: /api/v1/pods
	// config.APIPath = "apis" // deployments: /apis/apps/v1/namespaces/{namespace}/deployments/{deployment}

	// 配置请求的资源组/资源版本
	config.GroupVersion = &corev1.SchemeGroupVersion // 无名资源组, group: " ", version: "v1"

	// 配置数据的编解码工具（序列化和反序列化）
	config.NegotiatedSerializer = scheme.Codecs

	// 创建 RESTClient
	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err)
	}
	_ = restClient
}
