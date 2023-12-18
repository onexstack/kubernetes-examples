package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
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

	rc, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create a client set from config
	clientset, err := kubernetes.NewForConfig(rc)
	if err != nil {
		panic(err.Error())
	}

	// create a new instance of sharedInformerFactory for all namespaces
	factory := informers.NewSharedInformerFactory(clientset, 5*time.Second)

	// using this factory create an informer for `configmap` resources
	cmInformer := factory.Core().V1().ConfigMaps()

	// adds an event handler to the shared informer
	cmInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			cm := obj.(*corev1.ConfigMap)
			fmt.Printf("Informer event: ConfigMap ADDED %s/%s\n", cm.GetNamespace(), cm.GetName())
		},
		UpdateFunc: func(old, new interface{}) {
			cm := old.(*corev1.ConfigMap)
			fmt.Printf("Informer event: ConfigMap UPDATED %s/%s\n", cm.GetNamespace(), cm.GetName())
		},
		DeleteFunc: func(obj interface{}) {
			cm := obj.(*corev1.ConfigMap)
			fmt.Printf("Informer event: ConfigMap DELETED %s/%s\n", cm.GetNamespace(), cm.GetName())
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// starts the shared informers that have been created by the factory
	factory.Start(ctx.Done())

	// wait for the initial synchronization of the local cache
	if !cache.WaitForCacheSync(ctx.Done(), cmInformer.Informer().HasSynced) {
		panic("failed to sync")
	}

	// causes the goroutine to block (hit CTRL+C to exit)
	select {}
}
