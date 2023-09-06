package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/zhengyansheng/scheduler-extender/pkg/scheduler"
	"github.com/zhengyansheng/scheduler-extender/pkg/util"
	"k8s.io/client-go/informers"
)

var (
	kubeconfig = flag.String("kubeconfig", "", "paths to a kubeconfig")
	port       = flag.Int("port", 8000, "port is the port that the scheduler server serves at")
)

func main() {
	flag.Parse()

	clientSet, err := util.NewClientSet(*kubeconfig)
	if err != nil {
		log.Fatalf("Failed to build clientset: %v", err)
	}

	informerFactory := informers.NewSharedInformerFactory(clientSet, 0)

	s, err := scheduler.NewScheduleExtender(informerFactory)
	if err != nil {
		log.Fatalf("Failed to new schedule extender controller: %s", err)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	informerFactory.Start(ctx.Done())
	informerFactory.WaitForCacheSync(ctx.Done())

	s.Run(fmt.Sprintf("0.0.0.0:%d", *port))
}
