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

// Package ciscoxrqos translates ciscoxr qos to openconfig.
package ciscoxrqos

import (
	"fmt"
	"strings"

	log "github.com/golang/glog"
	ocqos "github.com/openconfig/functional-translators/ciscoxr/ciscoxrqos/yang/openconfig"
	"github.com/openconfig/functional-translators/ftconsts"
	"github.com/openconfig/functional-translators/ftutilities"
	"github.com/openconfig/functional-translators/translator"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

type outStats struct {
	droppedOctets  []uint64
	droppedPkts    []uint64
	transmitOctets []uint64
	transmitPkts   []uint64
	className      []string
}

type inStats struct {
	matchedOctets []uint64
	matchedPkts   []uint64
	className     []string
}

var (
	translateMap = map[string][]string{
		"/openconfig/qos/interfaces/interface/output/queues/queue/state/dropped-octets": {
			"/Cisco-IOS-XR-qos-ma-oper/qos/interface-table/interface/output/service-policy-names/service-policy-instance/statistics",
			"/Cisco-IOS-XR-qos-ma-oper/qos/interface-table/interface/member-interfaces/member-interface/output/service-policy-names/service-policy-instance/statistics",
		},
		"/openconfig/qos/interfaces/interface/output/queues/queue/state/dropped-pkts": {
			"/Cisco-IOS-XR-qos-ma-oper/qos/interface-table/interface/output/service-policy-names/service-policy-instance/statistics",
			"/Cisco-IOS-XR-qos-ma-oper/qos/interface-table/interface/member-interfaces/member-interface/output/service-policy-names/service-policy-instance/statistics",
		},
		"/openconfig/qos/interfaces/interface/output/queues/queue/state/transmit-octets": {
			"/Cisco-IOS-XR-qos-ma-oper/qos/interface-table/interface/output/service-policy-names/service-policy-instance/statistics",
			"/Cisco-IOS-XR-qos-ma-oper/qos/interface-table/interface/member-interfaces/member-interface/output/service-policy-names/service-policy-instance/statistics",
		},
		"/openconfig/qos/interfaces/interface/output/queues/queue/state/transmit-pkts": {
			"/Cisco-IOS-XR-qos-ma-oper/qos/interface-table/interface/output/service-policy-names/service-policy-instance/statistics",
			"/Cisco-IOS-XR-qos-ma-oper/qos/interface-table/interface/member-interfaces/member-interface/output/service-policy-names/service-policy-instance/statistics",
		},
		"/openconfig/qos/interfaces/interface/input/classifiers/classifier/terms/term/state/matched-octets": {
			"/Cisco-IOS-XR-qos-ma-oper/qos/interface-table/interface/input/service-policy-names/service-policy-instance/statistics",
		},
		"/openconfig/qos/interfaces/interface/input/classifiers/classifier/terms/term/state/matched-packets": {
			"/Cisco-IOS-XR-qos-ma-oper/qos/interface-table/interface/input/service-policy-names/service-policy-instance/statistics",
		},
	}
	paths       = ftutilities.MustStringMapPaths(translateMap)
	nativePaths = []*gnmipb.Path{
		{
			Origin: "Cisco-IOS-XR-qos-ma-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "qos"}, {Name: "interface-table"}, {Name: "interface"}, {Name: "output"},
				{Name: "service-policy-names"}, {Name: "service-policy-instance"},
				{Name: "statistics"}, {Name: "class-stats"}, {Name: "class-name"},
			},
		},
		{
			Origin: "Cisco-IOS-XR-qos-ma-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "qos"}, {Name: "interface-table"}, {Name: "interface"}, {Name: "member-interfaces"},
				{Name: "member-interface"}, {Name: "output"}, {Name: "service-policy-names"}, {Name: "service-policy-instance"},
				{Name: "statistics"}, {Name: "class-stats"}, {Name: "class-name"},
			},
		},
		{
			Origin: "Cisco-IOS-XR-qos-ma-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "qos"}, {Name: "interface-table"}, {Name: "interface"}, {Name: "input"},
				{Name: "service-policy-names"}, {Name: "service-policy-instance"},
				{Name: "statistics"}, {Name: "class-stats"}, {Name: "class-name"},
			},
		},
		{
			Origin: "Cisco-IOS-XR-qos-ma-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "qos"}, {Name: "interface-table"}, {Name: "interface"}, {Name: "output"},
				{Name: "service-policy-names"}, {Name: "service-policy-instance"},
				{Name: "statistics"}, {Name: "class-stats"}, {Name: "general-stats"},
				{Name: "total-drop-bytes"},
			},
		},
		{
			Origin: "Cisco-IOS-XR-qos-ma-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "qos"}, {Name: "interface-table"}, {Name: "interface"}, {Name: "output"},
				{Name: "service-policy-names"}, {Name: "service-policy-instance"},
				{Name: "statistics"}, {Name: "class-stats"}, {Name: "general-stats"},
				{Name: "total-drop-packets"},
			},
		},
		{
			Origin: "Cisco-IOS-XR-qos-ma-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "qos"}, {Name: "interface-table"}, {Name: "interface"}, {Name: "output"},
				{Name: "service-policy-names"}, {Name: "service-policy-instance"},
				{Name: "statistics"}, {Name: "class-stats"}, {Name: "general-stats"},
				{Name: "transmit-bytes"},
			},
		},
		{
			Origin: "Cisco-IOS-XR-qos-ma-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "qos"}, {Name: "interface-table"}, {Name: "interface"}, {Name: "output"},
				{Name: "service-policy-names"}, {Name: "service-policy-instance"},
				{Name: "statistics"}, {Name: "class-stats"}, {Name: "general-stats"},
				{Name: "transmit-packets"},
			},
		},
		{
			Origin: "Cisco-IOS-XR-qos-ma-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "qos"}, {Name: "interface-table"}, {Name: "interface"}, {Name: "member-interfaces"},
				{Name: "member-interface"}, {Name: "output"}, {Name: "service-policy-names"}, {Name: "service-policy-instance"},
				{Name: "statistics"}, {Name: "class-stats"}, {Name: "general-stats"},
				{Name: "total-drop-bytes"},
			},
		},
		{
			Origin: "Cisco-IOS-XR-qos-ma-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "qos"}, {Name: "interface-table"}, {Name: "interface"}, {Name: "member-interfaces"},
				{Name: "member-interface"}, {Name: "output"}, {Name: "service-policy-names"}, {Name: "service-policy-instance"},
				{Name: "statistics"}, {Name: "class-stats"}, {Name: "general-stats"},
				{Name: "total-drop-packets"},
			},
		},
		{
			Origin: "Cisco-IOS-XR-qos-ma-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "qos"}, {Name: "interface-table"}, {Name: "interface"}, {Name: "member-interfaces"},
				{Name: "member-interface"}, {Name: "output"}, {Name: "service-policy-names"}, {Name: "service-policy-instance"},
				{Name: "statistics"}, {Name: "class-stats"}, {Name: "general-stats"},
				{Name: "transmit-bytes"},
			},
		},
		{
			Origin: "Cisco-IOS-XR-qos-ma-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "qos"}, {Name: "interface-table"}, {Name: "interface"}, {Name: "member-interfaces"},
				{Name: "member-interface"}, {Name: "output"}, {Name: "service-policy-names"}, {Name: "service-policy-instance"},
				{Name: "statistics"}, {Name: "class-stats"}, {Name: "general-stats"},
				{Name: "transmit-packets"},
			},
		},
		{
			Origin: "Cisco-IOS-XR-qos-ma-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "qos"}, {Name: "interface-table"}, {Name: "interface"}, {Name: "input"},
				{Name: "service-policy-names"}, {Name: "service-policy-instance"},
				{Name: "statistics"}, {Name: "class-stats"}, {Name: "general-stats"},
				{Name: "pre-policy-matched-bytes"},
			},
		},
		{
			Origin: "Cisco-IOS-XR-qos-ma-oper",
			Elem: []*gnmipb.PathElem{
				{Name: "qos"}, {Name: "interface-table"}, {Name: "interface"}, {Name: "input"},
				{Name: "service-policy-names"}, {Name: "service-policy-instance"},
				{Name: "statistics"}, {Name: "class-stats"}, {Name: "general-stats"},
				{Name: "pre-policy-matched-packets"},
			},
		},
	}
)

