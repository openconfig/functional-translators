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

package ciscoxrfabric

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

func TestTranslate(t *testing.T) {
	successFabricPlaneSR := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 123,
				Prefix: &gnmipb.Path{
					Origin: "Cisco-IOS-XR-fabric-plane-health-oper",
					Target: "some_target",
					Elem: []*gnmipb.PathElem{
						{Name: "fabric"},
						{Name: "fabric-plane-ids"},
						{Name: "fabric-plane-id", Key: map[string]string{"fabric-plane-key": "1"}},
						{Name: "fabric-plane-stats"},
					},
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "asic-internal-drops"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 1,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "mcast-lost-cells"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 2,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "rx-pe-cells"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 3,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "rx-uce-cells"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 4,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "ucast-lost-cells"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 5,
							},
						},
					},
				},
			},
		},
	}
	successSwitchStatsSR := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 123,
				Prefix: &gnmipb.Path{
					Origin: "Cisco-IOS-XR-switch-oper",
					Elem: []*gnmipb.PathElem{
						{Name: "show-switch"},
						{Name: "statistics"},
						{Name: "statistics-detail"},
						{Name: "statistics-detail-instances"},
						{Name: "statistics-detail-instance", Key: map[string]string{"node-id": "0/0/RP0"}},
						{Name: "statistics-detail-port-numbers"},
						{Name: "statistics-detail-port-number", Key: map[string]string{"port": "1"}},
					},
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "ethsw-detailed-stat-info"},
								{Name: "rx-bad-crc"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 1,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "ethsw-detailed-stat-info"},
								{Name: "rx-errors"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 2,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "ethsw-detailed-stat-info"},
								{Name: "tx-fifo-unrun"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 3,
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
				Prefix: &gnmipb.Path{
					Origin: "Cisco-IOS-XR-fabric-plane-health-oper",
					Elem: []*gnmipb.PathElem{
						{Name: "fabric"},
						{Name: "fabric-plane-ids"},
						{Name: "fabric-plane-id", Key: map[string]string{"fabric-plane-key": "1"}},
					},
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "asic-internal-drops"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 1,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "mcast-lost-cells"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 2,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "rx-pe-cells"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 3,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "rx-uce-cells"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 4,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "ucast-lost-cells"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 5,
							},
						},
					},
				},
			},
		},
	}
	successFabricPlaneOutput := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 123,
				Prefix: &gnmipb.Path{
					Origin: "openconfig",
					Target: "some_target",
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "components"},
								{Name: "component", Key: map[string]string{"name": "1"}},
								{Name: "integrated-circuit"},
								{Name: "pipeline-counters"},
								{Name: "errors"},
								{Name: "fabric-block"},
								{Name: "fabric-block-error", Key: map[string]string{"name": "asic-internal-drops"}},
								{Name: "state"},
								{Name: "count"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 1,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "components"},
								{Name: "component", Key: map[string]string{"name": "1"}},
								{Name: "integrated-circuit"},
								{Name: "pipeline-counters"},
								{Name: "errors"},
								{Name: "fabric-block"},
								{Name: "fabric-block-error", Key: map[string]string{"name": "multicast-lost-cells"}},
								{Name: "state"},
								{Name: "count"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 2,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "components"},
								{Name: "component", Key: map[string]string{"name": "1"}},
								{Name: "integrated-circuit"},
								{Name: "pipeline-counters"},
								{Name: "errors"},
								{Name: "fabric-block"},
								{Name: "fabric-block-error", Key: map[string]string{"name": "parity-error-cells"}},
								{Name: "state"},
								{Name: "count"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 3,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "components"},
								{Name: "component", Key: map[string]string{"name": "1"}},
								{Name: "integrated-circuit"},
								{Name: "pipeline-counters"},
								{Name: "errors"},
								{Name: "fabric-block"},
								{Name: "fabric-block-error", Key: map[string]string{"name": "uncorrectable-error-cells"}},
								{Name: "state"},
								{Name: "count"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 4,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "components"},
								{Name: "component", Key: map[string]string{"name": "1"}},
								{Name: "integrated-circuit"},
								{Name: "pipeline-counters"},
								{Name: "errors"},
								{Name: "fabric-block"},
								{Name: "fabric-block-error", Key: map[string]string{"name": "unicast-lost-cells"}},
								{Name: "state"},
								{Name: "count"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 5,
							},
						},
					},
				},
			},
		},
	}
	successSwitchStatsOutput := &gnmipb.SubscribeResponse{
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
								{Name: "component", Key: map[string]string{"name": "0/0/RP0:1"}},
								{Name: "integrated-circuit"},
								{Name: "pipeline-counters"},
								{Name: "errors"},
								{Name: "fabric-block"},
								{Name: "fabric-block-error", Key: map[string]string{"name": "rx-bad-crc"}},
								{Name: "state"},
								{Name: "count"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 1,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "components"},
								{Name: "component", Key: map[string]string{"name": "0/0/RP0:1"}},
								{Name: "integrated-circuit"},
								{Name: "pipeline-counters"},
								{Name: "errors"},
								{Name: "fabric-block"},
								{Name: "fabric-block-error", Key: map[string]string{"name": "rx-errors"}},
								{Name: "state"},
								{Name: "count"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 2,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "components"},
								{Name: "component", Key: map[string]string{"name": "0/0/RP0:1"}},
								{Name: "integrated-circuit"},
								{Name: "pipeline-counters"},
								{Name: "errors"},
								{Name: "fabric-block"},
								{Name: "fabric-block-error", Key: map[string]string{"name": "tx-fifo-unrun"}},
								{Name: "state"},
								{Name: "count"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 3,
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
			name:  "successSwitchStats",
			input: successSwitchStatsSR,
			want:  successSwitchStatsOutput,
		},
		{
			name:  "successFabricPlane",
			input: successFabricPlaneSR,
			want:  successFabricPlaneOutput,
		},
		{
			name:    "invalid SR",
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
