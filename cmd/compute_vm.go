package cmd

import (
	"fmt"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
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

		client, err := GenerateComputeClient(Config)
		if err != nil {
			check.ExitError(err)
		}

		if VmName == "" || ResGroup == "" {
			check.ExitError(fmt.Errorf("please specify the name and resource group of the VM"))
		}

		vm, err := client.LoadVmByName(ResGroup, VmName)
		if err != nil {
			check.ExitError(fmt.Errorf("could not load vm: %w", err))
		}

		overall := result.Overall{}
		sc := vm.GetPartialResult()

		overall.AddSubcheck(sc)

		check.ExitRaw(overall.GetStatus(), overall.GetOutput())
	},
}

func init() {
	computeVmCmd.Flags().StringVarP(&VmName, "name", "n", "", "Look for vm by name")
	computeVmCmd.Flags().StringVarP(&ResGroup, "group", "r", "", "Azure resource group")

	computeCmd.AddCommand(computeVmCmd)
}
