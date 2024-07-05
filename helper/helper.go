package helper

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"k8s.io/client-go/tools/clientcmd"
)

func Prompt() {
	fmt.Printf("-> Press Return key to continue.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println()
}

func AddKubeconfigFlag() string {
	defaultKubeconfig := os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
	if defaultKubeconfig == "" {
		defaultKubeconfig = clientcmd.RecommendedHomeFile
	}

	kubeconfig := flag.String(clientcmd.RecommendedConfigPathFlag, defaultKubeconfig, "Absolute path to the kubeconfig file.")

	flag.Parse()

	return *kubeconfig
}
