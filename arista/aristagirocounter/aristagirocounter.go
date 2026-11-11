// Copyright 2026 Google LLC
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

// Package aristagirocounter translates Arista Girocounter PM data from native to openconfig.
// Tracking bug: b/514426299.
// Vendor commitment: Native OpenConfig support slated for EOS 4.37.1F-DPE-CTNR.
package aristagirocounter

import (
	"net"
	"strconv"
	"strings"

	log "github.com/golang/glog"
	"github.com/openconfig/functional-translators/ftconsts"
	"github.com/openconfig/functional-translators/ftutilities"
	"github.com/openconfig/functional-translators/translator"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

const aristaDecapTableID = "4294967295"

var (
	translateMap = map[string][]string{
		"/openconfig/network-instances/network-instance/policy-forwarding/policies/policy/rules/rule/state/matched-pkts": {
			"/eos_native/Smash/flexCounters/counterTable/DecapGroup/4294967295/decapGroupCounter",
		},
		"/openconfig/network-instances/network-instance/policy-forwarding/policies/policy/rules/rule/state/matched-octets": {
			"/eos_native/Smash/flexCounters/counterTable/DecapGroup/4294967295/decapGroupCounter",
		},
		"/openconfig/network-instances/network-instance/policy-forwarding/policies/policy/rules/rule/ipv4/state/destination-address": {
			"/eos_native/Smash/flexCounters/counterTable/DecapGroup/4294967295/decapGroupCounter",
		},
		"/openconfig/network-instances/network-instance/policy-forwarding/policies/policy/rules/rule/ipv6/state/destination-address": {
			"/eos_native/Smash/flexCounters/counterTable/DecapGroup/4294967295/decapGroupCounter",
		},
	}

	ocCommonPrefix = []*gnmipb.PathElem{
		{Name: "network-instances"},
		{
			Name: "network-instance",
			Key: map[string]string{
				"name": "default",
			},
		},
		{Name: "policy-forwarding"},
		{Name: "policies"},
	}

	decapCommonPrefixPattern = &gnmipb.Path{
		Origin: "eos_native",
		Elem: []*gnmipb.PathElem{
			{Name: "Smash"},
			{Name: "flexCounters"},
			{Name: "counterTable"},
			{Name: "DecapGroup"},
			{Name: aristaDecapTableID},
			{Name: "decapGroupCounter"},
			{Name: "*"},
		},
	}
)

type impl struct {
}

// New creates a functional translator.
func New() *translator.FunctionalTranslator {
	i := &impl{}
	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               ftconsts.AristaGirocounterFunctionalTranslator,
			Translate:        i.translate,
			OutputToInputMap: ftutilities.MustStringMapPaths(translateMap),
			Metadata: []*translator.FTMetadata{
				{
					// TODO(team): Set SoftwareVersionRange.ExclusiveMax to the first EOS
					// release that natively supports the policy-forwarding matched-pkts/octets paths.
					Vendor: ftconsts.VendorArista,
					SoftwareVersionRange: &translator.SWRange{
						InclusiveMin: "4.36.1FX-wbb",
						// TODO(team): Confirm 4.37.1 is the first EOS release that natively
						// supports the decap policy-forwarding OC paths.
						ExclusiveMax: "4.37.1",
					},
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("Failed to create Arista Girocounter functional translator: %v", err)
	}
	return ft
}

func getOCPktsPath(key string) *gnmipb.Path {
	return getOCPathForLeaf(key, "matched-pkts")
}

func getOCOctetsPath(key string) *gnmipb.Path {
	return getOCPathForLeaf(key, "matched-octets")
}

func getOCDestAddrPath(key string, isIPv6 bool) *gnmipb.Path {
	ipFamily := "ipv4"
	if isIPv6 {
		ipFamily = "ipv6"
	}
	return &gnmipb.Path{
		Elem: []*gnmipb.PathElem{
			{
				Name: "policy",
				Key: map[string]string{
					"policy-id": key,
				},
			},
			{Name: "rules"},
			{
				Name: "rule",
				Key: map[string]string{
					"sequence-id": getOCSequenceID(key),
				},
			},
			{Name: ipFamily},
			{Name: "state"},
			{Name: "destination-address"},
		},
	}
}

func parseIP(key string) (net.IP, bool) {
	// Try parsing as CIDR first
	ip, _, err := net.ParseCIDR(key)
	if err == nil {
		return ip, ip.To4() == nil
	}

	// Fallback to simple IP parsing
	ip = net.ParseIP(key)
	if ip != nil {
		return ip, ip.To4() == nil
	}

	return nil, false
}

// getOCSequenceID returns a valid OpenConfig uint32 string for the rule sequence-id.
// In OpenConfig policy-forwarding, sequence-id is defined as uint32. Decap endpoints
// are typically IP address strings (e.g. "10.61.61.148/30"), which map to policy-id
// and fall through to sequence-id "1". In synthetic test cases or numeric flexCounter
// groupings where the key is already a numeric uint32 string, that value is preserved.
func getOCSequenceID(key string) string {
	if _, err := strconv.ParseUint(key, 10, 32); err == nil {
		return key
	}
	return "1"
}

func getOCPathForLeaf(key, leafName string) *gnmipb.Path {
	return &gnmipb.Path{
		Elem: []*gnmipb.PathElem{
			// Arista Girocounter flexCounters expose a counter per decap group endpoint.
			// In OpenConfig policy-forwarding, the endpoint maps to policy-id, and
			// sequence-id is a uint32 identifying the rule within that policy.
			{
				Name: "policy",
				Key: map[string]string{
					"policy-id": key,
				},
			},
			{Name: "rules"},
			{
				Name: "rule",
				Key: map[string]string{
					"sequence-id": getOCSequenceID(key),
				},
			},
			{Name: "state"},
			{Name: leafName},
		},
	}
}

// extractEndpoint extracts the decap IP from an Arista Smash key formatted like "<decapIp>_<tunnelType>".
// If no "_" delimiter is found, it returns the raw key unchanged.
func extractEndpoint(rawKey string) string {
	if idx := strings.LastIndex(rawKey, "_"); idx != -1 {
		return rawKey[:idx]
	}
	return rawKey
}

func (i *impl) updateHandler(prefix *gnmipb.Path, update *gnmipb.Update) []*gnmipb.Update {
	fullPath := ftutilities.Join(prefix, update.GetPath())
	if fullPath.GetOrigin() != "eos_native" {
		return nil
	}
	elems := fullPath.GetElem()
	if len(elems) != 8 {
		return nil
	}

	leaf := elems[7].GetName()
	if leaf != "pkts" && leaf != "octets" {
		return nil
	}

	prefixPath := &gnmipb.Path{
		Origin: fullPath.GetOrigin(),
		Elem:   elems[:7],
	}
	if !ftutilities.MatchPath(prefixPath, decapCommonPrefixPattern) {
		return nil
	}

	rawKey := elems[6].GetName()
	key := extractEndpoint(rawKey)

	var ocPath *gnmipb.Path
	switch leaf {
	case "pkts":
		ocPath = getOCPktsPath(key)
	case "octets":
		ocPath = getOCOctetsPath(key)
	}

	updates := make([]*gnmipb.Update, 0, 2)
	updates = append(updates, &gnmipb.Update{
		Path: ocPath,
		Val:  update.GetVal(),
	})

	if ip, isIPv6 := parseIP(key); ip != nil {
		updates = append(updates, &gnmipb.Update{
			Path: getOCDestAddrPath(key, isIPv6),
			Val: &gnmipb.TypedValue{
				Value: &gnmipb.TypedValue_StringVal{
					StringVal: key,
				},
			},
		})
	}

	return updates
}

func (i *impl) deleteHandler(n *gnmipb.Notification) []*gnmipb.Path {
	prefix := n.GetPrefix()
	deletes := n.GetDelete()

	var returnDeletes []*gnmipb.Path
	for _, deletePath := range deletes {
		fullPath := ftutilities.Join(prefix, deletePath)
		if fullPath.GetOrigin() != "eos_native" {
			continue
		}
		elems := fullPath.GetElem()

		// If higher-level paths are deleted (len(elems) < 7), this functional translator currently
		// ignores it. That might cause derivative OpenConfig paths to linger. This is an
		// accepted stateless limitation.
		if len(elems) == 7 {
			if ftutilities.MatchPath(fullPath, decapCommonPrefixPattern) {
				key := extractEndpoint(elems[6].GetName())
				returnDeletes = append(returnDeletes, getOCPktsPath(key), getOCOctetsPath(key))
				if ip, isIPv6 := parseIP(key); ip != nil {
					returnDeletes = append(returnDeletes, getOCDestAddrPath(key, isIPv6))
				}
			}
			continue
		}

		// If 'pkts' or 'octets' leaves are deleted individually (len(elems) == 8), the 'destination-address'
		// OpenConfig path is not deleted. Because both derived paths share this destination address,
		// deleting it when only one of the counters goes away might be wrong. But if both are deleted
		// individually (and the parent node isn't explicitly deleted), the destination-address might
		// linger. This is an accepted stateless limitation.
		if len(elems) == 8 {
			leaf := elems[7].GetName()
			if leaf != "pkts" && leaf != "octets" {
				continue
			}
			prefixPath := &gnmipb.Path{
				Origin: fullPath.GetOrigin(),
				Elem:   elems[:7],
			}
			if !ftutilities.MatchPath(prefixPath, decapCommonPrefixPattern) {
				continue
			}

			key := extractEndpoint(elems[6].GetName())
			switch leaf {
			case "pkts":
				returnDeletes = append(returnDeletes, getOCPktsPath(key))
			case "octets":
				returnDeletes = append(returnDeletes, getOCOctetsPath(key))
			}
		}
	}
	return returnDeletes
}

func (i *impl) translate(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
	notification := sr.GetUpdate()
	if notification == nil {
		return nil, nil
	}
	prefix := notification.GetPrefix()
	// Each input update yields up to 2 OpenConfig updates: the counter and the destination-address.
	outgoingUpdates := make([]*gnmipb.Update, 0, len(notification.GetUpdate())*2)
	for _, update := range notification.GetUpdate() {
		if updates := i.updateHandler(prefix, update); len(updates) > 0 {
			outgoingUpdates = append(outgoingUpdates, updates...)
		}
	}

	outgoingDeletes := i.deleteHandler(notification)

	if len(outgoingUpdates) == 0 && len(outgoingDeletes) == 0 {
		return nil, nil
	}

	outgoingNotification := &gnmipb.Notification{
		Timestamp: notification.GetTimestamp(),
		Prefix: &gnmipb.Path{
			Origin: "openconfig",
			Target: notification.GetPrefix().GetTarget(),
			Elem:   ocCommonPrefix,
		},
		Update: outgoingUpdates,
		Delete: outgoingDeletes,
	}
	return &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: outgoingNotification,
		},
	}, nil
}
