package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhengyansheng/common"
	"github.com/zhengyansheng/scheduler-extender/algorithm"
	"github.com/zhengyansheng/scheduler-extender/pkg/util"
	"k8s.io/client-go/informers"
	v1 "k8s.io/client-go/listers/core/v1"
	extenderv1 "k8s.io/kube-scheduler/extender/v1"
)

const (
	apiPrefix = "/scheduler/extender"
)

type extender struct {
	http       *gin.Engine
	nodeLister v1.NodeLister
}

func initInformerFactory(kubeconfig string) (informers.SharedInformerFactory, error) {
	// Create a clientset
	clientSet, err := util.NewClientSet(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("new clientset: %v", err)
	}

	// Create a shared informer factory
	factory := informers.NewSharedInformerFactory(clientSet, 0)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the informer factory
	factory.Start(ctx.Done())
	factory.WaitForCacheSync(ctx.Done())

	return factory, nil

}

func NewScheduleExtender(kubeconfig string) (*extender, error) {
	// Create a new extender
	s := &extender{http: gin.New()}

	// Register the routes
	s.http.GET("/healthz", s.Ping)
	s.http.POST(apiPrefix+"/filter", s.FilterHandle)
	s.http.POST(apiPrefix+"/prioritize", s.ScoreHandle)

	// Create a shared informer factory
	factory, err := initInformerFactory(kubeconfig)
	if err != nil {
		return nil, err
	}
	s.nodeLister = factory.Core().V1().Nodes().Lister()

	// Return the extender
	return s, nil
}

func (s *extender) Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func (s *extender) FilterHandle(c *gin.Context) {
	var extenderArgs extenderv1.ExtenderArgs
	var extenderFilterResult *extenderv1.ExtenderFilterResult

	if err := c.ShouldBindJSON(&extenderArgs); err != nil {
		extenderFilterResult = &extenderv1.ExtenderFilterResult{Error: err.Error()}
	} else {
		p := algorithm.NewExtender(s.nodeLister)
		extenderFilterResult = p.Filter(extenderArgs)
	}
	common.Indent(extenderFilterResult)
	c.Header("Content-Type", "application/json")
	result, err := json.Marshal(extenderFilterResult)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, string(result))
}

func (s *extender) ScoreHandle(c *gin.Context) {
	var extenderArgs extenderv1.ExtenderArgs
	var hostPriorityList *extenderv1.HostPriorityList

	if err := c.ShouldBindJSON(&extenderArgs); err != nil {
		hostPriorityList = &extenderv1.HostPriorityList{}
	} else {
		p := algorithm.NewExtender(s.nodeLister)
		hostPriorityList = p.Score(extenderArgs)
	}
	common.Indent(hostPriorityList)
	c.Header("Content-Type", "application/json")
	result, err := json.Marshal(hostPriorityList)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, string(result))
	return
}

func (s *extender) Run(addr ...string) {
	s.http.Run(addr...)
}
