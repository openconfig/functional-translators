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

package ciscoxrqos

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
					Origin: "Cisco-IOS-XR-qos-ma-oper",
					Elem: []*gnmipb.PathElem{
						{Name: "qos"},
						{Name: "interface-table"},
						{Name: "interface", Key: map[string]string{"interface-name": "Bundle-Ether1"}},
					},
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "output"},
								{Name: "service-policy-names"},
								{Name: "service-policy-instance", Key: map[string]string{"service-policy-name": "INGRESS_POLICY"}},
								{Name: "statistics"},
								{Name: "class-stats"},
								{Name: "class-name"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "inet-mplsogre-classifier-nc1",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "output"},
								{Name: "service-policy-names"},
								{Name: "service-policy-instance", Key: map[string]string{"service-policy-name": "INGRESS_POLICY"}},
								{Name: "statistics"},
								{Name: "class-stats"},
								{Name: "general-stats"},
								{Name: "transmit-bytes"},
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
								{Name: "output"},
								{Name: "service-policy-names"},
								{Name: "service-policy-instance", Key: map[string]string{"service-policy-name": "INGRESS_POLICY"}},
								{Name: "statistics"},
								{Name: "class-stats"},
								{Name: "general-stats"},
								{Name: "transmit-packets"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 10,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "output"},
								{Name: "service-policy-names"},
								{Name: "service-policy-instance", Key: map[string]string{"service-policy-name": "INGRESS_POLICY"}},
								{Name: "statistics"},
								{Name: "class-stats"},
								{Name: "general-stats"},
								{Name: "total-drop-bytes"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 200,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "output"},
								{Name: "service-policy-names"},
								{Name: "service-policy-instance", Key: map[string]string{"service-policy-name": "INGRESS_POLICY"}},
								{Name: "statistics"},
								{Name: "class-stats"},
								{Name: "general-stats"},
								{Name: "total-drop-packets"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 20,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "input"},
								{Name: "service-policy-names"},
								{Name: "service-policy-instance", Key: map[string]string{"service-policy-name": "INGRESS_POLICY"}},
								{Name: "statistics"},
								{Name: "class-stats"},
								{Name: "class-name"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "inet-mplsogre-classifier-nc1",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "input"},
								{Name: "service-policy-names"},
								{Name: "service-policy-instance", Key: map[string]string{"service-policy-name": "INGRESS_POLICY"}},
								{Name: "statistics"},
								{Name: "class-stats"},
								{Name: "general-stats"},
								{Name: "pre-policy-matched-bytes"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 300,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "input"},
								{Name: "service-policy-names"},
								{Name: "service-policy-instance", Key: map[string]string{"service-policy-name": "INGRESS_POLICY"}},
								{Name: "statistics"},
								{Name: "class-stats"},
								{Name: "general-stats"},
								{Name: "pre-policy-matched-packets"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 30,
							},
						},
					},
				},
			},
		},
	}
	successMemberSR := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 123,
				Prefix: &gnmipb.Path{
					Origin: "Cisco-IOS-XR-qos-ma-oper",
					Elem: []*gnmipb.PathElem{
						{Name: "qos"},
						{Name: "interface-table"},
						{Name: "interface", Key: map[string]string{"interface-name": "Bundle-Ether1"}},
						{Name: "member-interfaces"},
						{Name: "member-interface", Key: map[string]string{"interface-name": "FourHundredGigE0/0/0/2"}},
						{Name: "output"},
						{Name: "service-policy-names"},
						{Name: "service-policy-instance", Key: map[string]string{"service-policy-name": "INGRESS_POLICY"}},
						{Name: "statistics"},
					},
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "class-stats"},
								{Name: "class-name"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "inet-mplsogre-classifier-nc1",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "class-stats"},
								{Name: "general-stats"},
								{Name: "transmit-bytes"},
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
								{Name: "class-stats"},
								{Name: "general-stats"},
								{Name: "transmit-packets"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 10,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "class-stats"},
								{Name: "general-stats"},
								{Name: "total-drop-bytes"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 200,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "class-stats"},
								{Name: "general-stats"},
								{Name: "total-drop-packets"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 20,
							},
						},
					},
				},
			},
		},
	}
	mismatchOpLenSR := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 123,
				Prefix: &gnmipb.Path{
					Origin: "Cisco-IOS-XR-qos-ma-oper",
					Elem: []*gnmipb.PathElem{
						{Name: "qos"},
						{Name: "interface-table"},
						{Name: "interface", Key: map[string]string{"interface-name": "Bundle-Ether1"}},
						{Name: "output"},
						{Name: "service-policy-names"},
						{Name: "service-policy-instance", Key: map[string]string{"service-policy-name": "INGRESS_POLICY"}},
						{Name: "statistics"},
					},
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "class-stats"},
								{Name: "class-name"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "inet-mplsogre-classifier-nc1",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "class-stats"},
								{Name: "general-stats"},
								{Name: "transmit-packets"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 10,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "class-stats"},
								{Name: "general-stats"},
								{Name: "total-drop-bytes"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 200,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "class-stats"},
								{Name: "general-stats"},
								{Name: "total-drop-packets"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 20,
							},
						},
					},
				},
			},
		},
	}
	mismatchInLenSR := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 123,
				Prefix: &gnmipb.Path{
					Origin: "Cisco-IOS-XR-qos-ma-oper",
					Elem: []*gnmipb.PathElem{
						{Name: "qos"},
						{Name: "interface-table"},
						{Name: "interface", Key: map[string]string{"interface-name": "Bundle-Ether1"}},
						{Name: "input"},
						{Name: "service-policy-names"},
						{Name: "service-policy-instance", Key: map[string]string{"service-policy-name": "INGRESS_POLICY"}},
						{Name: "statistics"},
					},
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "class-stats"},
								{Name: "class-name"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "inet-mplsogre-classifier-nc1",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "class-stats"},
								{Name: "general-stats"},
								{Name: "pre-policy-matched-bytes"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 100,
							},
						},
					},
				},
			},
		},
	}
	wrongClassNameSR := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 123,
				Prefix: &gnmipb.Path{
					Origin: "Cisco-IOS-XR-qos-ma-oper",
					Elem: []*gnmipb.PathElem{
						{Name: "qos"},
						{Name: "interface-table"},
						{Name: "interface", Key: map[string]string{"interface-name": "Bundle-Ether1"}},
						{Name: "input"},
						{Name: "service-policy-names"},
						{Name: "service-policy-instance", Key: map[string]string{"service-policy-name": "INGRESS_POLICY"}},
						{Name: "statistics"},
					},
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "class-stats"},
								{Name: "class-name"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "abcd1",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "class-stats"},
								{Name: "general-stats"},
								{Name: "pre-policy-matched-bytes"},
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
								{Name: "class-stats"},
								{Name: "general-stats"},
								{Name: "pre-policy-matched-packets"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 100,
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
								{Name: "qos"},
								{Name: "interfaces"},
								{Name: "interface", Key: map[string]string{"interface-id": "Bundle-Ether1"}},
								{Name: "output"},
								{Name: "queues"},
								{Name: "queue", Key: map[string]string{"name": "inet-mplsogre-classifier-nc1"}},
								{Name: "state"},
								{Name: "dropped-octets"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 200,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "qos"},
								{Name: "interfaces"},
								{Name: "interface", Key: map[string]string{"interface-id": "Bundle-Ether1"}},
								{Name: "output"},
								{Name: "queues"},
								{Name: "queue", Key: map[string]string{"name": "inet-mplsogre-classifier-nc1"}},
								{Name: "state"},
								{Name: "dropped-pkts"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 20,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "qos"},
								{Name: "interfaces"},
								{Name: "interface", Key: map[string]string{"interface-id": "Bundle-Ether1"}},
								{Name: "output"},
								{Name: "queues"},
								{Name: "queue", Key: map[string]string{"name": "inet-mplsogre-classifier-nc1"}},
								{Name: "state"},
								{Name: "transmit-octets"},
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
								{Name: "qos"},
								{Name: "interfaces"},
								{Name: "interface", Key: map[string]string{"interface-id": "Bundle-Ether1"}},
								{Name: "output"},
								{Name: "queues"},
								{Name: "queue", Key: map[string]string{"name": "inet-mplsogre-classifier-nc1"}},
								{Name: "state"},
								{Name: "transmit-pkts"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 10,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "qos"},
								{Name: "interfaces"},
								{Name: "interface", Key: map[string]string{"interface-id": "Bundle-Ether1"}},
								{Name: "input"},
								{Name: "classifiers"},
								{Name: "classifier", Key: map[string]string{"type": "IPV4"}},
								{Name: "terms"},
								{Name: "term", Key: map[string]string{"id": "nc1"}},
								{Name: "state"},
								{Name: "matched-octets"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 300,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "qos"},
								{Name: "interfaces"},
								{Name: "interface", Key: map[string]string{"interface-id": "Bundle-Ether1"}},
								{Name: "input"},
								{Name: "classifiers"},
								{Name: "classifier", Key: map[string]string{"type": "IPV4"}},
								{Name: "terms"},
								{Name: "term", Key: map[string]string{"id": "nc1"}},
								{Name: "state"},
								{Name: "matched-packets"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 30,
							},
						},
					},
				},
			},
		},
	}
	successMemberOutput := &gnmipb.SubscribeResponse{
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
								{Name: "qos"},
								{Name: "interfaces"},
								{Name: "interface", Key: map[string]string{"interface-id": "FourHundredGigE0/0/0/2"}},
								{Name: "output"},
								{Name: "queues"},
								{Name: "queue", Key: map[string]string{"name": "inet-mplsogre-classifier-nc1"}},
								{Name: "state"},
								{Name: "dropped-octets"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 200,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "qos"},
								{Name: "interfaces"},
								{Name: "interface", Key: map[string]string{"interface-id": "FourHundredGigE0/0/0/2"}},
								{Name: "output"},
								{Name: "queues"},
								{Name: "queue", Key: map[string]string{"name": "inet-mplsogre-classifier-nc1"}},
								{Name: "state"},
								{Name: "dropped-pkts"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 20,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "qos"},
								{Name: "interfaces"},
								{Name: "interface", Key: map[string]string{"interface-id": "FourHundredGigE0/0/0/2"}},
								{Name: "output"},
								{Name: "queues"},
								{Name: "queue", Key: map[string]string{"name": "inet-mplsogre-classifier-nc1"}},
								{Name: "state"},
								{Name: "transmit-octets"},
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
								{Name: "qos"},
								{Name: "interfaces"},
								{Name: "interface", Key: map[string]string{"interface-id": "FourHundredGigE0/0/0/2"}},
								{Name: "output"},
								{Name: "queues"},
								{Name: "queue", Key: map[string]string{"name": "inet-mplsogre-classifier-nc1"}},
								{Name: "state"},
								{Name: "transmit-pkts"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 10,
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
			name:  "success member interface",
			input: successMemberSR,
			want:  successMemberOutput,
		},
		{
			name:    "mismatch stats len for output qos",
			input:   mismatchOpLenSR,
			wantErr: true,
		},
		{
			name:    "mismatch stats len for input qos",
			input:   mismatchInLenSR,
			wantErr: true,
		},
		{
			name:  "class name does not have any parts separated by -",
			input: wrongClassNameSR,
			want:  nil,
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
