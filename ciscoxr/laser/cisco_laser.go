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

// Package ciscoxrlaser transslates laser native path to openconfig .
package ciscoxrlaser

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/openconfig/ygot/ytypes"
	"google3/third_party/openconfig/functional_translators/ftconsts"
	"github.com/openconfig/functional-translators/ftutilities/ftutilities"
	"github.com/openconfig/functional-translators"

	xr2431 "google3/third_party/openconfig/functional_translators/ciscoxr/laser/yang/native/gostructs"
	lc "google3/third_party/openconfig/functional_translators/ciscoxr/laser/yang/openconfig/gostructs"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

var (
	translateMap = map[string][]string{
		"/openconfig/components/component/transceiver/thresholds/threshold/state/module-temperature-upper": {
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/temp-high-threshold",
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/temp-high-warning-threshold",
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/derived-optics-type",
		},
		"/openconfig/components/component/transceiver/thresholds/threshold/state/module-temperature-lower": {
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/temp-low-threshold",
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/temp-low-warning-threshold",
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/derived-optics-type",
		},
		"/openconfig/components/component/transceiver/thresholds/threshold/state/input-power-upper": {
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/rx-high-threshold",
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/rx-high-warning-threshold",
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/derived-optics-type",
		},
		"/openconfig/components/component/transceiver/thresholds/threshold/state/input-power-lower": {
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/rx-low-threshold",
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/rx-low-warning-threshold",
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/derived-optics-type",
		},
		"/openconfig/components/component/transceiver/thresholds/threshold/state/output-power-upper": {
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/tx-high-threshold",
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/tx-high-warning-threshold",
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/derived-optics-type",
		},
		"/openconfig/components/component/transceiver/thresholds/threshold/state/output-power-lower": {
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/tx-low-threshold",
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/tx-low-warning-threshold",
			"/Cisco-IOS-XR-controller-optics-oper/optics-oper/optics-ports/optics-port/optics-info/derived-optics-type",
		},
	}
	paths = ftutilities.MustStringMapPaths(translateMap)
)

