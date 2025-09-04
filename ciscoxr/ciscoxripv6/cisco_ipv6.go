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

// Package ciscoxripv6 implements the CiscoXR specific subinterface translation.
package ciscoxripv6

import (
	log "github.com/golang/glog"
	"github.com/openconfig/functional-translators/ftconsts"
	"github.com/openconfig/functional-translators/ftutilities"
	"github.com/openconfig/functional-translators/translator"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

var (
	translateMap = map[string][]string{
		"/openconfig/interfaces/interface/subinterfaces/subinterface/ipv6/addresses/address/state/ip": {
			"/Cisco-IOS-XR-ipv6-ma-oper/ipv6-network/nodes/node/interface-data/vrfs/vrf/briefs/brief/address",
		},
		"/openconfig/interfaces/interface/subinterfaces/subinterface/ipv6/addresses/address/state/prefix-length": {
			"/Cisco-IOS-XR-ipv6-ma-oper/ipv6-network/nodes/node/interface-data/vrfs/vrf/briefs/brief/address",
			"/Cisco-IOS-XR-ipv6-ma-oper/ipv6-network/nodes/node/interface-data/vrfs/vrf/briefs/brief/prefix-length",
		},
	}
	paths = ftutilities.MustStringMapPaths(translateMap)
)

const (
	subinterfaceZero = "0"

	// These paths are stripped of their origin and both have length 10.
	// addressPath : "ipv6-network/nodes/node/interface-data/vrfs/vrf/briefs/brief/address/address"
	// prefixLengthPath : "ipv6-network/nodes/node/interface-data/vrfs/vrf/briefs/brief/address/prefix-length"
	validPathElemCount = 10

	addressPath      = "/Cisco-IOS-XR-ipv6-ma-oper/ipv6-network/nodes/node/interface-data/vrfs/vrf/briefs/brief/address/address"
	prefixLengthPath = "/Cisco-IOS-XR-ipv6-ma-oper/ipv6-network/nodes/node/interface-data/vrfs/vrf/briefs/brief/address/prefix-length"

	// Interface index is given under the `brief` element
	interfaceNameIndex = 7
	interfaceNameKey   = "interface-name"
)

type interfaceData struct {
	interfaceName  string
	address        string
	ipPrefixLength uint64 // uint8 in OC; encoded as uint64 in gNMI.
}

func (i *interfaceData) ip() *gnmipb.Update {
	return &gnmipb.Update{
		Path: &gnmipb.Path{
			Elem: []*gnmipb.PathElem{
				{Name: "interfaces"},
				{
					Name: "interface",
					Key: map[string]string{
						"name": i.interfaceName,
					},
				},
				{Name: "subinterfaces"},
				{
					Name: "subinterface",
					Key: map[string]string{
						"index": subinterfaceZero,
					},
				},
				{Name: "ipv6"},
				{Name: "addresses"},
				{
					Name: "address",
					Key: map[string]string{
						"address": i.address,
					},
				},
				{Name: "state"},
				{Name: "ip"},
			},
		},
		Val: &gnmipb.TypedValue{
			Value: &gnmipb.TypedValue_StringVal{
				StringVal: i.address,
			},
		},
	}
}

// prefixLength returns a gNMI update for the prefix length. No update is returned if the prefix
// length is 0.
func (i *interfaceData) prefixLength() *gnmipb.Update {
	if i.ipPrefixLength == 0 {
		return nil
	}
	return &gnmipb.Update{
		Path: &gnmipb.Path{
			Elem: []*gnmipb.PathElem{
				{Name: "interfaces"},
				{
					Name: "interface",
					Key: map[string]string{
						"name": i.interfaceName,
					},
				},
				{Name: "subinterfaces"},
				{
					Name: "subinterface",
					Key: map[string]string{
						"index": subinterfaceZero,
					},
				},
				{Name: "ipv6"},
				{Name: "addresses"},
				{
					Name: "address",
					Key: map[string]string{
						"address": i.address,
					},
				},
				{Name: "state"},
				{Name: "prefix-length"},
			},
		},
		Val: &gnmipb.TypedValue{
			Value: &gnmipb.TypedValue_UintVal{
				UintVal: i.ipPrefixLength,
			},
		},
	}
}

type pathResult int

const (
	invalid pathResult = iota
	address
	prefixLength
)

func isValid(path *gnmipb.Path) pathResult {
	if len(path.GetElem()) != validPathElemCount {
		return invalid
	}
	pathString := ftutilities.GNMIPathToSchemaString(path, false)
	switch pathString {
	case addressPath:
		return address
	case prefixLengthPath:
		return prefixLength
	}
	return invalid
}

func intfName(n *gnmipb.Notification) string {
	if len(n.GetPrefix().GetElem()) < interfaceNameIndex {
		return ""
	}
	return n.GetPrefix().GetElem()[interfaceNameIndex].GetKey()[interfaceNameKey]
}

// parse parses the notification and returns a list of interfaceData.
// The expectation is that the address must always come before the prefix length.
func parse(n *gnmipb.Notification) []*interfaceData {
	var out []*interfaceData

	prefix := n.GetPrefix()
	var curAddress string
	for _, update := range n.GetUpdate() {
		fullPath := ftutilities.Join(prefix, update.GetPath())
		switch r := isValid(fullPath); r {
		case address:
			curAddress = update.GetVal().GetStringVal()
			out = append(out, &interfaceData{
				interfaceName: intfName(n),
				address:       curAddress,
			})
		case prefixLength:
			if len(out) == 0 {
				log.Errorf("During translation of interface %s, prefix length update received before address update", intfName(n))
				continue
			}
			out[len(out)-1].ipPrefixLength = update.GetVal().GetUintVal()
		}
	}
	return out
}

// New creates a functional translator.
func New() *translator.FunctionalTranslator {
	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               ftconsts.CiscoXRIPv6Translator,
			Translate:        translate,
			OutputToInputMap: paths,
			Metadata: []*translator.DeviceMetadata{
				{
					Vendor: ftconsts.VendorCiscoXR,
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("Failed to create Cisco subinterface functional translator: %v", err)
	}
	return ft
}

func translate(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
	if sr.GetUpdate() == nil {
		return nil, nil
	}

	outgoingNotification := &gnmipb.Notification{
		Prefix: &gnmipb.Path{
			Origin: "openconfig",
			Target: sr.GetUpdate().GetPrefix().GetTarget(),
		},
		Timestamp: sr.GetUpdate().GetTimestamp(),
	}

	ids := parse(sr.GetUpdate())
	if len(ids) == 0 {
		return nil, nil
	}
	for _, i := range ids {
		outgoingNotification.Update = append(outgoingNotification.Update, i.ip())
		if u := i.prefixLength(); u != nil {
			outgoingNotification.Update = append(outgoingNotification.Update, u)
		}
	}
	return &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: outgoingNotification,
		},
	}, nil
}
