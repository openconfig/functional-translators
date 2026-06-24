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

// Package aristacfmpm translates the CFM PM data from native to openconfig.
package aristacfmpm

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	log "github.com/golang/glog"
	"github.com/openconfig/functional-translators/ftconsts"
	"github.com/openconfig/functional-translators/ftutilities"
	"github.com/openconfig/functional-translators/translator"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

const (
	rootSmash = "Smash"
	rootSysdb = "Sysdb"
	elemCfm   = "cfm"

	// Smash paths
	elemMepSmashTable     = "mepSmashTable"
	elemDmMiStatsCurrent  = "dmMiStatsCurrent"
	elemSlmMiStatsCurrent = "slmMiStatsCurrent"
	elemDmMiStats         = "dmMiStats"
	elemSlmMiStats        = "slmMiStats"
	elemDelayTwoWayAvg    = "delayTwoWayAvg"
	elemForwardAvgFlr     = "forwardAvgFlr"
	elemBackwardAvgFlr    = "backwardAvgFlr"

	// Sysdb paths
	elemConfig         = "config"
	elemMdConfig       = "mdConfig"
	elemMaConfig       = "maConfig"
	elemCfmProfileName = "cfmProfileName"

	// Metrics
	metricDelay    = "delay"
	metricLossFar  = "lossFar"
	metricLossNear = "lossNear"
)

