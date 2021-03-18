package cmd

import (
	"fmt"
	"github.com/NETWAYS/check_cloud_azure/internal/compute"
	"github.com/NETWAYS/go-check"
	"github.com/spf13/cobra"
)

var (
	VmName   string
	ResGroup string
)

var computeVmCmd = &cobra.Command{
	Use:   "vm",
	Short: "Checks a single Azure VM",
	Run: func(cmd *cobra.Command, args []string) {
		RequireComputeClient()

		var (
			err error
			vm  *compute.VirtualMachine
		)

		if VmName != "" && ResGroup != "" {
			vm, err = ComputeClient.LoadVmByName(ResGroup, VmName)
			if err != nil {
				check.ExitError(fmt.Errorf("could not load vm: %w", err))
			}
		} else {
			check.ExitError(fmt.Errorf("please specify the name and resource group of the VM"))
		}

		output := vm.GetOutput()
		output += "\n\n" + vm.GetLongOutput()

		check.Exit(vm.GetStatus(), output)
	},
}

func init() {
	computeVmCmd.Flags().StringVarP(&VmName, "name", "n", "", "Look for vm by name")
	computeVmCmd.Flags().StringVarP(&ResGroup, "group", "r", "", "Azure resource group")

	computeCmd.AddCommand(computeVmCmd)
}
