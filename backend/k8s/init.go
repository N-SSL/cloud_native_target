package k8s

import (
	"flag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
	"path/filepath"
)

var (
	ClientSet        *kubernetes.Clientset
	MetricsClientSet *metricsv.Clientset
)

func InitClient() {
	var k8sConfig *string
	if home := homedir.HomeDir(); home != "" {
		k8sConfig = flag.String("k8sConfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the k8sConfig file")
	} else {
		k8sConfig = flag.String("k8sConfig", "", "absolute path to the k8sConfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *k8sConfig)
	if err != nil {
		panic(err)
	}
	ClientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	MetricsClientSet, err = metricsv.NewForConfig(config)
	if err != nil {
		panic(err)
	}
}