var (
	translateMap = map[string][]string{
		"/openconfig/oam/cfm/domains/maintenance-domain/maintenance-associations/maintenance-association/mep-endpoints/mep-endpoint/pm-profiles/pm-profile/state/delay-measurement-state/frame-delay-two-way-average": {
			"/eos_native/Smash/cfm/mepSmashTable/dmMiStatsCurrent",
			"/eos_native/Sysdb/cfm/config/mdConfig",
			"/eos_native/Sysdb/cfm/status/mdStatus",
		},
		"/openconfig/oam/cfm/domains/maintenance-domain/maintenance-associations/maintenance-association/mep-endpoints/mep-endpoint/pm-profiles/pm-profile/state/loss-measurement-state/far-end-average-frame-loss-ratio": {
			"/eos_native/Smash/cfm/mepSmashTable/slmMiStatsCurrent",
			"/eos_native/Sysdb/cfm/config/mdConfig",
			"/eos_native/Sysdb/cfm/status/mdStatus",
		},
		"/openconfig/oam/cfm/domains/maintenance-domain/maintenance-associations/maintenance-association/mep-endpoints/mep-endpoint/pm-profiles/pm-profile/state/loss-measurement-state/near-end-average-frame-loss-ratio": {
			"/eos_native/Smash/cfm/mepSmashTable/slmMiStatsCurrent",
			"/eos_native/Sysdb/cfm/config/mdConfig",
			"/eos_native/Sysdb/cfm/status/mdStatus",
		},
	}
	updatePathPatterns = []*gnmipb.Path{
		{
			Origin: "eos_native",
			Elem: []*gnmipb.PathElem{
				{Name: rootSmash}, {Name: elemCfm}, {Name: elemMepSmashTable}, {Name: elemDmMiStatsCurrent},
				{Name: "*"}, // key containing domainID, assocID, and localMepID
				{Name: elemDmMiStats},
				{Name: elemDelayTwoWayAvg},
			},
		},
		{
			Origin: "eos_native",
			Elem: []*gnmipb.PathElem{
				{Name: rootSmash}, {Name: elemCfm}, {Name: elemMepSmashTable}, {Name: elemSlmMiStatsCurrent},
				{Name: "*"}, // key containing domainID, assocID, and localMepID
				{Name: elemSlmMiStats},
				{Name: elemForwardAvgFlr},
			},
		},
		{
			Origin: "eos_native",
			Elem: []*gnmipb.PathElem{
				{Name: rootSmash}, {Name: elemCfm}, {Name: elemMepSmashTable}, {Name: elemSlmMiStatsCurrent},
				{Name: "*"}, // key containing domainID, assocID, and localMepID
				{Name: elemSlmMiStats},
				{Name: elemBackwardAvgFlr},
			},
		},

		{
			Origin: "eos_native",
			Elem: []*gnmipb.PathElem{
				{Name: rootSysdb}, {Name: elemCfm}, {Name: elemConfig}, {Name: elemMdConfig},
				{Name: "*"}, // domain name
				{Name: elemMaConfig},
				{Name: "*"}, // association name
				{Name: elemCfmProfileName},
			},
		},
		{
			Origin: "eos_native",
			Elem: []*gnmipb.PathElem{
				{Name: rootSysdb}, {Name: elemCfm}, {Name: "status"}, {Name: "mdStatus"},
				{Name: "*"}, // domain
				{Name: "maStatus"},
				{Name: "*"}, // maName
				{Name: "localMepStatus"},
				{Name: "*"}, // localMepID
				{Name: "rdiTxCondition"},
				{Name: "rdi"},
			},
		},
	}
	deletePathPatterns = []*gnmipb.Path{
		{
			Origin: "eos_native",
			Elem: []*gnmipb.PathElem{
				{Name: rootSmash}, {Name: elemCfm}, {Name: elemMepSmashTable}, {Name: elemDmMiStatsCurrent},
				{Name: "*"}, // key containing domainID, assocID, and localMepID
			},
		},
		{
			Origin: "eos_native",
			Elem: []*gnmipb.PathElem{
				{Name: rootSmash}, {Name: elemCfm}, {Name: elemMepSmashTable}, {Name: elemSlmMiStatsCurrent},
				{Name: "*"}, // key containing domainID, assocID, and localMepID
			},
		},

		{
			Origin: "eos_native",
			Elem: []*gnmipb.PathElem{
				{Name: rootSysdb}, {Name: elemCfm}, {Name: elemConfig}, {Name: elemMdConfig},
				{Name: "*"}, // domain name
				{Name: elemMaConfig},
				{Name: "*"}, // association name
				{Name: elemCfmProfileName},
			},
		},
		{
			Origin: "eos_native",
			Elem: []*gnmipb.PathElem{
				{Name: rootSysdb}, {Name: elemCfm}, {Name: "status"}, {Name: "mdStatus"},
				{Name: "*"}, // domain
				{Name: "maStatus"},
				{Name: "*"}, // maName
				{Name: "localMepStatus"},
				{Name: "*"}, // localMepID
			},
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
			ID:               ftconsts.AristaCFMPMFunctionalTranslator,
			Translate:        i.translate,
			OutputToInputMap: ftutilities.MustStringMapPaths(translateMap),
			Metadata: []*translator.FTMetadata{
				{
					Vendor: ftconsts.VendorArista,
					SoftwareVersionRange: &translator.SWRange{
						InclusiveMin: "4.33.0F",
						ExclusiveMax: "4.37",
					},
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("Failed to create Arista CFM PM functional translator: %v", err)
	}
	return ft
}

// valueToString converts a TypedValue to a string.
// It strictly expects string or JSON values for the profile name.
func valueToString(tv *gnmipb.TypedValue) string {
	if tv == nil {
		return ""
	}
	var s string
	switch v := tv.GetValue().(type) {
	case *gnmipb.TypedValue_StringVal:
		return v.StringVal
	case *gnmipb.TypedValue_JsonVal:
		s = string(v.JsonVal)
	default:
		return ""
	}
	// Try parsing JSON object wrapper {"value": "..."}
	if strings.HasPrefix(strings.TrimSpace(s), "{") {
		var wrapper struct {
			Value string `json:"value"`
		}
		if err := json.Unmarshal([]byte(s), &wrapper); err == nil && wrapper.Value != "" {
			return wrapper.Value
		}
	}
	// Handle quoted strings
	if strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") {
		return s[1 : len(s)-1]
	}
	return s
}

// parseIntSlice converts a space-separated string of numbers into a byte slice string.
// It stops parsing if a zero (null terminator) is encountered.
func parseIntSlice(s string) (string, error) {
	parts := strings.Split(s, " ")
	var result strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		i, err := strconv.Atoi(p)
		if err != nil {
			return "", err
		}
		if i == 0 { // Null terminator
			break
		}
		result.WriteByte(byte(i))
	}
	return result.String(), nil
}

// parseSmashKey parses the complex Smash key format.
// Format: <id>_Array{base: 0, slice: [<assocID_ints>]}_maNameFormatShortInt_<id>_Array{base: 0, slice: [<domainID_ints>]}_mdNameFormatNoName_<localMepID>
func parseSmashKey(key string) (assocID string, domainID string, localMEPID string, err error) {
	const (
		arrayStart = "_Array{base: 0, slice: ["
		arrayEnd   = "]}"
		maFormat   = "_maNameFormatShortInt_"
		mdFormat   = "_mdNameFormatNoName_"
	)

	// Helper to extract content between start and end markers
	extract := func(s, start, end string) (string, string, bool) {
		idxStart := strings.Index(s, start)
		if idxStart == -1 {
			return "", "", false
		}
		rest := s[idxStart+len(start):]
		idxEnd := strings.Index(rest, end)
		if idxEnd == -1 {
			return "", "", false
		}
		return rest[:idxEnd], rest[idxEnd+len(end):], true
	}

	// 1. Skip first ID and find first array (AssocID)
	assocVal, rest, ok := extract(key, arrayStart, arrayEnd)
	if !ok {
		return "", "", "", fmt.Errorf("invalid key format: missing assocID array")
	}

	// 2. Expect maNameFormatShortInt and second array (DomainID)
	// rest should be: _maNameFormatShortInt_<id>_Array{base: 0, slice: [<domainID_ints>]}_mdNameFormatNoName_<localMepID>
	domainVal, rest, ok := extract(rest, arrayStart, arrayEnd)
	if !ok {
		return "", "", "", fmt.Errorf("invalid key format: missing domainID array")
	}

	// 3. Extract localMEPID from the end.
	idxLocal := strings.Index(rest, mdFormat)
	if idxLocal == -1 {
		return "", "", "", fmt.Errorf("invalid key format: missing localMEPID marker")
	}
	remaining := rest[idxLocal+len(mdFormat):]
	// Match strictly digits, similar to (\d+) in the original regex.
	endOfDigits := 0
	for i, r := range remaining {
		if r < '0' || r > '9' {
			endOfDigits = i
			break
		}
		endOfDigits = i + 1
	}
	localMEPID = remaining[:endOfDigits]
	if len(localMEPID) == 0 {
		return "", "", "", fmt.Errorf("invalid key format: missing or invalid localMEPID")
	}

	assocID, err = parseIntSlice(assocVal)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to parse assocID: %v", err)
	}
	domainID, err = parseIntSlice(domainVal)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to parse domainID: %v", err)
	}

	return assocID, domainID, localMEPID, nil
}

