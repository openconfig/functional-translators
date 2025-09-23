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

// Package ftutilities contains utility functions used by functional translators.
package ftutilities

import (
	"fmt"
	"maps"
	"os"
	"path"
	"strings"
	"sync"

	log "github.com/golang/glog"
	"google.golang.org/protobuf/encoding/prototext"
	"github.com/openconfig/ygot/ygot"
	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

// StripPathPrefix strips the prefix from the path.
func StripPathPrefix(s, prefix string) string {
	if s == "" {
		return ""
	}
	s = strings.TrimPrefix(s, "/")
	s = strings.TrimPrefix(s, prefix)
	return strings.TrimPrefix(s, "/")
}

// ForcePathPrefix adds the prefix to the string if it is not already present, even if the string
// is empty.
func ForcePathPrefix(s, prefix string) string {
	if prefix == "" {
		return s
	}
	prefix = strings.Trim(prefix, "/")
	s = StripPathPrefix(s, prefix)
	return fmt.Sprintf("/%s/%s", prefix, s)
}

// ConfigToState replaces "config" elements with "state" elements.
func ConfigToState(p *gnmipb.Path) *gnmipb.Path {
	returnPath := &gnmipb.Path{
		Origin: p.GetOrigin(),
		Target: p.GetTarget(),
	}
	for _, elem := range p.GetElem() {
		if elem.GetName() == "config" {
			returnPath.Elem = append(returnPath.Elem, &gnmipb.PathElem{
				Name: "state",
			})
		} else {
			returnPath.Elem = append(returnPath.Elem, elem)
		}
	}
	return returnPath
}

// Join returns a new gNMI path with the elements of p1 and p2 concatenated. The origin and target
// of p1 are used, if present, and replaced with the values from p2 otherwise.
func Join(p1, p2 *gnmipb.Path) *gnmipb.Path {
	if p1 == nil {
		return p2
	}
	origin := p1.GetOrigin()
	target := p1.GetTarget()
	if origin == "" {
		origin = p2.GetOrigin()
	}
	if target == "" {
		target = p2.GetTarget()
	}
	return &gnmipb.Path{
		Origin: origin,
		Target: target,
		Elem:   append(p1.GetElem(), p2.GetElem()...),
	}
}

// Filter returns a new notification with only the updates that return true from the provided fn.
func Filter(notification *gnmipb.Notification, fn func(path *gnmipb.Path, isDelete bool) bool) *gnmipb.Notification {
	var updates []*gnmipb.Update
	for _, update := range notification.GetUpdate() {
		if fn(Join(notification.GetPrefix(), update.GetPath()), false) {
			updates = append(updates, update)
		}
	}
	var deletes []*gnmipb.Path
	for _, delete := range notification.GetDelete() {
		if fn(Join(notification.GetPrefix(), delete), true) {
			deletes = append(deletes, delete)
		}
	}
	if len(updates) == 0 && len(deletes) == 0 {
		return nil
	}
	return &gnmipb.Notification{
		Prefix:    notification.GetPrefix(),
		Timestamp: notification.GetTimestamp(),
		Update:    updates,
		Delete:    deletes,
	}
}

// MatchPath returns true if path matches against the provided pattern.
// A wildcard character "*" in the pattern matches all path elements.
func MatchPath(path, pattern *gnmipb.Path) bool {
	wildcardMarker := "*"
	if len(path.GetElem()) != len(pattern.GetElem()) {
		return false
	}
	for ix, pathElem := range path.GetElem() {
		patternElem := pattern.GetElem()[ix]
		if patternElem.GetName() == wildcardMarker {
			continue
		}
		if pathElem.GetName() != patternElem.GetName() {
			return false
		}
	}
	return true
}

// SortByYgotString returns a function to sort gnmi paths by their stringified value.
func SortByYgotString(s []*gnmipb.Path) func(i, j int) bool {
	return func(i, j int) bool {
		l, err := ygot.PathToString(s[i])
		if err != nil {
			return false
		}
		r, err := ygot.PathToString(s[j])
		if err != nil {
			return false
		}
		return l < r
	}
}

// ValidOrigins is the set of valid origins for gNMI paths. If they occur as the first element of
// a path or as the first non-empty string in a stringified path, they are used to set the origin.
var ValidOrigins = map[string]struct{}{
	"openconfig": {},

	// Arista
	"eos_native": {},

	// Cisco XR-controller-optics-oper
	"Cisco-IOS-XR-controller-optics-oper": {},

	// Cisco XR-fabric-plane-health-oper
	"Cisco-IOS-XR-fabric-plane-health-oper": {},

	// Cisco XR-infra-statsd-oper
	"Cisco-IOS-XR-infra-statsd-oper": {},

	// Cisco XR-ipv4-arp-oper
	"Cisco-IOS-XR-ipv4-arp-oper": {},

	// Cisco XR-ipv4-io-oper
	"Cisco-IOS-XR-ipv4-io-oper": {},

	// Cisco XR-ipv6-ma-oper
	"Cisco-IOS-XR-ipv6-ma-oper": {},

	// Cisco XR-ipv6-nd-oper
	"Cisco-IOS-XR-ipv6-nd-oper": {},

	// Cisco XR-platforms-ofa-oper
	"Cisco-IOS-XR-platforms-ofa-oper": {},

	// Cisco XR-shellutil-filesystem-oper
	"Cisco-IOS-XR-shellutil-filesystem-oper": {},

	// Cisco XR-show-fpd-loc-ng-oper
	"Cisco-IOS-XR-show-fpd-loc-ng-oper": {},

	// Cisco-IOS-XR-switch-oper
	"Cisco-IOS-XR-switch-oper": {},

	// Cisco XR-qos-ma-oper
	"Cisco-IOS-XR-qos-ma-oper": {},
}

// StringToPath converts a string to a gNMI path, potentially including an origin.
// The string must be in the format "origin/path/to/element", "/origin/path/to/element",
// "path/to/element", or "/path/to/element".
// The origin is only parsed if it is in the ValidOrigins map.
func StringToPath(s string) (*gnmipb.Path, error) {
	if s == "" {
		return nil, fmt.Errorf("empty string provided")
	}
	s = strings.TrimPrefix(s, "/")
	var origin string
	var elems []*gnmipb.PathElem
	for ix, elemName := range strings.Split(s, "/") {
		if ix == 0 {
			if _, ok := ValidOrigins[elemName]; ok {
				origin = elemName
				continue
			}
		}
		elems = append(elems, &gnmipb.PathElem{
			Name: elemName,
		})
	}
	return &gnmipb.Path{
		Origin: origin,
		Elem:   elems,
	}, nil
}

// StringMapPaths converts each string in the slices, into a list of gnmi Paths.
// The lists are returned with the same keys as the input.
func StringMapPaths(stringPathMap map[string][]string) (map[string][]*gnmipb.Path, error) {
	m := make(map[string][]*gnmipb.Path)
	for k, paths := range stringPathMap {
		for _, s := range paths {
			p, err := StringToPath(s)
			if err != nil {
				return nil, fmt.Errorf("failed to convert string to path: %v", err)
			}
			m[k] = append(m[k], p)
		}
	}
	return m, nil
}

// MustStringMapPaths converts each string in the slices, into a list of gnmi Paths.
// it fails if there is an error.
func MustStringMapPaths(m map[string][]string) map[string][]*gnmipb.Path {
	p, err := StringMapPaths(m)
	if err != nil {
		log.Fatalf("map %#v cannot parse output paths into gNMI Paths", m)
	}
	return p
}

// PathInList returns True if the path is in the list of paths.
func PathInList(p *gnmipb.Path, paths []*gnmipb.Path) bool {
	for _, path := range paths {
		if MatchPath(p, path) {
			return true
		}
	}
	return false
}

// StateLeaves returns true if one of the path elements has name "state".
func StateLeaves(up *gnmipb.Update) bool {
	path := up.GetPath()
	for _, elem := range path.GetElem() {
		if elem.GetName() == "state" {
			return true
		}
	}
	return false
}

// FilterUpdates returns a slice containing updates that return true from the provided fn.
func FilterUpdates(update []*gnmipb.Update, fn func(up *gnmipb.Update) bool) []*gnmipb.Update {
	var returnUpdates []*gnmipb.Update
	for _, u := range update {
		if fn(u) {
			returnUpdates = append(returnUpdates, u)
		}
	}
	return returnUpdates
}

// GNMIPathToSchemaString converts a gNMI path to a string.
func GNMIPathToSchemaString(p *gnmipb.Path, setOCIfOriginMissing bool) string {
	s := GNMIPathToSchemaStrings(p, setOCIfOriginMissing)
	return "/" + path.Join(s...)
}

func pathNilOrEmpty(p *gnmipb.Path) bool {
	if p == nil {
		return true
	}
	if len(p.GetElem()) > 0 {
		return false
	}
	if p.GetOrigin() != "" {
		return false
	}
	if p.GetTarget() != "" {
		return false
	}
	return true
}

// GNMIPathToSchemaStrings extracts the schema path and converts it into a slice of strings.
// Origin is prepended as the first element, if present.
// setOCOriginIfMissing specifies if that path should include the origin. This is necessary because
// Arista does not set the origin, and we would have to forcibly set the origin as openConfig.
func GNMIPathToSchemaStrings(path *gnmipb.Path, setOCIfOriginMissing bool) []string {
	if pathNilOrEmpty(path) {
		return nil
	}
	var p []string
	if path.GetOrigin() != "" {
		p = append(p, path.GetOrigin())
	} else if setOCIfOriginMissing {
		// Happens on every message for Arista/Nokia.
		p = append(p, "openconfig")
	}
	for _, e := range path.Elem {
		p = append(p, e.Name)
	}
	return p
}

// FilterStructToState converts a ygot struct to a gNMI subscribe response.
func FilterStructToState(s ygot.GoStruct, ts int64, origin, target string) (*gnmipb.SubscribeResponse, error) {
	outgoingNotifications, err := ygot.TogNMINotifications(s, ts, ygot.GNMINotificationsConfig{UsePathElem: true})
	if err != nil || len(outgoingNotifications) != 1 {
		return nil, fmt.Errorf("failed to convert outgoing notification: %v", err)
	}
	n := outgoingNotifications[0]
	n.Update = FilterUpdates(n.GetUpdate(), StateLeaves)
	if len(n.GetUpdate()) == 0 {
		return nil, nil
	}
	sr := &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{
			Update: n,
		},
	}

	if n.GetPrefix() == nil {
		n.Prefix = &gnmipb.Path{}
	}
	n.Prefix.Origin = origin
	n.Prefix.Target = target
	return sr, nil
}

