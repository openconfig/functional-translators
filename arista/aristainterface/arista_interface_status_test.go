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
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"github.com/openconfig/functional-translators/ftutilities"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

func TestStatusTranslate(t *testing.T) {
	tests := []struct {
		name           string
		inputPath      string
		wantOutputPath string
		expectError    bool
	}{
		{
			name:           "AdminStatusNative",
			inputPath:      "testdata/interface_status_admin_native_input.txt",
			wantOutputPath: "testdata/interface_status_admin_native_output.txt",
		},
		{
			name:           "AdminStatusNativeShutdown",
			inputPath:      "testdata/interface_status_admin_native_shutdown_input.txt",
			wantOutputPath: "testdata/interface_status_admin_native_shutdown_output.txt",
		},
		{
			name:           "OperStatusNative",
			inputPath:      "testdata/interface_status_oper_native_input.txt",
			wantOutputPath: "testdata/interface_status_oper_native_output.txt",
		},
		{
			name:           "AdminStatusOCPassthrough",
			inputPath:      "testdata/interface_status_admin_oc_input.txt",
			wantOutputPath: "testdata/interface_status_admin_oc_output.txt",
		},
		{
			name:           "NativeStatusDelete",
			inputPath:      "testdata/interface_status_delete_native_input.txt",
			wantOutputPath: "testdata/interface_status_delete_native_output.txt",
		},
		{
			name:           "OpenConfigStatusDelete",
			inputPath:      "testdata/interface_status_delete_oc_input.txt",
			wantOutputPath: "testdata/interface_status_delete_oc_output.txt",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ft := NewStatusFT()
			inputSR, err := ftutilities.LoadSubscribeResponse(test.inputPath)
			if err != nil {
				t.Fatalf("Failed to load input message: %v", err)
			}
			beforeInput := proto.CloneOf(inputSR)
			gotSR, gotErr := ft.Translate(inputSR)
			if diff := cmp.Diff(beforeInput, inputSR, protocmp.Transform()); diff != "" {
				t.Errorf("ft.Translate(%s) mutated input message (-want +got):\n%s", test.inputPath, diff)
			}
			errorMatchesExpectation := (gotErr != nil) == test.expectError
			if !errorMatchesExpectation {
				t.Fatalf("Unexpected error result returned from translate() = %v, want error %t", gotErr, test.expectError)
			}
			if !test.expectError {
				if test.wantOutputPath == "" {
					if gotSR != nil {
						t.Fatalf("unexpected SubscribeResponse from translate() = %v, want nil", gotSR)
					}
					return
				}
				wantSR, err := ftutilities.LoadSubscribeResponse(test.wantOutputPath)
				if err != nil {
					t.Fatalf("Failed to load want message: %v", err)
				}
				if diff := cmp.Diff(wantSR, gotSR, protocmp.Transform()); diff != "" {
					t.Fatalf("unexpected diff from translate():\n%s", diff)
				}
			}
		})
	}
}

func TestStatusTranslatePrefixAndImmutability(t *testing.T) {
	ocLeafElems := []*gnmipb.PathElem{
		{Name: "interfaces"},
		{Name: "interface", Key: map[string]string{"name": "Vlan100"}},
		{Name: "state"},
		{Name: "admin-status"},
	}

	tests := []struct {
		name   string
		input  *gnmipb.SubscribeResponse
		wantSR *gnmipb.SubscribeResponse
	}{
		{
			name: "nil_prefix_update_and_delete",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 42,
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Origin: "openconfig",
									Target: "dut1",
									Elem:   ocLeafElems,
								},
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_StringVal{StringVal: "UP"},
								},
							},
						},
						Delete: []*gnmipb.Path{
							{
								Origin: "openconfig",
								Target: "dut1",
								Elem:   ocLeafElems,
							},
						},
					},
				},
			},
			wantSR: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 42,
						Prefix: &gnmipb.Path{
							Origin: "openconfig",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{Elem: ocLeafElems},
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_StringVal{StringVal: "UP"},
								},
							},
						},
						Delete: []*gnmipb.Path{
							{Elem: ocLeafElems},
						},
					},
				},
			},
		},
		{
			name: "non_nil_prefix_update_and_delete",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 42,
						Prefix: &gnmipb.Path{
							Origin: "openconfig",
							Target: "dut1",
							Elem:   []*gnmipb.PathElem{{Name: "interfaces"}},
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Origin: "openconfig",
									Target: "dut1",
									Elem:   ocLeafElems[1:],
								},
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_StringVal{StringVal: "UP"},
								},
							},
						},
						Delete: []*gnmipb.Path{
							{
								Origin: "openconfig",
								Target: "dut1",
								Elem:   ocLeafElems[1:],
							},
						},
					},
				},
			},
			wantSR: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 42,
						Prefix: &gnmipb.Path{
							Origin: "openconfig",
							Target: "dut1",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{Elem: ocLeafElems},
								Val: &gnmipb.TypedValue{
									Value: &gnmipb.TypedValue_StringVal{StringVal: "UP"},
								},
							},
						},
						Delete: []*gnmipb.Path{
							{Elem: ocLeafElems},
						},
					},
				},
			},
		},
	}

	ft := NewStatusFT()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			beforeInput := proto.CloneOf(tc.input)
			gotSR, err := ft.Translate(tc.input)
			if err != nil {
				t.Fatalf("ft.Translate(%v) returned unexpected error: %v", tc.input, err)
			}
			if diff := cmp.Diff(tc.wantSR, gotSR, protocmp.Transform()); diff != "" {
				t.Errorf("ft.Translate(%v) output diff (-want +got):\n%s", tc.input, diff)
			}
			if diff := cmp.Diff(beforeInput, tc.input, protocmp.Transform()); diff != "" {
				t.Errorf("ft.Translate(%v) mutated input message (-want +got):\n%s", tc.input, diff)
			}
		})
	}
}

func TestStatusTranslateConcurrentMarshal(t *testing.T) {
	sharedSR := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 99,
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Origin: "openconfig",
							Target: "dut1",
							Elem: []*gnmipb.PathElem{
								{Name: "interfaces"},
								{Name: "interface", Key: map[string]string{"name": "Vlan100"}},
								{Name: "state"},
								{Name: "admin-status"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{StringVal: "UP"},
						},
					},
				},
				Delete: []*gnmipb.Path{
					{
						Origin: "openconfig",
						Target: "dut1",
						Elem: []*gnmipb.PathElem{
							{Name: "interfaces"},
							{Name: "interface", Key: map[string]string{"name": "Vlan100"}},
							{Name: "state"},
							{Name: "oper-status"},
						},
					},
				},
			},
		},
	}

	ft := NewStatusFT()
	var wg sync.WaitGroup
	for range 8 {
		wg.Add(2)
		go func() {
			defer wg.Done()
			for range 100 {
				if _, err := ft.Translate(sharedSR); err != nil {
					t.Errorf("ft.Translate(%v) failed: %v", sharedSR, err)
					return
				}
			}
		}()
		go func() {
			defer wg.Done()
			for range 100 {
				if _, err := proto.Marshal(sharedSR); err != nil {
					t.Errorf("proto.Marshal(%v) failed: %v", sharedSR, err)
					return
				}
			}
		}()
	}
	wg.Wait()
}

