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

// Package aristapwstate implements the functional translator for Arista pseudowire state.
package aristapwstate

import (
	"fmt"

	log "github.com/golang/glog"
	"github.com/openconfig/functional-translators/ftconsts"
	"github.com/openconfig/functional-translators/ftutilities"
	"github.com/openconfig/functional-translators/translator"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

const (
	patchNameIndex       = 5
	minDeletePathElemLen = 6
)

var (
	translateMap = map[string][]string{
		"/openconfig/network-instances/network-instance/connection-points/connection-point/state/status": {
			"/eos_native/Sysdb/pseudowire/agent/pseudowirePatchInfoColl/patchInfo",
		},
	}
	updatePathPatterns = []*gnmipb.Path{
		{
			Origin: "eos_native",
			Elem: []*gnmipb.PathElem{
				{Name: "Sysdb"}, {Name: "pseudowire"}, {Name: "agent"}, {Name: "pseudowirePatchInfoColl"},
				{Name: "patchInfo"},
				{Name: "*"}, // patchName
				{Name: "patchState"},
			},
		},
	}
	deletePathPatterns = []*gnmipb.Path{
		{
			Origin: "eos_native",
			Elem: []*gnmipb.PathElem{
				{Name: "Sysdb"}, {Name: "pseudowire"}, {Name: "agent"}, {Name: "pseudowirePatchInfoColl"},
				{Name: "patchInfo"},
				{Name: "*"}, // patchName
			},
		},
	}
	updateValueMap = map[string]string{
		"patchUp":   "UP",
		"patchDown": "DOWN",
	}
)

// New creates a functional translator.
func New() *translator.FunctionalTranslator {
	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               ftconsts.AristaPWStateFunctionalTranslator,
			Translate:        translate,
			OutputToInputMap: ftutilities.MustStringMapPaths(translateMap),
			Metadata: []*translator.FTMetadata{
				{
					Vendor: ftconsts.VendorArista,
					SoftwareVersionRange: &translator.SWRange{
						InclusiveMin: "4.33.0F",
						/* TODO: bhageshbhutani - Get Vendor confirmation of the support for the below version
						for the required OC path starting from the below version.
						https://b.corp.google.com/issues/433946399#comment28
						*/
						ExclusiveMax: "4.36",
					},
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("Failed to create Arista PW state functional translator: %v", err)
	}
	return ft
}
func updateHandler(n *gnmipb.Notification) ([]*gnmipb.Update, error) {
	if len(n.GetUpdate()) == 0 {
		return nil, nil
	}
	prefix := n.GetPrefix()
	var updates []*gnmipb.Update
	for _, update := range n.GetUpdate() {
		fullPath := ftutilities.Join(prefix, update.GetPath())
		for _, path := range updatePathPatterns {
			if !ftutilities.MatchPath(fullPath, path) {
				continue
			}
			elems := fullPath.GetElem()
			patchName := elems[patchNameIndex].GetName()
			val, ok := updateValueMap[update.GetVal().GetStringVal()]
			if !ok {
				return nil, fmt.Errorf("path %v has unexpected value: %q", fullPath, update.GetVal().GetStringVal())
			}
			outgoingUpdate := &gnmipb.Update{
				Path: returnPath(patchName),
				Val: &gnmipb.TypedValue{
					Value: &gnmipb.TypedValue_StringVal{
						StringVal: val,
					},
				},
			}
			updates = append(updates, outgoingUpdate)
			break
		}
	}
	return updates, nil
}

// deleteHandler returns the deletes that should be sent to the target.
func deleteHandler(n *gnmipb.Notification) ([]*gnmipb.Path, error) {
	prefix := n.GetPrefix()
	var deletes []*gnmipb.Path
	for _, del := range n.GetDelete() {
		fullPath := ftutilities.Join(prefix, del)
		elems := fullPath.GetElem()
		for _, pattern := range deletePathPatterns {
			if ftutilities.MatchPath(fullPath, pattern) {
				patchName := elems[patchNameIndex].GetName()
				deletes = append(deletes, returnPath(patchName))
				break
			}
		}
	}
	return deletes, nil
}

// returnPath returns a gNMI path for the update.
// Does not set the origin or the target.
func returnPath(patchName string) *gnmipb.Path {
	return &gnmipb.Path{
		Elem: []*gnmipb.PathElem{
			{Name: "network-instances"},
			{
				Name: "network-instance",
				Key: map[string]string{
					// Arista uses "default" as the network-instance name instead of "DEFAULT", which is an unexpected deviation.
					"name": "default",
				},
			},
			{Name: "connection-points"},
			{
				Name: "connection-point",
				Key: map[string]string{
					"connection-point-id": patchName,
				},
			},
			{Name: "state"},
			{Name: "status"},
		},
	}
}
func translate(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
	notification := sr.GetUpdate()
	if notification == nil {
		return nil, nil
	}
	updates, err := updateHandler(notification)
	if err != nil {
		return nil, err
	}
	deletes, err := deleteHandler(notification)
	if err != nil {
		return nil, err
	}
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
