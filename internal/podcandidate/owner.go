package podcandidate

import "fmt"

const (
	WorkloadDeployment  WorkloadType = "deployment"
	WorkloadCronjob     WorkloadType = "cronjob"
	WorkloadJob         WorkloadType = "job"
	WorkloadStatefulset WorkloadType = "statefulset"
	WorkloadDaemonset   WorkloadType = "daemonset"
	WorkloadPod         WorkloadType = "pod"
)

type WorkloadType string

func getOwnerName(t WorkloadType, ns, name string) string {
	return fmt.Sprintf("%s/%s/%s", t, ns, name)
}
