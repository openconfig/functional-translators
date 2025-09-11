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

// Package ciscoxrtransceiver implements the translation of Cisco transceiver paths.
package ciscoxrtransceiver

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
		"/openconfig/components/component/transceiver/physical-channels/channel/state/index": {
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/lane-data/lane-index",
		},
		"/openconfig/components/component/transceiver/physical-channels/channel/state/input-power/instant": {
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/derived-optics-type",
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/lane-data/lane-index",
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/lane-data/receive-power",
		},
		"/openconfig/components/component/transceiver/physical-channels/channel/state/laser-bias-current/instant": {
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/lane-data/laser-bias-current-milli-amps",
		},
		"/openconfig/components/component/transceiver/physical-channels/channel/state/output-power/instant": {
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/derived-optics-type",
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/lane-data/lane-index",
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/lane-data/transmit-power",
		},
		"/openconfig/components/component/transceiver/state/form-factor": {
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/form-factor",
		},
		"/openconfig/components/component/transceiver/state/vendor": {
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/transceiver-info/vendor-name",
		},
		"/openconfig/components/component/transceiver/state/vendor-part": {
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/transceiver-info/optics-vendor-part",
		},
		"/openconfig/components/component/transceiver/state/vendor-rev": {
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/transceiver-info/optics-vendor-rev",
		},
	}
	expectedOpticsPrefix = &gnmipb.Path{
		Origin: "Cisco-IOS-XR-controller-optics-oper",
		Elem: []*gnmipb.PathElem{
			{Name: "optics-oper"},
			{Name: "optics-ports"},
			{Name: "optics-port"},
			{Name: "optics-info"},
		},
	}

	derivedOpticsType         = "derived-optics-type"
	formFactor                = "form-factor"
	laneIndex                 = "lane-index"
	laserBiasCurrentMilliAmps = "laser-bias-current-milli-amps"
	opticsVendorPart          = "optics-vendor-part"
	opticsVendorRev           = "optics-vendor-rev"
	receivePower              = "receive-power"
	transmitPower             = "transmit-power"
	vendorName                = "vendor-name"
	expectedLeaves            = map[string]bool{
		derivedOpticsType:         true,
		formFactor:                true,
		laneIndex:                 true,
		laserBiasCurrentMilliAmps: true,
		opticsVendorPart:          true,
		opticsVendorRev:           true,
		receivePower:              true,
		transmitPower:             true,
		vendorName:                true,
	}
)

func hasPrefix(path *gnmipb.Path, prefix *gnmipb.Path) bool {
	if len(path.GetElem()) < len(prefix.GetElem()) {
		return false
	}
	for i := 0; i < len(prefix.GetElem()); i++ {
		if path.GetElem()[i].GetName() != prefix.GetElem()[i].GetName() {
			return false
		}
	}
	return true
}

func pathExpected(path *gnmipb.Path) bool {
	return hasPrefix(path, expectedOpticsPrefix) && expectedLeaves[path.GetElem()[len(path.GetElem())-1].GetName()]
}

func index(componentName, laneID string) *gnmipb.Path {
	return &gnmipb.Path{
		Elem: []*gnmipb.PathElem{
			{Name: "components"},
			{Name: "component", Key: map[string]string{"name": componentName}},
			{Name: "transceiver"},
			{Name: "physical-channels"},
			{Name: "channel", Key: map[string]string{"index": laneID}},
			{Name: "state"},
			{Name: "index"},
		},
	}
}

func inputPower(componentName, laneID string) *gnmipb.Path {
	// Don't set origin, as it is set in the prefix.
	return &gnmipb.Path{
		Elem: []*gnmipb.PathElem{
			{Name: "components"},
			{Name: "component", Key: map[string]string{"name": componentName}},
			{Name: "transceiver"},
			{Name: "physical-channels"},
			{Name: "channel", Key: map[string]string{"index": laneID}},
			{Name: "state"},
			{Name: "input-power"},
			{Name: "instant"},
		},
	}
}

