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

// Package ciscoxrvendordrops translated cisco native path for vendordrops traps to openconfig path
package ciscoxrvendordrops

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/golang/glog"
	"github.com/openconfig/functional-translators/ftconsts"
	"github.com/openconfig/functional-translators/ftutilities"
	"github.com/openconfig/functional-translators/translator"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

// trap is struct to hold trap info from cisco native path ofa/stats/nodes/node/Cisco-IOS-XR-ofa-npu-stats-oper:npu-numbers/npu-number/display/trap-ids/trap-id.
// We are mainly interested in mtu exceed traps, we collect npu-id, node-name, trap-id to uniquely identify the trap, and to build oc path.
type trap struct {
	trapString    string
	nodeName      string
	trapID        uint64
	npuID         uint64
	packetDropped uint64
}

type stats struct {
	nodeName  string
	npuID     uint64
	blockName string
	counters  map[string]uint64
}

var (
	translateMap = map[string][]string{
		"/openconfig/components/component/integrated-circuit/pipeline-counters/drop/vendor": {
			"/Cisco-IOS-XR-platforms-ofa-oper/ofa/stats/nodes/node/Cisco-IOS-XR-ofa-npu-stats-oper:npu-numbers/npu-number/display/trap-ids/trap-id",
			"/Cisco-IOS-XR-platforms-ofa-oper/ofa/stats/nodes/node/Cisco-IOS-XR-ofa-npu-stats-oper:asic-statistics/asic-statistics-for-npu-ids/asic-statistics-for-npu-id",
		},
	}
	paths            = ftutilities.MustStringMapPaths(translateMap)
	nativeTrapsPaths = []*gnmipb.Path{
		{
			Origin: "Cisco-IOS-XR-platforms-ofa-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "ofa"}, {Name: "stats"},
				{Name: "nodes"}, {Name: "node"}, {Name: "Cisco-IOS-XR-ofa-npu-stats-oper:npu-numbers"},
				{Name: "npu-number"}, {Name: "display"}, {Name: "trap-ids"}, {Name: "trap-id"},
				{Name: "*"},
			},
		},
	}
	nativeStatsPaths = []*gnmipb.Path{
		{
			Origin: "Cisco-IOS-XR-platforms-ofa-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "ofa"}, {Name: "stats"},
				{Name: "nodes"}, {Name: "node"}, {Name: "Cisco-IOS-XR-ofa-npu-stats-oper:asic-statistics"},
				{Name: "asic-statistics-for-npu-ids"}, {Name: "asic-statistics-for-npu-id"},
				{Name: "npu-statistics"}, {Name: "block-info"}, {Name: "block-name"},
			},
		},
		{
			Origin: "Cisco-IOS-XR-platforms-ofa-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "ofa"}, {Name: "stats"},
				{Name: "nodes"}, {Name: "node"}, {Name: "Cisco-IOS-XR-ofa-npu-stats-oper:asic-statistics"},
				{Name: "asic-statistics-for-npu-ids"}, {Name: "asic-statistics-for-npu-id"},
				{Name: "npu-statistics"}, {Name: "block-info"}, {Name: "field-info"},
				{Name: "field-name"},
			},
		},
		{
			Origin: "Cisco-IOS-XR-platforms-ofa-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "ofa"}, {Name: "stats"},
				{Name: "nodes"}, {Name: "node"}, {Name: "Cisco-IOS-XR-ofa-npu-stats-oper:asic-statistics"},
				{Name: "asic-statistics-for-npu-ids"}, {Name: "asic-statistics-for-npu-id"},
				{Name: "npu-statistics"}, {Name: "block-info"}, {Name: "field-info"},
				{Name: "field-value"},
			},
		},
	}
	dropMap = map[string]bool{
		"IFGB_RX 0 partial drop":    true,
		"IFGB_RX 1 partial drop":    true,
		"IFGB_RX 2 partial drop":    true,
		"IFGB_RX 3 partial drop":    true,
		"IFGB_RX 4 partial drop":    true,
		"IFGB_RX 5 partial drop":    true,
		"IFGB_RX 6 partial drop":    true,
		"IFGB_RX 7 partial drop":    true,
		"IFGB_RX 8 partial drop":    true,
		"IFGB_RX 9 partial drop":    true,
		"IFGB_RX 10 partial drop":   true,
		"IFGB_RX 11 partial drop":   true,
		"IFGB_RX 0 full drop":       true,
		"IFGB_RX 1 full drop":       true,
		"IFGB_RX 2 full drop":       true,
		"IFGB_RX 3 full drop":       true,
		"IFGB_RX 4 full drop":       true,
		"IFGB_RX 5 full drop":       true,
		"IFGB_RX 6 full drop":       true,
		"IFGB_RX 7 full drop":       true,
		"IFGB_RX 8 full drop":       true,
		"IFGB_RX 9 full drop":       true,
		"IFGB_RX 10 full drop":      true,
		"IFGB_RX 11 full drop":      true,
		"IFGB_RX 0 undersize drop":  true,
		"IFGB_RX 1 undersize drop":  true,
		"IFGB_RX 2 undersize drop":  true,
		"IFGB_RX 3 undersize drop":  true,
		"IFGB_RX 4 undersize drop":  true,
		"IFGB_RX 5 undersize drop":  true,
		"IFGB_RX 6 undersize drop":  true,
		"IFGB_RX 7 undersize drop":  true,
		"IFGB_RX 8 undersize drop":  true,
		"IFGB_RX 9 undersize drop":  true,
		"IFGB_RX 10 undersize drop": true,
		"IFGB_RX 11 undersize drop": true,
		"PDVOQ drop packets":        true,
		"TXCGM drop":                true,
	}
)

