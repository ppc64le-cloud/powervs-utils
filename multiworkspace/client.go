package multiworkspace

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/ppc64le-cloud/powervs-utils/utils"
)

// newClients creates and returns a new PowerVS workspace clients.
func newClients(session *ibmpisession.IBMPISession, workspaceID string) Clients {
	ctx := context.Background()
	return &clients{
		instanceClient: instance.NewIBMPIInstanceClient(ctx, session, workspaceID),
		dhcpClient:     instance.NewIBMPIDhcpClient(ctx, session, workspaceID),
	}
}

// GetAll returns all the server instances(VM) in a workspace.
func (c *clients) GetAll() (*models.PVMInstances, error) {
	return c.instanceClient.GetAll()
}

// Get returns details about a given PowerVS server instance(VM) with given ID.
func (c *clients) Get(ID string) (*models.PVMInstance, error) {
	return c.instanceClient.Get(ID)
}

func (c *clients) GetDHCPServers() (models.DHCPServers, error) {
	return c.dhcpClient.GetAll()
}

func (c *clients) GetDHCPServerByID(id string) (*models.DHCPServerDetail, error) {
	return c.dhcpClient.Get(id)
}

// newResourceClient creates and returns a new resource client.
func newResourceClient(auth *core.IamAuthenticator) (ResourceClient, error) {
	client, err := resourcecontrollerv2.NewResourceControllerV2(&resourcecontrollerv2.ResourceControllerV2Options{
		Authenticator: auth,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create resource client: %w", err)
	}
	return &resourceClient{
		client: client,
	}, nil

}

// GetWorkspace return the workspace with given name in given zone.
func (c *resourceClient) GetWorkspace(name string, zone string) (*resourcecontrollerv2.ResourceInstance, error) {
	var serviceInstancesList []resourcecontrollerv2.ResourceInstance

	powerVSResourceID := powerVSResourceID
	powerVSResourcePlanID := powerVSResourcePlanID

	f := func(start string) (bool, string, error) {
		listServiceInstanceOptions := &resourcecontrollerv2.ListResourceInstancesOptions{
			ResourceID:     &powerVSResourceID,
			ResourcePlanID: &powerVSResourcePlanID,
		}
		if name != "" {
			listServiceInstanceOptions.Name = &name
		}
		if start != "" {
			listServiceInstanceOptions.Start = &start
		}

		serviceInstances, _, err := c.client.ListResourceInstances(listServiceInstanceOptions)
		if err != nil {
			return false, "", err
		}
		if serviceInstances != nil {
			if zone != "" {
				for _, resource := range serviceInstances.Resources {
					if *resource.RegionID == zone {
						serviceInstancesList = append(serviceInstancesList, resource)
					}
				}
			}

			nextURL, err := serviceInstances.GetNextStart()
			if err != nil {
				return false, "", err
			}
			if nextURL == nil {
				return true, "", nil
			}
			return false, *nextURL, nil
		}
		return true, "", nil
	}

	if err := utils.PagingHelper(f); err != nil {
		return nil, fmt.Errorf("error listing workspaces: %w", err)
	}
	switch len(serviceInstancesList) {
	case 0:
		return nil, nil
	case 1:
		return &serviceInstancesList[0], nil
	default:
		errStr := fmt.Errorf("there exist more than one workspace with same name %s, Try setting workspace.ID", name)
		return nil, errStr
	}
}

// newIBMPowerVSSession creates new PowerVS session.
func newIBMPowerVSSession(zone string, auth *core.IamAuthenticator) (*ibmpisession.IBMPISession, error) {
	accountID, err := getAccountID(auth)
	if err != nil {
		return nil, fmt.Errorf("failed to get account ID: %w", err)
	}

	return ibmpisession.NewIBMPISession(&ibmpisession.IBMPIOptions{
		Authenticator: auth,
		UserAccount:   accountID,
		Zone:          zone,
	})
}

// getAccountID parses the account ID from the token and returns it.
func getAccountID(auth core.Authenticator) (string, error) {
	// fake request to get a barer token from the request header
	ctx := context.TODO()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", http.NoBody)
	if err != nil {
		return "", err
	}
	err = auth.Authenticate(req)
	if err != nil {
		return "", err
	}
	bearerToken := req.Header.Get("Authorization")
	if strings.HasPrefix(bearerToken, "Bearer") {
		bearerToken = bearerToken[7:]
	}
	token, err := jwt.Parse(bearerToken, func(_ *jwt.Token) (interface{}, error) {
		return "", nil
	})
	if err != nil && !strings.Contains(err.Error(), "key is of invalid type") {
		return "", err
	}

	return token.Claims.(jwt.MapClaims)["account"].(map[string]interface{})["bss"].(string), nil
}
