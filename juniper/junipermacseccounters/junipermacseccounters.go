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

// Package junipermacseccounters translates MACsec interface counters from Juniper native to OpenConfig.
package junipermacseccounters

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
	// Juniper native paths to OpenConfig MACsec interface counters translation map.
	// Note: The origin (junos) is required by the framework but must be stripped
	// before subscribing on the device.
	translateMap = map[string][]string{
		"/openconfig/macsec/interfaces/interface/state/counters/rx-badicv-pkts": {
			"/junos/macsec/interfaces/interface/mka/state/counters/jnx-integrity-check-value-mismatch",
		},
		"/openconfig/macsec/interfaces/interface/state/counters/tx-pkts-ctrl": {
			"/junos/macsec/interfaces/interface/mka/state/counters/out-mkpdu",
		},
		"/openconfig/macsec/interfaces/interface/state/counters/rx-pkts-ctrl": {
			"/junos/macsec/interfaces/interface/mka/state/counters/in-mkpdu",
		},
		"/openconfig/macsec/interfaces/interface/state/counters/rx-unrecognized-ckn": {
			"/junos/macsec/interfaces/interface/mka/state/counters/jnx-cak-error",
		},
		"/openconfig/macsec/interfaces/interface/state/counters/tx-pkts-err-in": {
			"/junos/macsec/interfaces/interface/mka/state/counters/out-mkpdu",
		},
		"/openconfig/macsec/interfaces/interface/state/counters/tx-pkts-dropped": {
			"/junos/macsec/interfaces/interface/mka/state/counters/out-mkpdu",
		},
		"/openconfig/macsec/interfaces/interface/state/counters/rx-pkts-dropped": {
			"/junos/macsec/interfaces/interface/mka/state/counters/in-mkpdu",
		},
	}

	// Update path patterns for Juniper native counters
	updatePathPatterns = []*gnmipb.Path{
		{
			Origin: "junos",
			Elem: []*gnmipb.PathElem{
				{Name: "macsec"}, {Name: "interfaces"},
				{Name: "interface", Key: map[string]string{"name": "*"}},
				{Name: "mka"}, {Name: "state"}, {Name: "counters"},
				{Name: "jnx-integrity-check-value-mismatch"},
			},
		},
		{
			Origin: "junos",
			Elem: []*gnmipb.PathElem{
				{Name: "macsec"}, {Name: "interfaces"},
				{Name: "interface", Key: map[string]string{"name": "*"}},
				{Name: "mka"}, {Name: "state"}, {Name: "counters"},
				{Name: "out-mkpdu"},
			},
		},
		{
			Origin: "junos",
			Elem: []*gnmipb.PathElem{
				{Name: "macsec"}, {Name: "interfaces"},
				{Name: "interface", Key: map[string]string{"name": "*"}},
				{Name: "mka"}, {Name: "state"}, {Name: "counters"},
				{Name: "in-mkpdu"},
			},
		},
		{
			Origin: "junos",
			Elem: []*gnmipb.PathElem{
				{Name: "macsec"}, {Name: "interfaces"},
				{Name: "interface", Key: map[string]string{"name": "*"}},
				{Name: "mka"}, {Name: "state"}, {Name: "counters"},
				{Name: "jnx-cak-error"},
			},
		},
	}

	// Mapping of Juniper leaf names to OpenConfig leaf names (one native leaf may map to multiple OC leaves)
	vendorToOCLeaves = map[string][]string{
		"jnx-integrity-check-value-mismatch": {"rx-badicv-pkts"},
		"out-mkpdu":                          {"tx-pkts-ctrl", "tx-pkts-err-in", "tx-pkts-dropped"},
		"in-mkpdu":                           {"rx-pkts-ctrl", "rx-pkts-dropped"},
		"jnx-cak-error":                      {"rx-unrecognized-ckn"},
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
			ocLeaves, found := vendorToOCLeaves[vendorLeaf]
			if !found {
				return nil, fmt.Errorf("vendor leaf %q not found in mapping for path %v", vendorLeaf, fullPath)
			}

			// Extract interface name from the "interface" element key
			var intfID string
			for _, elem := range fullPath.GetElem() {
				if elem.GetName() == "interface" {
					intfID = elem.GetKey()["name"]
					break
				}
			}
			if intfID == "" {
				return nil, fmt.Errorf("interface name key not found in path %v", fullPath)
			}

			incomingVal := update.GetVal()
			outVal, err := outgoingVal(fullPath, incomingVal)
			if err != nil {
				return nil, fmt.Errorf("failed to get outgoing value for path %v: %v", fullPath, err)
			}

			for _, ocLeaf := range ocLeaves {
				outgoingUpdate := &gnmipb.Update{
					Path: returnPath(intfID, ocLeaf),
					Val:  outVal,
				}
				updates = append(updates, outgoingUpdate)
			}
			break
		}
	}
	return updates, nil
}

// returnPath returns a gNMI path for the update.
func returnPath(interfaceName, leaf string) *gnmipb.Path {
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
