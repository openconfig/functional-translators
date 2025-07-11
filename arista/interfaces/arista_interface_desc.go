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

// Package aristainterfacedesc implements a functional translator for interface description.
package aristainterfacedesc

import (
	"strings"

	"github.com/golang/glog"
	"github.com/openconfig/functional-translators/arista/interfaces/yang/openconfig"
	"google3/third_party/openconfig/functional_translators/ftconsts"
	"github.com/openconfig/functional-translators/ftutilities/ftutilities"
	"github.com/openconfig/functional-translators/simplemapper"
	"github.com/openconfig/functional-translators"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

func deleteHandler(n *gnmipb.Notification) ([]*gnmipb.Path, error) {
	prefix := n.GetPrefix()
	deletes := n.GetDelete()

	var returnDeletes []*gnmipb.Path
	for _, delete := range deletes {
		fullPath := ftutilities.Join(prefix, delete)
		pathSchemaString := ftutilities.GNMIPathToSchemaString(fullPath, true)
		descriptionPath := "/openconfig/interfaces/interface/config/description"
		if strings.HasPrefix(descriptionPath, pathSchemaString) {
			returnDeletes = append(returnDeletes, ftutilities.ConfigToState(delete))
		}
	}

	return returnDeletes, nil
}

// New returns a new FunctionalTranslator for Arista interface descriptions.
func New() *translator.FunctionalTranslator {
	m, err := simplemapper.NewSimpleMapper(interfaces.Schema, interfaces.Schema,
		map[string]string{
			"/openconfig/interfaces/interface[name=<interfaceName>]/state/description": "/openconfig/interfaces/interface[name=<interfaceName>]/config/description",
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
			Translate:        m.Handler,
			ID:               ftconsts.AristaInterfaceDescriptionFunctionalTranslator,
			OutputToInputMap: p,
			Metadata: []*translator.DeviceMetadata{
				&translator.DeviceMetadata{
					Vendor: "arista",
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("Failed to create Arista interface description functional translator: %v", err)
	}
	return ft
}
