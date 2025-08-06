// Package registrar provides a map with all functional translators.
package registrar

import (
	"github.com/openconfig/functional-translators/arista/interfaces/aristainterfacedesc"
	"github.com/openconfig/functional-translators/arista/interfaces/aristainterfacemac"
	"github.com/openconfig/functional-translators/ciscoxr/laser/ciscoxrlaser"
	"github.com/openconfig/functional-translators/ftconsts"
	"github.com/openconfig/functional-translators"
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
