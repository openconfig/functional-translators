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

package junipermacsecstate

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/functional-translators/ftconsts"
	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

func TestTranslateInterfaceState(t *testing.T) {
	tests := []struct {
		name           string
		notification   *gnmipb.Notification
		expectedLen    int
		expectedStatus []string
		expectedCKN    []string
		expectDelete   bool
	}{
		{
			name: "Interface enable state with complete data",
			notification: &gnmipb.Notification{
				Timestamp: 1000,
				Prefix: &gnmipb.Path{
					Origin: "junos",
					Target: "device1",
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "macsec"},
								{Name: "interfaces"},
								{Name: "interface", Key: map[string]string{"name": "et-0/0/0"}},
								{Name: "name"},
							},
						},
						Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "et-0/0/0"}},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "macsec"},
								{Name: "interfaces"},
								{Name: "interface", Key: map[string]string{"name": "et-0/0/0"}},
								{Name: "mka"},
								{Name: "state"},
								{Name: "jnx-cak-name"},
							},
						},
						Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "abc123def456"}},
					},
				},
			},
			expectedLen:    2,
			expectedStatus: []string{"Secured"},
			expectedCKN:    []string{"abc123def456"},
			expectDelete:   false,
		},
		{
			name: "Interface with MKA counters",
			notification: &gnmipb.Notification{
				Timestamp: 2000,
				Prefix: &gnmipb.Path{
					Origin: "junos",
					Target: "device2",
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "macsec"},
								{Name: "interfaces"},
								{Name: "interface", Key: map[string]string{"name": "et-0/0/1"}},
								{Name: "name"},
							},
						},
						Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "et-0/0/1"}},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "macsec"},
								{Name: "interfaces"},
								{Name: "interface", Key: map[string]string{"name": "et-0/0/1"}},
								{Name: "mka"},
								{Name: "state"},
								{Name: "jnx-cak-name"},
							},
						},
						Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "def789abc012"}},
					},
				},
			},
			expectedLen:    2,
			expectedStatus: []string{"Secured"},
			expectedCKN:    []string{"def789abc012"},
			expectDelete:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testImpl := &impl{state: newMacsecState()}

			sr := &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: test.notification,
				},
			}

			response, err := testImpl.translate(sr)
			if err != nil {
				t.Fatalf("translate() failed: %v", err)
			}

			if response == nil {
				t.Fatalf("expected non-nil response")
			}

			outNotif := response.GetUpdate()
			if outNotif == nil {
				t.Fatalf("expected non-nil notification in response")
			}

			if len(outNotif.GetUpdate()) != test.expectedLen {
				t.Errorf("translate() returned %d updates, want %d", len(outNotif.GetUpdate()), test.expectedLen)
			}

			// Verify origin is set to openconfig
			if outNotif.GetPrefix().GetOrigin() != "openconfig" {
				t.Errorf("translate() origin = %q, want %q", outNotif.GetPrefix().GetOrigin(), "openconfig")
			}

			// Verify status and CKN lists
			var gotStatus, gotCKN []string
			for _, mUpdate := range outNotif.GetUpdate() {
				lastElem := mUpdate.GetPath().GetElem()[len(mUpdate.GetPath().GetElem())-1].GetName()
				var vals []string
				for _, valElem := range mUpdate.GetVal().GetLeaflistVal().GetElement() {
					vals = append(vals, valElem.GetStringVal())
				}
				if lastElem == "status" {
					gotStatus = vals
				} else if lastElem == "ckn" {
					gotCKN = vals
				}
			}

			if diff := cmp.Diff(test.expectedStatus, gotStatus); diff != "" {
				t.Errorf("translate() status diff (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(test.expectedCKN, gotCKN); diff != "" {
				t.Errorf("translate() CKN diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDeleteHandler(t *testing.T) {
	tests := []struct {
		name          string
		notification  *gnmipb.Notification
		expectDeletes int
	}{
		{
			name: "Interface MACsec delete",
			notification: &gnmipb.Notification{
				Timestamp: 1000,
				Prefix: &gnmipb.Path{
					Origin: "junos",
					Target: "device1",
				},
				Delete: []*gnmipb.Path{
					{
						Elem: []*gnmipb.PathElem{
							{Name: "macsec"},
							{Name: "interfaces"},
							{Name: "interface", Key: map[string]string{"name": "et-0/0/0"}},
						},
					},
				},
			},
			expectDeletes: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testImpl := &impl{state: newMacsecState()}
			intfForDelete := testImpl.deleteHandler(test.notification)

			if len(intfForDelete) != test.expectDeletes {
				t.Errorf("deleteHandler() returned %d deletes, want %d", len(intfForDelete), test.expectDeletes)
			}
		})
	}
}

func TestExtractInterfaceAndSCI(t *testing.T) {
	tests := []struct {
		name          string
		path          *gnmipb.Path
		expIntf       string
		expCounter    string
		expectedError bool
	}{
		{
			name: "Interface state path",
			path: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{Name: "macsec"},
					{Name: "interfaces"},
					{Name: "interface", Key: map[string]string{"name": "et-0/0/0"}},
					{Name: "name"},
				},
			},
			expIntf:    "et-0/0/0",
			expCounter: "name",
		},
		{
			name: "MKA state path",
			path: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{Name: "macsec"},
					{Name: "interfaces"},
					{Name: "interface", Key: map[string]string{"name": "et-0/0/1"}},
					{Name: "mka"},
					{Name: "state"},
					{Name: "jnx-cak-name"},
				},
			},
			expIntf:    "et-0/0/1",
			expCounter: "jnx-cak-name",
		},
		{
			name: "Path with no interface key",
			path: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{Name: "macsec"},
					{Name: "interfaces"},
					{Name: "interface"},
					{Name: "name"},
				},
			},
			expectedError: true,
		},
		{
			name: "Interface element not preceded by interfaces parent",
			path: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{Name: "macsec"},
					{Name: "wrongparent"},
					{Name: "interface", Key: map[string]string{"name": "et-0/0/0"}},
					{Name: "name"},
				},
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			intfID, counterName, err := interfaceIDAndValue(test.path)

			if (err != nil) != test.expectedError {
				t.Errorf("interfaceIDAndValue() error = %v, wantErr %v", err, test.expectedError)
			}

			if intfID != test.expIntf {
				t.Errorf("interfaceIDAndValue() intfID = %q, want %q", intfID, test.expIntf)
			}

			if !test.expectedError && counterName != test.expCounter {
				t.Errorf("interfaceIDAndValue() counterName = %q, want %q", test.expCounter, counterName)
			}
		})
	}
}

