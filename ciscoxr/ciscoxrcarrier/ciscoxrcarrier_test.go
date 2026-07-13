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

package ciscoxrcarrier

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

func TestTranslate(t *testing.T) {
	successCarrierSR := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 123,
				Prefix: &gnmipb.Path{
					Origin: "Cisco-IOS-XR-infra-statsd-oper",
					Elem: []*gnmipb.PathElem{
						{Name: "infra-statistics"},
						{Name: "interfaces"},
					},
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "interface", Key: map[string]string{"interface-name": "Fh0/0/0/0"}},
								{Name: "generic-counters"},
								{Name: "carrier-transitions"},
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
								{Name: "interface", Key: map[string]string{"interface-name": "Fh0/0/0/1"}},
								{Name: "generic-counters"},
								{Name: "carrier-transitions"},
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
								{Name: "interface", Key: map[string]string{"interface-name": "Fh0/0/0/3"}},
								{Name: "generic-counters"},
								{Name: "carrier-transitions"},
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
	invalidSR := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 123,
				Prefix: &gnmipb.Path{
					Origin: "Cisco-IOS-XR-infra-statsd-oper",
					Elem: []*gnmipb.PathElem{
						{Name: "infra-statistics"},
						{Name: "interfacesssss"},
						{Name: "interfacexx", Key: map[string]string{"name": "1"}},
					},
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "carrier"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 1,
							},
						},
					},
				},
			},
		},
	}
	nilUpdateSR := &gnmipb.SubscribeResponse{}
	nilIntfsSR := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 123,
				Prefix: &gnmipb.Path{
					Origin: "Cisco-IOS-XR-infra-statsd-oper",
					Elem: []*gnmipb.PathElem{
						{Name: "infra-statistics"},
					},
				},
				Delete: []*gnmipb.Path{
					{
						Elem: []*gnmipb.PathElem{
							{Name: "interfaces"},
						},
					},
				},
			},
		},
	}
	nilStatsSR := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 123,
				Prefix: &gnmipb.Path{
					Origin: "Cisco-IOS-XR-infra-statsd-oper",
				},
				Delete: []*gnmipb.Path{
					{
						Elem: []*gnmipb.PathElem{
							{Name: "infra-statistics"},
						},
					},
				},
			},
		},
	}
	nilGenericCountersSR := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 123,
				Prefix: &gnmipb.Path{
					Origin: "Cisco-IOS-XR-infra-statsd-oper",
					Elem: []*gnmipb.PathElem{
						{Name: "infra-statistics"},
						{Name: "interfaces"},
					},
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "interface", Key: map[string]string{"interface-name": "Fh0/0/0/0"}},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_JsonIetfVal{
								JsonIetfVal: []byte("{}"),
							},
						},
					},
				},
				Delete: []*gnmipb.Path{
					{
						Elem: []*gnmipb.PathElem{
							{Name: "interface", Key: map[string]string{"interface-name": "Fh0/0/0/0"}},
							{Name: "generic-counters"},
						},
					},
				},
			},
		},
	}
	nilCarrierTransitionsSR := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 123,
				Prefix: &gnmipb.Path{
					Origin: "Cisco-IOS-XR-infra-statsd-oper",
					Elem: []*gnmipb.PathElem{
						{Name: "infra-statistics"},
						{Name: "interfaces"},
					},
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "interface", Key: map[string]string{"interface-name": "Fh0/0/0/0"}},
								{Name: "generic-counters"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_JsonIetfVal{
								JsonIetfVal: []byte("{}"),
							},
						},
					},
				},
				Delete: []*gnmipb.Path{
					{
						Elem: []*gnmipb.PathElem{
							{Name: "interface", Key: map[string]string{"interface-name": "Fh0/0/0/0"}},
							{Name: "generic-counters"},
							{Name: "carrier-transitions"},
						},
					},
				},
			},
		},
	}
	nilInterfaceSR := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 123,
				Prefix: &gnmipb.Path{
					Origin: "Cisco-IOS-XR-infra-statsd-oper",
					Elem: []*gnmipb.PathElem{
						{Name: "infra-statistics"},
						{Name: "interfaces"},
					},
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "interface", Key: map[string]string{"interface-name": "Fh0/0/0/0"}},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_JsonIetfVal{
								JsonIetfVal: []byte("null"),
							},
						},
					},
				},
			},
		},
	}
	emptyOutput := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 123,
				Prefix: &gnmipb.Path{
					Origin: "openconfig",
				},
				Update: []*gnmipb.Update{},
			},
		},
	}
	successCarrierOutput := &gnmipb.SubscribeResponse{
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
								{Name: "interfaces"},
								{Name: "interface", Key: map[string]string{"name": "Fh0/0/0/0"}},
								{Name: "ethernet"},
								{Name: "state"},
								{Name: "counters"},
								{Name: "phy-carrier-transitions"},
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
								{Name: "interfaces"},
								{Name: "interface", Key: map[string]string{"name": "Fh0/0/0/1"}},
								{Name: "ethernet"},
								{Name: "state"},
								{Name: "counters"},
								{Name: "phy-carrier-transitions"},
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
								{Name: "interfaces"},
								{Name: "interface", Key: map[string]string{"name": "Fh0/0/0/3"}},
								{Name: "ethernet"},
								{Name: "state"},
								{Name: "counters"},
								{Name: "phy-carrier-transitions"},
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
	tests := []struct {
		name    string
		input   *gnmipb.SubscribeResponse
		want    *gnmipb.SubscribeResponse
		wantErr bool
		wantNil bool
	}{
		{
			name:  "SuccessSwitchStats",
			input: successCarrierSR,
			want:  successCarrierOutput,
		},
		{
			name:    "InvalidSR",
			input:   invalidSR,
			wantErr: true,
		},
		{
			name:    "NilUpdate",
			input:   nilUpdateSR,
			wantNil: true,
		},
		{
			name:    "NilInterfaces",
			input:   nilIntfsSR,
			wantNil: true,
		},
		{
			name:    "NilStats",
			input:   nilStatsSR,
			wantNil: true,
		},
		{
			name:  "NilGenericCounters",
			input: nilGenericCountersSR,
			want:  emptyOutput,
		},
		{
			name:  "NilCarrierTransitions",
			input: nilCarrierTransitionsSR,
			want:  emptyOutput,
		},
		{
			name:  "NilInterface",
			input: nilInterfaceSR,
			want:  emptyOutput,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ft := New()
			sr, err := ft.Translate(test.input)
			errorMatchesExpectation := (err != nil) == test.wantErr
			if !errorMatchesExpectation {
				t.Fatalf("Translate(%v) error got %v, want error %t", test.input, err, test.wantErr)
			}
			if !test.wantErr {
				if test.wantNil {
					if sr != nil {
						t.Fatalf("Translate(%v) got %v, want nil", test.input, sr)
					}
				} else {
					if diff := cmp.Diff(test.want, sr, protocmp.Transform(), protocmp.SortRepeatedFields(&gnmipb.Notification{}, "update")); diff != "" {
						t.Fatalf("Translate(%v) got unexpected diff (-want +got):\n%s", test.input, diff)
					}
				}
			}
		})
	}
}
