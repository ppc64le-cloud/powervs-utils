package utils

import (
	"testing"
)

func TestGetSystemSpecifications(t *testing.T) {
	tests := []struct {
		testcasename string
		systemType   string
		wantSpec     SystemSpecifications
		wantErr      bool
	}{
		{
			testcasename: "known system type s922",
			systemType:   "s922",
			wantSpec:     SystemSpecifications{MaxMemoryGiB: 877, MaxProcessor: 15},
		},
		{
			testcasename: "known system type e1080",
			systemType:   "e1080",
			wantSpec:     SystemSpecifications{MaxMemoryGiB: 29802, MaxProcessor: 165},
		},
		{
			testcasename: "known system type s1122",
			systemType:   "s1122",
			wantSpec:     SystemSpecifications{MaxMemoryGiB: 31840, MaxProcessor: 229},
		},
		{
			testcasename: "unknown system type returns error",
			systemType:   "x9999",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testcasename, func(t *testing.T) {
			spec, err := GetSystemSpecifications(tt.systemType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSystemSpecifications(%q) error = %v, wantErr %v", tt.systemType, err, tt.wantErr)
				return
			}
			if !tt.wantErr && spec != tt.wantSpec {
				t.Errorf("GetSystemSpecifications(%q) = %+v, want %+v", tt.systemType, spec, tt.wantSpec)
			}
		})
	}
}

func TestSystemTypeNames(t *testing.T) {
	expected := map[string]bool{
		"s922":  true,
		"e980":  true,
		"s1022": true,
		"e1050": true,
		"e1080": true,
		"s1122": true,
	}

	types := SystemTypeNames()

	if len(types) != len(expected) {
		t.Errorf("SystemTypeNames() returned %d types, want %d", len(types), len(expected))
	}

	for _, st := range types {
		if !expected[st] {
			t.Errorf("SystemTypeNames() returned unexpected type %q", st)
		}
	}
}

func TestValidateSystemType(t *testing.T) {
	tests := []struct {
		testcasename string
		systemType   string
		want         bool
	}{
		{
			testcasename: "known type s922",
			systemType:   "s922",
			want:         true,
		},
		{
			testcasename: "known type e1080",
			systemType:   "e1080",
			want:         true,
		},
		{
			testcasename: "unknown type",
			systemType:   "x9999",
			want:         false,
		},
		{
			testcasename: "empty string",
			systemType:   "",
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testcasename, func(t *testing.T) {
			got := ValidateSystemType(tt.systemType)
			if got != tt.want {
				t.Errorf("ValidateSystemType(%q) = %v, want %v", tt.systemType, got, tt.want)
			}
		})
	}
}