// validates the leaves and builds stats structs
func buildStats(prefix *gnmipb.Path, leaves []*gnmipb.Update) (map[string]*outStats, map[string]*inStats) {
	intfOutStats := make(map[string]*outStats)
	intfInStats := make(map[string]*inStats)
	for _, leaf := range leaves {
		path := ftutilities.Join(prefix, leaf.GetPath())
		if !ftutilities.PathInList(path, nativePaths) {
			continue
		}
		elems := path.GetElem()
		switch {
		case elems[3].GetName() == "output" || elems[5].GetName() == "output":
			intfOutStats = buildOutputStats(prefix, leaf, intfOutStats)
		case elems[3].GetName() == "input":
			intfInStats = buildInputStats(prefix, leaf, intfInStats)
		}
	}
	return intfOutStats, intfInStats
}

// allEqual returns true if all the numbers are equal.
func allEqual(nums ...int) bool {
	if len(nums) == 0 {
		return true
	}
	for _, n := range nums {
		if n != nums[0] {
			return false
		}
	}
	return true
}

// build output stats structs
func buildOutputStats(prefix *gnmipb.Path, leaf *gnmipb.Update, intfOutStats map[string]*outStats) map[string]*outStats {
	path := ftutilities.Join(prefix, leaf.GetPath())
	elems := path.GetElem()
	var intfName string
	var startIndex int
	if elems[4].GetName() == "member-interface" {
		intfName = elems[4].GetKey()["interface-name"]
		startIndex = 10
	} else {
		intfName = elems[2].GetKey()["interface-name"]
		startIndex = 8
	}
	t, ok := intfOutStats[intfName]
	if !ok {
		t = &outStats{
			droppedOctets:  []uint64{},
			droppedPkts:    []uint64{},
			transmitOctets: []uint64{},
			transmitPkts:   []uint64{},
			className:      []string{},
		}
		intfOutStats[intfName] = t
	}
	switch elems[startIndex].GetName() {
	case "class-name":
		t.className = append(t.className, leaf.GetVal().GetStringVal())
	case "general-stats":
		switch elems[startIndex+1].GetName() {
		case "total-drop-bytes":
			t.droppedOctets = append(t.droppedOctets, leaf.GetVal().GetUintVal())
		case "total-drop-packets":
			t.droppedPkts = append(t.droppedPkts, leaf.GetVal().GetUintVal())
		case "transmit-bytes":
			t.transmitOctets = append(t.transmitOctets, leaf.GetVal().GetUintVal())
		case "transmit-packets":
			t.transmitPkts = append(t.transmitPkts, leaf.GetVal().GetUintVal())
		}
	}
	return intfOutStats
}

