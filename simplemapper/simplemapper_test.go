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
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"github.com/openconfig/ygot/ytypes"
	"github.com/openconfig/functional-translators/arista/aristainterface/yang/openconfig"
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
			name:     "success - constant key",
			bindings: map[string]string{"<var1>": "data1"},
			path: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "path",
						Key: map[string]string{
							"keyA": "<var1>",
							"keyB": "const",
						},
					},
				},
			},
			want: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "path",
						Key: map[string]string{
							"keyA": "data1",
							"keyB": "const",
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
		{
			name:    "error - nil value",
			val:     nil,
			wantErr: true,
		},
		{
			name:    "error - nil string value",
			val:     (*string)(nil),
			wantErr: true,
		},
		{
			name:    "error - nil bool value",
			val:     (*bool)(nil),
			wantErr: true,
		},
		{
			name:    "error - nil float64 value",
			val:     (*float64)(nil),
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
			inSchemaGen:  openconfig.Schema,
			outSchemaGen: openconfig.Schema,
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
			inSchemaGen:  openconfig.Schema,
			outSchemaGen: openconfig.Schema,
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
			inSchemaGen:  openconfig.Schema,
			outSchemaGen: openconfig.Schema,
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
			inSchemaGen:  openconfig.Schema,
			outSchemaGen: openconfig.Schema,
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
			inSchemaGen:  openconfig.Schema,
			outSchemaGen: openconfig.Schema,
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

func TestUpdateHandler(t *testing.T) {
	isc, err := openconfig.Schema()
	if err != nil {
		t.Fatalf("Failed to load input schema: %v", err)
	}
	osc, err := openconfig.Schema()
	if err != nil {
		t.Fatalf("Failed to load output schema: %v", err)
	}
	m, err := NewSimpleMapper(openconfig.Schema, openconfig.Schema, map[string]string{
		"/openconfig/interfaces/interface[name=<name>]/config/description": "/openconfig/interfaces/interface[name=<name>]/config/description",
	}, func(*gnmipb.Notification) ([]*gnmipb.Path, error) { return nil, nil })
	if err != nil {
		t.Fatalf("NewSimpleMapper() failed: %v", err)
	}

	notification := &gnmipb.Notification{
		Update: []*gnmipb.Update{
			{
				Path: &gnmipb.Path{
					Elem: []*gnmipb.PathElem{
						{Name: "interfaces"},
						{Name: "interface", Key: map[string]string{"name": "eth0"}},
						{Name: "config"},
						{Name: "description"},
					},
				},
				Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "foo"}},
			},
		},
	}

	got, err := m.updateHandler(isc, osc, notification)
	if err != nil {
		t.Fatalf("updateHandler() failed: %v", err)
	}

	if got == nil {
		t.Fatal("updateHandler() returned nil")
	}

	if got.Prefix == nil {
		t.Errorf("updateHandler() returned notification with nil prefix")
	}
	t.Run("no updates match", func(t *testing.T) {
		m, err := NewSimpleMapper(openconfig.Schema, openconfig.Schema, map[string]string{
			"/openconfig/interfaces/interface[name=<name>]/config/description": "/openconfig/interfaces/interface[name=<name>]/config/description",
		}, nil)
		if err != nil {
			t.Fatalf("NewSimpleMapper() failed: %v", err)
		}
		notification := &gnmipb.Notification{
			Update: []*gnmipb.Update{{
				Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "non-existent"}}},
				Val:  &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "foo"}},
			}},
		}
		// We need to use a valid path for unmarshal but not mapped.
		notification.Update[0].Path = &gnmipb.Path{Elem: []*gnmipb.PathElem{
			{Name: "interfaces"},
			{Name: "interface", Key: map[string]string{"name": "eth0"}},
			{Name: "config"},
			{Name: "name"},
		}}
		got, err := m.updateHandler(isc, osc, notification)
		if err != nil {
			t.Fatalf("updateHandler() failed: %v", err)
		}
		if got != nil {
			t.Errorf("updateHandler() = %v, want nil", got)
		}
	})
}

