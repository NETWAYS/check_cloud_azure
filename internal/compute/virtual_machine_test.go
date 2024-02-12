package compute_test

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// nolint: bodyclose
func TestVirtualMachine(t *testing.T) {
	c, cleanup := testClientWithMock()
	defer cleanup()

	httpmock.RegisterResponder("GET",
		withBaseURL("/resourceGroups/test-group/providers/Microsoft.Compute/virtualMachines/test-vm?%24expand=instanceView&api-version=2020-06-01"),
		newJsonFileResponder("./testdata/vmByName.json"))

	vm, err := c.LoadVmByName("test-group", "test-vm")
	require.NoError(t, err)

	assert.Equal(t, 0, vm.GetStatus())

	output := vm.GetOutput()
	assert.Contains(t, output, "power=running")
	assert.Contains(t, output, "provision=succeeded")
	assert.Contains(t, output, "agent=succeeded")

	long := vm.GetLongOutput()
	assert.Contains(t, long, "Size: Standard_D2s_v3")
	assert.Contains(t, long, "Location: germanywestcentral")
}

// nolint: bodyclose
func TestVirtualMachine_deallocated(t *testing.T) {
	c, cleanup := testClientWithMock()
	defer cleanup()

	httpmock.RegisterResponder("GET",
		withBaseURL("/resourceGroups/test-group/providers/Microsoft.Compute/virtualMachines/test-vm2?%24expand=instanceView&api-version=2020-06-01"),
		newJsonFileResponder("./testdata/vmByName-deallocated.json"))

	vm, err := c.LoadVmByName("test-group", "test-vm2")
	require.NoError(t, err)

	assert.Equal(t, 2, vm.GetStatus())

	output := vm.GetOutput()
	assert.Contains(t, output, "power=deallocated")
	assert.Contains(t, output, "provision=succeeded")
	assert.Contains(t, output, "agent=(none)")

	long := vm.GetLongOutput()
	assert.Contains(t, long, "Size: Standard_B1s")
	assert.Contains(t, long, "Location: germanywestcentral")
}
