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

package aristainterface

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"github.com/openconfig/functional-translators/ftutilities"
)

func TestStatusTranslate(t *testing.T) {
	tests := []struct {
		name           string
		inputPath      string
		wantOutputPath string
		expectError    bool
	}{
		{
			name:           "AdminStatusNative",
			inputPath:      "testdata/interface_status_admin_native_input.txt",
			wantOutputPath: "testdata/interface_status_admin_native_output.txt",
		},
		{
			name:           "AdminStatusNativeShutdown",
			inputPath:      "testdata/interface_status_admin_native_shutdown_input.txt",
			wantOutputPath: "testdata/interface_status_admin_native_shutdown_output.txt",
		},
		{
			name:           "OperStatusNative",
			inputPath:      "testdata/interface_status_oper_native_input.txt",
			wantOutputPath: "testdata/interface_status_oper_native_output.txt",
		},
		{
			name:           "AdminStatusOCPassthrough",
			inputPath:      "testdata/interface_status_admin_oc_input.txt",
			wantOutputPath: "testdata/interface_status_admin_oc_output.txt",
		},
		{
			name:           "NativeStatusDelete",
			inputPath:      "testdata/interface_status_delete_native_input.txt",
			wantOutputPath: "testdata/interface_status_delete_native_output.txt",
		},
		{
			name:           "OpenConfigStatusDelete",
			inputPath:      "testdata/interface_status_delete_oc_input.txt",
			wantOutputPath: "testdata/interface_status_delete_oc_output.txt",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ft := NewStatusFT()
			inputSR, err := ftutilities.LoadSubscribeResponse(test.inputPath)
			if err != nil {
				t.Fatalf("Failed to load input message: %v", err)
			}
			gotSR, gotErr := ft.Translate(inputSR)
			errorMatchesExpectation := (gotErr != nil) == test.expectError
			if !errorMatchesExpectation {
				t.Fatalf("Unexpected error result returned from translate() = %v, want error %t", gotErr, test.expectError)
			}
			if !test.expectError {
				if test.wantOutputPath == "" {
					if gotSR != nil {
						t.Fatalf("unexpected SubscribeResponse from translate() = %v, want nil", gotSR)
					}
					return
				}
				wantSR, err := ftutilities.LoadSubscribeResponse(test.wantOutputPath)
				if err != nil {
					t.Fatalf("Failed to load want message: %v", err)
				}
				if diff := cmp.Diff(wantSR, gotSR, protocmp.Transform()); diff != "" {
					t.Fatalf("unexpected diff from translate():\n%s", diff)
				}
			}
		})
	}
}
