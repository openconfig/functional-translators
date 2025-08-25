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

package simplemapper

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"google3/third_party/openconfig/functional_translators/arista/aristainterface/yang/openconfig/interfaces"
	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

func TestBindKeys(t *testing.T) {
	tests := []struct {
		name     string
		pathBind *gnmipb.Path
		path     *gnmipb.Path
		want     map[string]string
		wantErr  bool
	}{
		{
			name: "success",
			pathBind: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "path",
						Key: map[string]string{
							"keyA": "<var1>",
							"keyB": "<var2>",
						},
					},
					{
						Name: "foo",
						Key: map[string]string{
							"keyC": "<var3>",
						},
					},
				},
			},
			path: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "path",
						Key: map[string]string{
							"keyA": "data1",
							"keyB": "data2",
						},
					},
					{
						Name: "foo",
						Key: map[string]string{
							"keyC": "data3",
						},
					},
				},
			},
			want: map[string]string{
				"<var1>": "data1",
				"<var2>": "data2",
				"<var3>": "data3",
			},
		},
		{
			name: "success - constant keys",
			pathBind: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
						Key: map[string]string{
							"keyC": "data3",
						},
					},
					{
						Name: "path",
						Key: map[string]string{
							"keyA": "<var1>",
							"keyB": "<var2>",
						},
					},
				},
			},
			path: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
						Key: map[string]string{
							"keyC": "data3",
						},
					},
					{
						Name: "path",
						Key: map[string]string{
							"keyA": "data1",
							"keyB": "data2",
						},
					},
				},
			},
			want: map[string]string{
				"<var1>": "data1",
				"<var2>": "data2",
			},
		},
		{
			name: "error - duplicate vars",
			pathBind: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "path",
						Key: map[string]string{
							"keyA": "<var1>",
							"keyB": "<var2>",
						},
					},
					{
						Name: "foo",
						Key: map[string]string{
							"keyC": "<var1>",
						},
					},
				},
			},
			path: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "path",
						Key: map[string]string{
							"keyA": "data1",
							"keyB": "data2",
						},
					},
					{
						Name: "foo",
						Key: map[string]string{
							"keyC": "data3",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error - different element names",
			pathBind: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "path",
						Key: map[string]string{
							"keyA": "<var1>",
							"keyB": "<var2>",
						},
					},
				},
			},
			path: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "other",
					},
					{
						Name: "path",
						Key: map[string]string{
							"keyA": "data1",
							"keyB": "data2",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error - different element lengths",
			pathBind: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "path",
						Key: map[string]string{
							"keyA": "<var1>",
							"keyB": "<var2>",
						},
					},
					{
						Name: "leaf",
					},
				},
			},
			path: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "path",
						Key: map[string]string{
							"keyA": "data1",
							"keyB": "data2",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error - key not found",
			pathBind: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "path",
						Key: map[string]string{
							"keyA": "<var1>",
							"keyB": "<var2>",
						},
					},
				},
			},
			path: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "path",
						Key: map[string]string{
							"keyA": "data1",
							"keyC": "data2",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error - different key lengths",
			pathBind: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "path",
						Key: map[string]string{
							"keyA": "<var1>",
							"keyB": "<var2>",
						},
					},
				},
			},
			path: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "path",
						Key: map[string]string{
							"keyA": "data1",
							"keyB": "data2",
							"keyC": "data3",
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := bindKeys(tc.pathBind, tc.path)
			if tc.wantErr && err == nil {
				t.Errorf("bindKeys() returned nil for error, expected error")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("bindKeys() returned an unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("bindKeys() returned an unexpected diff (-want +got): %v", diff)
			}
		})
	}
}

func TestApplyBind(t *testing.T) {
	tests := []struct {
		name     string
		bindings map[string]string
		path     *gnmipb.Path
		want     *gnmipb.Path
		wantErr  bool
	}{
		{
			name: "success",
			bindings: map[string]string{
				"<var1>": "data1",
				"<var2>": "data2",
				"<var3>": "data3",
			},
			path: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "path",
						Key: map[string]string{
							"keyA": "<var1>",
							"keyB": "<var2>",
						},
					},
					{
						Name: "foo",
						Key: map[string]string{
							"keyC": "<var3>",
						},
					},
				},
			},
			want: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "path",
						Key: map[string]string{
							"keyA": "data1",
							"keyB": "data2",
						},
					},
					{
						Name: "foo",
						Key: map[string]string{
							"keyC": "data3",
						},
					},
				},
			},
		},
		{
			name: "error - var not found",
			bindings: map[string]string{
				"<var1>": "data1",
				"<var3>": "data3",
			},
			path: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "path",
						Key: map[string]string{
							"keyA": "<var1>",
							"keyB": "<var2>",
						},
					},
					{
						Name: "foo",
						Key: map[string]string{
							"keyC": "<var3>",
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := applyBind(tc.bindings, tc.path)
			if tc.wantErr && err == nil {
				t.Errorf("applyBind() returned nil for error, expected error")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("applyBind() returned an unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("applyBind() returned an unexpected diff (-want +got): %v", diff)
			}
		})
	}
}

