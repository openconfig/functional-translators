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

// counterType is an enum for different policer counter names.
type counterType int

const (
	byteDrop counterType = iota
	pktDrop
	yellowByteIn
	yellowPktIn
)

// counterMap maps string counter names to their counterType enum.
var counterMap = map[string]counterType{
	"byteDropCount":     byteDrop,
	"pktDropCount":      pktDrop,
	"yellowByteInCount": yellowByteIn,
	"yellowPktInCount":  yellowPktIn,
}

// memberCounters holds the 4 QoS counter metrics for a single physical member interface.
type memberCounters struct {
	byteDropCount     uint64
	pktDropCount      uint64
	yellowByteInCount uint64
	yellowPktInCount  uint64
}

// attachmentPointInfo holds drop counts for all physical members of a specific logical interface.
type attachmentPointInfo struct {
	name    string
	members map[string]*memberCounters // map[physicalInterfaceName]*memberCounters
}

func (a *attachmentPointInfo) setMemberCounter(member, counterName string, count uint64) bool {
	if a.members == nil {
		a.members = make(map[string]*memberCounters)
	}
	mCounters, ok := a.members[member]
	if !ok {
		mCounters = &memberCounters{}
		a.members[member] = mCounters
	}
	cType, ok := counterMap[counterName]
	if !ok {
		return false
	}

	switch cType {
	case byteDrop:
		mCounters.byteDropCount = count
	case pktDrop:
		mCounters.pktDropCount = count
	case yellowByteIn:
		mCounters.yellowByteInCount = count
	case yellowPktIn:
		mCounters.yellowPktInCount = count
	}
	return true
}

func (a *attachmentPointInfo) aggregateCounters() (byteDrop, pktDrop, yellowByte, yellowPkt uint64) {
	for _, m := range a.members {
		byteDrop += m.byteDropCount
		pktDrop += m.pktDropCount
		yellowByte += m.yellowByteInCount
		yellowPkt += m.yellowPktInCount
	}
	return
}

// policerState holds policer drop information for a specific target device.
type policerState struct {
	attachmentPoints map[string]*attachmentPointInfo // map[attachmentPointName]*attachmentPointInfo
}

func newPolicerState() *policerState {
	return &policerState{
		attachmentPoints: make(map[string]*attachmentPointInfo),
	}
}

func (p *policerState) createOrRetrieveAttachmentPoint(apName string) *attachmentPointInfo {
	if p.attachmentPoints == nil {
		p.attachmentPoints = make(map[string]*attachmentPointInfo)
	}
	if _, ok := p.attachmentPoints[apName]; !ok {
		p.attachmentPoints[apName] = &attachmentPointInfo{
			name:    apName,
			members: make(map[string]*memberCounters),
		}
	}
	return p.attachmentPoints[apName]
}

func (p *policerState) attachmentPointInfo(apName string) (*attachmentPointInfo, bool) {
	info, ok := p.attachmentPoints[apName]
	return info, ok
}

func (p *policerState) removeMemberAndCleanup(apName, member string) (int, bool) {
	apInfo, ok := p.attachmentPoints[apName]
	if !ok {
		return 0, false
	}
	if apInfo.members != nil {
		delete(apInfo.members, member)
	}
	remaining := len(apInfo.members)
	if remaining == 0 {
		delete(p.attachmentPoints, apName)
	}
	return remaining, true
}

type impl struct {
	state *policerState
}

// New returns a new FunctionalTranslator for Arista QoS policers.
func New() *translator.FunctionalTranslator {
	i := &impl{
		state: newPolicerState(),
	}
	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               "arista-qos-policers-ft",
			Translate:        i.translate,
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

func (i *impl) translate(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
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
	updatedAPs := make(map[string]bool)
	deletedAPs := make(map[string]bool)

	// Handle gNMI Updates (Incremental count updates)
	for _, update := range notification.GetUpdate() {
		fullPath := ftutilities.Join(prefix, update.GetPath())
		physIntf, attachmentPoint, ok := parseIntfPairKey(fullPath)
		if !ok {
			continue
		}

		apInfo := i.state.createOrRetrieveAttachmentPoint(attachmentPoint)
		elems := update.GetPath().GetElem()
		if len(elems) == 0 {
			continue
		}
		lastElem := elems[len(elems)-1].GetName()

		// Send counter update to the cache manager
		if apInfo.setMemberCounter(physIntf, lastElem, update.GetVal().GetUintVal()) {
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

		remainingMembers, ok := i.state.removeMemberAndCleanup(attachmentPoint, physIntf)
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
		if info, ok := i.state.attachmentPointInfo(ap); ok {
			byteDrop, pktDrop, yellowByte, yellowPkt := info.aggregateCounters()

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
