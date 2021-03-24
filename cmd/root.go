package cmd

import (
	"fmt"
	"github.com/Azure/go-autorest/autorest"
	"github.com/NETWAYS/check_cloud_azure/internal/common"
	"github.com/NETWAYS/go-check"
	"github.com/spf13/cobra"
	"os"
)

var (
	Timeout    = 30
	Authorizer autorest.Authorizer
	AuthFile   string
)

var rootCmd = &cobra.Command{
	Use:   "check_cloud_azure",
	Short: "Icinga check plugin to check Microsoft's Azure resources",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		go check.HandleTimeout(Timeout)

		if AuthFile != "" {
			// Inject AZURE_AUTH_LOCATION into the environment
			// There seem to be no other way to pass this to the SDK
			_ = os.Setenv("AZURE_AUTH_LOCATION", AuthFile)
		}
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

func RequireAuth() {
	var err error

	Authorizer, err = common.BuildAnyAuthorizer()
	if err != nil {
		check.ExitError(err)
	}
}

func init() {
	rootCmd.AddCommand(computeCmd)
	rootCmd.SetHelpFunc(Help)

	p := rootCmd.PersistentFlags()
	p.StringVar(&AuthFile, "auth-file", AuthFile, "Azure auth file (env:AZURE_AUTH_LOCATION)")
	p.IntVarP(&Timeout, "timeout", "t", Timeout, "Timeout for the check")
}