// New creates a functional translator.
func New() *translator.FunctionalTranslator {
	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               ftconsts.CiscoXRLaserTranslator,
			Translate:        translate,
			OutputToInputMap: paths,
			Metadata: []*translator.DeviceMetadata{
				&translator.DeviceMetadata{
					Vendor: ftconsts.VendorCiscoXR,
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("Failed to create Cisco laser functional translator: %v", err)
	}
	return ft
}

func translate(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
	schema, err := xr2431.Schema()
	if err != nil {
		return nil, fmt.Errorf("failed to get schema: %v", err)
	}
	n := sr.GetUpdate()

	if err := ytypes.UnmarshalNotifications(schema, []*gnmipb.Notification{n}, nil); err != nil {
		return nil, fmt.Errorf("failed to unmarshal notifications: %v", err)
	}

	portMap := schema.Root.(*xr2431.CiscoDevice).GetOpticsOper().GetOpticsPorts().GetOrCreateOpticsPortMap()
	if portMap == nil {
		return nil, fmt.Errorf("failed to get optics port map")
	}
	lcRoot := &lc.Device{}
	for portName, port := range portMap {
		if port == nil {
			continue
		}
		opticsInfo := port.GetOpticsInfo()
		if opticsInfo == nil {
			continue
		}
		modifiedPortName, wanted := ftutilities.MaybeConvertOptical(portName, *opticsInfo.DerivedOpticsType)
		if !wanted {
			continue
		}
		transceiver := lcRoot.GetOrCreateComponents().GetOrCreateComponent(modifiedPortName).GetOrCreateTransceiver()
		criticalThresholdState := transceiver.GetOrCreateThresholds().GetOrCreateThreshold(lc.OpenconfigAlarmTypes_OPENCONFIG_ALARM_SEVERITY_CRITICAL).GetOrCreateState()
		warningThresholdState := transceiver.GetOrCreateThresholds().GetOrCreateThreshold(lc.OpenconfigAlarmTypes_OPENCONFIG_ALARM_SEVERITY_WARNING).GetOrCreateState()
		criticalThresholdState.Severity = lc.OpenconfigAlarmTypes_OPENCONFIG_ALARM_SEVERITY_CRITICAL
		warningThresholdState.Severity = lc.OpenconfigAlarmTypes_OPENCONFIG_ALARM_SEVERITY_WARNING
		if opticsInfo.TempHighThreshold != nil {
			moduleTempAlarmUpper := float64(*opticsInfo.TempHighThreshold) * 0.01 // native unit is 0.01 degree c, openconfig unit is degree c.
			criticalThresholdState.ModuleTemperatureUpper = &moduleTempAlarmUpper
		}
		if opticsInfo.TempLowThreshold != nil {
			moduleTempAlarmLower := float64(*opticsInfo.TempLowThreshold) * 0.01 // native unit is 0.01 degree c, openconfig unit is degree c.
			criticalThresholdState.ModuleTemperatureLower = &moduleTempAlarmLower
		}
		if opticsInfo.RxHighThreshold != nil {
			inputPowerAlarmUpper := float64(*opticsInfo.RxHighThreshold) * 0.1 // native unit is 0.1 dbm, openconfig unit is dbm.
			criticalThresholdState.InputPowerUpper = &inputPowerAlarmUpper
		}
		if opticsInfo.RxLowThreshold != nil {
			inputPowerAlarmLower := float64(*opticsInfo.RxLowThreshold) * 0.1 // native unit is 0.1 dbm, openconfig unit is dbm.
			criticalThresholdState.InputPowerLower = &inputPowerAlarmLower
		}
		if opticsInfo.TxHighThreshold != nil {
			outputPowerAlarmUpper := float64(*opticsInfo.TxHighThreshold) * 0.1 // native unit is 0.1 dbm, openconfig unit is dbm.
			criticalThresholdState.OutputPowerUpper = &outputPowerAlarmUpper
		}
		if opticsInfo.TxLowThreshold != nil {
			outputPowerAlarmLower := float64(*opticsInfo.TxLowThreshold) * 0.1 // native unit is 0.1 dbm, openconfig unit is dbm.
			criticalThresholdState.OutputPowerLower = &outputPowerAlarmLower
		}
		if opticsInfo.TempHighWarningThreshold != nil {
			moduleTempWarningUpper := float64(*opticsInfo.TempHighWarningThreshold) * 0.01 // native unit is 0.01 degree c, openconfig unit is degree c.
			warningThresholdState.ModuleTemperatureUpper = &moduleTempWarningUpper
		}
		if opticsInfo.TempLowWarningThreshold != nil {
			moduleTempWarningLower := float64(*opticsInfo.TempLowWarningThreshold) * 0.01 // native unit is 0.01 degree c, openconfig unit is degree c.
			warningThresholdState.ModuleTemperatureLower = &moduleTempWarningLower
		}
		if opticsInfo.RxHighWarningThreshold != nil {
			inputPowerWarningUpper := float64(*opticsInfo.RxHighWarningThreshold) * 0.1 // native unit is 0.1 dbm, openconfig unit is dbm.
			warningThresholdState.InputPowerUpper = &inputPowerWarningUpper
		}
		if opticsInfo.RxLowWarningThreshold != nil {
			inputPowerWarningLower := float64(*opticsInfo.RxLowWarningThreshold) * 0.1 // native unit is 0.1 dbm, openconfig unit is dbm.
			warningThresholdState.InputPowerLower = &inputPowerWarningLower
		}
		if opticsInfo.TxHighWarningThreshold != nil {
			outputPowerWarningUpper := float64(*opticsInfo.TxHighWarningThreshold) * 0.1 // native unit is 0.1 dbm, openconfig unit is dbm.
			warningThresholdState.OutputPowerUpper = &outputPowerWarningUpper
		}
		if opticsInfo.TxLowWarningThreshold != nil {
			outputPowerWarningLower := float64(*opticsInfo.TxLowWarningThreshold) * 0.1 // native unit is 0.1 dbm, openconfig unit is dbm.
			warningThresholdState.OutputPowerLower = &outputPowerWarningLower
		}
	}

	return ftutilities.FilterStructToState(lcRoot, n.GetTimestamp(), "openconfig", n.GetPrefix().GetTarget())
}
