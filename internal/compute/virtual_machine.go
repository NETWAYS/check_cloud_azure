package compute

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
)

type VirtualMachine struct {
	VirtualMachine *armcompute.VirtualMachine
}

func (v *VirtualMachine) GetPartialResult() result.PartialResult {
	result := result.PartialResult{}

	result.SetState(v.GetStatus())

	result.Output = v.GetOutput() + "\n\n" + v.GetLongOutput()

	return result
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
	out += fmt.Sprintf("Size: %s\n", *v.VirtualMachine.Properties.HardwareProfile.VMSize)
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
	states := []int{
		check.Unknown,
		check.Unknown,
		check.Unknown,
	}

	// field Level StatusLevelTypes`json:"level,omitempty"` of InstanceViewStatus type
	// Level - The level code. Possible values include: 'Info', 'Warning', 'Error'
	_, provLevel := v.GetProvisioningState()
	states[0] = LevelToState(provLevel)

	_, agentLevel := v.GetAgentProvisioningState()
	states[2] = LevelToState(agentLevel)

	// https://docs.microsoft.com/en-us/dotnet/api/microsoft.azure.management.compute.fluent.powerstate?view=azure-dotnet#fields
	state, level := v.GetPowerState()
	if level == armcompute.StatusLevelTypesInfo {
		if state == "deallocated" {
			states[1] = check.Critical
		} else {
			states[1] = check.OK
		}
	} else {
		states[1] = LevelToState(provLevel)
	}

	return result.WorstState(states...)
}

func (v *VirtualMachine) GetPowerState() (string, armcompute.StatusLevelTypes) {
	for _, state := range v.VirtualMachine.Properties.InstanceView.Statuses {
		if strings.HasPrefix(*state.Code, "PowerState/") {
			return strings.SplitN(*state.Code, "/", 2)[1], *state.Level
		}
	}

	return "", ""
}

func (v *VirtualMachine) GetProvisioningState() (string, armcompute.StatusLevelTypes) {
	for _, state := range v.VirtualMachine.Properties.InstanceView.Statuses {
		if strings.HasPrefix(*state.Code, "ProvisioningState/") {
			return strings.SplitN(*state.Code, "/", 2)[1], *state.Level
		}
	}

	return "", ""
}

func (v *VirtualMachine) GetAgentProvisioningState() (string, armcompute.StatusLevelTypes) {
	if v.VirtualMachine.Properties.InstanceView.VMAgent == nil {
		return "(none)", armcompute.StatusLevelTypesError
	}

	for _, state := range v.VirtualMachine.Properties.InstanceView.VMAgent.Statuses {
		if strings.HasPrefix(*state.Code, "ProvisioningState/") {
			return strings.SplitN(*state.Code, "/", 2)[1], *state.Level
		}
	}

	return "", ""
}
