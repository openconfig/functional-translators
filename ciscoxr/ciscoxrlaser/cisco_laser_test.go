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
	"github.com/openconfig/functional-translators/ftutilities"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

func TestTranslate(t *testing.T) {
	tests := []struct {
		name           string
		inputPath      string
		wantOutputPath string
		expectError    bool
		wantNil        bool
	}{
		{
			name:           "Success",
			inputPath:      "testdata/success_input.txt",
			wantOutputPath: "testdata/success_output.txt",
		},
		{
			name:           "TempThresholdsNotPresent",
			inputPath:      "testdata/no_temp_input.txt",
			wantOutputPath: "testdata/no_temp_output.txt",
		},
		{
			name:           "RXThresholdsNotPresent",
			inputPath:      "testdata/no_rx_input.txt",
			wantOutputPath: "testdata/no_rx_output.txt",
		},
		{
			name:           "TXThresholdsNotPresent",
			inputPath:      "testdata/no_tx_input.txt",
			wantOutputPath: "testdata/no_tx_output.txt",
		},
		{
			name:        "PathWithNoInfo",
			inputPath:   "testdata/no_info_input.txt",
			expectError: true,
		},
		{
			name:        "PathWithNoPort",
			inputPath:   "testdata/no_port_input.txt",
			expectError: true,
		},
		{
			name:      "NilOpticsOper",
			inputPath: "testdata/nil_optics_oper_input.txt",
			wantNil:   true,
		},
		{
			name:      "NilOpticsPorts",
			inputPath: "testdata/nil_optics_ports_input.txt",
			wantNil:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inputSR, err := ftutilities.LoadSubscribeResponse(test.inputPath)
			if err != nil {
				t.Fatalf("ftutilities.LoadSubscribeResponse(%q) failed: %v", test.inputPath, err)
			}
			gotSR, err := New().Translate(inputSR)
			errorMatchesExpectation := (err != nil) == test.expectError
			if !errorMatchesExpectation {
				t.Fatalf("Translate(%q) error got %v, want error %t", test.inputPath, err, test.expectError)
			}
			if !test.expectError {
				if test.wantNil {
					if gotSR != nil {
						t.Fatalf("Translate(%q) got %v, want nil", test.inputPath, gotSR)
					}
				} else {
					wantSR, err := ftutilities.LoadSubscribeResponse(test.wantOutputPath)
					if err != nil {
						t.Fatalf("ftutilities.LoadSubscribeResponse(%q) failed: %v", test.wantOutputPath, err)
					}
					if diff := cmp.Diff(wantSR, gotSR, protocmp.Transform(), protocmp.SortRepeatedFields(&gnmipb.Notification{}, "update")); diff != "" {
						t.Fatalf("Translate(%q) got %v, want %v, diff(-want +got):\n%s", test.inputPath, gotSR, wantSR, diff)
					}
				}
			}
		})
	}

}
