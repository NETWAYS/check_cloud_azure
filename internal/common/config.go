package common

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

type Config struct {
	TenantID       string
	ClientID       string
	ClientSecret   string
	SubscriptionID string
}

func CreateCredential(conf Config) (azcore.TokenCredential, error) {
	clientSecretCredential, err := azidentity.NewClientSecretCredential(conf.TenantID, conf.ClientID, conf.ClientSecret, nil)
	if err != nil {
		return nil, err
	}

	return clientSecretCredential, nil
}
