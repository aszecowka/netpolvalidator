package main

import (
	"context"
	"fmt"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/aszecowka/netpolvalidator/internal"
	"github.com/aszecowka/netpolvalidator/internal/model"
	"github.com/aszecowka/netpolvalidator/internal/netpol"
	"github.com/aszecowka/netpolvalidator/internal/ns"
	"github.com/aszecowka/netpolvalidator/internal/podcandidate"
	"github.com/aszecowka/netpolvalidator/internal/rule"
	"github.com/aszecowka/netpolvalidator/internal/state"
)

func main() {
	cfg, err := internal.Load()
	if err != nil {
		panic(err)
	}

	generateReport(cfg)
}

func generateReport(cfg internal.Config) {

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", cfg.Kubeconfig)
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
	podCandidateProviders := make(map[string]state.PodCandidatesProvider)
	podCandidateProviders["cronjob"] = podcandidate.NewCronjobFetcher(clientset.BatchV1beta1())
	podCandidateProviders["daemonset"] = podcandidate.NewDaemonsetFetcher(clientset.AppsV1())
	podCandidateProviders["deployment"] = podcandidate.NewDeploymentsFetcher(clientset.AppsV1())
	podCandidateProviders["job"] = podcandidate.NewJobFetcher(clientset.BatchV1())
	podCandidateProviders["pod"] = podcandidate.NewPodsFetcher(clientset.CoreV1())
	podCandidateProviders["statefulset"] = podcandidate.NewStatefulsetsFetcher(clientset.AppsV1())

	clusterStateBuilder := state.NewBuilder(nsService, netpolService, podCandidateProviders)
	clusterState, err := clusterStateBuilder.Build(ctx)
	if err != nil {
		panic(err)
	}
	for namespaces, candidates := range clusterState.PodCandidates {
		fmt.Printf("ns: %s, candidates: %d\n", namespaces, len(candidates))
	}

	validators := make(map[string]rule.Validator)
	validators["label correctness"] = rule.NewLabelCorrectness()

	var allViolations []model.Violation
	for _, validator := range validators {
		violations, err := validator.Validate(*clusterState)
		if err != nil {
			panic(err)
		}
		allViolations = append(allViolations, violations...)
	}

	fmt.Printf("Found %d violations\n", len(allViolations))
	for _, v := range allViolations {
		fmt.Println(v)
	}
}
