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

// Package ciscoxrcarrier translates carrier native path to openconfig.
package ciscoxrcarrier

import (
	"fmt"

	log "github.com/golang/glog"
	"github.com/openconfig/ygot/ytypes"
	xr2431 "github.com/openconfig/functional-translators/ciscoxr/ciscoxrcarrier/yang/native"
	"github.com/openconfig/functional-translators/ftconsts"
	"github.com/openconfig/functional-translators/ftutilities"
	"github.com/openconfig/functional-translators/translator"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

var (
	translateMap = map[string][]string{
		"/openconfig/interfaces/interface/ethernet/state/counters/phy-carrier-transitions": {
			"/Cisco-IOS-XR-infra-statsd-oper/infra-statistics/interfaces/interface/generic-counters/carrier-transitions",
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
			ID:               ftconsts.CiscoXRCarrierTranslator,
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
		log.Fatalf("Failed to create Cisco carrier functional translator: %v", err)
	}
	schema, schemaErr = xr2431.Schema()
	if schemaErr != nil {
		log.Fatalf("Failed to get schema: %v", schemaErr)
	}
	return ft
}

func translate(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
	if sr.GetUpdate() == nil {
		return nil, nil
	}
	n := sr.GetUpdate()
	ts := n.GetTimestamp()
	target := sr.GetUpdate().GetPrefix().GetTarget()
	// Make a shallow copy of the schema and replace the root. This prevents state from one
	// unmarshal operation from leaking into subsequent operations.
	schemaCopy := *schema
	d := &xr2431.CiscoDevice{}
	schemaCopy.Root = d
	if err := ytypes.UnmarshalNotifications(&schemaCopy, []*gnmipb.Notification{n}, nil); err != nil {
		return nil, fmt.Errorf("failed to unmarshal notifications: %v", err)
	}
	carrierIntfMap := d.GetInfraStatistics().GetInterfaces().GetOrCreateInterfaceMap()
	updates := make([]*gnmipb.Update, 0)
	for name, intf := range carrierIntfMap {
		phyCarrierCounter := *intf.GetGenericCounters().CarrierTransitions
		update := &gnmipb.Update{
			Path: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{Name: "interfaces"},
					{Name: "interface", Key: map[string]string{"name": name}},
					{Name: "ethernet"},
					{Name: "state"},
					{Name: "counters"},
					{Name: "phy-carrier-transitions"},
				},
			},
			Val: &gnmipb.TypedValue{
				Value: &gnmipb.TypedValue_UintVal{
					UintVal: uint64(phyCarrierCounter),
				},
			},
		}
		updates = append(updates, update)
	}
	outgoingSR := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Timestamp: ts,
				Prefix: &gnmipb.Path{
					Target: target,
					Origin: "openconfig",
				},
				Update: updates,
			},
		},
	}
	return outgoingSR, nil
}