func (i *impl) updateHandler(notification *gnmipb.Notification) ([]*gnmipb.Update, error) {
	if len(notification.GetUpdate()) == 0 {
		return nil, nil
	}
	prefix := notification.GetPrefix()
	target := prefix.GetTarget()
	var updates []*gnmipb.Update
	targetCFM := ftutilities.CfmMap.CreateOrUpdateTargetCFMInfo(target)
	for _, update := range notification.GetUpdate() {
		fullPath := ftutilities.Join(prefix, update.GetPath())
		for _, path := range updatePathPatterns {
			if !ftutilities.MatchPath(fullPath, path) {
				continue
			}
			elems := fullPath.GetElem()
			// MatchPath guarantees element count and structure.

			switch elems[0].GetName() {
			case rootSmash:
				// Path: Smash/cfm/mepSmashTable/<statType>/<key>/<statGroup>/<metric>
				key := elems[4].GetName()
				assocID, domainID, parsedLocalMEPID, err := parseSmashKey(key)
				if err != nil {
					return nil, fmt.Errorf("error parsing smash key %q: %v", key, err)
				}

				cacheKey := ftutilities.CFMCacheKey{DomainID: domainID, AssocID: assocID}
				profileName, ok := targetCFM.ProfileName(cacheKey)
				if !ok {
					continue
				}
				localMEPID, ok := targetCFM.LocalMEPID(cacheKey)
				if !ok {
					localMEPID = parsedLocalMEPID
				}

				var metricType string
				switch elems[3].GetName() {
				case elemDmMiStatsCurrent:
					metricType = metricDelay
				case elemSlmMiStatsCurrent:
					if len(elems) < 7 {
						continue
					}
					switch elems[6].GetName() {
					case elemForwardAvgFlr:
						metricType = metricLossFar
					case elemBackwardAvgFlr:
						metricType = metricLossNear
					default:
						continue
					}
				default:
					continue
				}

				val, err := i.convertMetricValue(update.GetVal(), metricType)
				if err != nil {
					log.Errorf("Failed to convert metric value %v: %v", update.GetVal(), err)
					continue
				}

				outgoingUpdate := &gnmipb.Update{
					Path: pmPath(domainID, assocID, localMEPID, profileName, metricType),
					Val:  val,
				}
				updates = append(updates, outgoingUpdate)

			case rootSysdb:
				domainID := elems[4].GetName()
				assocKey := elems[6].GetName()
				const assocPrefix = "maNameFormatShortInt_"
				if !strings.HasPrefix(assocKey, assocPrefix) {
					log.Warningf("Unexpected assoc key format: %q", assocKey)
					continue
				}
				assocID := strings.TrimPrefix(assocKey, assocPrefix)
				cacheKey := ftutilities.CFMCacheKey{DomainID: domainID, AssocID: assocID}

				if elems[2].GetName() == elemConfig {
					// Path: Sysdb/cfm/config/mdConfig/<domain>/maConfig/<assoc>/cfmProfileName
					profileName := valueToString(update.GetVal())
					if profileName == "" {
						continue // Ignore empty or invalid values
					}
					targetCFM.SetProfileName(cacheKey, profileName)
				} else { // Path matches Sysdb status pattern based on updatePathPatterns
					// Path: Sysdb/cfm/status/mdStatus/<domain>/maStatus/<assoc>/localMepStatus/<id>/...
					if len(elems) >= 9 {
						localMEPID := elems[8].GetName()
						targetCFM.SetLocalMEPID(cacheKey, localMEPID)
					}
				}
			}
			break
		}
	}
	return updates, nil
}

