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

// Package ciscoxrfabric translates fabric native path to openconfig.
package ciscoxrfabric

import (
	"fmt"

	log "github.com/golang/glog"
	"github.com/openconfig/ygot/ytypes"
	xr2431 "github.com/openconfig/functional-translators/ciscoxr/ciscoxrfabric/yang/native"
	fc "github.com/openconfig/functional-translators/ciscoxr/ciscoxrfabric/yang/openconfig"
	"github.com/openconfig/functional-translators/ftconsts"
	"github.com/openconfig/functional-translators/ftutilities"
	"github.com/openconfig/functional-translators/translator"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

type errors struct {
	rxCRCError  uint64
	rxError     uint64
	txFIFOUnrun uint64
}

var (
	translateMap = map[string][]string{
		"/openconfig/components/component/integrated-circuit/pipeline-counters/errors": {
			"/Cisco-IOS-XR-fabric-plane-health-oper/fabric/fabric-plane-ids/fabric-plane-id/fabric-plane-stats/asic-internal-drops",
			"/Cisco-IOS-XR-fabric-plane-health-oper/fabric/fabric-plane-ids/fabric-plane-id/fabric-plane-stats/mcast-lost-cells",
			"/Cisco-IOS-XR-fabric-plane-health-oper/fabric/fabric-plane-ids/fabric-plane-id/fabric-plane-stats/rx-pe-cells",
			"/Cisco-IOS-XR-fabric-plane-health-oper/fabric/fabric-plane-ids/fabric-plane-id/fabric-plane-stats/rx-uce-cells",
			"/Cisco-IOS-XR-fabric-plane-health-oper/fabric/fabric-plane-ids/fabric-plane-id/fabric-plane-stats/ucast-lost-cells",
			"/Cisco-IOS-XR-switch-oper/show-switch/statistics/statistics-detail/statistics-detail-instances/statistics-detail-instance/statistics-detail-port-numbers/statistics-detail-port-number",
		},
	}
	paths                   = ftutilities.MustStringMapPaths(translateMap)
	nativeNonCompliantPaths = []*gnmipb.Path{
		{
			Origin: "Cisco-IOS-XR-switch-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "show-switch"}, {Name: "statistics"}, {Name: "statistics-detail"}, {Name: "statistics-detail-instances"},
				{Name: "statistics-detail-instance"}, {Name: "statistics-detail-port-numbers"}, {Name: "statistics-detail-port-number"},
				{Name: "ethsw-detailed-stat-info"}, {Name: "rx-bad-crc"},
			},
		},
		{
			Origin: "Cisco-IOS-XR-switch-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "show-switch"}, {Name: "statistics"}, {Name: "statistics-detail"}, {Name: "statistics-detail-instances"},
				{Name: "statistics-detail-instance"}, {Name: "statistics-detail-port-numbers"}, {Name: "statistics-detail-port-number"},
				{Name: "ethsw-detailed-stat-info"}, {Name: "rx-errors"},
			},
		},
		{
			Origin: "Cisco-IOS-XR-switch-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "show-switch"}, {Name: "statistics"}, {Name: "statistics-detail"}, {Name: "statistics-detail-instances"},
				{Name: "statistics-detail-instance"}, {Name: "statistics-detail-port-numbers"}, {Name: "statistics-detail-port-number"},
				{Name: "ethsw-detailed-stat-info"}, {Name: "tx-fifo-unrun"},
			},
		},
	}
	// schema is a package-level variable to optimize CiscoXR YANG schema
	// initialization.
	schema    *ytypes.Schema
	schemaErr error
)

// checkCompliantPath returns true if path is yang compliant and can be unmarshalled, false if otherwise and it would raw manipulations.
func checkCompliantPath(prefix *gnmipb.Path, leaves []*gnmipb.Update) bool {
	for _, leaf := range leaves {
		path := ftutilities.Join(prefix, leaf.GetPath())
		if ftutilities.PathInList(path, nativeNonCompliantPaths) {
			return false
		}
	}
	return true
}

