package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	defaultKubeconfig := os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
	if defaultKubeconfig == "" {
		defaultKubeconfig = clientcmd.RecommendedHomeFile
	}

	kubeconfig := flag.String(clientcmd.RecommendedConfigPathFlag, defaultKubeconfig, "Absolute path to the kubeconfig file.")

	flag.Parse()

	// 使用client-go的工具函数创建一个Kubernetes客户端
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// 创建一个Watcher，用于监视Pod资源对象的变更事件
	lw := cache.NewListWatchFromClient(
		clientset.CoreV1().RESTClient(),
		"pods",
		metav1.NamespaceAll, // 所有的namespace
		fields.Everything(),
	)

	// 创建一个Reflector，用于从Kubernetes API服务器获取Pod资源对象的列表，并将其进行本地缓存
	store := cache.NewStore(cache.MetaNamespaceKeyFunc)

	reflector := cache.NewReflector(lw, &corev1.Pod{}, store, 10*time.Second)
	fmt.Println(reflector.Name)

	// 启动Reflector，开始监听Kubernetes API服务器上Pod资源对象的变更事件
	stopCh := make(chan struct{})
	go reflector.Run(stopCh)

	var wg sync.WaitGroup
	wg.Add(1)
	// 测试：打印本地缓存中，缓存的一条 Key
	go func() {
		defer wg.Done()
		for {
			if len(store.ListKeys()) > 0 {
				fmt.Printf("Local store cached a key: %q\n", store.ListKeys()[0])
				return
			}
		}
	}()

	wg.Wait()
}
