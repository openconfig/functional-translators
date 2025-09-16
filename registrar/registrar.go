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
	"github.com/openconfig/functional-translators/arista/aristacfmstate"
	"github.com/openconfig/functional-translators/arista/aristainterface"
	"github.com/openconfig/functional-translators/ciscoxr/ciscoxr8000icresource"
	"github.com/openconfig/functional-translators/ciscoxr/ciscoxrarp"
	"github.com/openconfig/functional-translators/ciscoxr/ciscoxripv6"
	"github.com/openconfig/functional-translators/ciscoxr/ciscoxrlaser"
	"github.com/openconfig/functional-translators/ciscoxr/ciscoxrmount"
	"github.com/openconfig/functional-translators/ciscoxr/ciscoxrqos"
	"github.com/openconfig/functional-translators/ciscoxr/ciscoxrsubcounters"
	"github.com/openconfig/functional-translators/ciscoxr/ciscoxrtransceiver"
	"github.com/openconfig/functional-translators/ciscoxr/ciscoxrvendordrops"
	"github.com/openconfig/functional-translators/ftconsts"
	"github.com/openconfig/functional-translators/translator"
)

var (
	// FunctionalTranslatorRegistry is an eagerly initialized map with all functional translators. All
	// new functional translator IDs should be added here to be included.
	// TODO: Add the remaining functional translators already listed in ftconsts.go when released.
	FunctionalTranslatorRegistry = map[string]*translator.FunctionalTranslator{
		// go/keep-sorted start
		ftconsts.AristaCfmStateFunctionalTranslator:                       aristacfmstate.New(),
		ftconsts.AristaInterfaceDescriptionFunctionalTranslator:           aristainterface.NewDescFT(),
		ftconsts.AristaInterfaceMacFunctionalTranslator:                   aristainterface.NewMacFT(),
		ftconsts.CiscoXR8000IntegratedCircuitResourceFunctionalTranslator: ciscoxr8000icresource.New(),
		ftconsts.CiscoXRArpTranslator:                                     ciscoxrarp.New(),
		ftconsts.CiscoXRIPv6Translator:                                    ciscoxripv6.New(),
		ftconsts.CiscoXRLaserTranslator:                                   ciscoxrlaser.New(),
		ftconsts.CiscoXRMountTranslator:                                   ciscoxrmount.New(),
		ftconsts.CiscoXRQosTranslator:                                     ciscoxrqos.New(),
		ftconsts.CiscoXRSubinterfaceCounterTranslator:                     ciscoxrsubcounters.New(),
		ftconsts.CiscoXRTransceiverTranslator:                             ciscoxrtransceiver.New(),
		ftconsts.CiscoXRVendorDropsTranslator:                             ciscoxrvendordrops.New(),
		// go/keep-sorted end
	}
)