func TestHandler(t *testing.T) {
	m, err := NewSimpleMapper(openconfig.Schema, openconfig.Schema, map[string]string{
		"/openconfig/interfaces/interface[name=<name>]/config/description": "/openconfig/interfaces/interface[name=<name>]/config/description",
	}, func(*gnmipb.Notification) ([]*gnmipb.Path, error) { return nil, nil })
	if err != nil {
		t.Fatalf("NewSimpleMapper() failed: %v", err)
	}

	tests := []struct {
		name string
		in   *gnmipb.SubscribeResponse
		want *gnmipb.Notification
	}{
		{
			name: "single update",
			in: &gnmipb.SubscribeResponse{
				Response: &gnmipb.SubscribeResponse_Update{
					Update: &gnmipb.Notification{
						Update: []*gnmipb.Update{{
							Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{
								{Name: "interfaces"},
								{Name: "interface", Key: map[string]string{"name": "eth0"}},
								{Name: "config"},
								{Name: "description"},
							}},
							Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "foo"}},
						}},
					},
				},
			},
			want: &gnmipb.Notification{
				Prefix: &gnmipb.Path{Origin: "openconfig"},
				Update: []*gnmipb.Update{{
					Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{
						{Name: "interfaces"},
						{Name: "interface", Key: map[string]string{"name": "eth0"}},
						{Name: "config"},
						{Name: "description"},
					}},
					Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotSR, err := m.Handler(tc.in)
			if err != nil {
				t.Fatalf("Handler() failed: %v", err)
			}
			got := gotSR.GetUpdate()
			if diff := cmp.Diff(tc.want, got, protocmp.Transform(), protocmp.IgnoreFields(&gnmipb.Notification{}, "timestamp")); diff != "" {
				t.Errorf("Handler() returned unexpected diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandler_Deletes(t *testing.T) {
	m, err := NewSimpleMapper(openconfig.Schema, openconfig.Schema, map[string]string{
		"/openconfig/interfaces/interface[name=<name>]/config/description": "/openconfig/interfaces/interface[name=<name>]/config/description",
	}, func(n *gnmipb.Notification) ([]*gnmipb.Path, error) {
		return []*gnmipb.Path{{
			Elem: []*gnmipb.PathElem{
				{Name: "interfaces"},
				{Name: "interface", Key: map[string]string{"name": "eth0"}},
				{Name: "config"},
				{Name: "description"},
			},
		}}, nil
	})
	if err != nil {
		t.Fatalf("NewSimpleMapper() failed: %v", err)
	}

	in := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{},
		},
	}

	want := &gnmipb.Notification{
		Prefix: &gnmipb.Path{Origin: "openconfig"},
		Delete: []*gnmipb.Path{{
			Elem: []*gnmipb.PathElem{
				{Name: "interfaces"},
				{Name: "interface", Key: map[string]string{"name": "eth0"}},
				{Name: "config"},
				{Name: "description"},
			},
		}},
	}

	gotSR, err := m.Handler(in)
	if err != nil {
		t.Fatalf("Handler() failed: %v", err)
	}
	got := gotSR.GetUpdate()
	if diff := cmp.Diff(want, got, protocmp.Transform(), protocmp.IgnoreFields(&gnmipb.Notification{}, "timestamp")); diff != "" {
		t.Errorf("Handler() returned unexpected diff (-want +got):\n%s", diff)
	}
}

func TestHandler_Filtering(t *testing.T) {
	m, err := NewSimpleMapper(openconfig.Schema, openconfig.Schema, map[string]string{
		"/openconfig/interfaces/interface[name=<name>]/config/description": "/openconfig/interfaces/interface[name=<name>]/config/description",
	}, func(n *gnmipb.Notification) ([]*gnmipb.Path, error) {
		return []*gnmipb.Path{
			{ // Kept
				Elem: []*gnmipb.PathElem{
					{Name: "interfaces"},
					{Name: "interface", Key: map[string]string{"name": "eth0"}},
					{Name: "config"},
				},
			},
			{ // Filtered
				Elem: []*gnmipb.PathElem{
					{Name: "system"},
					{Name: "config"},
				},
			},
		}, nil
	})
	if err != nil {
		t.Fatalf("NewSimpleMapper() failed: %v", err)
	}

	in := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Update: []*gnmipb.Update{
					{ // Filtered because not in mappings
						Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{
							{Name: "interfaces"},
							{Name: "interface", Key: map[string]string{"name": "eth0"}},
							{Name: "config"},
							{Name: "name"},
						}},
						Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "eth0"}},
					},
				},
			},
		},
	}

	want := &gnmipb.Notification{
		Prefix: &gnmipb.Path{Origin: "openconfig"},
		Delete: []*gnmipb.Path{{
			Elem: []*gnmipb.PathElem{
				{Name: "interfaces"},
				{Name: "interface", Key: map[string]string{"name": "eth0"}},
				{Name: "config"},
			},
		}},
	}

	gotSR, err := m.Handler(in)
	if err != nil {
		t.Fatalf("Handler() failed: %v", err)
	}
	got := gotSR.GetUpdate()
	if diff := cmp.Diff(want, got, protocmp.Transform(), protocmp.IgnoreFields(&gnmipb.Notification{}, "timestamp")); diff != "" {
		t.Errorf("Handler() returned unexpected diff (-want +got):\n%s", diff)
	}
}

