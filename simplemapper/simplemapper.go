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

// Package simplemapper provides a simple 1-1 mapper for functional translators.
package simplemapper

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	log "github.com/golang/glog"
	"google.golang.org/protobuf/proto"
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
	"github.com/openconfig/functional-translators/ftutilities"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

// bindKeys extracts variable bindings from a concrete path by comparing it with a
// path containing variables. For example, given a varPath
// /interfaces/interface[name=<ifname>]/state/oper-status and a concretePath
// /interfaces/interface[name=eth0]/state/oper-status, it returns {"<ifname>": "eth0"}.
// varPath is a gNMI path where some key values are variables (e.g. "<ifname>").
// concretePath is a gNMI path with specific values for keys, usually from a gNMI notification.
func bindKeys(varPath *gnmipb.Path, concretePath *gnmipb.Path) (map[string]string, error) {
	ret := make(map[string]string)
	if len(varPath.GetElem()) != len(concretePath.GetElem()) {
		return nil, fmt.Errorf("path with variables and concrete path have different lengths: %d vs %d", len(varPath.GetElem()), len(concretePath.GetElem()))
	}
	for i, elemBind := range varPath.GetElem() {
		elem := concretePath.GetElem()[i]
		if elemBind.Name != elem.Name {
			return nil, fmt.Errorf("path with variables and concrete path have different elem names: %q vs %q", elemBind.Name, elem.Name)
		}
		if len(elemBind.GetKey()) != len(elem.GetKey()) {
			return nil, fmt.Errorf("path with variables and concrete path have different key lengths: %d vs %d", len(elemBind.GetKey()), len(elem.GetKey()))
		}
		for key, valBind := range elemBind.GetKey() {
			if !isVar(valBind) { // Constant key, e.g. /afi-safis/afi-safi[afi-safi-name=IPV4_UNICAST]/
				continue
			}
			val, ok := elem.GetKey()[key]
			if !ok {
				return nil, fmt.Errorf("key from path with variables not found in concrete path: %q", key)
			}
			if _, ok := ret[valBind]; ok {
				return nil, fmt.Errorf("duplicate var %q", valBind)
			}
			ret[valBind] = val
		}
	}
	return ret, nil
}

// applyBind substitutes variables in a path with their bound values.
// For example, given bindings {"<ifname>": "eth0"} and a varPath
// /interfaces/interface[name=<ifname>]/config/description, it returns a new path
// /interfaces/interface[name=eth0]/config/description.
func applyBind(bindings map[string]string, varPath *gnmipb.Path) (*gnmipb.Path, error) {
	ret := proto.Clone(varPath).(*gnmipb.Path)
	for _, elem := range ret.GetElem() {
		newKeys := make(map[string]string)
		for key, name := range elem.GetKey() {
			if isVar(name) {
				val, ok := bindings[name]
				if !ok {
					return nil, fmt.Errorf("variable %q not found in bindings", name)
				}
				newKeys[key] = val
			} else { // Constant key, e.g. /afi-safis/afi-safi[afi-safi-name=IPV4_UNICAST]/
				newKeys[key] = name
			}
		}
		elem.Key = newKeys
	}
	return ret, nil
}

func isVar(s string) bool {
	return strings.HasPrefix(s, "<") && strings.HasSuffix(s, ">")
}

func varsToWildcards(path *gnmipb.Path) *gnmipb.Path {
	ret := proto.Clone(path).(*gnmipb.Path)
	for _, elem := range ret.GetElem() {
		newKeys := make(map[string]string)
		for key, name := range elem.GetKey() {
			if !isVar(name) {
				newKeys[key] = name
				continue
			}
			newKeys[key] = "*"
		}
		elem.Key = newKeys
	}
	return ret
}

var errNilValue = errors.New("nil value")

// wrap is a generic helper that handles the nil pointer check and the gNMI wrapping logic.
func wrap[T any](ptr *T, wrapFn func(T) *gnmipb.TypedValue) (*gnmipb.TypedValue, error) {
	if ptr == nil {
		return nil, errNilValue
	}
	return wrapFn(*ptr), nil
}

// TODO(team): Support the rest of the types.
// TODO(team): Support automatic casting when types are different in input and output schemas.
func yangValToGNMIVal(val any) (*gnmipb.TypedValue, error) {
	switch v := val.(type) {
	case *string:
		return wrap(v, func(x string) *gnmipb.TypedValue {
			return &gnmipb.TypedValue{Value: &gnmipb.TypedValue_StringVal{StringVal: x}}
		})
	case *bool:
		return wrap(v, func(x bool) *gnmipb.TypedValue {
			return &gnmipb.TypedValue{Value: &gnmipb.TypedValue_BoolVal{BoolVal: x}}
		})
	case *float64:
		return wrap(v, func(x float64) *gnmipb.TypedValue {
			return &gnmipb.TypedValue{Value: &gnmipb.TypedValue_DoubleVal{DoubleVal: x}}
		})
	case nil:
		return nil, errNilValue
	default:
		return nil, fmt.Errorf("unsupported type: %T", v)
	}
}

