package compute_test

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/NETWAYS/check_cloud_azure/internal/common"
	"github.com/NETWAYS/check_cloud_azure/internal/compute"
	"github.com/NETWAYS/go-check/http/mock"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

var (
	AzureAuthorizer     autorest.Authorizer
	AzureSubscriptionId = "SUB-UUID"
)

// allowMocking reconfigures an autorest.Client to not retry and use a basic autorest.Sender implementation
//
// This is required to override any internal transports, so the normal http.DefaultClient is used, and httpmock
// can engage.
func allowMocking(c *autorest.Client) {
	// Disable retry logic to avoid blocking on errors
	c.RetryAttempts = 1
	c.RetryDuration = 1

	c.Sender = http.DefaultClient
}

func handleCredentials() {
	if os.Getenv("AZURE_TENANT_ID") != "" {
		auth, err := common.BuildAnyAuthorizer()
		if err != nil {
			panic(err)
		}
		AzureAuthorizer = auth
	} else {
		AzureAuthorizer = autorest.NullAuthorizer{}
	}

	if subEnv := os.Getenv("AZURE_SUBSCRIPTION_ID"); subEnv != "" {
		AzureSubscriptionId = subEnv
	}
}

func testClientWithMock() (client *compute.Client, cleanup func()) {
	handleCredentials()

	client = compute.NewClient(AzureAuthorizer, AzureSubscriptionId)

	// Reconfigure clients to use a basic http.Client implementation
	allowMocking(&client.GetVMClient().BaseClient.Client)
	allowMocking(&client.GetGroupClient().BaseClient.Client)

	httpmock.Activate()
	cleanup = httpmock.DeactivateAndReset

	checkhttpmock.ActivateRecorder()

	return
}

func withBaseURL(url string) string {
	return "https://management.azure.com/subscriptions/" + AzureSubscriptionId + url
}

func newJsonFileResponder(fileName string) func(request *http.Request) (*http.Response, error) {
	return func(request *http.Request) (*http.Response, error) {
		data, err := ioutil.ReadFile(fileName)
		if err != nil {
			return nil, err
		}

		return httpmock.NewBytesResponse(200, data), nil
	}
}

func TestClient_LoadVmByName(t *testing.T) {
	c, cleanup := testClientWithMock()
	defer cleanup()

	httpmock.RegisterResponder("GET",
		withBaseURL("/resourceGroups/test-group/providers/Microsoft.Compute/virtualMachines/test-vm?%24expand=instanceView&api-version=2020-06-01"),
		newJsonFileResponder("./testdata/vmByName.json"))

	vm, err := c.LoadVmByName("test-group", "test-vm")
	assert.NoError(t, err)
	assert.Equal(t, "test-vm", *vm.VirtualMachine.Name)
}

func TestClient_LoadResourceGroup(t *testing.T) {
	c, cleanup := testClientWithMock()
	defer cleanup()

	httpmock.RegisterResponder("GET",
		withBaseURL("/resourcegroups/test-group?api-version=2020-10-01"),
		newJsonFileResponder("./testdata/resourceGroup.json"))

	r, err := c.LoadResourceGroup("test-group")
	assert.NoError(t, err)
	assert.Equal(t, "test-group", *r.Name)
}

func TestClient_LoadResourceGroupsByFilter(t *testing.T) {
	c, cleanup := testClientWithMock()
	defer cleanup()

	httpmock.RegisterResponder("GET",
		withBaseURL("/resourcegroups?%24filter=tagName+eq+%27Abteilung%27+and+tagValue+eq+%27Development%27&api-version=2020-10-01"),
		newJsonFileResponder("./testdata/resourceGroups.json"))

	groups, err := c.LoadResourceGroupsByFilter("Abteilung", "Development")
	assert.NoError(t, err)
	assert.Len(t, groups, 1)
	assert.Equal(t, "test-group", *groups[0].Name)

	httpmock.RegisterResponder("GET",
		withBaseURL("/resourcegroups?api-version=2020-10-01"),
		newJsonFileResponder("./testdata/resourceGroups.json"))

	groups, err = c.LoadResourceGroupsByFilter("", "")
	assert.NoError(t, err)
	assert.Len(t, groups, 1)
	assert.Equal(t, "test-group", *groups[0].Name)
}

func TestClient_LoadVmsByResourceGroup(t *testing.T) {
	c, cleanup := testClientWithMock()
	defer cleanup()

	httpmock.RegisterResponder("GET",
		withBaseURL("/resourceGroups/test-group/providers/Microsoft.Compute/virtualMachines?api-version=2020-06-01"),
		newJsonFileResponder("./testdata/vmsByResourceGroup.json"))

	httpmock.RegisterResponder("GET",
		withBaseURL("/resourceGroups/test-group/providers/Microsoft.Compute/virtualMachines/test-vm?%24expand=instanceView&api-version=2020-06-01"),
		newJsonFileResponder("./testdata/vmByName.json"))

	httpmock.RegisterResponder("GET",
		withBaseURL("/resourceGroups/test-group/providers/Microsoft.Compute/virtualMachines/test-vm2?%24expand=instanceView&api-version=2020-06-01"),
		newJsonFileResponder("./testdata/vmByName2.json"))

	vms, err := c.LoadVmsByResourceGroup("test-group")
	assert.NoError(t, err)
	assert.Len(t, vms.VirtualMachines, 2)
	assert.Equal(t, "test-vm", *vms.VirtualMachines[0].VirtualMachine.Name)
	assert.Equal(t, "test-vm2", *vms.VirtualMachines[1].VirtualMachine.Name)
}