// builds port errors map
func buildPorts(prefix *gnmipb.Path, leaves []*gnmipb.Update) map[string]*errors {
	portsErrorsMap := make(map[string]*errors)
	for _, leaf := range leaves {
		path := ftutilities.Join(prefix, leaf.GetPath())
		elems := path.GetElem()
		nodeName := elems[4].GetKey()["node-id"]
		portID := elems[6].GetKey()["port"]
		fullPortName := fmt.Sprintf("%s:%s", nodeName, portID)
		t, ok := portsErrorsMap[fullPortName]
		if !ok {
			t = &errors{
				rxCRCError:  uint64(0),
				rxError:     uint64(0),
				txFIFOUnrun: uint64(0),
			}
			portsErrorsMap[fullPortName] = t
		}
		switch elems[8].GetName() {
		case "rx-bad-crc":
			t.rxCRCError = leaf.GetVal().GetUintVal()
		case "rx-errors":
			t.rxError = leaf.GetVal().GetUintVal()
		case "tx-fifo-unrun":
			t.txFIFOUnrun = leaf.GetVal().GetUintVal()
		}
	}
	return portsErrorsMap
}

// New creates a functional translator.
func New() *translator.FunctionalTranslator {
	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               ftconsts.CiscoXRFabricTranslator,
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
		log.Fatalf("Failed to create Cisco fabric functional translator: %v", err)
	}
	schema, schemaErr = xr2431.Schema()
	if schemaErr != nil {
		log.Fatalf("Failed to get schema: %v", schemaErr)
	}
	return ft
}

func translate(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
	if sr.GetUpdate() == nil {
		return nil, nil
	}
	fcRoot := &fc.Device{}
	n := sr.GetUpdate()
	if !checkCompliantPath(sr.GetUpdate().GetPrefix(), sr.GetUpdate().GetUpdate()) {
		portsErrorsMap := buildPorts(sr.GetUpdate().GetPrefix(), sr.GetUpdate().GetUpdate())
		for port, portErrors := range portsErrorsMap {
			fabricBlock := fcRoot.GetOrCreateComponents().GetOrCreateComponent(port).GetOrCreateIntegratedCircuit().GetOrCreatePipelineCounters().GetOrCreateErrors().GetOrCreateFabricBlock()
			fabricBlock.GetOrCreateFabricBlockError("rx-bad-crc").GetOrCreateState().Count = &portErrors.rxCRCError
			fabricBlock.GetOrCreateFabricBlockError("rx-errors").GetOrCreateState().Count = &portErrors.rxError
			fabricBlock.GetOrCreateFabricBlockError("tx-fifo-unrun").GetOrCreateState().Count = &portErrors.txFIFOUnrun
		}
	} else {
		// Make a shallow copy of the schema and replace the root. This prevents state from one
		// unmarshal operation from leaking into subsequent operations.
		schemaCopy := *schema
		d := &xr2431.CiscoDevice{}
		schemaCopy.Root = d
		if err := ytypes.UnmarshalNotifications(&schemaCopy, []*gnmipb.Notification{n}, nil); err != nil {
			return nil, fmt.Errorf("failed to unmarshal notifications: %v", err)
		}
		fabricPlaneIDMap := d.GetFabric().GetFabricPlaneIds().GetOrCreateFabricPlaneIdMap()
		for fabricPlaneID, fabricPlane := range fabricPlaneIDMap {
			fabricPlaneStats := fabricPlane.GetFabricPlaneStats()
			componentName := fmt.Sprintf("%d", fabricPlaneID)
			fabricBlock := fcRoot.GetOrCreateComponents().GetOrCreateComponent(componentName).GetOrCreateIntegratedCircuit().GetOrCreatePipelineCounters().GetOrCreateErrors().GetOrCreateFabricBlock()
			rxUceCellsError := fabricBlock.GetOrCreateFabricBlockError("uncorrectable-error-cells").GetOrCreateState()
			rxUceCellsError.Count = fabricPlaneStats.RxUceCells
			ucastLostCellsError := fabricBlock.GetOrCreateFabricBlockError("unicast-lost-cells").GetOrCreateState()
			ucastErrors := uint64(*fabricPlaneStats.UcastLostCells)
			ucastLostCellsError.Count = &ucastErrors
			mcastLostCellsError := fabricBlock.GetOrCreateFabricBlockError("multicast-lost-cells").GetOrCreateState()
			mcastErrors := uint64(*fabricPlaneStats.McastLostCells)
			mcastLostCellsError.Count = &mcastErrors
			rxPeCellsError := fabricBlock.GetOrCreateFabricBlockError("parity-error-cells").GetOrCreateState()
			rxPeCellsError.Count = fabricPlaneStats.RxPeCells
			aiDropsError := fabricBlock.GetOrCreateFabricBlockError("asic-internal-drops").GetOrCreateState()
			aiDrops := uint64(*fabricPlaneStats.AsicInternalDrops)
			aiDropsError.Count = &aiDrops
		}
	}
	return ftutilities.FilterStructToState(fcRoot, n.GetTimestamp(), "openconfig", n.GetPrefix().GetTarget())
}
