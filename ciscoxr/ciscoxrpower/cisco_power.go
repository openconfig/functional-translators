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

// Package ciscoxrpower translates pem status native path to openconfig component pem oper-status.
package ciscoxrpower

import (
	"fmt"
	"slices"

	log "github.com/golang/glog"
	"github.com/openconfig/functional-translators/ftconsts"
	"github.com/openconfig/functional-translators/ftutilities"
	"github.com/openconfig/functional-translators/translator"

	oc "github.com/openconfig/functional-translators/ciscoxr/ciscoxrpower/yang/openconfig"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

var (
	translateMap = map[string][]string{
		"/openconfig/components/component/state/oper-status": {
			"/openconfig/components/component/state/oper-status",
			"/Cisco-IOS-XR-envmon-oper/power-management/rack/producers/producer-nodes/producer-node/pem-info-array",
		},
	}
	outputToInputMap = ftutilities.MustStringMapPaths(translateMap)
	nativeSrcPaths   = []*gnmipb.Path{
		{
			Origin: "Cisco-IOS-XR-envmon-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "power-management"}, {Name: "rack"}, {Name: "producers"}, {Name: "producer-nodes"},
				{Name: "producer-node"}, {Name: "pem-info-array"}, {Name: "node-name"},
			},
		},
		{
			Origin: "Cisco-IOS-XR-envmon-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "power-management"}, {Name: "rack"}, {Name: "producers"}, {Name: "producer-nodes"},
				{Name: "producer-node"}, {Name: "pem-info-array"}, {Name: "node-status"},
			},
		},
	}
	ocSrcPaths = []*gnmipb.Path{
		{
			Origin: "openconfig",
			Elem: []*gnmipb.PathElem{
				{Name: "components"}, {Name: "component"}, {Name: "state"}, {Name: "oper-status"},
			},
		},
	}
	// These are the PEM names on 8808, the FT will filter the oc oper-status leave for these PEMs, and translate the native path instead.
	pems = []string{"0/PT0-PM0", "0/PT0-PM1", "0/PT0-PM2", "0/PT1-PM0", "0/PT1-PM1", "0/PT1-PM2", "0/PT2-PM0", "0/PT2-PM1", "0/PT2-PM2"}
)

// checkPathInUpdate returns true if one of checked paths exists in the received gnmi update.
func checkPathInUpdate(prefix *gnmipb.Path, leaves []*gnmipb.Update, paths []*gnmipb.Path) bool {
	for _, leaf := range leaves {
		path := ftutilities.Join(prefix, leaf.GetPath())
		if ftutilities.PathInList(path, paths) {
			return true
		}
	}
	return false
}

// buildPowerStatus builds PEM oper-status map (PEM name to PEM oper-status)
func buildPowerStatus(prefix *gnmipb.Path, leaves []*gnmipb.Update) (map[string]string, error) {
	powerStatusMap := make(map[string]string)
	var pemNames, pemStatuses []string
	for _, leaf := range leaves {
		path := ftutilities.Join(prefix, leaf.GetPath())
		switch path.GetElem()[6].GetName() {
		case "node-name":
			pemNames = append(pemNames, leaf.GetVal().GetStringVal())
		case "node-status":
			pemStatuses = append(pemStatuses, leaf.GetVal().GetStringVal())
		}
	}
	if len(pemNames) != len(pemStatuses) {
		return nil, fmt.Errorf("number of power module names and statuses are not equal: %v, %v", len(pemNames), len(pemStatuses))
	}
	for i, name := range pemNames {
		if slices.Contains(pems, name) {
			powerStatusMap[name] = pemStatuses[i]
		}
	}
	return powerStatusMap, nil
}

// deletePowerStatusLeaves deletes PEM oper-status leaves from the gnmi update.
func deletePowerStatusLeaves(prefix *gnmipb.Path, leaves []*gnmipb.Update) []*gnmipb.Update {
	var filteredLeaves []*gnmipb.Update
	for _, leaf := range leaves {
		path := ftutilities.Join(prefix, leaf.GetPath())
		elems := path.GetElem()
		componentName := elems[1].GetKey()["name"]
		if slices.Contains(pems, componentName) {
			continue
		}
		filteredLeaves = append(filteredLeaves, leaf)
	}
	return filteredLeaves
}

// New creates a functional translator.
func New() *translator.FunctionalTranslator {
	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               ftconsts.CiscoXRPowerTranslator,
			Translate:        translate,
			OutputToInputMap: outputToInputMap,
			Metadata: []*translator.FTMetadata{
				{
					Vendor: ftconsts.VendorCiscoXR,
					SoftwareVersionRange: &translator.SWRange{
						InclusiveMin: "24.3.20",
						ExclusiveMax: "25.4.1",
					},
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("Failed to create Cisco Power functional translator: %v", err)
	}
	return ft
}

func translate(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
	if sr.GetUpdate() == nil {
		return nil, nil
	}
	ocRoot := &oc.Device{}
	n := sr.GetUpdate()
	switch {
	case checkPathInUpdate(sr.GetUpdate().GetPrefix(), sr.GetUpdate().GetUpdate(), nativeSrcPaths):
		powerStatusMap, err := buildPowerStatus(sr.GetUpdate().GetPrefix(), sr.GetUpdate().GetUpdate())
		if err != nil {
			return nil, err
		}
		if len(powerStatusMap) == 0 {
			return nil, nil
		}
		for pmName, status := range powerStatusMap {
			switch status {
			case "OK":
				ocRoot.GetOrCreateComponents().GetOrCreateComponent(pmName).GetOrCreateState().OperStatus = oc.OpenconfigPlatformTypes_COMPONENT_OPER_STATUS_ACTIVE
			default:
				ocRoot.GetOrCreateComponents().GetOrCreateComponent(pmName).GetOrCreateState().OperStatus = oc.OpenconfigPlatformTypes_COMPONENT_OPER_STATUS_INACTIVE
			}
		}
		return ftutilities.FilterStructToState(ocRoot, n.GetTimestamp(), "openconfig", n.GetPrefix().GetTarget())
	case checkPathInUpdate(sr.GetUpdate().GetPrefix(), sr.GetUpdate().GetUpdate(), ocSrcPaths):
		filteredLeaves := deletePowerStatusLeaves(sr.GetUpdate().GetPrefix(), sr.GetUpdate().GetUpdate())
		if filteredLeaves == nil {
			return nil, nil
		}
		filteredSR := &gnmipb.SubscribeResponse{
			Response: &gnmipb.SubscribeResponse_Update{
				Update: &gnmipb.Notification{
					Timestamp: n.GetTimestamp(),
					Prefix:    sr.GetUpdate().GetPrefix(),
					Update:    filteredLeaves,
				},
			},
		}
		return filteredSR, nil
	}
	return nil, nil
}
