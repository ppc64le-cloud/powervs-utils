package multiworkspace

import (
	"errors"
	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
)

var (
	// make sure types correctly implement the interfaces.
	_ MultiWorkspace = &multiWorkspace{}
	_ Clients        = &clients{}
	_ ResourceClient = &resourceClient{}
)

const (
	// PowerVSResourceID is Power VS power-iaas service id, can be retrieved using ibmcloud cli
	// ibmcloud catalog service power-iaas.
	powerVSResourceID = "abd259f0-9990-11e8-acc8-b9f54a8f1661"

	// PowerVSResourcePlanID is Power VS power-iaas plan id, can be retrieved using ibmcloud cli
	// ibmcloud catalog service power-iaas.
	powerVSResourcePlanID = "f165dd34-3a40-423b-9d95-e90a23f724dd"
)

var (
	// errWorkspaceNotSet is used to inform users that no PowerVS workspaces is passed while creating multi-workspace client.
	errWorkspaceNotSet = errors.New("PowerVS workspace not set")

	// errAuthenticatorNotSet is used to inform users that authenticator is not passed while creating multi-workspace client.
	errAuthenticatorNotSet = errors.New("IBM Cloud authenticator not set")

	// errZoneNotSet is used to inform users that zone is not set for workspace.
	errZoneNotSet = errors.New("zone not set for workspace")

	// errWorkspaceNameAndIDNotSet is used to inform users that both workspace name and id are not set.
	errWorkspaceNameAndIDNotSet = errors.New("both workspace name and id are not set")

	// errCreatingWorkspaceClients is used to inform users that there is an error while creating clients for provided workspaces.
	errCreatingWorkspaceClients = errors.New("failed to create workspace client")

	// errCreatingResourceClients is used to inform users that there is an error while creating resource client.
	errCreatingResourceClients = errors.New("failed to create resource clients")

	// ErrInstanceNotFound is used to inform users that no instance is found in all configured workspaces.
	ErrInstanceNotFound = errors.New("instance not found")
)

// MultiWorkspace exposes all the functionalities provided by multi-workspace client.
type MultiWorkspace interface {
	GetInstanceDetails(instanceName string) (*InstanceDetails, error)
	GetInstanceDHCPIP(networkName, instanceMac string, workspace Workspace) (string, error)
}

// Clients provides PowerVS workspace services.
type Clients interface {
	GetAll() (*models.PVMInstances, error)
	Get(ID string) (*models.PVMInstance, error)

	GetDHCPServers() (models.DHCPServers, error)
	GetDHCPServerByID(string) (*models.DHCPServerDetail, error)
}

// ResourceClient provides resource client services.
type ResourceClient interface {
	GetWorkspace(name string, zone string) (*resourcecontrollerv2.ResourceInstance, error)
}

// Options are used to configure the multi-workspace client.
type Options struct {
	// List of PowerVS workspaces.
	Workspaces []Workspace

	// IBM Cloud authenticator.
	Authenticator *core.IamAuthenticator
}

// Workspace holds the details about IBM PowerVS workspace.
type Workspace struct {
	// PowerVS workspace name.
	// Either Name or ID is mandatory.
	Name string

	// PowerVS Workspace ID.
	// Either Name or ID is mandatory.
	ID string

	// PowerVS zone.
	// Required.
	Zone string
}

// multiWorkspace holds all the workspace clients.
type multiWorkspace struct {
	workspaceClients []workspaceClient
}

// multiWorkspace implements MultiWorkspace interface.
type workspaceClient struct {
	Workspace
	clients Clients
}

// clients implements Clients interface.
type clients struct {
	instanceClient *instance.IBMPIInstanceClient
	dhcpClient     *instance.IBMPIDhcpClient
}

// resourceClient implements ResourceClient interface.
type resourceClient struct {
	client *resourcecontrollerv2.ResourceControllerV2
}

// InstanceDetails contains the instance details and its corresponding workspace details.
type InstanceDetails struct {
	// Instance details
	Instance *models.PVMInstance

	// The workspace details in which the Instance exists.
	Workspace Workspace
}
