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

// Package ciscoxrfragment translated cisco native path for fragment traps to openconfig path
package ciscoxrfragment

import (
	"fmt"
	"strconv"

	log "github.com/golang/glog"
	ocfrag "github.com/openconfig/functional-translators/ciscoxr/ciscoxrfragment/yang/openconfig"
	"github.com/openconfig/functional-translators/ftconsts"
	"github.com/openconfig/functional-translators/ftutilities"
	"github.com/openconfig/functional-translators/translator"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

// trap is struct to hold trap info from cisco native path ofa/stats/nodes/node/Cisco-IOS-XR-ofa-npu-stats-oper:npu-numbers/npu-number/display/trap-ids/trap-id.
// We are mainly interested in mtu exceed traps, we collect npu-id, node-name, trap-id to uniquely identify the trap, and to build oc path
type trap struct {
	trapString     string
	nodeName       string
	trapID         uint64
	npuID          uint64
	packetAccepted uint64
	packetDropped  uint64
}

var (
	translateMap = map[string][]string{
		"/openconfig/components/component/integrated-circuit/pipeline-counters/drop/host-interface-block/state/fragment-punt": {
			"/Cisco-IOS-XR-platforms-ofa-oper/ofa/stats/nodes/node/Cisco-IOS-XR-ofa-npu-stats-oper:npu-numbers/npu-number/display/trap-ids/trap-id",
		},
		"/openconfig/components/component/integrated-circuit/pipeline-counters/packet/host-interface-block/state/fragment-punt-pkts": {
			"/Cisco-IOS-XR-platforms-ofa-oper/ofa/stats/nodes/node/Cisco-IOS-XR-ofa-npu-stats-oper:npu-numbers/npu-number/display/trap-ids/trap-id",
		},
	}
	paths      = mustStringMapPaths(translateMap)
	nativePath = &gnmipb.Path{
		Origin: "Cisco-IOS-XR-platforms-ofa-oper",
		Elem: []*gnmipb.PathElem{
			{Name: "ofa"}, {Name: "stats"},
			{Name: "nodes"}, {Name: "node"}, {Name: "Cisco-IOS-XR-ofa-npu-stats-oper:npu-numbers"},
			{Name: "npu-number"}, {Name: "display"}, {Name: "trap-ids"}, {Name: "trap-id"},
			{Name: "*"},
		},
	}
)

func mustStringMapPaths(m map[string][]string) map[string][]*gnmipb.Path {
	p, err := ftutilities.StringMapPaths(m)
	if err != nil {
		log.Fatalf("map %#v cannot parse output paths into gNMI Paths", m)
	}
	return p
}

// validates the leaves and build trap structs
func buildTraps(prefix *gnmipb.Path, leafs []*gnmipb.Update) (map[string]*trap, error) {
	traps := make(map[string]*trap)
	for _, leaf := range leafs {
		path := ftutilities.Join(prefix, leaf.GetPath())
		if !ftutilities.MatchPath(path, nativePath) {
			return nil, fmt.Errorf("invalid path: %v", path)
		}
		elems := path.GetElem()
		nodeName := elems[3].GetKey()["node-name"]
		npuID, err := strconv.Atoi(elems[5].GetKey()["npu-id"])
		if err != nil {
			return nil, fmt.Errorf("failed to convert npu-id to int: %v", err)
		}
		trapID, err := strconv.Atoi(elems[8].GetKey()["trap-id"])
		if err != nil {
			return nil, fmt.Errorf("failed to convert trap-id to int: %v", err)
		}
		fragmentKey := fmt.Sprintf("%s-%v-%v", nodeName, npuID, trapID)
		if _, ok := traps[fragmentKey]; !ok {
			traps[fragmentKey] = &trap{
				nodeName: nodeName,
				trapID:   uint64(trapID),
				npuID:    uint64(npuID),
			}
		}
		switch elems[9].GetName() {
		case "packet-dropped":
			traps[fragmentKey].packetDropped = leaf.GetVal().GetUintVal()
		case "packet-accepted":
			traps[fragmentKey].packetAccepted = leaf.GetVal().GetUintVal()
		case "trap-string":
			traps[fragmentKey].trapString = leaf.GetVal().GetStringVal()
		default:
			// We do not need other leaves, so we skip them.
			continue
		}
	}
	return traps, nil
}

// New creates a functional translator.
func New() *translator.FunctionalTranslator {
	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               ftconsts.CiscoXRFragmentTranslator,
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
		log.Fatalf("Failed to create Cisco fragment functional translator: %v", err)
	}
	return ft
}

func translate(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
	traps, err := buildTraps(sr.GetUpdate().GetPrefix(), sr.GetUpdate().GetUpdate())
	if err != nil {
		return nil, fmt.Errorf("failed to validate path: %v", err)
	}
	fragRoot := &ocfrag.Device{}
	for _, trap := range traps {
		if trap.trapString == "MTU_EXCEEDED" {
			componentName := fmt.Sprintf("%s:%d", trap.nodeName, trap.npuID)
			pipelineCounter := fragRoot.GetOrCreateComponents().GetOrCreateComponent(componentName).GetOrCreateIntegratedCircuit().GetOrCreatePipelineCounters()
			fragmentDroppedState := pipelineCounter.GetOrCreateDrop().GetOrCreateHostInterfaceBlock().GetOrCreateState()
			fragmentAcceptedState := pipelineCounter.GetOrCreatePacket().GetOrCreateHostInterfaceBlock().GetOrCreateState()
			fragmentDroppedState.FragmentPunt = &trap.packetDropped
			fragmentAcceptedState.FragmentPuntPkts = &trap.packetAccepted
		}
	}

	return ftutilities.FilterStructToState(fragRoot, sr.GetUpdate().GetTimestamp(), "openconfig", sr.GetUpdate().GetPrefix().GetTarget())
}
