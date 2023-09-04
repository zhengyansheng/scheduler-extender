package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zhengyansheng/scheduler-extender/handlers"
)

const (
	apiPrefix = "/scheduler/extender"
)

var (
	port = flag.Int("port", 8000, "port is the port that the scheduler server serves at")
)

func main() {
	r := gin.Default()
	r.GET("/healthz", handlers.PingHandle)
	r.POST(apiPrefix+"/filter", handlers.FilterHandle)
	r.POST(apiPrefix+"/prioritize", handlers.ScoreHandle)

	log.Fatal(r.Run(":" + string(*port)))
}
