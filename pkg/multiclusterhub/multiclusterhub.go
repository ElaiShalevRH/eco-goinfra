package multiclusterhub

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/msg"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mchv1 "github.com/stolostron/multiclusterhub-operator/api/v1"
	goclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// MultiClusterHubBuilder provides struct for the MultiClusterHub object containing connection to
// the cluster and the MultiClusterHub definitions.
type MultiClusterHubBuilder struct {
	Definition *mchv1.MultiClusterHub
	Object     *mchv1.MultiClusterHub
	errorMsg   string
	apiClient  goclient.Client
}

// NewMultiClusterHubBuilder creates a new instance of MultiClusterHubBuilder.
func NewMultiClusterHubBuilder(apiClient *clients.Settings, name, namespace string) *MultiClusterHubBuilder {
	glog.V(100).Infof(
		`Initializing new MultiClusterHub structure with the following params: name: %s, namespace:`,
		name, namespace)

	builder := MultiClusterHubBuilder{
		apiClient: apiClient.Client,
		Definition: &mchv1.MultiClusterHub{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		},
	}

	if name == "" {
		glog.V(100).Infof("The name of the MultiClusterHub is empty")

		builder.errorMsg = "MultiClusterHub 'name' cannot be empty"
	}

	if namespace == "" {
		glog.V(100).Infof("The namespace of the MultiClusterHub is empty")

		builder.errorMsg = "MultiClusterHub 'namespace' cannot be empty"
	}

	return &builder
}

// PullMultiClusterHub loads an existing MultiClusterHub into MultiClusterHubBuilder struct.
func PullMultiClusterHub(apiClient *clients.Settings, name, namespace string) (*MultiClusterHubBuilder, error) {
	glog.V(100).Infof("Pulling existing MultiClusterHub name: %s", name)

	builder := MultiClusterHubBuilder{
		apiClient: apiClient.Client,
		Definition: &mchv1.MultiClusterHub{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		},
	}

	if name == "" {
		builder.errorMsg = "MultiClusterHub 'name' cannot be empty"
	}

	if namespace == "" {
		builder.errorMsg = "MultiClusterHub 'namespace' cannot be empty"
	}

	if !builder.Exists() {
		return nil, fmt.Errorf("MultiClusterHub object %s does not exist", name)
	}

	builder.Definition = builder.Object

	return &builder, nil
}

// Get fetches the defined MultiClusterHub from the cluster.
func (builder *MultiClusterHubBuilder) Get() (*mchv1.MultiClusterHub, error) {
	if valid, err := builder.validate(); !valid {
		return nil, err
	}

	glog.V(100).Infof("Getting MultiClusterHub %s", builder.Definition.Name)

	multiClusterHub := &mchv1.MultiClusterHub{}
	err := builder.apiClient.Get(context.TODO(), goclient.ObjectKey{
		Name:      builder.Definition.Name,
		Namespace: builder.Definition.Namespace,
	}, multiClusterHub)

	if err != nil {
		return nil, err
	}

	return multiClusterHub, err
}

// Update modifies an existing MultiClusterHub on the cluster.
func (builder *MultiClusterHubBuilder) Update() (*MultiClusterHubBuilder, error) {
	if valid, err := builder.validate(); !valid {
		return builder, err
	}

	glog.V(100).Infof("Updating MultiClusterHub %s in the namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	err := builder.apiClient.Update(context.TODO(), builder.Definition)
	builder.Object = builder.Definition

	return builder, err
}

// Delete removes a MultiClusterHub from the cluster.
func (builder *MultiClusterHubBuilder) Delete() error {
	if valid, err := builder.validate(); !valid {
		return err
	}

	glog.V(100).Infof("Deleting the MultiClusterHub %s in the namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	if !builder.Exists() {
		return fmt.Errorf("MultiClusterHub cannot be deleted because it does not exist")
	}

	err := builder.apiClient.Delete(context.TODO(), builder.Definition)

	if err != nil {
		return fmt.Errorf("cannot delete MultiClusterHub: %w", err)
	}

	builder.Object = nil
	builder.Definition.ResourceVersion = ""
	builder.Definition.CreationTimestamp = metav1.Time{}

	return nil
}

// Exists checks if the defined MultiClusterHub has already been created.
func (builder *MultiClusterHubBuilder) Exists() bool {
	if valid, _ := builder.validate(); !valid {
		return false
	}

	glog.V(100).Infof("Checking if MultiClusterHub %s in namespace %s exists",
		builder.Definition.Name, builder.Definition.Namespace)

	var err error
	builder.Object, err = builder.Get()

	return err == nil || !k8serrors.IsNotFound(err)
}

// validate will check that the builder and builder definition are properly initialized before
// accessing any member fields.
func (builder *MultiClusterHubBuilder) validate() (bool, error) {
	resourceCRD := "MultiClusterHub"

	if builder == nil {
		glog.V(100).Infof("The %s builder is uninitialized", resourceCRD)

		return false, fmt.Errorf("error: received nil %s builder", resourceCRD)
	}

	if builder.Definition == nil {
		glog.V(100).Infof("The %s is undefined", resourceCRD)

		builder.errorMsg = msg.UndefinedCrdObjectErrString(resourceCRD)
	}

	if builder.apiClient == nil {
		glog.V(100).Infof("The %s builder apiclient is nil", resourceCRD)

		builder.errorMsg = fmt.Sprintf("%s builder cannot have nil apiClient", resourceCRD)
	}

	if builder.errorMsg != "" {
		glog.V(100).Infof("The %s builder has error message: %s", resourceCRD, builder.errorMsg)

		return false, fmt.Errorf(builder.errorMsg)
	}

	return true, nil
}
