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

package ftutilities

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

func TestJoin(t *testing.T) {
	tests := []struct {
		name string
		p1   *gnmipb.Path
		p2   *gnmipb.Path
		want *gnmipb.Path
	}{
		{
			name: "success p1 override",
			p1: &gnmipb.Path{
				Origin: "p1_origin",
				Target: "p1_target",
				Elem:   []*gnmipb.PathElem{{Name: "path"}, {Name: "one"}},
			},
			p2: &gnmipb.Path{
				Origin: "p2_origin",
				Target: "p2_target",
				Elem:   []*gnmipb.PathElem{{Name: "path"}, {Name: "two"}},
			},
			want: &gnmipb.Path{
				Origin: "p1_origin",
				Target: "p1_target",
				Elem:   []*gnmipb.PathElem{{Name: "path"}, {Name: "one"}, {Name: "path"}, {Name: "two"}},
			},
		},
		{
			name: "success p1 missing origin",
			p1: &gnmipb.Path{
				Target: "p1_target",
				Elem:   []*gnmipb.PathElem{{Name: "path"}, {Name: "one"}},
			},
			p2: &gnmipb.Path{
				Origin: "p2_origin",
				Target: "p2_target",
				Elem:   []*gnmipb.PathElem{{Name: "path"}, {Name: "two"}},
			},
			want: &gnmipb.Path{
				Origin: "p2_origin",
				Target: "p1_target",
				Elem:   []*gnmipb.PathElem{{Name: "path"}, {Name: "one"}, {Name: "path"}, {Name: "two"}},
			},
		},
		{
			name: "success p1 missing target",
			p1: &gnmipb.Path{
				Origin: "p1_origin",
				Elem:   []*gnmipb.PathElem{{Name: "path"}, {Name: "one"}},
			},
			p2: &gnmipb.Path{
				Origin: "p2_origin",
				Target: "p2_target",
				Elem:   []*gnmipb.PathElem{{Name: "path"}, {Name: "two"}},
			},
			want: &gnmipb.Path{
				Origin: "p1_origin",
				Target: "p2_target",
				Elem:   []*gnmipb.PathElem{{Name: "path"}, {Name: "one"}, {Name: "path"}, {Name: "two"}},
			},
		},
		{
			name: "success target and origin missing",
			p1: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{{Name: "path"}, {Name: "one"}},
			},
			p2: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{{Name: "path"}, {Name: "two"}},
			},
			want: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{{Name: "path"}, {Name: "one"}, {Name: "path"}, {Name: "two"}},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Join(tc.p1, tc.p2)
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Join(%v, %v) returned an unexpected diff (-want +got): %v", tc.p1, tc.p2, diff)
			}
		})
	}
}

func TestStringToPath(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    *gnmipb.Path
		wantErr bool
	}{
		{
			name: "success with openconfig origin",
			s:    "/openconfig/some/path",
			want: &gnmipb.Path{
				Origin: "openconfig",
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "path",
					},
				},
			},
		},
		{
			name: "success with eos_native origin",
			s:    "/eos_native/some/other/path",
			want: &gnmipb.Path{
				Origin: "eos_native",
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "other",
					},
					{
						Name: "path",
					},
				},
			},
		},
		{
			name: "missing leading slash",
			s:    "openconfig/some/path",
			want: &gnmipb.Path{
				Origin: "openconfig",
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "path",
					},
				},
			},
		},
		{
			name: "no origin",
			s:    "some/path",
			want: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "path",
					},
				},
			},
		},
		{
			name:    "empty string",
			s:       "",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := StringToPath(tc.s)
			if tc.wantErr != (err != nil) {
				t.Fatalf("stringToPath(%q) returned an unexpected error: %v\n want err %t", tc.s, err, tc.wantErr)
			}

			if tc.wantErr {
				return
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("StringToPath(%q) returned an unexpected diff (-want +got): %v", tc.s, diff)
			}
		})
	}
}

