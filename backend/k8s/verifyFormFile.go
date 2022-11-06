package k8s

import (
	"bytes"
	"errors"
	"io"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/util/yaml"
	"mime/multipart"
	"strings"
)

func VerifyDeploy(FileHeader *multipart.FileHeader) ([]byte,error) {
	deployData,err := verifyDeploymentFile(FileHeader)
	if err != nil{
		return nil,err
	}
	deployBytes,err2 := json.Marshal(deployData)
	if err2!=nil{
		return nil,err2
	}
	return deployBytes,err2
}

func verifyDeploymentFile(FileHeader *multipart.FileHeader )(appsv1.Deployment,error){
	deployData, derr1 := FileHeader.Open()
	defer deployData.Close()
	if derr1 !=nil {
		return appsv1.Deployment{},derr1
	}
	deployBuf :=bytes.NewBuffer(nil)
	_ , derr2 :=io.Copy(deployBuf,deployData)

	if derr2 !=nil {
		return appsv1.Deployment{},derr2
	}
	var deployment appsv1.Deployment
	if strings.Contains(FileHeader.Filename,".yaml") {
		deployJson,err := yaml.ToJSON(deployBuf.Bytes())
		if err !=nil {
			return appsv1.Deployment{},err
		}

		if err := json.Unmarshal(deployJson, &deployment); err != nil {
			return appsv1.Deployment{}, err
		}
		return deployment, nil
	}
	if strings.Contains(FileHeader.Filename,".json") {
		if err := json.Unmarshal(deployBuf.Bytes(), &deployment); err != nil {
			return appsv1.Deployment{}, err
		}
		return deployment, nil
	}
	return appsv1.Deployment{},errors.New("no correct Suffix")
}

func VerifyService(FileHeader *multipart.FileHeader) ([]byte,error) {
	serviceData,err := verifyServiceFile(FileHeader)
	if err != nil{
		return nil,err
	}
	serviceBytes,err2 := json.Marshal(serviceData)
	if err2!=nil{
		return nil,err2
	}
	return serviceBytes,err2
}

func verifyServiceFile(FileHeader *multipart.FileHeader )(corev1.Service,error){
	serviceData, serr1 := FileHeader.Open()
	defer serviceData.Close()
	if serr1 !=nil {
		return corev1.Service{},serr1
	}
	serviceBuf :=bytes.NewBuffer(nil)
	_ , serr2 :=io.Copy(serviceBuf,serviceData)

	if serr2 !=nil {
		return corev1.Service{},serr2
	}
	var service corev1.Service
	if strings.Contains(FileHeader.Filename,".yaml") {
		serviceJson,err := yaml.ToJSON(serviceBuf.Bytes())
		if err !=nil {
			return corev1.Service{},err
		}

		if err := json.Unmarshal(serviceJson, &service); err != nil {
			return corev1.Service{}, err
		}
		return service, nil
	}
	if strings.Contains(FileHeader.Filename,".json") {
		if err := json.Unmarshal(serviceBuf.Bytes(), &service); err != nil {
			return corev1.Service{}, err
		}
		return service, nil
	}
	return corev1.Service{},errors.New("no correct Suffix")
}

