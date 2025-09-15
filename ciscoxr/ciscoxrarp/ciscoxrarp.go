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

// Package ciscoxrarp translates arp native path to openconfig.
package ciscoxrarp

import (
	"fmt"

	log "github.com/golang/glog"
	"github.com/openconfig/ygot/ytypes"
	xr2431 "github.com/openconfig/functional-translators/ciscoxr/ciscoxrarp/yang/native"
	lc "github.com/openconfig/functional-translators/ciscoxr/ciscoxrarp/yang/openconfig"
	"github.com/openconfig/functional-translators/ftconsts"
	"github.com/openconfig/functional-translators/ftutilities"
	"github.com/openconfig/functional-translators/translator"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

var (
	translateMap = map[string][]string{
		"/openconfig/interfaces/interface/subinterfaces/subinterface/ipv4/neighbors/neighbor/state/ip": []string{
			"/Cisco-IOS-XR-ipv4-arp-oper/arp/nodes/node/entries/entry/address",
		},
		"/openconfig/interfaces/interface/subinterfaces/subinterface/ipv4/neighbors/neighbor/state/link-layer-address": {
			"/Cisco-IOS-XR-ipv4-arp-oper/arp/nodes/node/entries/entry/hardware-address",
		},
		"/openconfig/interfaces/interface/subinterfaces/subinterface/ipv6/neighbors/neighbor/state/ip": []string{
			"/Cisco-IOS-XR-ipv6-nd-oper/ipv6-node-discovery/nodes/node/neighbor-interfaces/neighbor-interface/host-addresses/host-address/link-layer-address",
		},
		"/openconfig/interfaces/interface/subinterfaces/subinterface/ipv6/neighbors/neighbor/state/link-layer-address": {
			"Cisco-IOS-XR-ipv6-nd-oper/ipv6-node-discovery/nodes/node/neighbor-interfaces/neighbor-interface/host-addresses/host-address/link-layer-address",
		},
	}
	paths = ftutilities.MustStringMapPaths(translateMap)
)

// New creates a functional translator.
func New() *translator.FunctionalTranslator {
	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               ftconsts.CiscoXRArpTranslator,
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
		log.Fatalf("Failed to create Cisco ARP functional translator: %v", err)
	}
	return ft
}

func translate(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
	if sr.GetUpdate() == nil {
		return nil, nil
	}
	schema, err := xr2431.Schema()
	if err != nil {
		return nil, fmt.Errorf("failed to get schema: %v", err)
	}
	n := sr.GetUpdate()

	if err := ytypes.UnmarshalNotifications(schema, []*gnmipb.Notification{n}, nil); err != nil {
		return nil, fmt.Errorf("failed to unmarshal notifications: %v", err)
	}
	lcRoot := &lc.Device{}
	if schema.Root.(*xr2431.CiscoDevice).GetArp() != nil {
		nodeMap := schema.Root.(*xr2431.CiscoDevice).GetArp().GetNodes().GetOrCreateNodeMap()
		for _, node := range nodeMap {
			entryMap := node.GetOrCreateEntries().GetOrCreateEntryMap()
			for _, entry := range entryMap {
				subState := lcRoot.GetOrCreateInterfaces().GetOrCreateInterface(*entry.InterfaceName).
					GetOrCreateSubinterfaces().
					GetOrCreateSubinterface(0).
					GetOrCreateIpv4().
					GetOrCreateNeighbors().GetOrCreateNeighbor(*entry.Address).
					GetOrCreateState()
				subState.LinkLayerAddress = entry.HardwareAddress
				subState.Ip = entry.Address
			}
		}
	}
	if schema.Root.(*xr2431.CiscoDevice).GetIpv6NodeDiscovery() != nil {
		ndNodeMap := schema.Root.(*xr2431.CiscoDevice).GetIpv6NodeDiscovery().GetNodes().GetOrCreateNodeMap()
		for _, node := range ndNodeMap {
			neighborMap := node.GetNeighborInterfaces().GetOrCreateNeighborInterfaceMap()
			for _, neighbor := range neighborMap {
				addressMap := neighbor.GetHostAddresses().GetOrCreateHostAddressMap()
				for _, address := range addressMap {
					subState := lcRoot.GetOrCreateInterfaces().GetOrCreateInterface(*neighbor.InterfaceName).
						GetOrCreateSubinterfaces().
						GetOrCreateSubinterface(0).
						GetOrCreateIpv6().
						GetOrCreateNeighbors().
						GetOrCreateNeighbor(*address.HostAddress).
						GetOrCreateState()
					subState.LinkLayerAddress = address.LinkLayerAddress
					subState.Ip = address.HostAddress
				}
			}
		}
	}

	return ftutilities.FilterStructToState(lcRoot, n.GetTimestamp(), "openconfig", n.GetPrefix().GetTarget())
}
