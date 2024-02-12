package compute

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

type Client struct {
	Token          azcore.TokenCredential
	SubscriptionID string
	VMClient       *armcompute.VirtualMachinesClient
	GroupClient    *armresources.ResourceGroupsClient
	Context        context.Context
}

func NewClient(token azcore.TokenCredential, subscriptionId string) *Client {
	return &Client{
		Token:          token,
		SubscriptionID: subscriptionId,
		Context:        context.Background(),
	}
}

func (c *Client) GetVMClient(cred azcore.TokenCredential) (*armcompute.VirtualMachinesClient, error) {
	if c.VMClient == nil {
		client, err := armcompute.NewVirtualMachinesClient(c.SubscriptionID, cred, nil)
		if err != nil {
			return client, err
		}

		c.VMClient = client
	}

	return c.VMClient, nil
}

func (c *Client) LoadVmByName(groupName string, vmName string) (*VirtualMachine, error) {
	fmt.Printf("%v\n", c.VMClient)
	local, err := c.VMClient.Get(c.Context, groupName, vmName, nil)

	if err != nil {
		err = fmt.Errorf("could not load vm '%s': %w", vmName, err)
		return nil, err
	}

	result := VirtualMachine{
		VirtualMachine: &local.VirtualMachine,
	}

	return &result, nil
}

func (c *Client) GetGroupClient() *armresources.ResourceGroupsClient {
	if c.GroupClient == nil {
		resourcesClientFactory, err := armresources.NewClientFactory(c.SubscriptionID, c.Token, nil)
		if err != nil {
			return nil
		}

		resourceGroupClient := resourcesClientFactory.NewResourceGroupsClient()
		c.GroupClient = resourceGroupClient
	}

	return c.GroupClient
}

func (c *Client) LoadResourceGroup(name string) (armresources.ResourceGroup, error) {
	groupResponse, err := c.GetGroupClient().Get(c.Context, name, nil)
	if err != nil {
		err = fmt.Errorf("could not load group '%s': %w", name, err)
		return armresources.ResourceGroup{}, err
	}

	return groupResponse.ResourceGroup, nil
}

func (c *Client) LoadResourceGroupsByFilter(name, value string) ([]armresources.ResourceGroup, error) {
	pager := c.GetGroupClient().NewListPager(nil)
	result := []armresources.ResourceGroup{}

	for pager.More() {
		nextPage, err := pager.NextPage(c.Context)
		if err != nil {
			return nil, err
		}

		for _, resGrp := range nextPage.Value {
			val, ok := resGrp.Tags[name]

			if ok && *val == value {
				result = append(result, *resGrp)
			}
		}
	}

	return result, nil
}

func (c *Client) LoadVmsByResourceGroup(group string) (*VirtualMachines, error) {
	vms := &VirtualMachines{}
	pager := c.VMClient.NewListPager(group, nil)

	for pager.More() {
		nextPage, err := pager.NextPage(c.Context)
		if err != nil {
			return nil, err
		}

		for _, vm := range nextPage.VirtualMachineListResult.Value {
			tmp := VirtualMachine{
				VirtualMachine: vm,
			}
			vms.VirtualMachines = append(vms.VirtualMachines, &tmp)
		}
	}

	return vms, nil
}
