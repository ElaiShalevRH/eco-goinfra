package olm

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListInstallPlan returns a list of installplans found for specific namespace.
func ListInstallPlan(
	apiClient *clients.Settings, nsname string, options ...v1.ListOptions) ([]*InstallPlanBuilder, error) {
	if nsname == "" {
		glog.V(100).Info("The nsname of the installplan is empty")

		return nil, fmt.Errorf("the nsname of the installplan is empty")
	}

	passedOptions := v1.ListOptions{}

	if len(options) == 1 {
		passedOptions = options[0]
	} else if len(options) > 1 {

		return nil, fmt.Errorf("error: more than one ListOptions was passed")
	}

	installPlanList, err := apiClient.InstallPlans(nsname).List(context.Background(), passedOptions)

	if err != nil {
		glog.V(100).Infof("Failed to list all installplan in namespace %s due to %s",
			nsname, err.Error())

		return nil, err
	}

	var installPlanObjects []*InstallPlanBuilder

	for _, foundCsv := range installPlanList.Items {
		copiedCsv := foundCsv
		csvBuilder := &InstallPlanBuilder{
			apiClient:  apiClient,
			Object:     &copiedCsv,
			Definition: &copiedCsv,
		}

		installPlanObjects = append(installPlanObjects, csvBuilder)
	}

	if len(installPlanObjects) == 0 {
		return nil, fmt.Errorf("installplan not found in namespace %s", nsname)
	}

	return installPlanObjects, nil
}
