package compute

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
	"strings"
)

type VirtualMachine struct {
	VirtualMachine *compute.VirtualMachine
}

func (v *VirtualMachine) GetOutput() (out string) {
	instance := v.VirtualMachine

	powerState, _ := v.GetPowerState()
	provState, _ := v.GetProvisioningState()
	agentState, _ := v.GetAgentProvisioningState()

	out = fmt.Sprintf(`"%s" provision=%s power=%s agent=%s`,
		*instance.Name, provState, powerState, agentState)

	return
}

func (v *VirtualMachine) GetLongOutput() (out string) {
	out += fmt.Sprintf("Size: %s\n", v.VirtualMachine.HardwareProfile.VMSize)
	out += "Location: " + *v.VirtualMachine.Location + "\n"

	/* Maybe helpful later
	if v.VirtualMachine.HostGroup != nil {
		hostGroup := *v.VirtualMachine.HostGroup
		out += fmt.Sprintf("Hostgroup: %s\n", *hostGroup.ID)
	} else {
		out += "Hostgroup: (none)"
	}
	*/

	return
}

func (v *VirtualMachine) GetStatus() int {
	states := []int{3, 3, 3}

	// field Level StatusLevelTypes`json:"level,omitempty"` of InstanceViewStatus type
	// Level - The level code. Possible values include: 'Info', 'Warning', 'Error'
	_, provLevel := v.GetProvisioningState()
	states[0] = LevelToState(provLevel)

	_, agentLevel := v.GetAgentProvisioningState()
	states[2] = LevelToState(agentLevel)

	// https://docs.microsoft.com/en-us/dotnet/api/microsoft.azure.management.compute.fluent.powerstate?view=azure-dotnet#fields
	switch state, level := v.GetPowerState(); level {
	case compute.Info:
		if state == "deallocated" {
			states[1] = check.Critical
		} else {
			states[1] = check.OK
		}
	default:
		states[1] = LevelToState(provLevel)
	}

	return result.WorstState(states...)
}

func (v *VirtualMachine) GetPowerState() (string, compute.StatusLevelTypes) {
	for _, state := range *v.VirtualMachine.InstanceView.Statuses {
		if strings.HasPrefix(*state.Code, "PowerState/") {
			return strings.SplitN(*state.Code, "/", 2)[1], state.Level
		}
	}

	return "", ""
}

func (v *VirtualMachine) GetProvisioningState() (string, compute.StatusLevelTypes) {
	for _, state := range *v.VirtualMachine.InstanceView.Statuses {
		if strings.HasPrefix(*state.Code, "ProvisioningState/") {
			return strings.SplitN(*state.Code, "/", 2)[1], state.Level
		}
	}

	return "", ""
}

func (v *VirtualMachine) GetAgentProvisioningState() (string, compute.StatusLevelTypes) {
	if v.VirtualMachine.InstanceView.VMAgent == nil {
		return "(none)", compute.Error
	}

	for _, state := range *v.VirtualMachine.InstanceView.VMAgent.Statuses {
		if strings.HasPrefix(*state.Code, "ProvisioningState/") {
			return strings.SplitN(*state.Code, "/", 2)[1], state.Level
		}
	}

	return "", ""
}
