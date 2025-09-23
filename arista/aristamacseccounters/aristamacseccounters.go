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

// Package aristamacseccounters translates the interface MACSec counters from native to openconfig.
package aristamacseccounters

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
	translateMap = map[string][]string{
		// This leaf will be in SAMPLED mode because of the `counters` keyword - updates every 30 seconds.
		"/openconfig/macsec/interfaces/interface/state/counters/rx-badicv-pkts": {
			"/eos_native/Smash/macsec/counters/msgCounter",
		},
		"/openconfig/macsec/interfaces/interface/state/counters/rx-unrecognized-ckn": {
			"/eos_native/Smash/macsec/counters/msgCounter",
		},
		"/openconfig/macsec/interfaces/interface/state/counters/tx-pkts-err-in": {
			"/eos_native/Smash/hardware/counter/macsec",
		},
		"/openconfig/macsec/interfaces/interface/state/counters/tx-pkts-ctrl": {
			"/eos_native/Smash/hardware/counter/macsec",
		},
		"/openconfig/macsec/interfaces/interface/state/counters/rx-pkts-ctrl": {
			"/eos_native/Smash/hardware/counter/macsec",
		},
		"/openconfig/macsec/interfaces/interface/state/counters/tx-pkts-dropped": {
			"/eos_native/Smash/counters/ethIntf/SandCounters/current/counter",
		},
		"/openconfig/macsec/interfaces/interface/state/counters/rx-pkts-dropped": {
			"/eos_native/Smash/counters/ethIntf/SandCounters/current/counter",
		},
	}
	updatePathPatterns = []*gnmipb.Path{
		{
			Origin: "eos_native",
			Elem: []*gnmipb.PathElem{
				{Name: "Smash"}, {Name: "macsec"}, {Name: "counters"}, {Name: "msgCounter"},
				{Name: "*"}, // interface-id
				{Name: "rxMsgCounter"},
				{Name: "icvValidationErr"},
			},
		},
		{
			Origin: "eos_native",
			Elem: []*gnmipb.PathElem{
				{Name: "Smash"}, {Name: "macsec"}, {Name: "counters"}, {Name: "msgCounter"},
				{Name: "*"}, // interface-id
				{Name: "rxMsgCounter"},
				{Name: "unrecognizedCkn"},
			},
		},
		{
			Origin: "eos_native",
			Elem: []*gnmipb.PathElem{
				{Name: "Smash"}, {Name: "hardware"}, {Name: "counter"}, {Name: "macsec"},
				{Name: "*"}, // slice-id
				{Name: "current"},
				{Name: "*"}, // interface-id
				{Name: "macSecCounters"},
				{Name: "outPktsTooLong"},
			},
		},
		{
			Origin: "eos_native",
			Elem: []*gnmipb.PathElem{
				{Name: "Smash"}, {Name: "hardware"}, {Name: "counter"}, {Name: "macsec"},
				{Name: "*"}, // slice-id
				{Name: "current"},
				{Name: "*"}, // interface-id
				{Name: "macSecCounters"},
				{Name: "outPktCtrl"},
			},
		},
		{
			Origin: "eos_native",
			Elem: []*gnmipb.PathElem{
				{Name: "Smash"}, {Name: "hardware"}, {Name: "counter"}, {Name: "macsec"},
				{Name: "*"}, // slice-id
				{Name: "current"},
				{Name: "*"}, // interface-id
				{Name: "macSecCounters"},
				{Name: "inPktsCtrl"},
			},
		},
		{
			Origin: "eos_native",
			Elem: []*gnmipb.PathElem{
				{Name: "Smash"}, {Name: "counters"}, {Name: "ethIntf"}, {Name: "SandCounters"}, {Name: "current"}, {Name: "counter"},
				{Name: "*"}, // interface-id
				{Name: "statistics"},
				{Name: "outDiscards"},
			},
		},
		{
			Origin: "eos_native",
			Elem: []*gnmipb.PathElem{
				{Name: "Smash"}, {Name: "counters"}, {Name: "ethIntf"}, {Name: "SandCounters"}, {Name: "current"}, {Name: "counter"},
				{Name: "*"}, // interface-id
				{Name: "statistics"},
				{Name: "inDiscards"},
			},
		},
	}

	vendorToOCLeaf = map[string]string{
		"icvValidationErr": "rx-badicv-pkts",
		"unrecognizedCkn":  "rx-unrecognized-ckn",
		"outPktsTooLong":   "tx-pkts-err-in",
		"outPktCtrl":       "tx-pkts-ctrl",
		"inPktsCtrl":       "rx-pkts-ctrl",
		"outDiscards":      "tx-pkts-dropped",
		"inDiscards":       "rx-pkts-dropped",
	}
)

// New creates a functional translator.
func New() *translator.FunctionalTranslator {
	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               ftconsts.AristaMacsecCountersTranslator,
			Translate:        translate,
			OutputToInputMap: ftutilities.MustStringMapPaths(translateMap),
			Metadata: []*translator.FTMetadata{
				{
					Vendor: ftconsts.VendorArista,
					SoftwareVersionRange: &translator.SWRange{
						InclusiveMin: "4.33.0F",
						ExclusiveMax: "4.35",
					},
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("Failed to create Arista MACSec counters functional translator: %v", err)
	}
	return ft
}

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
		return nil, fmt.Errorf("value has unexpected Json %v, Skipping update for this path: %v", val, fullPath)
	}
	return &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: uint64(v)}}, nil
}

// updateHandler returns the deletes that should be sent to the target.
func updateHandler(n *gnmipb.Notification) ([]*gnmipb.Update, error) {
	prefix := n.GetPrefix()
	var updates []*gnmipb.Update

	for _, update := range n.GetUpdate() {
		fullPath := ftutilities.Join(prefix, update.GetPath())

		for _, path := range updatePathPatterns {
			matched := ftutilities.MatchPath(fullPath, path)
			if !matched {
				continue
			}
			// Get the name of last element in the path.
			vendorLeaf := fullPath.GetElem()[len(fullPath.GetElem())-1].GetName()
			// Lookup the mapping of vendor leaf to oc leaf.
			ocLeaf, found := vendorToOCLeaf[vendorLeaf]
			if !found {
				return nil, fmt.Errorf("interesting update but skipping due to missing OC leaf mapping: path %v", fullPath)
			}

			// The interface ID is always the 3rd last element in the path.
			lastElemIndex := len(fullPath.GetElem()) - 1
			intfID := fullPath.GetElem()[lastElemIndex-2].GetName()
			// One of the vendor paths `outDiscards` is sending deprecated Json value
			// instead of gnmi TypeValue.
			incomingVal := update.GetVal()
			outVal, err := outgoingVal(fullPath, incomingVal)
			if err != nil {
				return nil, fmt.Errorf("failed to get outgoing value for path %v: %v", fullPath, err)
			}
			outgoingUpdate := &gnmipb.Update{
				Path: returnPath(intfID, ocLeaf),
				Val:  outVal,
			}
			updates = append(updates, outgoingUpdate)
		}
	}
	return updates, nil
}

// returnPath returns a gNMI path for the update.
// Does not set the origin or the target.
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

	// Explicit delete handler is not required for this translator.
	// This leaf will be in SAMPLED mode so explicit deletes are not supported and will instead rely
	// on the implicit delete handler.

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