func TestHandler_EmptyInput(t *testing.T) {
	m, _ := NewSimpleMapper(openconfig.Schema, openconfig.Schema, nil, nil)
	got, err := m.Handler(&gnmipb.SubscribeResponse{})
	if err != nil {
		t.Errorf("Handler() failed: %v", err)
	}
	if got != nil {
		t.Errorf("Handler() = %v, want nil", got)
	}
}

func TestHandler_NoOutput(t *testing.T) {
	m, _ := NewSimpleMapper(openconfig.Schema, openconfig.Schema, map[string]string{
		"/openconfig/interfaces/interface[name=<name>]/config/description": "/openconfig/interfaces/interface[name=<name>]/config/description",
	}, func(*gnmipb.Notification) ([]*gnmipb.Path, error) { return nil, nil })

	in := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Update: []*gnmipb.Update{{
					Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{
						{Name: "interfaces"},
						{Name: "interface", Key: map[string]string{"name": "eth0"}},
						{Name: "config"},
						{Name: "name"},
					}},
					Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "eth0"}},
				}},
			},
		},
	}
	got, err := m.Handler(in)
	if err != nil {
		t.Errorf("Handler() failed: %v", err)
	}
	if got != nil {
		t.Errorf("Handler() = %v, want nil", got)
	}
}

func TestHandler_UpdateError(t *testing.T) {
	m, _ := NewSimpleMapper(openconfig.Schema, openconfig.Schema, map[string]string{
		"/openconfig/interfaces/interface[name=<name>]/config/description": "/openconfig/interfaces/interface[name=<name>]/config/description",
	}, func(*gnmipb.Notification) ([]*gnmipb.Path, error) { return nil, nil })

	// Trigger updateHandler error (unmarshal failure - wrong type)
	in := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{
				Update: []*gnmipb.Update{{
					Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{
						{Name: "interfaces"},
						{Name: "interface", Key: map[string]string{"name": "eth0"}},
						{Name: "config"},
						{Name: "enabled"},
					}},
					Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "not-a-bool"}},
				}},
			},
		},
	}
	_, err := m.Handler(in)
	if err == nil {
		t.Errorf("Handler() succeeded for bad input, want error")
	}
}

func TestHandler_DeleteError(t *testing.T) {
	// Trigger deleteHandler error
	m, _ := NewSimpleMapper(openconfig.Schema, openconfig.Schema, nil, func(*gnmipb.Notification) ([]*gnmipb.Path, error) {
		return nil, fmt.Errorf("delete error")
	})
	in := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{},
		},
	}
	_, err := m.Handler(in)
	if err == nil {
		t.Errorf("Handler() succeeded for delete error, want error")
	}
}

