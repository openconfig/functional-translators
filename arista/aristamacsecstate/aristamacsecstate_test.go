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

package aristamacsecstate

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"github.com/openconfig/functional-translators/ftutilities"
)

func TestTranslate(t *testing.T) {
	tests := []struct {
		name           string
		inputPath      string
		wantOutputPath string
		wantErr        bool
		wantNil        bool
	}{
		{
			name:           "secured_status_success",
			inputPath:      "testdata/secured_status_success_input.txt",
			wantOutputPath: "testdata/secured_status_success_output.txt",
		},
		{
			name:           "secured_status_success_for_target2",
			inputPath:      "testdata/secured_status_success_for_target2_input.txt",
			wantOutputPath: "testdata/secured_status_success_for_target2_output.txt",
		},
		{
			name:           "unencrypted_allowed_status_success",
			inputPath:      "testdata/unencrypted_allowed_status_success_input.txt",
			wantOutputPath: "testdata/unencrypted_allowed_status_success_output.txt",
		},
		{
			name:           "unencrypted_dropped_status_success",
			inputPath:      "testdata/unencrypted_dropped_status_success_input.txt",
			wantOutputPath: "testdata/unencrypted_dropped_status_success_output.txt",
		},
		{
			name:           "unknown_status_success",
			inputPath:      "testdata/unknown_status_success_input.txt",
			wantOutputPath: "testdata/unknown_status_success_output.txt",
		},
		{
			name:           "delete_cpstatus_only",
			inputPath:      "testdata/delete_cpstatus_only_input.txt",
			wantOutputPath: "testdata/delete_cpstatus_only_output.txt",
		},
		{
			name:           "delete_last_CKN",
			inputPath:      "testdata/delete_last_CKN_input.txt",
			wantOutputPath: "testdata/delete_last_CKN_output.txt",
		},
		{
			name:           "delete_cpstatus_and_port_status",
			inputPath:      "testdata/delete_cpstatus_and_port_status_input.txt",
			wantOutputPath: "testdata/delete_cpstatus_and_port_status_output.txt",
		},
		{
			name:           "update_eth22_and_delete_cpstatus_eth12",
			inputPath:      "testdata/update_eth22_and_delete_cpstatus_eth12_input.txt",
			wantOutputPath: "testdata/update_eth22_and_delete_cpstatus_eth12_output.txt",
		},
		{
			name:           "update_eth22_and_delete_port_status_eth12",
			inputPath:      "testdata/update_eth22_and_delete_port_status_eth12_input.txt",
			wantOutputPath: "testdata/update_eth22_and_delete_port_status_eth12_output.txt",
		},
		{
			name:           "unknown_status_breakout_interface_success",
			inputPath:      "testdata/unknown_status_breakout_interface_success_input.txt",
			wantOutputPath: "testdata/unknown_status_breakout_interface_success_output.txt",
		},
		{
			name:           "delete_one_CKN_validate_others_remain",
			inputPath:      "testdata/delete_one_CKN_validate_others_remain_input.txt",
			wantOutputPath: "testdata/delete_one_CKN_validate_others_remain_output.txt",
		},
		{
			name:           "CKN_map_not_initialized_and_incomplete",
			inputPath:      "testdata/CKN_map_not_initialized_and_incomplete_input.txt",
			wantOutputPath: "testdata/CKN_map_not_initialized_and_incomplete_output.txt",
			wantNil:        true,
		}, {
			name:           "verify_CKN_map_is_initialized_and_usable",
			inputPath:      "testdata/verify_CKN_map_is_initialized_and_usable_input.txt",
			wantOutputPath: "testdata/verify_CKN_map_is_initialized_and_usable_output.txt",
			wantNil:        true,
		},
		{
			name:      "invalidpath",
			inputPath: "testdata/invalidpath_input.txt",
			wantNil:   true,
		},
		{
			name:      "invaliddeletepath",
			inputPath: "testdata/invaliddeletepath_input.txt",
			wantNil:   true,
		},
		{
			name:      "invalidshortpath",
			inputPath: "testdata/invalidshortpath_input.txt",
			wantErr:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			inputSR, err := ftutilities.LoadSubscribeResponse(test.inputPath)
			if err != nil {
				t.Fatalf("ftutilities.LoadSubscribeResponse(%s) failed for input SubscribeRequest with err: %v", test.inputPath, err)
			}
			ft := New()
			gotSR, err := ft.Translate(inputSR)
			if err != nil {
				if !test.wantErr {
					t.Fatalf("ft.Translate(%v) failed with unexpected err: %v", inputSR, err)
				}
				// Expected error received.
				return
			}
			if test.wantErr {
				t.Fatalf("ft.Translate(%v) expected an error, but got nil", inputSR)
			}
			if test.wantNil {
				if gotSR != nil {
					t.Fatalf("ft.Translate(%v) expected nil, but got %v", inputSR, gotSR)
				}
				return // Expected nil output received.
			}
			wantSR, err := ftutilities.LoadSubscribeResponse(test.wantOutputPath)
			if err != nil {
				t.Fatalf("ftutilities.LoadSubscribeResponse(%s) failed for output SubscribeRequest with err: %v", test.wantOutputPath, err)
			}
			if diff := cmp.Diff(wantSR, gotSR, protocmp.Transform()); diff != "" {
				t.Fatalf("ft.Translate(%v) returned unexpected diff (-want +got):\n%s", inputSR, diff)
			}

			// inputSR, err := ftutilities.LoadSubscribeResponse(test.inputPath)
			// if err != nil {
			// 	t.Fatalf("ftutilities.LoadSubscribeResponse(%s) failed for input SubscribeRequest with err: %v", test.inputPath, err)
			// }
			// ft := New()

			// wantSR := &gnmipb.SubscribeResponse{}
			// if !test.wantNil && test.expectedOutput != nil {
			// 	wantSR = prepareExpectedSR(t, test.expectedOutput)
			// 	fmt.Println("output: ", wantSR)
			// }
			// sr := prepareInputSR(t, test.lazySR, test.extraUpdates)
			// fmt.Println("input: ", sr)
			// ft := New()
			// srCopy := proto.Clone(sr).(*gnmipb.SubscribeResponse)
			// gotSR, err := ft.Translate(srCopy)
			// if gotNil, gotErr := gotSR == nil, err != nil; gotNil || gotErr {
			// 	switch {
			// 	case gotErr != test.wantErr:
			// 		t.Fatalf("unexpected error result returned from translate() = %v, want error %t", gotErr, test.wantErr)
			// 	case err != nil:
			// 		return
			// 	case gotNil != test.wantNil:
			// 		t.Fatalf("unexpected nil result returned from translate() = %v, want nil %t", gotSR, test.wantNil)
			// 	default:
			// 		return
			// 	}
			// }
			// if gotNil := gotSR == nil; gotNil != test.wantNil {
			// 	t.Fatalf("unexpected nil result returned from translate() = %v, want nil %t", gotSR, test.wantNil)
			// }
			// if diff := cmp.Diff(wantSR, gotSR, protocmp.Transform()); diff != "" {
			// 	t.Fatalf("unexpected diff from translate (-want +got): %s", diff)
			// }
		})
	}
}