func (i *impl) deleteHandler(notification *gnmipb.Notification) ([]*gnmipb.Path, error) {
	prefix := notification.GetPrefix()
	target := prefix.GetTarget()
	var deletes []*gnmipb.Path
	targetCFM := ftutilities.CfmMap.CreateOrUpdateTargetCFMInfo(target)
	// Loop 1: Process Smash deletes first (telemetry)
	for _, del := range notification.GetDelete() {
		fullPath := ftutilities.Join(prefix, del)
		for _, pattern := range deletePathPatterns {
			if !ftutilities.MatchPath(fullPath, pattern) {
				continue
			}
			elems := fullPath.GetElem()
			if elems[0].GetName() != rootSmash {
				continue
			}

			if len(elems) < 5 {
				continue
			}

			// Process Smash delete
			key := elems[4].GetName()
			assocID, domainID, parsedLocalMEPID, err := parseSmashKey(key)
			if err != nil {
				log.Warningf("Error parsing smash key %q: %v", key, err)
				continue
			}

			cacheKey := ftutilities.CFMCacheKey{DomainID: domainID, AssocID: assocID}
			profileName, ok := targetCFM.ProfileName(cacheKey)
			if !ok {
				continue
			}
			localMEPID, ok := targetCFM.LocalMEPID(cacheKey)
			if !ok {
				localMEPID = parsedLocalMEPID
			}
			switch elems[3].GetName() {
			case elemDmMiStatsCurrent:
				deletes = append(deletes, pmPath(domainID, assocID, localMEPID, profileName, metricDelay))
			case elemSlmMiStatsCurrent:
				deletes = append(deletes, pmPath(domainID, assocID, localMEPID, profileName, metricLossFar))
				deletes = append(deletes, pmPath(domainID, assocID, localMEPID, profileName, metricLossNear))
			default:
				continue
			}
			break
		}
	}

	// Loop 2: Process Sysdb deletes last (metadata/cache)
	for _, del := range notification.GetDelete() {
		fullPath := ftutilities.Join(prefix, del)
		for _, pattern := range deletePathPatterns {
			if !ftutilities.MatchPath(fullPath, pattern) {
				continue
			}
			elems := fullPath.GetElem()
			if elems[0].GetName() != rootSysdb {
				continue
			}

			// Process Sysdb delete
			if len(elems) < 7 {
				log.Warningf("Path too short to extract domain and assoc IDs: %v", fullPath)
				continue
			}
			domainID := elems[4].GetName()
			assocKey := elems[6].GetName()
			const assocPrefix = "maNameFormatShortInt_"
			if !strings.HasPrefix(assocKey, assocPrefix) {
				log.Warningf("Unexpected assoc key format: %q", assocKey)
				continue
			}
			assocID := strings.TrimPrefix(assocKey, assocPrefix)

			cacheKey := ftutilities.CFMCacheKey{DomainID: domainID, AssocID: assocID}
			if elems[2].GetName() == elemConfig {
				targetCFM.DeleteProfileName(cacheKey)
			} else { // Path matches Sysdb status pattern based on deletePathPatterns
				if len(elems) >= 9 {
					deletedMEPID := elems[8].GetName()
					cachedMEPID, ok := targetCFM.LocalMEPID(cacheKey)
					if ok && cachedMEPID == deletedMEPID {
						targetCFM.DeleteLocalMEPID(cacheKey)
					}
				}
			}
			break
		}
	}
	return deletes, nil
}

