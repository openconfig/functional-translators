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

package aristamacseccounters

import (
	"maps"
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
		leafMapping    map[string]string
		wantErr        bool
		wantNil        bool
	}{
		{
			name:           "mismatch_password_translation_success",
			inputPath:      "testdata/mismatch_password_translation_success_input.txt",
			wantOutputPath: "testdata/mismatch_password_translation_success_output.txt",
		},
		{
			name:           "unknown_ckn_translation_success",
			inputPath:      "testdata/unknown_ckn_translation_success_input.txt",
			wantOutputPath: "testdata/unknown_ckn_translation_success_output.txt",
		},
		{
			name:           "leaf_mapping_missing",
			inputPath:      "testdata/leaf_mapping_missing_input.txt",
			wantOutputPath: "testdata/leaf_mapping_missing_output.txt",
			leafMapping:    map[string]string{"icvValidationErr": "rx-badicv-pkts"},
			wantErr:        true,
		},
		{
			name:           "txPktsErrin_translation_success",
			inputPath:      "testdata/txpktserrin_translation_success_input.txt",
			wantOutputPath: "testdata/txpktserrin_translation_success_output.txt",
		},
		{
			name:           "txPktsCtrl_translation_success",
			inputPath:      "testdata/txpktsctrl_translation_success_input.txt",
			wantOutputPath: "testdata/txpktsctrl_translation_success_output.txt",
		},
		{
			name:           "rxPktsCtrl_translation_success",
			inputPath:      "testdata/rxpktsctrl_translation_success_input.txt",
			wantOutputPath: "testdata/rxpktsctrl_translation_success_output.txt",
		},
		{
			name:           "txPktsDropped_translation_success",
			inputPath:      "testdata/txpktsdropped_translation_success_input.txt",
			wantOutputPath: "testdata/txpktsdropped_translation_success_output.txt",
		},
		{
			name:           "txPktsDropped_translation_success_with_json",
			inputPath:      "testdata/txpktsdropped_translation_success_with_json_input.txt",
			wantOutputPath: "testdata/txpktsdropped_translation_success_with_json_output.txt",
		},
		{
			name:           "rxPktsDropped_translation_success",
			inputPath:      "testdata/rxpktsdropped_translation_success_input.txt",
			wantOutputPath: "testdata/rxpktsdropped_translation_success_output.txt",
		},
		{
			name:      "invalid_json_format",
			inputPath: "testdata/invalid_json_format_input.txt",
			wantErr:   true,
		},
		{
			name:      "json_missing_value_key",
			inputPath: "testdata/json_missing_value_key_input.txt",
			wantErr:   true,
		},
		{
			name:      "json_value_is_not_a_float",
			inputPath: "testdata/json_value_is_not_a_float_input.txt",
			wantErr:   true,
		},
		{
			name:      "invalid_input_path",
			inputPath: "testdata/invalid_input_path_input.txt",
			wantNil:   true,
		},
		{
			name:      "invalid_delete_path",
			inputPath: "testdata/invalid_delete_path_input.txt",
			wantNil:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if len(test.leafMapping) > 0 {
				origMap := maps.Clone(vendorToOCLeaf)
				defer func() {
					vendorToOCLeaf = origMap
				}()
				vendorToOCLeaf = test.leafMapping
			}
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
		})
	}
}
