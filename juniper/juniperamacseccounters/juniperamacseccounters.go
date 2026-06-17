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

// Package juniperamacseccounters translates MACsec interface counters and SC/SA statistics from Juniper native to OpenConfig.
package juniperamacseccounters

import (
	"encoding/json"
	"fmt"

	log "github.com/golang/glog"
	"github.com/openconfig/functional-translators/ftconsts"
	"github.com/openconfig/functional-translators/ftutilities"
	"github.com/openconfig/functional-translators/translator"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

var (
	// Juniper native paths to OpenConfig MACsec interface counters translation map
	translateMap = map[string][]string{
		// Interface level MACsec counters
		"/openconfig/macsec/interfaces/interface/state/counters/tx-untagged-pkts": {
			"/junos/interfaces/interface/macsec/state/counters/tx-untagged-packets",
		},
		"/openconfig/macsec/interfaces/interface/state/counters/rx-untagged-pkts": {
			"/junos/interfaces/interface/macsec/state/counters/rx-untagged-packets",
		},
		"/openconfig/macsec/interfaces/interface/state/counters/rx-badtag-pkts": {
			"/junos/interfaces/interface/macsec/state/counters/rx-badtag-packets",
		},
		"/openconfig/macsec/interfaces/interface/state/counters/rx-unknownsci-pkts": {
			"/junos/interfaces/interface/macsec/state/counters/rx-unknownsci-packets",
		},
		"/openconfig/macsec/interfaces/interface/state/counters/rx-nosci-pkts": {
			"/junos/interfaces/interface/macsec/state/counters/rx-nosci-packets",
		},
		// Secure Channel TX statistics
		"/openconfig/macsec/interfaces/interface/scsa-tx/scsa-tx/state/counters/sc-auth-only": {
			"/junos/interfaces/interface/macsec/scsa-tx/state/counters/sc-auth-only-packets",
		},
		"/openconfig/macsec/interfaces/interface/scsa-tx/scsa-tx/state/counters/sc-encrypted": {
			"/junos/interfaces/interface/macsec/scsa-tx/state/counters/sc-encrypted-packets",
		},
		"/openconfig/macsec/interfaces/interface/scsa-tx/scsa-tx/state/counters/sa-auth-only": {
			"/junos/interfaces/interface/macsec/scsa-tx/state/counters/sa-auth-only-packets",
		},
		"/openconfig/macsec/interfaces/interface/scsa-tx/scsa-tx/state/counters/sa-encrypted": {
			"/junos/interfaces/interface/macsec/scsa-tx/state/counters/sa-encrypted-packets",
		},
		// Secure Channel RX statistics
		"/openconfig/macsec/interfaces/interface/scsa-rx/scsa-rx/state/counters/sc-invalid": {
			"/junos/interfaces/interface/macsec/scsa-rx/state/counters/sc-invalid-packets",
		},
		"/openconfig/macsec/interfaces/interface/scsa-rx/scsa-rx/state/counters/sc-valid": {
			"/junos/interfaces/interface/macsec/scsa-rx/state/counters/sc-valid-packets",
		},
		"/openconfig/macsec/interfaces/interface/scsa-rx/scsa-rx/state/counters/sa-invalid": {
			"/junos/interfaces/interface/macsec/scsa-rx/state/counters/sa-invalid-packets",
		},
		"/openconfig/macsec/interfaces/interface/scsa-rx/scsa-rx/state/counters/sa-valid": {
			"/junos/interfaces/interface/macsec/scsa-rx/state/counters/sa-valid-packets",
		},
	}

	// Update path patterns for Juniper native counters
	updatePathPatterns = []*gnmipb.Path{
		// Interface level counters
		{
			Origin: "junos",
			Elem: []*gnmipb.PathElem{
				{Name: "interfaces"}, {Name: "interface"},
				{Name: "*"}, // interface-name
				{Name: "macsec"}, {Name: "state"}, {Name: "counters"},
				{Name: "tx-untagged-packets"},
			},
		},
		{
			Origin: "junos",
			Elem: []*gnmipb.PathElem{
				{Name: "interfaces"}, {Name: "interface"},
				{Name: "*"}, // interface-name
				{Name: "macsec"}, {Name: "state"}, {Name: "counters"},
				{Name: "rx-untagged-packets"},
			},
		},
		{
			Origin: "junos",
			Elem: []*gnmipb.PathElem{
				{Name: "interfaces"}, {Name: "interface"},
				{Name: "*"}, // interface-name
				{Name: "macsec"}, {Name: "state"}, {Name: "counters"},
				{Name: "rx-badtag-packets"},
			},
		},
		{
			Origin: "junos",
			Elem: []*gnmipb.PathElem{
				{Name: "interfaces"}, {Name: "interface"},
				{Name: "*"}, // interface-name
				{Name: "macsec"}, {Name: "state"}, {Name: "counters"},
				{Name: "rx-unknownsci-packets"},
			},
		},
		{
			Origin: "junos",
			Elem: []*gnmipb.PathElem{
				{Name: "interfaces"}, {Name: "interface"},
				{Name: "*"}, // interface-name
				{Name: "macsec"}, {Name: "state"}, {Name: "counters"},
				{Name: "rx-nosci-packets"},
			},
		},
		// TX SC/SA counters
		{
			Origin: "junos",
			Elem: []*gnmipb.PathElem{
				{Name: "interfaces"}, {Name: "interface"},
				{Name: "*"}, // interface-name
				{Name: "macsec"}, {Name: "scsa-tx"}, {Name: "state"}, {Name: "counters"},
				{Name: "*"}, // sci-tx
				{Name: "sc-auth-only-packets"},
			},
		},
		{
			Origin: "junos",
			Elem: []*gnmipb.PathElem{
				{Name: "interfaces"}, {Name: "interface"},
				{Name: "*"}, // interface-name
				{Name: "macsec"}, {Name: "scsa-tx"}, {Name: "state"}, {Name: "counters"},
				{Name: "*"}, // sci-tx
				{Name: "sc-encrypted-packets"},
			},
		},
		{
			Origin: "junos",
			Elem: []*gnmipb.PathElem{
				{Name: "interfaces"}, {Name: "interface"},
				{Name: "*"}, // interface-name
				{Name: "macsec"}, {Name: "scsa-tx"}, {Name: "state"}, {Name: "counters"},
				{Name: "*"}, // sci-tx
				{Name: "sa-auth-only-packets"},
			},
		},
		{
			Origin: "junos",
			Elem: []*gnmipb.PathElem{
				{Name: "interfaces"}, {Name: "interface"},
				{Name: "*"}, // interface-name
				{Name: "macsec"}, {Name: "scsa-tx"}, {Name: "state"}, {Name: "counters"},
				{Name: "*"}, // sci-tx
				{Name: "sa-encrypted-packets"},
			},
		},
		// RX SC/SA counters
		{
			Origin: "junos",
			Elem: []*gnmipb.PathElem{
				{Name: "interfaces"}, {Name: "interface"},
				{Name: "*"}, // interface-name
				{Name: "macsec"}, {Name: "scsa-rx"}, {Name: "state"}, {Name: "counters"},
				{Name: "*"}, // sci-rx
				{Name: "sc-invalid-packets"},
			},
		},
		{
			Origin: "junos",
			Elem: []*gnmipb.PathElem{
				{Name: "interfaces"}, {Name: "interface"},
				{Name: "*"}, // interface-name
				{Name: "macsec"}, {Name: "scsa-rx"}, {Name: "state"}, {Name: "counters"},
				{Name: "*"}, // sci-rx
				{Name: "sc-valid-packets"},
			},
		},
		{
			Origin: "junos",
			Elem: []*gnmipb.PathElem{
				{Name: "interfaces"}, {Name: "interface"},
				{Name: "*"}, // interface-name
				{Name: "macsec"}, {Name: "scsa-rx"}, {Name: "state"}, {Name: "counters"},
				{Name: "*"}, // sci-rx
				{Name: "sa-invalid-packets"},
			},
		},
		{
			Origin: "junos",
			Elem: []*gnmipb.PathElem{
				{Name: "interfaces"}, {Name: "interface"},
				{Name: "*"}, // interface-name
				{Name: "macsec"}, {Name: "scsa-rx"}, {Name: "state"}, {Name: "counters"},
				{Name: "*"}, // sci-rx
				{Name: "sa-valid-packets"},
			},
		},
	}

	// Mapping of Juniper leaf names to OpenConfig leaf names
	vendorToOCLeaf = map[string]string{
		"tx-untagged-packets":   "tx-untagged-pkts",
		"rx-untagged-packets":   "rx-untagged-pkts",
		"rx-badtag-packets":     "rx-badtag-pkts",
		"rx-unknownsci-packets": "rx-unknownsci-pkts",
		"rx-nosci-packets":      "rx-nosci-pkts",
		"sc-auth-only-packets":  "sc-auth-only",
		"sc-encrypted-packets":  "sc-encrypted",
		"sa-auth-only-packets":  "sa-auth-only",
		"sa-encrypted-packets":  "sa-encrypted",
		"sc-invalid-packets":    "sc-invalid",
		"sc-valid-packets":      "sc-valid",
		"sa-invalid-packets":    "sa-invalid",
		"sa-valid-packets":      "sa-valid",
	}
)

// New creates a functional translator for Juniper MACsec counters.
func New() *translator.FunctionalTranslator {
	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               ftconsts.JuniperMacsecCountersTranslator,
			Translate:        translate,
			OutputToInputMap: ftutilities.MustStringMapPaths(translateMap),
			Metadata: []*translator.FTMetadata{
				{
					Vendor: ftconsts.VendorJuniper,
					SoftwareVersionRange: &translator.SWRange{
						InclusiveMin: "21.4",
						ExclusiveMax: "24.0",
					},
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("Failed to create Juniper MACsec counters functional translator: %v", err)
	}
	return ft
}

// outgoingVal converts Juniper native counter value to OpenConfig typed value.
func outgoingVal(fullPath *gnmipb.Path, incomingVal *gnmipb.TypedValue) (*gnmipb.TypedValue, error) {
	jsonVal := incomingVal.GetJsonVal()
	if jsonVal == nil {
		return incomingVal, nil
	}
	var jsonValMap map[string]any
	if err := json.Unmarshal(jsonVal, &jsonValMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON value for path %v: %v. JSON: %s. Skipping this update", fullPath, err, string(jsonVal))
	}

	val, ok := jsonValMap["value"]
	if !ok {
		return nil, fmt.Errorf("value not found in JSON for path: %v", fullPath)
	}
	v, ok := val.(float64)
	if !ok {
		return nil, fmt.Errorf("value has unexpected JSON type %v, skipping update for path: %v", val, fullPath)
	}
	return &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: uint64(v)}}, nil
}

