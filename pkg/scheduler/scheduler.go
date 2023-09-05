package scheduler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhengyansheng/common"
	"github.com/zhengyansheng/scheduler-extender/algorithm"
	"k8s.io/client-go/informers"
	v1 "k8s.io/client-go/listers/core/v1"
	extenderv1 "k8s.io/kube-scheduler/extender/v1"
)

const (
	apiPrefix = "/scheduler/extender"
)

type ScheduleExtender struct {
	http       *gin.Engine
	nodeLister v1.NodeLister
}

func NewScheduleExtender(factory informers.SharedInformerFactory) (*ScheduleExtender, error) {
	s := &ScheduleExtender{
		http: gin.New(),
	}

	s.http.GET("/healthz", s.Ping)
	s.http.POST(apiPrefix+"/filter", s.FilterHandle)
	s.http.POST(apiPrefix+"/prioritize", s.ScoreHandle)

	s.nodeLister = factory.Core().V1().Nodes().Lister()

	return s, nil
}

func (s *ScheduleExtender) Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func (s *ScheduleExtender) FilterHandle(c *gin.Context) {
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

func (s *ScheduleExtender) ScoreHandle(c *gin.Context) {
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

func (s *ScheduleExtender) Run(addr ...string) {
	s.http.Run(addr...)
}
