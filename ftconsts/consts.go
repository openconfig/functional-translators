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

// Package ftconsts contains the constants used in the functional translators.
package ftconsts

const (
	// VendorArista is the vendor string used for Arista WBB devices.
	VendorArista = "ARISTA"
	// VendorCiscoXR is the vendor string used for Cisco WBB devices.
	VendorCiscoXR = "CISCOXR"
	// VendorJuniper is the vendor string used for Juniper WBB devices.
	VendorJuniper = "JUNIPER"
	// VendorNokia is the vendor string used for Nokia WBB devices.
	VendorNokia = "NOKIA"
)

// Functional Translator Keys.
const (
	// IdentityFunctionalTranslator is the name of the identity functional translator.
	IdentityFunctionalTranslator = "identity-ft"

	// AristaInterfaceDescriptionFunctionalTranslator is the name of the Arista BGP neighbor enabled functional translator.
	AristaBGPNeighborEnabledFunctionalTranslator = "arista-bgp-neighbor-enabled-ft"

	// AristaCfmStateFunctionalTranslator is the name of the Arista CFM state functional translator.
	AristaCfmStateFunctionalTranslator = "arista-cfm-state-ft"

	// AristaInterfaceDescriptionFunctionalTranslator is the name of the Arista interface description functional translator.
	AristaInterfaceDescriptionFunctionalTranslator = "arista-interface-description-ft"

	// AristaDecimalToDoubleFunctionalTranslator is the name of the Arista decimal to double functional translator.
	AristaDecimalToDoubleFunctionalTranslator = "arista-decimal-to-double-ft"

	// AristaInterfaceMacFunctionalTranslator is the name of the Arista interface mac address functional translator.
	AristaInterfaceMacFunctionalTranslator = "arista-interface-mac-ft"

	// AristaMacsecStateFunctionalTranslator is the name of the Arista MACSec ckn and status functional translator.
	AristaMacsecStateFunctionalTranslator = "arista-macsec-state-ft"

	// AristaMacsecCountersTranslator is the name of the Arista macsec counters functional translator.
	AristaMacsecCountersTranslator = "arista-macsec-counters-ft"

	// AristaTransceiverPowerFunctionalTranslator is the name of the Arista transceiver input power functional translator.
	AristaTransceiverPowerFunctionalTranslator = "arista-transceiver-input-power-ft"

	// CiscoXR8000IntegratedCircuitResourceFunctionalTranslator is the name of the identity functional translator.
	CiscoXR8000IntegratedCircuitResourceFunctionalTranslator = "ciscoxr-8000-integrated-circuit-resource-ft"

	// CiscoXRArpTranslator is the name of a translator that provides arp information.
	CiscoXRArpTranslator = "ciscoxr-arp-ft"

	// CiscoXRCarrierTranslator is the name of a translator that provides phy-carrier-transitions information.
	CiscoXRCarrierTranslator = "ciscoxr-carrier-ft"

	// CiscoXRFabricTranslator is the name of a translator that provides fabric information.
	CiscoXRFabricTranslator = "ciscoxr-fabric-ft"

	// CiscoXRFragmentTranslator is the name of a translator that provides fragment packet drops and packets accepted.
	CiscoXRFragmentTranslator = "ciscoxr-fragment-ft"

	// CiscoXRFpdTranslator is the name of a translator that provides fpd status translations.
	CiscoXRFpdTranslator = "ciscoxr-fpd-ft"

	// CiscoXRLagMacFunctionalTranslator is the name of a translator that provides lag mac address translations.
	CiscoXRLagMacFunctionalTranslator = "ciscoxr-lagmac-ft"

	// CiscoXRLaserTranslator is the name of a translator that provides laser information.
	CiscoXRLaserTranslator = "ciscoxr-laser-ft"

	// CiscoXRMountTranslator is the name of a translator that provides mount information.
	CiscoXRMountTranslator = "ciscoxr-mount-ft"

	// CiscoXRQosTranslator is the name of a translator that provides QOS information.
	CiscoXRQosTranslator = "ciscoxr-qos-ft"

	// CiscoXRSubinterfaceCounterTranslator is the name of a translator that provides subinterface
	// counter information, as well as IPv4 address information.
	CiscoXRSubinterfaceCounterTranslator = "ciscoxr-subinterface-counter-ft"

	// CiscoXRTransceiverTranslator is the name of a translator that provides transceiver information.
	CiscoXRTransceiverTranslator = "ciscoxr-transceiver-ft"

	// CiscoXRVendorTranslator is the name of a translator that provides Vendor information.
	CiscoXRVendorDropsTranslator = "ciscoxr-vendordrops-ft"
)
