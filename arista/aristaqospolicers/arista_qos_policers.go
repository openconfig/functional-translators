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

// Package aristaqospolicers implements the functional translator for Arista QoS interface policer drops.
package aristaqospolicers

import (
	"strings"

	log "github.com/golang/glog"
	"github.com/openconfig/functional-translators/ftutilities"
	"github.com/openconfig/functional-translators/translator"
	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

const (
	violatingOctets = "violating-octets"
	violatingPkts   = "violating-pkts"
	exceedingOctets = "exceeding-octets"
	exceedingPkts   = "exceeding-pkts"
)

var (
	translateMap = map[string][]string{
		"/openconfig/qos/interfaces/interface/input/scheduler-policy/schedulers/scheduler/state/violating-octets": {
			"/eos_native/Sysdb/qos/status/policingStatus/intfPolicer",
		},
		"/openconfig/qos/interfaces/interface/input/scheduler-policy/schedulers/scheduler/state/violating-pkts": {
			"/eos_native/Sysdb/qos/status/policingStatus/intfPolicer",
		},
		"/openconfig/qos/interfaces/interface/input/scheduler-policy/schedulers/scheduler/state/exceeding-octets": {
			"/eos_native/Sysdb/qos/status/policingStatus/intfPolicer",
		},
		"/openconfig/qos/interfaces/interface/input/scheduler-policy/schedulers/scheduler/state/exceeding-pkts": {
			"/eos_native/Sysdb/qos/status/policingStatus/intfPolicer",
		},
		"/openconfig/qos/interfaces/interface/input/scheduler-policy/schedulers/scheduler/state/conforming-octets": {
			"/openconfig/qos/interfaces/interface/input/scheduler-policy/schedulers/scheduler/state/conforming-octets",
		},
		"/openconfig/qos/interfaces/interface/input/scheduler-policy/schedulers/scheduler/state/conforming-pkts": {
			"/openconfig/qos/interfaces/interface/input/scheduler-policy/schedulers/scheduler/state/conforming-pkts",
		},
		"/openconfig/qos/interfaces/interface/input/scheduler-policy/schedulers/scheduler/state/sequence": {
			"/openconfig/qos/interfaces/interface/input/scheduler-policy/schedulers/scheduler/state/sequence",
		},
	}
)

// New returns a new FunctionalTranslator for Arista QoS policers.
func New() *translator.FunctionalTranslator {
	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               "arista-qos-policers-ft",
			Translate:        translate,
			OutputToInputMap: ftutilities.MustStringMapPaths(translateMap),
			Metadata: []*translator.FTMetadata{
				{
					Vendor: "Arista",
					SoftwareVersionRange: &translator.SWRange{
						InclusiveMin: "4.33.0F",
						ExclusiveMax: "4.37",
					},
				},
			},
		})
	if err != nil {
		log.Fatalf("failed to create arista qos policers functional translator: %v", err)
	}
	return ft
}

