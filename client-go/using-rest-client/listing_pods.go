package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/pflag"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	defaultKubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")

	kubeconfig = pflag.StringP("kubeconfig", "c", defaultKubeconfig, "Absolute path to the kubeconfig file.")
)

func main() {
	pflag.Parse()

	// 1. 加载配置文件，生成config对象
	cfg, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	// the base API path "/api" (legacy resource)
	cfg.APIPath = "api"
	// the Pod group and version "/v1" (group name is empty for legacy resources)
	cfg.GroupVersion = &corev1.SchemeGroupVersion
	// specify the serializer
	cfg.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	// create a RESTClient instance, using the the
	// configuration object as input parameter
	rc, err := rest.RESTClientFor(cfg)
	if err != nil {
		panic(err.Error())
	}

	// the list of Pods (the result)
	res := &corev1.PodList{}

	// fluent interface to setup and perform
	// the GET /api/v1/pods request
	err = rc.Get().
		Namespace(metav1.NamespaceSystem).
		Resource("pods").
		Do(context.TODO()).
		Into(res)
	if err != nil {
		panic(err.Error())
	}

	// print the results on the terminal in the form of a table
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 5, 0, 3, ' ', 0)

	dorow := func(cells []string) {
		fmt.Fprintln(w, strings.Join(cells, "\t"))
	}

	dorow([]string{"NAME", "STATUS", "AGE"})

	for _, p := range res.Items {
		age := time.Since(p.CreationTimestamp.Time).Round(time.Second)
		dorow([]string{p.Name, string(p.Status.Phase), fmt.Sprintf("%dm", int(age.Minutes()))})
	}

	w.Flush()
}
