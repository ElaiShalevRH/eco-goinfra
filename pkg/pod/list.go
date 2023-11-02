package pod

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// List returns pod inventory in the given namespace.
func List(apiClient *clients.Settings, nsname string, options ...v1.ListOptions) ([]*Builder, error) {
	if nsname == "" {
		glog.V(100).Infof("pod 'nsname' parameter can not be empty")

		return nil, fmt.Errorf("failed to list pods, 'nsname' parameter is empty")
	}

	logMessage := fmt.Sprintf("Listing pods in the nsname %s", nsname)
	passedOptions := v1.ListOptions{}

	if len(options) == 1 {
		passedOptions = options[0]
		logMessage += fmt.Sprintf(" with the options %v", passedOptions)
	} else if len(options) > 1 {
		glog.V(100).Infof("'options' parameter must be empty or single-valued")

		return nil, fmt.Errorf("error: more than one ListOptions was passed")
	}

	glog.V(100).Infof(logMessage)

	podList, err := apiClient.Pods(nsname).List(context.Background(), passedOptions)

	if err != nil {
		glog.V(100).Infof("Failed to list pods in the nsname %s due to %s", nsname, err.Error())

		return nil, err
	}

	var podObjects []*Builder

	for _, runningPod := range podList.Items {
		copiedPod := runningPod
		podBuilder := &Builder{
			apiClient:  apiClient,
			Object:     &copiedPod,
			Definition: &copiedPod,
		}

		podObjects = append(podObjects, podBuilder)
	}

	return podObjects, nil
}

// ListInAllNamespaces returns a cluster-wide pod inventory.
func ListInAllNamespaces(apiClient *clients.Settings, options ...v1.ListOptions) ([]*Builder, error) {
	logMessage := "Listing all pods in all namespaces"
	passedOptions := v1.ListOptions{}

	if len(options) == 1 {
		passedOptions = options[0]
		logMessage += fmt.Sprintf(" with the options %v", passedOptions)
	} else if len(options) > 1 {
		glog.V(100).Infof("'options' parameter must be empty or single-valued")

		return nil, fmt.Errorf("error: more than one ListOptions was passed")
	}

	glog.V(100).Infof(logMessage)

	podList, err := apiClient.Pods("").List(context.Background(), passedOptions)

	if err != nil {
		glog.V(100).Infof("Failed to list all pods due to %s", err.Error())

		return nil, err
	}

	var podObjects []*Builder

	for _, runningPod := range podList.Items {
		copiedPod := runningPod
		podBuilder := &Builder{
			apiClient:  apiClient,
			Object:     &copiedPod,
			Definition: &copiedPod,
		}

		podObjects = append(podObjects, podBuilder)
	}

	return podObjects, nil
}

// WaitForAllPodsInNamespaceRunning check that all pods in namespace that match options are in running state.
func WaitForAllPodsInNamespaceRunning(
	apiClient *clients.Settings,
	nsname string,
	timeout time.Duration,
	options ...v1.ListOptions) (bool, error) {
	if nsname == "" {
		glog.V(100).Infof("'nsname' parameter can not be empty")

		return false, fmt.Errorf("failed to list pods, 'nsname' parameter is empty")
	}

	logMessage := fmt.Sprintf("Waiting for all pods in %s namespace", nsname)
	passedOptions := v1.ListOptions{}

	if len(options) == 1 {
		passedOptions = options[0]
		logMessage += fmt.Sprintf(" with the options %v", passedOptions)
	} else if len(options) > 1 {
		glog.V(100).Infof("'options' parameter must be empty or single-valued")

		return false, fmt.Errorf("error: more than one ListOptions was passed")
	}

	glog.V(100).Infof(logMessage + " are in running state")

	podList, err := List(apiClient, nsname, passedOptions)
	if err != nil {
		glog.V(100).Infof("Failed to list all pods due to %s", err.Error())

		return false, err
	}

	for _, podObj := range podList {
		err = podObj.WaitUntilRunning(timeout)
		if err != nil {
			glog.V(100).Infof("Timout was reached while waiting for all pods in running state: %s", err.Error())

			return false, err
		}
	}

	return true, nil
}
