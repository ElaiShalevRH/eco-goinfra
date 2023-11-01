package olm

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListPackageManifest returns PackageManifest inventory in the given namespace.
func ListPackageManifest(
	apiClient *clients.Settings,
	nsname string,
	options ...metaV1.ListOptions) ([]*PackageManifestBuilder, error) {
	passedOptions := metaV1.ListOptions{}

	if len(options) == 1 {
		passedOptions = options[0]
	} else if len(options) > 1 {

		return nil, fmt.Errorf("error: more than one ListOptions was passed")
	}

	glog.V(100).Infof("Listing PackageManifests in the namespace %s with the options %v", nsname, passedOptions)

	if nsname == "" {
		glog.V(100).Infof("packagemanifest 'nsname' parameter can not be empty")

		return nil, fmt.Errorf("failed to list packagemanifests, 'nsname' parameter is empty")
	}

	pkgManifestList, err := apiClient.PackageManifestInterface.PackageManifests(nsname).List(context.Background(),
		passedOptions)

	if err != nil {
		glog.V(100).Infof("Failed to list PackageManifests in the namespace %s due to %s",
			nsname, err.Error())

		return nil, err
	}

	var pkgManifestObjects []*PackageManifestBuilder

	for _, runningPkgManifest := range pkgManifestList.Items {
		copiedPkgManifest := runningPkgManifest
		pkgManifestBuilder := &PackageManifestBuilder{
			apiClient:  apiClient,
			Object:     &copiedPkgManifest,
			Definition: &copiedPkgManifest,
		}

		pkgManifestObjects = append(pkgManifestObjects, pkgManifestBuilder)
	}

	return pkgManifestObjects, nil
}
