package k8s

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedAppsV1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	typedCoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func EndEnv(name string) error {
	deploymentClient := ClientSet.AppsV1().Deployments(corev1.NamespaceDefault)
	serviceClient := ClientSet.CoreV1().Services(corev1.NamespaceDefault)

	if err := deleteDeployment(deploymentClient, name); err != nil {
		return err
	}
	if err := deleteService(serviceClient, name+"-svc"); err != nil {
		return err
	}
	return nil
}

func deleteDeployment(deploymentClient typedAppsV1.DeploymentInterface, deploymentName string) error {
	fmt.Printf("Start deleting deployment %s...\n", deploymentName)
	err := deploymentClient.Delete(context.TODO(), deploymentName, metav1.DeleteOptions{})
	return err
}

func deleteService(serviceClient typedCoreV1.ServiceInterface, serviceName string) error {
	fmt.Printf("Start deleting service %s...\n", serviceName)
	err := serviceClient.Delete(context.TODO(), serviceName, metav1.DeleteOptions{})
	return err
}
