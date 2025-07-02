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

package ciscoxrlaser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"google3/third_party/openconfig/functional_translators/ftutilities"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

func TestTranslate(t *testing.T) {
	tests := []struct {
		name           string
		inputPath      string
		wantOutputPath string
		expectError    bool
	}{
		{
			name:           "success",
			inputPath:      "testdata/success_input.txt",
			wantOutputPath: "testdata/success_output.txt",
		},
		{
			name:           "Temp thresholds not present",
			inputPath:      "testdata/no_temp_input.txt",
			wantOutputPath: "testdata/no_temp_output.txt",
		},
		{
			name:           "RX thresholds not present",
			inputPath:      "testdata/no_rx_input.txt",
			wantOutputPath: "testdata/no_rx_output.txt",
		},
		{
			name:           "TX thresholds not present",
			inputPath:      "testdata/no_tx_input.txt",
			wantOutputPath: "testdata/no_tx_output.txt",
		},
		{
			name:        "path with no info",
			inputPath:   "testdata/no_info_input.txt",
			expectError: true,
		},
		{
			name:        "path with no port",
			inputPath:   "testdata/no_port_input.txt",
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inputSR, err := ftutilities.LoadSubscribeResponse(test.inputPath)
			if err != nil {
				t.Fatalf("Failed to load input message: %v", err)
			}
			gotSR, err := translate(inputSR)
			errorMatchesExpectation := (err != nil) == test.expectError
			if !errorMatchesExpectation {
				t.Fatalf("Unexpected error result returned from translate() = %v, want error %t", err, test.expectError)
			}
			if !test.expectError {
				wantSR, err := ftutilities.LoadSubscribeResponse(test.wantOutputPath)
				if err != nil {
					t.Fatalf("Failed to load want message: %v", err)
				}
				if diff := cmp.Diff(wantSR, gotSR, protocmp.Transform(), protocmp.SortRepeatedFields(&gnmipb.Notification{}, "update")); diff != "" {
					t.Fatalf("Unexpected diff from translate() = %v, want %v, diff(-want +got):\n%s", gotSR, wantSR, diff)
				}
			}
		})
	}

}
