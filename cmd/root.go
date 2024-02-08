package cmd

import (
	"fmt"
	"os"

	"github.com/NETWAYS/check_cloud_azure/internal/common"
	"github.com/NETWAYS/go-check"
	"github.com/spf13/cobra"
)

var (
	Timeout = 30
	Config  common.Config
)

var rootCmd = &cobra.Command{
	Use:   "check_cloud_azure",
	Short: "Icinga check plugin to check Microsoft's Azure resources",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		go check.HandleTimeout(Timeout)
	},
	Run: Help,
}

func Execute(version string) {
	defer check.CatchPanic()

	rootCmd.Version = version

	if err := rootCmd.Execute(); err != nil {
		check.ExitError(err)
	}
}

func Help(cmd *cobra.Command, _ []string) {
	fmt.Println(cmd.Short)
	fmt.Println()

	_ = cmd.Usage()

	os.Exit(3)
}

func init() {
	rootCmd.AddCommand(computeCmd)
	rootCmd.SetHelpFunc(Help)

	p := rootCmd.PersistentFlags()
	p.IntVarP(&Timeout, "timeout", "t", Timeout, "Timeout for the check")

	p.StringVar(&Config.TenantID, "tenant-id", "", "Azure tenant ID")
	p.StringVar(&Config.ClientID, "client-id", "", "Azure client ID")
	p.StringVar(&Config.ClientSecret, "client-secret", "", "Azure client secret")
}
