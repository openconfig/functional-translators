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
