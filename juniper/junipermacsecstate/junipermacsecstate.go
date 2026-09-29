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

// Package junipermacsecstate translates MACsec interface state and MKA state from Juniper native to OpenConfig.
package junipermacsecstate

import (
	"fmt"
	"sort"
	"strings"

	log "github.com/golang/glog"
	"github.com/openconfig/functional-translators/ftconsts"
	"github.com/openconfig/functional-translators/ftutilities"
	"github.com/openconfig/functional-translators/translator"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

var (
	// Translation map from OpenConfig MACsec state paths to Juniper native paths.
	// Juniper does not support `*` subscription for the native paths.
	// Therefore, we need to subscribe to the longest prefix/container of a path.
	// Example:
	// for native path: /junos/macsec/interfaces/interface[name=<intf>]/mka/state/jnx-cak-name
	// Subscribe to: /junos/macsec/interfaces/interface[name=<intf>]/mka/state

	translateMap = map[string][]string{
		"/openconfig/macsec/interfaces/interface/state/status": {
			"/junos/macsec/interfaces/interface/name",
		},
		"/openconfig/macsec/interfaces/interface/state/ckn": {
			"/junos/macsec/interfaces/interface/mka/state/jnx-cak-name",
		},
	}

	paths        = ftutilities.MustStringMapPaths(translateMap)
	pathPatterns = []*gnmipb.Path{
		// Interface name (presence indicates MACsec is enabled)
		{
			Origin: "junos",
			Elem: []*gnmipb.PathElem{
				{Name: "macsec"}, {Name: "interfaces"},
				{Name: "interface", Key: map[string]string{"name": "*"}},
				{Name: "name"},
			},
		},
		// CKN value from MKA state
		{
			Origin: "junos",
			Elem: []*gnmipb.PathElem{
				{Name: "macsec"}, {Name: "interfaces"},
				{Name: "interface", Key: map[string]string{"name": "*"}},
				{Name: "mka"}, {Name: "state"},
				{Name: "jnx-cak-name"},
			},
		},
	}

	deletePathPatterns = []*gnmipb.Path{
		{
			Origin: "junos",
			Elem: []*gnmipb.PathElem{
				{Name: "macsec"}, {Name: "interfaces"},
				{Name: "interface", Key: map[string]string{"name": "*"}},
			},
		},
	}
)

type cknInfo struct {
	principal    bool
	success      bool
	principalSet bool
	successSet   bool
}

type interfaceMacSecInfo struct {
	interfaceName string
	cpStatus      bool
	cpStatusSet   bool
	cknStatuses   map[string]*cknInfo
}

func (i *interfaceMacSecInfo) createOrGetCKN(ckn string) *cknInfo {
	if i.cknStatuses == nil {
		i.cknStatuses = make(map[string]*cknInfo)
	}
	if _, ok := i.cknStatuses[ckn]; !ok {
		i.cknStatuses[ckn] = new(cknInfo)
	}
	return i.cknStatuses[ckn]
}

func (c *cknInfo) isComplete() bool {
	return c != nil && c.principalSet && c.successSet
}

type macsecState struct {
	interfaces map[string]*interfaceMacSecInfo
}

func newMacsecState() *macsecState {
	return &macsecState{
		interfaces: make(map[string]*interfaceMacSecInfo),
	}
}

func (m *macsecState) createOrGetInterface(interfaceName string) *interfaceMacSecInfo {
	if _, ok := m.interfaces[interfaceName]; !ok {
		m.interfaces[interfaceName] = &interfaceMacSecInfo{
			interfaceName: interfaceName,
			cknStatuses:   make(map[string]*cknInfo),
		}
	}
	return m.interfaces[interfaceName]
}

func (m *macsecState) interfaceInfo(intf string) (*interfaceMacSecInfo, bool) {
	info, ok := m.interfaces[intf]
	return info, ok
}

func (m *macsecState) clearInterfaceInfo(intf string) {
	delete(m.interfaces, intf)
}

type impl struct {
	state *macsecState
}

// New creates a functional translator for Juniper MACsec state.
func New() *translator.FunctionalTranslator {
	i := &impl{
		state: newMacsecState(),
	}
	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               ftconsts.JuniperMacsecStateFunctionalTranslator,
			Translate:        i.translate,
			OutputToInputMap: paths,
			Metadata: []*translator.FTMetadata{
				{
					Vendor: ftconsts.VendorJuniper,
					// TODO: Consider whether software image conditions are needed for non-OC paths.
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("Failed to create Juniper MACsec state functional translator: %v", err)
	}
	return ft
}

// interfaceIDAndValue extracts the interface name from the key and the leaf name from the path.
func interfaceIDAndValue(path *gnmipb.Path) (intfID, counterName string, err error) {
	if len(path.GetElem()) < 2 {
		return "", "", fmt.Errorf("path %v has fewer than 2 elements", path)
	}

	// Find interface element and extract name from key
	var intfElem *gnmipb.PathElem
	for i, elem := range path.GetElem() {
		if elem.GetName() == "interface" && i > 0 && path.GetElem()[i-1].GetName() == "interfaces" {
			intfElem = elem
			break
		}
	}

	if intfElem == nil {
		return "", "", fmt.Errorf("interface element not found in path %v", path)
	}

	intfID = intfElem.GetKey()["name"]
	if intfID == "" {
		return "", "", fmt.Errorf("interface name key not found in path %v", path)
	}
	counterName = path.GetElem()[len(path.GetElem())-1].GetName()

	return intfID, counterName, nil
}

// deleteInfo holds extracted information from a delete path.
type deleteInfo struct {
	intfID     string
	deleteType string
}

// extractDeleteInfo returns delete information from the delete path.
func extractDeleteInfo(path *gnmipb.Path) (*deleteInfo, error) {
	lastElemIndex := len(path.GetElem()) - 1

	for _, pattern := range deletePathPatterns {
		if ftutilities.MatchPath(path, pattern) {
			// Interface level delete
			if path.GetElem()[lastElemIndex].GetName() == "interface" {
				intfID, _, err := interfaceIDAndValue(path)
				if err != nil {
					return nil, err
				}
				return &deleteInfo{
					intfID:     intfID,
					deleteType: "interface",
				}, nil
			}
		}
	}

	return nil, nil
}

// isKnownLeaf reports whether leaf is a leaf name that this translator processes.
// Used for early-exit before expensive path joining and matching.
func isKnownLeaf(leaf string) bool {
	switch leaf {
	case "name", "jnx-cak-name":
		return true
	}
	return false
}

// metadata populates the MACsec map with the native paths that contribute to the derived MACsec status.
func (i *impl) metadata(prefix *gnmipb.Path, update *gnmipb.Update) (string, error) {
	// O(1) early exit: skip updates whose leaf name cannot match any pathPattern.
	updateElems := update.GetPath().GetElem()
	if len(updateElems) == 0 {
		return "", nil
	}
	if !isKnownLeaf(updateElems[len(updateElems)-1].GetName()) {
		return "", nil
	}

	fullPath := ftutilities.Join(prefix, update.GetPath())
	matched := false
	interfaceName, counterName, err := interfaceIDAndValue(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to get interface ID from path %v: %v", fullPath, err)
	}
	for _, pattern := range pathPatterns {
		if ftutilities.MatchPath(fullPath, pattern) {
			matched = true
			ifaceInfo := i.state.createOrGetInterface(interfaceName)

			switch counterName {
			case "name":
				// Presence of the interface name in macsec tree indicates MACsec is enabled
				ifaceInfo.cpStatus = true
				ifaceInfo.cpStatusSet = true
			case "jnx-cak-name":
				// Extract the actual CKN value from the update
				cknVal := strings.ToLower(update.GetVal().GetStringVal())
				if cknVal == "" {
					log.Warningf("Empty jnx-cak-name value for interface %q on path %v", interfaceName, fullPath)
					break
				}
				cknInfo := ifaceInfo.createOrGetCKN(cknVal)
				cknInfo.success = true
				cknInfo.successSet = true
				cknInfo.principal = true
				cknInfo.principalSet = true
			default:
				log.Warningf("Unknown counter name %q encountered for MACsec status processing for path %v", counterName, fullPath)
			}
		}
	}
	if !matched {
		return "", nil
	}
	return interfaceName, nil
}

// translateMACSecState returns the MACsec status and ckn for the given interface.
func (i *impl) translateMACSecState(interfaceName string) (intfMACSecStatuses, cknKeys []string, skip bool) {
	ifaceInfo, ok := i.state.interfaceInfo(interfaceName)
	if !ok {
		log.V(1).Infof("interface %q not found for status translation.", interfaceName)
		return nil, nil, true
	}
	if !ifaceInfo.cpStatusSet {
		log.V(1).Infof("cpStatusSet is false for interface %q. Returning empty CKN and status list.", interfaceName)
		return nil, nil, true
	}

	if len(ifaceInfo.cknStatuses) == 0 {
		log.V(1).Infof("no CKNs found for interface %q. Returning empty CKN and status list.", interfaceName)
		return nil, nil, true
	}
	cknNamesToSort := make([]string, 0, len(ifaceInfo.cknStatuses))
	for ckn := range ifaceInfo.cknStatuses {
		cknNamesToSort = append(cknNamesToSort, ckn)
	}
	if len(cknNamesToSort) == 0 {
		return nil, nil, true
	}
	sort.Strings(cknNamesToSort)
	for _, c := range cknNamesToSort {
		cknKeys = append(cknKeys, c)
		cknStatus := "Unknown"

		cknInfo := ifaceInfo.cknStatuses[c]
		if !cknInfo.isComplete() {
			log.V(1).Infof("macsec data for interface %q (CKN %q) is not yet complete. PrincipalSet: %t, SuccessSet: %t",
				interfaceName, c, cknInfo.principalSet, cknInfo.successSet)
			return nil, nil, true
		}

		controlledPortEnabled := ifaceInfo.cpStatus
		success := cknInfo.success
		principal := cknInfo.principal

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
		intfMACSecStatuses = append(intfMACSecStatuses, cknStatus)
	}
	return intfMACSecStatuses, cknKeys, false
}

// returnPathForMACSecStatus returns OpenConfig path for MACsec status of the given interface.
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

// returnPathForMACSecCKN returns OpenConfig path for MACsec CKN of the given interface.
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

// deleteHandler processes delete operations from the native notification and returns
// a map of interface names that should be deleted in OpenConfig.
func (i *impl) deleteHandler(n *gnmipb.Notification) map[string]bool {
	interfacesForDelete := make(map[string]bool)
	prefix := n.GetPrefix()

	for _, delPath := range n.GetDelete() {
		fullPath := ftutilities.Join(prefix, delPath)
		delInfo, err := extractDeleteInfo(fullPath)
		if err != nil {
			log.Warningf("Failed to extract delete info from path %v: %v", fullPath, err)
			continue
		}

		if delInfo != nil && delInfo.deleteType == "interface" {
			interfacesForDelete[delInfo.intfID] = true
		}
	}

	return interfacesForDelete
}

func (i *impl) translate(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
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
		interfaceName, err := i.metadata(prefix, update)
		if err != nil {
			return nil, fmt.Errorf("failed to populate MACsec map: %v", err)
		}
		if interfaceName != "" {
			interfaceSeen[interfaceName] = true
		}
	}

	finalInterfacesForOCUpdate := make(map[string]bool)
	// Determine final set of interfaces for OC Update
	interfacesForOCDelete := i.deleteHandler(notification)
	for intfName := range interfaceSeen {
		if !interfacesForOCDelete[intfName] {
			finalInterfacesForOCUpdate[intfName] = true
		}
	}

	// Evict cached state for deleted interfaces to prevent stale data.
	if len(interfacesForOCDelete) > 0 {
		for intfName := range interfacesForOCDelete {
			i.state.clearInterfaceInfo(intfName)
		}
	}

	// Generate final set of deletes
	for intfName := range interfacesForOCDelete {
		outgoingDeletes = append(outgoingDeletes, returnPathForMACSecStatus(intfName), returnPathForMACSecCKN(intfName))
	}

	for interfaceName := range finalInterfacesForOCUpdate {
		intfMACSecStatuses, ckns, skip := i.translateMACSecState(interfaceName)
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
		log.V(1).Infof("target %s: Generating OC Update for status (len %d) and CKN (len %d) leaf-lists for interface %q.", target, len(intfMACSecStatuses), len(ckns), interfaceName)
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
