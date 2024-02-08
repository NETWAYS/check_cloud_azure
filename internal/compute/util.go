package compute

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/NETWAYS/go-check"
)

func LevelToState(level armcompute.StatusLevelTypes) int {
	switch level {
	case armcompute.StatusLevelTypesInfo:
		return check.OK
	case armcompute.StatusLevelTypesWarning:
		return check.Warning
	case armcompute.StatusLevelTypesError:
		return check.Critical
	default:
		return check.Unknown
	}
}
