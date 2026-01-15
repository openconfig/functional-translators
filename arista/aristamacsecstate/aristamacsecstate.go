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

// Package aristamacsecstate translates the interface MACSec state from native to openconfig.
package aristamacsecstate

import (
	"fmt"
	"maps"
	"sort"

	log "github.com/golang/glog"
	"github.com/openconfig/functional-translators/ftconsts"
	"github.com/openconfig/functional-translators/ftutilities"
	"github.com/openconfig/functional-translators/translator"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

var (
	// Arista does not support `*` subscription for the native paths.
	// Therefore, we need to subscribe to the longest prefix/container of a path.
	// Example:
	// for native path: /eos_native/Sysdb/macsec/status/cpStatus/<interface-id>/controlledPortEnabled
	// Subscribe to: /eos_native/Sysdb/macsec/status/cpStatus
	translateMap = map[string][]string{
		"/openconfig/macsec/interfaces/interface/state/status": {
			"/eos_native/Sysdb/macsec/status/cpStatus",
			"/eos_native/Sysdb/macsec/mkaStatus/portStatus",
		},
		"/openconfig/macsec/interfaces/interface/state/ckn": {
			"/eos_native/Sysdb/macsec/mkaStatus/portStatus",
		},
	}
	paths        = ftutilities.MustStringMapPaths(translateMap)
	pathPatterns = []*gnmipb.Path{
		{
			Origin: "eos_native",
			Elem: []*gnmipb.PathElem{
				{Name: "Sysdb"}, {Name: "macsec"}, {Name: "status"}, {Name: "cpStatus"},
				{Name: "*"}, // interface-id
				{Name: "controlledPortEnabled"},
			},
		},
		{
			Origin: "eos_native",
			Elem: []*gnmipb.PathElem{
				{Name: "Sysdb"}, {Name: "macsec"}, {Name: "mkaStatus"}, {Name: "portStatus"},
				{Name: "*"}, // interface-id
				{Name: "actorStatus"},
				{Name: "*"}, // CKN
				{Name: "success"},
			},
		},
		{
			Origin: "eos_native",
			Elem: []*gnmipb.PathElem{
				{Name: "Sysdb"}, {Name: "macsec"}, {Name: "mkaStatus"}, {Name: "portStatus"},
				{Name: "*"}, // interface-id
				{Name: "actorStatus"},
				{Name: "*"}, // CKN
				{Name: "principal"},
			},
		},
	}
	deletePathPatterns = []*gnmipb.Path{
		{
			Origin: "eos_native",
			Elem: []*gnmipb.PathElem{
				{Name: "Sysdb"}, {Name: "macsec"}, {Name: "status"}, {Name: "cpStatus"},
				{Name: "*"}, // interface-id
			},
		},
		{
			Origin: "eos_native",
			Elem: []*gnmipb.PathElem{
				{Name: "Sysdb"}, {Name: "macsec"}, {Name: "mkaStatus"}, {Name: "portStatus"},
				{Name: "*"}, // interface-id
				{Name: "actorStatus"},
				{Name: "*"}, // CKN
			},
		},
		{
			Origin: "eos_native",
			Elem: []*gnmipb.PathElem{
				{Name: "Sysdb"}, {Name: "macsec"}, {Name: "mkaStatus"}, {Name: "portStatus"},
				{Name: "*"}, // interface-id,
			},
		},
	}
)

// New creates a functional translator.
func New() *translator.FunctionalTranslator {
	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               ftconsts.AristaMacsecStateFunctionalTranslator,
			Translate:        translate,
			OutputToInputMap: paths,
			Metadata: []*translator.FTMetadata{
				{
					Vendor: ftconsts.VendorArista,
					SoftwareVersionRange: &translator.SWRange{
						InclusiveMin: "4.33.0F",
						ExclusiveMax: "4.36",
					},
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("Failed to create Arista MACsec state functional translator: %v", err)
	}
	return ft
}

// interfaceIDAndCKN returns the interface ID and CKN from the path for an update.
// If the path is for controlledPortEnabled, then the interface ID is the second last element and CKN is empty.
// If the path is for success or principal, then the interface ID is the fourth last element and CKN is the second last element.
func interfaceIDAndCKN(path *gnmipb.Path) (intfID, ckn string, err error) {
	if len(path.GetElem()) < 6 {
		return "", "", fmt.Errorf("path %v has fewer than 6 elements", path)
	}
	lastElemIndex := len(path.GetElem()) - 1
	lastElemName := path.GetElem()[lastElemIndex].GetName()
	if lastElemName == "controlledPortEnabled" {
		return path.GetElem()[lastElemIndex-1].GetName(), "", nil
	}
	if lastElemName == "success" || lastElemName == "principal" {
		return path.GetElem()[lastElemIndex-3].GetName(), path.GetElem()[lastElemIndex-1].GetName(), nil
	}
	return "", "", nil
}

// deleteInfo holds extracted information from a delete path.
type deleteInfo struct {
	intfID     string
	ckn        string
	deleteType string
}

// extractDeleteInfo returns the interface ID or CKN (if applicable) from the delete path.
// Returns: deleteInfo, error.
// The last element is the interface ID or the CKN.
func extractDeleteInfo(path *gnmipb.Path) (*deleteInfo, error) {
	if len(path.GetElem()) != 5 && len(path.GetElem()) != 7 {
		return nil, fmt.Errorf("received delete path %v with unexpected number of %d elements", path, len(path.GetElem()))
	}
	matched := false
	for _, pattern := range deletePathPatterns {
		if ftutilities.MatchPath(path, pattern) {
			matched = true
			break
		}
	}
	if !matched {
		return nil, nil
	}
	lastElemIndex := len(path.GetElem()) - 1
	// CKN level delete:
	// Sysdb/macsec/mkaStatus/portStatus/<interface-id>/actorStatus/<CKN>
	if path.GetElem()[lastElemIndex-1].GetName() == "actorStatus" {
		return &deleteInfo{
			intfID:     path.GetElem()[lastElemIndex-2].GetName(), // <interface-id>
			ckn:        path.GetElem()[lastElemIndex].GetName(),   // <CKN>
			deleteType: "ckn-delete",
		}, nil
	}
	// Interface level delete for portStatus:
	// Sysdb/macsec/mkaStatus/portStatus/<interface-id>
	if path.GetElem()[lastElemIndex-1].GetName() == "portStatus" {
		return &deleteInfo{
			intfID:     path.GetElem()[lastElemIndex].GetName(), // <interface-id>
			deleteType: "intf-delete",
		}, nil
	}
	// Interface level delete for all CKNs:
	// Sysdb/macsec/status/cpStatus/<interface-id>
	if path.GetElem()[lastElemIndex-1].GetName() == "cpStatus" {
		return &deleteInfo{
			intfID:     path.GetElem()[lastElemIndex].GetName(), // <interface-id>
			deleteType: "cpStatus-delete",
		}, nil
	}
	return nil, fmt.Errorf("delete path %v matched a pattern but has an unrecognized structure", path)
}

// deleteHandler updates the cache based on delete notifications.
// It returns a map of interfaces that need an OC Delete, and a map of interfaces that need an OC Update.
func deleteHandler(n *gnmipb.Notification) (interfacesForOCDelete, interfacesForOCUpdate map[string]bool) {
	prefix := n.GetPrefix()
	deletes := n.GetDelete()
	target := prefix.GetTarget()

	interfacesForOCDelete = make(map[string]bool)
	interfacesForOCUpdate = make(map[string]bool)

	for _, del := range deletes {
		fullPath := ftutilities.Join(prefix, del)
		deleteInfo, err := extractDeleteInfo(fullPath)
		if err != nil {
			log.Errorf("failed to extract interface ID or CKN from delete path %v after matching a pattern: %v", fullPath, err)
			continue
		}
		if deleteInfo == nil {
			log.V(1).Infof("delete path %v did not match known MACsec delete patterns.", fullPath)
			continue
		}

		// Check if the map for the specific target exists
		if targetInfo, targetExists := ftutilities.AristaMACSecMap.RetrieveTargetMacSecInfo(target); targetExists {
			ifaceInfo, ok := targetInfo.InterfaceInfo(deleteInfo.intfID)
			if !ok {
				log.V(1).Infof("interface '%s' on target '%s' not found for delete handler.", deleteInfo.intfID, target)
				continue
			}
			switch deleteInfo.deleteType {
			case "intf-delete":
				targetInfo.ClearInterfaceInfo(deleteInfo.intfID)
				interfacesForOCDelete[deleteInfo.intfID] = true
				if len(targetInfo.Interfaces) == 0 {
					log.V(1).Infof("no more interfaces for target '%s', removing target from map.", target)
					ftutilities.AristaMACSecMap.DeleteTargetMacSecInfo(target)
				}
			case "ckn-delete":
				ifaceInfo.RemoveCkn(deleteInfo.ckn)
				interfacesForOCUpdate[deleteInfo.intfID] = true
				if len(ifaceInfo.CloneStatuses()) == 0 {
					log.V(1).Infof("no more CKNs for interface '%s' on target '%s', removing interface from map.", deleteInfo.intfID, target)
					targetInfo.ClearInterfaceInfo(deleteInfo.intfID)
					interfacesForOCDelete[deleteInfo.intfID] = true
				}
			case "cpStatus-delete":
				ifaceInfo.ResetCPStatus()
				interfacesForOCDelete[deleteInfo.intfID] = true
			}
		}
	}
	return interfacesForOCDelete, interfacesForOCUpdate
}

// returnPathForMACSecStatus returns a gNMI path for MACSec status of the given interface.
// Does not set the origin or the target.
func returnPathForMACSecStatus(interfaceName string) *gnmipb.Path {
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
			{Name: "status"},
		},
	}
}

