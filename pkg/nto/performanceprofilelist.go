package nto //nolint:misspell

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	v2 "github.com/openshift/cluster-node-tuning-operator/pkg/apis/performanceprofile/v2"
	goclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// ListProfiles returns a list of all installed PerformanceProfiles.
func ListProfiles(apiClient *clients.Settings, options ...goclient.ListOptions) ([]*Builder, error) {
	glog.V(100).Infof("Listing PerformanceProfiles on cluster")

	passedOptions := goclient.ListOptions{}

	if len(options) == 1 {
		passedOptions = options[0]
	} else if len(options) > 1 {

		return nil, fmt.Errorf("error: more than one ListOptions was passed")
	}

	var performanceProfiles v2.PerformanceProfileList
	err := apiClient.List(context.TODO(), &performanceProfiles, &passedOptions)

	if err != nil {
		glog.V(100).Infof("Failed to list PerformanceProfiles due to %s", err.Error())

		return nil, err
	}

	var perfProfilesObjects []*Builder

	for _, perfProfile := range performanceProfiles.Items {
		copiedPerfProfile := perfProfile
		perfProfileBuilder := &Builder{
			apiClient:  apiClient,
			Object:     &copiedPerfProfile,
			Definition: &copiedPerfProfile,
		}

		perfProfilesObjects = append(perfProfilesObjects, perfProfileBuilder)
	}

	return perfProfilesObjects, nil
}

// CleanAllPerformanceProfiles removes all PerformanceProfiles installed on a cluster.
func CleanAllPerformanceProfiles(apiClient *clients.Settings, options ...goclient.ListOptions) error {
	glog.V(100).Infof("Cleaning up PerformanceProfiles")

	policies, err := ListProfiles(apiClient, options...)

	if err != nil {
		glog.V(100).Infof("Failed to list PerformanceProfiles")

		return err
	}

	for _, policy := range policies {
		_, err = policy.Delete()

		if err != nil {
			glog.V(100).Infof("Failed to delete PerformanceProfiles: %s", policy.Object.Name)

			return err
		}
	}

	return nil
}
