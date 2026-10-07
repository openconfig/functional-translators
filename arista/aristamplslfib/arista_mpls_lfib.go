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

// Package aristamplslfib implements the functional translator for Arista MPLS LFIB counters.
//
// While Arista EOS supports the standard OpenConfig AFT MPLS path, it does not populate all
// counter information required by downstream telemetry systems (specifically per-label forwarding
// counters for static VRF decap labels prior to EOS 4.38.0; see b/570223062). Instead, these
// hardware ASIC counters are exposed via EOS Native Smash:
//
//	/eos_native/Smash/flexCounters/counterTable/MplsLfib/4294967295/counterEntry/<key>/{pkts, octets}
//
// where <key> is formatted as "<label>_1048576" (with 1048576 being the Bottom-of-Stack marker).
// This translator converts these native Smash notifications into standard OpenConfig AFT:
//
//	/network-instances/network-instance[name=default]/afts/mpls/label-entry[label=<label>]/state/counters/{packets-forwarded, octets-forwarded}
package aristamplslfib

import (
	"strconv"
	"strings"

	log "github.com/golang/glog"
	"github.com/openconfig/functional-translators/ftconsts"
	"github.com/openconfig/functional-translators/ftutilities"
	"github.com/openconfig/functional-translators/translator"
	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

const (
	// Output OpenConfig leaf names according to openconfig-aft.yang.
	ocPacketsForwarded = "packets-forwarded"
	ocOctetsForwarded  = "octets-forwarded"

	// Input Arista EOS Native Smash leaf names.
	nativePkts   = "pkts"
	nativeOctets = "octets"

	// nativeCounterTablePath is the base Smash path subscribed to on Arista EOS.
	// 4294967295 (0xFFFFFFFF) is a sentinel value representing the aggregate
	// counter table for all MPLS labels during LFIB lookups.
	nativeCounterTablePath = "/eos_native/Smash/flexCounters/counterTable/MplsLfib/4294967295/counterEntry"

	// Maximum valid 20-bit MPLS label value ((1 << 20) - 1).
	maxValidMPLSLabel = 1048575

	// singleLabelSentinel is the Smash MplsLfib label1 value (1 << 20, i.e., maxValidMPLSLabel + 1)
	// indicating a single-label LFIB lookup entry (no second MPLS label in the lookup tuple).
	singleLabelSentinel = "1048576"

	// Supported EOS software version range for Smash flexCounters MPLS LFIB table.
	// TODO(team): Verify native OpenConfig gNMI counter support for MPLS
	// static vrf-label decap entries once EOS 4.38.0 (W release) is available.
	maxEOSVersion = "4.38.0"
)

var (
	translateMap = map[string][]string{
		"/openconfig/network-instances/network-instance/afts/mpls/label-entry/state/counters/packets-forwarded": {
			nativeCounterTablePath,
		},
		"/openconfig/network-instances/network-instance/afts/mpls/label-entry/state/counters/octets-forwarded": {
			nativeCounterTablePath,
		},
	}

	smashLfibUpdatePattern = &gnmipb.Path{
		Origin: "eos_native",
		Elem: []*gnmipb.PathElem{
			{Name: "Smash"},
			{Name: "flexCounters"},
			{Name: "counterTable"},
			{Name: "MplsLfib"},
			{Name: "4294967295"},
			{Name: "counterEntry"},
			{Name: "*"}, // <key>: e.g., 80978_1048576
			{Name: "*"}, // <leaf>: pkts or octets
		},
	}

	smashLfibDeletePattern = &gnmipb.Path{
		Origin: "eos_native",
		Elem: []*gnmipb.PathElem{
			{Name: "Smash"},
			{Name: "flexCounters"},
			{Name: "counterTable"},
			{Name: "MplsLfib"},
			{Name: "4294967295"},
			{Name: "counterEntry"},
			{Name: "*"}, // <key>: e.g., 80978_1048576
		},
	}

	vendorToOCLeaf = map[string]string{
		nativePkts:   ocPacketsForwarded,
		nativeOctets: ocOctetsForwarded,
	}
)

// New creates and registers the Arista MPLS LFIB functional translator.
func New() *translator.FunctionalTranslator {
	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               ftconsts.AristaMplsLfibFunctionalTranslator,
			Translate:        translate,
			OutputToInputMap: ftutilities.MustStringMapPaths(translateMap),
			Metadata: []*translator.FTMetadata{
				{
					Vendor: ftconsts.VendorArista,
					SoftwareVersionRange: &translator.SWRange{
						// TODO(team): Native OpenConfig support for MPLS static
						// vrf-label counters is expected in EOS 4.38.0.
						ExclusiveMax: maxEOSVersion,
					},
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("Failed to create Arista MPLS LFIB functional translator: %v", err)
	}
	return ft
}

// parseKey extracts and validates the MPLS label from the Smash counterEntry key.
// The Smash key represents a (label0, label1) lookup tuple formatted as "<label0>_<label1>".
// Only single-label entries where label1 equals singleLabelSentinel ("1048576") are matched.
func parseKey(key string) (string, bool) {
	labelStr, bosStr, found := strings.Cut(key, "_")
	if !found || bosStr != singleLabelSentinel {
		return "", false
	}

	label64, err := strconv.ParseUint(labelStr, 10, 32)
	if err != nil || label64 > maxValidMPLSLabel {
		return "", false
	}

	return labelStr, true
}

// extractUint64 safely extracts an unsigned 64-bit integer from a gNMI TypedValue.
func extractUint64(tv *gnmipb.TypedValue) (uint64, bool) {
	if tv == nil {
		return 0, false
	}
	switch v := tv.GetValue().(type) {
	case *gnmipb.TypedValue_UintVal:
		return v.UintVal, true
	case *gnmipb.TypedValue_IntVal:
		if v.IntVal < 0 {
			return 0, false
		}
		return uint64(v.IntVal), true
	default:
		return 0, false
	}
}

// returnOCPath constructs the OpenConfig path for a label counter relative to the
// notification prefix (/network-instances/network-instance[name=default]/afts/mpls).
func returnOCPath(label string, ocLeaf string) *gnmipb.Path {
	return &gnmipb.Path{
		Elem: []*gnmipb.PathElem{
			{
				Name: "label-entry",
				Key: map[string]string{
					"label": label,
				},
			},
			{Name: "state"},
			{Name: "counters"},
			{Name: ocLeaf},
		},
	}
}

// updateHandler processes incoming Smash counter updates and translates them to OpenConfig.
func updateHandler(n *gnmipb.Notification) ([]*gnmipb.Update, error) {
	if n == nil {
		return nil, nil
	}

	prefix := n.GetPrefix()
	var outgoingUpdates []*gnmipb.Update

	for _, u := range n.GetUpdate() {
		if u == nil {
			continue
		}
		fullPath := ftutilities.Join(prefix, u.GetPath())
		if !ftutilities.MatchPath(fullPath, smashLfibUpdatePattern) {
			continue
		}

		elems := fullPath.GetElem()
		key := elems[len(elems)-2].GetName()
		label, ok := parseKey(key)
		if !ok {
			log.V(2).Infof("Skipping malformed Smash MPLS LFIB key %q in path %v", key, fullPath)
			continue
		}

		vendorLeaf := elems[len(elems)-1].GetName()
		ocLeaf, found := vendorToOCLeaf[vendorLeaf]
		if !found {
			continue
		}

		val, ok := extractUint64(u.GetVal())
		if !ok {
			log.V(2).Infof("Skipping non-integer value for Smash MPLS LFIB path %v: %v", fullPath, u.GetVal())
			continue
		}

		outgoingUpdates = append(outgoingUpdates, &gnmipb.Update{
			Path: returnOCPath(label, ocLeaf),
			Val: &gnmipb.TypedValue{
				Value: &gnmipb.TypedValue_UintVal{
					UintVal: val,
				},
			},
		})
	}

	return outgoingUpdates, nil
}

// deleteHandler processes incoming 7-element Smash counterEntry deletes and translates them to
// OpenConfig counter leaf deletes. Note that in SAMPLE mode, stale labels are evicted via implicit
// deletion; this handler supports entry-level deletes if subscribed in ON_CHANGE mode.
func deleteHandler(n *gnmipb.Notification) []*gnmipb.Path {
	prefix := n.GetPrefix()
	var outgoingDeletes []*gnmipb.Path

	for _, d := range n.GetDelete() {
		if d == nil {
			continue
		}
		fullPath := ftutilities.Join(prefix, d)
		if !ftutilities.MatchPath(fullPath, smashLfibDeletePattern) {
			continue
		}

		elems := fullPath.GetElem()
		key := elems[len(elems)-1].GetName()
		label, ok := parseKey(key)
		if !ok {
			continue
		}

		outgoingDeletes = append(
			outgoingDeletes,
			returnOCPath(label, ocPacketsForwarded),
			returnOCPath(label, ocOctetsForwarded),
		)
	}

	return outgoingDeletes
}

func translate(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
	if sr == nil {
		return nil, nil
	}

	notification := sr.GetUpdate()
	if notification == nil {
		return nil, nil
	}

	updates, err := updateHandler(notification)
	if err != nil {
		return nil, err
	}
	deletes := deleteHandler(notification)

	if len(updates) == 0 && len(deletes) == 0 {
		return nil, nil
	}

	target := notification.GetPrefix().GetTarget()
	return &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: notification.GetTimestamp(),
				Prefix: &gnmipb.Path{
					Origin: "openconfig",
					Target: target,
					Elem: []*gnmipb.PathElem{
						{Name: "network-instances"},
						{
							Name: "network-instance",
							// MPLS LFIB is a device-global forwarding table residing in the default
							// network instance (even for labels that decapsulate into non-default VRFs).
							// Arista EOS uses "default" as the canonical default network-instance name.
							Key: map[string]string{
								"name": "default",
							},
						},
						{Name: "afts"},
						{Name: "mpls"},
					},
				},
				Update: updates,
				Delete: deletes,
			},
		},
	}, nil
}
