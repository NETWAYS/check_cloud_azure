package cmd

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
	"github.com/spf13/cobra"
)

var (
	TagName  string
	TagValue string
)

var computeVmsCmd = &cobra.Command{
	Use:   "vms",
	Short: "Checks multiple Azure VMs in a resource group",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := GenerateComputeClient(Config)
		if err != nil {
			check.ExitError(err)
		}

		var (
			groups []armresources.ResourceGroup
		)

		// Load the ressource groups first
		if ResGroup != "" {
			var group armresources.ResourceGroup

			group, err = client.LoadResourceGroup(ResGroup)
			if err != nil {
				check.ExitError(err)
			}

			groups = []armresources.ResourceGroup{group}
		} else {
			groups, err = client.LoadResourceGroupsByFilter(TagName, TagValue)
			if err != nil {
				check.ExitError(fmt.Errorf("could not load resource groups: %w", err))
			} else if len(groups) == 0 {
				check.ExitError(fmt.Errorf("no groups found for %s=%s", TagName, TagValue))
			}
		}

		overall := result.Overall{}

		// Iterate through groups
		// Load all VMs for each group
		// Build one result state from that
		for _, group := range groups {
			vms, err := client.LoadVmsByResourceGroup(*group.Name)
			if err != nil {
				check.ExitError(err)
			}

			groupPartial := vms.GetPartialResult()
			groupPartial.Output += "Group: " + *group.Name

			overall.AddSubcheck(groupPartial)
		}
		check.ExitRaw(overall.GetStatus(), overall.GetOutput())
	},
}

func init() {
	computeVmsCmd.Flags().StringVarP(&ResGroup, "group", "r", "", "Azure resource group")
	computeVmsCmd.Flags().StringVarP(&TagName, "tagname", "n", "", "Filter resource group by tag (e.g. tag1)")
	computeVmsCmd.Flags().StringVarP(&TagValue, "tagvalue", "v", "", "Tag value of resource group (e.g. value1)")

	computeCmd.AddCommand(computeVmsCmd)
}
