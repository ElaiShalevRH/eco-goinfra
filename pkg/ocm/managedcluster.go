package ocm

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/msg"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ocmV1Client "open-cluster-management.io/api/client/cluster/clientset/versioned/typed/cluster/v1"
	mcV1 "open-cluster-management.io/api/cluster/v1"
)

// ManagedClusterBuilder provides struct for the ManagedCluster object containing connection to
// the cluster and the ManagedCluster definitions.
type ManagedClusterBuilder struct {
	Definition *mcV1.ManagedCluster
	Object     *mcV1.ManagedCluster
	errorMsg   string
	apiClient  ocmV1Client.ClusterV1Interface
}

// ManagedClusterAdditionalOptions additional options for ManagedCluster object.
type ManagedClusterAdditionalOptions func(builder *ManagedClusterBuilder) (*ManagedClusterBuilder, error)

// NewManagedClusterBuilder creates a new instance of ManagedClusterBuilder.
func NewManagedClusterBuilder(apiClient *clients.Settings, name string) *ManagedClusterBuilder {
	glog.V(100).Infof(
		`Initializing new ManagedCluster structure with the following params: name: %s`, name)

	if apiClient == nil {
		glog.V(100).Infof("The apiClient is empty")
		return nil
	}

	builder := ManagedClusterBuilder{
		apiClient: &apiClient.ClusterV1Client,
		Definition: &mcV1.ManagedCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
			Spec: mcV1.ManagedClusterSpec{},
		},
	}

	if name == "" {
		glog.V(100).Infof("The name of the ManagedCluster is empty")

		builder.errorMsg = "managedCluster 'name' cannot be empty"
	}

	return &builder
}

// WithOptions creates ManagedCluster with generic mutation options.
func (builder *ManagedClusterBuilder) WithOptions(
	options ...ManagedClusterAdditionalOptions) *ManagedClusterBuilder {
	if valid, _ := builder.validate(); !valid {
		return builder
	}

	glog.V(100).Infof("Setting ManagedCluster additional options")

	for _, option := range options {
		if option != nil {
			builder, err := option(builder)

			if err != nil {
				glog.V(100).Infof("Error occurred in mutation function")

				builder.errorMsg = err.Error()

				return builder
			}
		}
	}

	return builder
}

// PullManagedCluster loads an existing ManagedCluster into ManagedClusterBuilder struct.
func PullManagedCluster(apiClient *clients.Settings, name string) (*ManagedClusterBuilder, error) {
	glog.V(100).Infof("Pulling existing ManagedCluster name: %s", name)

	if apiClient == nil {
		glog.V(100).Infof("The apiClient is empty")
		return nil, fmt.Errorf("the apiClient is empty")
	}

	builder := ManagedClusterBuilder{
		apiClient: &apiClient.ClusterV1Client,
		Definition: &mcV1.ManagedCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		},
	}

	if name == "" {
		builder.errorMsg = "managedcluster 'name' cannot be empty"
	}

	if !builder.Exists() {
		return nil, fmt.Errorf("managedcluster object %s doesn't exist", name)
	}

	builder.Definition = builder.Object

	return &builder, nil
}

// Update modifies an existing ManagedCluster on the cluster.
func (builder *ManagedClusterBuilder) Update() (*ManagedClusterBuilder, error) {
	if valid, err := builder.validate(); !valid {
		return builder, err
	}

	glog.V(100).Infof("Updating ManagedCluster %s", builder.Definition.Name)

	var err error
	builder.Object, err = builder.apiClient.ManagedClusters().Update(context.TODO(), builder.Definition, metav1.UpdateOptions{})

	return builder, err
}

// Delete removes a ManagedCluster from the cluster.
func (builder *ManagedClusterBuilder) Delete() error {
	if valid, err := builder.validate(); !valid {
		return err
	}

	glog.V(100).Infof("Deleting the ManagedCluster %s", builder.Definition.Name)

	if !builder.Exists() {
		return fmt.Errorf("managedcluster cannot be deleted because it does not exist")
	}

	err := builder.apiClient.ManagedClusters().Delete(context.TODO(), builder.Definition.Name, metav1.DeleteOptions{})

	if err != nil {
		return fmt.Errorf("cannot delete managedcluster: %w", err)
	}

	builder.Object = nil
	builder.Definition.ResourceVersion = ""
	builder.Definition.CreationTimestamp = metav1.Time{}

	return nil
}

// Exists checks if the defined ManagedCluster has already been created.
func (builder *ManagedClusterBuilder) Exists() bool {
	if valid, _ := builder.validate(); !valid {
		return false
	}

	glog.V(100).Infof("Checking if managedcluster %s exists", builder.Definition.Name)

	var err error
	builder.Object, err = builder.apiClient.ManagedClusters().Get(context.TODO(), builder.Definition.Name, metav1.GetOptions{})

	return err == nil || !k8serrors.IsNotFound(err)
}

// validate will check that the builder and builder definition are properly initialized before
// accessing any member fields.
func (builder *ManagedClusterBuilder) validate() (bool, error) {
	resourceCRD := "ManagedCluster"

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