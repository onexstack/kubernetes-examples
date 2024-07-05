package main

import (
	"flag"
	"fmt"
	"os"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

func main() {
	defaultKubeconfig := os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
	if len(defaultKubeconfig) == 0 {
		defaultKubeconfig = clientcmd.RecommendedHomeFile
	}

	kubeconfig := flag.String(
		clientcmd.RecommendedConfigPathFlag,
		defaultKubeconfig,
		"Absolute path to the kubeconfig file.",
	)
	flag.Parse()

	// creates the connection
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		klog.Fatal(err)
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}

	// create the pod watcher
	podListWatcher := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", v1.NamespaceDefault, fields.Everything())

	// Bind the workqueue to a cache with the help of an informer. This way we make sure that
	// whenever the cache is updated, the pod key is added to the workqueue.
	// Note that when we finally process the item from the workqueue, we might see a newer version
	// of the Pod than the version which was responsible for triggering the update.
	indexer, informer := cache.NewIndexerInformer(podListWatcher, &v1.Pod{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println("New object added:", obj.(*v1.Pod).Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Println("Object updated. Old:", oldObj.(*v1.Pod).Name, "New:", newObj.(*v1.Pod).Name)
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("Object deleted:", obj.(*v1.Pod).Name)
		},
	}, cache.Indexers{})

	// 启动 Informer
	stopCh := make(chan struct{})
	defer close(stopCh)
	go informer.Run(stopCh)
	if !cache.WaitForCacheSync(stopCh, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}

	// Cache 缓存完，执行业务逻辑
	go func() {
		// 不断轮训 Indexer，当有数据时打印一条数据，并退出
		for {
			if len(indexer.ListKeys()) > 0 {
				obj, _, err := indexer.GetByKey(indexer.ListKeys()[0])
				if err != nil {
					panic(err)
				}
				accessor, _ := meta.Accessor(obj)
				fmt.Printf("Resource name: %s\n", accessor.GetName())
				return
			}
		}
	}()

	// Wait forever
	select {}
}