func TestConfigToState(t *testing.T) {
	tests := []struct {
		name string
		p    *gnmipb.Path
		want *gnmipb.Path
	}{
		{
			name: "success",
			p: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "config",
					},
					{
						Name: "leaf",
					},
				},
			},
			want: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "state",
					},
					{
						Name: "leaf",
					},
				},
			},
		},
		{
			name: "no config",
			p: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "leaf",
					},
				},
			},
			want: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "leaf",
					},
				},
			},
		},
		{
			name: "no elements",
			p:    &gnmipb.Path{},
			want: &gnmipb.Path{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ConfigToState(tc.p)
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("ConfigToState(%v) returned an unexpected diff (-want +got): %v", tc.p, diff)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	twoNotifs := &gnmipb.Notification{
		Prefix: &gnmipb.Path{
			Elem: []*gnmipb.PathElem{
				{
					Name: "some",
				},
				{
					Name: "prefix",
				},
			},
		},
		Timestamp: 123,
		Update: []*gnmipb.Update{
			{
				Path: &gnmipb.Path{
					Elem: []*gnmipb.PathElem{
						{
							Name: "some",
						},
						{
							Name: "path",
						},
					},
				},
				Val: &gnmipb.TypedValue{
					Value: &gnmipb.TypedValue_StringVal{
						StringVal: "some_value",
					},
				},
			},
			{
				Path: &gnmipb.Path{
					Elem: []*gnmipb.PathElem{
						{
							Name: "some",
						},
						{
							Name: "other",
						},
						{
							Name: "path",
						},
					},
				},
				Val: &gnmipb.TypedValue{
					Value: &gnmipb.TypedValue_StringVal{
						StringVal: "some_other_value",
					},
				},
			},
		},
	}
	updateAndDeleteNotif := &gnmipb.Notification{
		Prefix: &gnmipb.Path{
			Elem: []*gnmipb.PathElem{
				{
					Name: "some",
				},
				{
					Name: "prefix",
				},
			},
		},
		Timestamp: 123,
		Update: []*gnmipb.Update{
			{
				Path: &gnmipb.Path{
					Elem: []*gnmipb.PathElem{
						{
							Name: "some",
						},
						{
							Name: "path",
						},
					},
				},
				Val: &gnmipb.TypedValue{
					Value: &gnmipb.TypedValue_StringVal{
						StringVal: "some_value",
					},
				},
			},
			{
				Path: &gnmipb.Path{
					Elem: []*gnmipb.PathElem{
						{
							Name: "some",
						},
						{
							Name: "additional",
						},
						{
							Name: "other",
						},
						{
							Name: "path",
						},
					},
				},
				Val: &gnmipb.TypedValue{
					Value: &gnmipb.TypedValue_IntVal{
						IntVal: 1,
					},
				},
			},
		},
		Delete: []*gnmipb.Path{
			{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "other",
					},
					{
						Name: "path",
					},
				},
			},
		},
	}

	originInPrefix := &gnmipb.Notification{
		Prefix: &gnmipb.Path{
			Origin: "openconfig",
			Elem: []*gnmipb.PathElem{
				{
					Name: "some",
				},
				{
					Name: "prefix",
				},
			},
		},
		Timestamp: 123,
		Update: []*gnmipb.Update{
			{
				Path: &gnmipb.Path{
					Elem: []*gnmipb.PathElem{
						{
							Name: "suffix",
						},
					},
				},
				Val: &gnmipb.TypedValue{
					Value: &gnmipb.TypedValue_StringVal{
						StringVal: "string_val",
					},
				},
			},
		},
	}

	updateWithOrigins := &gnmipb.Notification{
		Prefix: &gnmipb.Path{
			Elem: []*gnmipb.PathElem{
				{
					Name: "some",
				},
				{
					Name: "prefix",
				},
			},
		},
		Timestamp: 123,
		Update: []*gnmipb.Update{
			{
				Path: &gnmipb.Path{
					Origin: "openconfig",
					Elem: []*gnmipb.PathElem{
						{
							Name: "some",
						},
						{
							Name: "path",
						},
					},
				},
				Val: &gnmipb.TypedValue{
					Value: &gnmipb.TypedValue_StringVal{
						StringVal: "some_value",
					},
				},
			},
			{
				Path: &gnmipb.Path{
					Origin: "eos_native",
					Elem: []*gnmipb.PathElem{
						{
							Name: "some",
						},
						{
							Name: "other",
						},
						{
							Name: "path",
						},
					},
				},
				Val: &gnmipb.TypedValue{
					Value: &gnmipb.TypedValue_IntVal{
						IntVal: 1,
					},
				},
			},
		},
	}

	removeOther := func(path *gnmipb.Path, isDelete bool) bool {
		for _, elem := range path.GetElem() {
			if elem.GetName() == "other" {
				return false
			}
		}
		return true
	}
	removeOtherFromDelete := func(path *gnmipb.Path, isDelete bool) bool {
		if isDelete {
			for _, elem := range path.GetElem() {
				if elem.GetName() == "other" {
					return false
				}
			}
		}
		return true
	}

	tests := []struct {
		name         string
		notification *gnmipb.Notification
		fn           func(path *gnmipb.Path, isDelete bool) bool
		want         *gnmipb.Notification
	}{
		{
			name:         "success",
			notification: twoNotifs,
			fn: func(path *gnmipb.Path, isDelete bool) bool {
				for _, elem := range path.GetElem() {
					if elem.GetName() == "other" {
						return false
					}
				}
				return true
			},
			want: &gnmipb.Notification{
				Prefix: &gnmipb.Path{
					Elem: []*gnmipb.PathElem{
						{
							Name: "some",
						},
						{
							Name: "prefix",
						},
					},
				},
				Timestamp: 123,
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{
									Name: "some",
								},
								{
									Name: "path",
								},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "some_value",
							},
						},
					},
				},
			},
		},
		{
			name: "success, remove all items with origin in prefix",
			fn: func(path *gnmipb.Path, isDelete bool) bool {
				return path.GetOrigin() == ""
			},
			notification: originInPrefix,
			want:         nil,
		},
		{
			name: "success, remove all items with origin in path",
			fn: func(path *gnmipb.Path, isDelete bool) bool {
				return path.GetOrigin() == ""
			},
			notification: updateWithOrigins,
			want:         nil,
		},
		{
			name: "no updates",
			notification: &gnmipb.Notification{
				Prefix: &gnmipb.Path{
					Elem: []*gnmipb.PathElem{
						{
							Name: "some",
						},
						{
							Name: "prefix",
						},
					},
				},
				Delete: []*gnmipb.Path{
					{
						Elem: []*gnmipb.PathElem{
							{
								Name: "some",
							},
							{
								Name: "path",
							},
						},
					},
					{
						Elem: []*gnmipb.PathElem{
							{
								Name: "some",
							},
							{
								Name: "other",
							},
							{
								Name: "path",
							},
						},
					},
				},
				Timestamp: 123,
			},
			fn: removeOther,
			want: &gnmipb.Notification{
				Prefix: &gnmipb.Path{
					Elem: []*gnmipb.PathElem{
						{
							Name: "some",
						},
						{
							Name: "prefix",
						},
					},
				},
				Delete: []*gnmipb.Path{
					{
						Elem: []*gnmipb.PathElem{
							{
								Name: "some",
							},
							{
								Name: "path",
							},
						},
					},
				},
				Timestamp: 123,
			},
		},
		{
			name:         "remove other deletes",
			notification: updateAndDeleteNotif,
			fn:           removeOtherFromDelete,
			want: &gnmipb.Notification{
				Prefix: &gnmipb.Path{
					Elem: []*gnmipb.PathElem{
						{
							Name: "some",
						},
						{
							Name: "prefix",
						},
					},
				},
				Timestamp: 123,
				Update: []*gnmipb.Update{
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{
									Name: "some",
								},
								{
									Name: "path",
								},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_StringVal{
								StringVal: "some_value",
							},
						},
					},
					{
						Path: &gnmipb.Path{
							Elem: []*gnmipb.PathElem{
								{
									Name: "some",
								},
								{
									Name: "additional",
								},
								{
									Name: "other",
								},
								{
									Name: "path",
								},
							},
						},
						Val: &gnmipb.TypedValue{
							Value: &gnmipb.TypedValue_IntVal{
								IntVal: 1,
							},
						},
					},
				},
			},
		},
		{
			name:         "remove all",
			notification: twoNotifs,
			fn:           func(path *gnmipb.Path, isDelete bool) bool { return false },
			want:         nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Filter(tc.notification, tc.fn)
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Filter() returned an unexpected diff (-want +got): %v", diff)
			}
		})
	}
}