// returnPathForMACSecCKN returns a gNMI path for MACSec CKN of the given interface.
// Does not set the origin or the target.
func returnPathForMACSecCKN(interfaceName string) *gnmipb.Path {
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
			{Name: "ckn"},
		},
	}
}

// metadata populates the MACSec map with the native paths that contribute to the derived MACSec status.
func metadata(prefix *gnmipb.Path, update *gnmipb.Update, target string) (string, error) {
	fullPath := ftutilities.Join(prefix, update.GetPath())
	matched := false
	interfaceName, ckn, err := interfaceIDAndCKN(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to get interface ID and CKN from path %v: %v", fullPath, err)
	}
	for _, pattern := range pathPatterns {
		if ftutilities.MatchPath(fullPath, pattern) {
			matched = true
			targetInfo := ftutilities.AristaMACSecMap.CreateOrUpdateTargetMacSecInfo(target)
			ifaceInfo := targetInfo.CreateOrGetInterface(interfaceName)

			leafName := fullPath.GetElem()[len(fullPath.GetElem())-1].GetName()
			boolVal := update.GetVal().GetBoolVal()

			switch leafName {
			case "controlledPortEnabled":
				ifaceInfo.SetIntfCPStatus(boolVal)
			case "success":
				ifaceInfo.SetIntfSuccess(ckn, boolVal)
			case "principal":
				ifaceInfo.SetIntfPrincipal(ckn, boolVal)
			default:
				log.Warningf("Unknown leaf name '%s' encountered for MACsec status processing for path %v", leafName, fullPath)
			}
		}
	}
	if !matched {
		return "", nil
	}
	return interfaceName, nil
}

