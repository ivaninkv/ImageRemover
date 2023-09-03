package kube

import (
	"context"
	"imageRemover/config"
	"imageRemover/logger"
	"imageRemover/output"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func GetImages(cfg config.Config) map[string]bool {
	logger.Log.Debug().Msg("Getting images from deployments")

	err, clientset := createClientset(cfg.KubeCluster.ServerUrl, cfg.KubeCluster.Token)
	if err != nil {
		panic(err)
	}

	// Мапа для хранения образов
	images := make(map[string]bool)

	// Сохранение деплойментов в срез
	nsList, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, ns := range nsList.Items {
		logger.Log.Debug().Str("namespace", ns.Name).Msg("Handling namespace:")
		deploymentList, err := clientset.AppsV1().Deployments(ns.Name).List(context.TODO(), metav1.ListOptions{})
		if err == nil {
			for _, dpl := range deploymentList.Items {
				logger.Log.Debug().Str("deployment", dpl.Name).Msg("Handling deployment:")
				for _, container := range dpl.Spec.Template.Spec.Containers {
					logger.Log.Debug().Str("container", container.Image).Msg("Handling container:")
					images[container.Image] = true
				}
			}
		}
	}

	if cfg.Output.WriteToTXT {
		output.WriteToTXT(cfg.Output.KubeFileName, images)
	}
	return images
}

func createClientset(serverURL string, token string) (error, *kubernetes.Clientset) {
	// Создание конфигурации OpenShift с адресом сервера и токеном
	cfg, err := clientcmd.BuildConfigFromFlags(serverURL, "")
	if err != nil {
		panic(err.Error())
	}
	cfg.BearerToken = token
	cfg.Insecure = true // Если требуется отключить проверку сертификата сервера

	// Создание клиента OpenShift
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err.Error())
	}
	return err, clientset
}
