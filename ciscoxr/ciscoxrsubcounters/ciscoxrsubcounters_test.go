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

package ciscoxrsubcounters

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
	}{
		{
			name:           "ipv4",
			inputPath:      "testdata/ipv4_input.txt",
			wantOutputPath: "testdata/ipv4_output.txt",
		},
		{
			name:           "infra",
			inputPath:      "testdata/infra_input.txt",
			wantOutputPath: "testdata/infra_output.txt",
		},
		{
			name:           "ipv6",
			inputPath:      "testdata/ipv6_input.txt",
			wantOutputPath: "testdata/ipv6_output.txt",
		},
		{
			name:      "ipv4_nil_vrf",
			inputPath: "testdata/ipv4_nil_vrf_input.txt",
		},
		{
			name:      "ipv4_nil_detail",
			inputPath: "testdata/ipv4_nil_detail_input.txt",
		},
		{
			name:      "ipv4_nil_node",
			inputPath: "testdata/ipv4_nil_node_input.txt",
		},
		{
			name:      "infra_nil_protocol",
			inputPath: "testdata/infra_nil_protocol_input.txt",
		},
		{
			name:      "nil_notification",
			inputPath: "testdata/nil_notification_input.txt",
		},
		{
			name:      "ipv4_nil_interface_data",
			inputPath: "testdata/ipv4_nil_interface_data_input.txt",
		},
		{
			name:      "ipv4_nil_vrfs",
			inputPath: "testdata/ipv4_nil_vrfs_input.txt",
		},
		{
			name:      "ipv4_nil_details",
			inputPath: "testdata/ipv4_nil_details_input.txt",
		},
		{
			name:      "infra_multiast_counters",
			inputPath: "testdata/infra_multicast_counters_input.txt",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inputSR, err := ftutilities.LoadSubscribeResponse(test.inputPath)
			if err != nil {
				t.Fatalf("Failed to load input message: %v", err)
			}
			gotSR, err := New().Translate(inputSR)
			if err != nil {
				t.Fatalf("translate() failed: %v", err)
			}
			var wantSR *gnmipb.SubscribeResponse
			if test.wantOutputPath != "" {
				var err error
				wantSR, err = ftutilities.LoadSubscribeResponse(test.wantOutputPath)
				if err != nil {
					t.Fatalf("Failed to load want message: %v", err)
				}
			}
			if diff := cmp.Diff(wantSR, gotSR, protocmp.Transform(), protocmp.SortRepeatedFields(&gnmipb.Notification{}, "update")); diff != "" {
				t.Errorf("translate() returned diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTranslateNil(t *testing.T) {
	if _, err := New().Translate(nil); err != nil {
		t.Errorf("Translate(nil) return unexpected error: %v", err)
	}
	if _, err := New().Translate(&gnmipb.SubscribeResponse{}); err != nil {
		t.Errorf("Translate(nil) return unexpected error: %v", err)
	}
	invalidSR := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Prefix: &gnmipb.Path{
					Origin: "Cisco-IOS-XR-ipv4-io-oper",
					Target: "dev",
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{
									Name: "invalid",
								},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "invalid",
							},
						},
					},
				},
			},
		},
	}
	if got, err := New().Translate(invalidSR); err == nil {
		t.Errorf("Translate(invalid) failed to return an error: %v. Got %v", err, got)
	}
}