func TestFilterUpdates(t *testing.T) {
	singleUpdate := []*gnmipb.Update{
		{
			Path: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "state",
					},
					{
						Name: "leaf",
					},
				},
			},
			Val: &gnmipb.TypedValue{
				Value: &gnmipb.TypedValue_UintVal{UintVal: 0},
			},
		},
	}
	twoUpdates := []*gnmipb.Update{
		{
			Path: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "config",
					},
					{
						Name: "leaf",
					},
				},
			},
			Val: &gnmipb.TypedValue{
				Value: &gnmipb.TypedValue_UintVal{UintVal: 0},
			},
		},
		{
			Path: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "state",
					},
					{
						Name: "leaf",
					},
				},
			},
			Val: &gnmipb.TypedValue{
				Value: &gnmipb.TypedValue_UintVal{UintVal: 0},
			},
		},
	}
	tests := []struct {
		name   string
		update []*gnmipb.Update
		fn     func(up *gnmipb.Update) bool
		want   []*gnmipb.Update
	}{
		{
			name:   "return all",
			update: singleUpdate,
			fn:     func(up *gnmipb.Update) bool { return true },
			want:   singleUpdate,
		},
		{
			name:   "return none",
			update: singleUpdate,
			fn:     func(up *gnmipb.Update) bool { return false },
			want:   nil,
		},
		{
			name:   "filter non-state leaves",
			update: twoUpdates,
			fn:     StateLeaves,
			want:   singleUpdate,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := FilterUpdates(tc.update, tc.fn)
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("FilterUpdates(%v) returned an unexpected diff (-want +got): %v", tc.update, diff)
			}
		})
	}
}