// translateMACSecState returns the MACSec ckn and status for the given interface.
func translateMACSecState(interfaceName string, target string) (intfMACSecStatus, cknKeys []string, skip bool) {
	var success, principal bool
	targetInfo, ok := ftutilities.AristaMACSecMap.RetrieveTargetMacSecInfo(target)
	if !ok {
		log.V(1).Infof("target '%s' not found in AristaMACSecMap for status translation.", target)
		return nil, nil, true
	}

	ifaceInfo, ok := targetInfo.InterfaceInfo(interfaceName)
	if !ok {
		log.V(1).Infof("interface '%s' on target '%s' not found for status translation.", interfaceName, target)
		return nil, nil, true
	}
	controlledPortEnabled, cpStatusSet := ifaceInfo.IntfCPStatus()
	if !cpStatusSet {
		log.V(1).Infof("cpStatusSet is false for interface '%s' on target '%s'. Returning empty CKN and status list.", interfaceName, target)
		return nil, nil, true
	}

	if len(ifaceInfo.CloneStatuses()) == 0 {
		log.V(1).Infof("no CKNs found for interface '%s' on target '%s'. Returning empty CKN and status list.", interfaceName, target)
		return nil, nil, true
	}
	var cknNamesToSort []string
	for ckn := range maps.Keys(ifaceInfo.CloneStatuses()) {
		cknNamesToSort = append(cknNamesToSort, ckn)
	}
	if len(cknNamesToSort) == 0 {
		return nil, nil, true
	}
	sort.Strings(cknNamesToSort)
	for _, c := range cknNamesToSort {
		cknKeys = append(cknKeys, c)
		cknStatus := "Unknown"

		success, _ = ifaceInfo.IntfSuccess(c)
		principal, _ = ifaceInfo.IntfPrincipal(c)

		// Check if all required values have been set using the IsComplete method.
		// If not, we can't determine a definitive status yet.
		if !ifaceInfo.IsComplete(c) {
			log.V(1).Infof("macsec data for interface '%s' on target '%s' is not yet complete. CPStatusSet: %t, PrincipalSet: %t, SuccessSet: %t",
				interfaceName, target, controlledPortEnabled, principal, success)
			return nil, nil, true
		}

		switch {
		case controlledPortEnabled && success && principal:
			cknStatus = "Secured"
		case controlledPortEnabled && !success && !principal:
			cknStatus = "Unencrypted Allowed"
		case !controlledPortEnabled && !success && !principal:
			cknStatus = "Unencrypted Dropped"
		default:
			// This should not happen if macsec is working as expected and the native paths are correctly populated.
			cknStatus = "Unknown"
		}
		intfMACSecStatus = append(intfMACSecStatus, cknStatus)
	}
	return intfMACSecStatus, cknKeys, false
}

