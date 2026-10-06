// Copyright 2026 Google LLC
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

package aristagirocounter

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

func TestTranslate(t *testing.T) {
	tests := []struct {
		desc      string
		in        *gnmipb.SubscribeResponse
		wantOut   *gnmipb.SubscribeResponse
		wantError bool
	}{
		{
			desc: "valid pkts update with real device telemetry (cx01.apa98)",
			in: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1726943800000000000,
						Prefix: &gnmipb.Path{
							Origin: "eos_native",
							Target: "cx01.apa98",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "Smash"},
										{Name: "flexCounters"},
										{Name: "counterTable"},
										{Name: "DecapGroup"},
										{Name: "4294967295"},
										{Name: "decapGroupCounter"},
										{Name: "10.61.61.148/30_udp"},
										{Name: "pkts"},
									},
								},
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_UintVal{
										UintVal: 96971532,
									},
								},
							},
						},
					},
				},
			},
			wantOut: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1726943800000000000,
						Prefix:    wantPrefix("cx01.apa98"),
						Update: []*gnmipb.Update{
							{
								Path: wantLeafPath("10.61.61.148/30", "1", "matched-pkts"),
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_UintVal{
										UintVal: 96971532,
									},
								},
							},
							{
								Path: wantDestAddrPath("10.61.61.148/30", "1", "ipv4"),
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_StringVal{
										StringVal: "10.61.61.148/30",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			desc: "valid octets update with real device telemetry (cx01.apa98)",
			in: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1726943800000000000,
						Prefix: &gnmipb.Path{
							Origin: "eos_native",
							Target: "cx01.apa98",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "Smash"},
										{Name: "flexCounters"},
										{Name: "counterTable"},
										{Name: "DecapGroup"},
										{Name: "4294967295"},
										{Name: "decapGroupCounter"},
										{Name: "10.61.61.148/30_udp"},
										{Name: "octets"},
									},
								},
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_UintVal{
										UintVal: 49649424384,
									},
								},
							},
						},
					},
				},
			},
			wantOut: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1726943800000000000,
						Prefix:    wantPrefix("cx01.apa98"),
						Update: []*gnmipb.Update{
							{
								Path: wantLeafPath("10.61.61.148/30", "1", "matched-octets"),
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_UintVal{
										UintVal: 49649424384,
									},
								},
							},
							{
								Path: wantDestAddrPath("10.61.61.148/30", "1", "ipv4"),
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_StringVal{
										StringVal: "10.61.61.148/30",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			desc: "valid pkts update with mocked IPv6 telemetry",
			// Note: This is mocked data, not real device telemetry.
			in: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1726943800000000000,
						Prefix: &gnmipb.Path{
							Origin: "eos_native",
							Target: "cx01.apa98",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "Smash"},
										{Name: "flexCounters"},
										{Name: "counterTable"},
										{Name: "DecapGroup"},
										{Name: "4294967295"},
										{Name: "decapGroupCounter"},
										{Name: "2001:db8::1/128_udp"},
										{Name: "pkts"},
									},
								},
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_UintVal{
										UintVal: 12345,
									},
								},
							},
						},
					},
				},
			},
			wantOut: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1726943800000000000,
						Prefix:    wantPrefix("cx01.apa98"),
						Update: []*gnmipb.Update{
							{
								Path: wantLeafPath("2001:db8::1/128", "1", "matched-pkts"),
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_UintVal{
										UintVal: 12345,
									},
								},
							},
							{
								Path: wantDestAddrPath("2001:db8::1/128", "1", "ipv6"),
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_StringVal{
										StringVal: "2001:db8::1/128",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			desc: "valid octets update with mocked IPv6 telemetry",
			// Note: This is mocked data, not real device telemetry.
			in: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1726943800000000000,
						Prefix: &gnmipb.Path{
							Origin: "eos_native",
							Target: "cx01.apa98",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "Smash"},
										{Name: "flexCounters"},
										{Name: "counterTable"},
										{Name: "DecapGroup"},
										{Name: "4294967295"},
										{Name: "decapGroupCounter"},
										{Name: "2001:db8::1/128_udp"},
										{Name: "octets"},
									},
								},
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_UintVal{
										UintVal: 54321,
									},
								},
							},
						},
					},
				},
			},
			wantOut: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1726943800000000000,
						Prefix:    wantPrefix("cx01.apa98"),
						Update: []*gnmipb.Update{
							{
								Path: wantLeafPath("2001:db8::1/128", "1", "matched-octets"),
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_UintVal{
										UintVal: 54321,
									},
								},
							},
							{
								Path: wantDestAddrPath("2001:db8::1/128", "1", "ipv6"),
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_StringVal{
										StringVal: "2001:db8::1/128",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			desc: "unmatched path",
			in: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1234,
						Prefix: &gnmipb.Path{
							Origin: "eos_native",
							Target: "target1",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "Smash"},
										{Name: "flexCounters"},
										{Name: "counterTable"},
										{Name: "DecapGroup"},
										{Name: "4294967295"},
										{Name: "otherCounter"},
										{Name: "1.2.3.4"},
										{Name: "pkts"},
									},
								},
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_UintVal{
										UintVal: 42,
									},
								},
							},
						},
					},
				},
			},
			wantOut: nil,
		},
		{
			desc: "delete prefix path",
			in: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1234,
						Prefix: &gnmipb.Path{
							Origin: "eos_native",
							Target: "target1",
						},
						Delete: []*gnmipb.Path{
							{
								Elem: []*gnmipb.PathElem{
									{Name: "Smash"},
									{Name: "flexCounters"},
									{Name: "counterTable"},
									{Name: "DecapGroup"},
									{Name: "4294967295"},
									{Name: "decapGroupCounter"},
									{Name: "1.2.3.4"},
								},
							},
						},
					},
				},
			},
			wantOut: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1234,
						Prefix:    wantPrefix("target1"),
						Delete: []*gnmipb.Path{
							wantLeafPath("1.2.3.4", "1", "matched-pkts"),
							wantLeafPath("1.2.3.4", "1", "matched-octets"),
							wantDestAddrPath("1.2.3.4", "1", "ipv4"),
						},
					},
				},
			},
		},
		{
			desc:    "nil notification",
			in:      &gnmipb.SubscribeResponse{},
			wantOut: nil,
		},
		{
			desc: "update with prefix elements",
			in: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1234,
						Prefix: &gnmipb.Path{
							Origin: "eos_native",
							Target: "target1",
							Elem: []*gnmipb.PathElem{
								{Name: "Smash"},
								{Name: "flexCounters"},
							},
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "counterTable"},
										{Name: "DecapGroup"},
										{Name: "4294967295"},
										{Name: "decapGroupCounter"},
										{Name: "1.2.3.4"},
										{Name: "pkts"},
									},
								},
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_UintVal{
										UintVal: 42,
									},
								},
							},
						},
					},
				},
			},
			wantOut: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1234,
						Prefix:    wantPrefix("target1"),
						Update: []*gnmipb.Update{
							{
								Path: wantLeafPath("1.2.3.4", "1", "matched-pkts"),
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_UintVal{
										UintVal: 42,
									},
								},
							},
							{
								Path: wantDestAddrPath("1.2.3.4", "1", "ipv4"),
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_StringVal{
										StringVal: "1.2.3.4",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			desc: "pkts and octets delete and unmatched delete",
			in: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1234,
						Prefix: &gnmipb.Path{
							Origin: "eos_native",
							Target: "target1",
						},
						Delete: []*gnmipb.Path{
							{
								Elem: []*gnmipb.PathElem{
									{Name: "Smash"},
									{Name: "flexCounters"},
									{Name: "counterTable"},
									{Name: "DecapGroup"},
									{Name: "4294967295"},
									{Name: "decapGroupCounter"},
									{Name: "1.2.3.4"},
									{Name: "pkts"},
								},
							},
							{
								Elem: []*gnmipb.PathElem{
									{Name: "Smash"},
									{Name: "flexCounters"},
									{Name: "counterTable"},
									{Name: "DecapGroup"},
									{Name: "4294967295"},
									{Name: "decapGroupCounter"},
									{Name: "1.2.3.4"},
									{Name: "octets"},
								},
							},
							{
								Elem: []*gnmipb.PathElem{
									{Name: "unmatched"},
								},
							},
						},
					},
				},
			},
			wantOut: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1234,
						Prefix:    wantPrefix("target1"),
						Delete: []*gnmipb.Path{
							wantLeafPath("1.2.3.4", "1", "matched-pkts"),
							wantLeafPath("1.2.3.4", "1", "matched-octets"),
						},
					},
				},
			},
		},
		{
			desc: "update with nil prefix",
			in: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1234,
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Origin: "eos_native",
									Elem: []*gnmipb.PathElem{
										{Name: "Smash"},
										{Name: "flexCounters"},
										{Name: "counterTable"},
										{Name: "DecapGroup"},
										{Name: "4294967295"},
										{Name: "decapGroupCounter"},
										{Name: "1.2.3.4"},
										{Name: "pkts"},
									},
								},
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_UintVal{
										UintVal: 42,
									},
								},
							},
						},
					},
				},
			},
			wantOut: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1234,
						Prefix:    wantPrefix(""),
						Update: []*gnmipb.Update{
							{
								Path: wantLeafPath("1.2.3.4", "1", "matched-pkts"),
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_UintVal{
										UintVal: 42,
									},
								},
							},
							{
								Path: wantDestAddrPath("1.2.3.4", "1", "ipv4"),
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_StringVal{
										StringVal: "1.2.3.4",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			desc: "non-empty prefix elem update",
			in: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1234,
						Prefix: &gnmipb.Path{
							Origin: "eos_native",
							Target: "target1",
							Elem: []*gnmipb.PathElem{
								{Name: "Smash"},
							},
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "flexCounters"},
										{Name: "counterTable"},
										{Name: "DecapGroup"},
										{Name: "4294967295"},
										{Name: "decapGroupCounter"},
										{Name: "1.2.3.4"},
										{Name: "pkts"},
									},
								},
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_UintVal{
										UintVal: 42,
									},
								},
							},
						},
					},
				},
			},
			wantOut: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1234,
						Prefix:    wantPrefix("target1"),
						Update: []*gnmipb.Update{
							{
								Path: wantLeafPath("1.2.3.4", "1", "matched-pkts"),
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_UintVal{
										UintVal: 42,
									},
								},
							},
							{
								Path: wantDestAddrPath("1.2.3.4", "1", "ipv4"),
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_StringVal{
										StringVal: "1.2.3.4",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			desc: "valid real device output format with 4294967295 and decapIp_tunnelType key",
			in: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1234,
						Prefix: &gnmipb.Path{
							Origin: "eos_native",
							Target: "target1",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "Smash"},
										{Name: "flexCounters"},
										{Name: "counterTable"},
										{Name: "DecapGroup"},
										{Name: "4294967295"},
										{Name: "decapGroupCounter"},
										{Name: "22.22.22.2/32_udp"},
										{Name: "pkts"},
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
										{Name: "Smash"},
										{Name: "flexCounters"},
										{Name: "counterTable"},
										{Name: "DecapGroup"},
										{Name: "4294967295"},
										{Name: "decapGroupCounter"},
										{Name: "22.22.22.2/32_udp"},
										{Name: "octets"},
									},
								},
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_UintVal{
										UintVal: 166,
									},
								},
							},
						},
					},
				},
			},
			wantOut: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1234,
						Prefix:    wantPrefix("target1"),
						Update: []*gnmipb.Update{
							{
								Path: wantLeafPath("22.22.22.2/32", "1", "matched-pkts"),
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_UintVal{
										UintVal: 1,
									},
								},
							},
							{
								Path: wantDestAddrPath("22.22.22.2/32", "1", "ipv4"),
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_StringVal{
										StringVal: "22.22.22.2/32",
									},
								},
							},
							{
								Path: wantLeafPath("22.22.22.2/32", "1", "matched-octets"),
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_UintVal{
										UintVal: 166,
									},
								},
							},
							{
								Path: wantDestAddrPath("22.22.22.2/32", "1", "ipv4"),
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_StringVal{
										StringVal: "22.22.22.2/32",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			desc: "valid numeric rule id preserves sequence-id uint32",
			in: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1726943800000000000,
						Prefix: &gnmipb.Path{
							Origin: "eos_native",
							Target: "target1",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "Smash"},
										{Name: "flexCounters"},
										{Name: "counterTable"},
										{Name: "DecapGroup"},
										{Name: "4294967295"},
										{Name: "decapGroupCounter"},
										{Name: "42"},
										{Name: "pkts"},
									},
								},
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_UintVal{
										UintVal: 500,
									},
								},
							},
						},
					},
				},
			},
			wantOut: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1726943800000000000,
						Prefix:    wantPrefix("target1"),
						Update: []*gnmipb.Update{
							{
								Path: wantLeafPath("42", "42", "matched-pkts"),
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_UintVal{
										UintVal: 500,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			desc: "unmatched non-eos_native origin in update",
			in: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Prefix: &gnmipb.Path{Origin: "openconfig"},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "Smash"},
										{Name: "flexCounters"},
										{Name: "counterTable"},
										{Name: "DecapGroup"},
										{Name: "4294967295"},
										{Name: "decapGroupCounter"},
										{Name: "10.61.61.148/30"},
										{Name: "pkts"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 10}},
							},
						},
					},
				},
			},
			wantOut: nil,
		},
		{
			desc: "unmatched path element count != 8",
			in: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Prefix: &gnmipb.Path{Origin: "eos_native"},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "Smash"},
										{Name: "flexCounters"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 10}},
							},
						},
					},
				},
			},
			wantOut: nil,
		},
		{
			desc: "unmatched leaf name",
			in: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Prefix: &gnmipb.Path{Origin: "eos_native"},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "Smash"},
										{Name: "flexCounters"},
										{Name: "counterTable"},
										{Name: "DecapGroup"},
										{Name: "4294967295"},
										{Name: "decapGroupCounter"},
										{Name: "10.61.61.148/30"},
										{Name: "unsupported-leaf"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 10}},
							},
						},
					},
				},
			},
			wantOut: nil,
		},
	}

	ft := New()

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			gotOut, err := ft.Translate(tt.in)
			if (err != nil) != tt.wantError {
				t.Errorf("Translate() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if diff := cmp.Diff(tt.wantOut, gotOut, protocmp.Transform()); diff != "" {
				t.Errorf("Translate() output mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func BenchmarkTranslate(b *testing.B) {
	in := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 1726943800000000000,
				Prefix: &gnmipb.Path{
					Origin: "eos_native",
					Target: "cx01.apa98",
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "Smash"},
								{Name: "flexCounters"},
								{Name: "counterTable"},
								{Name: "DecapGroup"},
								{Name: "4294967295"},
								{Name: "decapGroupCounter"},
								{Name: "10.61.61.148/30_udp"},
								{Name: "pkts"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 96971532,
							},
						},
					},
				},
			},
		},
	}

	ft := New()

	b.ResetTimer()
	b.ReportAllocs()
	for range b.N {
		if _, err := ft.Translate(in); err != nil {
			b.Fatalf("Translate failed: %v", err)
		}
	}
}

