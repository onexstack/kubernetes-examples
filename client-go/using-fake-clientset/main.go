package main

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func main() {
	// utility function to create a int32 pointer
	i32Ptr := func(i int32) *int32 { return &i }

	// the required request body (a deployment object)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: i32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:1.21.6",
						},
					},
				},
			},
		},
	}

	cs := fake.NewSimpleClientset()

	// create the deployment in the specified namespace
	res, err := cs.AppsV1().Deployments(metav1.NamespaceDefault).
		Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("deployment.apps/%s created\n", res.Name)

	deploy, err := cs.AppsV1().Deployments(metav1.NamespaceDefault).Get(context.TODO(), deployment.Name, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("get deployment.apps/%s\n", deploy.Name)
}
