package cmd

import (
	"fmt"
	"os"

	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/NETWAYS/check_cloud_azure/internal/common"
	"github.com/NETWAYS/check_cloud_azure/internal/compute"
	"github.com/NETWAYS/go-check"
	"github.com/spf13/cobra"
)

var (
	Subscription  string
	ComputeClient *compute.Client
)

var computeCmd = &cobra.Command{
	Use:   "compute",
	Short: "Checks in the compute context",
	Run:   Help,
}

func RequireComputeClient() {
	// Lookup environment or auth-file for subscription
	if sub := os.Getenv("AZURE_SUBSCRIPTION_ID"); sub != "" {
		Subscription = sub
	} else if settings, err := auth.GetSettingsFromFile(); err == nil {
		if sub, ok := settings.Values["AZURE_SUBSCRIPTION_ID"]; ok {
			Subscription = sub
		}
	}

	if Subscription == "" {
		check.ExitError(fmt.Errorf("please specify Azure subscription id"))
	}

	token, err := common.CreateCredential(Config)
	if err != nil {
		check.ExitError(fmt.Errorf("Failed to create token"))
	}

	ComputeClient = compute.NewClient(token, Subscription)
}

func init() {
	p := computeCmd.PersistentFlags()
	p.StringVarP(&Subscription, "sub", "s", Subscription, "Azure Subscription ID (env:AZURE_SUBSCRIPTION_ID)")
}
