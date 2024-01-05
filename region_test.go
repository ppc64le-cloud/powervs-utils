package utils

import (
	"testing"
)

func TestGetRegion(t *testing.T) {
	type args struct {
		zone string
	}
	tests := []struct {
		name       string
		args       args
		wantRegion string
		wantErr    bool
	}{
		{
			"London",
			args{"lon06"},
			"lon",
			false,
		},
		{
			"Dallas",
			args{"us-south"},
			"us-south",
			false,
		},
		{
			"Dallas 12",
			args{"dal12"},
			"dal",
			false,
		},
		{
			"Sao Paulo 01",
			args{"sao01"},
			"sao",
			false,
		},
		{
			"Sao Paulo 04",
			args{"sao04"},
			"sao",
			false,
		},
		{
			"Washington DC",
			args{"us-east"},
			"us-east",
			false,
		},
		{
			"Washington DC 06",
			args{"wdc06"},
			"wdc",
			false,
		},
		{
			"Washington DC 07",
			args{"wdc07"},
			"wdc",
			false,
		},
		{
			"Toronto",
			args{"tor01"},
			"tor",
			false,
		},
		{
			"Frankfurt",
			args{"eu-de-1"},
			"eu-de",
			false,
		},
		{
			"Sydney",
			args{"syd01"},
			"syd",
			false,
		},
		{
			"India",
			args{"blr01"},
			"",
			true,
		},
		{
			"Tokyo",
			args{"tok04"},
			"tok",
			false,
		},
		{
			"Montreal",
			args{"mon01"},
			"mon",
			false,
		},
		{
			"Osaka",
			args{"osa21"},
			"osa",
			false,
		},
		{
			"Madrid 02",
			args{"mad02"},
			"mad",
			false,
		},
		{
			"Madrid 04",
			args{"mad04"},
			"mad",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRegion, err := GetRegion(tt.args.zone)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRegion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRegion != tt.wantRegion {
				t.Errorf("GetRegion() gotRegion = %v, want %v", gotRegion, tt.wantRegion)
			}
		})
	}
}

func TestVPCRegionForPowerVSRegion(t *testing.T) {
	type args struct {
		region string
	}
	tests := []struct {
		name       string
		args       args
		wantRegion string
		wantErr    bool
	}{
		{
			"Dallas",
			args{"dal"},
			"us-south",
			false,
		},
		{
			"Osaka",
			args{"eu-de1"},
			"eu-de",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vpcRegion, err := VPCRegionForPowerVSRegion(tt.args.region)

			if err != nil {
				if !tt.wantErr {
					t.Errorf("VPCRegionForPowerVSRegion() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if vpcRegion != tt.wantRegion {
				t.Errorf("VPCRegionForPowerVSRegion() gotRegion = %v, want %v", vpcRegion, tt.wantRegion)
			}
		})
	}
}
