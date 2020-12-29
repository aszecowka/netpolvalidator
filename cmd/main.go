package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/aszecowka/netpolvalidator/internal/netpol"
	"github.com/aszecowka/netpolvalidator/internal/ns"
	"github.com/aszecowka/netpolvalidator/internal/podcandidate"
	"github.com/aszecowka/netpolvalidator/internal/rule"
	"github.com/aszecowka/netpolvalidator/internal/state"
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

	nsService := ns.New(clientset.CoreV1().Namespaces())
	netpolService := netpol.NewService(clientset.NetworkingV1())
	deployService := podcandidate.NewDeploymentService(clientset.AppsV1())
	podCandidateProviders := make(map[string]state.PodCandidatesProvider)
	podCandidateProviders["deploy"] = deployService
	clusterStateBuilder := state.NewBuilder(nsService, netpolService, podCandidateProviders)
	clusterState, err := clusterStateBuilder.Build(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("ns", len(clusterState.Namespaces))
	fmt.Println("net pol", len(clusterState.NetworkPolicies))
	for namespaces, candidates := range clusterState.PodCandidates {
		fmt.Printf("ns: %s, candidates: %d\n", namespaces, len(candidates))
	}

	validators := make(map[string]rule.Validator)
	validators["label correctness"] = rule.NewLabelCorrectness()

	for name, validator := range validators {
		err := validator.Validate(*clusterState)
		if err != nil {
			fmt.Println(name, err)
		}
	}
}