func TestStringMapPaths(t *testing.T) {
	tests := []struct {
		name          string
		stringPathMap map[string][]string
		want          map[string][]*gnmipb.Path
		wantErr       bool
	}{
		{
			name: "success",
			stringPathMap: map[string][]string{
				"one":     {"a/b/one", "c/d/one"},
				"two":     {"c/d/two", "a/b/two"},
				"origins": {"openconfig/a/b/c", "eos_native/a/b/c"},
			},
			want: map[string][]*gnmipb.Path{
				"one": {
					{
						Elem: []*gnmipb.PathElem{
							{
								Name: "a",
							},
							{
								Name: "b",
							},
							{
								Name: "one",
							},
						},
					},
					{
						Elem: []*gnmipb.PathElem{
							{
								Name: "c",
							},
							{
								Name: "d",
							},
							{
								Name: "one",
							},
						},
					},
				},
				"two": {
					{
						Elem: []*gnmipb.PathElem{
							{
								Name: "a",
							},
							{
								Name: "b",
							},
							{
								Name: "two",
							},
						},
					},
					{
						Elem: []*gnmipb.PathElem{
							{
								Name: "c",
							},
							{
								Name: "d",
							},
							{
								Name: "two",
							},
						},
					},
				},
				"origins": {
					{
						Origin: "openconfig",
						Elem: []*gnmipb.PathElem{
							{
								Name: "a",
							},
							{
								Name: "b",
							},
							{
								Name: "c",
							},
						},
					},
					{
						Origin: "eos_native",
						Elem: []*gnmipb.PathElem{
							{
								Name: "a",
							},
							{
								Name: "b",
							},
							{
								Name: "c",
							},
						},
					},
				},
			},
		},
		{
			name: "empty path",
			stringPathMap: map[string][]string{
				"one": {""},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := StringMapPaths(tc.stringPathMap)
			if tc.wantErr != (err != nil) {
				t.Fatalf("StringMapPaths(%v) returned an unexpected error: %v", tc.stringPathMap, err)
			}
			if tc.wantErr {
				return
			}
			for k, gotPath := range got {
				// Sort returned paths by ygot string.
				sort.SliceStable(gotPath, SortByYgotString(gotPath))
				if diff := cmp.Diff(tc.want[k], gotPath, protocmp.Transform()); diff != "" {
					t.Errorf("StringMapPaths(%v) returned an unexpected diff (-want +got): %v", tc.stringPathMap, diff)
				}
			}
		})
	}
}

func TestMatchPath(t *testing.T) {
	tests := []struct {
		name    string
		path    *gnmipb.Path
		pattern *gnmipb.Path
		want    bool
	}{
		{
			name:    "success",
			path:    &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "just"}, {Name: "some"}, {Name: "paths"}}},
			pattern: &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "just"}, {Name: "some"}, {Name: "paths"}}},
			want:    true,
		},
		{
			name:    "path mismatch",
			path:    &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "just"}, {Name: "some"}, {Name: "paths"}}},
			pattern: &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "just"}, {Name: "some"}, {Name: "other"}}},
			want:    false,
		},
		{
			name:    "wildcard match",
			path:    &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "just"}, {Name: "some"}, {Name: "paths"}}},
			pattern: &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "just"}, {Name: "some"}, {Name: "*"}}},
			want:    true,
		},
		{
			name: "match ignores keys",
			path: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "path",
						Key: map[string]string{
							"key": "value",
						},
					},
				},
			},
			pattern: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "path",
					},
				},
			},
			want: true,
		},
		{
			name: "wildcard match ignores keys",
			path: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "path",
						Key: map[string]string{
							"key": "value",
						},
					},
				},
			},
			pattern: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "*",
					},
				},
			},
			want: true,
		},
		{
			name:    "nil",
			path:    nil,
			pattern: nil,
			want:    true,
		},
		{
			name:    "path nil",
			path:    nil,
			pattern: &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "just"}, {Name: "some"}, {Name: "paths"}}},
			want:    false,
		},
		{
			name:    "pattern nil",
			path:    &gnmipb.Path{Elem: []*gnmipb.PathElem{{Name: "just"}, {Name: "some"}, {Name: "paths"}}},
			pattern: nil,
			want:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := MatchPath(tc.path, tc.pattern)
			if got != tc.want {
				t.Errorf("MatchPath(%v, %v) = %v, want: %v", tc.path, tc.pattern, got, tc.want)
			}
		})
	}
}