func TestHandler_NilDeletes(t *testing.T) {
	m, _ := NewSimpleMapper(openconfig.Schema, openconfig.Schema, map[string]string{
		"/openconfig/interfaces/interface[name=<name>]/config/description": "/openconfig/interfaces/interface[name=<name>]/config/description",
	}, nil)

	in := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: &gnmipb.Notification{},
		},
	}
	got, err := m.Handler(in)
	if err != nil {
		t.Errorf("Handler() failed: %v", err)
	}
	if got != nil {
		t.Errorf("Handler() = %v, want nil", got)
	}
}

func TestNewSimpleMapper_Errors(t *testing.T) {
	badSchema := func() (*ytypes.Schema, error) { return nil, fmt.Errorf("bad schema") }
	tests := []struct {
		name      string
		inSchema  SchemaFn
		outSchema SchemaFn
		mappings  map[string]string
	}{
		{
			name:     "bad input schema",
			inSchema: badSchema,
			outSchema: func() (*ytypes.Schema, error) {
				s, _ := openconfig.Schema()
				return s, nil
			},
		},
		{
			name: "bad output schema",
			inSchema: func() (*ytypes.Schema, error) {
				s, _ := openconfig.Schema()
				return s, nil
			},
			outSchema: badSchema,
		},
		{
			name:      "bad output path",
			inSchema:  openconfig.Schema,
			outSchema: openconfig.Schema,
			mappings:  map[string]string{"/interfaces/interface[name]": "/valid/path"},
		},
		{
			name:      "bad input path",
			inSchema:  openconfig.Schema,
			outSchema: openconfig.Schema,
			mappings:  map[string]string{"/valid/path": "/interfaces/interface[name]"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewSimpleMapper(tc.inSchema, tc.outSchema, tc.mappings, nil)
			if err == nil {
				t.Errorf("NewSimpleMapper() for %s succeeded, want error", tc.name)
			}
		})
	}
}

func TestUpdateHandler_Errors(t *testing.T) {
	isc, err := openconfig.Schema()
	if err != nil {
		t.Fatalf("Failed to load input schema: %v", err)
	}
	osc, err := openconfig.Schema()
	if err != nil {
		t.Fatalf("Failed to load output schema: %v", err)
	}
	m, _ := NewSimpleMapper(openconfig.Schema, openconfig.Schema, map[string]string{
		"/openconfig/interfaces/interface[name=<name>]/config/description": "/openconfig/interfaces/interface[name=<name>]/config/description",
	}, nil)

	tests := []struct {
		name         string
		notification *gnmipb.Notification
	}{
		{
			name: "unmarshal failure",
			notification: &gnmipb.Notification{
				Update: []*gnmipb.Update{{
					Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "non-existent"}}},
					Val:  &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "foo"}},
				}},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := m.updateHandler(isc, osc, tc.notification)
			if err == nil {
				t.Errorf("updateHandler() for %s succeeded, want error", tc.name)
			}
		})
	}

	t.Run("applyBind failure", func(t *testing.T) {
		m, _ := NewSimpleMapper(openconfig.Schema, openconfig.Schema, map[string]string{
			"/openconfig/interfaces/interface[name=<name>]/config/description": "/openconfig/interfaces/interface[name=<missing>]/config/description",
		}, nil)
		notification := &gnmipb.Notification{
			Update: []*gnmipb.Update{{
				Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{
					{Name: "interfaces"},
					{Name: "interface", Key: map[string]string{"name": "eth0"}},
					{Name: "config"},
					{Name: "description"},
				}},
				Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "foo"}},
			}},
		}
		_, err := m.updateHandler(isc, osc, notification)
		if err == nil {
			t.Errorf("updateHandler() succeeded for missing binding in output, want error")
		}
	})

	t.Run("GetNode error skipping", func(t *testing.T) {
		m, _ := NewSimpleMapper(openconfig.Schema, openconfig.Schema, map[string]string{
			"/openconfig/interfaces/interface[name=<name>]/config/description": "/interfaces/interface[name=eth0]/config/non-existent",
		}, nil)
		notification := &gnmipb.Notification{
			Update: []*gnmipb.Update{{
				Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{
					{Name: "interfaces"},
					{Name: "interface", Key: map[string]string{"name": "eth0"}},
					{Name: "config"},
					{Name: "description"},
				}},
				Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "foo"}},
			}},
		}
		_, err := m.updateHandler(isc, osc, notification)
		if err != nil {
			t.Errorf("updateHandler() failed for non-existent input path, want success (skip): %v", err)
		}
	})

	t.Run("yangValToGNMIVal failure", func(t *testing.T) {
		m, _ := NewSimpleMapper(openconfig.Schema, openconfig.Schema, map[string]string{
			"/openconfig/interfaces/interface[name=<name>]/state/last-change": "/openconfig/interfaces/interface[name=<name>]/config/description",
		}, nil)
		// last-change is uint64 in OC, which is supported by Unmarshal but not by yangValToGNMIVal.
		notification := &gnmipb.Notification{
			Update: []*gnmipb.Update{{
				Path: &gnmipb.Path{Elem: []*gnmipb.PathElem{
					{Name: "interfaces"},
					{Name: "interface", Key: map[string]string{"name": "eth0"}},
					{Name: "state"},
					{Name: "last-change"},
				}},
				Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: 12345}},
			}},
		}
		_, err := m.updateHandler(isc, osc, notification)
		if err == nil {
			t.Errorf("updateHandler() succeeded for unsupported type, want error")
		}
	})
}

