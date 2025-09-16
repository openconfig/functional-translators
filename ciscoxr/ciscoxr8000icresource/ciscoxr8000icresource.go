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

// Package ciscoxr8000icresource provides a functional translator for CiscoXR 8000 series integrated circuit resources.
package ciscoxr8000icresource

import (
	"fmt"
	"strconv"

	log "github.com/golang/glog"
	ic "github.com/openconfig/functional-translators/ciscoxr/ciscoxr8000icresource/yang/openconfig"
	"github.com/openconfig/functional-translators/ftconsts"
	"github.com/openconfig/functional-translators/ftutilities"
	"github.com/openconfig/functional-translators/translator"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

type resource struct {
	name               string
	npuID              uint64
	bankID             uint64
	bankName           string
	nodeName           string
	maxEntries         uint64
	inUseEntries       uint64
	redOORThreshold    uint64
	yellowOORThreshold uint64
}

var (
	translateMap = map[string][]string{
		"/openconfig/components/component/integrated-circuit/utilization/resources/resource/state/max-limit": {
			"/Cisco-IOS-XR-platforms-ofa-oper/ofa/stats/nodes/node/Cisco-IOS-XR-8000-platforms-npu-resources-oper:hw-resources-datas/hw-resources-data/npu-hwr",
		},
		"/openconfig/components/component/integrated-circuit/utilization/resources/resource/state/name": {
			"/Cisco-IOS-XR-platforms-ofa-oper/ofa/stats/nodes/node/Cisco-IOS-XR-8000-platforms-npu-resources-oper:hw-resources-datas/hw-resources-data/npu-hwr",
		},
		"/openconfig/components/component/integrated-circuit/utilization/resources/resource/state/used-threshold-upper-clear": {
			"/Cisco-IOS-XR-platforms-ofa-oper/ofa/stats/nodes/node/Cisco-IOS-XR-8000-platforms-npu-resources-oper:hw-resources-datas/hw-resources-data/npu-hwr",
		},
		"/openconfig/components/component/integrated-circuit/utilization/resources/resource/state/used-threshold-upper": {
			"/Cisco-IOS-XR-platforms-ofa-oper/ofa/stats/nodes/node/Cisco-IOS-XR-8000-platforms-npu-resources-oper:hw-resources-datas/hw-resources-data/npu-hwr",
		},
		"/openconfig/components/component/integrated-circuit/utilization/resources/resource/state/used": {
			"/Cisco-IOS-XR-platforms-ofa-oper/ofa/stats/nodes/node/Cisco-IOS-XR-8000-platforms-npu-resources-oper:hw-resources-datas/hw-resources-data/npu-hwr",
		},
	}
	wantedCounters = map[string]bool{
		"inuse-entries": true,
		"max-entries":   true,
	}
	wantedOORStates = map[string]bool{
		"red-oor-threshold":    true,
		"yellow-oor-threshold": true,
	}
	paths       = ftutilities.MustStringMapPaths(translateMap)
	nativePaths = []*gnmipb.Path{
		{
			Origin: "Cisco-IOS-XR-platforms-ofa-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "ofa"}, {Name: "stats"}, {Name: "nodes"}, {Name: "node"},
				{Name: "Cisco-IOS-XR-8000-platforms-npu-resources-oper:hw-resources-datas"},
				{Name: "hw-resources-data"}, {Name: "npu-hwr"}, {Name: "bank"},
				{Name: "counter"}, {Name: "inuse-entries"},
			},
		},
		{
			Origin: "Cisco-IOS-XR-platforms-ofa-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "ofa"}, {Name: "stats"}, {Name: "nodes"}, {Name: "node"},
				{Name: "Cisco-IOS-XR-8000-platforms-npu-resources-oper:hw-resources-datas"},
				{Name: "hw-resources-data"}, {Name: "npu-hwr"}, {Name: "bank"},
				{Name: "oor-state"}, {Name: "red-oor-threshold"},
			},
		},
		{
			Origin: "Cisco-IOS-XR-platforms-ofa-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "ofa"}, {Name: "stats"}, {Name: "nodes"}, {Name: "node"},
				{Name: "Cisco-IOS-XR-8000-platforms-npu-resources-oper:hw-resources-datas"},
				{Name: "hw-resources-data"}, {Name: "npu-hwr"}, {Name: "bank"},
				{Name: "oor-state"}, {Name: "yellow-oor-threshold"},
			},
		},
		{
			Origin: "Cisco-IOS-XR-platforms-ofa-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "ofa"}, {Name: "stats"}, {Name: "nodes"}, {Name: "node"},
				{Name: "Cisco-IOS-XR-8000-platforms-npu-resources-oper:hw-resources-datas"},
				{Name: "hw-resources-data"}, {Name: "npu-hwr"}, {Name: "bank"},
				{Name: "bank-name"},
			},
		},
		{
			Origin: "Cisco-IOS-XR-platforms-ofa-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "ofa"}, {Name: "stats"}, {Name: "nodes"}, {Name: "node"},
				{Name: "Cisco-IOS-XR-8000-platforms-npu-resources-oper:hw-resources-datas"},
				{Name: "hw-resources-data"}, {Name: "npu-hwr"}, {Name: "bank"},
				{Name: "counter"}, {Name: "max-entries"},
			},
		},
	}
)

