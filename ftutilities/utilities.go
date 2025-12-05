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

	// Cisco XR-envmon-oper
	"Cisco-IOS-XR-envmon-oper": {},
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

// stringMapPaths converts each string in the slices, into a list of gnmi Paths.
// The lists are returned with the same keys as the input.
func stringMapPaths(stringPathMap map[string][]string) (map[string][]*gnmipb.Path, error) {
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
	p, err := stringMapPaths(m)
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

// TargetQoSInfo holds QoS information for all port-channels on a target.
type TargetQoSInfo struct {
	mu                  sync.Mutex
	TargetHostname      string
	PortChannels        map[string]*PortChannelInfo     // map[PortChannelName]*PortChannelInfo
	MemberToPCMap       map[string]string               // map[InterfaceName]PortChannelName
	UnassociatedMembers map[string]*MemberInterfaceInfo // "Waiting room"
}

// PortChannelInfo holds QoS information for a specific port-channel,
// including all its member interfaces.
type PortChannelInfo struct {
	mu              sync.Mutex
	portChannelName string
	Members         map[string]*MemberInterfaceInfo // map[InterfaceName]*MemberInterfaceInfo
}

// MemberInterfaceInfo holds QoS queue information for a specific member interface.
type MemberInterfaceInfo struct {
	mu            sync.Mutex
	interfaceName string
	Queues        map[string]*QueueCounters // map[QueueID]*QueueCounters
}

// NewMemberInterfaceInfo creates a new MemberInterfaceInfo instance.
func NewMemberInterfaceInfo(interfaceName string) *MemberInterfaceInfo {
	return &MemberInterfaceInfo{
		interfaceName: interfaceName,
		Queues:        make(map[string]*QueueCounters),
	}
}

// QueueCounters holds the 4 counters for a specific QoS queue.
// It also tracks whether each counter has been explicitly set.
type QueueCounters struct {
	TxPackets      uint64
	TxBytes        uint64
	DroppedPackets uint64
	DroppedBytes   uint64

	TxPacketsSet      bool
	TxBytesSet        bool
	DroppedPacketsSet bool
	DroppedBytesSet   bool
}

// QoSAggregationMapCache is a thread-safe cache for TargetQoSInfo.
// It stores cached QoS counter values from distinct OC paths
// per target/port-channel/interface/queue.
type QoSAggregationMapCache struct {
	mu   sync.Mutex
	data map[string]*TargetQoSInfo // map[TargetHostname]*TargetQoSInfo
}

// QoSAggMap is the global instance of the QoSAggregationMapCache.
var (
	QoSAggMap = &QoSAggregationMapCache{
		data: make(map[string]*TargetQoSInfo),
	}
)

// createOrGetQueue is an internal helper that assumes the lock is held.
func (m *MemberInterfaceInfo) createOrGetQueue(queueID string) *QueueCounters {
	if m.Queues == nil {
		m.Queues = make(map[string]*QueueCounters)
	}
	if _, ok := m.Queues[queueID]; !ok {
		m.Queues[queueID] = new(QueueCounters)
	}
	return m.Queues[queueID]
}

// SetTxPackets sets the TxPackets counter for a given queue and marks it as set.
func (m *MemberInterfaceInfo) SetTxPackets(queueID string, val uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	q := m.createOrGetQueue(queueID)
	q.TxPackets = val
	q.TxPacketsSet = true
}

// SetTxBytes sets the TxBytes counter for a given queue and marks it as set.
func (m *MemberInterfaceInfo) SetTxBytes(queueID string, val uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	q := m.createOrGetQueue(queueID)
	q.TxBytes = val
	q.TxBytesSet = true
}

// SetDroppedPackets sets the DroppedPackets counter for a given queue and marks it as set.
func (m *MemberInterfaceInfo) SetDroppedPackets(queueID string, val uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	q := m.createOrGetQueue(queueID)
	q.DroppedPackets = val
	q.DroppedPacketsSet = true
}

// SetDroppedBytes sets the DroppedBytes counter for a given queue and marks it as set.
func (m *MemberInterfaceInfo) SetDroppedBytes(queueID string, val uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	q := m.createOrGetQueue(queueID)
	q.DroppedBytes = val
	q.DroppedBytesSet = true
}

// CloneQueues returns a copy of the Queues map.
func (m *MemberInterfaceInfo) CloneQueues() map[string]*QueueCounters {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Clone the map to avoid concurrent access issues during iteration
	clonedMap := make(map[string]*QueueCounters, len(m.Queues))
	for k, v := range m.Queues {
		// Clone the counter struct itself
		if v != nil {
			c := *v
			clonedMap[k] = &c
		}
	}
	return clonedMap
}

// --- PortChannelInfo Methods ---

// CreateOrRetrieveMember returns the MemberInterfaceInfo for the given interface name,
// creating it if it doesn't exist.
func (p *PortChannelInfo) CreateOrRetrieveMember(interfaceName string) *MemberInterfaceInfo {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.Members == nil {
		p.Members = make(map[string]*MemberInterfaceInfo)
	}
	if _, ok := p.Members[interfaceName]; !ok {
		p.Members[interfaceName] = &MemberInterfaceInfo{
			interfaceName: interfaceName,
			Queues:        make(map[string]*QueueCounters),
		}
	}
	return p.Members[interfaceName]
}

// AddMemberInfo adds a pre-existing MemberInterfaceInfo object to the PortChannel.
// This is used to move a member from the "waiting room" into the PortChannel.
func (p *PortChannelInfo) AddMemberInfo(memberInfo *MemberInterfaceInfo) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.Members == nil {
		p.Members = make(map[string]*MemberInterfaceInfo)
	}
	p.Members[memberInfo.interfaceName] = memberInfo
}

