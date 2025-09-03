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

package aristacfmstate

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
		wantNil        bool
		wantErr        bool
	}{
		{
			name:           "RdiTxCondition translation success",
			inputPath:      "testdata/rdi_tx_condition_input.txt",
			wantOutputPath: "testdata/rdi_tx_condition_output.txt",
		},
		{
			name:           "Delete translation success",
			inputPath:      "testdata/delete_input.txt",
			wantOutputPath: "testdata/delete_output.txt",
		},
		{
			name:      "Invalid input path",
			inputPath: "testdata/invalid_input_path_input.txt",
			wantNil:   true,
		},
		{
			name:      "Invalid delete path",
			inputPath: "testdata/invalid_delete_path_input.txt",
			wantNil:   true,
		},
		{
			name:      "Short delete path",
			inputPath: "testdata/short_delete_path_input.txt",
			wantErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inputSR, err := ftutilities.LoadSubscribeResponse(test.inputPath)
			if err != nil {
				t.Fatalf("Failed to load input message: %v", err)
			}
			ft := New()
			gotSR, err := ft.Translate(inputSR)
			if gotNil, gotErr := gotSR == nil, err != nil; gotNil || gotErr {
				switch {
				case gotErr != test.wantErr:
					t.Fatalf("Unexpected error result returned from translate() = %v, want error %t", err, test.wantErr)
				case err != nil:
					return
				case gotNil != test.wantNil:
					t.Fatalf("Unexpected nil result returned from translate() = %t, want nil %t", gotNil, test.wantNil)
				default:
					return
				}
			}

			if test.wantNil || test.wantErr {
				return
			}

			wantSR, err := ftutilities.LoadSubscribeResponse(test.wantOutputPath)
			if err != nil {
				t.Fatalf("Failed to load want message: %v", err)
			}
			if diff := cmp.Diff(wantSR, gotSR, protocmp.Transform()); diff != "" {
				t.Fatalf("Unexpected diff from translate (-want +got): %s", diff)
			}
		})
	}
}
