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

// Package translator holds base logic for functional translators.
package translator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/openconfig/functional-translators"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

// DeviceMetadata contains metadata to identify a type of device.
type DeviceMetadata struct {
	Vendor          string
	HardwareModel   string
	SoftwareVersion string
}

// FunctionalTranslatorOptions contains the options for a FunctionalTranslator.
type FunctionalTranslatorOptions struct {
	ID               string
	Translate        func(*gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error)
	OutputToInputMap map[string][]*gnmipb.Path
	Metadata         []*DeviceMetadata
	// MatchPaths is a function when given a superset of output paths and device metadata, returns
	// a MatchedPaths which contains the subset of output paths supported by the FT (OutputPaths)
	// and a set of paths (InputPaths) needed to provide those paths as output.
	MatchPaths func(map[string]*gnmipb.Path, *DeviceMetadata) (*MatchedPaths, error)
}

// FunctionalTranslator is a per-platform (vendor/hw_model/sw_model) struct, which handles the
// logic of converting openconfig output path -> []input_paths, as well as any translation
// from []input_path -> output_path(s).
type FunctionalTranslator struct {
	ID               string
	translate        func(*gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error)
	OutputToInputMap map[string][]*gnmipb.Path
	Metadata         []*DeviceMetadata
	matchPaths       func(map[string]*gnmipb.Path, *DeviceMetadata) (*MatchedPaths, error)
}

// NewFunctionalTranslator returns a FunctionalTranslator initialized with provided information.
func NewFunctionalTranslator(opts FunctionalTranslatorOptions) (*FunctionalTranslator, error) {
	// Validate the options.
	if opts.ID == "" {
		return nil, fmt.Errorf("Functional Translator ID is nil")
	}
	if opts.Translate == nil {
		return nil, fmt.Errorf("%s has a nil Translate() function", opts.ID)
	}

	ft := &FunctionalTranslator{
		ID:               opts.ID,
		translate:        opts.Translate,
		OutputToInputMap: opts.OutputToInputMap,
		Metadata:         opts.Metadata,
		matchPaths:       opts.MatchPaths,
	}

	// Apply default values if not provided.
	if ft.matchPaths == nil {
		ft.matchPaths = ft.defaultPathMatcher
	}

	return ft, nil
}

// Translate translates vendor notifications to notifications OpenConfig-compliant notifications.
func (ft *FunctionalTranslator) Translate(input *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
	return ft.translate(input)
}

// MatchPaths is a function when given a superset of output paths and device metadata, returns
// a MatchedPaths which contains the subset of output paths supported by the FT (OutputPaths)
// and a set of paths (InputPaths) needed to provide those paths as output.
func (ft *FunctionalTranslator) MatchPaths(outputSuperset map[string]*gnmipb.Path, deviceMetadata *DeviceMetadata) (*MatchedPaths, error) {
	return ft.matchPaths(outputSuperset, deviceMetadata)
}

// OutputToInput returns a bool indicating if the given output path is supported by the FT, and
// if so, returns the input paths that are needed to provide the output path.
func (ft *FunctionalTranslator) OutputToInput(output *gnmipb.Path) (bool, []*gnmipb.Path, error) {
	outputKey := ftutilities.GNMIPathToSchemaString(output, false)
	inputs, ok := ft.OutputToInputMap[outputKey]
	return ok, inputs, nil
}

// MatchedPaths include the set of paths supported by the FT (OutputPaths) and the set of required
// paths needed to be able to provide those paths as output.
type MatchedPaths struct {
	InputPaths []*gnmipb.Path
	// OutputToInput contains the information describing which input paths are needed to provide
	// which output paths. It is a map from output path to a list of input paths.
	OutputToInput map[string][]string
}

func (ft *FunctionalTranslator) metadataMatch(got *DeviceMetadata) bool {
	if len(ft.Metadata) == 0 {
		return true
	}
	for _, m := range ft.Metadata {
		if m.Vendor != "" && !strings.EqualFold(m.Vendor, got.Vendor) {
			continue
		}
		if m.HardwareModel != "" && !strings.EqualFold(m.HardwareModel, got.HardwareModel) {
			continue
		}
		if m.SoftwareVersion != "" && !strings.EqualFold(m.SoftwareVersion, got.SoftwareVersion) {
			continue
		}
		return true
	}
	return false
}

// defaultPathMatcher returns the required set of telemetry paths that should be a part of a subscription
// to the device in order to provide the set out requested output paths. Returned list is not
// guaranteed to be free of duplicates. The keys of the paths should be the result of
// ygot.PathToString.
func (ft *FunctionalTranslator) defaultPathMatcher(outputSuperset map[string]*gnmipb.Path, deviceMetadata *DeviceMetadata) (*MatchedPaths, error) {
	if deviceMetadata == nil {
		return nil, fmt.Errorf("deviceMetadata cannot be nil, got %v", deviceMetadata)
	}
	if !ft.metadataMatch(deviceMetadata) {
		return nil, nil
	}
	var returnInputPaths []*gnmipb.Path
	returnOutputToInput := map[string][]string{}
	// Most often, we expect len(desiredOutputPaths) > len(ft.outputToInput); so, we iterate through
	// the shorter list.
	for key, inputs := range ft.OutputToInputMap {
		if outputPath, ok := outputSuperset[key]; ok {
			// Check that the key matches the path.
			outputKey := ftutilities.GNMIPathToSchemaString(outputPath, false)
			if outputKey != key {
				return nil, fmt.Errorf("GNMIPathPathToString(path) = %s does not match desired output path %s", outputKey, key)
			}
			returnInputPaths = append(returnInputPaths, inputs...)
			inputKeys := make([]string, 0, len(inputs))
			for _, input := range inputs {
				inputKeys = append(inputKeys, ftutilities.GNMIPathToSchemaString(input, false))
			}
			returnOutputToInput[key] = inputKeys
		}
	}
	// Sort for consistent return ordering.
	sort.Slice(returnInputPaths, ftutilities.SortByYgotString(returnInputPaths))
	return &MatchedPaths{
		InputPaths:    returnInputPaths,
		OutputToInput: returnOutputToInput,
	}, nil
}
