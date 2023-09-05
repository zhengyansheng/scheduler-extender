package util

import (
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
)

func BuildRestConfig(configFile string) (*rest.Config, error) {
	if len(configFile) != 0 {
		klog.Infof("kubeconfig specified. building kube config from that")
		return clientcmd.BuildConfigFromFlags("", configFile)
	}

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err == nil {
		klog.Infof("kubeconfig not specified. try to building kube config from ~/.kube/config")
		return kubeConfig, nil
	}

	klog.Infof("Building kube configs for running in cluster...")
	return rest.InClusterConfig()
}

func NewClientSet(configFile string) (kubernetes.Interface, error) {
	restConfig, err := BuildRestConfig(configFile)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(restConfig)

}
