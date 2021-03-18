package compute

import (
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/NETWAYS/go-check"
)

func LevelToState(level compute.StatusLevelTypes) int {
	switch level {
	case compute.Info:
		return check.OK
	case compute.Warning:
		return check.Warning
	case compute.Error:
		return check.Critical
	default:
		return check.Unknown
	}
}
