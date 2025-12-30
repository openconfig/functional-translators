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

// Package ciscoxrsubcounters implements the CiscoXR specific subinterface counters translation,
// as well as ipv4 address and prefix length translation.
package ciscoxrsubcounters

import (
	"fmt"
	"math"
	"strings"

	log "github.com/golang/glog"
	"github.com/openconfig/ygot/ytypes"
	"github.com/openconfig/functional-translators/ftconsts"
	"github.com/openconfig/functional-translators/ftutilities"
	"github.com/openconfig/functional-translators/translator"

	xr "github.com/openconfig/functional-translators/ciscoxr/ciscoxrsubcounters/yang/native"
	oc "github.com/openconfig/functional-translators/ciscoxr/ciscoxrsubcounters/yang/openconfig"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

const (
	ipv4ProtocolName = "ipv4_unicast"
	ipv6ProtocolName = "ipv6_unicast"
)

var (
	translateMap = map[string][]string{
		// TODO: b/441745512 - Add easy check to alert that the format of these paths is correct for better error messages.
		"/openconfig/interfaces/interface/subinterfaces/subinterface/ipv4/addresses/address/state/ip": {
			"/Cisco-IOS-XR-ipv4-io-oper/ipv4-network/nodes/node/interface-data/vrfs/vrf/details/detail/primary-address",
			"/Cisco-IOS-XR-ipv4-io-oper/ipv4-network/nodes/node/interface-data/vrfs/vrf/details/detail/secondary-address",
		},
		"/openconfig/interfaces/interface/subinterfaces/subinterface/ipv4/addresses/address/state/prefix-length": {
			"/Cisco-IOS-XR-ipv4-io-oper/ipv4-network/nodes/node/interface-data/vrfs/vrf/details/detail/prefix-length",
		},
		"/openconfig/interfaces/interface/subinterfaces/subinterface/ipv4/state/counters/in-pkts": {
			"/Cisco-IOS-XR-infra-statsd-oper/infra-statistics/interfaces/interface/protocols/protocol/packets-received",
		},
		"/openconfig/interfaces/interface/subinterfaces/subinterface/ipv4/state/counters/out-pkts": {
			"/Cisco-IOS-XR-infra-statsd-oper/infra-statistics/interfaces/interface/protocols/protocol/packets-sent",
		},
		"/openconfig/interfaces/interface/subinterfaces/subinterface/ipv6/state/counters/in-pkts": {
			"/Cisco-IOS-XR-infra-statsd-oper/infra-statistics/interfaces/interface/protocols/protocol/packets-received",
		},
		"/openconfig/interfaces/interface/subinterfaces/subinterface/ipv6/state/counters/out-pkts": {
			"/Cisco-IOS-XR-infra-statsd-oper/infra-statistics/interfaces/interface/protocols/protocol/packets-sent",
		},
	}
	paths = ftutilities.MustStringMapPaths(translateMap)
	// schema is a package-level variable to optimize CiscoXR YANG schema
	// initialization.
	schema    *ytypes.Schema
	schemaErr error
)

// New creates a functional translator.
func New() *translator.FunctionalTranslator {
	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               ftconsts.CiscoXRSubinterfaceCounterTranslator,
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
		log.Fatalf("Failed to create Cisco subinterface functional translator: %v", err)
	}
	schema, schemaErr = xr.Schema()
	if schemaErr != nil {
		log.Fatalf("Failed to get schema: %v", schemaErr)
	}
	return ft
}

