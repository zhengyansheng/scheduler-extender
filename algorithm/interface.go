package algorithm

import (
	"fmt"
	"math/rand"

	"k8s.io/klog/v2"
	extenderv1 "k8s.io/kube-scheduler/extender/v1"
)

type Extender interface {
	Ping() bool
	Filter(extenderv1.ExtenderArgs) *extenderv1.ExtenderFilterResult
	Score(extenderv1.ExtenderArgs) *extenderv1.HostPriorityList
}

type extender struct {
}

func NewExtender() Extender {
	return &extender{}
}

func (e *extender) Ping() bool {
	return true
}

func (e *extender) Filter(args extenderv1.ExtenderArgs) *extenderv1.ExtenderFilterResult {
	pod := args.Pod
	if pod == nil {
		return &extenderv1.ExtenderFilterResult{Error: fmt.Sprintf("pod is nil")}
	}
	scheduleNodes := make([]string, 0)
	failedNodes := make(map[string]string)
	for _, nodeName := range *args.NodeNames {
		scheduleNodes = append(scheduleNodes, nodeName)
	}

	klog.Infof("filter localstorage pods on nodes: %v", scheduleNodes)
	return &extenderv1.ExtenderFilterResult{
		NodeNames:   &scheduleNodes,
		Nodes:       args.Nodes,
		FailedNodes: failedNodes,
	}
}

func (e *extender) Score(args extenderv1.ExtenderArgs) *extenderv1.HostPriorityList {
	pod := args.Pod
	if pod == nil {
		klog.Errorf("pod is nil")
		return nil
	}

	nodeNames := *args.NodeNames
	klog.Infof("scoring nodes %v", nodeNames)

	hostPriorityList := make(extenderv1.HostPriorityList, len(nodeNames))

	for i, nodeName := range nodeNames {
		hostPriorityList[i] = extenderv1.HostPriority{
			Host:  nodeName,
			Score: int64(rand.Intn(10)),
		}
	}

	klog.Infof("score localstorage pods on nodes: %v", hostPriorityList)
	return &hostPriorityList

}
