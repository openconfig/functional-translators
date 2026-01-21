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

package translator

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
	"github.com/openconfig/functional-translators/ftutilities"
	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

func TestMetadataMatch(t *testing.T) {
	tests := []struct {
		name       string
		ftMetaData []*FTMetadata
		wantMatch  bool
	}{
		{
			name:       "empty requirement catching all",
			ftMetaData: nil,
			wantMatch:  true,
		},
		{
			name: "partial requirement",
			ftMetaData: []*FTMetadata{
				{
					Vendor: "vendor",
				},
			},
			wantMatch: true,
		},
		{
			name: "requirement exact match",
			ftMetaData: []*FTMetadata{
				{
					Vendor:          "vendor",
					SoftwareVersion: "sw",
					HardwareModel:   "hw",
				},
			},
			wantMatch: true,
		},
		{
			name: "requirement case insensitive match",
			ftMetaData: []*FTMetadata{
				{
					Vendor:          "Vendor",
					SoftwareVersion: "SW",
					HardwareModel:   "HW",
				},
			},
			wantMatch: true,
		},
		{
			name: "match one",
			ftMetaData: []*FTMetadata{
				{
					Vendor: "non-vendor",
				},
				{
					Vendor: "vendor",
				},
			},
			wantMatch: true,
		},
		{
			name: "match all",
			ftMetaData: []*FTMetadata{
				{
					Vendor: "vendor",
				},
				{
					SoftwareVersion: "sw",
				},
			},
			wantMatch: true,
		},
		{
			name: "match none",
			ftMetaData: []*FTMetadata{
				{
					Vendor: "non-vendor",
				},
				{
					SoftwareVersion: "non-sw",
				},
				{
					HardwareModel: "non-hw",
				},
			},
			wantMatch: false,
		},
		{
			name: "vendor mismatch",
			ftMetaData: []*FTMetadata{
				{
					Vendor: "non-vendor",
				},
			},
			wantMatch: false,
		},
		{
			name: "sw mismatch",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersion: "non-sw",
				},
			},
			wantMatch: false,
		},
		{
			name: "hw mismatch",
			ftMetaData: []*FTMetadata{
				{
					HardwareModel: "non-hw",
				},
			},
			wantMatch: false,
		},
	}

	inputMetaData := &DeviceMetadata{
		Vendor:          "vendor",
		SoftwareVersion: "sw",
		HardwareModel:   "hw",
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ft, err := NewFunctionalTranslator(FunctionalTranslatorOptions{
				ID: "test",
				Translate: func(*gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
					return nil, nil
				},
				Metadata: tc.ftMetaData,
			})
			if err != nil {
				t.Fatalf("NewFunctionalTranslator got error: %v, want error: nil", err)
			}
			got := ft.metadataMatch(inputMetaData)
			if got != tc.wantMatch {
				t.Errorf("metadataMatch got %v, want %v", got, tc.wantMatch)
			}
		})
	}
}