func translate(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
	if sr.GetUpdate() == nil {
		return nil, nil
	}
	// Make a shallow copy of the schema and replace the root. This prevents state from one
	// unmarshal operation from leaking into subsequent operations.
	schemaCopy := *schema
	schemaCopy.Root = &xr.CiscoDevice{}
	n := sr.GetUpdate()

	if err := ytypes.UnmarshalNotifications(&schemaCopy, []*gnmipb.Notification{n}, nil); err != nil {
		return nil, fmt.Errorf("failed to unmarshal notifications: %v", err)
	}
	if schemaCopy.Root == nil || schemaCopy.Root.(*xr.CiscoDevice) == nil {
		return nil, nil
	}

	root := &oc.Device{}
	handleIPv4(&schemaCopy, root)
	handleCounters(&schemaCopy, root)

	return ftutilities.FilterStructToState(root, n.GetTimestamp(), "openconfig", n.GetPrefix().GetTarget())
}

func handleIPv4(schema *ytypes.Schema, root *oc.Device) {
	ipv4Map := schema.Root.(*xr.CiscoDevice).GetIpv4Network().GetNodes()
	if ipv4Map == nil || ipv4Map.Node == nil {
		return
	}
	for _, node := range ipv4Map.Node {
		if node == nil || node.InterfaceData == nil || node.InterfaceData.Vrfs == nil || node.InterfaceData.Vrfs.Vrf == nil {
			continue
		}
		for _, vrf := range node.InterfaceData.Vrfs.Vrf {
			if vrf == nil || vrf.Details == nil || vrf.Details.Detail == nil {
				continue
			}
			for _, detail := range vrf.Details.Detail {
				if detail == nil {
					continue
				}
				if detail.InterfaceName == nil {
					continue
				}
				subinterface := root.GetOrCreateInterfaces().
					GetOrCreateInterface(*detail.InterfaceName).
					GetOrCreateSubinterfaces().
					GetOrCreateSubinterface(0)
				if detail.PrimaryAddress != nil {
					address := subinterface.GetOrCreateIpv4().
						GetOrCreateAddresses().
						GetOrCreateAddress(*detail.PrimaryAddress)
					address.Ip = detail.PrimaryAddress

					state := address.GetOrCreateState()
					state.Ip = detail.PrimaryAddress
					if detail.PrefixLength != nil && *detail.PrefixLength <= math.MaxUint8 {
						u8 := uint8(*detail.PrefixLength)
						state.PrefixLength = &u8
					}
				}
			}
		}
	}
}

func handleCounters(schema *ytypes.Schema, root *oc.Device) {
	infraMap := schema.Root.(*xr.CiscoDevice).GetInfraStatistics().GetInterfaces()
	if infraMap == nil || infraMap.Interface == nil {
		return
	}
	for _, iface := range infraMap.Interface {
		if iface == nil || iface.Protocols == nil || iface.Protocols.Protocol == nil || iface.InterfaceName == nil {
			continue
		}
		for _, protocol := range iface.Protocols.Protocol {
			if protocol == nil || protocol.ProtocolName == nil {
				continue
			}
			// Only create the subinterface counters if the protocol is ipv4 or ipv6 unicast.
			switch strings.ToLower(*protocol.ProtocolName) {
			case ipv4ProtocolName:
				counters := root.GetOrCreateInterfaces().
					GetOrCreateInterface(*iface.InterfaceName).
					GetOrCreateSubinterfaces().
					GetOrCreateSubinterface(0).
					GetOrCreateIpv4().GetOrCreateState().GetOrCreateCounters()
				if protocol.PacketsReceived != nil {
					counters.InPkts = protocol.PacketsReceived
				}
				if protocol.PacketsSent != nil {
					counters.OutPkts = protocol.PacketsSent
				}
			case ipv6ProtocolName:
				counters := root.GetOrCreateInterfaces().
					GetOrCreateInterface(*iface.InterfaceName).
					GetOrCreateSubinterfaces().
					GetOrCreateSubinterface(0).
					GetOrCreateIpv6().GetOrCreateState().GetOrCreateCounters()
				if protocol.PacketsReceived != nil {
					counters.InPkts = protocol.PacketsReceived
				}
				if protocol.PacketsSent != nil {
					counters.OutPkts = protocol.PacketsSent
				}
			default:
			}
		}
	}
}
