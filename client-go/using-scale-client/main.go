package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/scale"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute pathto the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err)
	}
	gr, err := restmapper.GetAPIGroupResources(dc)
	if err != nil {
		panic(err)
	}
	mapper := restmapper.NewDiscoveryRESTMapper(gr)

	resolver := scale.NewDiscoveryScaleKindResolver(dc)

	client, err := scale.NewForConfig(config, mapper, dynamic.LegacyAPIPathResolverFunc, resolver)
	if err != nil {
		panic(err)
	}

	resource := schema.GroupResource{Group: "apps", Resource: "deployments"}

	result, err := client.Scales(corev1.NamespaceDefault).
		Get(context.TODO(), resource, "nginx", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deployment %s has %d replica(s).\n", result.GetName(), result.Spec.Replicas)

	result.Spec.Replicas = 2
	updated, err := client.Scales(corev1.NamespaceDefault).
		Update(context.TODO(), resource, result, metav1.UpdateOptions{})
	fmt.Printf("Deployment %s replicas updated to %d.\n", updated.GetName(), updated.Spec.Replicas)
}