func laserBiasCurrent(componentName, laneID string) *gnmipb.Path {
	return &gnmipb.Path{
		Elem: []*gnmipb.PathElem{
			{Name: "components"},
			{Name: "component", Key: map[string]string{"name": componentName}},
			{Name: "transceiver"},
			{Name: "physical-channels"},
			{Name: "channel", Key: map[string]string{"index": laneID}},
			{Name: "state"},
			{Name: "laser-bias-current"},
			{Name: "instant"},
		},
	}
}

func outputPower(componentName, laneID string) *gnmipb.Path {
	// Don't set origin, as it is set in the prefix.
	return &gnmipb.Path{
		Elem: []*gnmipb.PathElem{
			{Name: "components"},
			{Name: "component", Key: map[string]string{"name": componentName}},
			{Name: "transceiver"},
			{Name: "physical-channels"},
			{Name: "channel", Key: map[string]string{"index": laneID}},
			{Name: "state"},
			{Name: "output-power"},
			{Name: "instant"},
		},
	}
}

func formFactorPath(componentName string) *gnmipb.Path {
	return &gnmipb.Path{
		Elem: []*gnmipb.PathElem{
			{Name: "components"},
			{Name: "component", Key: map[string]string{"name": componentName}},
			{Name: "transceiver"},
			{Name: "state"},
			{Name: "form-factor"},
		},
	}
}

func vendorPath(componentName string) *gnmipb.Path {
	return &gnmipb.Path{
		Elem: []*gnmipb.PathElem{
			{Name: "components"},
			{Name: "component", Key: map[string]string{"name": componentName}},
			{Name: "transceiver"},
			{Name: "state"},
			{Name: "vendor"},
		},
	}
}

func vendorPartPath(componentName string) *gnmipb.Path {
	return &gnmipb.Path{
		Elem: []*gnmipb.PathElem{
			{Name: "components"},
			{Name: "component", Key: map[string]string{"name": componentName}},
			{Name: "transceiver"},
			{Name: "state"},
			{Name: "vendor-part"},
		},
	}
}

func vendorRevPath(componentName string) *gnmipb.Path {
	return &gnmipb.Path{
		Elem: []*gnmipb.PathElem{
			{Name: "components"},
			{Name: "component", Key: map[string]string{"name": componentName}},
			{Name: "transceiver"},
			{Name: "state"},
			{Name: "vendor-rev"},
		},
	}
}

type update struct {
	fullPath   *gnmipb.Path
	laneIndex  string
	opticsType string
	value      *gnmipb.TypedValue
}

func (u *update) leaf() string {
	return u.fullPath.GetElem()[len(u.fullPath.GetElem())-1].GetName()
}

func (u *update) componentName() (name string, wanted bool) {
	return ftutilities.MaybeConvertOptical(u.fullPath.GetElem()[2].GetKey()["name"], u.opticsType)
}

func (u *update) toOpenConfig() *gnmipb.Update {
	var outgoingPath *gnmipb.Path
	name, wanted := u.componentName()
	if !wanted {
		return nil
	}
	switch u.leaf() {
	case laneIndex:
		outgoingPath = index(name, u.laneIndex)
	case receivePower:
		outgoingPath = inputPower(name, u.laneIndex)
	case laserBiasCurrentMilliAmps:
		outgoingPath = laserBiasCurrent(name, u.laneIndex)
	case transmitPower:
		outgoingPath = outputPower(name, u.laneIndex)
	case formFactor:
		outgoingPath = formFactorPath(name)
	case vendorName:
		outgoingPath = vendorPath(name)
	case opticsVendorRev:
		outgoingPath = vendorRevPath(name)
	case opticsVendorPart:
		outgoingPath = vendorPartPath(name)
	default:
		// This should never happen, as we filter out unexpected paths.
		return nil
	}
	return &gnmipb.Update{
		Path: outgoingPath,
		Val:  u.value,
	}
}

