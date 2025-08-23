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

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

func TestTranslate(t *testing.T) {
	target := "dev1"
	tests := []struct {
		name string
		in   *gnmipb.SubscribeResponse
		want *gnmipb.SubscribeResponse
	}{
		{
			name: "translate-lacp-mac-to-intf-mac-update-and-delete",
			in: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1,
						Prefix: &gnmipb.Path{
							Target: target,
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lacp"},
										{Name: "interfaces"},
										{
											Name: "interface",
											Key: map[string]string{
												"name": "Port-Channel1",
											},
										},
										{Name: "state"},
										{Name: "system-id-mac"},
									},
								},
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_StringVal{
										StringVal: "c4:ca:2b:80:fb:7e",
									},
								},
							},
						},
						Delete: []*gnmipb.Path{
							{
								Elem: []*gnmipb.PathElem{
									{Name: "lacp"},
									{Name: "interfaces"},
									{
										Name: "interface",
										Key: map[string]string{
											"name": "Port-Channel2",
										},
									},
									{Name: "state"},
									{Name: "system-id-mac"},
								},
							},
							{
								Elem: []*gnmipb.PathElem{
									{Name: "lacp"},
									{Name: "interfaces"},
									{
										Name: "interface",
										Key: map[string]string{
											"name": "Port-Channel3",
										},
									},
								},
							},
						},
					},
				},
			},
			want: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1,
						Prefix: &gnmipb.Path{
							Origin: "openconfig",
							Target: target,
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "interfaces"},
										{
											Name: "interface",
											Key: map[string]string{
												"name": "Port-Channel1",
											},
										},
										{Name: "ethernet"},
										{Name: "state"},
										{Name: "mac-address"},
									},
								},
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_StringVal{
										StringVal: "c4:ca:2b:80:fb:7e",
									},
								},
							},
						},
						Delete: []*gnmipb.Path{
							{
								Elem: []*gnmipb.PathElem{
									{Name: "interfaces"},
									{
										Name: "interface",
										Key: map[string]string{
											"name": "Port-Channel2",
										},
									},
									{Name: "ethernet"},
									{Name: "state"},
									{Name: "mac-address"},
								},
							},
							{
								Elem: []*gnmipb.PathElem{
									{Name: "interfaces"},
									{
										Name: "interface",
										Key: map[string]string{
											"name": "Port-Channel3",
										},
									},
									{Name: "ethernet"},
									{Name: "state"},
									{Name: "mac-address"},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "translate-lacp-mac-to-intf-mac-update-only",
			in: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1,
						Prefix: &gnmipb.Path{
							Target: target,
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "lacp"},
										{Name: "interfaces"},
										{
											Name: "interface",
											Key: map[string]string{
												"name": "Port-Channel1",
											},
										},
										{Name: "state"},
										{Name: "system-id-mac"},
									},
								},
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_StringVal{
										StringVal: "c4:ca:2b:80:fb:7e",
									},
								},
							},
						},
					},
				},
			},
			want: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1,
						Prefix: &gnmipb.Path{
							Origin: "openconfig",
							Target: target,
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "interfaces"},
										{
											Name: "interface",
											Key: map[string]string{
												"name": "Port-Channel1",
											},
										},
										{Name: "ethernet"},
										{Name: "state"},
										{Name: "mac-address"},
									},
								},
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_StringVal{
										StringVal: "c4:ca:2b:80:fb:7e",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "translate-lacp-mac-to-intf-mac-delete-only",
			in: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1,
						Prefix: &gnmipb.Path{
							Target: target,
						},
						Delete: []*gnmipb.Path{
							{
								Elem: []*gnmipb.PathElem{
									{Name: "lacp"},
									{Name: "interfaces"},
									{
										Name: "interface",
										Key: map[string]string{
											"name": "Port-Channel2",
										},
									},
									{Name: "state"},
									{Name: "system-id-mac"},
								},
							},
							{
								Elem: []*gnmipb.PathElem{
									{Name: "lacp"},
									{Name: "interfaces"},
									{
										Name: "interface",
										Key: map[string]string{
											"name": "Port-Channel3",
										},
									},
								},
							},
						},
					},
				},
			},
			want: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1,
						Prefix: &gnmipb.Path{
							Origin: "openconfig",
							Target: target,
						},
						Delete: []*gnmipb.Path{
							{
								Elem: []*gnmipb.PathElem{
									{Name: "interfaces"},
									{
										Name: "interface",
										Key: map[string]string{
											"name": "Port-Channel2",
										},
									},
									{Name: "ethernet"},
									{Name: "state"},
									{Name: "mac-address"},
								},
							},
							{
								Elem: []*gnmipb.PathElem{
									{Name: "interfaces"},
									{
										Name: "interface",
										Key: map[string]string{
											"name": "Port-Channel3",
										},
									},
									{Name: "ethernet"},
									{Name: "state"},
									{Name: "mac-address"},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "update-pass-through-ehternet-mac-to-ethernet-mac",
			in: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 6,
						Prefix: &gnmipb.Path{
							Target: target,
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "interfaces"},
										{
											Name: "interface",
											Key: map[string]string{
												"name": "Ethernet1",
											},
										},
										{Name: "ethernet"},
										{Name: "state"},
										{Name: "mac-address"},
									},
								},
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_StringVal{
										StringVal: "mac-for-ethernet1-to-be-pass-through",
									},
								},
							},
						},
					},
				},
			},
			want: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 6,
						Prefix: &gnmipb.Path{
							Origin: "openconfig",
							Target: target,
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "interfaces"},
										{
											Name: "interface",
											Key: map[string]string{
												"name": "Ethernet1",
											},
										},
										{Name: "ethernet"},
										{Name: "state"},
										{Name: "mac-address"},
									},
								},
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_StringVal{
										StringVal: "mac-for-ethernet1-to-be-pass-through",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "delete-pass-through-ethernet-mac",
			in: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 6,
						Prefix: &gnmipb.Path{
							Target: target,
						},
						Delete: []*gnmipb.Path{
							{
								Elem: []*gnmipb.PathElem{
									{Name: "interfaces"},
									{
										Name: "interface",
										Key: map[string]string{
											"name": "Ethernet1",
										},
									},
									{Name: "ethernet"},
									{Name: "state"},
									{Name: "mac-address"},
								},
							},
						},
					},
				},
			},
			want: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 6,
						Prefix: &gnmipb.Path{
							Origin: "openconfig",
							Target: target,
						},

						Delete: []*gnmipb.Path{
							{
								Elem: []*gnmipb.PathElem{
									{Name: "interfaces"},
									{
										Name: "interface",
										Key: map[string]string{
											"name": "Ethernet1",
										},
									},
									{Name: "ethernet"},
									{Name: "state"},
									{Name: "mac-address"},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "delete-pass-through-ehternet-intf",
			in: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 6,
						Prefix: &gnmipb.Path{
							Target: target,
						},
						Delete: []*gnmipb.Path{
							{
								Elem: []*gnmipb.PathElem{
									{Name: "interfaces"},
									{
										Name: "interface",
										Key: map[string]string{
											"name": "Ethernet1",
										},
									},
								},
							},
						},
					},
				},
			},
			want: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 6,
						Prefix: &gnmipb.Path{
							Origin: "openconfig",
							Target: target,
						},
						Delete: []*gnmipb.Path{
							{
								Elem: []*gnmipb.PathElem{
									{Name: "interfaces"},
									{
										Name: "interface",
										Key: map[string]string{
											"name": "Ethernet1",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	ft := NewMacFT()
	for _, tc := range tests {

		t.Run(tc.name+"-without-origin", func(t *testing.T) {
			tc.in.GetUpdate().GetPrefix().Origin = ""
			gotSR, gotErr := ft.Translate(tc.in)
			if gotErr != nil {
				t.Fatalf("translate returned an unexpected error: %v", gotErr)
			}
			if diff := cmp.Diff(tc.want, gotSR, protocmp.Transform()); diff != "" {
				t.Errorf("Translate() returned an unexpected diff (-want +got):\n%s", diff)
			}
		})

		t.Run(tc.name+"-with-openconfig-origin", func(t *testing.T) {
			tc.in.GetUpdate().GetPrefix().Origin = "openconfig"
			gotSR, gotErr := ft.Translate(tc.in)
			if gotErr != nil {
				t.Fatalf("translate returned an unexpected error: %v", gotErr)
			}
			if diff := cmp.Diff(tc.want, gotSR, protocmp.Transform()); diff != "" {
				t.Errorf("Translate() returned an unexpected diff (-want +got):\n%s", diff)
			}
		})

	}
}

func BenchmarkTranslate(b *testing.B) {
	in := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 6,
				Prefix: &gnmipb.Path{
					Target: "cx01.fra01",
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "interfaces"},
								{
									Name: "interface",
									Key: map[string]string{
										"name": "Ethernet1",
									},
								},
								{Name: "ethernet"},
								{Name: "state"},
								{Name: "mac-address"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "mac-for-ethernet1-to-be-pass-through",
							},
						},
					},
				},
			},
		},
	}

	ft := NewMacFT()

	b.ResetTimer()
	for range b.N {
		ft.Translate(in)
	}
}
