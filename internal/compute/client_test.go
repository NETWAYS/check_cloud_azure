// nolint:bodyclose
package compute_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/NETWAYS/check_cloud_azure/internal/compute"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	AzureToken          azcore.TokenCredential
	AzureSubscriptionId = "SUB-UUID"
)

func testClientWithMock() (client *compute.Client, cleanup func()) {
	client = compute.NewClient(AzureToken, AzureSubscriptionId)

	httpmock.Activate()
	cleanup = httpmock.DeactivateAndReset

	return
}

func withBaseURL(url string) string {
	return "https://management.azure.com/subscriptions/" + AzureSubscriptionId + url
}

func newJsonFileResponder(fileName string) func(request *http.Request) (*http.Response, error) {
	return func(request *http.Request) (*http.Response, error) {
		data, err := os.ReadFile(fileName)
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
	require.NoError(t, err)
	assert.Equal(t, "test-vm", *vm.VirtualMachine.Name)
}

func TestClient_LoadResourceGroup(t *testing.T) {
	c, cleanup := testClientWithMock()
	defer cleanup()

	httpmock.RegisterResponder("GET",
		withBaseURL("/resourcegroups/test-group?api-version=2020-10-01"),
		newJsonFileResponder("./testdata/resourceGroup.json"))

	r, err := c.LoadResourceGroup("test-group")
	require.NoError(t, err)
	assert.Equal(t, "test-group", *r.Name)
}

func TestClient_LoadResourceGroupsByFilter(t *testing.T) {
	c, cleanup := testClientWithMock()
	defer cleanup()

	httpmock.RegisterResponder("GET",
		withBaseURL("/resourcegroups?%24filter=tagName+eq+%27Abteilung%27+and+tagValue+eq+%27Development%27&api-version=2020-10-01"),
		newJsonFileResponder("./testdata/resourceGroups.json"))

	groups, err := c.LoadResourceGroupsByFilter("Abteilung", "Development")
	require.NoError(t, err)
	assert.Len(t, groups, 1)
	assert.Equal(t, "test-group", *groups[0].Name)

	httpmock.RegisterResponder("GET",
		withBaseURL("/resourcegroups?api-version=2020-10-01"),
		newJsonFileResponder("./testdata/resourceGroups.json"))

	groups, err = c.LoadResourceGroupsByFilter("", "")
	require.NoError(t, err)
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
	require.NoError(t, err)
	assert.Len(t, vms.VirtualMachines, 2)
	assert.Equal(t, "test-vm", *vms.VirtualMachines[0].VirtualMachine.Name)
	assert.Equal(t, "test-vm2", *vms.VirtualMachines[1].VirtualMachine.Name)
}