func TestGNMIPathToSchemaStrings(t *testing.T) {
	tests := []struct {
		name string
		path *gnmipb.Path
		// Expected results when setOCOriginIfMissing is true/false.
		wantSetOriginTrue  []string
		wantSetOriginFalse []string
	}{
		{
			name: "success with origin",
			path: &gnmipb.Path{
				Origin: "openconfig",
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "path",
					},
				},
			},
			wantSetOriginTrue:  []string{"openconfig", "some", "path"},
			wantSetOriginFalse: []string{"openconfig", "some", "path"},
		},
		{
			name: "success without origin, add openconfig",
			path: &gnmipb.Path{
				Elem: []*gnmipb.PathElem{
					{
						Name: "some",
					},
					{
						Name: "path",
					},
				},
			},
			wantSetOriginTrue:  []string{"openconfig", "some", "path"},
			wantSetOriginFalse: []string{"some", "path"},
		},
		{
			name:               "success origin set",
			path:               &gnmipb.Path{Origin: "eos_native"},
			wantSetOriginTrue:  []string{"eos_native"},
			wantSetOriginFalse: []string{"eos_native"},
		},
		{
			name:               "success target set",
			path:               &gnmipb.Path{Target: "some_target"},
			wantSetOriginTrue:  []string{"openconfig"},
			wantSetOriginFalse: nil,
		},
		{
			name:               "success empty notitification",
			path:               &gnmipb.Path{},
			wantSetOriginTrue:  nil,
			wantSetOriginFalse: nil,
		},
		{
			name:               "success nil",
			path:               nil,
			wantSetOriginTrue:  nil,
			wantSetOriginFalse: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotSetOriginTrue := GNMIPathToSchemaStrings(tc.path, true)
			gotSetOriginFalse := GNMIPathToSchemaStrings(tc.path, false)
			if diff := cmp.Diff(tc.wantSetOriginTrue, gotSetOriginTrue); diff != "" {
				t.Errorf("GNMIPathToSchemaStrings(%v, true) returned an unexpected diff (-want +got): %v", tc.path, diff)
			}
			if diff := cmp.Diff(tc.wantSetOriginFalse, gotSetOriginFalse); diff != "" {
				t.Errorf("GNMIPathToSchemaStrings(%v, false) returned an unexpected diff (-want +got): %v", tc.path, diff)
			}
		})
	}
}

func TestCreateOrGetCKN(t *testing.T) {
	tests := []struct {
		name                  string
		CKN                   string
		setup                 func(intfInfo *InterfaceMacSecInfo)
		expectNewCKNMap       bool
		expectCKNInMap        bool
		expectedTotalCKNCount int
	}{
		{
			name: "Create CKN2 in map that already has CKN1",
			CKN:  "CKN2",
			setup: func(intfInfo *InterfaceMacSecInfo) {
				intfInfo.cknStatuses = map[string]*CKNInfo{
					"CKN1": {
						principal:    true,
						success:      true,
						principalSet: true,
						successSet:   true,
					},
				}
			},
			expectNewCKNMap:       true,
			expectCKNInMap:        true,
			expectedTotalCKNCount: 2,
		},
		{
			name: "Retrieve existing CKN1",
			CKN:  "CKN1",
			setup: func(intfInfo *InterfaceMacSecInfo) {
				intfInfo.cknStatuses = map[string]*CKNInfo{
					"CKN1": {
						principal:    true,
						success:      true,
						principalSet: true,
						successSet:   true,
					},
				}
			},
			expectNewCKNMap:       false,
			expectCKNInMap:        true,
			expectedTotalCKNCount: 1,
		},
		{
			name: "Create CKN when cknStatuses map is initially nil",
			CKN:  "CKN3",
			setup: func(intfInfo *InterfaceMacSecInfo) {
				intfInfo.cknStatuses = nil
			},
			expectNewCKNMap:       true,
			expectCKNInMap:        true,
			expectedTotalCKNCount: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testInterfaceInfo := &InterfaceMacSecInfo{
				interfaceName: "Ethernet22",
			}
			if tc.setup != nil {
				tc.setup(testInterfaceInfo)
			}

			var preExistingCKN1Info *CKNInfo
			if testInterfaceInfo.cknStatuses != nil {
				preExistingCKN1Info = testInterfaceInfo.cknStatuses["CKN1"]
			}

			gotCKNInfo := testInterfaceInfo.CreateOrGetCKN(tc.CKN)

			if gotCKNInfo == nil {
				t.Fatalf("CreateOrGetCKN(%q) returned nil, expected non-nil", tc.CKN)
			}
			if tc.expectNewCKNMap {
				expectedFreshCKN := &CKNInfo{}
				if diff := cmp.Diff(expectedFreshCKN, gotCKNInfo, cmp.AllowUnexported(CKNInfo{})); diff != "" {
					t.Errorf("CreateOrGetCKN(%q) returned a CKNInfo with non-default values for a new CKN (-want +got):\n%s", tc.CKN, diff)
				}
			}
			if tc.CKN == "CKN1" && preExistingCKN1Info != nil {
				if gotCKNInfo != preExistingCKN1Info {
					t.Errorf("CreateOrGetCKN(%q) expected to return existing instance %p, but got %p", tc.CKN, preExistingCKN1Info, gotCKNInfo)
				}
			}
			finalCKNMap := testInterfaceInfo.CloneStatuses()
			if tc.expectCKNInMap {
				if _, ok := finalCKNMap[tc.CKN]; !ok {
					t.Errorf("ckn %q not found in internal map after CreateOrGetCKN", tc.CKN)
				}
			}
			if len(finalCKNMap) != tc.expectedTotalCKNCount {
				t.Errorf("expected total %d CKNs in map, got %d", tc.expectedTotalCKNCount, len(finalCKNMap))
			}
		})
	}
}

