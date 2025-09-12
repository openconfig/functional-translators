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

package ciscoxrvendordrops

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
								StringVal: "L3_ROUTE_LOOKUP_FAILED",
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
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "trap-id", Key: map[string]string{"trap-id": "0"}},
								{Name: "packet-accepted"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 1000000,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "trap-id", Key: map[string]string{"trap-id": "1"}},
								{Name: "trap-string"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "L3_NULL_ADJ(D*)",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "trap-id", Key: map[string]string{"trap-id": "1"}},
								{Name: "packet-dropped"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 33,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "trap-id", Key: map[string]string{"trap-id": "2"}},
								{Name: "trap-string"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "MPLS_TE_MIDPOINT_LDP_LABELS_MISS(D*)",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "trap-id", Key: map[string]string{"trap-id": "2"}},
								{Name: "packet-dropped"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 34,
							},
						},
					},
				},
			},
		},
	}
	successStatsSR := &gnmipb.SubscribeResponse{
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
						{Name: "Cisco-IOS-XR-ofa-npu-stats-oper:asic-statistics"},
						{Name: "asic-statistics-for-npu-ids"},
						{Name: "asic-statistics-for-npu-id", Key: map[string]string{"npu-id": "0"}},
					},
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "valid"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_BoolVal{
								BoolVal: true,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "rack-number"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 0,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "npu-statistics"},
								{Name: "block-info"},
								{Name: "block-name"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "block 1 Summary",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "npu-statistics"},
								{Name: "block-info"},
								{Name: "field-info"},
								{Name: "is-overflow"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_BoolVal{
								BoolVal: false,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "npu-statistics"},
								{Name: "block-info"},
								{Name: "field-info"},
								{Name: "field-name"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "IFGB_RX 0 partial drop",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "npu-statistics"},
								{Name: "block-info"},
								{Name: "field-info"},
								{Name: "field-value"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 100,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "npu-statistics"},
								{Name: "block-info"},
								{Name: "field-info"},
								{Name: "field-name"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "IFGB_RX 1 partial drop",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "npu-statistics"},
								{Name: "block-info"},
								{Name: "field-info"},
								{Name: "field-value"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 101,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "valid"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_BoolVal{
								BoolVal: true,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "rack-number"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 0,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "npu-statistics"},
								{Name: "block-info"},
								{Name: "block-name"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "block 2 Summary",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "npu-statistics"},
								{Name: "block-info"},
								{Name: "field-info"},
								{Name: "is-overflow"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_BoolVal{
								BoolVal: false,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "npu-statistics"},
								{Name: "block-info"},
								{Name: "field-info"},
								{Name: "field-name"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "IFGB_RX 0 partial drop",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "npu-statistics"},
								{Name: "block-info"},
								{Name: "field-info"},
								{Name: "field-value"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 102,
							},
						},
					},
				},
			},
		},
	}
	successStatsIgnoredSR := &gnmipb.SubscribeResponse{
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
						{Name: "Cisco-IOS-XR-ofa-npu-stats-oper:asic-statistics"},
						{Name: "asic-statistics-for-npu-ids"},
						{Name: "asic-statistics-for-npu-id", Key: map[string]string{"npu-id": "0"}},
					},
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "valid"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_BoolVal{
								BoolVal: true,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "rack-number"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 0,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "npu-statistics"},
								{Name: "block-info"},
								{Name: "block-name"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "block 1 Summary",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "npu-statistics"},
								{Name: "block-info"},
								{Name: "field-info"},
								{Name: "is-overflow"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_BoolVal{
								BoolVal: false,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "npu-statistics"},
								{Name: "block-info"},
								{Name: "field-info"},
								{Name: "field-name"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "IFGB_RX 0 partial drop",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "npu-statistics"},
								{Name: "block-info"},
								{Name: "field-info"},
								{Name: "field-value"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 100,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "npu-statistics"},
								{Name: "block-info"},
								{Name: "field-info"},
								{Name: "field-name"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "IFGB_RX 1 partial drop",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "npu-statistics"},
								{Name: "block-info"},
								{Name: "field-info"},
								{Name: "field-value"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 101,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "valid"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_BoolVal{
								BoolVal: true,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "rack-number"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 0,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "npu-statistics"},
								{Name: "block-info"},
								{Name: "block-name"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "block 2 Summary",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "npu-statistics"},
								{Name: "block-info"},
								{Name: "field-info"},
								{Name: "is-overflow"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_BoolVal{
								BoolVal: false,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "npu-statistics"},
								{Name: "block-info"},
								{Name: "field-info"},
								{Name: "field-name"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "IFGB_RX 0 partial drop",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "npu-statistics"},
								{Name: "block-info"},
								{Name: "field-info"},
								{Name: "field-value"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 102,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "npu-statistics"},
								{Name: "block-info"},
								{Name: "block-name"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "block 1 ",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "npu-statistics"},
								{Name: "block-info"},
								{Name: "field-info"},
								{Name: "is-overflow"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_BoolVal{
								BoolVal: false,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "npu-statistics"},
								{Name: "block-info"},
								{Name: "field-info"},
								{Name: "field-name"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "IFGB_RX 0 partial drop",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "npu-statistics"},
								{Name: "block-info"},
								{Name: "field-info"},
								{Name: "field-value"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 100,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "npu-statistics"},
								{Name: "block-info"},
								{Name: "field-info"},
								{Name: "field-name"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "IFGB_RX 1 partial drop",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "npu-statistics"},
								{Name: "block-info"},
								{Name: "field-info"},
								{Name: "field-value"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 110,
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
								{Name: "vendor"},
								{Name: "CiscoXR"},
								{Name: "spitfire"},
								{Name: "packet-processing"},
								{Name: "state"},
								{Name: "L3_ROUTE_LOOKUP_FAILED"},
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
								{Name: "drop"},
								{Name: "vendor"},
								{Name: "CiscoXR"},
								{Name: "spitfire"},
								{Name: "packet-processing"},
								{Name: "state"},
								{Name: "L3_NULL_ADJ(D*)"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 33,
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
								{Name: "drop"},
								{Name: "vendor"},
								{Name: "CiscoXR"},
								{Name: "spitfire"},
								{Name: "packet-processing"},
								{Name: "state"},
								{Name: "MPLS_TE_MIDPOINT_LDP_LABELS_MISS(D*)"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 34,
							},
						},
					},
				},
			},
		},
	}
	successStatsOutput := &gnmipb.SubscribeResponse{
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
								{Name: "component", Key: map[string]string{"name": "0/RP0/CPU0:0:block_1_Summary"}},
								{Name: "integrated-circuit"},
								{Name: "pipeline-counters"},
								{Name: "drop"},
								{Name: "vendor"},
								{Name: "CiscoXR"},
								{Name: "spitfire"},
								{Name: "adverse"},
								{Name: "state"},
								{Name: "IFGB_RX_0_partial_drop"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 100,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "components"},
								{Name: "component", Key: map[string]string{"name": "0/RP0/CPU0:0:block_1_Summary"}},
								{Name: "integrated-circuit"},
								{Name: "pipeline-counters"},
								{Name: "drop"},
								{Name: "vendor"},
								{Name: "CiscoXR"},
								{Name: "spitfire"},
								{Name: "adverse"},
								{Name: "state"},
								{Name: "IFGB_RX_1_partial_drop"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 101,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "components"},
								{Name: "component", Key: map[string]string{"name": "0/RP0/CPU0:0:block_2_Summary"}},
								{Name: "integrated-circuit"},
								{Name: "pipeline-counters"},
								{Name: "drop"},
								{Name: "vendor"},
								{Name: "CiscoXR"},
								{Name: "spitfire"},
								{Name: "adverse"},
								{Name: "state"},
								{Name: "IFGB_RX_0_partial_drop"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 102,
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
					Origin: "Cisco-IOS-XR-platforms-ofa-oper",
					Elem: []*gnmipb.PathElem{
						{Name: "ofa"},
						{Name: "stats"},
						{Name: "nodes"},
						{Name: "node", Key: map[string]string{"node-name": "0/RP0/CPU0"}},
						{Name: "Cisco-IOS-XR-ofa-npu-stats-oper:npu-numbers"},
						{Name: "npu-number", Key: map[string]string{"npu-id": "xyzaz"}},
						{Name: "display"},
						{Name: "trap-ids"},
					},
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "trap-id", Key: map[string]string{"trap-id": "xyx"}},
								{Name: "trap-string"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "L3_ROUTE_LOOKUP_FAILED",
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
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "trap-id", Key: map[string]string{"trap-id": "1"}},
								{Name: "trap-string"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "L3_NULL_ADJ(D*)",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "trap-id", Key: map[string]string{"trap-id": "1"}},
								{Name: "packet-dropped"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 33,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "trap-id", Key: map[string]string{"trap-id": "2"}},
								{Name: "trap-string"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "MPLS_TE_MIDPOINT_LDP_LABELS_MISS(D*)",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "trap-id", Key: map[string]string{"trap-id": "2"}},
								{Name: "packet-dropped"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 34,
							},
						},
					},
				},
			},
		},
	}
	notMatchedSR := &gnmipb.SubscribeResponse{
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
						{Name: "npu-numberssss", Key: map[string]string{"npu-id": "xyzaz"}},
						{Name: "display"},
						{Name: "trap-ids"},
					},
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "trap-id", Key: map[string]string{"trap-id": "xyx"}},
								{Name: "trap-string"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "L3_ROUTE_LOOKUP_FAILED",
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
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "trap-id", Key: map[string]string{"trap-id": "1"}},
								{Name: "trap-string"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "L3_NULL_ADJ(D*)",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "trap-id", Key: map[string]string{"trap-id": "1"}},
								{Name: "packet-dropped"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 33,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "trap-id", Key: map[string]string{"trap-id": "2"}},
								{Name: "trap-string"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "MPLS_TE_MIDPOINT_LDP_LABELS_MISS(D*)",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "trap-id", Key: map[string]string{"trap-id": "2"}},
								{Name: "packet-dropped"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 34,
							},
						},
					},
				},
			},
		},
	}
	notMatchedOutput := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 123,
				Prefix: &gnmipb.Path{
					Origin: "openconfig",
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
		{
			name:  "not matched path",
			input: notMatchedSR,
			want:  notMatchedOutput,
		},
		{
			name:  "success stats",
			input: successStatsSR,
			want:  successStatsOutput,
		},
		{
			name:  "success stats with ignored non summayblock",
			input: successStatsIgnoredSR,
			want:  successStatsOutput,
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
