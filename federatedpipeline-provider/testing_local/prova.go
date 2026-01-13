package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	//connecting to k8s cluster
	config := &rest.Config{
		Host: "http://dtazzioli-edge.cloudmmwunibo.it:8080",
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("Error in creating clientset")
		fmt.Println(err)
		return
	}
	//get deployment
	deployment, err := clientset.AppsV1().Deployments("default").Get(context.TODO(), "flclient-deployment", metav1.GetOptions{})
	if err != nil {
		fmt.Println("Error in getting deployment")
		fmt.Println(err)
		return
	}
	fmt.Println(deployment)


}