func TestResetCPStatus(t *testing.T) {
	tests := []struct {
		name                string
		initialCPStatus     bool
		initialCPStatusSet  bool
		expectedCPStatus    bool
		expectedCPStatusSet bool
	}{
		{
			name:                "cpStatus true, cpStatusSet true",
			initialCPStatus:     true,
			initialCPStatusSet:  true,
			expectedCPStatus:    false,
			expectedCPStatusSet: false,
		},
		{
			name:                "cpStatus false, cpStatusSet true",
			initialCPStatus:     false,
			initialCPStatusSet:  true,
			expectedCPStatus:    false,
			expectedCPStatusSet: false,
		},
		{
			name:                "cpStatus false, cpStatusSet false",
			initialCPStatus:     false,
			initialCPStatusSet:  false,
			expectedCPStatus:    false,
			expectedCPStatusSet: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ifaceInfo := &InterfaceMacSecInfo{}
			if tc.initialCPStatusSet {
				ifaceInfo.SetIntfCPStatus(tc.initialCPStatus)
			}
			ifaceInfo.ResetCPStatus()
			finalStatus, finalSet := ifaceInfo.IntfCPStatus()

			if finalStatus != tc.expectedCPStatus {
				t.Errorf("After ResetCPStatus(), cpStatus = %v, want %v (initial cpStatus=%v, initial cpStatusSet=%v)", finalStatus, tc.expectedCPStatus, tc.initialCPStatus, tc.initialCPStatusSet)
			}
			if finalSet != tc.expectedCPStatusSet {
				t.Errorf("After ResetCPStatus(), cpStatusSet = %v, want %v (initial cpStatus=%v, initial cpStatusSet=%v)", finalSet, tc.expectedCPStatusSet, tc.initialCPStatus, tc.initialCPStatusSet)
			}
		})
	}
}

func TestIntfSuccess(t *testing.T) {
	tests := []struct {
		name                 string
		CKN                  string
		setup                func(intfInfo *InterfaceMacSecInfo)
		expectedCKNStatus    bool
		expectedCKNStatusSet bool
	}{
		{
			name: "CKN exists, success is true",
			CKN:  "CKN1",
			setup: func(intfInfo *InterfaceMacSecInfo) {
				intfInfo.cknStatuses = map[string]*CKNInfo{
					"CKN1": {
						principal:    true,
						success:      true,
						principalSet: true,
						successSet:   true,
					},
				}
			},
			expectedCKNStatus:    true,
			expectedCKNStatusSet: true,
		},
		{
			name: "CKN exists, map is empty",
			CKN:  "CKN1",
			setup: func(intfInfo *InterfaceMacSecInfo) {
				intfInfo.cknStatuses = map[string]*CKNInfo{}
			},
			expectedCKNStatus:    false,
			expectedCKNStatusSet: false,
		},
		{
			name: "CKN does not exist in the map",
			CKN:  "CKN2",
			setup: func(intfInfo *InterfaceMacSecInfo) {
				intfInfo.cknStatuses = map[string]*CKNInfo{
					"CKN1": {
						principal:    true,
						success:      true,
						principalSet: true,
						successSet:   true,
					},
				}
			},
			expectedCKNStatus:    false,
			expectedCKNStatusSet: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testInterfaceInfo := &InterfaceMacSecInfo{
				interfaceName: "Ethernet22",
			}
			if tc.setup != nil {
				tc.setup(testInterfaceInfo)
			}
			gotCKNStatus, gotCKNStatusSet := testInterfaceInfo.IntfSuccess(tc.CKN)
			if gotCKNStatus != tc.expectedCKNStatus {
				t.Errorf("IntfSuccess(%q) = %v, want %v", tc.CKN, gotCKNStatus, tc.expectedCKNStatus)
			}
			if gotCKNStatusSet != tc.expectedCKNStatusSet {
				t.Errorf("IntfSuccess(%q) = %v, want %v", tc.CKN, gotCKNStatusSet, tc.expectedCKNStatusSet)
			}
		})
	}
}

