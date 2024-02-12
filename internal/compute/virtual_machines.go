package compute

import (
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
)

type VirtualMachines struct {
	VirtualMachines []*VirtualMachine
}

func (m *VirtualMachines) GetPartialResult() result.PartialResult {
	result := result.PartialResult{}
	result.SetDefaultState(check.Unknown)

	if m.IsEmpty() {
		result.Output = "no VMs found"
		return result
	}

	for idx := range m.VirtualMachines {
		result.AddSubcheck(m.VirtualMachines[idx].GetPartialResult())
	}

	return result
}

func (m VirtualMachines) IsEmpty() bool {
	return len(m.VirtualMachines) == 0
}
