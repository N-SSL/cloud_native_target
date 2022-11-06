package k8s

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	typedAppsV1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	typedCoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"log"
)


func StartEnvNow(deployFile []byte,serviceFile []byte, extName string) (string, int32, error) {
	fmt.Println("start ENV:" + extName)
	var deployment appsv1.Deployment
	var service  corev1.Service
	_ = json.Unmarshal(deployFile, &deployment)
	_ = json.Unmarshal(serviceFile, &service)

	fmt.Println("get config file")
	deploymentClient := ClientSet.AppsV1().Deployments(metav1.NamespaceDefault)
	serviceClient := ClientSet.CoreV1().Services(metav1.NamespaceDefault)

	fmt.Println("deploy")
	deployment = renameDeployment(deployment, extName)
	if err := createDeployment(deploymentClient, deployment); err != nil {
		log.Println(err)
		return "",0, err
	}

	fmt.Println("service")
	service = renameService(service, extName)
	ip,port, err := createService(serviceClient, service)
	if err != nil {
		return "",0, err
	}
	return ip, port, nil
}
func renameDeployment(deployment appsv1.Deployment, name string) appsv1.Deployment {
	deployment.ObjectMeta.Name = name
	deployment.Spec.Selector.MatchLabels = map[string]string{"app": name}
	deployment.Spec.Template.Labels = map[string]string{"app": name}
  oriname:="container.apparmor.security.beta.kubernetes.io/"+deployment.Spec.Template.Spec.Containers[0].Name
	deployment.Spec.Template.Spec.Containers[0].Name = name
	//deployment.Spec.Template.Annotations = strings.Replace(deployment.Spec.Template.ObjectMeta.Annotations, "cpt-gotty", name, 1)
	secname:="container.apparmor.security.beta.kubernetes.io/"+name
  fmt.Printf("Sec name: \n  %s \n", secname)
  fmt.Printf("Ori name: \n  %s \n", oriname)
  fmt.Printf("Annotation content is: \n  %s \n", deployment.Spec.Template.ObjectMeta.Annotations)
  _, ok := deployment.Spec.Template.ObjectMeta.Annotations[oriname]
  if ok {
      deployment.Spec.Template.ObjectMeta.Annotations=map[string]string{secname: "unconfined"}
  } else {
      fmt.Println("apparmor is not contained in annotations")
  }
	return deployment
}


func createDeployment(deploymentClient typedAppsV1.DeploymentInterface, deployment appsv1.Deployment) error {
	if _, err := deploymentClient.Get(context.TODO(), deployment.Name, metav1.GetOptions{}); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		fmt.Printf("Start creating deployment...\n")
    fmt.Printf("Deployment content is: \n  %s \n", deployment)
		result, err := deploymentClient.Create(context.TODO(), &deployment, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		fmt.Printf("Deployment %s created!\n", result.Name)
	} else {
		fmt.Printf("Start updating deployment...\n")
		if _, err = deploymentClient.Update(context.TODO(), &deployment, metav1.UpdateOptions{}); err != nil {
			return err
		}
		fmt.Printf("Deployment %s updated!\n", deployment.Name)
	}
	return nil
}

func createService(serviceClient typedCoreV1.ServiceInterface, service corev1.Service) (string, int32, error) {
	if _, err := serviceClient.Get(context.TODO(), service.Name, metav1.GetOptions{}); err != nil {
		if !errors.IsNotFound(err) {
			return "",0, err
		}
		fmt.Printf("Start creating service...\n")
		result, err := serviceClient.Create(context.TODO(), &service, metav1.CreateOptions{})
		if err != nil {
			return "",0, err
		}
		fmt.Printf("Service %s created!\n", result.GetObjectMeta().GetName())
		return result.Spec.ClusterIP,result.Spec.Ports[0].Port, nil
	} else {
		fmt.Printf("Start updating service...\n")
		var result *corev1.Service
		if result, err = serviceClient.Update(context.TODO(), &service, metav1.UpdateOptions{}); err != nil {
			return "", 0, err
		}
		fmt.Printf("Service %s updated!\n", result.Name)
		return result.Spec.ClusterIP,result.Spec.Ports[0].Port, nil
	}
}

func renameService(svc corev1.Service, name string) corev1.Service {
	svc.Spec.Selector = map[string]string{"app": name}
	name += "-svc"
	svc.ObjectMeta.Name = name
	return svc
}