func TestUpdateHandler_PrefixFlattening(t *testing.T) {
	isc, err := openconfig.Schema()
	if err != nil {
		t.Fatalf("Failed to load input schema: %v", err)
	}
	osc, err := openconfig.Schema()
	if err != nil {
		t.Fatalf("Failed to load output schema: %v", err)
	}
	// Map leaves that share a common prefix.
	m, err := NewSimpleMapper(openconfig.Schema, openconfig.Schema, map[string]string{
		"/openconfig/interfaces/interface[name=<name>]/config/description": "/openconfig/interfaces/interface[name=<name>]/config/description",
	}, nil)
	if err != nil {
		t.Fatalf("NewSimpleMapper() failed: %v", err)
	}

	notification := &gnmipb.Notification{
		Update: []*gnmipb.Update{
			{
				Path: &gnmipb.Path{
					Elem: []*gnmipb.PathElem{
						{Name: "interfaces"},
						{Name: "interface", Key: map[string]string{"name": "eth0"}},
						{Name: "config"},
						{Name: "description"},
					},
				},
				Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "foo"}},
			},
			{
				Path: &gnmipb.Path{
					Elem: []*gnmipb.PathElem{
						{Name: "interfaces"},
						{Name: "interface", Key: map[string]string{"name": "eth1"}},
						{Name: "config"},
						{Name: "description"},
					},
				},
				Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: "bar"}},
			},
		},
	}

	got, err := m.updateHandler(isc, osc, notification)
	if err != nil {
		t.Fatalf("updateHandler() failed: %v", err)
	}

	if got == nil {
		t.Fatal("updateHandler() returned nil")
	}

	// Verify that the prefix is empty (flattened)
	if len(got.Prefix.GetElem()) != 0 {
		t.Errorf("updateHandler() prefix not flattened: %v", got.Prefix)
	}

	// Verify that update paths are full paths.
	// Since we send /interfaces/interface/config/description, full path has 4 elements.
	// /interfaces/interface/name is also a full path but has only 3 elements, so we check for prefix instead.
	for _, u := range got.Update {
		if u.Path.GetElem()[0].Name != "interfaces" {
			t.Errorf("updateHandler() update path is not a full path (not flattened?): %v", u.Path)
		}
	}
}