func translate(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
	if sr.GetUpdate() == nil {
		return nil, nil
	}
	notification := sr.GetUpdate()
	prefix := notification.GetPrefix()
	target := prefix.GetTarget()
	var outgoingNotification *gnmipb.Notification
	var outgoingUpdates []*gnmipb.Update
	var outgoingDeletes []*gnmipb.Path
	interfaceSeen := make(map[string]bool)

	for _, update := range notification.GetUpdate() {
		interfaceName, err := metadata(prefix, update, target)
		if err != nil {
			return nil, fmt.Errorf("failed to populate MACSec map: %v", err)
		}
		if interfaceName != "" {
			interfaceSeen[interfaceName] = true
		}
	}
	finalInterfacesForOCUpdate := make(map[string]bool)
	// Determine final set of interfaces for OC Update
	// These are interfaces affected by native updates or "modifying" native deletes,
	interfacesForOCDelete, interfacesForOCUpdate := deleteHandler(notification)
	for intfName := range interfaceSeen {
		if !interfacesForOCDelete[intfName] {
			finalInterfacesForOCUpdate[intfName] = true
		}
	}
	for intfName := range interfacesForOCUpdate {
		if !interfacesForOCDelete[intfName] {
			finalInterfacesForOCUpdate[intfName] = true
		}
	}
	// Generate final set of deletes
	for intfName := range interfacesForOCDelete {
		outgoingDeletes = append(outgoingDeletes, returnPathForMACSecStatus(intfName), returnPathForMACSecCKN(intfName))
	}
	for interfaceName := range finalInterfacesForOCUpdate {
		intfMACSecStatuses, ckns, skip := translateMACSecState(interfaceName, target)
		if skip {
			continue
		}
		var statusElements []*gnmipb.TypedValue
		for _, statusStr := range intfMACSecStatuses {
			statusElements = append(statusElements, &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: statusStr}})
		}
		statusUpdate := &gnmipb.Update{
			Path: returnPathForMACSecStatus(interfaceName),
			Val:  &gnmipb.TypedValue{Value: &gnmipb.TypedValue_LeaflistVal{LeaflistVal: &gnmipb.ScalarArray{Element: statusElements}}},
		}
		outgoingUpdates = append(outgoingUpdates, statusUpdate)

		var cknElements []*gnmipb.TypedValue
		for _, cknNameStr := range ckns {
			cknElements = append(cknElements, &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: cknNameStr}})
		}
		cknUpdate := &gnmipb.Update{
			Path: returnPathForMACSecCKN(interfaceName),
			Val:  &gnmipb.TypedValue{Value: &gnmipb.TypedValue_LeaflistVal{LeaflistVal: &gnmipb.ScalarArray{Element: cknElements}}},
		}
		outgoingUpdates = append(outgoingUpdates, cknUpdate)
		log.V(1).Infof("target %s: Generating OC Update for status (len %d) and CKN (len %d) leaf-lists for interface '%s'.", target, len(intfMACSecStatuses), len(ckns), interfaceName)
	}

	if len(outgoingUpdates) == 0 && len(outgoingDeletes) == 0 {
		return nil, nil
	}

	outgoingNotification = &gnmipb.Notification{
		Timestamp: notification.GetTimestamp(),
		Prefix:    &gnmipb.Path{Origin: "openconfig", Target: target},
		Update:    outgoingUpdates,
		Delete:    outgoingDeletes,
	}

	return &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{Update: outgoingNotification},
	}, nil
}
