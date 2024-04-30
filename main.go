package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/zhengyansheng/scheduler-extender/pkg/scheduler"
)

var (
	kubeconfig = flag.String("kubeconfig", "~/.kube/config", "paths to a kubeconfig")
	port       = flag.Int("port", 8000, "port is the port that the scheduler server serves at")
)

func main() {
	flag.Parse()

	s, err := scheduler.NewScheduleExtender(*kubeconfig)
	if err != nil {
		log.Fatalf("Failed to new schedule extender controller: %s", err)
	}

	s.Run(fmt.Sprintf("0.0.0.0:%d", *port))
}