// ClearMemberInfo removes QoS information for a specific member interface.
func (p *PortChannelInfo) ClearMemberInfo(intf string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.Members, intf)
}

// CloneMembers returns a thread-safe copy of the Members map.
func (p *PortChannelInfo) CloneMembers() map[string]*MemberInterfaceInfo {
	p.mu.Lock()
	defer p.mu.Unlock()
	// Clone the map to avoid race conditions during iteration.
	clonedMap := make(map[string]*MemberInterfaceInfo, len(p.Members))
	for k, v := range p.Members {
		clonedMap[k] = v
	}
	return clonedMap
}

// --- TargetQoSInfo Methods ---

// newTargetQoSInfo creates a new TargetQoSInfo for the given target hostname.
func newTargetQoSInfo(targetHostname string) *TargetQoSInfo {
	return &TargetQoSInfo{
		TargetHostname:      targetHostname,
		PortChannels:        make(map[string]*PortChannelInfo),
		MemberToPCMap:       make(map[string]string),
		UnassociatedMembers: make(map[string]*MemberInterfaceInfo),
	}
}

// CreateOrRetrievePortChannel returns the PortChannelInfo for the given port-channel name,
// creating it if it doesn't exist.
func (t *TargetQoSInfo) CreateOrRetrievePortChannel(pcName string) *PortChannelInfo {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.PortChannels == nil {
		t.PortChannels = make(map[string]*PortChannelInfo)
	}
	if _, ok := t.PortChannels[pcName]; !ok {
		t.PortChannels[pcName] = &PortChannelInfo{
			portChannelName: pcName,
			Members:         make(map[string]*MemberInterfaceInfo),
		}
	}
	return t.PortChannels[pcName]
}

// PortChannelInfo retrieves the QoS info for a specific port-channel.
func (t *TargetQoSInfo) PortChannelInfo(pcName string) (*PortChannelInfo, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	info, ok := t.PortChannels[pcName]
	return info, ok
}

// --- NEW Reverse Map Methods for TargetQoSInfo ---

// RetrievePortChannelForMember returns the parent Port-Channel name for a member.
func (t *TargetQoSInfo) RetrievePortChannelForMember(memberName string) (string, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	pcName, ok := t.MemberToPCMap[memberName]
	return pcName, ok
}

// SetPortChannelForMember sets the parent Port-Channel for a member.
func (t *TargetQoSInfo) SetPortChannelForMember(memberName, pcName string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.MemberToPCMap[memberName] = pcName
}

// FindAndRemoveMember finds a member, removes it from its old PortChannel's
// member list, and removes it from the reverse map.
// It returns the old Port-Channel name and true if found.
func (t *TargetQoSInfo) FindAndRemoveMember(memberName string) (string, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	oldPCName, ok := t.MemberToPCMap[memberName]
	if !ok {
		return "", false
	}

	// Remove from reverse map
	delete(t.MemberToPCMap, memberName)

	// Remove from the PortChannel's member list
	if pcInfo, pcOK := t.PortChannels[oldPCName]; pcOK {
		pcInfo.ClearMemberInfo(memberName) // This locks pcInfo.mu
	}
	return oldPCName, true
}

// --- QoSAggregationMapCache Methods ---

// RetrieveTargetQoSInfo fetches the TargetQoSInfo for a given target hostname.
// It returns the info and a boolean indicating if the target was found.
func (c *QoSAggregationMapCache) RetrieveTargetQoSInfo(targetHostname string) (*TargetQoSInfo, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	info, ok := c.data[targetHostname]
	return info, ok
}

// ClearAllTargetQoSInfo removes all entries from the cache.
func (c *QoSAggregationMapCache) ClearAllTargetQoSInfo() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]*TargetQoSInfo)
}

// CreateOrUpdateTargetQoSInfo retrieves an existing TargetQoSInfo for the given target
// or creates a new one if it doesn't exist, then stores it in the cache.
func (c *QoSAggregationMapCache) CreateOrUpdateTargetQoSInfo(targetHostname string) *TargetQoSInfo {
	c.mu.Lock()
	defer c.mu.Unlock()
	info, ok := c.data[targetHostname]
	if !ok {
		info = newTargetQoSInfo(targetHostname)
		c.data[targetHostname] = info
	}
	return info
}
