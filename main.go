package main

import (
	"context"
	"flag"
	"fmt"
	"go-k8s-definitions/models"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	clusterDefinitions := make(map[string]*models.ClusterDefinition)

	cluster1Definition := getClusterDefinition("cluster1")
	clusterDefinitions["cluster1"] = cluster1Definition

	for name, definition := range clusterDefinitions {
		fmt.Printf("%s : %+v \n", name, definition)
	}
}

func getClusterDefinition(clusterName string) *models.ClusterDefinition {
	kubeConfig := flag.String(fmt.Sprintf("%s-kubeconfig", clusterName), fmt.Sprintf("./configs/%s.conf", clusterName), "")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &models.ClusterDefinition{
		Name:       clusterName,
		Pods:       getPods(clientset),
		Nodes:      getNodes(clientset),
		Namespaces: getNamespaces(clientset),
		Endpoints:  getEndpoints(clientset),
	}
}

func getPods(clientset *kubernetes.Clientset) []*models.PodDefinition {
	var podDefinitions []*models.PodDefinition
	pods, err := clientset.CoreV1().Pods("logistics-matching").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, pod := range pods.Items {
		podDefinitions = append(podDefinitions, &models.PodDefinition{Name: pod.Name, IP: pod.Status.PodIP})
	}
	return podDefinitions
}

func getNamespaces(clientset *kubernetes.Clientset) []*models.NamespaceDefinition {
	var namespaceDefinitions []*models.NamespaceDefinition
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, namespace := range namespaces.Items {
		namespaceDefinitions = append(namespaceDefinitions, &models.NamespaceDefinition{Name: namespace.Name})
	}
	return namespaceDefinitions
}

func getNodes(clientset *kubernetes.Clientset) []*models.NodeDefinition {
	var nodeDefinitions []*models.NodeDefinition
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, node := range nodes.Items {
		nodeDefinitions = append(nodeDefinitions, &models.NodeDefinition{Name: node.Name})
	}
	return nodeDefinitions
}

func getEndpoints(clientset *kubernetes.Clientset) map[string][]*models.SubsetAddressDefinition {
	endpointDefinitions := make(map[string][]*models.SubsetAddressDefinition)
	endpoints, err := clientset.CoreV1().Endpoints("namespace").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, endpoint := range endpoints.Items {
		var subSetAddresses []*models.SubsetAddressDefinition
		for _, subset := range endpoint.Subsets {
			for _, subsetAddress := range subset.Addresses {
				subSetAddresses = append(subSetAddresses, &models.SubsetAddressDefinition{IP: subsetAddress.IP})
			}
		}
		endpointDefinitions[endpoint.Name] = subSetAddresses
	}
	return endpointDefinitions
}