// build input stats structs
func buildInputStats(prefix *gnmipb.Path, leaf *gnmipb.Update, intfInStats map[string]*inStats) map[string]*inStats {
	path := ftutilities.Join(prefix, leaf.GetPath())
	elems := path.GetElem()
	intfName := elems[2].GetKey()["interface-name"]
	t, ok := intfInStats[intfName]
	if !ok {
		t = &inStats{
			matchedOctets: []uint64{},
			matchedPkts:   []uint64{},
			className:     []string{},
		}
		intfInStats[intfName] = t
	}
	switch elems[8].GetName() {
	case "class-name":
		t.className = append(t.className, leaf.GetVal().GetStringVal())
	case "general-stats":
		switch elems[9].GetName() {
		case "pre-policy-matched-bytes":
			t.matchedOctets = append(t.matchedOctets, leaf.GetVal().GetUintVal())
		case "pre-policy-matched-packets":
			t.matchedPkts = append(t.matchedPkts, leaf.GetVal().GetUintVal())
		}
	}
	return intfInStats
}

// New creates a functional translator.
func New() *translator.FunctionalTranslator {
	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               ftconsts.CiscoXRQosTranslator,
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
		log.Fatalf("Failed to create Cisco QoS functional translator: %v", err)
	}
	return ft
}