func TestMetadataMatchSWRange(t *testing.T) {
	tests := []struct {
		name       string
		swVersion  string
		ftMetaData []*FTMetadata
		wantMatch  bool
	}{
		{
			name:      "sw range match",
			swVersion: "1.5",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.0",
						ExclusiveMax: "2.0",
					},
				},
			},
			wantMatch: true,
		},
		{
			name:      "sw range match, sw version longer than max",
			swVersion: "2.0.1",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.0",
						ExclusiveMax: "2.0",
					},
				},
			},
			wantMatch: false,
		},
		{
			name:      "sw range match, numeric suffix",
			swVersion: "1.5.1",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.0",
						ExclusiveMax: "2.0",
					},
				},
			},
			wantMatch: true,
		},
		{
			name:      "sw range match, letter suffix",
			swVersion: "4.34.2F-(some-random-suffix)",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "4.34",
						ExclusiveMax: "4.35",
					},
				},
			},
			wantMatch: true,
		},
		{
			name:      "sw range mismatch",
			swVersion: "2.5",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.0",
						ExclusiveMax: "2.0",
					},
				},
			},
			wantMatch: false,
		},

		{
			name:      "sw range mismatch, numeric suffix",
			swVersion: "2.5.1",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.0",
						ExclusiveMax: "2.0",
					},
				},
			},
			wantMatch: false,
		},
		{
			name:      "sw range match, letter suffix",
			swVersion: "4.34.2F",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "4.34",
						ExclusiveMax: "4.35",
					},
				},
			},
			wantMatch: true,
		},
		{
			name:      "sw range match, equal min",
			swVersion: "1.0",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.0",
						ExclusiveMax: "2.0",
					},
				},
			},
			wantMatch: true,
		},
		{
			name:      "sw range mismatch, equal max",
			swVersion: "2.0",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.0",
						ExclusiveMax: "2.0",
					},
				},
			},
			wantMatch: false,
		},
		{
			name:      "sw range match, letter comparison",
			swVersion: "1.5B",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.5A",
						ExclusiveMax: "1.5C",
					},
				},
			},
			wantMatch: true,
		},
		{
			name:      "sw range mismatch, after, letter comparison",
			swVersion: "1.5D",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.5A",
						ExclusiveMax: "1.5C",
					},
				},
			},
			wantMatch: false,
		},
		{
			name:      "sw range mismatch, before, letter comparison",
			swVersion: "1.5A",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.5B",
						ExclusiveMax: "1.5C",
					},
				},
			},
			wantMatch: false,
		},
		{
			name:      "sw range match, letter is greater than number",
			swVersion: "1.A",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.9",
						ExclusiveMax: "1.B",
					},
				},
			},
			wantMatch: true,
		},
		{
			name:      "sw outside range, letter is greater than number",
			swVersion: "1.B",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.8",
						ExclusiveMax: "1.9",
					},
				},
			},
			wantMatch: false,
		},
		{
			name:      "sw range mismatch, less than min",
			swVersion: "0.9",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.0",
						ExclusiveMax: "2.0",
					},
				},
			},
			wantMatch: false,
		},
		{
			name:      "longer version in range, inside range",
			swVersion: "1.3",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.2.3.4.5",
						ExclusiveMax: "2.1.0",
					},
				},
			},
			wantMatch: true,
		},
		{
			name:      "longer version in range, outside range",
			swVersion: "1.3",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.2.3.4.5",
						ExclusiveMax: "1.3.0",
					},
				},
			},
			wantMatch: false,
		},
		{
			name:      "sw range match, identical letter components",
			swVersion: "1.A.2",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.A.1",
						ExclusiveMax: "1.A.3",
					},
				},
			},
			wantMatch: true,
		},
		{
			name:      "sw version is prefix of min version",
			swVersion: "1.2",
			ftMetaData: []*FTMetadata{
				{
					SoftwareVersionRange: &SWRange{
						InclusiveMin: "1.2.3",
						ExclusiveMax: "2.0",
					},
				},
			},
			wantMatch: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ft, err := NewFunctionalTranslator(FunctionalTranslatorOptions{
				ID: "test",
				Translate: func(*gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
					return nil, nil
				},
				Metadata: tc.ftMetaData,
			})
			if err != nil {
				t.Fatalf("NewFunctionalTranslator got error: %v, want error: nil", err)
			}
			got := ft.metadataMatch(&DeviceMetadata{SoftwareVersion: tc.swVersion})
			if got != tc.wantMatch {
				t.Errorf("swVersionMatch got %v, want %v for swVersion %s against ft: %v", got, tc.wantMatch, tc.swVersion, ft.Metadata())
			}
		})
	}
}