func dbmValue(u *gnmipb.Update) (*gnmipb.TypedValue, error) {
	dbmFactor := 100.0
	switch t := u.GetVal().GetValue().(type) {
	case *gnmipb.TypedValue_IntVal:
		return &gnmipb.TypedValue{Value: &gnmipb.TypedValue_DoubleVal{DoubleVal: float64(t.IntVal) / dbmFactor}}, nil
	default:
		return nil, fmt.Errorf("unexpected value type %T received in update %v", t, u)
	}
}

func isLaneIndex(u *gnmipb.Update) bool {
	return u.GetPath().GetElem()[len(u.GetPath().GetElem())-1].GetName() == laneIndex
}

func isOpticsType(u *gnmipb.Update) bool {
	return u.GetPath().GetElem()[len(u.GetPath().GetElem())-1].GetName() == derivedOpticsType
}

func translate(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
	// Silently ignore deletes and paths we don't care about.
	var outgoingUpdates []*gnmipb.Update
	srPrefix := sr.GetUpdate().GetPrefix()
	var (
		extractedLaneValue  string
		extractedOpticsType string
	)
	for _, u := range sr.GetUpdate().GetUpdate() {
		fullPath := ftutilities.Join(srPrefix, u.GetPath())
		if pathExpected(fullPath) {
			if isOpticsType(u) {
				extractedOpticsType = u.GetVal().GetStringVal()
				continue
			}
			if isLaneIndex(u) {
				ix := fmt.Sprintf("%d", u.GetVal().GetUintVal())
				extractedLaneValue = ix
			}
			var (
				v   *gnmipb.TypedValue
				err error
			)
			switch u.GetPath().GetElem()[len(u.GetPath().GetElem())-1].GetName() {
			case receivePower, transmitPower:
				v, err = dbmValue(u)
				if err != nil {
					log.Errorf("Failed to translate update %v: %v", u, err)
					continue
				}
			default:
				v = u.GetVal()
			}
			up := &update{
				fullPath:   fullPath,
				laneIndex:  extractedLaneValue,
				opticsType: extractedOpticsType,
				value:      v,
			}
			// This assumes that we always get the lane index and optics type before we get the power data.
			// We have to make this assumption because the data is ordered and may contain multiple
			// lane index leaves.
			// We also collect the lane index under the assumption that the optics type always comes
			// before it.
			if oc := up.toOpenConfig(); oc != nil {
				outgoingUpdates = append(outgoingUpdates, oc)
			}
		}
	}
	if len(outgoingUpdates) == 0 {
		return nil, nil
	}
	outgoingSR := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Prefix: &gnmipb.Path{
					Origin: "openconfig",
					Target: srPrefix.GetTarget(),
				},
				Update:    outgoingUpdates,
				Timestamp: sr.GetUpdate().GetTimestamp(),
			},
		},
	}
	return outgoingSR, nil // End translate.
}

// New creates a functional translator.
func New() *translator.FunctionalTranslator {
	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               ftconsts.CiscoXRTransceiverTranslator,
			Translate:        translate,
			OutputToInputMap: ftutilities.MustStringMapPaths(translateMap),
			Metadata: []*translator.FTMetadata{
				{
					Vendor:          ftconsts.VendorCiscoXR,
					SoftwareVersion: "24.3.2",
				},
				{
					Vendor:          ftconsts.VendorCiscoXR,
					SoftwareVersion: "24.3.20",
				},
				{
					Vendor:          ftconsts.VendorCiscoXR,
					SoftwareVersion: "24.3.30.06I-EFT1LabOnly",
				},
				{
					Vendor:          ftconsts.VendorCiscoXR,
					SoftwareVersion: "24.3.30",
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("Failed to create Cisco transceiver functional translator: %v", err)
	}
	return ft
}
