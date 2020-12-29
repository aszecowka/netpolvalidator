package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()

	nsList, err := clientset.CoreV1().Namespaces().List(ctx, v1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("namespaces", len(nsList.Items))

	netPolicies, err := clientset.NetworkingV1().NetworkPolicies("default").List(ctx, v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("net policies in default ns", len(netPolicies.Items))

	// Deployments
	deployments, err := clientset.AppsV1().Deployments("ts").List(ctx, v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	// deployments.Items[0].Spec.Template == PodTemplateSpec
	objMetaString := deployments.Items[0].ObjectMeta
	fmt.Println("Obj meta string: ", objMetaString)

}