// TestSensibility tests that the FunctionalTranslator reasonably initializes and handles FTs to
// avoid panics.
func TestSensibility(t *testing.T) {
	t.Run("simple-ft-new", func(t *testing.T) {
		ft, err := NewFunctionalTranslator(FunctionalTranslatorOptions{
			ID: "test",
			Translate: func(*gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
				return nil, nil
			},
			Metadata: []*FTMetadata{
				{
					Vendor: "test",
				},
			},
		})
		if err != nil {
			t.Fatalf("NewFunctionalTranslator got unexpected error %v", err)
		}
		if ft == nil {
			t.Fatalf("NewFunctionalTranslator got nil, want non-nil")
		}
	})

	t.Run("nil-translate-new", func(t *testing.T) {
		_, err := NewFunctionalTranslator(FunctionalTranslatorOptions{
			ID: "test",
			Metadata: []*FTMetadata{
				{
					Vendor: "test",
				},
			},
		})
		if err == nil {
			t.Errorf("NewFunctionalTranslator got no error, want error")
		}
	})

	t.Run("nil-output-to-input", func(t *testing.T) {
		nilOutputToInputFT, err := NewFunctionalTranslator(FunctionalTranslatorOptions{
			ID: "test",
			Translate: func(*gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
				return nil, nil
			},
			Metadata: []*FTMetadata{
				{
					Vendor: "test",
				},
			},
		})
		if err != nil {
			t.Fatalf("NewFunctionalTranslator got error: %v, want error: nil", err)
		}
		if _, _, err := nilOutputToInputFT.OutputToInput(nil); err == nil {
			t.Errorf("nilOutputToInputFT.OutputToInput(nil) got no error, want error")
		}
	})
}

func TestNewFunctionalTranslator(t *testing.T) {
	ftMetadata := FTMetadata{Vendor: "vendor"}
	tests := []struct {
		name      string
		opts      FunctionalTranslatorOptions
		stringMap map[string][]string
		wantErr   bool
	}{
		{
			name: "valid_options",
			opts: FunctionalTranslatorOptions{
				ID: "test-id",
				Translate: func(*gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
					return nil, nil
				},
				Metadata: []*FTMetadata{&ftMetadata},
			},
		},
		{
			name: "valid_options_with_OutputToInputMap",
			opts: FunctionalTranslatorOptions{
				ID: "test-id-with-map",
				Translate: func(*gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
					return nil, nil
				},
				Metadata: []*FTMetadata{&ftMetadata},
				OutputToInputMap: ftutilities.MustStringMapPaths(map[string][]string{
					"/openconfig/components/component/properties/property/state/value": {
						"/Cisco-IOS-XR-show-fpd-loc-ng-oper/show-fpd/hw-module-fpd",
					},
				}),
			},
			stringMap: map[string][]string{
				"/openconfig/components/component/properties/property/state/value": {
					"/Cisco-IOS-XR-show-fpd-loc-ng-oper/show-fpd/hw-module-fpd",
				},
			},
		},
		{
			name: "no_slash_output_OutputToInputMap",
			opts: FunctionalTranslatorOptions{
				ID: "test-id-with-map",
				Translate: func(*gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
					return nil, nil
				},
				Metadata: []*FTMetadata{&ftMetadata},
				OutputToInputMap: ftutilities.MustStringMapPaths(map[string][]string{
					"openconfig/components/component/properties/property/state/value": {
						"/Cisco-IOS-XR-show-fpd-loc-ng-oper/show-fpd/hw-module-fpd",
					},
				}),
			},
			wantErr: true,
		},
		{
			name: "incorrect_origin_output_OutputToInputMap",
			opts: FunctionalTranslatorOptions{
				ID: "test-id-with-map",
				Translate: func(*gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
					return nil, nil
				},
				Metadata: []*FTMetadata{&ftMetadata},
				OutputToInputMap: ftutilities.MustStringMapPaths(map[string][]string{
					"/open/components/component/properties/property/state/value": {
						"/Cisco-IOS-XR-show-fpd-loc-ng-oper/show-fpd/hw-module-fpd",
					},
				}),
			},
			wantErr: true,
		},
		{
			name: "incorrect_origin_input_OutputToInputMap",
			opts: FunctionalTranslatorOptions{
				ID: "test-id-with-map",
				Translate: func(*gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
					return nil, nil
				},
				Metadata: []*FTMetadata{&ftMetadata},
				OutputToInputMap: ftutilities.MustStringMapPaths(map[string][]string{
					"/openconfig/components/component/properties/property/state/value": {
						"/Cisco-XR-show-fpd-loc-ng-oper/show-fpd/hw-module-fpd",
					},
				}),
			},
			wantErr: true,
		},
		{
			name: "missing_ID",
			opts: FunctionalTranslatorOptions{
				Translate: func(*gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
					return nil, nil
				},
				Metadata: []*FTMetadata{&ftMetadata},
			},
			wantErr: true,
		},
		{
			name: "nil_Translate",
			opts: FunctionalTranslatorOptions{
				ID:       "test-id",
				Metadata: []*FTMetadata{&ftMetadata},
			},
			wantErr: true,
		},
		{
			name: "empty_Metadata_is_valid",
			opts: FunctionalTranslatorOptions{
				ID: "test-id",
				Translate: func(*gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
					return nil, nil
				},
				Metadata: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ft, err := NewFunctionalTranslator(tc.opts)
			if (err != nil) != tc.wantErr {
				t.Errorf("NewFunctionalTranslator(%v) got error: %v, want error: %v", tc.opts, err, tc.wantErr)
			}
			if !tc.wantErr && ft == nil {
				t.Errorf("NewFunctionalTranslator(%v) got nil ft, want non-nil", tc.opts)
			}
			if !tc.wantErr && ft.ID() != tc.opts.ID {
				t.Errorf("NewFunctionalTranslator(%v) got ID %q, want %q", tc.opts, ft.ID(), tc.opts.ID)
			}
			if !tc.wantErr {
				outputSuperset := make(map[string]*gnmipb.Path)
				for output := range tc.opts.OutputToInputMap {
					p, err := ftutilities.StringToPath(output)
					if err != nil {
						t.Fatalf("Failed to parse path %q: %v", output, err)
					}
					outputSuperset[output] = p
				}
				deviceMetadata := &DeviceMetadata{
					Vendor: ftMetadata.Vendor,
				}
				matched, err := ft.MatchPaths(outputSuperset, deviceMetadata)
				if err != nil {
					t.Fatalf("ft.MatchPaths(%v, %v) got unexpected error: %v", outputSuperset, deviceMetadata, err)
				}
				if diff := cmp.Diff(tc.stringMap, matched.OutputToInput, cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("ft.MatchPoints(%v, %v) returned diff (-want +got):\n%s", outputSuperset, deviceMetadata, diff)
				}
			}
		})
	}
}

