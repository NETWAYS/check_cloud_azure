package common

import (
	"fmt"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

// BuildAnyAuthorizer tries to build an Authorizer
//
// * Supporting Azure CLI for login
// * Trying environment variables
func BuildAnyAuthorizer() (authorizer autorest.Authorizer, err error) {
	authorizer, err = auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
	if err == nil {
		return
	}

	authorizer, err = auth.NewAuthorizerFromCLI()
	if err == nil {
		return
	}

	authorizer, err = auth.NewAuthorizerFromEnvironment()
	if err == nil {
		return
	}

	err = fmt.Errorf("could not build an Azure authorizer: %w", err)
	return
}
