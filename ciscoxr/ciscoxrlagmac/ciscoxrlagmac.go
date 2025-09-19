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

// Package ciscoxrlagmac implements a functional translator for the Lag MAC address.
package ciscoxrlagmac

import (
	"strings"

	log "github.com/golang/glog"
	oc "github.com/openconfig/functional-translators/ciscoxr/ciscoxrlagmac/yang/openconfig"
	"github.com/openconfig/functional-translators/ftconsts"
	"github.com/openconfig/functional-translators/ftutilities"
	"github.com/openconfig/functional-translators/simplemapper"
	"github.com/openconfig/functional-translators/translator"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

const (
	intfNameIdx   = 2
	lacpMacSchema = "/openconfig/lacp/interfaces/interface/state/system-id-mac"
	intfMacSchema = "/openconfig/interfaces/interface/ethernet/state/mac-address"
)

func intfMacPath(intfName string) *gnmipb.Path {
	return &gnmipb.Path{
		Elem: []*gnmipb.PathElem{
			{Name: "interfaces"},
			{
				Name: "interface",
				Key: map[string]string{
					"name": intfName,
				},
			},
			{Name: "ethernet"},
			{Name: "state"},
			{Name: "mac-address"},
		},
	}
}

func deleteHandler(n *gnmipb.Notification) ([]*gnmipb.Path, error) {
	prefix := n.GetPrefix()
	deletes := n.GetDelete()
	var returnDeletes []*gnmipb.Path
	for _, d := range deletes {
		fullPath := ftutilities.Join(prefix, d)
		gotSchema := ftutilities.GNMIPathToSchemaString(fullPath, true)

		switch {
		case strings.HasPrefix(lacpMacSchema, gotSchema):
			intfName := d.GetElem()[intfNameIdx].GetKey()["name"]
			returnDeletes = append(returnDeletes, intfMacPath(intfName))
		case strings.HasPrefix(intfMacSchema, gotSchema):
			returnDeletes = append(returnDeletes, d)
		default:
			continue
		}
	}
	return returnDeletes, nil
}

// New returns a new FunctionalTranslator for Cisco interface descriptions.
func New() *translator.FunctionalTranslator {
	m, err := simplemapper.NewSimpleMapper(oc.Schema, oc.Schema,
		map[string]string{
			"/openconfig/interfaces/interface[name=<lagIntfName>]/ethernet/state/mac-address":      "/openconfig/lacp/interfaces/interface[name=<lagIntfName>]/state/system-id-mac",
			"/openconfig/interfaces/interface[name=<ethernetIntfName>]/ethernet/state/mac-address": "/openconfig/interfaces/interface[name=<ethernetIntfName>]/ethernet/state/mac-address",
		},
		deleteHandler,
	)
	if err != nil {
		log.Fatalf("Failed to create mapper: %v", err)
	}

	p, err := ftutilities.StringMapPaths(m.OutputToInputSchemaStrings())
	if err != nil {
		log.Fatalf("map %#v cannot parse output paths into gNMI Paths", m.OutputToInputSchemaStrings())
	}

	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               ftconsts.CiscoXRLagMacFunctionalTranslator,
			Translate:        m.Handler,
			OutputToInputMap: p,
			Metadata: []*translator.FTMetadata{
				{
					Vendor: ftconsts.VendorCiscoXR,
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("Failed to create Cisco LAG MAC functional translator: %v", err)
	}
	return ft
}
