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

package aristamplslfib

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"github.com/openconfig/functional-translators/ftconsts"
	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

func TestNew_MetadataAndMap(t *testing.T) {
	ft := New()
	if got, want := ft.ID(), ftconsts.AristaMplsLfibFunctionalTranslator; got != want {
		t.Errorf("ft.ID() = %q, want %q", got, want)
	}

	metadata := ft.Metadata()
	if len(metadata) != 1 {
		t.Fatalf("expected 1 metadata entry, got %d", len(metadata))
	}
	if got, want := metadata[0].Vendor, ftconsts.VendorArista; got != want {
		t.Errorf("metadata[0].Vendor = %q, want %q", got, want)
	}
	if got, want := metadata[0].SoftwareVersionRange.InclusiveMin, ""; got != want {
		t.Errorf("metadata[0].SoftwareVersionRange.InclusiveMin = %q, want empty", got)
	}
	if got, want := metadata[0].SoftwareVersionRange.ExclusiveMax, "4.38.0"; got != want {
		t.Errorf("metadata[0].SoftwareVersionRange.ExclusiveMax = %q, want %q", got, want)
	}

	inputMap := ft.OutputToInputMap()
	wantOCPackets := "/openconfig/network-instances/network-instance/afts/mpls/label-entry/state/counters/packets-forwarded"
	wantOCOctets := "/openconfig/network-instances/network-instance/afts/mpls/label-entry/state/counters/octets-forwarded"

	if _, ok := inputMap[wantOCPackets]; !ok {
		t.Errorf("missing %q in OutputToInputMap", wantOCPackets)
	}
	if _, ok := inputMap[wantOCOctets]; !ok {
		t.Errorf("missing %q in OutputToInputMap", wantOCOctets)
	}
}

func TestTranslate_Success(t *testing.T) {
	sr := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 1726680000,
				Prefix: &gnmipb.Path{
					Origin: "eos_native",
					Target: "pf93.sql14",
				},
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "Smash"},
								{Name: "flexCounters"},
								{Name: "counterTable"},
								{Name: "MplsLfib"},
								{Name: "4294967295"},
								{Name: "counterEntry"},
								{Name: "80978_1048576"},
								{Name: "pkts"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{UintVal: 12859},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "Smash"},
								{Name: "flexCounters"},
								{Name: "counterTable"},
								{Name: "MplsLfib"},
								{Name: "4294967295"},
								{Name: "counterEntry"},
								{Name: "80978_1048576"},
								{Name: "octets"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{UintVal: 1024000},
						},
					},
					{
						// Test positive IntVal conversion and 0 value.
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "Smash"},
								{Name: "flexCounters"},
								{Name: "counterTable"},
								{Name: "MplsLfib"},
								{Name: "4294967295"},
								{Name: "counterEntry"},
								{Name: "1002_1048576"},
								{Name: "pkts"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_IntVal{IntVal: 550},
						},
					},
					{
						// Test boundary max label and max uint64 value.
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{Name: "Smash"},
								{Name: "flexCounters"},
								{Name: "counterTable"},
								{Name: "MplsLfib"},
								{Name: "4294967295"},
								{Name: "counterEntry"},
								{Name: "1048575_1048576"},
								{Name: "octets"},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{UintVal: 18446744073709551615},
						},
					},
				},
			},
		},
	}

	want := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 1726680000,
				Prefix: &gnmipb.Path{
					Origin: "openconfig",
					Target: "pf93.sql14",
					Elem: []*gnmipb.PathElem{
						{Name: "network-instances"},
						{
							Name: "network-instance",
							Key: map[string]string{
								"name": "default",
							},
						},
						{Name: "afts"},
						{Name: "mpls"},
					},
				},
				Update: []*gnmipb.Update{
					{
						Path: returnOCPath("80978", ocPacketsForwarded),
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{UintVal: 12859},
						},
					},
					{
						Path: returnOCPath("80978", ocOctetsForwarded),
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{UintVal: 1024000},
						},
					},
					{
						Path: returnOCPath("1002", ocPacketsForwarded),
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{UintVal: 550},
						},
					},
					{
						Path: returnOCPath("1048575", ocOctetsForwarded),
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_UintVal{UintVal: 18446744073709551615},
						},
					},
				},
			},
		},
	}

	got, err := translate(sr)
	if err != nil {
		t.Fatalf("translate() returned unexpected error: %v", err)
	}

	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("translate() returned diff (-want +got):\n%s", diff)
	}
}