// MaybeConvertOptical returns the modified port name based on the optics type. Breakout child
// interfaces are ignored, as telemetry is provided through the parent interface.
// This is used by CISCOXR WBB devices when using the native path, which is of the form
// "Optics0/0/0/0" and the openconfig path is of the form "HundredGigE0/0/0/0", etc.
func MaybeConvertOptical(portName string, opticsType string) (newPortName string, wanted bool) {
	portSplits := strings.Split(portName, "/")
	if len(portSplits) != 4 {
		return portName, false
	}
	prefix := "Optics"
	opticsType = strings.ToLower(opticsType)
	switch {
	case strings.HasPrefix(opticsType, "100g"):
		prefix = "HundredGigE"
	case strings.HasPrefix(opticsType, "10g"):
		prefix = "TenGigE"
	case strings.HasPrefix(opticsType, "2x100g"):
		prefix = "TwoHundredGigE"
	case strings.HasPrefix(opticsType, "400g"):
		prefix = "FourHundredGigE"
	case strings.HasPrefix(opticsType, "40g"):
		prefix = "FortyGigE"
	case strings.HasPrefix(opticsType, "4x100g"):
		prefix = "FourHundredGigE"
	case strings.HasPrefix(opticsType, "4x10g"):
		prefix = "FortyGigE"
	}
	return strings.Replace(portName, "Optics", prefix, 1), true
}