// extractInterfaceAndSCI extracts interface name and SCI (if applicable) from the path.
func extractInterfaceAndSCI(fullPath *gnmipb.Path) (interfaceName, sci string, isScsa bool, err error) {
	elems := fullPath.GetElem()
	if len(elems) < 3 {
		return "", "", false, fmt.Errorf("path %v has fewer than 3 elements", fullPath)
	}

	// Find interface element - it's after "interfaces/interface"
	var intfIdx int
	for i, elem := range elems {
		if elem.GetName() == "interface" && i > 0 && elems[i-1].GetName() == "interfaces" {
			intfIdx = i + 1
			break
		}
	}

	if intfIdx >= len(elems) {
		return "", "", false, fmt.Errorf("could not find interface element in path %v", fullPath)
	}

	interfaceName = elems[intfIdx].GetName()

	// Check if this is an SCSA path (contains scsa-tx or scsa-rx)
	for _, elem := range elems {
		if elem.GetName() == "scsa-tx" || elem.GetName() == "scsa-rx" {
			isScsa = true
			break
		}
	}

	// If SCSA path, extract SCI - it's the first wildcard match after scsa-tx/scsa-rx
	if isScsa {
		for i, elem := range elems {
			if (elem.GetName() == "scsa-tx" || elem.GetName() == "scsa-rx") && i+2 < len(elems) {
				// Skip "state" and get the next element which should be a counter container
				if elems[i+1].GetName() == "state" && elems[i+2].GetName() == "counters" && i+3 < len(elems) {
					sci = elems[i+3].GetName()
					break
				}
			}
		}
	}

	return interfaceName, sci, isScsa, nil
}

