package multiworkspace

import (
	"errors"
	"fmt"
	"log"

	"github.com/IBM-Cloud/power-go-client/power/models"
)

// New creates a new multi-workspace client with given options.
func New(options Options) (MultiWorkspace, error) {
	if options.Authenticator == nil {
		return nil, errAuthenticatorNotSet
	}

	if len(options.Workspaces) == 0 {
		return nil, errWorkspaceNotSet
	}

	if err := options.validateWorkspaces(); err != nil {
		return nil, err
	}

	clients, err := options.constructWorkspaceClients()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errCreatingWorkspaceClients, err)
	}

	return multiWorkspace{
		workspaceClients: clients,
	}, nil
}

func (o Options) validateWorkspaces() error {
	for _, workspace := range o.Workspaces {
		return validateWorkspace(workspace)
	}
	return nil
}

func (o Options) constructWorkspaceClients() ([]workspaceClient, error) {
	resourceClient, err := newResourceClient(o.Authenticator)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errCreatingResourceClients, err)
	}

	var workspaceClients []workspaceClient

	for _, workspace := range o.Workspaces {
		if workspace.Name != "" {
			// If the workspace name is set, fetch the workspace id.
			workspaceDetails, err := resourceClient.GetWorkspace(workspace.Name, workspace.Zone)
			if err != nil {
				return nil, fmt.Errorf("failed to get workspace ID for workspace: %s: %w", workspace.Name, err)
			}
			if workspaceDetails == nil || workspaceDetails.GUID == nil {
				return nil, fmt.Errorf("failed to get workspace ID for workspace: %s", workspace.Name)
			}
			workspace.ID = *workspaceDetails.GUID
		}

		session, err := newIBMPowerVSSession(workspace.Zone, o.Authenticator)
		if err != nil {
			return nil, fmt.Errorf("failed to create PowerVS session: %w", err)
		}

		client := workspaceClient{
			Workspace: workspace,
			clients:   newClients(session, workspace.ID),
		}
		workspaceClients = append(workspaceClients, client)
	}

	return workspaceClients, nil
}

// GetInstanceDetails tries to find PowerVS instance(server instance / VM) with a given `instanceName` in all configured workspaces.
// In case of error during a search operation, it will return the error.
// If the instance with a given `instanceName` not found, it will return an ErrInstanceNotFound error.
// If the instance with a giver `instanceName` found, it will return InstanceDetails response.
func (c multiWorkspace) GetInstanceDetails(instanceName string) (*InstanceDetails, error) {
	log.Println("Searching instance in all workspaces, instanceName:", instanceName)

	// Note: Possible future enhancements
	// 1. Cache the instanceName with workspace so, Instead of iterating over all the workspaces, We can directly look for the instanceName in a particular workspace.
	// 2. Concurrently, look for instanceName in all the workspaces.

	var allErrors []error
	for _, client := range c.workspaceClients {
		instanceDetails, err := findInstance(client, instanceName)
		if err == nil {
			// If the instance is found in workspace return the results.
			return &InstanceDetails{
				Instance:  instanceDetails,
				Workspace: client.Workspace,
			}, nil
		}
		if errors.Is(err, ErrInstanceNotFound) {
			// If the instance not found in workspace, continue with the next workspace.
			continue
		}
		// Any error apart from ErrInstanceNotFound, store them to return at a later stage.
		allErrors = append(allErrors, err)
	}
	if len(allErrors) > 0 {
		return nil, errors.Join(allErrors...)
	}
	return nil, ErrInstanceNotFound
}

// GetInstanceDHCPIP fetches and returns DHCP IP of an instance with `instanceMac` of a DHCP network identified by
// `networkName` in a given `workspace`.
func (c multiWorkspace) GetInstanceDHCPIP(networkName, instanceMac string, workspace Workspace) (string, error) {
	log.Printf("Fetching DHCP IP for instance MAC: %s, network: %s, worksapce: %+v", instanceMac, networkName, workspace)
	if err := validateWorkspace(workspace); err != nil {
		return "", fmt.Errorf("workspace validation error: %w", err)
	}

	var workSpaceClient *workspaceClient
	for _, client := range c.workspaceClients {
		if workspace.ID != "" && workspace.ID == client.ID {
			workSpaceClient = &client
			break
		}
		if workspace.Name != "" && workspace.Name == client.Name {
			workSpaceClient = &client
			break
		}
	}
	if workSpaceClient == nil {
		return "", fmt.Errorf("failed to find workspace client for workspace: %s", workspace.ID)
	}

	// Fetch the DHCP server ID.
	dhcpServerID, err := workSpaceClient.getDHCPServerID(networkName)
	if err != nil {
		return "", fmt.Errorf("failed to get DHCP server id for network: %s: %w", networkName, err)
	}

	// Fetch DHCP server details.
	dhcpServerDetails, err := workSpaceClient.clients.GetDHCPServerByID(dhcpServerID)
	if err != nil {
		return "", fmt.Errorf("failed to get DHCP server details for network: %s: %w", networkName, err)
	}

	var instanceIP string
	for _, lease := range dhcpServerDetails.Leases {
		if instanceMac == *lease.InstanceMacAddress {
			instanceIP = *lease.InstanceIP
		}
	}

	if instanceIP == "" {
		return "", fmt.Errorf("failed to find instance IP for instance MAC: %s", instanceMac)
	}

	return instanceIP, nil
}

// getDHCPServerID fetches and returns the DHCP server ID associated with the `networkName`.
func (client *workspaceClient) getDHCPServerID(networkName string) (string, error) {
	// Fetch all the DHCP servers.
	dhcpServers, err := client.clients.GetDHCPServers()
	if err != nil {
		return "", fmt.Errorf("failed to get DHCP servers: %w", err)
	}

	// Get the DHCP server ID associated with the `networkName`.
	for _, server := range dhcpServers {
		if server.Network != nil && server.Network.Name != nil && *server.Network.Name == networkName {
			return *server.ID, nil
		}
	}
	return "", fmt.Errorf("not able to get DHCP server ID for network %s", networkName)
}

func findInstance(workspaceClient workspaceClient, instanceName string) (*models.PVMInstance, error) {
	log.Println("Searching instance in workspace:", workspaceClient.Name)
	instances, err := workspaceClient.clients.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get the instance list: %w", err)
	}
	for _, pvmInstance := range instances.PvmInstances {
		if pvmInstance != nil && *pvmInstance.ServerName == instanceName {
			log.Println("Found instance in workspace:", workspaceClient.Name)
			return workspaceClient.clients.Get(*pvmInstance.PvmInstanceID)
		}
	}
	return nil, ErrInstanceNotFound
}

func validateWorkspace(workspace Workspace) error {
	if workspace.Zone == "" {
		return fmt.Errorf("%s: %s", errZoneNotSet, workspace.Name)
	}
	if workspace.Name == "" && workspace.ID == "" {
		return errWorkspaceNameAndIDNotSet
	}
	return nil
}
