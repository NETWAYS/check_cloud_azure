package cmd

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-10-01/resources"
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
		RequireComputeClient()

		var (
			err           error
			overallStatus int
			total         int
			counters      = map[string]int{}
			output        string
			groups        []resources.Group
		)

		if ResGroup != "" {
			var group resources.Group

			group, err = ComputeClient.LoadResourceGroup(ResGroup)
			if err != nil {
				check.ExitError(err)
			}

			groups = []resources.Group{group}
		} else {
			groups, err = ComputeClient.LoadResourceGroupsByFilter(TagName, TagValue)
			if err != nil {
				check.ExitError(fmt.Errorf("could not load resource groups: %w", err))
			} else if len(groups) == 0 {
				check.ExitError(fmt.Errorf("no groups found for %s=%s", TagName, TagValue))
			}
		}

		for _, group := range groups {
			vms, err := ComputeClient.LoadVmsByResourceGroup(*group.Name)
			if err != nil {
				check.ExitError(err)
			}

			if vms.IsEmpty() {
				// TODO: should this be a failure?
				continue
			}

			output += "\n## Group: " + *group.Name + "\n\n"

			overallStatus = result.WorstState(overallStatus, vms.GetStatus())
			output += vms.GetOutput()

			for _, vm := range vms.VirtualMachines {
				power, _ := vm.GetPowerState()
				counters[power] += 1
			}

			total++
		}

		summary := fmt.Sprintf("%d VMs found", total)

		if total == 0 {
			overallStatus = check.Unknown
		}

		for state, count := range counters {
			summary += fmt.Sprintf(" - %d %s", count, state)
		}

		check.ExitRaw(overallStatus, summary+"\n"+output)
	},
}

func init() {
	computeVmsCmd.Flags().StringVarP(&ResGroup, "group", "r", "", "Azure resource group")
	computeVmsCmd.Flags().StringVarP(&TagName, "tagname", "n", "", "Filter resource group by tag (e.g. tag1)")
	computeVmsCmd.Flags().StringVarP(&TagValue, "tagvalue", "v", "", "Tag value of resource group (e.g. value1)")

	computeCmd.AddCommand(computeVmsCmd)
}
