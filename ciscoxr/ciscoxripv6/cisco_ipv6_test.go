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

package ciscoxripv6

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"github.com/openconfig/functional-translators/ftutilities"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

func TestTranslate(t *testing.T) {
	tests := []struct {
		name string
		// inSR is used if inPath is empty.
		inSR   *gnmipb.SubscribeResponse
		inPath string
		// wantSR is used if wantPath is empty.
		wantSR   *gnmipb.SubscribeResponse
		wantPath string
		wantErr  bool
	}{
		{
			name:     "success from file",
			inPath:   "testdata/success_input.txt",
			wantPath: "testdata/success_output.txt",
		},
		{
			name: "no updates",
			inSR: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{},
				},
			},
			wantSR: nil,
		},
		{
			name: "prefix length update before address",
			inSR: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Prefix: &gnmipb.Path{
							Origin: "Cisco-IOS-XR-ipv6-ma-oper",
							Elem: []*gnmipb.PathElem{
								{Name: "ipv6-network"},
								{Name: "nodes"},
								{Name: "node", Key: map[string]string{"node-name": "0/0/CPU0"}},
								{Name: "interface-data"},
								{Name: "vrfs"},
								{Name: "vrf", Key: map[string]string{"vrf-name": "default"}},
								{Name: "briefs"},
								{Name: "brief", Key: map[string]string{"interface-name": "GigabitEthernet0/0/0/0.1"}},
							},
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "address"}, {Name: "prefix-length"}}},
								Val:  &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 64}},
							},
						},
					},
				},
			},
			wantSR: nil,
		},
		{
			name: "multiple addresses",
			inSR: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Prefix: &gnmipb.Path{
							Origin: "Cisco-IOS-XR-ipv6-ma-oper",
							Elem: []*gnmipb.PathElem{
								{Name: "ipv6-network"},
								{Name: "nodes"},
								{Name: "node", Key: map[string]string{"node-name": "0/0/CPU0"}},
								{Name: "interface-data"},
								{Name: "vrfs"},
								{Name: "vrf", Key: map[string]string{"vrf-name": "default"}},
								{Name: "briefs"},
								{Name: "brief", Key: map[string]string{"interface-name": "GigabitEthernet0/0/0/0.1"}},
							},
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "address"}, {Name: "address"}}},
								Val:  &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "fe80::20c:29ff:fe9a:4b1a"}},
							},
							{
								Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "address"}, {Name: "prefix-length"}}},
								Val:  &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 64}},
							},
							{
								Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "address"}, {Name: "address"}}},
								Val:  &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "2001:db8::1"}},
							},
							{
								Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "address"}, {Name: "prefix-length"}}},
								Val:  &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 128}},
							},
						},
					},
				},
			},
			wantSR: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Prefix: &gnmipb.Path{Origin: "openconfig"},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "interfaces"},
										{Name: "interface", Key: map[string]string{"name": "GigabitEthernet0/0/0/0.1"}},
										{Name: "subinterfaces"},
										{Name: "subinterface", Key: map[string]string{"index": "0"}},
										{Name: "ipv6"},
										{Name: "addresses"},
										{Name: "address", Key: map[string]string{"address": "fe80::20c:29ff:fe9a:4b1a"}},
										{Name: "state"},
										{Name: "ip"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "fe80::20c:29ff:fe9a:4b1a"}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "interfaces"},
										{Name: "interface", Key: map[string]string{"name": "GigabitEthernet0/0/0/0.1"}},
										{Name: "subinterfaces"},
										{Name: "subinterface", Key: map[string]string{"index": "0"}},
										{Name: "ipv6"},
										{Name: "addresses"},
										{Name: "address", Key: map[string]string{"address": "fe80::20c:29ff:fe9a:4b1a"}},
										{Name: "state"},
										{Name: "prefix-length"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 64}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "interfaces"},
										{Name: "interface", Key: map[string]string{"name": "GigabitEthernet0/0/0/0.1"}},
										{Name: "subinterfaces"},
										{Name: "subinterface", Key: map[string]string{"index": "0"}},
										{Name: "ipv6"},
										{Name: "addresses"},
										{Name: "address", Key: map[string]string{"address": "2001:db8::1"}},
										{Name: "state"},
										{Name: "ip"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "2001:db8::1"}},
							},
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "interfaces"},
										{Name: "interface", Key: map[string]string{"name": "GigabitEthernet0/0/0/0.1"}},
										{Name: "subinterfaces"},
										{Name: "subinterface", Key: map[string]string{"index": "0"}},
										{Name: "ipv6"},
										{Name: "addresses"},
										{Name: "address", Key: map[string]string{"address": "2001:db8::1"}},
										{Name: "state"},
										{Name: "prefix-length"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 128}},
							},
						},
					},
				},
			},
		},
		{
			name: "invalid path",
			inSR: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Prefix: &gnmipb.Path{
							Origin: "Cisco-IOS-XR-ipv6-ma-oper",
							Elem: []*gnmipb.PathElem{
								{Name: "ipv6-network"},
								{Name: "nodes"},
								{Name: "node", Key: map[string]string{"node-name": "0/0/CPU0"}},
								{Name: "interface-data"},
								{Name: "vrfs"},
								{Name: "vrf", Key: map[string]string{"vrf-name": "default"}},
								{Name: "briefs"},
								{Name: "brief", Key: map[string]string{"interface-name": "GigabitEthernet0/0/0/0.1"}},
							},
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "address"}, {Name: "not-a-valid-element"}}},
								Val:  &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "fe80::20c:29ff:fe9a:4b1a"}},
							},
						},
					},
				},
			},
			wantSR: nil,
		},
		{
			name: "missing interface name",
			inSR: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Prefix: &gnmipb.Path{
							Origin: "Cisco-IOS-XR-ipv6-ma-oper",
							Elem: []*gnmipb.PathElem{
								{Name: "ipv6-network"},
								{Name: "nodes"},
								{Name: "node", Key: map[string]string{"node-name": "0/0/CPU0"}},
								{Name: "interface-data"},
								{Name: "vrfs"},
								{Name: "vrf", Key: map[string]string{"vrf-name": "default"}},
								// Missing briefs and brief elements.
							},
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "briefs"}, {Name: "brief"}, {Name: "address"}, {Name: "address"}}},
								Val:  &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "fe80::20c:29ff:fe9a:4b1a"}},
							},
						},
					},
				},
			},
			// The translator will produce an update with an empty interface name.
			wantSR: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Prefix: &gnmipb.Path{Origin: "openconfig"},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "interfaces"},
										{Name: "interface", Key: map[string]string{"name": ""}},
										{Name: "subinterfaces"},
										{Name: "subinterface", Key: map[string]string{"index": "0"}},
										{Name: "ipv6"},
										{Name: "addresses"},
										{Name: "address", Key: map[string]string{"address": "fe80::20c:29ff:fe9a:4b1a"}},
										{Name: "state"},
										{Name: "ip"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "fe80::20c:29ff:fe9a:4b1a"}},
							},
						},
					},
				},
			},
		},
		{
			name: "zero prefix length",
			inSR: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Prefix: &gnmipb.Path{
							Origin: "Cisco-IOS-XR-ipv6-ma-oper",
							Elem: []*gnmipb.PathElem{
								{Name: "ipv6-network"},
								{Name: "nodes"},
								{Name: "node", Key: map[string]string{"node-name": "0/0/CPU0"}},
								{Name: "interface-data"},
								{Name: "vrfs"},
								{Name: "vrf", Key: map[string]string{"vrf-name": "default"}},
								{Name: "briefs"},
								{Name: "brief", Key: map[string]string{"interface-name": "GigabitEthernet0/0/0/0.1"}},
							},
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "address"}, {Name: "address"}}},
								Val:  &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "fe80::20c:29ff:fe9a:4b1a"}},
							},
							{
								Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "address"}, {Name: "prefix-length"}}},
								Val:  &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 0}},
							},
						},
					},
				},
			},
			wantSR: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Prefix: &gnmipb.Path{Origin: "openconfig"},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "interfaces"},
										{Name: "interface", Key: map[string]string{"name": "GigabitEthernet0/0/0/0.1"}},
										{Name: "subinterfaces"},
										{Name: "subinterface", Key: map[string]string{"index": "0"}},
										{Name: "ipv6"},
										{Name: "addresses"},
										{Name: "address", Key: map[string]string{"address": "fe80::20c:29ff:fe9a:4b1a"}},
										{Name: "state"},
										{Name: "ip"},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "fe80::20c:29ff:fe9a:4b1a"}},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var inSR *gnmipb.SubscribeResponse
			if test.inPath != "" {
				var err error
				inSR, err = ftutilities.LoadSubscribeResponse(test.inPath)
				if err != nil {
					t.Fatalf("Failed to load input message: %v", err)
				}
			} else {
				inSR = test.inSR
			}

			gotSR, err := New().Translate(inSR)
			if (err != nil) != test.wantErr {
				t.Fatalf("translate() got error %v, wantErr %v", err, test.wantErr)
			}
			if test.wantErr {
				return
			}

			var wantSR *gnmipb.SubscribeResponse
			if test.wantPath != "" {
				var err error
				wantSR, err = ftutilities.LoadSubscribeResponse(test.wantPath)
				if err != nil {
					t.Fatalf("Failed to load want message: %v", err)
				}
			} else {
				wantSR = test.wantSR
			}

			if diff := cmp.Diff(wantSR, gotSR, protocmp.Transform(), protocmp.SortRepeatedFields(&gnmipb.Notification{}, "update")); diff != "" {
				t.Errorf("translate() returned diff (-want +got):\n%s", diff)
			}
		})
	}
}