func TestTranslate_Deletes(t *testing.T) {
	sr := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 1726680000,
				Prefix: &gnmipb.Path{
					Origin: "eos_native",
					Target: "pf93.sql14",
				},
				Delete: []*gnmipb.Path{
					{
						Elem: []*gnmipb.PathElem{
							{Name: "Smash"},
							{Name: "flexCounters"},
							{Name: "counterTable"},
							{Name: "MplsLfib"},
							{Name: "4294967295"},
							{Name: "counterEntry"},
							{Name: "80978_1048576"},
						},
					},
				},
			},
		},
	}

	want := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: 1726680000,
				Prefix: &gnmipb.Path{
					Origin: "openconfig",
					Target: "pf93.sql14",
					Elem: []*gnmipb.PathElem{
						{Name: "network-instances"},
						{
							Name: "network-instance",
							Key: map[string]string{
								"name": "default",
							},
						},
						{Name: "afts"},
						{Name: "mpls"},
					},
				},
				Delete: []*gnmipb.Path{
					returnOCPath("80978", ocPacketsForwarded),
					returnOCPath("80978", ocOctetsForwarded),
				},
			},
		},
	}

	got, err := translate(sr)
	if err != nil {
		t.Fatalf("translate() returned unexpected error: %v", err)
	}

	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("translate() returned diff (-want +got):\n%s", diff)
	}
}

func TestTranslate_MalformedKeys(t *testing.T) {
	tests := []struct {
		name string
		key  string
		leaf string
	}{
		{
			name: "missing_underscore",
			key:  "80978",
			leaf: "pkts",
		},
		{
			name: "too_many_tokens",
			key:  "80978_1048576_extra",
			leaf: "pkts",
		},
		{
			name: "non_numeric_label",
			key:  "invalid_1048576",
			leaf: "pkts",
		},
		{
			name: "label_exceeds_20_bits_boundary",
			key:  "1048576_1048576",
			leaf: "pkts",
		},
		{
			name: "label_exceeds_20_bits",
			key:  "2000000_1048576",
			leaf: "pkts",
		},
		{
			name: "non_numeric_bos_marker",
			key:  "80978_invalid",
			leaf: "pkts",
		},
		{
			name: "two_label_entry_non_sentinel_label1",
			key:  "80978_100",
			leaf: "pkts",
		},
		{
			name: "unrecognized_leaf",
			key:  "80978_1048576",
			leaf: "invalidLeaf",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var deletes []*gnmipb.Path
			if tc.leaf == "pkts" {
				deletes = []*gnmipb.Path{
					{
						Elem: []*gnmipb.PathElem{
							{Name: "Smash"},
							{Name: "flexCounters"},
							{Name: "counterTable"},
							{Name: "MplsLfib"},
							{Name: "4294967295"},
							{Name: "counterEntry"},
							{Name: tc.key},
						},
					},
				}
			}
			sr := &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1726680000,
						Prefix: &gnmipb.Path{
							Origin: "eos_native",
							Target: "pf93.sql14",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "Smash"},
										{Name: "flexCounters"},
										{Name: "counterTable"},
										{Name: "MplsLfib"},
										{Name: "4294967295"},
										{Name: "counterEntry"},
										{Name: tc.key},
										{Name: tc.leaf},
									},
								},
								Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 100}},
							},
						},
						Delete: deletes,
					},
				},
			}

			got, err := translate(sr)
			if err != nil {
				t.Fatalf("translate() returned error: %v", err)
			}
			if got != nil {
				t.Errorf("translate() expected nil for malformed key %s, got: %v", tc.name, got)
			}
		})
	}
}

func TestTranslate_InvalidValues(t *testing.T) {
	tests := []struct {
		name string
		val  *gnmipb.TypedValue
	}{
		{
			name: "nil_typed_value",
			val:  nil,
		},
		{
			name: "negative_int_val",
			val:  &gnmipb.TypedValue{Value: &gnmipb.TypedValue_IntVal{IntVal: -50}},
		},
		{
			name: "string_typed_value",
			val:  &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "not_a_num"}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sr := &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Timestamp: 1726680000,
						Prefix: &gnmipb.Path{
							Origin: "eos_native",
							Target: "pf93.sql14",
						},
						Update: []*gnmipb.Update{
							{
								Path: &gnmipb.Path{
									Elem: []*gnmipb.PathElem{
										{Name: "Smash"},
										{Name: "flexCounters"},
										{Name: "counterTable"},
										{Name: "MplsLfib"},
										{Name: "4294967295"},
										{Name: "counterEntry"},
										{Name: "80978_1048576"},
										{Name: "pkts"},
									},
								},
								Val: tc.val,
							},
						},
					},
				},
			}

			got, err := translate(sr)
			if err != nil {
				t.Fatalf("translate() returned error: %v", err)
			}
			if got != nil {
				t.Errorf("translate() expected nil for invalid value %s, got: %v", tc.name, got)
			}
		})
	}
}