// wantPrefix returns the OpenConfig prefix the translator is expected to emit.
func wantPrefix(target string) *gnmipb.Path {
	return &gnmipb.Path{
		Origin: "openconfig",
		Target: target,
		Elem: []*gnmipb.PathElem{
			{Name: "network-instances"},
			{Name: "network-instance", Key: map[string]string{"name": "default"}},
			{Name: "policy-forwarding"},
			{Name: "policies"},
		},
	}
}

// wantRuleElems returns the expected policy/rule elements, relative to wantPrefix.
func wantRuleElems(policyID, sequenceID string) []*gnmipb.PathElem {
	return []*gnmipb.PathElem{
		{Name: "policy", Key: map[string]string{"policy-id": policyID}},
		{Name: "rules"},
		{Name: "rule", Key: map[string]string{"sequence-id": sequenceID}},
	}
}

// wantLeafPath returns the expected relative path of a rule state leaf (e.g. matched-pkts).
func wantLeafPath(policyID, sequenceID, leaf string) *gnmipb.Path {
	elems := append(wantRuleElems(policyID, sequenceID), &gnmipb.PathElem{Name: "state"}, &gnmipb.PathElem{Name: leaf})
	return &gnmipb.Path{Elem: elems}
}

// wantDestAddrPath returns the expected relative path of the rule's ipv4/ipv6 destination-address leaf.
func wantDestAddrPath(policyID, sequenceID, family string) *gnmipb.Path {
	elems := append(wantRuleElems(policyID, sequenceID),
		&gnmipb.PathElem{Name: family}, &gnmipb.PathElem{Name: "state"}, &gnmipb.PathElem{Name: "destination-address"})
	return &gnmipb.Path{Elem: elems}
}