// SchemaFn is a function that returns a ygot schema. These are typically auto-generated by ygot, e.g. interfaces.Schema.
type SchemaFn func() (*ytypes.Schema, error)

type pathMapping struct {
	input         *gnmipb.Path
	output        *gnmipb.Path
	inputWildcard *gnmipb.Path
}

// parseMapperPath converts a path string to a gNMI path and a schema path. It takes care of
// populating gnmi origin if the paths tarts with a "valid origin".
func parseMapperPath(pathStr string) (*gnmipb.Path, string, error) {
	pathStr = strings.TrimPrefix(pathStr, "/")
	origin, elems, _ := strings.Cut(pathStr, "/")
	if _, ok := ftutilities.ValidOrigins[origin]; !ok {
		origin = ""
		elems = pathStr
	}
	path, err := ygot.StringToStructuredPath(elems)
	if err != nil {
		return nil, "", fmt.Errorf("failed to convert path to gnmi path: %v", err)
	}
	path.Origin = origin
	schemaPath := ftutilities.GNMIPathToSchemaString(path, false)
	schemaPath = ftutilities.ForcePathPrefix(schemaPath, origin)
	return path, schemaPath, nil
}

// NewSimpleMapper creates a new simple mapper.
func NewSimpleMapper(inSchema, outSchema SchemaFn, outputToInput map[string]string, deleteHandler func(*gnmipb.Notification) ([]*gnmipb.Path, error)) (*SimpleMapper, error) {
	var mappings []pathMapping
	outputToInputSchemaStrings := make(map[string][]string)
	outputToInputSchemaMap := make(map[string]map[string]bool)
	for o, i := range outputToInput {
		oPath, oSchemaPath, err := parseMapperPath(o)
		if err != nil {
			return nil, err
		}
		iPath, iSchemaPath, err := parseMapperPath(i)
		if err != nil {
			return nil, err
		}
		mappings = append(mappings, pathMapping{
			input:         iPath,
			output:        oPath,
			inputWildcard: varsToWildcards(iPath),
		})
		if _, ok := outputToInputSchemaMap[oSchemaPath]; !ok {
			outputToInputSchemaMap[oSchemaPath] = make(map[string]bool)
		}
		outputToInputSchemaMap[oSchemaPath][iSchemaPath] = true
	}
	for o, is := range outputToInputSchemaMap {
		for i := range is {
			outputToInputSchemaStrings[o] = append(outputToInputSchemaStrings[o], i)
		}
		sort.Strings(outputToInputSchemaStrings[o])
	}

	isc, err := inSchema()
	if err != nil {
		return nil, fmt.Errorf("cannot load input schema, %w", err)
	}
	osc, err := outSchema()
	if err != nil {
		return nil, fmt.Errorf("cannot load output schema, %w", err)
	}

	return &SimpleMapper{
		inSchema:                   isc,
		outSchema:                  osc,
		mapEntries:                 mappings,
		deleteHandler:              deleteHandler,
		outputToInputSchemaStrings: outputToInputSchemaStrings,
	}, nil
}

// SimpleMapper objects translate notifications with simple 1-1 leaf relabeling.
type SimpleMapper struct {
	inSchema                   *ytypes.Schema
	outSchema                  *ytypes.Schema
	mapEntries                 []pathMapping
	deleteHandler              func(*gnmipb.Notification) ([]*gnmipb.Path, error)
	outputToInputSchemaStrings map[string][]string
}

// OutputToInputSchemaStrings returns the parsed schema path strings for use in the functional translator constructor.
func (m *SimpleMapper) OutputToInputSchemaStrings() map[string][]string {
	return m.outputToInputSchemaStrings
}

