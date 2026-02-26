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

package aristacfmpm

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"github.com/openconfig/functional-translators/ftutilities"
	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

func TestTranslate(t *testing.T) {
	tests := []struct {
		name           string
		inputPath      string
		wantOutputPath string
		wantNil        bool
		wantErr        bool
	}{
		{
			name:           "CFMPM translation success",
			inputPath:      "testdata/pm_input.txt",
			wantOutputPath: "testdata/pm_output.txt",
		},
		{
			name:           "CFMPM delete success",
			inputPath:      "testdata/delete_input.txt",
			wantOutputPath: "testdata/delete_output.txt",
		},
		{
			name:      "Empty profile name -> no output",
			inputPath: "testdata/empty_profile_input.txt",
			wantNil:   true,
		},
		{
			name:           "JSON profile name support",
			inputPath:      "testdata/json_profile_input.txt",
			wantOutputPath: "testdata/json_profile_output.txt",
		},
		{
			name:           "JSON wrapper profile name support",
			inputPath:      "testdata/json_wrapper_profile_input.txt",
			wantOutputPath: "testdata/json_wrapper_profile_output.txt",
		},
		{
			name:           "Loss Measurement translation success",
			inputPath:      "testdata/loss_input.txt",
			wantOutputPath: "testdata/loss_output.txt",
		},
		{
			name:           "SLM stats support",
			inputPath:      "testdata/slm_input.txt",
			wantOutputPath: "testdata/slm_output.txt",
		},
		{
			name:      "SLM stats malformed key -> no output",
			inputPath: "testdata/slm_malformed_key_input.txt",
			wantErr:   true,
		},
		{
			name:      "SLM stats invalid/empty MEP ID -> no output",
			inputPath: "testdata/slm_invalid_mep_id_input.txt",
			wantErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inputSR, err := ftutilities.LoadSubscribeResponse(test.inputPath)
			if err != nil {
				t.Fatalf("Failed to load input message: %v", err)
			}
			ft := New()
			gotSR, err := ft.Translate(inputSR)
			if gotNil, gotErr := gotSR == nil, err != nil; gotNil || gotErr {
				switch {
				case gotErr != test.wantErr:
					t.Fatalf("Unexpected error result returned from translate() = %v, want error %t", err, test.wantErr)
				case err != nil:
					return
				case gotNil != test.wantNil:
					t.Fatalf("Unexpected nil result returned from Translate() = %t, want nil %t", gotNil, test.wantNil)
				default:
					return
				}
			}

			if test.wantNil || test.wantErr {
				return
			}

			wantSR, err := ftutilities.LoadSubscribeResponse(test.wantOutputPath)
			if err != nil {
				t.Fatalf("Failed to load want message: %v", err)
			}
			if diff := cmp.Diff(wantSR, gotSR, protocmp.Transform()); diff != "" {
				t.Fatalf("Unexpected diff from translate (-want +got): %s", diff)
			}
		})
	}
}

func TestValueToString(t *testing.T) {
	tests := []struct {
		name  string
		input *gnmipb.TypedValue
		want  string
	}{
		{
			name:  "nil input",
			input: nil,
			want:  "",
		},
		{
			name:  "int value",
			input: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_IntVal{IntVal: 10}},
			want:  "",
		},
		{
			name:  "json unmarshal error",
			input: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_JsonVal{JsonVal: []byte(`{invalid}`)}},
			want:  "{invalid}",
		},
		{
			name:  "json string without wrapper",
			input: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_JsonVal{JsonVal: []byte(`invalid json but valid string`)}},
			want:  "invalid json but valid string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := valueToString(tt.input); got != tt.want {
				t.Errorf("valueToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseIntSlice(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "invalid int",
			input:   "110 invalid",
			wantErr: true,
		},
		{
			name:  "null terminator",
			input: "110 0 111",
			want:  "n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseIntSlice(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseIntSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseIntSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseSmashKey(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "missing assocID start",
			input:   "invalid key format",
			wantErr: true,
		},
		{
			name:    "missing assocID end",
			input:   "id_Array{base: 0, slice: [110",
			wantErr: true,
		},
		{
			name:    "missing domainID array format",
			input:   "id_Array{base: 0, slice: [110]}_maNameFormatShortInt_<id>_",
			wantErr: true,
		},
		{
			name:    "missing localMEPID marker",
			input:   "id_Array{base: 0, slice: [110]}_maNameFormatShortInt_<id>_Array{base: 0, slice: [111]}",
			wantErr: true,
		},
		{
			name:    "missing localMEPID value",
			input:   "id_Array{base: 0, slice: [110]}_maNameFormatShortInt_<id>_Array{base: 0, slice: [111]}_mdNameFormatNoName_",
			wantErr: true,
		},
		{
			name:    "invalid assoc array parsing",
			input:   "id_Array{base: 0, slice: [invalid]}_maNameFormatShortInt_<id>_Array{base: 0, slice: [111]}_mdNameFormatNoName_1",
			wantErr: true,
		},
		{
			name:    "invalid domain array parsing",
			input:   "id_Array{base: 0, slice: [110]}_maNameFormatShortInt_<id>_Array{base: 0, slice: [invalid]}_mdNameFormatNoName_1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, err := parseSmashKey(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSmashKey() expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestConvertMetricValue(t *testing.T) {
	i := &impl{}
	tests := []struct {
		name       string
		tv         *gnmipb.TypedValue
		metricType string
		wantErr    bool
	}{
		{
			name:    "nil value",
			tv:      nil,
			wantErr: true,
		},
		{
			name:       "unsupported typed value",
			tv:         &gnmipb.TypedValue{Value: &gnmipb.TypedValue_BoolVal{BoolVal: true}},
			metricType: metricDelay,
			wantErr:    true,
		},
		{
			name:       "invalid string metric",
			tv:         &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "abc"}},
			metricType: metricDelay,
			wantErr:    true,
		},
		{
			name:       "unknown metric type",
			tv:         &gnmipb.TypedValue{Value: &gnmipb.TypedValue_IntVal{IntVal: 10}},
			metricType: "unknown",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := i.convertMetricValue(tt.tv, tt.metricType)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertMetricValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
