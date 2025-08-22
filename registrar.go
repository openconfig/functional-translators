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

// Package registrar provides a map with all functional translators.
package registrar

import (
	"google3/third_party/openconfig/functional_translators/arista/interfaces/aristainterfacedesc"
	"google3/third_party/openconfig/functional_translators/arista/interfaces/aristainterfacemac"
	"google3/third_party/openconfig/functional_translators/ciscoxr/laser/ciscoxrlaser"
	"github.com/openconfig/functional-translators/ftconsts"
	"google3/third_party/openconfig/functional_translators/translator"
)

var (
	// FunctionalTranslatorRegistry is an eagerly initialized map with all functional translators. All
	// new functional translator IDs should be added here to be included.
	// TODO: Add the remaining functional translators already listed in ftconsts.go when released.
	FunctionalTranslatorRegistry = map[string]*translator.FunctionalTranslator{
		// go/keep-sorted start
		ftconsts.AristaInterfaceDescriptionFunctionalTranslator: aristainterfacedesc.New(),
		ftconsts.AristaInterfaceMacFunctionalTranslator:         aristainterfacemac.New(),
		ftconsts.CiscoXRLaserTranslator:                         ciscoxrlaser.New(),
		// go/keep-sorted end
	}
)
