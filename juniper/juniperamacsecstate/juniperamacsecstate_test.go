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

package juniperamacsecstate

import (
	"testing"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

func TestTranslateInterfaceState(t *testing.T) {
	tests := []struct {
		name         string
		notification *gnmipb.Notification
		expectedLen  int
		expectDelete bool
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
								{Name: "interfaces"},
								{Name: "interface"},
								{Name: "et-0/0/0"},
								{Name: "macsec"},
								{Name: "state"},
								{Name: "enable"},
							},
						},
						Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_BoolVal{BoolVal: true}},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "interfaces"},
								{Name: "interface"},
								{Name: "et-0/0/0"},
								{Name: "macsec"},
								{Name: "mka"},
								{Name: "state"},
								{Name: "counters"},
								{Name: "in-mkpdu"},
							},
						},
						Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 10}},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "interfaces"},
								{Name: "interface"},
								{Name: "et-0/0/0"},
								{Name: "macsec"},
								{Name: "mka"},
								{Name: "state"},
								{Name: "counters"},
								{Name: "in-sak-mkpdu"},
							},
						},
						Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 5}},
					},
				},
			},
			expectedLen:  2,
			expectDelete: false,
		},
		{
			name: "Interface with MKA counters",
			notification: &gnmipb.Notification{
				Timestamp: 2000,
				Prefix: &gnmipb.Path{
					Origin: "junos",
					Target: "device1",
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "interfaces"},
								{Name: "interface"},
								{Name: "et-0/0/1"},
								{Name: "macsec"},
								{Name: "state"},
								{Name: "enable"},
							},
						},
						Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_BoolVal{BoolVal: false}},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "interfaces"},
								{Name: "interface"},
								{Name: "et-0/0/1"},
								{Name: "macsec"},
								{Name: "mka"},
								{Name: "state"},
								{Name: "counters"},
								{Name: "in-mkpdu"},
							},
						},
						Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 0}},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "interfaces"},
								{Name: "interface"},
								{Name: "et-0/0/1"},
								{Name: "macsec"},
								{Name: "mka"},
								{Name: "state"},
								{Name: "counters"},
								{Name: "in-sak-mkpdu"},
							},
						},
						Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 0}},
					},
				},
			},
			expectedLen:  2,
			expectDelete: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sr := &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: test.notification,
				},
			}

			response, err := translate(sr)
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
				t.Errorf("expected %d updates, got %d", test.expectedLen, len(outNotif.GetUpdate()))
			}

			// Verify origin is set to openconfig
			if outNotif.GetPrefix().GetOrigin() != "openconfig" {
				t.Errorf("expected origin 'openconfig', got '%s'", outNotif.GetPrefix().GetOrigin())
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
							{Name: "interfaces"},
							{Name: "interface"},
							{Name: "et-0/0/0"},
							{Name: "macsec"},
						},
					},
				},
			},
			expectDeletes: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			intfForDelete := deleteHandler(test.notification)

			if len(intfForDelete) != test.expectDeletes {
				t.Errorf("expected %d deletes, got %d", test.expectDeletes, len(intfForDelete))
			}
		})
	}
}

func TestExtractInterfaceAndSCI(t *testing.T) {
	tests := []struct {
		name          string
		path          *gnmipb.Path
		expIntf       string
		expSCI        string
		expIsScsa     bool
		expectedError bool
	}{
		{
			name: "Interface state path",
			path: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{Name: "interfaces"},
					{Name: "interface"},
					{Name: "et-0/0/0"},
					{Name: "macsec"},
					{Name: "state"},
					{Name: "enable"},
				},
			},
			expIntf:   "et-0/0/0",
			expSCI:    "",
			expIsScsa: false,
		},
		{
			name: "SCSA RX state path",
			path: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{Name: "interfaces"},
					{Name: "interface"},
					{Name: "et-0/0/1"},
					{Name: "macsec"},
					{Name: "scsa-rx"},
					{Name: "state"},
					{Name: "0011223344556677"},
					{Name: "sci-rx"},
				},
			},
			expIntf:   "et-0/0/1",
			expSCI:    "0011223344556677",
			expIsScsa: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			intfID, counterName, err := interfaceIDAndValue(test.path)

			if (err != nil) != test.expectedError {
				t.Errorf("unexpected error state: got %v, expected error: %v", err, test.expectedError)
			}

			if intfID != test.expIntf {
				t.Errorf("expected interface '%s', got '%s'", test.expIntf, intfID)
			}

			// Verify counterName is not empty for valid paths
			if !test.expectedError && counterName == "" {
				t.Errorf("expected non-empty counterName for valid path")
			}
		})
	}
}
