// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package translator

import (
	"testing"
)

func TestMetadataMatch(t *testing.T) {
	tests := []struct {
		name       string
		ftMetaData []*FTMetadata
		wantMatch  bool
	}{
		{
			name:       "empty requirement catching all",
			ftMetaData: nil,
			wantMatch:  true,
		},
		{
			name: "partial requirement",
			ftMetaData: []*FTMetadata{
				{
					Vendor: "vendor",
				},
			},
			wantMatch: true,
		},
		{
			name: "requirement exact match",
			ftMetaData: []*FTMetadata{
				{
					Vendor:          "vendor",
					SoftwareVersion: "sw",
					HardwareModel:   "hw",
				},
			},
			wantMatch: true,
		},
		{
			name: "requirement case insensitive match",
			ftMetaData: []*FTMetadata{
				{
					Vendor:          "Vendor",
					SoftwareVersion: "SW",
					HardwareModel:   "HW",
				},
			},
			wantMatch: true,
		},
		{
			name: "match one",
			ftMetaData: []*FTMetadata{
				{
					Vendor: "non-vendor",
				},
				{
					Vendor: "vendor",
				},
			},
			wantMatch: true,
		},
		{
			name: "match all",
			ftMetaData: []*FTMetadata{
				{
					Vendor: "vendor",
				},
				{
					SoftwareVersion: "sw",
				},
			},
			wantMatch: true,
		},
		{
			name: "match none",
			ftMetaData: []*FTMetadata{
				{
					Vendor: "non-vendor",
				},
				{
					SoftwareVersion: "non-sw",
				},
				{
					HardwareModel: "non-hw",
				},
			},
			wantMatch: false,
		},
		{
			name: "vendor mismatch",
			ftMetaData: []*FTMetadata{
				{
					Vendor: "non-vendor",
				},
			},
			wantMatch: false,
		},
		{
			name: "sw mismatch",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersion: "non-sw",
				},
			},
			wantMatch: false,
		},
		{
			name: "hw mismatch",
			ftMetaData: []*FTMetadata{
				{
					HardwareModel: "non-hw",
				},
			},
			wantMatch: false,
		},
	}

	inputMetaData := &DeviceMetadata{
		Vendor:          "vendor",
		SoftwareVersion: "sw",
		HardwareModel:   "hw",
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var ft FunctionalTranslator
			ft.Metadata = tc.ftMetaData
			got := ft.metadataMatch(inputMetaData)
			if got != tc.wantMatch {
				t.Errorf("metadataMatch got %v, want %v", got, tc.wantMatch)
			}
		})
	}
}

func TestMetadataMatchSWRange(t *testing.T) {
	tests := []struct {
		name       string
		swVersion  string
		ftMetaData []*FTMetadata
		wantMatch  bool
	}{
		{
			name:      "sw range match",
			swVersion: "1.5",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.0",
						ExclusiveMax: "2.0",
					},
				},
			},
			wantMatch: true,
		},
		{
			name:      "sw range match, sw version longer than max",
			swVersion: "2.0.1",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.0",
						ExclusiveMax: "2.0",
					},
				},
			},
			wantMatch: false,
		},
		{
			name:      "sw range match, numeric suffix",
			swVersion: "1.5.1",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.0",
						ExclusiveMax: "2.0",
					},
				},
			},
			wantMatch: true,
		},
		{
			name:      "sw range match, letter suffix",
			swVersion: "4.34.2F-(some-random-suffix)",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "4.34",
						ExclusiveMax: "4.35",
					},
				},
			},
			wantMatch: true,
		},
		{
			name:      "sw range mismatch",
			swVersion: "2.5",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.0",
						ExclusiveMax: "2.0",
					},
				},
			},
			wantMatch: false,
		},

		{
			name:      "sw range mismatch, numeric suffix",
			swVersion: "2.5.1",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.0",
						ExclusiveMax: "2.0",
					},
				},
			},
			wantMatch: false,
		},
		{
			name:      "sw range match, letter suffix",
			swVersion: "4.34.2F",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "4.34",
						ExclusiveMax: "4.35",
					},
				},
			},
			wantMatch: true,
		},
		{
			name:      "sw range match, equal min",
			swVersion: "1.0",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.0",
						ExclusiveMax: "2.0",
					},
				},
			},
			wantMatch: true,
		},
		{
			name:      "sw range mismatch, equal max",
			swVersion: "2.0",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.0",
						ExclusiveMax: "2.0",
					},
				},
			},
			wantMatch: false,
		},
		{
			name:      "sw range match, letter comparison",
			swVersion: "1.5B",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.5A",
						ExclusiveMax: "1.5C",
					},
				},
			},
			wantMatch: true,
		},
		{
			name:      "sw range mismatch, after, letter comparison",
			swVersion: "1.5D",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.5A",
						ExclusiveMax: "1.5C",
					},
				},
			},
			wantMatch: false,
		},
		{
			name:      "sw range mismatch, before, letter comparison",
			swVersion: "1.5A",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.5B",
						ExclusiveMax: "1.5C",
					},
				},
			},
			wantMatch: false,
		},
		{
			name:      "sw range match, letter is greater than number",
			swVersion: "1.A",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.9",
						ExclusiveMax: "1.B",
					},
				},
			},
			wantMatch: true,
		},
		{
			name:      "sw outside range, letter is greater than number",
			swVersion: "1.B",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.8",
						ExclusiveMax: "1.9",
					},
				},
			},
			wantMatch: false,
		},
		{
			name:      "sw range mismatch, less than min",
			swVersion: "0.9",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.0",
						ExclusiveMax: "2.0",
					},
				},
			},
			wantMatch: false,
		},
		{
			name:      "longer version in range, inside range",
			swVersion: "1.3",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.2.3.4.5",
						ExclusiveMax: "2.1.0",
					},
				},
			},
			wantMatch: true,
		},
		{
			name:      "longer version in range, outside range",
			swVersion: "1.3",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.2.3.4.5",
						ExclusiveMax: "1.3.0",
					},
				},
			},
			wantMatch: false,
		},
		{
			name:      "sw range match, identical letter components",
			swVersion: "1.A.2",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.A.1",
						ExclusiveMax: "1.A.3",
					},
				},
			},
			wantMatch: true,
		},
		{
			name:      "sw version is prefix of min version",
			swVersion: "1.2",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.2.3",
						ExclusiveMax: "2.0",
					},
				},
			},
			wantMatch: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var ft FunctionalTranslator
			ft.Metadata = tc.ftMetaData
			got := ft.metadataMatch(&DeviceMetadata{SoftwareVersion: tc.swVersion})
			if got != tc.wantMatch {
				t.Errorf("swVersionMatch got %v, want %v for swVersion %s against ft: %v", got, tc.wantMatch, tc.swVersion, ft.Metadata)
			}
		})
	}
}
