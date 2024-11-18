package compute_test

import (
	"github.com/NETWAYS/check_cloud_azure/internal/compute"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// nolint: bodyclose
func TestVirtualMachines(t *testing.T) {
	c, cleanup := testClientWithMock()
	defer cleanup()

	httpmock.RegisterResponder("GET",
		withBaseURL("/resourceGroups/test-group/providers/Microsoft.Compute/virtualMachines/test-vm?%24expand=instanceView&api-version=2020-06-01"),
		newJsonFileResponder("./testdata/vmByName.json"))

	httpmock.RegisterResponder("GET",
		withBaseURL("/resourceGroups/test-group/providers/Microsoft.Compute/virtualMachines/test-vm2?%24expand=instanceView&api-version=2020-06-01"),
		newJsonFileResponder("./testdata/vmByName-deallocated.json"))

	running, err := c.LoadVmByName("test-group", "test-vm")
	require.NoError(t, err)

	deallocated, err := c.LoadVmByName("test-group", "test-vm2")
	require.NoError(t, err)

	vms := compute.VirtualMachines{}

	assert.Equal(t, 2, vms.GetStatus())
	assert.Contains(t, vms.GetOutput(), "no VMs found")

	vms.VirtualMachines = append(vms.VirtualMachines, running)

	assert.Equal(t, 0, vms.GetStatus())
	assert.Contains(t, vms.GetOutput(), "[OK]")

	vms.VirtualMachines = append(vms.VirtualMachines, deallocated)

	assert.Equal(t, 2, vms.GetStatus())
	assert.Contains(t, vms.GetOutput(), "[OK]")
	assert.Contains(t, vms.GetOutput(), "[CRITICAL]")
}