func TestTranslate_NilAndNonMatchingInputs(t *testing.T) {
	t.Run("nil_subscribe_response", func(t *testing.T) {
		got, err := translate(nil)
		if err != nil || got != nil {
			t.Errorf("translate(nil) = (%v, %v), want (nil, nil)", got, err)
		}
	})

	t.Run("nil_notification", func(t *testing.T) {
		sr := &gnmipb.SubscribeResponse{Response: &gnmipb.SubscribeResponse_Update{Update: nil}}
		got, err := translate(sr)
		if err != nil || got != nil {
			t.Errorf("translate(nil notification) = (%v, %v), want (nil, nil)", got, err)
		}
	})

	t.Run("empty_notification", func(t *testing.T) {
		sr := &gnmipb.SubscribeResponse{
			Response: &gnmipb.SubscribeResponse_Update{
				Update: &gnmipb.Notification{
					Prefix: &gnmipb.Path{Origin: "eos_native", Target: "pf93.sql14"},
				},
			},
		}
		got, err := translate(sr)
		if err != nil || got != nil {
			t.Errorf("translate(empty notification) = (%v, %v), want (nil, nil)", got, err)
		}
	})

	t.Run("non_matching_path", func(t *testing.T) {
		sr := &gnmipb.SubscribeResponse{
			Response: &gnmipb.SubscribeResponse_Update{
				Update: &gnmipb.Notification{
					Prefix: &gnmipb.Path{Origin: "eos_native", Target: "pf93.sql14"},
					Update: []*gnmipb.Update{
						{
							Path: &gnmipb.Path{
								Elem: []*gnmipb.PathElem{
									{Name: "interfaces"},
									{Name: "interface"},
								},
							},
							Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 100}},
						},
						nil, // Verify nil update is skipped safely
					},
					Delete: []*gnmipb.Path{
						{
							Elem: []*gnmipb.PathElem{
								{Name: "interfaces"},
								{Name: "interface"},
							},
						},
						nil, // Verify nil delete is skipped safely
					},
				},
			},
		}
		got, err := translate(sr)
		if err != nil || got != nil {
			t.Errorf("translate(non matching path) = (%v, %v), want (nil, nil)", got, err)
		}
	})

	t.Run("non_matching_similar_smash_path", func(t *testing.T) {
		// Tests an 8-element Smash update path and 7-element Smash delete path from a different
		// counter table (Ipv4Fib instead of MplsLfib). The key and leaf tokens are well-formed,
		// so this specifically tests that ftutilities.MatchPath rejects non-matching tables and
		// prevents translation mutants.
		sr := &gnmipb.SubscribeResponse{
			Response: &gnmipb.SubscribeResponse_Update{
				Update: &gnmipb.Notification{
					Prefix: &gnmipb.Path{Origin: "eos_native", Target: "pf93.sql14"},
					Update: []*gnmipb.Update{
						{
							Path: &gnmipb.Path{
								Elem: []*gnmipb.PathElem{
									{Name: "Smash"},
									{Name: "flexCounters"},
									{Name: "counterTable"},
									{Name: "Ipv4Fib"},
									{Name: "4294967295"},
									{Name: "counterEntry"},
									{Name: "80978_1048576"},
									{Name: "pkts"},
								},
							},
							Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 100}},
						},
					},
					Delete: []*gnmipb.Path{
						{
							Elem: []*gnmipb.PathElem{
								{Name: "Smash"},
								{Name: "flexCounters"},
								{Name: "counterTable"},
								{Name: "Ipv4Fib"},
								{Name: "4294967295"},
								{Name: "counterEntry"},
								{Name: "80978_1048576"},
							},
						},
					},
				},
			},
		}
		got, err := translate(sr)
		if err != nil || got != nil {
			t.Errorf("translate(similar non-matching smash path) = (%v, %v), want (nil, nil)", got, err)
		}
	})
}