func TestOutputToInput(t *testing.T) {
	outputPath := "/openconfig/interfaces/interface/state/counters/in-octets"
	translateMap := map[string][]string{
		outputPath: {
			"/openconfig/interfaces/interface/state/counters/in-octets",
		},
	}
	stringPaths := ftutilities.MustStringMapPaths(translateMap)
	ft, err := NewFunctionalTranslator(FunctionalTranslatorOptions{
		ID:               "test-ft",
		Translate:        func(*gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) { return nil, nil },
		Metadata:         []*FTMetadata{{Vendor: "vendor"}},
		OutputToInputMap: stringPaths,
	})
	if err != nil {
		t.Fatalf("NewFunctionalTranslator got unexpected error: %v", err)
	}

	tests := []struct {
		name string
		path string
		want []*gnmipb.Path
		ok   bool
	}{
		{
			name: "path-found",
			path: outputPath,
			want: stringPaths[outputPath],
			ok:   true,
		},
		{
			name: "path-not-found",
			path: "/openconfig/interfaces/interface/state/counters/",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path, err := ftutilities.StringToPath(tc.path)
			if err != nil {
				t.Fatalf("ftutilities.StringToPath(%s) got unexpected error: %v", outputPath, err)
			}
			ok, got, err := ft.OutputToInput(path)
			if err != nil {
				t.Fatalf("OutputToInput(%v) got error: %v, want error: nil", tc.path, err)
			}
			if ok != tc.ok {
				t.Fatalf("OutputToInput(%v) ok got: %v, want: %v", tc.path, ok, tc.ok)
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform(), cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("OutputToInput(%v) returned diff (-want +got):\n%s", tc.path, diff)
			}
		})
	}
}