// validates the leaves and build trap structs
func buildTraps(prefix *gnmipb.Path, leaves []*gnmipb.Update) (map[string]*trap, error) {
	traps := make(map[string]*trap)
	for _, leaf := range leaves {
		path := ftutilities.Join(prefix, leaf.GetPath())
		if !ftutilities.PathInList(path, nativeTrapsPaths) {
			continue
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
		trapKey := fmt.Sprintf("%s-%v-%v", nodeName, npuID, trapID)
		if _, ok := traps[trapKey]; !ok {
			traps[trapKey] = &trap{
				nodeName: nodeName,
				trapID:   uint64(trapID),
				npuID:    uint64(npuID),
			}
		}
		switch elems[9].GetName() {
		case "packet-dropped":
			traps[trapKey].packetDropped = leaf.GetVal().GetUintVal()
		case "trap-string":
			traps[trapKey].trapString = leaf.GetVal().GetStringVal()
		default:
			// TODO(team): Expose the number of times we end up here as a statistic to help us expose how much data is being wasted / dropped.
			// We do not need other leaves, so we skip them.
			continue
		}
	}
	return traps, nil
}

// validates the leaves and build stats structs
func buildStats(prefix *gnmipb.Path, leaves []*gnmipb.Update) (map[string]*stats, error) {
	var statsKey string
	statsMap := make(map[string]*stats)
	i := 0
	for i < len(leaves) {
		path := ftutilities.Join(prefix, leaves[i].GetPath())
		if !ftutilities.PathInList(path, nativeStatsPaths) {
			i++
			continue
		}
		elems := path.GetElem()
		nodeName := elems[3].GetKey()["node-name"]
		npuID, err := strconv.Atoi(elems[6].GetKey()["npu-id"])
		if err != nil {
			return nil, fmt.Errorf("failed to convert npu-id to int: %v", err)
		}
		if elems[9].GetName() == "block-name" {
			newBlockName := strings.ReplaceAll(leaves[i].GetVal().GetStringVal(), " ", "_")
			statsKey = fmt.Sprintf("%s:%v:%s", nodeName, npuID, newBlockName)
			if _, ok := statsMap[statsKey]; !ok {
				statsMap[statsKey] = &stats{
					nodeName:  nodeName,
					npuID:     uint64(npuID),
					blockName: newBlockName,
					counters:  make(map[string]uint64),
				}
			}
			i++
			continue
		}
		if elems[10].GetName() == "field-name" && i+1 < len(leaves) {
			fieldName := leaves[i].GetVal().GetStringVal()
			fieldValuePath := ftutilities.Join(prefix, leaves[i+1].GetPath())
			if !ftutilities.MatchPath(fieldValuePath, nativeStatsPaths[2]) {
				return nil, fmt.Errorf("field value path did not come directly after field name path in subscribe response")
			}
			fieldValue := leaves[i+1].GetVal().GetUintVal()
			if _, ok := dropMap[fieldName]; ok {
				newFieldName := strings.ReplaceAll(fieldName, " ", "_")
				statsMap[statsKey].counters[newFieldName] = fieldValue
			}
			i += 2
			continue
		}
		return nil, fmt.Errorf("incoming path did not match any of the native input paths")
	}
	return statsMap, nil
}

// New creates a functional translator.
func New() *translator.FunctionalTranslator {
	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               ftconsts.CiscoXRVendorDropsTranslator,
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
		log.Fatalf("Failed to create Cisco vendor drops functional translator: %v", err)
	}
	return ft
}

