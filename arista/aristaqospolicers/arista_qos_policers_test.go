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

package aristaqospolicers

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
		wantNil        bool
		wantErr        bool
	}{
		{
			name:           "Logical interface first member update",
			inputPath:      "testdata/logical_interface_member1_input.txt",
			wantOutputPath: "testdata/logical_interface_member1_output.txt",
		},
		{
			name:           "Logical interface second member update calculates aggregate sum",
			inputPath:      "testdata/logical_interface_member2_input.txt",
			wantOutputPath: "testdata/logical_interface_member2_output.txt",
		},
		{
			name:           "Logical interface member deletion reduces sum",
			inputPath:      "testdata/logical_interface_delete1_input.txt",
			wantOutputPath: "testdata/logical_interface_delete1_output.txt",
		},
		{
			name:           "Logical interface member deletion removes OC path",
			inputPath:      "testdata/logical_interface_delete2_input.txt",
			wantOutputPath: "testdata/logical_interface_delete2_output.txt",
		},
		{
			name:      "Skip update with empty path elements",
			inputPath: "testdata/empty_update_path_input.txt",
			wantNil:   true,
		},
		{
			name:      "Skip invalid key without underscore (len != 2)",
			inputPath: "testdata/invalid_key_no_underscore_input.txt",
			wantNil:   true,
		},
		{
			name:      "Skip physical interface with empty attachment point (parts[1] == \"\")",
			inputPath: "testdata/physical_interface_skip_input.txt",
			wantNil:   true,
		},
	}

	// Initialize the FT once so state is maintained chronologically across the updates
	ft := New()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inputSR, err := ftutilities.LoadSubscribeResponse(test.inputPath)
			if err != nil {
				t.Fatalf("Failed to load input message: %v", err)
			}

			gotSR, err := ft.Translate(inputSR)
			if (err != nil) != test.wantErr {
				t.Fatalf("Unexpected error result returned from Translate() = %v, want error %t", err, test.wantErr)
			}
			if err != nil {
				return
			}
			if (gotSR == nil) != test.wantNil {
				t.Fatalf("Unexpected nil result returned from Translate() = %t, want nil %t", gotSR == nil, test.wantNil)
			}
			if gotSR == nil {
				return
			}

			wantSR, err := ftutilities.LoadSubscribeResponse(test.wantOutputPath)
			if err != nil {
				t.Fatalf("Failed to load want message: %v", err)
			}
			if diff := cmp.Diff(wantSR, gotSR, protocmp.Transform()); diff != "" {
				t.Fatalf("Unexpected diff from Translate() (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCacheCleanupOnLastMemberDeletion(t *testing.T) {
	ft := New()
	target := "target-cache-test"
	attachmentPoint := "Port-Channel99"
	physIntf := "Ethernet99/1"

	// Send an update to create the attachment point in the cache
	addInput := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 1000,
				Prefix:    &gnmipb.Path{Origin: "eos_native", Target: target},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "Sysdb"}, {Name: "qos"}, {Name: "status"}, {Name: "policingStatus"},
								{Name: "intfPolicer"}, {Name: physIntf + "_" + attachmentPoint}, {Name: "counter"}, {Name: "byteDropCount"},
							},
						},
						Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 100}},
					},
				},
			},
		},
	}
	_, _ = ft.Translate(addInput)

	// Verify the attachment point was successfully created in the cache
	targetInfo := ftutilities.PolicerAggMap.CreateOrUpdateTargetPolicerInfo(target)
	if _, ok := targetInfo.AttachmentPointInfo(attachmentPoint); !ok {
		t.Fatalf("Setup failed: Expected %s to exist in cache", attachmentPoint)
	}

	// Send a delete notification to remove the only physical member
	delInput := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 2000,
				Prefix:    &gnmipb.Path{Origin: "eos_native", Target: target},
				Delete: []*gnmipb.Path{
					{
						Elem: []*gnmipb.PathElem{
							{Name: "Sysdb"}, {Name: "qos"}, {Name: "status"}, {Name: "policingStatus"},
							{Name: "intfPolicer"}, {Name: physIntf + "_" + attachmentPoint},
						},
					},
				},
			},
		},
	}
	_, _ = ft.Translate(delInput)

	// Verify the attachment point was actually deleted from the cache
	if _, ok := targetInfo.AttachmentPointInfo(attachmentPoint); ok {
		t.Errorf("Expected %s to be completely deleted from cache after its last member was removed to prevent memory leaks", attachmentPoint)
	}
}
