package deployment

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// List returns deployment inventory in the given namespace.
func List(apiClient *clients.Settings, nsname string, options ...metaV1.ListOptions) ([]*Builder, error) {
	if nsname == "" {
		glog.V(100).Infof("deployment 'nsname' parameter can not be empty")

		return nil, fmt.Errorf("failed to list deployments, 'nsname' parameter is empty")
	}

	passedOptions := metaV1.ListOptions{}
	logMessage := fmt.Sprintf("Listing deployments in the namespace %s", nsname)

	if len(options) == 1 {
		passedOptions = options[0]
		logMessage += fmt.Sprintf(" with the options %v", passedOptions)
	} else if len(options) > 1 {
		glog.V(100).Infof("'options' parameter must be empty or single-valued")

		return nil, fmt.Errorf("error: more than one ListOptions was passed")
	}

	glog.V(100).Infof(logMessage)

	deploymentList, err := apiClient.Deployments(nsname).List(context.Background(), passedOptions)

	if err != nil {
		glog.V(100).Infof("Failed to list deployments in the namespace %s due to %s", nsname, err.Error())

		return nil, err
	}

	var deploymentObjects []*Builder

	for _, runningDeployment := range deploymentList.Items {
		copiedDeployment := runningDeployment
		deploymentBuilder := &Builder{
			apiClient:  apiClient,
			Object:     &copiedDeployment,
			Definition: &copiedDeployment,
		}

		deploymentObjects = append(deploymentObjects, deploymentBuilder)
	}

	return deploymentObjects, nil
}
