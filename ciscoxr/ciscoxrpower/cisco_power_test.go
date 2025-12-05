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

package ciscoxrpower

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
		expectNil      bool
	}{
		{
			name:           "successNative_to_OC",
			inputPath:      "testdata/success_native_to_oc_input.txt",
			wantOutputPath: "testdata/success_native_to_oc_output.txt",
		},
		{
			name:      "successOC_filter",
			inputPath: "testdata/success_oc_filter_input.txt",
			expectNil: true,
		},
		{
			name:           "successOC_no_filter",
			inputPath:      "testdata/success_oc_no_filter_input.txt",
			wantOutputPath: "testdata/success_oc_no_filter_output.txt",
		},
		{
			name:      "no_translation_non_8808",
			inputPath: "testdata/no_translation_non_8808_input.txt",
			expectNil: true,
		},
		{
			name:        "failureNative_to_OC",
			inputPath:   "testdata/failure_native_to_oc_input.txt",
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
			if test.expectError && !test.expectNil {
				errorMatchesExpectation := (err != nil) == test.expectError
				if !errorMatchesExpectation {
					t.Fatalf("Unexpected error result returned from translate() = %v, want error %t", err, test.expectError)
				}

			}
			if test.expectNil && !test.expectError {
				nilMatchesExpectation := (gotSR == nil) == test.expectNil
				if !nilMatchesExpectation {
					t.Fatalf("Unexpected nil result returned from translate() = %v, want nil %t", gotSR, test.expectNil)
				}
			}
			if !test.expectError && !test.expectNil {
				wantSR, err := ftutilities.LoadSubscribeResponse(test.wantOutputPath)
				if err != nil {
					t.Fatalf("Failed to load want message: %v", err)
				}
				if diff := cmp.Diff(wantSR, gotSR, protocmp.Transform(), protocmp.SortRepeatedFields(&gnmipb.Notification{}, "update")); diff != "" {
					t.Fatalf("translate() returned unexpected diff (-want +got):\n%s", diff)
				}
			}
		})
	}

}