// updateHandler processes updates and transforms them to OpenConfig format.
func updateHandler(n *gnmipb.Notification) ([]*gnmipb.Update, error) {
	prefix := n.GetPrefix()
	var updates []*gnmipb.Update

	for _, update := range n.GetUpdate() {
		fullPath := ftutilities.Join(prefix, update.GetPath())

		for _, pattern := range updatePathPatterns {
			matched := ftutilities.MatchPath(fullPath, pattern)
			if !matched {
				continue
			}

			vendorLeaf := fullPath.GetElem()[len(fullPath.GetElem())-1].GetName()
			ocLeaf, found := vendorToOCLeaf[vendorLeaf]
			if !found {
				return nil, fmt.Errorf("vendor leaf '%s' not found in mapping for path %v", vendorLeaf, fullPath)
			}

			interfaceName, sci, isScsa, err := extractInterfaceAndSCI(fullPath)
			if err != nil {
				return nil, fmt.Errorf("failed to extract interface and SCI from path %v: %v", fullPath, err)
			}

			incomingVal := update.GetVal()
			outVal, err := outgoingVal(fullPath, incomingVal)
			if err != nil {
				return nil, fmt.Errorf("failed to get outgoing value for path %v: %v", fullPath, err)
			}

			var outgoingUpdate *gnmipb.Update
			if isScsa {
				// Determine if TX or RX from the path
				var scaType string
				for _, elem := range fullPath.GetElem() {
					if elem.GetName() == "scsa-tx" {
						scaType = "scsa-tx"
						break
					} else if elem.GetName() == "scsa-rx" {
						scaType = "scsa-rx"
						break
					}
				}

				outgoingUpdate = &gnmipb.Update{
					Path: returnScSAPath(interfaceName, scaType, sci, ocLeaf),
					Val:  outVal,
				}
			} else {
				outgoingUpdate = &gnmipb.Update{
					Path: returnInterfacePath(interfaceName, ocLeaf),
					Val:  outVal,
				}
			}

			updates = append(updates, outgoingUpdate)
		}
	}
	return updates, nil
}