func TestInterfaceMacSecInfo(t *testing.T) {
	info := &InterfaceMacSecInfo{
		interfaceName: "Ethernet1",
		cknStatuses:   make(map[string]*CKNInfo),
	}
	ckn1 := "ckn1"
	ckn2 := "ckn2"

	// This section tests the CP Status.
	if _, ok := info.IntfCPStatus(); ok {
		t.Error("IntfCPStatus should not be set initially")
	}
	info.SetIntfCPStatus(true)
	if status, ok := info.IntfCPStatus(); !ok || !status {
		t.Errorf("IntfCPStatus() = %v, %v, want true, true", status, ok)
	}
	info.ResetCPStatus()
	if _, ok := info.IntfCPStatus(); ok {
		t.Error("IntfCPStatus should be unset after ResetCPStatus")
	}

	// This section tests the CreateOrGetCKN function.
	cknInfo1 := info.CreateOrGetCKN(ckn1)
	if cknInfo1 == nil {
		t.Fatalf("CreateOrGetCKN(%q) returned nil", ckn1)
	}
	if cknInfo1Again := info.CreateOrGetCKN(ckn1); cknInfo1Again != cknInfo1 {
		t.Errorf("CreateOrGetCKN(%q) returned a new instance, want the existing one", ckn1)
	}

	// This section tests the Principal status.
	if _, ok := info.IntfPrincipal(ckn1); ok {
		t.Errorf("IntfPrincipal(%q) should not be set initially", ckn1)
	}
	info.SetIntfPrincipal(ckn1, true)
	if principal, ok := info.IntfPrincipal(ckn1); !ok || !principal {
		t.Errorf("IntfPrincipal(%q) = %v, %v, want true, true", ckn1, principal, ok)
	}

	// This section tests the Success status.
	if _, ok := info.IntfSuccess(ckn1); ok {
		t.Errorf("IntfSuccess(%q) should not be set initially", ckn1)
	}
	info.SetIntfSuccess(ckn1, true)
	if success, ok := info.IntfSuccess(ckn1); !ok || !success {
		t.Errorf("IntfSuccess(%q) = %v, %v, want true, true", ckn1, success, ok)
	}

	// This section tests the IsComplete function.
	if !info.IsComplete(ckn1) {
		t.Errorf("IsComplete(%q) = false, want true", ckn1)
	}
	info.CreateOrGetCKN(ckn2)
	info.SetIntfPrincipal(ckn2, true)
	if info.IsComplete(ckn2) {
		t.Errorf("IsComplete(%q) = true, want false", ckn2)
	}

	// This section tests the CloneStatuses function.
	cloned := info.CloneStatuses()
	if len(cloned) != 2 {
		t.Errorf("len(CloneStatuses()) = %d, want 2", len(cloned))
	}
	if _, ok := cloned[ckn1]; !ok {
		t.Errorf("Cloned statuses missing ckn %q", ckn1)
	}
	// This modifies the original and checks that the clone is unaffected.
	info.cknStatuses[ckn1].success = false
	if cloned[ckn1].success != false {
		t.Error("CloneStatuses() did not create a shallow copy")
	}

	// This section tests the RemoveCkn function.
	info.RemoveCkn(ckn1)
	if _, ok := info.cknStatuses[ckn1]; ok {
		t.Errorf("CKN %q should have been removed", ckn1)
	}
	if _, ok := info.IntfPrincipal(ckn1); ok {
		t.Errorf("IntfPrincipal(%q) should be unset after removal", ckn1)
	}
}