func TestNew(t *testing.T) {
	ft := New()
	if ft == nil {
		t.Fatalf("New() returned nil")
	}
	if ft.ID() != ftconsts.JuniperMacsecStateFunctionalTranslator {
		t.Errorf("New() ID = %q, want %q", ft.ID(), ftconsts.JuniperMacsecStateFunctionalTranslator)
	}
}

func TestTranslate_Errors(t *testing.T) {
	tests := []struct {
		name        string
		input       *gnmipb.SubscribeResponse
		wantErr     bool
		wantNilResp bool
		check       func(*testing.T, *gnmipb.SubscribeResponse)
	}{
		{
			name: "nil update returns nil",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{Update: nil},
			},
			wantNilResp: true,
		},
		{
			name: "delete only notification returns deletes and non-nil response",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Prefix: &gnmipb.Path{Origin: "junos", Target: "device1"},
						Delete: []*gnmipb.Path{
							{Elem: []*gnmipb.PathElem{{Name: "macsec"}, {Name: "interfaces"}, {Name: "interface", Key: map[string]string{"name": "et-0/0/0"}}}},
						},
					},
				},
			},
			check: func(t *testing.T, res *gnmipb.SubscribeResponse) {
				if res.GetUpdate() == nil {
					t.Fatalf("expected non-nil response for delete-only notification")
				}
				if len(res.GetUpdate().GetDelete()) != 2 {
					t.Errorf("Translate() returned %d deletes, want 2", len(res.GetUpdate().GetDelete()))
				}
				if len(res.GetUpdate().GetUpdate()) != 0 {
					t.Errorf("Translate() returned %d updates, want 0", len(res.GetUpdate().GetUpdate()))
				}
			},
		},
		{
			name: "both update and delete notification returns both updates and deletes",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Prefix: &gnmipb.Path{Origin: "junos", Target: "device1"},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "macsec"}, {Name: "interfaces"}, {Name: "interface", Key: map[string]string{"name": "et-0/0/0"}}, {Name: "name"}}},
								Val:  &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "et-0/0/0"}},
							},
							{
								Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "macsec"}, {Name: "interfaces"}, {Name: "interface", Key: map[string]string{"name": "et-0/0/0"}}, {Name: "mka"}, {Name: "state"}, {Name: "jnx-cak-name"}}},
								Val:  &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "ckn1"}},
							},
						},
						Delete: []*gnmipb.Path{
							{Elem: []*gnmipb.PathElem{{Name: "macsec"}, {Name: "interfaces"}, {Name: "interface", Key: map[string]string{"name": "et-0/0/1"}}}},
						},
					},
				},
			},
			check: func(t *testing.T, res *gnmipb.SubscribeResponse) {
				if res.GetUpdate() == nil {
					t.Fatalf("expected non-nil response for mixed notification")
				}
				if len(res.GetUpdate().GetDelete()) != 2 {
					t.Errorf("Translate() returned %d deletes, want 2", len(res.GetUpdate().GetDelete()))
				}
				if len(res.GetUpdate().GetUpdate()) != 2 {
					t.Errorf("Translate() returned %d updates, want 2", len(res.GetUpdate().GetUpdate()))
				}
			},
		},
		{
			name: "update and delete on same interface returns only delete",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Prefix: &gnmipb.Path{Origin: "junos", Target: "device1"},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "macsec"}, {Name: "interfaces"}, {Name: "interface", Key: map[string]string{"name": "et-0/0/0"}}, {Name: "name"}}},
								Val:  &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "et-0/0/0"}},
							},
						},
						Delete: []*gnmipb.Path{
							{Elem: []*gnmipb.PathElem{{Name: "macsec"}, {Name: "interfaces"}, {Name: "interface", Key: map[string]string{"name": "et-0/0/0"}}}},
						},
					},
				},
			},
			check: func(t *testing.T, res *gnmipb.SubscribeResponse) {
				if res.GetUpdate() == nil {
					t.Fatalf("expected non-nil response for same-interface delete/update notification")
				}
				if len(res.GetUpdate().GetDelete()) != 2 {
					t.Errorf("Translate() returned %d deletes, want 2", len(res.GetUpdate().GetDelete()))
				}
				if len(res.GetUpdate().GetUpdate()) != 0 {
					t.Errorf("Translate() returned %d updates, want 0", len(res.GetUpdate().GetUpdate()))
				}
			},
		},
		{
			name: "path length less than 2 in interfaceIDAndValue",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Prefix: &gnmipb.Path{Origin: "junos", Target: "device1"},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "name"}}},
								Val:  &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "et-0/0/0"}},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "interface element not found",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Prefix: &gnmipb.Path{Origin: "junos", Target: "device1"},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "macsec"}, {Name: "interfaces"}, {Name: "name"}}},
								Val:  &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "et-0/0/0"}},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "interface element wrong parent container",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Prefix: &gnmipb.Path{Origin: "junos", Target: "device1"},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "macsec"}, {Name: "wrongcontainer"}, {Name: "interface", Key: map[string]string{"name": "et-0/0/0"}}, {Name: "name"}}},
								Val:  &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "et-0/0/0"}},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "delete path unmatched pattern",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Prefix: &gnmipb.Path{Origin: "junos", Target: "device1"},
						Delete: []*gnmipb.Path{{Elem: []*gnmipb.PathElem{{Name: "unmatched"}}}},
					},
				},
			},
			wantNilResp: true,
		},
		{
			name: "metadata empty updateElems",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Prefix: &gnmipb.Path{Origin: "junos", Target: "device1"},
						Update: []*gnmipb.Update{{Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{}}}},
					},
				},
			},
			wantNilResp: true,
		},
		{
			name: "metadata unknown leaf name",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Prefix: &gnmipb.Path{Origin: "junos", Target: "device1"},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "macsec"}, {Name: "interfaces"}, {Name: "interface", Key: map[string]string{"name": "et-0/0/0"}}, {Name: "unknown-leaf"}}},
							},
						},
					},
				},
			},
			wantNilResp: true,
		},
		{
			name: "cak name empty value on update",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Prefix: &gnmipb.Path{Origin: "junos", Target: "device1"},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "macsec"}, {Name: "interfaces"}, {Name: "interface", Key: map[string]string{"name": "et-0/0/0"}}, {Name: "mka"}, {Name: "state"}, {Name: "jnx-cak-name"}}},
								Val:  &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: ""}},
							},
						},
					},
				},
			},
			wantNilResp: true,
		},
		{
			name: "metadata unmatched pattern path",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Prefix: &gnmipb.Path{Origin: "junos", Target: "device1"},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "macsec"}, {Name: "interfaces"}, {Name: "interface", Key: map[string]string{"name": "et-0/0/0"}}, {Name: "unmatched-section"}, {Name: "name"}}},
								Val:  &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "et-0/0/0"}},
							},
						},
					},
				},
			},
			wantNilResp: true,
		},
		{
			name:        "nil subscribe response",
			input:       nil,
			wantNilResp: true,
		},
		{
			name: "sync response returns nil",
			input: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_SyncResponse{
					SyncResponse: true,
				},
			},
			wantNilResp: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testImpl := &impl{state: newMacsecState()}
			res, err := testImpl.translate(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("translate() error = %v, wantErr %v", err, tc.wantErr)
			}
			if tc.wantNilResp && res != nil {
				t.Errorf("translate() returned %v, want nil", res)
			}
			if tc.check != nil {
				tc.check(t, res)
			}
		})
	}
}