// LoadSubscribeResponse loads a subscribe response from a file.
func LoadSubscribeResponse(path string) (*gnmipb.SubscribeResponse, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}
	sr := &gnmipb.SubscribeResponse{}
	if err := prototext.Unmarshal(b, sr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal SubscribeResponse: %v", err)
	}
	return sr, nil
}

// InterfaceMacSecInfo holds MACsec status information for a specific interface.
type InterfaceMacSecInfo struct {
	mu            sync.Mutex
	interfaceName string

	cpStatus    bool
	cpStatusSet bool

	// Stores principal and success status per CKN.
	cknStatuses map[string]*CKNInfo // map[CKN_string]*CKNInfo
}

// CKNInfo holds principal and success status for a specific CKN.
type CKNInfo struct {
	principal    bool
	success      bool
	principalSet bool
	successSet   bool
}

// CreateOrGetCKN returns the CKNInfo for the given CKN, creating it if it doesn't exist.
func (i *InterfaceMacSecInfo) CreateOrGetCKN(ckn string) *CKNInfo {
	if i.cknStatuses == nil {
		i.cknStatuses = make(map[string]*CKNInfo)
	}
	if _, ok := i.cknStatuses[ckn]; !ok {
		i.cknStatuses[ckn] = new(CKNInfo)
	}
	return i.cknStatuses[ckn]
}