// returnInterfacePath returns the OpenConfig path for interface level counters.
func returnInterfacePath(interfaceName, leaf string) *gnmipb.Path {
	return &gnmipb.Path{
		Elem: []*gnmipb.PathElem{
			{Name: "macsec"},
			{Name: "interfaces"},
			{
				Name: "interface",
				Key: map[string]string{
					"name": interfaceName,
				},
			},
			{Name: "state"},
			{Name: "counters"},
			{Name: leaf},
		},
	}
}

// returnScSAPath returns the OpenConfig path for SC/SA level counters.
func returnScSAPath(interfaceName, scaType, sci, leaf string) *gnmipb.Path {
	return &gnmipb.Path{
		Elem: []*gnmipb.PathElem{
			{Name: "macsec"},
			{Name: "interfaces"},
			{
				Name: "interface",
				Key: map[string]string{
					"name": interfaceName,
				},
			},
			{Name: scaType},
			{
				Name: scaType,
				Key: map[string]string{
					scaType[5:]: sci, // "sci-tx" or "sci-rx"
				},
			},
			{Name: "state"},
			{Name: "counters"},
			{Name: leaf},
		},
	}
}

func translate(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
	notification := sr.GetUpdate()

	updates, err := updateHandler(notification)
	if err != nil {
		return nil, err
	}

	// Return early if there are no updates.
	if len(updates) == 0 {
		return nil, nil
	}

	outgoingNotification := &gnmipb.Notification{
		Timestamp: notification.GetTimestamp(),
		Prefix: &gnmipb.Path{
			Origin: "openconfig",
			Target: notification.GetPrefix().GetTarget(),
		},
		Update: updates,
	}
	return &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: outgoingNotification,
		},
	}, nil
}
