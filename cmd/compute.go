package cmd

import (
	"errors"
	"fmt"

	"github.com/NETWAYS/check_cloud_azure/internal/common"
	"github.com/NETWAYS/check_cloud_azure/internal/compute"
	"github.com/NETWAYS/go-check"
	"github.com/spf13/cobra"
)

var computeCmd = &cobra.Command{
	Use:   "compute",
	Short: "Checks in the compute context",
	Run:   Help,
}

func GenerateComputeClient(config common.Config) (*compute.Client, error) {
	if Config.SubscriptionID == "" {
		retError := errors.New("Subscription is missing")
		return nil, retError
	}

	token, err := common.CreateCredential(Config)
	if err != nil {
		check.ExitError(fmt.Errorf("Failed to create token"))
	}

	client := compute.NewClient(token, Config.SubscriptionID)

	client.VMClient, err = client.GetVMClient(token)
	if err != nil {
		check.ExitError(fmt.Errorf("Failed to create VMClient"))
	}

	return client, nil
}

func init() {
	p := computeCmd.PersistentFlags()
	p.StringVarP(&Config.SubscriptionID, "sub", "s", "", "Azure Subscription ID")
}
