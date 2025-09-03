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

// Package aristacfmstate translates the CFM state data from native to openconfig.
package aristacfmstate

import (
	"fmt"

	log "github.com/golang/glog"
	"github.com/openconfig/functional-translators/ftconsts"
	"github.com/openconfig/functional-translators/ftutilities"
	"github.com/openconfig/functional-translators/translator"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

var (
	translateMap = map[string][]string{
		"/openconfig/oam/cfm/domains/maintenance-domain/maintenance-associations/maintenance-association/mep-endpoints/mep-endpoint/state/present-rdi": {
			"/eos_native/Sysdb/cfm/status/mdStatus",
		},
	}
	updatePathPatterns = []*gnmipb.Path{
		{
			Origin: "eos_native",
			Elem: []*gnmipb.PathElem{
				{Name: "Sysdb"}, {Name: "cfm"}, {Name: "status"}, {Name: "mdStatus"},
				{Name: "*"}, // domain
				{Name: "maStatus"},
				{Name: "*"}, // maName
				{Name: "localMepStatus"},
				{Name: "*"}, // localMepID
				{Name: "rdiTxCondition"},
				{Name: "rdi"},
			},
		},
	}
	deletePathPatterns = []*gnmipb.Path{
		{
			Origin: "eos_native",
			Elem: []*gnmipb.PathElem{
				{Name: "Sysdb"}, {Name: "cfm"}, {Name: "status"}, {Name: "mdStatus"},
				{Name: "*"}, // domain
				{Name: "maStatus"},
				{Name: "*"}, // maName
				{Name: "localMepStatus"},
				{Name: "*"}, // localMepID
			},
		},
	}
)

// New creates a functional translator.
func New() *translator.FunctionalTranslator {
	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               ftconsts.AristaCfmStateFunctionalTranslator,
			Translate:        translate,
			OutputToInputMap: ftutilities.MustStringMapPaths(translateMap),
			Metadata: []*translator.DeviceMetadata{
				{
					Vendor:          ftconsts.VendorArista,
					SoftwareVersion: "4.33.0F",
				},
				{
					Vendor:          ftconsts.VendorArista,
					SoftwareVersion: "4.33.0F-DPE",
				},
				{
					Vendor:          ftconsts.VendorArista,
					SoftwareVersion: "4.34.0F",
				},
				{
					Vendor:          ftconsts.VendorArista,
					SoftwareVersion: "4.34.0F-DPE",
				},
				{
					Vendor:          ftconsts.VendorArista,
					SoftwareVersion: "4.34.0F-DPE-42073180.narmadawbbrel (engineering build)",
				},
				{
					Vendor:          ftconsts.VendorArista,
					SoftwareVersion: "4.34.0F-DPE-41679138.narmadawbbrel (engineering build)",
				},
				{
					Vendor:          ftconsts.VendorArista,
					SoftwareVersion: "4.34.0F-DPE-41297507.narmadawbbrel (engineering build)",
				},
				{
					Vendor:          ftconsts.VendorArista,
					SoftwareVersion: "4.34.0FX-wbb-DPE",
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("Failed to create Arista CFM state functional translator: %v", err)
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
			if len(elems) < 9 {
				return nil, fmt.Errorf("path %v is too short", fullPath)
			}
			domainID := elems[4].GetName()
			assocID := elems[6].GetName()
			localMepID := elems[8].GetName()
			outgoingUpdate := &gnmipb.Update{
				Path: returnPath(domainID, assocID, localMepID),
				Val:  update.GetVal(),
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
		if len(elems) > 0 && elems[0].GetName() == "Sysdb" {
			if len(elems) < 9 {
				return nil, fmt.Errorf("delete path %v is too short", fullPath)
			}
		}
		for _, pattern := range deletePathPatterns {
			if ftutilities.MatchPath(fullPath, pattern) {
				domainID := elems[4].GetName()
				assocID := elems[6].GetName()
				localMepID := elems[8].GetName()
				deletes = append(deletes, returnPath(domainID, assocID, localMepID))
				break
			}
		}
	}
	return deletes, nil
}

// returnPath returns a gNMI path for the update.
// Does not set the origin or the target.
func returnPath(domainID, assocID, localMepID string) *gnmipb.Path {
	return &gnmipb.Path{
		Elem: []*gnmipb.PathElem{
			{Name: "oam"},
			{Name: "cfm"},
			{Name: "domains"},
			{
				Name: "maintenance-domain",
				Key: map[string]string{
					"domain-name": domainID,
				},
			},
			{Name: "maintenance-associations"},
			{
				Name: "maintenance-association",
				Key: map[string]string{
					"association-name": assocID,
				},
			},
			{Name: "mep-endpoints"},
			{
				Name: "mep-endpoint",
				Key: map[string]string{
					"mep-id": localMepID,
				},
			},
			{Name: "state"},
			{Name: "present-rdi"},
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
