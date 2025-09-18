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

// Package ciscoxrfpd translates fpd native path to openconfig.
package ciscoxrfpd

import (
	"fmt"

	log "github.com/golang/glog"
	fc "github.com/openconfig/functional-translators/ciscoxr/ciscoxrfpd/yang/openconfig"
	"github.com/openconfig/functional-translators/ftconsts"
	"github.com/openconfig/functional-translators/ftutilities"
	"github.com/openconfig/functional-translators/translator"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

type fpdStatus struct {
	status   []string
	name     []string
	location []string
}

var (
	// TODO(team): Find replacement oc path for FPD status.
	translateMap = map[string][]string{
		"/openconfig/components/component/properties/property/state/value": {
			"/Cisco-IOS-XR-show-fpd-loc-ng-oper/show-fpd/hw-module-fpd",
		},
	}
	paths       = ftutilities.MustStringMapPaths(translateMap)
	nativePaths = []*gnmipb.Path{
		{
			Origin: "Cisco-IOS-XR-show-fpd-loc-ng-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "show-fpd"}, {Name: "hw-module-fpd"}, {Name: "fpd-info-detail"}, {Name: "location"},
			},
		},
		{
			Origin: "Cisco-IOS-XR-show-fpd-loc-ng-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "show-fpd"}, {Name: "hw-module-fpd"}, {Name: "fpd-info-detail"}, {Name: "fpd-name"},
			},
		},
		{
			Origin: "Cisco-IOS-XR-show-fpd-loc-ng-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "show-fpd"}, {Name: "hw-module-fpd"}, {Name: "fpd-info-detail"}, {Name: "status"},
			},
		},
	}
)

// allEqual returns true if all the numbers are equal.
func allEqual(nums ...int) bool {
	for _, n := range nums {
		if n != nums[0] {
			return false
		}
	}
	return true
}

// builds fpdStatus map
func buildFPDStatus(prefix *gnmipb.Path, leaves []*gnmipb.Update) *fpdStatus {
	componentFPDStatus := &fpdStatus{}
	for _, leaf := range leaves {
		path := ftutilities.Join(prefix, leaf.GetPath())
		if !ftutilities.PathInList(path, nativePaths) {
			continue
		}
		elems := path.GetElem()
		switch elems[3].GetName() {
		case "status":
			componentFPDStatus.status = append(componentFPDStatus.status, leaf.GetVal().GetStringVal())
		case "fpd-name":
			componentFPDStatus.name = append(componentFPDStatus.name, leaf.GetVal().GetStringVal())
		case "location":
			componentFPDStatus.location = append(componentFPDStatus.location, leaf.GetVal().GetStringVal())
		}
	}
	return componentFPDStatus
}

// New creates a functional translator.
func New() *translator.FunctionalTranslator {
	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               ftconsts.CiscoXRFpdTranslator,
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
		log.Fatalf("Failed to create Cisco FPD functional translator: %v", err)
	}
	return ft
}

func translate(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
	if sr.GetUpdate() == nil {
		return nil, nil
	}
	fcRoot := &fc.Device{}
	n := sr.GetUpdate()
	fpdStatusList := buildFPDStatus(n.GetPrefix(), n.GetUpdate())
	if !allEqual(len(fpdStatusList.status), len(fpdStatusList.name), len(fpdStatusList.location)) {
		return nil, fmt.Errorf("faulty response:fpdStatusList length mismatch: %v, %v, %v", len(fpdStatusList.status), len(fpdStatusList.name), len(fpdStatusList.location))
	}
	for i, location := range fpdStatusList.location {
		componentName := fmt.Sprintf("%s_%s", location, fpdStatusList.name[i])
		fpdProperty := fcRoot.GetOrCreateComponents().GetOrCreateComponent(componentName).GetOrCreateProperties().GetOrCreateProperty("fpd-status")
		fpdProperty.GetOrCreateState().Value = fc.UnionString(fpdStatusList.status[i])
	}
	return ftutilities.FilterStructToState(fcRoot, n.GetTimestamp(), "openconfig", n.GetPrefix().GetTarget())
}
