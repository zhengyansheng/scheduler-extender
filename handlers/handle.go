package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhengyansheng/common"
	"github.com/zhengyansheng/scheduler-extender/algorithm"
	extenderv1 "k8s.io/kube-scheduler/extender/v1"
)

func PingHandle(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}

func FilterHandle(c *gin.Context) {
	var extenderArgs extenderv1.ExtenderArgs
	var extenderFilterResult *extenderv1.ExtenderFilterResult

	if err := c.ShouldBindJSON(&extenderArgs); err != nil {
		extenderFilterResult = &extenderv1.ExtenderFilterResult{Error: err.Error()}
	} else {
		p := algorithm.NewExtender()
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

func ScoreHandle(c *gin.Context) {
	var extenderArgs extenderv1.ExtenderArgs
	var hostPriorityList *extenderv1.HostPriorityList

	if err := c.ShouldBindJSON(&extenderArgs); err != nil {
		hostPriorityList = &extenderv1.HostPriorityList{}
	} else {
		p := algorithm.NewExtender()
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
