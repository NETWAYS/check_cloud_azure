package compute

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-10-01/resources"
	"github.com/Azure/go-autorest/autorest"
)

type Client struct {
	Authorizer     autorest.Authorizer
	SubscriptionID string
	VMClient       *compute.VirtualMachinesClient
	GroupClient    *resources.GroupsClient
	Context        context.Context
}

func NewClient(authorizer autorest.Authorizer, subscriptionId string) *Client {
	return &Client{
		Authorizer:     authorizer,
		SubscriptionID: subscriptionId,
		Context:        context.Background(),
	}
}

func (c *Client) GetVMClient() *compute.VirtualMachinesClient {
	if c.VMClient == nil {
		client := compute.NewVirtualMachinesClient(c.SubscriptionID)
		client.Authorizer = c.Authorizer
		c.VMClient = &client
	}

	return c.VMClient
}

func (c *Client) LoadVmByName(groupName string, vmName string) (vm *VirtualMachine, err error) {
	local, err := c.GetVMClient().Get(c.Context, groupName, vmName, compute.InstanceView)
	if err != nil {
		err = fmt.Errorf("could not load vm '%s': %w", vmName, err)
		return
	}

	vm = &VirtualMachine{VirtualMachine: &local}

	return
}

func (c *Client) GetGroupClient() *resources.GroupsClient {
	if c.GroupClient == nil {
		client := resources.NewGroupsClient(c.SubscriptionID)
		client.Authorizer = c.Authorizer
		c.GroupClient = &client
	}

	return c.GroupClient
}

func (c *Client) LoadResourceGroup(name string) (group resources.Group, err error) {
	group, err = c.GetGroupClient().Get(c.Context, name)
	if err != nil {
		err = fmt.Errorf("could not load group '%s': %w", name, err)
	}

	return
}

func (c *Client) LoadResourceGroupsByFilter(name, value string) (groups []resources.Group, err error) {
	iter, err := c.GetGroupClient().ListComplete(c.Context, createFilter(name, value), nil)
	if err != nil {
		err = fmt.Errorf("could not load resource group: %w", err)
		return
	}

	for ; iter.NotDone(); err = iter.NextWithContext(c.Context) {
		if err != nil {
			err = fmt.Errorf("could not load resource group: %w", err)
			return
		}

		if iter.Value().ID != nil {
			groups = append(groups, iter.Value())
		}
	}

	return
}

func (c *Client) LoadVmsByResourceGroup(group string) (vms *VirtualMachines, err error) {
	vms = &VirtualMachines{}

	iter, err := c.GetVMClient().ListComplete(c.Context, group)
	if err != nil {
		err = fmt.Errorf("could not load virtual machines by group: %w", err)
		return
	}

	for ; iter.NotDone(); err = iter.NextWithContext(c.Context) {
		if err != nil {
			err = fmt.Errorf("could not load virtual machines by group: %w", err)
			return
		}

		if iter.Value().ID != nil {
			name := *iter.Value().Name
			var vmFull *VirtualMachine
			vmFull, err = c.LoadVmByName(group, name)
			if err != nil {
				return
			}

			vms.VirtualMachines = append(vms.VirtualMachines, vmFull)
		}
	}

	return
}

func createFilter(name, value string) string {
	// $filter=tagName eq 'tag1' and tagValue eq 'Value1'
	if name == "" && value == "" {
		return ""
	} else {
		return fmt.Sprintf("tagName eq '%s' and tagValue eq '%s'", name, value)
	}
}
