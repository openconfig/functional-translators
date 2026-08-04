// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aristainterface

import (
	log "github.com/golang/glog"
	"github.com/openconfig/functional-translators/ftconsts"
	"github.com/openconfig/functional-translators/ftutilities"
	"github.com/openconfig/functional-translators/translator"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

const (
	vlanNameIndex = 6
)

var (
	translateMap = map[string][]string{
		"/openconfig/interfaces/interface/state/admin-status": {
			"/openconfig/interfaces/interface/state/admin-status",
			"/eos_native/Sysdb/interface/config/eth/vlan/intfConfig",
		},
		"/openconfig/interfaces/interface/state/oper-status": {
			"/openconfig/interfaces/interface/state/oper-status",
			"/eos_native/Sysdb/interface/status/eth/vlan/intfStatus",
		},
	}
	updatePathPatterns = map[string][]*gnmipb.Path{
		"admin-status": {
			{
				Origin: "openconfig",
				Elem: []*gnmipb.PathElem{
					{Name: "interfaces"},
					{Name: "interface", Key: map[string]string{"name": "*"}},
					{Name: "state"},
					{Name: "admin-status"},
				},
			},
			{
				Origin: "eos_native",
				Elem: []*gnmipb.PathElem{
					{Name: "Sysdb"}, {Name: "interface"}, {Name: "config"}, {Name: "eth"}, {Name: "vlan"},
					{Name: "intfConfig"},
					{Name: "*"}, // Vlan name
					{Name: "adminEnabledStateLocal"},
				},
			},
		},
		"oper-status": {
			{
				Origin: "openconfig",
				Elem: []*gnmipb.PathElem{
					{Name: "interfaces"},
					{Name: "interface", Key: map[string]string{"name": "*"}},
					{Name: "state"},
					{Name: "oper-status"},
				},
			},
			{
				Origin: "eos_native",
				Elem: []*gnmipb.PathElem{
					{Name: "Sysdb"}, {Name: "interface"}, {Name: "status"}, {Name: "eth"}, {Name: "vlan"},
					{Name: "intfStatus"},
					{Name: "*"}, // Vlan name
					{Name: "operStatus"},
				},
			},
		},
	}
	updateValueMap = map[string]map[string]string{
		"admin-status": {
			"enabled":  "UP",
			"shutdown": "DOWN",
			// Arista treats "unknownEnabledState" as effectively UP for SVIs when not explicitly shutdown.
			"unknownEnabledState": "UP",
		},
		"oper-status": {
			"intfOperDormant":        "DORMANT",
			"intfOperDown":           "DOWN",
			"intfOperLowerLayerDown": "LOWER_LAYER_DOWN",
			"intfOperNotPresent":     "NOT_PRESENT",
			"intfOperTesting":        "TESTING",
			"intfOperUnknown":        "UNKNOWN",
			"intfOperUp":             "UP",
		},
	}

	deletePathPatterns = map[string][]*gnmipb.Path{
		"admin-status": {
			{
				Origin: "eos_native",
				Elem: []*gnmipb.PathElem{
					{Name: "Sysdb"}, {Name: "interface"}, {Name: "config"}, {Name: "eth"}, {Name: "vlan"},
					{Name: "intfConfig"},
					{Name: "*"}, // Vlan name
				},
			},
			{
				Origin: "openconfig",
				Elem: []*gnmipb.PathElem{
					{Name: "interfaces"},
					{Name: "interface", Key: map[string]string{"name": "*"}},
					{Name: "state"},
					{Name: "admin-status"},
				},
			},
		},
		"oper-status": {
			{
				Origin: "eos_native",
				Elem: []*gnmipb.PathElem{
					{Name: "Sysdb"}, {Name: "interface"}, {Name: "status"}, {Name: "eth"}, {Name: "vlan"},
					{Name: "intfStatus"},
					{Name: "*"}, // Vlan name
				},
			},
			{
				Origin: "openconfig",
				Elem: []*gnmipb.PathElem{
					{Name: "interfaces"},
					{Name: "interface", Key: map[string]string{"name": "*"}},
					{Name: "state"},
					{Name: "oper-status"},
				},
			},
		},
	}
)

// NewStatusFT creates a functional translator.
func NewStatusFT() *translator.FunctionalTranslator {
	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               ftconsts.AristaInterfaceStatusFunctionalTranslator,
			Translate:        translate,
			OutputToInputMap: ftutilities.MustStringMapPaths(translateMap),
			Metadata: []*translator.FTMetadata{
				{
					Vendor: ftconsts.VendorArista,
					SoftwareVersionRange: &translator.SWRange{
						InclusiveMin: "4.35.2F",
						// TODO: b/500399464
						ExclusiveMax: "4.37",
					},
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("Failed to create Arista routed-vlan interface state functional translator: %v", err)
	}
	return ft
}

// matchPath matches a path against a map of patterns and returns the status, origin, and whether a match was found.
func matchPath(path *gnmipb.Path, patterns map[string][]*gnmipb.Path) (string, string, bool) {
	for status, paths := range patterns {
		for _, pattern := range paths {
			if ftutilities.MatchPath(path, pattern) {
				return status, pattern.GetOrigin(), true
			}
		}
	}
	return "", "", false
}

func updateHandler(n *gnmipb.Notification) []*gnmipb.Update {
	if len(n.GetUpdate()) == 0 {
		return nil
	}
	prefix := n.GetPrefix()
	updates := make([]*gnmipb.Update, 0, len(n.GetUpdate()))
	for _, update := range n.GetUpdate() {
		fullPath := ftutilities.Join(prefix, update.GetPath())
		leaf, origin, ok := matchPath(fullPath, updatePathPatterns)
		if !ok {
			continue
		}
		// Passthrough for OC updates.
		if origin == "openconfig" {
			// Clear the Origin and Target from the fullPath as they will be set in the Notification Prefix.
			fullPath.Origin = ""
			fullPath.Target = ""
			updates = append(updates, &gnmipb.Update{
				Path: fullPath,
				Val:  update.GetVal(),
			})
			continue
		}
		// Translate vendor updates.
		elems := fullPath.GetElem()
		vlanName := elems[vlanNameIndex].GetName()
		val, ok := updateValueMap[leaf][update.GetVal().GetStringVal()]
		if !ok {
			log.Warningf("Path %v has unexpected value: %q, skipping this update.", fullPath, update.GetVal().GetStringVal())
			continue
		}
		outgoingUpdate := &gnmipb.Update{
			Path: returnPath(vlanName, leaf),
			Val: &gnmipb.TypedValue{
				Value: &gnmipb.TypedValue_StringVal{
					StringVal: val,
				},
			},
		}
		updates = append(updates, outgoingUpdate)
	}
	return updates
}

// deleteHandler returns the deletes that should be sent to the target.
func deleteHandler(n *gnmipb.Notification) []*gnmipb.Path {
	if len(n.GetDelete()) == 0 {
		return nil
	}
	prefix := n.GetPrefix()
	var deletes []*gnmipb.Path

	for _, del := range n.GetDelete() {
		fullPath := ftutilities.Join(prefix, del)
		status, origin, ok := matchPath(fullPath, deletePathPatterns)
		if !ok {
			continue
		}
		if origin == "openconfig" {
			// Clear the Origin and Target from the fullPath as they will be set in the Notification Prefix.
			fullPath.Origin = ""
			fullPath.Target = ""
			deletes = append(deletes, fullPath)
			continue
		}
		vlanName := fullPath.GetElem()[vlanNameIndex].GetName()
		deletes = append(deletes, returnPath(vlanName, status))
	}
	return deletes
}

// Does not set the origin or the target.
func returnPath(intfName string, status string) *gnmipb.Path {
	return &gnmipb.Path{
		Elem: []*gnmipb.PathElem{
			{Name: "interfaces"},
			{
				Name: "interface",
				Key: map[string]string{
					"name": intfName,
				},
			},
			{Name: "state"},
			{Name: status},
		},
	}
}

func translate(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
	notification := sr.GetUpdate()
	if notification == nil {
		return nil, nil
	}
	updates := updateHandler(notification)
	deletes := deleteHandler(notification)
	if len(updates) == 0 && len(deletes) == 0 {
		return nil, nil
	}
	outgoingNotification := &gnmipb.Notification{
		Timestamp: notification.GetTimestamp(),
		Prefix: &gnmipb.Path{
			Origin: "openconfig",
			Target: notification.GetPrefix().GetTarget(),
		},
		Update: updates,
		Delete: deletes,
	}
	return &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: outgoingNotification,
		},
	}, nil
}