// validates the leaves and builds resource structs
func buildResources(prefix *gnmipb.Path, leafs []*gnmipb.Update) map[string]*resource {
	resources := make(map[string]*resource)
	for _, leaf := range leafs {
		matchContinue := true
		path := ftutilities.Join(prefix, leaf.GetPath())
		for _, p := range nativePaths {
			if ftutilities.MatchPath(path, p) {
				matchContinue = false
				break
			}
		}
		if matchContinue {
			continue
		}
		elems := path.GetElem()
		nodeName := elems[3].GetKey()["node-name"]
		resourceName := elems[5].GetKey()["resource"]
		npuID, err := strconv.Atoi(elems[6].GetKey()["npu-id"])
		if err != nil {
			log.Errorf("failed to parse npu-id: %v", err)
			continue
		}
		bankID, err := strconv.Atoi(elems[7].GetKey()["bank-id"])
		if err != nil {
			log.Errorf("failed to parse bank-id: %v", err)
			continue
		}
		bankName := fmt.Sprintf("%s-%d", resourceName, bankID)
		resourceKey := fmt.Sprintf("%s-%s-%d-%d", nodeName, resourceName, npuID, bankID)
		if _, ok := resources[resourceKey]; !ok {
			resources[resourceKey] = &resource{
				name:     resourceName,
				nodeName: nodeName,
				bankName: bankName,
				npuID:    uint64(npuID),
				bankID:   uint64(bankID),
			}
		}
		switch elems[8].GetName() {
		case "counter":
			if elems[9].GetName() == "inuse-entries" {
				resources[resourceKey].inUseEntries = leaf.GetVal().GetUintVal()
			} else if elems[9].GetName() == "max-entries" {
				resources[resourceKey].maxEntries = leaf.GetVal().GetUintVal()
			} else {
				continue
			}
		case "oor-state":
			if elems[9].GetName() == "red-oor-threshold" {
				resources[resourceKey].redOORThreshold = leaf.GetVal().GetUintVal()
			} else if elems[9].GetName() == "yellow-oor-threshold" {
				resources[resourceKey].yellowOORThreshold = leaf.GetVal().GetUintVal()
			} else {
				// We do not need other leaves, so we skip them.
				continue
			}
		default:
			// We do not need other leaves, so we skip them.
			continue
		}
	}
	return resources
}

// New creates a functional translator.
func New() *translator.FunctionalTranslator {
	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               ftconsts.CiscoXR8000IntegratedCircuitResourceFunctionalTranslator,
			Translate:        translate,
			OutputToInputMap: paths,
			Metadata: []*translator.FTMetadata{
				{
					Vendor: ftconsts.VendorCiscoXR,
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("Failed to create Cisco HW resource 7.10.2 functional translator: %v", err)
	}
	return ft
}

func translate(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
	resources := buildResources(sr.GetUpdate().GetPrefix(), sr.GetUpdate().GetUpdate())
	icRoot := &ic.Device{}
	for _, resource := range resources {
		componentName := fmt.Sprintf("%s-NPU-%d", resource.nodeName, resource.npuID)
		icResource := icRoot.GetOrCreateComponents().GetOrCreateComponent(componentName).GetOrCreateIntegratedCircuit().GetOrCreateUtilization().GetOrCreateResources().GetOrCreateResource(resource.bankName).GetOrCreateState()
		icResource.Name = &resource.bankName
		icResource.Used = &resource.inUseEntries
		icResource.MaxLimit = &resource.maxEntries
		if resource.redOORThreshold != 0 && resource.redOORThreshold <= 255 { // uint8 max
			redOOR := uint8(resource.redOORThreshold)
			icResource.UsedThresholdUpper = &redOOR
		} else {
			log.Errorf("Threshold for redOOR %d > 255 or < 0", resource.redOORThreshold)
		}
		if resource.yellowOORThreshold != 0 && resource.yellowOORThreshold <= 255 { // uint8 max
			yellowOOR := uint8(resource.yellowOORThreshold)
			icResource.UsedThresholdUpperClear = &yellowOOR
		} else {
			log.Errorf("Threshold for yellowOOR %d > 255 or < 0", resource.yellowOORThreshold)
		}
	}
	return ftutilities.FilterStructToState(icRoot, sr.GetUpdate().GetTimestamp(), "openconfig", sr.GetUpdate().GetPrefix().GetTarget())
}