func translate(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
	notification := sr.GetUpdate()
	if notification == nil {
		return nil, nil
	}

	prefix := notification.GetPrefix()
	target := notification.GetPrefix().GetTarget()
	if target == "" {
		return nil, nil
	}

	if prefix.GetOrigin() == "openconfig" {
		updatedIntfs := make(map[string]bool)

		// Find which interfaces are receiving updates in this notification
		for _, u := range notification.GetUpdate() {
			fullPath := ftutilities.Join(prefix, u.GetPath())
			if len(fullPath.GetElem()) > 2 && fullPath.GetElem()[2].GetName() == "interface" {
				intfID := fullPath.GetElem()[2].GetKey()["interface-id"]
				if intfID != "" {
					updatedIntfs[intfID] = true
				}
			}
		}

		var validDeletes []*gnmipb.Path
		// Filter the deletes
		for _, d := range notification.GetDelete() {
			fullPath := ftutilities.Join(prefix, d)
			elems := fullPath.GetElem()
			// Is this Arista's destructive "Replace" delete?
			// b/508518501#comment9, point 2.
			// (It specifically targets the "state" container)
			if len(elems) > 0 && elems[len(elems)-1].GetName() == "state" {
				if len(elems) > 2 && elems[2].GetName() == "interface" && elems[len(elems)-2].GetName() == "scheduler" {
					intfID := elems[2].GetKey()["interface-id"]
					// If the interface is being updated in the same message, strip the state delete!
					if updatedIntfs[intfID] {
						continue
					}
				}
			}
			// Otherwise, it's a legitimate structural delete (e.g., scheduler-policy). Keep it.
			validDeletes = append(validDeletes, d)
		}

		// Drop notifications that are completely empty after filtering
		if len(notification.GetUpdate()) == 0 && len(validDeletes) == 0 {
			return nil, nil
		}

		// Forward the clean payload
		return &gnmipb.SubscribeResponse{
			Response: &gnmipb.SubscribeResponse_Update{
				Update: &gnmipb.Notification{
					Timestamp: notification.GetTimestamp(),
					Prefix:    prefix,
					Update:    notification.GetUpdate(),
					Delete:    validDeletes,
				},
			},
		}, nil
	}

	var outgoingUpdates []*gnmipb.Update
	var outgoingDeletes []*gnmipb.Path
	targetInfo := ftutilities.PolicerAggMap.CreateOrUpdateTargetPolicerInfo(target)
	updatedAPs := make(map[string]bool)
	deletedAPs := make(map[string]bool)

	// Handle gNMI Updates (Incremental count updates)
	for _, update := range notification.GetUpdate() {
		fullPath := ftutilities.Join(prefix, update.GetPath())
		physIntf, attachmentPoint, ok := parseIntfPairKey(fullPath)
		if !ok {
			continue
		}

		apInfo := targetInfo.CreateOrRetrieveAttachmentPoint(attachmentPoint)
		elems := update.GetPath().GetElem()
		if len(elems) == 0 {
			continue
		}
		lastElem := elems[len(elems)-1].GetName()

		// Send counter update to the cache manager
		if apInfo.SetMemberCounter(physIntf, lastElem, update.GetVal().GetUintVal()) {
			updatedAPs[attachmentPoint] = true
		}
	}

	// Handle gNMI Deletes (Physical member goes down or is removed from bundle)
	for _, path := range notification.GetDelete() {
		fullPath := ftutilities.Join(prefix, path)
		physIntf, attachmentPoint, ok := parseIntfPairKey(fullPath)
		if !ok {
			continue
		}

		remainingMembers, ok := targetInfo.RemoveMemberAndCleanup(attachmentPoint, physIntf)
		if !ok {
			continue
		}
		if remainingMembers == 0 {
			deletedAPs[attachmentPoint] = true
		} else {
			updatedAPs[attachmentPoint] = true
		}
	}

	for ap := range updatedAPs {
		if info, ok := targetInfo.AttachmentPointInfo(ap); ok {
			byteDrop, pktDrop, yellowByte, yellowPkt := info.AggregateCounters()

			outgoingUpdates = append(outgoingUpdates, buildOCUpdate(ap, violatingOctets, byteDrop))
			outgoingUpdates = append(outgoingUpdates, buildOCUpdate(ap, violatingPkts, pktDrop))
			outgoingUpdates = append(outgoingUpdates, buildOCUpdate(ap, exceedingOctets, yellowByte))
			outgoingUpdates = append(outgoingUpdates, buildOCUpdate(ap, exceedingPkts, yellowPkt))
		}
	}

	for ap := range deletedAPs {
		outgoingDeletes = append(outgoingDeletes, buildRelativeOCPath(ap, violatingOctets))
		outgoingDeletes = append(outgoingDeletes, buildRelativeOCPath(ap, violatingPkts))
		outgoingDeletes = append(outgoingDeletes, buildRelativeOCPath(ap, exceedingOctets))
		outgoingDeletes = append(outgoingDeletes, buildRelativeOCPath(ap, exceedingPkts))
	}

	if len(outgoingUpdates) == 0 && len(outgoingDeletes) == 0 {
		return nil, nil
	}

	return &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: notification.GetTimestamp(),
				Prefix: &gnmipb.Path{
					Origin: "openconfig",
					Target: target,
					Elem: []*gnmipb.PathElem{
						{Name: "qos"},
						{Name: "interfaces"},
					},
				},
				Update: outgoingUpdates,
				Delete: outgoingDeletes,
			},
		},
	}, nil
}

// extractIntfPairKey finds the intfPolicer element and returns the subsequent key.
func extractIntfPairKey(path *gnmipb.Path) string {
	elems := path.GetElem()
	for i, elem := range elems {
		if elem.GetName() == "intfPolicer" && i+1 < len(elems) {
			return elems[i+1].GetName()
		}
	}
	return ""
}

// buildRelativeOCPath creates the path to the specific QoS leaf, relative to /qos/interfaces.
func buildRelativeOCPath(attachmentPoint, leafName string) *gnmipb.Path {
	return &gnmipb.Path{
		Elem: []*gnmipb.PathElem{
			{
				Name: "interface",
				Key:  map[string]string{"interface-id": attachmentPoint},
			},
			{Name: "input"},
			{Name: "scheduler-policy"},
			{Name: "schedulers"},
			{
				Name: "scheduler",
				Key:  map[string]string{"sequence": "0"},
			},
			{Name: "state"},
			{Name: leafName},
		},
	}
}

// buildOCUpdate creates the fully formed OpenConfig gNMI Update.
func buildOCUpdate(attachmentPoint, leafName string, val uint64) *gnmipb.Update {
	return &gnmipb.Update{
		Path: buildRelativeOCPath(attachmentPoint, leafName),
		Val:  &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: val}},
	}
}

// parseIntfPairKey extracts and parses the intfPolicer key.
// The key is formatted as <IngressPhysicalPort>_<PolicyAttachmentPort>.
// Returns physIntf, attachmentPoint, and a boolean indicating if parsing was successful.
func parseIntfPairKey(path *gnmipb.Path) (string, string, bool) {
	intfPairKey := extractIntfPairKey(path)
	if intfPairKey == "" {
		return "", "", false
	}

	parts := strings.SplitN(intfPairKey, "_", 2)
	if len(parts) != 2 || parts[1] == "" {
		// Per requirements: only logical interfaces/subints have policers.
		// Skip if it doesn't have an attachment point after the underscore.
		return "", "", false
	}
	return parts[0], parts[1], true
}
