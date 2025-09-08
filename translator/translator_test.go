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