// SetIntfCPStatus sets the cpStatus and marks it as set.
func (i *InterfaceMacSecInfo) SetIntfCPStatus(b bool) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.cpStatus = b
	i.cpStatusSet = true
}

// IntfCPStatus returns the cpStatus and a boolean indicating if it has been set.
func (i *InterfaceMacSecInfo) IntfCPStatus() (bool, bool) {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.cpStatus, i.cpStatusSet
}

// ResetCPStatus marks the cpStatus as not set and resets its value.
// This is used when the native source for cpStatus is deleted.
func (i *InterfaceMacSecInfo) ResetCPStatus() {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.cpStatus = false // Reset to a default value
	i.cpStatusSet = false
}

// SetIntfPrincipal sets the principal status for a given CKN and marks it as set.
func (i *InterfaceMacSecInfo) SetIntfPrincipal(ckn string, b bool) {
	cknInfo := i.CreateOrGetCKN(ckn)
	cknInfo.principal = b
	cknInfo.principalSet = true
}

// IntfPrincipal returns the principal status for a given CKN and a boolean indicating if it has been set.
func (i *InterfaceMacSecInfo) IntfPrincipal(ckn string) (bool, bool) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if cknInfo, ok := i.cknStatuses[ckn]; ok && cknInfo.principalSet {
		return cknInfo.principal, true
	}
	return false, false
}

// SetIntfSuccess sets the success status for a given CKN and marks it as set.
func (i *InterfaceMacSecInfo) SetIntfSuccess(ckn string, b bool) {
	i.mu.Lock()
	defer i.mu.Unlock()
	cknInfo := i.CreateOrGetCKN(ckn)
	cknInfo.success = b
	cknInfo.successSet = true
}

// IntfSuccess returns the success status for a given CKN and a boolean indicating if it has been set.
func (i *InterfaceMacSecInfo) IntfSuccess(ckn string) (bool, bool) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if cknInfo, ok := i.cknStatuses[ckn]; ok && cknInfo.successSet {
		return cknInfo.success, true
	}
	return false, false
}

// CloneStatuses returns a copy of the CKN statuses map.
func (i *InterfaceMacSecInfo) CloneStatuses() map[string]*CKNInfo {
	i.mu.Lock()
	defer i.mu.Unlock()
	return maps.Clone(i.cknStatuses)
}

// RemoveCkn removes MACsec information for a specific CKN.
func (i *InterfaceMacSecInfo) RemoveCkn(ckn string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if i.cknStatuses != nil {
		delete(i.cknStatuses, ckn)
	}
}

// IsComplete checks if all necessary MACsec CKN status values have been set.
func (i *InterfaceMacSecInfo) IsComplete(ckn string) bool {
	i.mu.Lock()
	defer i.mu.Unlock()
	if cknInfo, ok := i.cknStatuses[ckn]; ok && cknInfo.principalSet && cknInfo.successSet {
		return true
	}
	return false
}