func translate(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
	if sr.GetUpdate() == nil {
		return nil, nil
	}
	intfOutStats, intfInStats := buildStats(sr.GetUpdate().GetPrefix(), sr.GetUpdate().GetUpdate())
	n := sr.GetUpdate()
	qosRoot := &ocqos.Device{}
	for intfName, intfOutStat := range intfOutStats {
		if !allEqual(len(intfOutStat.className), len(intfOutStat.droppedOctets), len(intfOutStat.droppedPkts), len(intfOutStat.transmitOctets), len(intfOutStat.transmitPkts)) {
			return nil, fmt.Errorf("mismatch len for interface %s output queue stats class", intfName)
		}
		intfOutput := qosRoot.GetOrCreateQos().GetOrCreateInterfaces().GetOrCreateInterface(intfName).GetOrCreateOutput()
		for i, className := range intfOutStat.className {
			queue := intfOutput.GetOrCreateQueues().GetOrCreateQueue(className)
			queue.GetOrCreateState().DroppedOctets = &intfOutStat.droppedOctets[i]
			queue.GetOrCreateState().DroppedPkts = &intfOutStat.droppedPkts[i]
			queue.GetOrCreateState().TransmitOctets = &intfOutStat.transmitOctets[i]
			queue.GetOrCreateState().TransmitPkts = &intfOutStat.transmitPkts[i]
		}
	}
	for intfName, intfInStat := range intfInStats {
		if !allEqual(len(intfInStat.className), len(intfInStat.matchedOctets), len(intfInStat.matchedPkts)) {
			return nil, fmt.Errorf("mismatch len for interface %s input queue stats class", intfName)
		}
		intfinput := qosRoot.GetOrCreateQos().GetOrCreateInterfaces().GetOrCreateInterface(intfName).GetOrCreateInput()
		for i, className := range intfInStat.className {
			nameParts := strings.Split(className, "-")
			if len(nameParts) < 2 {
				log.Warningf("wrong className %s does not have any parts separated by -", className)
				continue
			}
			switch nameParts[0] {
			case "inet6":
				classifier := intfinput.GetOrCreateClassifiers().GetOrCreateClassifier(ocqos.OpenconfigQos_Qos_Interfaces_Interface_Input_Classifiers_Classifier_Config_Type_IPV6)
				classifier.GetOrCreateTerms().GetOrCreateTerm(nameParts[len(nameParts)-1]).GetOrCreateState().MatchedOctets = &intfInStat.matchedOctets[i]
				classifier.GetOrCreateTerms().GetOrCreateTerm(nameParts[len(nameParts)-1]).GetOrCreateState().MatchedPackets = &intfInStat.matchedPkts[i]
			case "inet":
				classifier := intfinput.GetOrCreateClassifiers().GetOrCreateClassifier(ocqos.OpenconfigQos_Qos_Interfaces_Interface_Input_Classifiers_Classifier_Config_Type_IPV4)
				classifier.GetOrCreateTerms().GetOrCreateTerm(nameParts[len(nameParts)-1]).GetOrCreateState().MatchedOctets = &intfInStat.matchedOctets[i]
				classifier.GetOrCreateTerms().GetOrCreateTerm(nameParts[len(nameParts)-1]).GetOrCreateState().MatchedPackets = &intfInStat.matchedPkts[i]
			case "exp":
				classifier := intfinput.GetOrCreateClassifiers().GetOrCreateClassifier(ocqos.OpenconfigQos_Qos_Interfaces_Interface_Input_Classifiers_Classifier_Config_Type_MPLS)
				classifier.GetOrCreateTerms().GetOrCreateTerm(nameParts[len(nameParts)-1]).GetOrCreateState().MatchedOctets = &intfInStat.matchedOctets[i]
				classifier.GetOrCreateTerms().GetOrCreateTerm(nameParts[len(nameParts)-1]).GetOrCreateState().MatchedPackets = &intfInStat.matchedPkts[i]
			default:
				classifier := intfinput.GetOrCreateClassifiers().GetOrCreateClassifier(ocqos.OpenconfigQos_Qos_Interfaces_Interface_Input_Classifiers_Classifier_Config_Type_UNSET)
				classifier.GetOrCreateTerms().GetOrCreateTerm(nameParts[len(nameParts)-1]).GetOrCreateState().MatchedOctets = &intfInStat.matchedOctets[i]
				classifier.GetOrCreateTerms().GetOrCreateTerm(nameParts[len(nameParts)-1]).GetOrCreateState().MatchedPackets = &intfInStat.matchedPkts[i]
			}
		}
	}
	return ftutilities.FilterStructToState(qosRoot, n.GetTimestamp(), "openconfig", n.GetPrefix().GetTarget())
}