func TestTargetMacSecInfo(t *testing.T) {
	targetInfo := NewTargetMacSecInfo("hostname1")
	if targetInfo.TargetHostname != "hostname1" {
		t.Errorf("NewTargetMacSecInfo() hostname = %q, want %q", targetInfo.TargetHostname, "hostname1")
	}
	if targetInfo.Interfaces == nil {
		t.Error("NewTargetMacSecInfo() Interfaces map is nil")
	}

	intf1 := "Ethernet1"
	intf2 := "Ethernet2"

	// Test CreateOrGetInterface for a new interface.
	ifaceInfo1 := targetInfo.CreateOrGetInterface(intf1)
	if ifaceInfo1 == nil {
		t.Fatalf("CreateOrGetInterface(%q) returned nil", intf1)
	}
	if ifaceInfo1.interfaceName != intf1 {
		t.Errorf("ifaceInfo1.interfaceName = %q, want %q", ifaceInfo1.interfaceName, intf1)
	}

	// Test CreateOrGetInterface for an existing interface.
	ifaceInfo1Again := targetInfo.CreateOrGetInterface(intf1)
	if ifaceInfo1Again != ifaceInfo1 {
		t.Errorf("CreateOrGetInterface(%q) returned a new instance, want the existing one", intf1)
	}

	// Test InterfaceInfo for an existing interface.
	retrievedIface, ok := targetInfo.InterfaceInfo(intf1)
	if !ok {
		t.Errorf("InterfaceInfo(%q): ok = false, want true", intf1)
	}
	if retrievedIface != ifaceInfo1 {
		t.Errorf("InterfaceInfo(%q) returned wrong info", intf1)
	}

	// Test InterfaceInfo for a non-existent interface.
	_, ok = targetInfo.InterfaceInfo(intf2)
	if ok {
		t.Errorf("InterfaceInfo(%q): ok = true, want false", intf2)
	}

	// Test ClearInterfaceInfo.
	targetInfo.CreateOrGetInterface(intf2) // Add a second interface to ensure it's not deleted.
	targetInfo.ClearInterfaceInfo(intf1)
	_, ok = targetInfo.InterfaceInfo(intf1)
	if ok {
		t.Errorf("InterfaceInfo(%q) after Clear: ok = true, want false", intf1)
	}

	// Verify the other interface is still present.
	if _, ok := targetInfo.InterfaceInfo(intf2); !ok {
		t.Errorf("InterfaceInfo(%q) was deleted unexpectedly", intf2)
	}
}

func TestAristaMACSecMapCache(t *testing.T) {
	cache := &AristaMACSecMapCache{
		data: make(map[string]*TargetMacSecInfo),
	}

	target1 := "host1"
	target2 := "host2"

	// Test CreateOrUpdateTargetMacSecInfo for a new target.
	info1 := cache.CreateOrUpdateTargetMacSecInfo(target1)
	if info1 == nil {
		t.Fatalf("CreateOrUpdateTargetMacSecInfo(%q) returned nil, want non-nil", target1)
	}
	if info1.TargetHostname != target1 {
		t.Errorf("info1.TargetHostname = %q, want %q", info1.TargetHostname, target1)
	}

	// Test CreateOrUpdateTargetMacSecInfo for an existing target.
	info1Again := cache.CreateOrUpdateTargetMacSecInfo(target1)
	if info1Again != info1 {
		t.Errorf("CreateOrUpdateTargetMacSecInfo(%q) returned a new instance, want the existing one", target1)
	}

	// Test RetrieveTargetMacSecInfo for an existing target.
	retrievedInfo1, ok := cache.RetrieveTargetMacSecInfo(target1)
	if !ok {
		t.Errorf("RetrieveTargetMacSecInfo(%q): ok = false, want true", target1)
	}
	if retrievedInfo1 != info1 {
		t.Errorf("RetrieveTargetMacSecInfo(%q) returned wrong info", target1)
	}

	// Test RetrieveTargetMacSecInfo for a non-existent target.
	_, ok = cache.RetrieveTargetMacSecInfo(target2)
	if ok {
		t.Errorf("RetrieveTargetMacSecInfo(%q): ok = true, want false", target2)
	}

	// Test SetTargetMacSecInfo.
	info2 := NewTargetMacSecInfo(target2)
	cache.SetTargetMacSecInfo(target2, info2)
	retrievedInfo2, ok := cache.RetrieveTargetMacSecInfo(target2)
	if !ok {
		t.Errorf("RetrieveTargetMacSecInfo(%q) after Set: ok = false, want true", target2)
	}
	if retrievedInfo2 != info2 {
		t.Errorf("RetrieveTargetMacSecInfo(%q) after Set returned wrong info", target2)
	}

	// Test DeleteTargetMacSecInfo.
	cache.DeleteTargetMacSecInfo(target1)
	_, ok = cache.RetrieveTargetMacSecInfo(target1)
	if ok {
		t.Errorf("RetrieveTargetMacSecInfo(%q) after Delete: ok = true, want false", target1)
	}

	// Verify target2 is still there.
	_, ok = cache.RetrieveTargetMacSecInfo(target2)
	if !ok {
		t.Errorf("RetrieveTargetMacSecInfo(%q) was deleted unexpectedly", target2)
	}

	// Test ClearAllTargetMacSecInfo.
	cache.CreateOrUpdateTargetMacSecInfo(target1) // Add target1 back.
	cache.ClearAllTargetMacSecInfo()
	if len(cache.data) != 0 {
		t.Errorf("cache.data length after ClearAll = %d, want 0", len(cache.data))
	}
	_, ok = cache.RetrieveTargetMacSecInfo(target1)
	if ok {
		t.Errorf("RetrieveTargetMacSecInfo(%q) after ClearAll: ok = true, want false", target1)
	}
	_, ok = cache.RetrieveTargetMacSecInfo(target2)
	if ok {
		t.Errorf("RetrieveTargetMacSecInfo(%q) after ClearAll: ok = true, want false", target2)
	}
}
