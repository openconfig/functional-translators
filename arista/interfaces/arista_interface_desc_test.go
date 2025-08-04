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

package aristainterfacedesc

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
		expectError    bool
	}{
		{
			name:           "success",
			inputPath:      "testdata/interface_desc_success_input.txt",
			wantOutputPath: "testdata/interface_desc_success_output.txt",
		},
		{
			name:           "delete only",
			inputPath:      "testdata/interface_desc_delete_only_input.txt",
			wantOutputPath: "testdata/interface_desc_delete_only_output.txt",
		},
		{
			name:           "delete interface",
			inputPath:      "testdata/interface_desc_delete_interface_input.txt",
			wantOutputPath: "testdata/interface_desc_delete_interface_output.txt",
		},
		{
			name:           "delete interfaces",
			inputPath:      "testdata/interface_desc_delete_interfaces_input.txt",
			wantOutputPath: "testdata/interface_desc_delete_interfaces_output.txt",
		},
		{
			name:        "invalid type",
			inputPath:   "testdata/interface_desc_invalid_type_input.txt",
			expectError: true,
		},
		{
			name:        "invalid path",
			inputPath:   "testdata/interface_desc_invalid_path_input.txt",
			expectError: true,
		},
		{
			name:        "unrelated delete",
			inputPath:   "testdata/interface_desc_unrelated_delete_input.txt",
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ft := New()
			inputSR, err := ftutilities.LoadSubscribeResponse(test.inputPath)
			if err != nil {
				t.Fatalf("Failed to load input message: %v", err)
			}
			gotSR, gotErr := ft.Translate(inputSR)
			errorMatchesExpectation := (gotErr != nil) == test.expectError
			if !errorMatchesExpectation {
				t.Fatalf("unexpected error result returned from translate() = %v, want error %t", gotErr, test.expectError)
			}
			if !test.expectError {
				wantSR, err := ftutilities.LoadSubscribeResponse(test.wantOutputPath)
				if err != nil {
					t.Fatalf("Failed to load want message: %v", err)
				}
				if diff := cmp.Diff(wantSR, gotSR, protocmp.Transform()); diff != "" {
					t.Fatalf("unexpected diff from translate() = %v, want %v:\n%s", gotSR, wantSR, diff)
				}
			}
		})
	}
}
