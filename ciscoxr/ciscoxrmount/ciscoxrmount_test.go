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

package ciscoxrmount

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
					Origin: "openconfig",
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "file-system"},
								{Name: "node", Key: map[string]string{"node-name": "0/RP0/CPU0"}},
								{Name: "file-system"},
								{Name: "size"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 2000000,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "file-system"},
								{Name: "node", Key: map[string]string{"node-name": "0/RP0/CPU0"}},
								{Name: "file-system"},
								{Name: "free"},
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
								{Name: "file-system"},
								{Name: "node", Key: map[string]string{"node-name": "0/RP0/CPU0"}},
								{Name: "file-system"},
								{Name: "prefixes"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "harddisk",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "file-system"},
								{Name: "node", Key: map[string]string{"node-name": "0/RP0/CPU0"}},
								{Name: "file-system"},
								{Name: "size"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 2000000,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "file-system"},
								{Name: "node", Key: map[string]string{"node-name": "0/RP0/CPU0"}},
								{Name: "file-system"},
								{Name: "free"},
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
								{Name: "file-system"},
								{Name: "node", Key: map[string]string{"node-name": "0/RP0/CPU0"}},
								{Name: "file-system"},
								{Name: "prefixes"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "disk1",
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
								{Name: "system"},
								{Name: "mount-points"},
								{Name: "mount-point", Key: map[string]string{"name": "0/RP0/CPU0-harddisk"}},
								{Name: "state"},
								{Name: "size"},
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
								{Name: "system"},
								{Name: "mount-points"},
								{Name: "mount-point", Key: map[string]string{"name": "0/RP0/CPU0-harddisk"}},
								{Name: "state"},
								{Name: "available"},
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
								{Name: "system"},
								{Name: "mount-points"},
								{Name: "mount-point", Key: map[string]string{"name": "0/RP0/CPU0-harddisk"}},
								{Name: "state"},
								{Name: "utilized"},
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
								{Name: "system"},
								{Name: "mount-points"},
								{Name: "mount-point", Key: map[string]string{"name": "0/RP0/CPU0-harddisk"}},
								{Name: "state"},
								{Name: "name"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "0/RP0/CPU0-harddisk",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "system"},
								{Name: "mount-points"},
								{Name: "mount-point", Key: map[string]string{"name": "0/RP0/CPU0-disk1"}},
								{Name: "state"},
								{Name: "size"},
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
								{Name: "system"},
								{Name: "mount-points"},
								{Name: "mount-point", Key: map[string]string{"name": "0/RP0/CPU0-disk1"}},
								{Name: "state"},
								{Name: "available"},
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
								{Name: "system"},
								{Name: "mount-points"},
								{Name: "mount-point", Key: map[string]string{"name": "0/RP0/CPU0-disk1"}},
								{Name: "state"},
								{Name: "utilized"},
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
								{Name: "system"},
								{Name: "mount-points"},
								{Name: "mount-point", Key: map[string]string{"name": "0/RP0/CPU0-disk1"}},
								{Name: "state"},
								{Name: "name"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "0/RP0/CPU0-disk1",
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
								{Name: "file-system"},
								{Name: "node", Key: map[string]string{"node-name": "0/RP0/CPU0"}},
								{Name: "file-system"},
								{Name: "size"},
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
								{Name: "file-system"},
								{Name: "node", Key: map[string]string{"node-name": "0/RP0/CPU0"}},
								{Name: "file-system"},
								{Name: "free"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 2000000,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "file-system"},
								{Name: "node", Key: map[string]string{"node-name": "0/RP0/CPU0"}},
								{Name: "file-system"},
								{Name: "prefixes"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "harddisk",
							},
						},
					},
				},
			},
		},
	}
	inconsistentFileSystemSR := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 123,
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "file-system"},
								{Name: "node", Key: map[string]string{"node-name": "0/RP0/CPU0"}},
								{Name: "file-system"},
								{Name: "size"},
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
								{Name: "file-system"},
								{Name: "node", Key: map[string]string{"node-name": "0/RP0/CPU0"}},
								{Name: "file-system"},
								{Name: "free"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{
								UintVal: 2000000,
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "file-system"},
								{Name: "node", Key: map[string]string{"node-name": "0/RP0/CPU0"}},
								{Name: "file-system"},
								{Name: "prefixes"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "harddisk",
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
			name:    "invalid SR free space > size",
			input:   invalidSR,
			wantErr: true,
		},
		{
			name:    "inconsistent SR free do not exist",
			input:   inconsistentFileSystemSR,
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
					t.Fatalf("Unexpected diff from translate() = %v, want %v:\n%s", sr, test.want, diff)
				}
			}
		})
	}
}
