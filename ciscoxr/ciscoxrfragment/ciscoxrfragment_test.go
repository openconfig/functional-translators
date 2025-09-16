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

package ciscoxrfragment

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

func TestTranslate(t *testing.T) {
	successSR := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 123,
				Prefix: &gnmipb.Path{
					Origin: "Cisco-IOS-XR-platforms-ofa-oper",
					Elem: []*gnmipb.PathElem{
						{Name: "ofa"},
						{Name: "stats"},
						{Name: "nodes"},
						{Name: "node", Key: map[string]string{"node-name": "0/RP0/CPU0"}},
						{Name: "Cisco-IOS-XR-ofa-npu-stats-oper:npu-numbers"},
						{Name: "npu-number", Key: map[string]string{"npu-id": "0"}},
						{Name: "display"},
						{Name: "trap-ids"},
					},
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "trap-id", Key: map[string]string{"trap-id": "0"}},
								{Name: "trap-string"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "MTU_EXCEEDED",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "trap-id", Key: map[string]string{"trap-id": "0"}},
								{Name: "packet-accepted"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 32,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "trap-id", Key: map[string]string{"trap-id": "0"}},
								{Name: "packet-dropped"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 32,
							},
						},
					},
				},
			},
		},
	}
	successOutput := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 123,
				Prefix: &gnmipb.Path{
					Origin: "openconfig",
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "components"},
								{Name: "component", Key: map[string]string{"name": "0/RP0/CPU0:0"}},
								{Name: "integrated-circuit"},
								{Name: "pipeline-counters"},
								{Name: "drop"},
								{Name: "host-interface-block"},
								{Name: "state"},
								{Name: "fragment-punt"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 32,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "components"},
								{Name: "component", Key: map[string]string{"name": "0/RP0/CPU0:0"}},
								{Name: "integrated-circuit"},
								{Name: "pipeline-counters"},
								{Name: "packet"},
								{Name: "host-interface-block"},
								{Name: "state"},
								{Name: "fragment-punt-pkts"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 32,
							},
						},
					},
				},
			},
		},
	}
	invalidSR := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 123,
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "invalid"},
								{Name: "path"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_IntVal{
								IntVal: 123,
							},
						},
					},
				},
			},
		},
	}
	tests := []struct {
		name    string
		input   *gnmipb.SubscribeResponse
		want    *gnmipb.SubscribeResponse
		wantErr bool
	}{
		{
			name:  "success",
			input: successSR,
			want:  successOutput,
		},
		{
			name:    "invalid path",
			input:   invalidSR,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ft := New()
			sr, err := ft.Translate(test.input)
			errorMatchesExpectation := (err != nil) == test.wantErr
			if !errorMatchesExpectation {
				t.Fatalf("Unexpected error result returned from translate() = %v, want error %t", err, test.wantErr)
			}
			if !test.wantErr {
				if diff := cmp.Diff(test.want, sr, protocmp.Transform(), protocmp.SortRepeatedFields(&gnmipb.Notification{}, "update")); diff != "" {
					t.Fatalf("Unexpected diff from translate() = %v, want %v, diff(-want +got):\n%s", sr, test.want, diff)
				}
			}
		})
	}
}
