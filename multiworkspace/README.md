# Multi-workspace PowerVS Client

Multi-workspace PowerVS client helps in performing VM search operation for a given name across multiple workspaces.
It abstracts session and client creation per workspace and provides a simple multi-workspace client.


## Example Usage

```go
package main

import (
	"errors"
	"fmt"
	
	"github.com/IBM/go-sdk-core/v5/core"

	"github.com/ppc64le-cloud/powervs-utils/multiworkspace"
)

// IBMCloudAPIKEY is the IBM Cloud API Key.
const IBMCloudAPIKEY = ""

// powerVSMachineName to search across workspaces.
var powerVSMachineName = "instance-name"

func main() {

	options := multiworkspace.Options{
		Workspaces: []multiworkspace.Workspace{
			{
				Name: "workspace-1",
				Zone: "mad02",
			},
			{
				Name: "workspace-2",
				Zone: "osa21",
			},
		},
		Authenticator: &core.IamAuthenticator{
			ApiKey: IBMCloudAPIKEY,
		},
	}
	client, err := multiworkspace.New(options)
	if err != nil {
		panic(err)
	}
	instanceDetails, err := client.GetInstanceDetails(powerVSMachineName)
	if err != nil {
		if errors.Is(err, multiworkspace.ErrInstanceNotFound) {
			fmt.Println("Machine not found")
		}
		panic(err)
	}
	fmt.Println("Instance ID", *instanceDetails.Instance.PvmInstanceID)
	fmt.Println("Instance Workspace ID", instanceDetails.Workspace.ID)

	// Given that the instance is attached to only one DHCP network.
	ip, err := client.GetInstanceDHCPIP(instanceDetails.Instance.Networks[0].NetworkName, instanceDetails.Instance.Networks[0].MacAddress, instanceDetails.Workspace)
	if err != nil {
		panic(err)
	}
	fmt.Println("DHCP IP of instance", ip)
}
```