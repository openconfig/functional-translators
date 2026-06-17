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

package juniperamacseccounters

import (
	"testing"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

func TestTranslateInterfaceCounters(t *testing.T) {
	tests := []struct {
		name         string
		notification *gnmipb.Notification
		expectedLen  int
	}{
		{
			name: "TX untagged packets counter",
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
								{Name: "counters"},
								{Name: "tx-untagged-packets"},
							},
						},
						Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 1000}},
					},
				},
			},
			expectedLen: 1,
		},
		{
			name: "RX bad tag packets counter",
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
								{Name: "et-0/0/0"},
								{Name: "macsec"},
								{Name: "state"},
								{Name: "counters"},
								{Name: "rx-badtag-packets"},
							},
						},
						Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 50}},
					},
				},
			},
			expectedLen: 1,
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
			name: "Interface level counter path",
			path: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{Name: "interfaces"},
					{Name: "interface"},
					{Name: "et-0/0/0"},
					{Name: "macsec"},
					{Name: "state"},
					{Name: "counters"},
					{Name: "tx-untagged-packets"},
				},
			},
			expIntf:   "et-0/0/0",
			expSCI:    "",
			expIsScsa: false,
		},
		{
			name: "SCSA TX path",
			path: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{Name: "interfaces"},
					{Name: "interface"},
					{Name: "et-0/0/1"},
					{Name: "macsec"},
					{Name: "scsa-tx"},
					{Name: "state"},
					{Name: "counters"},
					{Name: "0011223344556677"},
					{Name: "sc-auth-only-packets"},
				},
			},
			expIntf:   "et-0/0/1",
			expSCI:    "0011223344556677",
			expIsScsa: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			intfName, sci, isScsa, err := extractInterfaceAndSCI(test.path)

			if (err != nil) != test.expectedError {
				t.Errorf("unexpected error state: got %v, expected error: %v", err, test.expectedError)
			}

			if intfName != test.expIntf {
				t.Errorf("expected interface '%s', got '%s'", test.expIntf, intfName)
			}

			if sci != test.expSCI {
				t.Errorf("expected SCI '%s', got '%s'", test.expSCI, sci)
			}

			if isScsa != test.expIsScsa {
				t.Errorf("expected isScsa %v, got %v", test.expIsScsa, isScsa)
			}
		})
	}
}
