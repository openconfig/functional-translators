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

package ciscoxrtransceiver

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

func TestTranslate(t *testing.T) {
	tests := []struct {
		name    string
		input   *gnmipb.SubscribeResponse
		want    *gnmipb.SubscribeResponse
		wantErr bool
	}{
		{
			name: "success, no optics info length != 4, don't report",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "Cisco-IOS-XR-controller-optics-oper",
							Elem: []*gnmipb.PathElem{
								{Name: "optics-oper"},
								{Name: "optics-ports"},
								{Name: "optics-port", Key: map[string]string{"name": "Optics0/0/0/0/1"}},
								{Name: "optics-info"},
							},
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "lane-index"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 0}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "transmit-power"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_IntVal{IntVal: 147}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "receive-power"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_IntVal{IntVal: -17}},
							},
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "success, no optics info length 4",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "Cisco-IOS-XR-controller-optics-oper",
							Elem: []*gnmipb.PathElem{
								{Name: "optics-oper"},
								{Name: "optics-ports"},
								{Name: "optics-port", Key: map[string]string{"name": "Optics0/0/0/0"}},
								{Name: "optics-info"},
							},
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "lane-index"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 0}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "transmit-power"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_IntVal{IntVal: 147}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "receive-power"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_IntVal{IntVal: -17}},
							},
						},
					},
				},
			},
			want: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "openconfig",
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "Optics0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "physical-channels"},
										{Name: "channel", Key: map[string]string{"index": "0"}},
										{Name: "state"},
										{Name: "index"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 0}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "Optics0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "physical-channels"},
										{Name: "channel", Key: map[string]string{"index": "0"}},
										{Name: "state"},
										{Name: "output-power"},
										{Name: "instant"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_DoubleVal{DoubleVal: 1.47}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "Optics0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "physical-channels"},
										{Name: "channel", Key: map[string]string{"index": "0"}},
										{Name: "state"},
										{Name: "input-power"},
										{Name: "instant"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_DoubleVal{DoubleVal: -0.17}},
							},
						},
					},
				},
			},
		},
		{
			name: "success, modified port name",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "Cisco-IOS-XR-controller-optics-oper",
							Elem: []*gnmipb.PathElem{
								{Name: "optics-oper"},
								{Name: "optics-ports"},
								{Name: "optics-port", Key: map[string]string{"name": "Optics0/0/0/0"}},
								{Name: "optics-info"},
							},
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "derived-optics-type"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "400G"}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "lane-index"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 0}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "transmit-power"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_IntVal{IntVal: 147}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "receive-power"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_IntVal{IntVal: -17}},
							},
						},
					},
				},
			},
			want: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "openconfig",
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "FourHundredGigE0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "physical-channels"},
										{Name: "channel", Key: map[string]string{"index": "0"}},
										{Name: "state"},
										{Name: "index"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 0}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "FourHundredGigE0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "physical-channels"},
										{Name: "channel", Key: map[string]string{"index": "0"}},
										{Name: "state"},
										{Name: "output-power"},
										{Name: "instant"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_DoubleVal{DoubleVal: 1.47}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "FourHundredGigE0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "physical-channels"},
										{Name: "channel", Key: map[string]string{"index": "0"}},
										{Name: "state"},
										{Name: "input-power"},
										{Name: "instant"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_DoubleVal{DoubleVal: -0.17}},
							},
						},
					},
				},
			},
		},
		{
			name: "success, modified port name HundredGigE",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "Cisco-IOS-XR-controller-optics-oper",
							Elem: []*gnmipb.PathElem{
								{Name: "optics-oper"},
								{Name: "optics-ports"},
								{Name: "optics-port", Key: map[string]string{"name": "Optics0/0/0/0"}},
								{Name: "optics-info"},
							},
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "derived-optics-type"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "100G"}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "lane-index"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 1}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "transmit-power"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_IntVal{IntVal: 123}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "receive-power"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_IntVal{IntVal: -45}},
							},
						},
					},
				},
			},
			want: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "openconfig",
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "HundredGigE0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "physical-channels"},
										{Name: "channel", Key: map[string]string{"index": "1"}},
										{Name: "state"},
										{Name: "index"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 1}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "HundredGigE0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "physical-channels"},
										{Name: "channel", Key: map[string]string{"index": "1"}},
										{Name: "state"},
										{Name: "output-power"},
										{Name: "instant"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_DoubleVal{DoubleVal: 1.23}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "HundredGigE0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "physical-channels"},
										{Name: "channel", Key: map[string]string{"index": "1"}},
										{Name: "state"},
										{Name: "input-power"},
										{Name: "instant"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_DoubleVal{DoubleVal: -0.45}},
							},
						},
					},
				},
			},
		},
		{
			name: "success, modified port name FourHundredGigE (4x100G)",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "Cisco-IOS-XR-controller-optics-oper",
							Elem: []*gnmipb.PathElem{
								{Name: "optics-oper"},
								{Name: "optics-ports"},
								{Name: "optics-port", Key: map[string]string{"name": "Optics0/0/0/0"}},
								{Name: "optics-info"},
							},
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "derived-optics-type"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "4x100G_LR4"}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "lane-index"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 2}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "transmit-power"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_IntVal{IntVal: 100}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "receive-power"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_IntVal{IntVal: -20}},
							},
						},
					},
				},
			},
			want: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "openconfig",
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "FourHundredGigE0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "physical-channels"},
										{Name: "channel", Key: map[string]string{"index": "2"}},
										{Name: "state"},
										{Name: "index"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 2}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "FourHundredGigE0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "physical-channels"},
										{Name: "channel", Key: map[string]string{"index": "2"}},
										{Name: "state"},
										{Name: "output-power"},
										{Name: "instant"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_DoubleVal{DoubleVal: 1.00}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "FourHundredGigE0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "physical-channels"},
										{Name: "channel", Key: map[string]string{"index": "2"}},
										{Name: "state"},
										{Name: "input-power"},
										{Name: "instant"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_DoubleVal{DoubleVal: -0.20}},
							},
						},
					},
				},
			},
		},
		{
			name: "success, modified port name FortyGigE (4x10G)",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "Cisco-IOS-XR-controller-optics-oper",
							Elem: []*gnmipb.PathElem{
								{Name: "optics-oper"},
								{Name: "optics-ports"},
								{Name: "optics-port", Key: map[string]string{"name": "Optics0/0/0/0"}},
								{Name: "optics-info"},
							},
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "derived-optics-type"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "4x10G_LR"}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "lane-index"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 3}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "transmit-power"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_IntVal{IntVal: 50}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "receive-power"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_IntVal{IntVal: -10}},
							},
						},
					},
				},
			},
			want: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "openconfig",
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "FortyGigE0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "physical-channels"},
										{Name: "channel", Key: map[string]string{"index": "3"}},
										{Name: "state"},
										{Name: "index"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 3}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "FortyGigE0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "physical-channels"},
										{Name: "channel", Key: map[string]string{"index": "3"}},
										{Name: "state"},
										{Name: "output-power"},
										{Name: "instant"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_DoubleVal{DoubleVal: 0.50}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "FortyGigE0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "physical-channels"},
										{Name: "channel", Key: map[string]string{"index": "3"}},
										{Name: "state"},
										{Name: "input-power"},
										{Name: "instant"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_DoubleVal{DoubleVal: -0.10}},
							},
						},
					},
				},
			},
		},
		{
			name: "success, modified port name TenGigE",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "Cisco-IOS-XR-controller-optics-oper",
							Elem: []*gnmipb.PathElem{
								{Name: "optics-oper"},
								{Name: "optics-ports"},
								{Name: "optics-port", Key: map[string]string{"name": "Optics0/0/0/0"}},
								{Name: "optics-info"},
							},
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "derived-optics-type"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "10G_SR"}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "lane-index"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 0}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "transmit-power"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_IntVal{IntVal: 10}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "receive-power"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_IntVal{IntVal: -5}},
							},
						},
					},
				},
			},
			want: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "openconfig",
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "TenGigE0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "physical-channels"},
										{Name: "channel", Key: map[string]string{"index": "0"}},
										{Name: "state"},
										{Name: "index"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 0}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "TenGigE0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "physical-channels"},
										{Name: "channel", Key: map[string]string{"index": "0"}},
										{Name: "state"},
										{Name: "output-power"},
										{Name: "instant"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_DoubleVal{DoubleVal: 0.10}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "TenGigE0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "physical-channels"},
										{Name: "channel", Key: map[string]string{"index": "0"}},
										{Name: "state"},
										{Name: "input-power"},
										{Name: "instant"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_DoubleVal{DoubleVal: -0.05}},
							},
						},
					},
				},
			},
		},
		{
			name: "success, modified port name TwoHundredGigE (2x100G)",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "Cisco-IOS-XR-controller-optics-oper",
							Elem: []*gnmipb.PathElem{
								{Name: "optics-oper"},
								{Name: "optics-ports"},
								{Name: "optics-port", Key: map[string]string{"name": "Optics0/0/0/0"}},
								{Name: "optics-info"},
							},
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "derived-optics-type"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "2x100G_LR4"}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "lane-index"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 2}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "transmit-power"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_IntVal{IntVal: 100}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "receive-power"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_IntVal{IntVal: -20}},
							},
						},
					},
				},
			},
			want: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "openconfig",
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "TwoHundredGigE0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "physical-channels"},
										{Name: "channel", Key: map[string]string{"index": "2"}},
										{Name: "state"},
										{Name: "index"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 2}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "TwoHundredGigE0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "physical-channels"},
										{Name: "channel", Key: map[string]string{"index": "2"}},
										{Name: "state"},
										{Name: "output-power"},
										{Name: "instant"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_DoubleVal{DoubleVal: 1.00}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "TwoHundredGigE0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "physical-channels"},
										{Name: "channel", Key: map[string]string{"index": "2"}},
										{Name: "state"},
										{Name: "input-power"},
										{Name: "instant"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_DoubleVal{DoubleVal: -0.20}},
							},
						},
					},
				},
			},
		},
		{
			name: "success, modified port name FortyGigE (40G)",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "Cisco-IOS-XR-controller-optics-oper",
							Elem: []*gnmipb.PathElem{
								{Name: "optics-oper"},
								{Name: "optics-ports"},
								{Name: "optics-port", Key: map[string]string{"name": "Optics0/0/0/0"}},
								{Name: "optics-info"},
							},
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "derived-optics-type"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "40G_LR4"}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "lane-index"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 1}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "transmit-power"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_IntVal{IntVal: 200}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "receive-power"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_IntVal{IntVal: -30}},
							},
						},
					},
				},
			},
			want: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "openconfig",
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "FortyGigE0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "physical-channels"},
										{Name: "channel", Key: map[string]string{"index": "1"}},
										{Name: "state"},
										{Name: "index"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 1}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "FortyGigE0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "physical-channels"},
										{Name: "channel", Key: map[string]string{"index": "1"}},
										{Name: "state"},
										{Name: "output-power"},
										{Name: "instant"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_DoubleVal{DoubleVal: 2.00}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "FortyGigE0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "physical-channels"},
										{Name: "channel", Key: map[string]string{"index": "1"}},
										{Name: "state"},
										{Name: "input-power"},
										{Name: "instant"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_DoubleVal{DoubleVal: -0.30}},
							},
						},
					},
				},
			},
		},
		{
			name: "no updates",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "Cisco-IOS-XR-controller-optics-oper",
							Elem: []*gnmipb.PathElem{
								{Name: "optics-oper"},
								{Name: "optics-ports"},
								{Name: "optics-port", Key: map[string]string{"name": "Optics0/0/0/0"}},
								{Name: "optics-info"},
							},
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{},
					},
				},
			},
			want: nil,
		},
		{
			name: "unexpected path",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183681554965,
						Prefix: &gnmipb.Path{
							Origin: "meta",
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "connectedAddress"},
									},
								},
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_StringVal{StringVal: "[2607:f8b0:8092:c4::12]:57400"},
								},
							},
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "unexpected value type",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "Cisco-IOS-XR-controller-optics-oper",
							Elem: []*gnmipb.PathElem{
								{Name: "optics-oper"},
								{Name: "optics-ports"},
								{Name: "optics-port", Key: map[string]string{"name": "Optics0/0/0/0/1"}},
								{Name: "optics-info"},
							},
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "lane-index"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 0}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "transmit-power"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_DoubleVal{DoubleVal: 1.47}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "receive-power"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "-17"}},
							},
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "success, laser bias current",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "Cisco-IOS-XR-controller-optics-oper",
							Elem: []*gnmipb.PathElem{
								{Name: "optics-oper"},
								{Name: "optics-ports"},
								{Name: "optics-port", Key: map[string]string{"name": "Optics0/0/0/0"}},
								{Name: "optics-info"},
							},
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "lane-index"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 0}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lane-data"},
										{Name: "laser-bias-current-milli-amps"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 12345}},
							},
						},
					},
				},
			},
			want: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "openconfig",
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "Optics0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "physical-channels"},
										{Name: "channel", Key: map[string]string{"index": "0"}},
										{Name: "state"},
										{Name: "index"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 0}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "Optics0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "physical-channels"},
										{Name: "channel", Key: map[string]string{"index": "0"}},
										{Name: "state"},
										{Name: "laser-bias-current"},
										{Name: "instant"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 12345}},
							},
						},
					},
				},
			},
		},
		{
			name: "success, form factor",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "Cisco-IOS-XR-controller-optics-oper",
							Elem: []*gnmipb.PathElem{
								{Name: "optics-oper"},
								{Name: "optics-ports"},
								{Name: "optics-port", Key: map[string]string{"name": "Optics0/0/0/0"}},
								{Name: "optics-info"},
							},
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "form-factor"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "qsfp"}},
							},
						},
					},
				},
			},
			want: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "openconfig",
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "Optics0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "state"},
										{Name: "form-factor"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "qsfp"}},
							},
						},
					},
				},
			},
		},
		{
			name: "success, vendor name",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "Cisco-IOS-XR-controller-optics-oper",
							Elem: []*gnmipb.PathElem{
								{Name: "optics-oper"},
								{Name: "optics-ports"},
								{Name: "optics-port", Key: map[string]string{"name": "Optics0/0/0/0"}},
								{Name: "optics-info"},
								{Name: "transceiver-info"},
							},
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "vendor-name"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "vendor-a"}},
							},
						},
					},
				},
			},
			want: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "openconfig",
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "Optics0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "state"},
										{Name: "vendor"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "vendor-a"}},
							},
						},
					},
				},
			},
		},
		{
			name: "success, vendor part",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "Cisco-IOS-XR-controller-optics-oper",
							Elem: []*gnmipb.PathElem{
								{Name: "optics-oper"},
								{Name: "optics-ports"},
								{Name: "optics-port", Key: map[string]string{"name": "Optics0/0/0/0"}},
								{Name: "optics-info"},
								{Name: "transceiver-info"},
							},
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "optics-vendor-part"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "part-123"}},
							},
						},
					},
				},
			},
			want: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "openconfig",
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "Optics0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "state"},
										{Name: "vendor-part"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "part-123"}},
							},
						},
					},
				},
			},
		},
		{
			name: "success, vendor rev",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "Cisco-IOS-XR-controller-optics-oper",
							Elem: []*gnmipb.PathElem{
								{Name: "optics-oper"},
								{Name: "optics-ports"},
								{Name: "optics-port", Key: map[string]string{"name": "Optics0/0/0/0"}},
								{Name: "optics-info"},
								{Name: "transceiver-info"},
							},
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "optics-vendor-rev"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "rev-b"}},
							},
						},
					},
				},
			},
			want: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1749043183927000000,
						Prefix: &gnmipb.Path{
							Origin: "openconfig",
							Target: "dx05.sql85-laarz",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "components"},
										{Name: "component", Key: map[string]string{"name": "Optics0/0/0/0"}},
										{Name: "transceiver"},
										{Name: "state"},
										{Name: "vendor-rev"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "rev-b"}},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ft := New()
			sr, err := ft.Translate(test.input)
			if (err != nil) != test.wantErr {
				t.Fatalf("Unexpected error result returned from translate() = %v, want error %t", err, test.wantErr)
			}
			if test.wantErr {
				return
			}
			if diff := cmp.Diff(test.want, sr, protocmp.Transform(), protocmp.SortRepeatedFields(&gnmipb.Notification{}, "update")); diff != "" {
				t.Fatalf("Unexpected diff from translate() = %v, want %v, diff(-want +got):\n%s", sr, test.want, diff)
			}
		})
	}
}
