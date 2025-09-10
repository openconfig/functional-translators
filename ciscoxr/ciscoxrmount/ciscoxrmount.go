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

// Package ciscoxrmount translates ciscoxr file systems to openconfig mount points.
package ciscoxrmount

import (
	"fmt"

	log "github.com/golang/glog"
	ocmount "github.com/openconfig/functional-translators/ciscoxr/ciscoxrmount/yang/openconfig"
	"github.com/openconfig/functional-translators/ftconsts"
	"github.com/openconfig/functional-translators/ftutilities"
	"github.com/openconfig/functional-translators/translator"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

type fileSystems struct {
	prefixes []string
	free     []uint64
	size     []uint64
	utilized []uint64
}

const convertBytesToMB = 1000000

var (
	translateMap = map[string][]string{
		"/openconfig/system/mount-points/mount-point/state/available": {
			"/Cisco-IOS-XR-shellutil-filesystem-oper/file-system/node/file-system",
		},
		"/openconfig/system/mount-points/mount-point/state/size": {
			"/Cisco-IOS-XR-shellutil-filesystem-oper/file-system/node/file-system",
		},
		"/openconfig/system/mount-points/mount-point/state/name": {
			"/Cisco-IOS-XR-shellutil-filesystem-oper/file-system/node/file-system",
		},
		"/openconfig/system/mount-points/mount-point/state/utilized": {
			"/Cisco-IOS-XR-shellutil-filesystem-oper/file-system/node/file-system",
		},
	}
	paths      = mustStringMapPaths(translateMap)
	nativePath = &gnmipb.Path{
		Origin: "Cisco-IOS-XR-shellutil-filesystem-oper",
		Elem: []*gnmipb.PathElem{
			{Name: "file-system"}, {Name: "node"},
			{Name: "file-system"}, {Name: "*"},
		},
	}
)

func mustStringMapPaths(m map[string][]string) map[string][]*gnmipb.Path {
	p, err := ftutilities.StringMapPaths(m)
	if err != nil {
		log.Fatalf("map %#v cannot parse output paths into gNMI Paths", m)
	}
	return p
}

// validates the leaves and build trap structs
func buildFileSystems(prefix *gnmipb.Path, leaves []*gnmipb.Update) (map[string]*fileSystems, error) {
	nodeFileSystems := make(map[string]*fileSystems)
	// under file-system we have a keyless list for each file-system leaves(size, free, prefixes, etc) first size, belongs to first file-system,
	// first free, belongs to first file-system, etc. There is no otherway to identify a specific file-system except for the order in the list.
	for _, leaf := range leaves {
		path := ftutilities.Join(prefix, leaf.GetPath())
		if !ftutilities.MatchPath(path, nativePath) {
			// This is not a path we are interested in.
			continue
		}
		elems := path.GetElem()
		nodeName := elems[1].GetKey()["node-name"]
		t, ok := nodeFileSystems[nodeName]
		if !ok {
			t = &fileSystems{
				prefixes: make([]string, 0),
				free:     make([]uint64, 0),
				size:     make([]uint64, 0),
				utilized: make([]uint64, 0),
			}
			nodeFileSystems[nodeName] = t
		}
		switch elems[3].GetName() {
		case "size":
			t.size = append(t.size, leaf.GetVal().GetUintVal()/convertBytesToMB)
		case "free":
			t.free = append(t.free, leaf.GetVal().GetUintVal()/convertBytesToMB)
		case "prefixes":
			t.prefixes = append(t.prefixes, leaf.GetVal().GetStringVal())
		default:
			// We do not need other leaves, so we skip them.
			continue
		}
	}
	for node, fs := range nodeFileSystems {
		if len(fs.free) != len(fs.size) || len(fs.size) != len(fs.prefixes) {
			return nil, fmt.Errorf("inconsistent number of leaves under same node %s", node)
		}
		for i := range fs.free {
			if fs.free[i] > fs.size[i] {
				return nil, fmt.Errorf("free %d is greater than size %d", fs.free[i], fs.size[i])
			}
			fs.utilized = append(fs.utilized, fs.size[i]-fs.free[i])
		}
	}
	return nodeFileSystems, nil
}

// New creates a functional translator.
func New() *translator.FunctionalTranslator {
	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               ftconsts.CiscoXRMountTranslator,
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
		log.Fatalf("Failed to create Cisco mount functional translator: %v", err)
	}
	return ft
}

func translate(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
	if sr.GetUpdate() == nil {
		return nil, nil
	}
	nodeFileSystems, err := buildFileSystems(sr.GetUpdate().GetPrefix(), sr.GetUpdate().GetUpdate())
	if err != nil {
		return nil, fmt.Errorf("failed to validate path: %v", err)
	}
	mountRoot := &ocmount.Device{}
	for node, fs := range nodeFileSystems {
		// /file-system/node[node-name]
		for i := range fs.size {
			mountPointName := fmt.Sprintf("%s-%s", node, fs.prefixes[i])
			mountPointState := mountRoot.GetOrCreateSystem().GetOrCreateMountPoints().GetOrCreateMountPoint(mountPointName).GetOrCreateState()
			mountPointState.Available = &fs.free[i]
			mountPointState.Size = &fs.size[i]
			mountPointState.Name = &mountPointName
			mountPointState.Utilized = &fs.utilized[i]
		}
	}
	return ftutilities.FilterStructToState(mountRoot, sr.GetUpdate().GetTimestamp(), "openconfig", sr.GetUpdate().GetPrefix().GetTarget())
}