func translate(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
	if sr.GetUpdate() == nil {
		return nil, nil
	}
	var err error
	var traps map[string]*trap
	var statsMap map[string]*stats
	n := sr.GetUpdate()
	ts := n.GetTimestamp()
	target := sr.GetUpdate().GetPrefix().GetTarget()
	traps, err = buildTraps(n.GetPrefix(), n.GetUpdate())
	if err != nil {
		return nil, fmt.Errorf("failed to validate path: %v", err)
	}
	statsMap, err = buildStats(n.GetPrefix(), n.GetUpdate())
	if err != nil {
		return nil, fmt.Errorf("failed to validate path: %v", err)
	}
	updates := make([]*gnmipb.Update, 0)
	for _, trap := range traps {
		componentName := fmt.Sprintf("%s:%d", trap.nodeName, trap.npuID)
		//  The path is build based on the rules defined in https://github.com/openconfig/public/blob/master/doc/vendor_counter_guide.md
		switch trap.trapString {
		case "L3_ROUTE_LOOKUP_FAILED":
			update := &gnmipb.Update{
				Path: &gnmipb.Path{
					Elem: []*gnmipb.PathElem{
						{Name: "components"},
						{Name: "component", Key: map[string]string{"name": componentName}},
						{Name: "integrated-circuit"},
						{Name: "pipeline-counters"},
						{Name: "drop"},
						{Name: "vendor"},
						{Name: "CiscoXR"},
						{Name: "spitfire"},
						{Name: "packet-processing"},
						{Name: "state"},
						{Name: "L3_ROUTE_LOOKUP_FAILED"},
					},
				},
				Val: &gnmipb.TypedValue{
					Value: &gnmipb.TypedValue_UintVal{
						UintVal: uint64(trap.packetDropped),
					},
				},
			}
			updates = append(updates, update)
		case "L3_NULL_ADJ(D*)":
			update := &gnmipb.Update{
				Path: &gnmipb.Path{
					Elem: []*gnmipb.PathElem{
						{Name: "components"},
						{Name: "component", Key: map[string]string{"name": componentName}},
						{Name: "integrated-circuit"},
						{Name: "pipeline-counters"},
						{Name: "drop"},
						{Name: "vendor"},
						{Name: "CiscoXR"},
						{Name: "spitfire"},
						{Name: "packet-processing"},
						{Name: "state"},
						{Name: "L3_NULL_ADJ(D*)"},
					},
				},
				Val: &gnmipb.TypedValue{
					Value: &gnmipb.TypedValue_UintVal{
						UintVal: uint64(trap.packetDropped),
					},
				},
			}
			updates = append(updates, update)
		case "MPLS_TE_MIDPOINT_LDP_LABELS_MISS(D*)":
			update := &gnmipb.Update{
				Path: &gnmipb.Path{
					Elem: []*gnmipb.PathElem{
						{Name: "components"},
						{Name: "component", Key: map[string]string{"name": componentName}},
						{Name: "integrated-circuit"},
						{Name: "pipeline-counters"},
						{Name: "drop"},
						{Name: "vendor"},
						{Name: "CiscoXR"},
						{Name: "spitfire"},
						{Name: "packet-processing"},
						{Name: "state"},
						{Name: "MPLS_TE_MIDPOINT_LDP_LABELS_MISS(D*)"},
					},
				},
				Val: &gnmipb.TypedValue{
					Value: &gnmipb.TypedValue_UintVal{
						UintVal: uint64(trap.packetDropped),
					},
				},
			}
			updates = append(updates, update)
		}
	}

	for key, stats := range statsMap {
		//  The path is build based on the rules defined in https://github.com/openconfig/public/blob/master/doc/vendor_counter_guide.md
		for counter, value := range stats.counters {
			if !strings.Contains(stats.blockName, "Summary") {
				continue
			}
			update := &gnmipb.Update{
				Path: &gnmipb.Path{
					Elem: []*gnmipb.PathElem{
						{Name: "components"},
						{Name: "component", Key: map[string]string{"name": key}},
						{Name: "integrated-circuit"},
						{Name: "pipeline-counters"},
						{Name: "drop"},
						{Name: "vendor"},
						{Name: "CiscoXR"},
						{Name: "spitfire"},
						{Name: "adverse"},
						{Name: "state"},
						{Name: counter},
					},
				},
				Val: &gnmipb.TypedValue{
					Value: &gnmipb.TypedValue_UintVal{
						UintVal: uint64(value),
					},
				},
			}
			updates = append(updates, update)
		}
	}
	outgoingSR := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: ts,
				Prefix: &gnmipb.Path{
					Target: target,
					Origin: "openconfig",
				},
				Update: updates,
			},
		},
	}
	return outgoingSR, nil
}
