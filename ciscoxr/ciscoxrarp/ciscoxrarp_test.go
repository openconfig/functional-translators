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

package ciscoxrarp

import (
	"sync"
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
	}{
		{
			name:           "success_arp",
			inputPath:      "testdata/arp_input.txt",
			wantOutputPath: "testdata/arp_output.txt",
		},
		{
			name:           "success_np",
			inputPath:      "testdata/nd_input.txt",
			wantOutputPath: "testdata/nd_output.txt",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inputSR, err := ftutilities.LoadSubscribeResponse(test.inputPath)
			if err != nil {
				t.Fatalf("Failed to load input message: %v", err)
			}
			gotSR, err := New().Translate(inputSR)
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

func TestConcurrentTranslate(t *testing.T) {
	tests := []struct {
		name           string
		inputPath      string
		wantOutputPath string
	}{
		{
			name:           "SuccessARP",
			inputPath:      "testdata/arp_input.txt",
			wantOutputPath: "testdata/arp_output.txt",
		},
		{
			name:           "SuccessND",
			inputPath:      "testdata/nd_input.txt",
			wantOutputPath: "testdata/nd_output.txt",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Do not call t.Parallel(): resetting package-level schema state requires sequential subtests.
			schemaOnce = sync.Once{}
			schema = nil
			schemaErr = nil

			wantSR, err := ftutilities.LoadSubscribeResponse(tc.wantOutputPath)
			if err != nil {
				t.Fatalf("LoadSubscribeResponse(%q) failed: %v", tc.wantOutputPath, err)
			}

			var wg sync.WaitGroup
			for range 16 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					inputSR, err := ftutilities.LoadSubscribeResponse(tc.inputPath)
					if err != nil {
						t.Errorf("LoadSubscribeResponse(%q) failed: %v", tc.inputPath, err)
						return
					}
					gotSR, err := New().Translate(inputSR)
					if err != nil {
						t.Errorf("Translate(%v) got error %v, want nil", inputSR, err)
						return
					}
					if diff := cmp.Diff(wantSR, gotSR, protocmp.Transform(), protocmp.SortRepeatedFields(&gnmipb.Notification{}, "update")); diff != "" {
						t.Errorf("Translate(%v) returned unexpected diff (-want +got):\n%s", inputSR, diff)
					}
				}()
			}
			wg.Wait()
		})
	}
}
