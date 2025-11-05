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

package aristainterface

import (
	"strings"

	log "github.com/golang/glog"
	"github.com/openconfig/functional-translators/arista/aristainterface/yang/openconfig"
	"github.com/openconfig/functional-translators/ftconsts"
	"github.com/openconfig/functional-translators/ftutilities"
	"github.com/openconfig/functional-translators/simplemapper"
	"github.com/openconfig/functional-translators/translator"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

var intfNameIdx = 2
var lacpMacSchema = "/openconfig/lacp/interfaces/interface/state/system-id-mac"
var intfMacSchema = "/openconfig/interfaces/interface/ethernet/state/mac-address"

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

func macDeleteHandler(n *gnmipb.Notification) ([]*gnmipb.Path, error) {
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

// NewMacFT returns a new FunctionalTranslator for Arista interface mac addresses.
func NewMacFT() *translator.FunctionalTranslator {
	m, err := simplemapper.NewSimpleMapper(openconfig.Schema, openconfig.Schema,
		map[string]string{
			"/openconfig/interfaces/interface[name=<lagIntfName>]/ethernet/state/mac-address":      "/openconfig/lacp/interfaces/interface[name=<lagIntfName>]/state/system-id-mac",
			"/openconfig/interfaces/interface[name=<ethernetIntfName>]/ethernet/state/mac-address": "/openconfig/interfaces/interface[name=<ethernetIntfName>]/ethernet/state/mac-address",
		},
		macDeleteHandler,
	)
	if err != nil {
		log.Fatalf("Failed to create mapper: %v", err)
	}

	p := ftutilities.MustStringMapPaths(m.OutputToInputSchemaStrings())

	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			Translate:        m.Handler,
			ID:               ftconsts.AristaInterfaceMacFunctionalTranslator,
			OutputToInputMap: p,
			Metadata: []*translator.FTMetadata{
				{
					Vendor: "arista",
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("Failed to create Arista interface MAC functional translator: %v", err)
	}
	return ft
}