// TargetMacSecInfo holds MACsec information for all interfaces on a target.
type TargetMacSecInfo struct {
	mu             sync.Mutex
	TargetHostname string
	Interfaces     map[string]*InterfaceMacSecInfo // map[InterfaceName]*InterfaceMacSecInfo
}

// NewTargetMacSecInfo creates a new TargetMacSecInfo for the given target hostname.
func NewTargetMacSecInfo(targetHostname string) *TargetMacSecInfo {
	return &TargetMacSecInfo{
		TargetHostname: targetHostname,
		Interfaces:     make(map[string]*InterfaceMacSecInfo),
	}
}

// CreateOrGetInterface returns the InterfaceMacSecInfo for the given interface name, creating it if it doesn't exist.
// It also initializes the CKN statuses map if it doesn't exist.
func (t *TargetMacSecInfo) CreateOrGetInterface(interfaceName string) *InterfaceMacSecInfo {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.Interfaces[interfaceName]; !ok {
		t.Interfaces[interfaceName] = &InterfaceMacSecInfo{
			interfaceName: interfaceName,
			cknStatuses:   make(map[string]*CKNInfo),
		}
	}
	return t.Interfaces[interfaceName]
}

// InterfaceInfo retrieves the MACsec info for a specific interface.
func (t *TargetMacSecInfo) InterfaceInfo(intf string) (*InterfaceMacSecInfo, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	info, ok := t.Interfaces[intf]
	return info, ok
}

// ClearInterfaceInfo removes MACsec information for a specific interface.
func (t *TargetMacSecInfo) ClearInterfaceInfo(intf string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.Interfaces, intf)
}

// AristaMACSecMapCache is a thread-safe cache for AristaMACSecMap.
// It stores cached boolean values from distinct native Arista MACsec paths per target/interface/CKN.
// Although Functional Translators (FTs) are typically stateless, this map is required as an exception
// to hold values from these multiple source paths, necessary for deriving the single OpenConfig MACsec status.
// Declaring it here allows access by both the FT logic and the FT registration process,
// where it is cleared to prevent using stale information between registrations or updates.
type AristaMACSecMapCache struct {
	mu   sync.Mutex
	data map[string]*TargetMacSecInfo
}

// Global instance of the AristaMACSecMapCache.
var (
	AristaMACSecMap = &AristaMACSecMapCache{
		data: make(map[string]*TargetMacSecInfo),
	}
)

// SetTargetMacSecInfo adds or updates the TargetMacSecInfo for a given target hostname.
func (c *AristaMACSecMapCache) SetTargetMacSecInfo(targetHostname string, info *TargetMacSecInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[targetHostname] = info
}

// RetrieveTargetMacSecInfo fetches the TargetMacSecInfo for a given target hostname.
// It returns the info and a boolean indicating if the target was found.
func (c *AristaMACSecMapCache) RetrieveTargetMacSecInfo(targetHostname string) (*TargetMacSecInfo, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	info, ok := c.data[targetHostname]
	return info, ok
}

// DeleteTargetMacSecInfo removes the TargetMacSecInfo for a given target hostname.
func (c *AristaMACSecMapCache) DeleteTargetMacSecInfo(targetHostname string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, targetHostname)
}

// ClearAllTargetMacSecInfo removes all entries from the cache.
func (c *AristaMACSecMapCache) ClearAllTargetMacSecInfo() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]*TargetMacSecInfo)
}

// CreateOrUpdateTargetMacSecInfo retrieves an existing TargetMacSecInfo for the given target
// or creates a new one if it doesn't exist, then stores it in the cache.
func (c *AristaMACSecMapCache) CreateOrUpdateTargetMacSecInfo(targetHostname string) *TargetMacSecInfo {
	c.mu.Lock()
	defer c.mu.Unlock()
	info, ok := c.data[targetHostname]
	if !ok {
		info = NewTargetMacSecInfo(targetHostname)
		c.data[targetHostname] = info
	}
	return info
}