func TestTranslateMACSecState_EdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(state *macsecState)
		intf       string
		wantSkip   bool
		wantStatus []string
	}{
		{
			name:     "interface not found in targetInfo",
			setup:    func(state *macsecState) {},
			intf:     "et-0/0/1",
			wantSkip: true,
		},
		{
			name: "cpStatusSet is false",
			setup: func(state *macsecState) {
				ifaceInfo := state.createOrGetInterface("et-0/0/1")
				ckn1 := ifaceInfo.createOrGetCKN("ckn1")
				ckn1.success = true
				ckn1.successSet = true
				ckn1.principal = true
				ckn1.principalSet = true
			},
			intf:     "et-0/0/1",
			wantSkip: true,
		},
		{
			name: "CloneStatuses is empty",
			setup: func(state *macsecState) {
				ifaceInfo := state.createOrGetInterface("et-0/0/1")
				ifaceInfo.cpStatus = true
				ifaceInfo.cpStatusSet = true
				for k := range ifaceInfo.cknStatuses {
					delete(ifaceInfo.cknStatuses, k)
				}
			},
			intf:     "et-0/0/1",
			wantSkip: true,
		},
		{
			name: "CKN data is incomplete",
			setup: func(state *macsecState) {
				ifaceInfo := state.createOrGetInterface("et-0/0/1")
				ifaceInfo.cpStatus = true
				ifaceInfo.cpStatusSet = true
				ifaceInfo.createOrGetCKN("ckn1")
			},
			intf:     "et-0/0/1",
			wantSkip: true,
		},
		{
			name: "Unknown case non-matching switch branch",
			setup: func(state *macsecState) {
				ifaceInfo := state.createOrGetInterface("et-0/0/1")
				ifaceInfo.cpStatus = false
				ifaceInfo.cpStatusSet = true
				ckn1 := ifaceInfo.createOrGetCKN("ckn1")
				ckn1.success = true
				ckn1.successSet = true
				ckn1.principal = false
				ckn1.principalSet = true
			},
			intf:       "et-0/0/1",
			wantSkip:   false,
			wantStatus: []string{"Unknown"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			state := newMacsecState()
			testImpl := &impl{state: state}
			if tc.setup != nil {
				tc.setup(state)
			}
			statuses, _, skip := testImpl.translateMACSecState(tc.intf)
			if skip != tc.wantSkip {
				t.Fatalf("translateMACSecState() skip = %v, want %v", skip, tc.wantSkip)
			}
			if !skip && tc.wantStatus != nil {
				if diff := cmp.Diff(tc.wantStatus, statuses); diff != "" {
					t.Errorf("translateMACSecState() status diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestInterfaceMacSecInfoCreateOrGetCKNNilMap(t *testing.T) {
	info := &interfaceMacSecInfo{interfaceName: "et-0/0/1"}
	ckn := info.createOrGetCKN("ckn1")
	if ckn == nil {
		t.Fatal("createOrGetCKN() returned nil")
	}
	if got := info.createOrGetCKN("ckn1"); got != ckn {
		t.Errorf("createOrGetCKN() returned %v, want %v", got, ckn)
	}
}

func TestCKNInfoIsCompleteNil(t *testing.T) {
	var c *cknInfo
	if c.isComplete() {
		t.Errorf("nil cknInfo.isComplete() = true, want false")
	}
}

func TestInterfaceIDAndValue_MissingInterfaceNameKey(t *testing.T) {
	path := &gnmipb.Path{
		Elem: []*gnmipb.PathElem{
			{Name: "macsec"},
			{Name: "interfaces"},
			{Name: "interface", Key: map[string]string{}},
			{Name: "name"},
		},
	}
	_, _, err := interfaceIDAndValue(path)
	if err == nil {
		t.Errorf("interfaceIDAndValue() expected error for missing name key, got nil")
	}
}

func TestExtractDeleteInfo_InvalidInterface(t *testing.T) {
	path := &gnmipb.Path{
		Origin: "junos",
		Elem: []*gnmipb.PathElem{
			{Name: "macsec"},
			{Name: "interfaces"},
			{Name: "interface", Key: map[string]string{}},
		},
	}
	_, err := extractDeleteInfo(path)
	if err == nil {
		t.Errorf("extractDeleteInfo() expected error for invalid interface, got nil")
	}
}

func TestIsKnownLeaf(t *testing.T) {
	if !isKnownLeaf("name") {
		t.Errorf("isKnownLeaf(\"name\") = false, want true")
	}
	if !isKnownLeaf("jnx-cak-name") {
		t.Errorf("isKnownLeaf(\"jnx-cak-name\") = false, want true")
	}
	if isKnownLeaf("unknown") {
		t.Errorf("isKnownLeaf(\"unknown\") = true, want false")
	}
}