func TestYangValToGNMIVal(t *testing.T) {
	testStr := "string"
	testBool := true
	testFloat := 0.065999276
	tests := []struct {
		name    string
		val     any
		want    *gnmipb.TypedValue
		wantErr bool
	}{
		{
			name: "success - string",
			val:  &testStr,
			want: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "string"}},
		},
		{
			name: "success - bool",
			val:  &testBool,
			want: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_BoolVal{BoolVal: true}},
		},
		{
			name: "success - float",
			val:  &testFloat,
			want: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_DoubleVal{DoubleVal: 0.065999276}},
		},
		{
			name:    "error - unsupported type, empty struct",
			val:     struct{}{},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := yangValToGNMIVal(tc.val)
			if tc.wantErr && err == nil {
				t.Errorf("yangValToGNMIVal(%v) returned nil for error, expected error", tc.val)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("yangValToGNMIVal(%v) returned an unexpected error: %v", tc.val, err)
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("yangValToGNMIVal(%v) returned an unexpected diff (-want +got): %v", tc.val, diff)
			}
		})
	}
}

func TestNewSimpleMapper(t *testing.T) {
	// Only tests that the schema paths for the functional translator are generated and returned through OutputToInputSchemaStrings() correctly.
	tests := []struct {
		name          string
		inSchemaGen   SchemaFn
		outSchemaGen  SchemaFn
		outputToInput map[string]string
		deleteHandler func(*gnmipb.Notification) ([]*gnmipb.Path, error)
		want          map[string][]string
	}{
		{
			name:         "success schema path starts with slash valid origin",
			inSchemaGen:  interfaces.Schema,
			outSchemaGen: interfaces.Schema,
			outputToInput: map[string]string{
				"/openconfig/interfaces/interface[name=<interfaceName>]/state/description": "/openconfig/interfaces/interface[name=<interfaceName>]/config/description",
			},
			deleteHandler: func(*gnmipb.Notification) ([]*gnmipb.Path, error) {
				return nil, nil
			},
			want: map[string][]string{
				"/openconfig/interfaces/interface/state/description": {
					"/openconfig/interfaces/interface/config/description",
				},
			},
		},
		{
			name:         "success schema path starts with valid origin",
			inSchemaGen:  interfaces.Schema,
			outSchemaGen: interfaces.Schema,
			outputToInput: map[string]string{
				"openconfig/interfaces/interface[name=<interfaceName>]/state/description": "openconfig/interfaces/interface[name=<interfaceName>]/config/description",
			},
			deleteHandler: func(*gnmipb.Notification) ([]*gnmipb.Path, error) {
				return nil, nil
			},
			want: map[string][]string{
				"/openconfig/interfaces/interface/state/description": {
					"/openconfig/interfaces/interface/config/description",
				},
			},
		},
		{
			name:         "success schema path without valid origin",
			inSchemaGen:  interfaces.Schema,
			outSchemaGen: interfaces.Schema,
			outputToInput: map[string]string{
				"interfaces/interface[name=<interfaceName>]/state/description": "interfaces/interface[name=<interfaceName>]/config/description",
			},
			deleteHandler: func(*gnmipb.Notification) ([]*gnmipb.Path, error) {
				return nil, nil
			},
			want: map[string][]string{
				"/interfaces/interface/state/description": {
					"/interfaces/interface/config/description",
				},
			},
		},
		{
			name:         "success schema path without valid origin starts with slash",
			inSchemaGen:  interfaces.Schema,
			outSchemaGen: interfaces.Schema,
			outputToInput: map[string]string{
				"/interfaces/interface[name=<interfaceName>]/state/description": "/interfaces/interface[name=<interfaceName>]/config/description",
			},
			deleteHandler: func(*gnmipb.Notification) ([]*gnmipb.Path, error) {
				return nil, nil
			},
			want: map[string][]string{
				"/interfaces/interface/state/description": {
					"/interfaces/interface/config/description",
				},
			},
		},
		{
			name:         "success many input to one output",
			inSchemaGen:  interfaces.Schema,
			outSchemaGen: interfaces.Schema,
			outputToInput: map[string]string{
				"/openconfig/interfaces/interface[name=<lagIntfName>]/ethernet/state/mac-address":      "/openconfig/lacp/interfaces/interface[name=<lagIntfName>]/state/system-id-mac",
				"/openconfig/interfaces/interface[name=<ethernetIntfName>]/ethernet/state/mac-address": "/openconfig/interfaces/interface[name=<ethernetIntfName>]/ethernet/state/mac-address",
			},
			deleteHandler: func(*gnmipb.Notification) ([]*gnmipb.Path, error) {
				return nil, nil
			},
			want: map[string][]string{
				"/openconfig/interfaces/interface/ethernet/state/mac-address": {
					"/openconfig/interfaces/interface/ethernet/state/mac-address",
					"/openconfig/lacp/interfaces/interface/state/system-id-mac",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sm, err := NewSimpleMapper(tc.inSchemaGen, tc.outSchemaGen, tc.outputToInput, tc.deleteHandler)
			if err != nil {
				t.Fatalf("NewSimpleMapper() returned an unexpected error: %v", err)
			}

			got := sm.OutputToInputSchemaStrings()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("OutputToInputSchemaStrings() returned an unexpected diff (-want +got): %v", diff)
			}
		})
	}
}

func TestVarsToWildcards(t *testing.T) {
	tests := []struct {
		desc string
		path *gnmipb.Path
		want *gnmipb.Path
	}{
		{
			desc: "Simple, including constant key",
			path: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
						Key: map[string]string{
							"keyA": "<var1>",
							"keyB": "const",
						},
					},
					{
						Name: "path",
						Key: map[string]string{
							"keyC": "<var3>",
						},
					},
				},
			},
			want: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
						Key: map[string]string{
							"keyA": "*",
							"keyB": "const",
						},
					},
					{
						Name: "path",
						Key: map[string]string{
							"keyC": "*",
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := varsToWildcards(tc.path)
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("varsToWildcards() returned an unexpected diff (-want +got): %v", diff)
			}
		})
	}
}
