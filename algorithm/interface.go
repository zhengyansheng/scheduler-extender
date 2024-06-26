package algorithm

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	v1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/klog/v2"
	extenderv1 "k8s.io/kube-scheduler/extender/v1"
)

type Extender interface {
	Filter(extenderv1.ExtenderArgs) *extenderv1.ExtenderFilterResult
	Score(extenderv1.ExtenderArgs) *extenderv1.HostPriorityList
}

type extender struct {
	nodeLister v1.NodeLister
}

func NewExtender(nodeLister v1.NodeLister) Extender {
	return &extender{nodeLister: nodeLister}
}

func (e *extender) Filter(args extenderv1.ExtenderArgs) *extenderv1.ExtenderFilterResult {
	pod := args.Pod
	if pod == nil {
		return &extenderv1.ExtenderFilterResult{Error: fmt.Sprintf("pod is nil")}
	}
	scheduleNodes := make([]string, 0)
	failedNodes := make(map[string]string)
	for _, nodeName := range *args.NodeNames {
		if strings.Contains(nodeName, "control") {
			continue
		}
		scheduleNodes = append(scheduleNodes, nodeName)
	}

	klog.Infof("filter pods on nodes: %v", scheduleNodes)
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
		var scoreValue = int64(rand.Intn(10))
		node, err := e.nodeLister.Get(nodeName)
		if err != nil {
			klog.Errorf("get node %s error: %v", nodeName, err)
			continue
		}
		if annotations := node.GetAnnotations(); annotations != nil {
			_, ok := annotations["score"]
			if ok {
				scoreValue = 11
			}
		}
		if labels := node.GetLabels(); labels != nil {
			if value, ok := labels["score"]; ok {
				atoi, err := strconv.Atoi(value)
				if err != nil {
					klog.Errorf("convert score %s to int error: %v", value, err)
					scoreValue = 0
				} else {
					scoreValue = int64(atoi)
				}
			}
		}

		hostPriorityList[i] = extenderv1.HostPriority{
			Host:  nodeName,
			Score: scoreValue,
		}
	}

	klog.Infof("score pods on nodes: %+v", hostPriorityList)
	return &hostPriorityList

}