func TestFinalInterfacesForOCUpdate(t *testing.T) {
	tests := []struct {
		name                                  string
		interfaceSeenInput                    map[string]bool
		interfacesForOCDeleteInput            map[string]bool
		interfacesForOCUpdateFromDeletesInput map[string]bool
		expectedFinalOCUpdate                 map[string]bool
	}{
		{
			name:                                  "Interface seen, not for delete -> should be sent for update",
			interfaceSeenInput:                    map[string]bool{"Ethernet1": true},
			interfacesForOCDeleteInput:            map[string]bool{},
			interfacesForOCUpdateFromDeletesInput: map[string]bool{},
			expectedFinalOCUpdate:                 map[string]bool{"Ethernet1": true},
		},
		{
			name:                                  "Interface seen, is for delete -> should NOT be sent for update",
			interfaceSeenInput:                    map[string]bool{"Ethernet1": true},
			interfacesForOCDeleteInput:            map[string]bool{"Ethernet1": true},
			interfacesForOCUpdateFromDeletesInput: map[string]bool{},
			expectedFinalOCUpdate:                 map[string]bool{},
		},
		{
			name:                                  "Interface needs update from delete, not for full delete -> should be sent for update",
			interfaceSeenInput:                    map[string]bool{},
			interfacesForOCDeleteInput:            map[string]bool{},
			interfacesForOCUpdateFromDeletesInput: map[string]bool{"Ethernet2": true},
			expectedFinalOCUpdate:                 map[string]bool{"Ethernet2": true},
		},
		{
			name:                                  "Interface needs update from delete, is for full delete -> should NOT be sent for update",
			interfaceSeenInput:                    map[string]bool{},
			interfacesForOCDeleteInput:            map[string]bool{"Ethernet2": true},
			interfacesForOCUpdateFromDeletesInput: map[string]bool{"Ethernet2": true},
			expectedFinalOCUpdate:                 map[string]bool{},
		},
		{
			name:                                  "Interface seen and needs update from delete, not for full delete -> should be sent for update",
			interfaceSeenInput:                    map[string]bool{"Ethernet3": true},
			interfacesForOCDeleteInput:            map[string]bool{},
			interfacesForOCUpdateFromDeletesInput: map[string]bool{"Ethernet3": true},
			expectedFinalOCUpdate:                 map[string]bool{"Ethernet3": true},
		},
		{
			name:                                  "Interface seen and needs update from delete, is for full delete -> should NOT be sent for update",
			interfaceSeenInput:                    map[string]bool{"Ethernet3": true},
			interfacesForOCDeleteInput:            map[string]bool{"Ethernet3": true},
			interfacesForOCUpdateFromDeletesInput: map[string]bool{"Ethernet3": true},
			expectedFinalOCUpdate:                 map[string]bool{},
		},
		{
			name:                                  "All input maps are empty",
			interfaceSeenInput:                    map[string]bool{},
			interfacesForOCDeleteInput:            map[string]bool{},
			interfacesForOCUpdateFromDeletesInput: map[string]bool{},
			expectedFinalOCUpdate:                 map[string]bool{},
		},
		{
			name:                                  "Only interfacesForOCDeleteInput has entries",
			interfaceSeenInput:                    map[string]bool{},
			interfacesForOCDeleteInput:            map[string]bool{"Ethernet3": true},
			interfacesForOCUpdateFromDeletesInput: map[string]bool{},
			expectedFinalOCUpdate:                 map[string]bool{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actualFinalInterfacesForOCUpdate := make(map[string]bool)
			// Process interfaces seen from native updates
			for intfName := range tc.interfaceSeenInput {
				if !tc.interfacesForOCDeleteInput[intfName] {
					actualFinalInterfacesForOCUpdate[intfName] = true
				}
			}

			// Process interfaces needing update due to modifying deletes
			for intfName := range tc.interfacesForOCUpdateFromDeletesInput {
				if !tc.interfacesForOCDeleteInput[intfName] {
					actualFinalInterfacesForOCUpdate[intfName] = true
				}
			}

			if diff := cmp.Diff(tc.expectedFinalOCUpdate, actualFinalInterfacesForOCUpdate); diff != "" {
				t.Errorf("finalInterfacesForOCUpdate produced unexpected diff (-want +got):\n%s", diff)
			}
		})
	}
}
