package utils

import "fmt"

// SystemSpecifications holds the hardware capacity limits for a PowerVS system type.
// Source: https://www.ibm.com/docs/en/power-virtual-server?topic=pricing-power-virtual-server-in-data-centers
type SystemSpecifications struct {
	// MaxMemoryGiB is the maximum allocatable memory in GiB for the system type.
	MaxMemoryGiB int
	// MaxProcessor is the maximum number of vCPUs allocatable for the system type.
	MaxProcessor float64
}

// powerVSSystemConfigurations contains the hardware capacity limits for each known PowerVS system type.
// Source: https://www.ibm.com/docs/en/power-virtual-server?topic=pricing-power-virtual-server-in-data-centers
var powerVSSystemConfigurations = map[string]SystemSpecifications{
	"s922": {MaxMemoryGiB: 877, MaxProcessor: 15},
	"e980": {MaxMemoryGiB: 14255, MaxProcessor: 143},
	"s1022": {MaxMemoryGiB: 1847, MaxProcessor: 33},
	"e1050": {MaxMemoryGiB: 7629, MaxProcessor: 87},
	"e1080": {MaxMemoryGiB: 29802, MaxProcessor: 165},
	"s1122": {MaxMemoryGiB: 31840, MaxProcessor: 229},
}

// GetSystemSpecifications returns the SystemSpecifications for the given PowerVS system type.
// An error is returned when the system type is not known.
func GetSystemSpecifications(systemType string) (SystemSpecifications, error) {
	spec, ok := powerVSSystemConfigurations[systemType]
	if !ok {
		return SystemSpecifications{}, fmt.Errorf("unknown system type %q: no specifications found", systemType)
	}
	return spec, nil
}

// SystemTypeNames returns the list of PowerVS system type names that have hardware specifications defined.
func SystemTypeNames() []string {
	types := make([]string, 0, len(powerVSSystemConfigurations))
	for k := range powerVSSystemConfigurations {
		types = append(types, k)
	}
	return types
}

// ValidateSystemType reports whether the given PowerVS system type has hardware specifications defined.
func ValidateSystemType(systemType string) bool {
	_, ok := powerVSSystemConfigurations[systemType]
	return ok
}