func (m *SimpleMapper) updateHandler(inSchema, outSchema *ytypes.Schema, notification *gnmipb.Notification) (*gnmipb.Notification, error) {
	blankRoot, err := ygot.DeepCopy(inSchema.Root)
	if err != nil {
		return nil, fmt.Errorf("failed to deep copy input schema root, %v", err)
	}
	defer func() { inSchema.Root = blankRoot }()

	if err := ytypes.UnmarshalNotifications(inSchema, []*gnmipb.Notification{notification}, nil); err != nil {
		return nil, fmt.Errorf("failed to unmarshal notifications with input schema: %v", err)
	}

	outRoot, err := ygot.DeepCopy(outSchema.Root)
	if err != nil {
		return nil, fmt.Errorf("failed to deep copy output schema root: %v", err)
	}
	for _, mapEntry := range m.mapEntries {
		nodes, err := ytypes.GetNode(inSchema.RootSchema(), inSchema.Root, mapEntry.inputWildcard, &ytypes.GetHandleWildcards{})
		if err != nil {
			// We get an error if the path doesn't exist, which is benign, so we log and continue for all errors.
			// TODO(team): Consider returning other types of errors if we can distinguish them.
			log.V(1).Infof("Entry skipped, no nodes found: %v", err)
			continue
		}

		for _, tn := range nodes {
			if tn.Data == nil {
				continue
			}

			val, err := yangValToGNMIVal(tn.Data)
			if err != nil {
				if errors.Is(err, errNilValue) {
					// Skip nil pointers inside interfaces.
					continue
				}
				return nil, fmt.Errorf("failed to convert yang val to gNMI val: %v", err)
			}

			bindings, err := bindKeys(mapEntry.input, tn.Path)
			if err != nil {
				return nil, fmt.Errorf("failed to bind keys for input path: %v", err)
			}
			outPath, err := applyBind(bindings, mapEntry.output)
			if err != nil {
				return nil, fmt.Errorf("failed to apply bindings to output path: %v", err)
			}
			if _, _, err := ytypes.GetOrCreateNode(outSchema.RootSchema(), outRoot, outPath); err != nil {
				return nil, fmt.Errorf("failed to get or create node for output path: %v", err)
			}
			if err := ytypes.SetNode(outSchema.RootSchema(), outRoot, outPath, val); err != nil {
				return nil, fmt.Errorf("failed to set node for output path: %v", err)
			}
		}
	}

	outgoingNotifications, err := ygot.TogNMINotifications(outRoot, notification.GetTimestamp(), ygot.GNMINotificationsConfig{UsePathElem: true})
	if err != nil {
		return nil, fmt.Errorf("failed to convert outgoing notification: %v", err)
	}
	if len(outgoingNotifications) != 1 {
		return nil, fmt.Errorf("received %d notifications, expected only one: %v", len(outgoingNotifications), outgoingNotifications)
	}
	if len(outgoingNotifications[0].GetUpdate()) == 0 {
		return nil, nil
	}
	out := outgoingNotifications[0]
	// If ygot calculated a prefix containing elements, flatten it into the
	// update paths. This ensures that all updates and deletes in the
	// notification are relative to a prefix that contains no elements
	// (i.e., they are full paths), which is required for deletes to be
	// handled correctly.
	if out.GetPrefix() != nil && len(out.GetPrefix().GetElem()) > 0 {
		prefixElems := out.GetPrefix().GetElem()
		for _, u := range out.GetUpdate() {
			if u.Path != nil {
				updateElems := u.Path.GetElem()
				newElems := make([]*gnmipb.PathElem, 0, len(prefixElems)+len(updateElems))
				newElems = append(newElems, prefixElems...)
				newElems = append(newElems, updateElems...)
				u.Path.Elem = newElems
			}
		}
		out.GetPrefix().Elem = nil
	}
	if out.GetPrefix() == nil {
		out.Prefix = &gnmipb.Path{}
	}
	return out, nil
}

// Handler translates gNMI notifications. This should be used as the Translate function for a functional translator.
func (m *SimpleMapper) Handler(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
	if m.deleteHandler == nil {
		m.deleteHandler = func(*gnmipb.Notification) ([]*gnmipb.Path, error) { return nil, nil }
	}
	if sr.GetUpdate() == nil {
		return nil, nil
	}
	notification := sr.GetUpdate()
	outgoingNotification, err := m.updateHandler(m.inSchema, m.outSchema, notification)
	if err != nil {
		return nil, fmt.Errorf("failed to handle updates: %v", err)
	}
	deletes, err := m.deleteHandler(notification)
	if err != nil {
		return nil, fmt.Errorf("failed to handle deletes: %v", err)
	}
	// Per gNMI spec, delete paths should not contain origin/target if
	// prefix is used.
	var cleanedDeletes []*gnmipb.Path
	for _, d := range deletes {
		// Clone the path to avoid modifying the original in-place, which might be
		// shared if deleteHandler returns paths from the original notification.
		cloned := proto.Clone(d).(*gnmipb.Path)
		cloned.Origin = ""
		cloned.Target = ""
		cleanedDeletes = append(cleanedDeletes, cloned)
	}
	deletes = cleanedDeletes

	if outgoingNotification == nil && len(deletes) == 0 {
		return nil, nil
	}

	target := notification.GetPrefix().GetTarget()
	if outgoingNotification == nil {
		outgoingNotification = &gnmipb.Notification{
			Prefix:    &gnmipb.Path{Target: target},
			Timestamp: notification.GetTimestamp(),
		}
	} else {
		outgoingNotification.Prefix.Target = target
	}
	outgoingNotification.Prefix.Origin = "openconfig"
	outgoingNotification.Delete = deletes

	outgoingNotification = ftutilities.Filter(outgoingNotification, func(path *gnmipb.Path, isDelete bool) bool {
		notificationSchema := ftutilities.GNMIPathToSchemaString(path, false)
		// If update, exact match
		if !isDelete {
			_, ok := m.outputToInputSchemaStrings[notificationSchema]
			return ok
		}

		// If delete, prefix match
		for ftOutputSchema := range m.outputToInputSchemaStrings {
			if strings.HasPrefix(ftOutputSchema, notificationSchema) {
				return true
			}
		}
		return false
	})
	if outgoingNotification == nil {
		return nil, nil
	}
	return &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: outgoingNotification,
		},
	}, nil
}