// pmPath returns a gNMI path for the PM update.
func pmPath(domainID, assocID, localMEPID, profileName, metricType string) *gnmipb.Path {
	path := &gnmipb.Path{
		Elem: []*gnmipb.PathElem{
			{Name: "oam"},
			{Name: "cfm"},
			{Name: "domains"},
			{
				Name: "maintenance-domain",
				Key: map[string]string{
					"domain-name": domainID,
				},
			},
			{Name: "maintenance-associations"},
			{
				Name: "maintenance-association",
				Key: map[string]string{
					"association-name": assocID,
				},
			},
			{Name: "mep-endpoints"},
			{
				Name: "mep-endpoint",
				Key: map[string]string{
					"mep-id": localMEPID,
				},
			},
			{Name: "pm-profiles"},
			{
				Name: "pm-profile",
				Key: map[string]string{
					"profile-name": profileName,
				},
			},
			{Name: "state"},
		},
	}

	switch metricType {
	case metricDelay:
		path.Elem = append(path.Elem, &gnmipb.PathElem{Name: "delay-measurement-state"})
		path.Elem = append(path.Elem, &gnmipb.PathElem{Name: "frame-delay-two-way-average"})
	case metricLossFar, metricLossNear:
		path.Elem = append(path.Elem, &gnmipb.PathElem{Name: "loss-measurement-state"})
		var leafName string
		if metricType == metricLossNear {
			leafName = "near-end-average-frame-loss-ratio"
		} else {
			leafName = "far-end-average-frame-loss-ratio"
		}
		path.Elem = append(path.Elem, &gnmipb.PathElem{Name: leafName})
	}

	return path
}

func (i *impl) translate(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
	notification := sr.GetUpdate()
	if notification == nil {
		return nil, nil
	}
	updates, err := i.updateHandler(notification)
	if err != nil {
		return nil, err
	}
	deletes, err := i.deleteHandler(notification)
	if err != nil {
		return nil, err
	}
	if len(updates) == 0 && len(deletes) == 0 {
		return nil, nil
	}
	outgoingNotification := &gnmipb.Notification{
		Timestamp: notification.GetTimestamp(),
		Prefix: &gnmipb.Path{
			Origin: "openconfig",
			Target: notification.GetPrefix().GetTarget(),
		},
		Update: updates,
		Delete: deletes,
	}
	return &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: outgoingNotification,
		},
	}, nil
}

// convertMetricValue parses native telemetry value to uint64 for openconfig.
func (i *impl) convertMetricValue(tv *gnmipb.TypedValue, metricType string) (*gnmipb.TypedValue, error) {
	if tv == nil {
		return nil, fmt.Errorf("empty metric value")
	}

	var valFloat float64
	switch v := tv.GetValue().(type) {
	case *gnmipb.TypedValue_DoubleVal:
		valFloat = v.DoubleVal
	case *gnmipb.TypedValue_IntVal:
		valFloat = float64(v.IntVal)
	case *gnmipb.TypedValue_UintVal:
		valFloat = float64(v.UintVal)
	case *gnmipb.TypedValue_StringVal:
		f, err := strconv.ParseFloat(v.StringVal, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse string metric value %q: %w", v.StringVal, err)
		}
		valFloat = f
	default:
		return nil, fmt.Errorf("unsupported metric value type: %T", v)
	}

	if valFloat < 0 {
		log.Warningf("Negative metric value received: %f for type %s. Setting to 0.", valFloat, metricType)
		valFloat = 0
	}

	var uintVal uint64
	switch metricType {
	case metricDelay:
		// Convert seconds to microseconds
		uintVal = uint64(math.Round(valFloat * 1e6))
	case metricLossFar, metricLossNear:
		// native string/number is a raw ratio like "0" or "123".
		uintVal = uint64(math.Round(valFloat))
	default:
		return nil, fmt.Errorf("unknown metric type: %s", metricType)
	}

	return &gnmipb.TypedValue{
		Value: &gnmipb.TypedValue_UintVal{
			UintVal: uintVal,
		},
	}, nil
}
