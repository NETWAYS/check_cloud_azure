package compute

import (
	"fmt"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
)

type VirtualMachines struct {
	VirtualMachines []*VirtualMachine
}

func (m VirtualMachines) GetStatus() int {
	if m.IsEmpty() {
		return check.Critical
	}

	var states []int

	for _, vm := range m.VirtualMachines {
		states = append(states, vm.GetStatus())
	}

	return result.WorstState(states...)
}

func (m VirtualMachines) GetOutput() (output string) {
	if m.IsEmpty() {
		return "no VMs found"
	}

	for _, vm := range m.VirtualMachines {
		output += fmt.Sprintf("[%s] %s\n", check.StatusText(vm.GetStatus()), vm.GetOutput())
	}

	return
}

func (m VirtualMachines) IsEmpty() bool {
	return len(m.VirtualMachines) == 0
}
