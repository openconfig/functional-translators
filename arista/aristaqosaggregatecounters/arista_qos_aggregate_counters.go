// Package aristaqosaggregatecounters implements a functional translator for Arista QOS aggregate counters.
package aristaqosaggregatecounters

import (
	"fmt"
	"strings"

	log "github.com/golang/glog"
	"github.com/openconfig/functional-translators/ftconsts"
	"github.com/openconfig/functional-translators/ftutilities"
	"github.com/openconfig/functional-translators/translator"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

var (
	translateMap = map[string][]string{
		"/openconfig/qos/interfaces/interface/output/queues/queue/state/transmit-octets": {
			"/openconfig/interfaces/interface/ethernet/state/aggregate-id",
			"/openconfig/qos/interfaces/interface/output/queues/queue/state/transmit-octets",
		},
		"/openconfig/qos/interfaces/interface/output/queues/queue/state/transmit-pkts": {
			"/openconfig/interfaces/interface/ethernet/state/aggregate-id",
			"/openconfig/qos/interfaces/interface/output/queues/queue/state/transmit-pkts",
		},
		"/openconfig/qos/interfaces/interface/output/queues/queue/state/dropped-octets": {
			"/openconfig/interfaces/interface/ethernet/state/aggregate-id",
			"/openconfig/qos/interfaces/interface/output/queues/queue/state/dropped-octets",
		},
		"/openconfig/qos/interfaces/interface/output/queues/queue/state/dropped-pkts": {
			"/openconfig/interfaces/interface/ethernet/state/aggregate-id",
			"/openconfig/qos/interfaces/interface/output/queues/queue/state/dropped-pkts",
		},
	}
	updatePathPatterns = []*gnmipb.Path{
		{
			Elem: []*gnmipb.PathElem{
				{Name: "interfaces"},
				{Name: "interface", Key: map[string]string{"name": "*"}},
				{Name: "ethernet"}, {Name: "state"}, {Name: "aggregate-id"},
			},
		},
		{
			Elem: []*gnmipb.PathElem{
				{Name: "qos"}, {Name: "interfaces"},
				{Name: "interface", Key: map[string]string{"interface-id": "*"}},
				{Name: "output"}, {Name: "queues"},
				{Name: "queue", Key: map[string]string{"name": "*"}},
				{Name: "state"}, {Name: "transmit-octets"},
			},
		},
		{
			Elem: []*gnmipb.PathElem{
				{Name: "qos"}, {Name: "interfaces"},
				{Name: "interface", Key: map[string]string{"interface-id": "*"}},
				{Name: "output"}, {Name: "queues"},
				{Name: "queue", Key: map[string]string{"name": "*"}},
				{Name: "state"}, {Name: "transmit-pkts"},
			},
		},
		{
			Elem: []*gnmipb.PathElem{
				{Name: "qos"}, {Name: "interfaces"},
				{Name: "interface", Key: map[string]string{"interface-id": "*"}},
				{Name: "output"}, {Name: "queues"},
				{Name: "queue", Key: map[string]string{"name": "*"}},
				{Name: "state"}, {Name: "dropped-octets"},
			},
		},
		{
			Elem: []*gnmipb.PathElem{
				{Name: "qos"}, {Name: "interfaces"},
				{Name: "interface", Key: map[string]string{"interface-id": "*"}},
				{Name: "output"}, {Name: "queues"},
				{Name: "queue", Key: map[string]string{"name": "*"}},
				{Name: "state"}, {Name: "dropped-pkts"},
			},
		},
	}
	deletePathPatterns = []*gnmipb.Path{
		{
			Elem: []*gnmipb.PathElem{
				{Name: "interfaces"},
				{Name: "interface", Key: map[string]string{"name": "*"}},
				{Name: "ethernet"}, {Name: "state"}, {Name: "aggregate-id"},
			},
		},
	}

	pathElemQoS            = &gnmipb.PathElem{Name: "qos"}
	pathElemInterfaces     = &gnmipb.PathElem{Name: "interfaces"}
	pathElemOutput         = &gnmipb.PathElem{Name: "output"}
	pathElemQueues         = &gnmipb.PathElem{Name: "queues"}
	pathElemState          = &gnmipb.PathElem{Name: "state"}
	pathElemTransmitOctets = &gnmipb.PathElem{Name: leafTransmitOctets}
	pathElemTransmitPkts   = &gnmipb.PathElem{Name: leafTransmitPkts}
	pathElemDroppedOctets  = &gnmipb.PathElem{Name: leafDroppedOctets}
	pathElemDroppedPkts    = &gnmipb.PathElem{Name: leafDroppedPkts}
	qosInterfacesPrefix    = []*gnmipb.PathElem{pathElemQoS, pathElemInterfaces}
)

const (
	leafAggregateID    = "aggregate-id"
	leafTransmitOctets = "transmit-octets"
	leafTransmitPkts   = "transmit-pkts"
	leafDroppedOctets  = "dropped-octets"
	leafDroppedPkts    = "dropped-pkts"
)

// queueCounters holds the counters for a specific QoS queue.
type queueCounters struct {
	txPackets      uint64
	txBytes        uint64
	droppedPackets uint64
	droppedBytes   uint64
}

// memberInterfaceInfo holds QoS queue information for a specific member interface.
type memberInterfaceInfo struct {
	interfaceName string
	queues        map[string]*queueCounters // map[queueID]*queueCounters
}

func newMemberInterfaceInfo(interfaceName string) *memberInterfaceInfo {
	return &memberInterfaceInfo{
		interfaceName: interfaceName,
		queues:        make(map[string]*queueCounters),
	}
}

func (m *memberInterfaceInfo) createOrGetQueue(queueID string) *queueCounters {
	if m.queues == nil {
		m.queues = make(map[string]*queueCounters)
	}
	if _, ok := m.queues[queueID]; !ok {
		m.queues[queueID] = new(queueCounters)
	}
	return m.queues[queueID]
}

// portChannelInfo holds QoS information for a specific port-channel, including all its member interfaces.
type portChannelInfo struct {
	portChannelName    string
	members            map[string]*memberInterfaceInfo // map[interfaceName]*memberInterfaceInfo
	aggregatedCounters map[string]queueCounters        // Reusable map for zero-allocation aggregation: map[queueID]queueCounters
}

func (p *portChannelInfo) createOrRetrieveMember(interfaceName string) *memberInterfaceInfo {
	if p.members == nil {
		p.members = make(map[string]*memberInterfaceInfo)
	}
	if _, ok := p.members[interfaceName]; !ok {
		p.members[interfaceName] = newMemberInterfaceInfo(interfaceName)
	}
	return p.members[interfaceName]
}

func (p *portChannelInfo) addMemberInfo(memberInfo *memberInterfaceInfo) {
	if p.members == nil {
		p.members = make(map[string]*memberInterfaceInfo)
	}
	p.members[memberInfo.interfaceName] = memberInfo
}

func (p *portChannelInfo) clearMemberInfo(intf string) {
	delete(p.members, intf)
}

func (p *portChannelInfo) aggregateCounters() map[string]queueCounters {
	if p.aggregatedCounters == nil {
		p.aggregatedCounters = make(map[string]queueCounters)
	} else {
		clear(p.aggregatedCounters)
	}
	for _, memberInfo := range p.members {
		for queueID, counters := range memberInfo.queues {
			agg := p.aggregatedCounters[queueID]
			agg.txBytes += counters.txBytes
			agg.txPackets += counters.txPackets
			agg.droppedBytes += counters.droppedBytes
			agg.droppedPackets += counters.droppedPackets
			p.aggregatedCounters[queueID] = agg
		}
	}
	return p.aggregatedCounters
}

// qosState holds QoS information for all port-channels on a target.
type qosState struct {
	portChannels        map[string]*portChannelInfo     // map[portChannelName]*portChannelInfo
	memberToPCMap       map[string]string               // map[interfaceName]portChannelName
	unassociatedMembers map[string]*memberInterfaceInfo // "Waiting room"
}

func newQoSState() *qosState {
	return &qosState{
		portChannels:        make(map[string]*portChannelInfo),
		memberToPCMap:       make(map[string]string),
		unassociatedMembers: make(map[string]*memberInterfaceInfo),
	}
}

func (s *qosState) createOrRetrievePortChannel(pcName string) *portChannelInfo {
	if s.portChannels == nil {
		s.portChannels = make(map[string]*portChannelInfo)
	}
	if _, ok := s.portChannels[pcName]; !ok {
		s.portChannels[pcName] = &portChannelInfo{
			portChannelName:    pcName,
			members:            make(map[string]*memberInterfaceInfo),
			aggregatedCounters: make(map[string]queueCounters),
		}
	}
	return s.portChannels[pcName]
}

func (s *qosState) portChannelInfo(pcName string) (*portChannelInfo, bool) {
	info, ok := s.portChannels[pcName]
	return info, ok
}

func (s *qosState) retrievePortChannelForMember(memberName string) (string, bool) {
	pcName, ok := s.memberToPCMap[memberName]
	return pcName, ok
}

func (s *qosState) setPortChannelForMember(memberName, pcName string) {
	s.memberToPCMap[memberName] = pcName
}

func (s *qosState) findAndRemoveMember(memberName string) (string, bool) {
	oldPCName, ok := s.memberToPCMap[memberName]
	if !ok {
		return "", false
	}
	delete(s.memberToPCMap, memberName)
	if pcInfo, pcOK := s.portChannels[oldPCName]; pcOK {
		pcInfo.clearMemberInfo(memberName)
	}
	return oldPCName, true
}

type impl struct {
	state *qosState
}

// New creates a functional translator.
func New() *translator.FunctionalTranslator {
	i := &impl{
		state: newQoSState(),
	}
	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               ftconsts.AristaQoSAggregateCountersTranslator,
			Translate:        i.translate,
			OutputToInputMap: ftutilities.MustStringMapPaths(translateMap),
			Metadata: []*translator.FTMetadata{
				{
					Vendor: ftconsts.VendorArista,
					SoftwareVersionRange: &translator.SWRange{
						InclusiveMin: "4.33.0F",
						ExclusiveMax: "4.37",
					},
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("failed to create Arista QOS aggregate counters functional translator: %v", err)
	}
	return ft
}

// parsePath extracts key information from the matched path.
func parsePath(path *gnmipb.Path) (interfaceName, simpleQueueName, leafName string, err error) {
	if path == nil || len(path.GetElem()) == 0 {
		return "", "", "", fmt.Errorf("path is nil or empty")
	}

	elems := path.GetElem()
	leafName = elems[len(elems)-1].GetName()

	// Use the leaf name to segregate the parsing logic.
	switch leafName {
	case leafAggregateID:
		// Handle only aggregate-id paths
		for _, elem := range elems {
			if elem.GetName() == "interface" {
				if id, ok := elem.GetKey()["name"]; ok {
					interfaceName = id
				}
			}
		}
		if interfaceName == "" {
			return "", "", "", fmt.Errorf("could not find key 'name' for aggregate-id path: %v", path)
		}
		return interfaceName, "", leafName, nil

	case leafTransmitOctets, leafTransmitPkts, leafDroppedOctets, leafDroppedPkts:
		// Handle only QoS counter paths
		// compositeQueueID format: <interface-id>-<queue-id> (e.g., "Ethernet33/5-0", "Ethernet33/5-MC-0")
		var compositeQueueID string
		for _, elem := range elems {
			switch elem.GetName() {
			case "interface":
				if id, ok := elem.GetKey()["interface-id"]; ok {
					interfaceName = id
				}
			case "queue":
				if name, ok := elem.GetKey()["name"]; ok {
					compositeQueueID = name
				}
			}
		}

		if interfaceName == "" {
			return "", "", "", fmt.Errorf("could not find key 'interface-id' for QoS path: %v", path)
		}
		if compositeQueueID == "" {
			return "", "", "", fmt.Errorf("could not find key 'name' for queue in QoS path: %v", path)
		}

		// Derive the simple queue name (0,1,2... MC-0, MC-1... etc) from the composite ID.
		expectedPrefix := interfaceName + "-"
		if !strings.HasPrefix(compositeQueueID, expectedPrefix) {
			return "", "", "", fmt.Errorf("queue ID %q does not start with the expected interface prefix %q", compositeQueueID, expectedPrefix)
		}
		return interfaceName, strings.TrimPrefix(compositeQueueID, expectedPrefix), leafName, nil

	default:
		// Path was matched by the patterns but has an unrecognized leaf.
		return "", "", "", fmt.Errorf("unrecognized leaf %q for a matched QoS aggregation path", leafName)
	}
}

// newCounterUpdate creates a gNMI update relative to the notification prefix
// (/qos/interfaces) for a given interface PathElem, queue PathElem, and counter
// leaf PathElem.
//
// Only /qos/interfaces is factored into Notification.Prefix because a single
// outgoing Notification combines updates across multiple interfaces (both
// member-port passthrough updates and Port-Channel aggregate updates) and
// multiple queues.
func newCounterUpdate(interfaceElem, queueElem, leafElem *gnmipb.PathElem, value uint64) *gnmipb.Update {
	return &gnmipb.Update{
		Path: &gnmipb.Path{
			Elem: []*gnmipb.PathElem{
				interfaceElem,
				pathElemOutput,
				pathElemQueues,
				queueElem,
				pathElemState,
				leafElem,
			},
		},
		Val: &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: value}},
	}
}

// aggregateAndBuildUpdates calculates the sum of counters for a port-channel and creates gNMI updates.
func (i *impl) aggregateAndBuildUpdates(pcName string) []*gnmipb.Update {
	pcInfo, ok := i.state.portChannelInfo(pcName)
	if !ok {
		return nil
	}

	// aggregatedCounters maps a simple queue name to its summed counters.
	aggregatedCounters := pcInfo.aggregateCounters()
	if len(aggregatedCounters) == 0 {
		return nil
	}
	interfaceElem := &gnmipb.PathElem{
		Name: "interface",
		Key: map[string]string{
			"interface-id": pcName,
		},
	}
	outgoingUpdates := make([]*gnmipb.Update, 0, len(aggregatedCounters)*4)
	// Build gNMI Updates from the aggregated results.
	for simpleQueueName, counters := range aggregatedCounters {
		// Create the new composite queue ID for the aggregated path.
		// e.g., if pcName is "Port-Channel10" and simpleQueueName is "0", this becomes "Port-Channel10-0".
		queueElem := &gnmipb.PathElem{
			Name: "queue",
			Key: map[string]string{
				"name": pcName + "-" + simpleQueueName,
			},
		}

		outgoingUpdates = append(outgoingUpdates,
			newCounterUpdate(interfaceElem, queueElem, pathElemTransmitOctets, counters.txBytes),
			newCounterUpdate(interfaceElem, queueElem, pathElemTransmitPkts, counters.txPackets),
			newCounterUpdate(interfaceElem, queueElem, pathElemDroppedOctets, counters.droppedBytes),
			newCounterUpdate(interfaceElem, queueElem, pathElemDroppedPkts, counters.droppedPackets),
		)
	}
	return outgoingUpdates
}

// deleteHandler processes delete notifications to handle port-channel member removals.
// It updates the cache and returns a map of port-channel names that were affected.
func (i *impl) deleteHandler(n *gnmipb.Notification) (impactedPortChannels map[string]bool) {
	impactedPortChannels = make(map[string]bool)
	prefix := n.GetPrefix()

	for _, delPath := range n.GetDelete() {
		fullPath := ftutilities.Join(prefix, delPath)
		interfaceName, _, _, err := parsePath(fullPath)
		if err != nil {
			log.Warningf("failed to parse delete path %v: %v", fullPath, err)
			continue
		}

		if oldPCName, removed := i.processMemberRemoval(interfaceName); removed {
			impactedPortChannels[oldPCName] = true
		}
	}
	return impactedPortChannels
}

// processMemberRemoval updates the cache and returns the name of the impacted port-channel and whether a removal occurred.
func (i *impl) processMemberRemoval(interfaceName string) (string, bool) {
	// Also clean up the member if it's in the "waiting room"
	if _, found := i.state.unassociatedMembers[interfaceName]; found {
		delete(i.state.unassociatedMembers, interfaceName)
		log.V(2).Infof("removed unassociated member %q from the waiting room.", interfaceName)
		return "", true
	}

	// findAndRemoveMember cleans up the cache and returns the old port-channel.
	if oldPCName, found := i.state.findAndRemoveMember(interfaceName); found {
		log.V(2).Infof("processed removal of member %q from Port-Channel %q.", interfaceName, oldPCName)
		return oldPCName, true
	}
	return "", false
}

// handleAggregateIDUpdate processes membership changes based on an aggregate-id update.
// It updates the cache and the set of impacted port-channels.
func (i *impl) handleAggregateIDUpdate(interfaceName, newPCName string, impactedPortChannels map[string]bool) {
	// Process the implicit removal from any old port-channel.
	if oldPCName, removed := i.state.findAndRemoveMember(interfaceName); removed {
		impactedPortChannels[oldPCName] = true
	}

	// If newPCName is empty, it signifies a removal.
	// We have no new port-channel to add to, so we can return early.
	if newPCName == "" {
		return
	}

	// Handle the addition to the new port-channel if one is specified.
	pcInfo := i.state.createOrRetrievePortChannel(newPCName)
	i.state.setPortChannelForMember(interfaceName, newPCName)

	// Check the "waiting room" for pending counters for this interface.
	if unassociatedMember, found := i.state.unassociatedMembers[interfaceName]; found {
		// Move the cached member info into the port-channel's member list.
		pcInfo.addMemberInfo(unassociatedMember)
		// Clean up the waiting room.
		delete(i.state.unassociatedMembers, interfaceName)
		log.V(2).Infof("moved pending counters for %s to Port-Channel %s", interfaceName, newPCName)
	} else {
		// If no pending data, just ensure the member struct exists.
		pcInfo.createOrRetrieveMember(interfaceName)
	}

	// Mark the port-channel as impacted to trigger an immediate aggregation.
	impactedPortChannels[newPCName] = true
	log.V(1).Infof("member %s assigned to new Port-Channel %s", interfaceName, newPCName)
}

// handleQoSUpdate processes a QoS counter update by updating the appropriate cache location.
func (i *impl) handleQoSUpdate(interfaceName, simpleQueueName, leafName string, value uint64, impactedPortChannels map[string]bool) {
	var memberInfo *memberInterfaceInfo

	// The member is already part of a known port-channel. Update its state directly.
	pcName, ok := i.state.retrievePortChannelForMember(interfaceName)
	if ok {
		impactedPortChannels[pcName] = true

		// Get the PortChannelInfo and then the MemberInterfaceInfo.
		pcInfo, pcOk := i.state.portChannelInfo(pcName)
		if !pcOk {
			log.Errorf("cache inconsistency: port-channel %s not found for member %s", pcName, interfaceName)
			return
		}
		memberInfo = pcInfo.createOrRetrieveMember(interfaceName)
	} else {
		// We don't know its LAG yet. Put its data in the "waiting room".
		var found bool
		memberInfo, found = i.state.unassociatedMembers[interfaceName]
		if !found {
			// Member is not in the waiting room yet, so create it.
			memberInfo = newMemberInterfaceInfo(interfaceName)
			i.state.unassociatedMembers[interfaceName] = memberInfo
		}
	}

	q := memberInfo.createOrGetQueue(simpleQueueName)
	switch leafName {
	case leafTransmitOctets:
		q.txBytes = value
	case leafTransmitPkts:
		q.txPackets = value
	case leafDroppedOctets:
		q.droppedBytes = value
	case leafDroppedPkts:
		q.droppedPackets = value
	}
}

func (i *impl) translate(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
	notification := sr.GetUpdate()
	if notification == nil {
		return nil, nil
	}

	prefix := notification.GetPrefix()
	target := prefix.GetTarget()
	timestamp := notification.GetTimestamp()

	// Handle deletes first and get the initial set of impacted port-channels.
	impactedPortChannels := i.deleteHandler(notification)

	// A slice to store all original QOS counter updates for singleton ports that should be passed through.
	passthroughUpdates := make([]*gnmipb.Update, 0, len(notification.GetUpdate()))

	// Handle both membership changes and counter updates.
	for _, update := range notification.GetUpdate() {
		fullPath := ftutilities.Join(prefix, update.GetPath())
		var matched bool
		for _, pattern := range updatePathPatterns {
			if ftutilities.MatchPath(fullPath, pattern) {
				matched = true
				break
			}
		}

		if !matched {
			if log.V(2) {
				log.Infof("path not matched by patterns: %v", fullPath)
			}
			continue
		}

		interfaceName, queueIDStr, leafName, err := parsePath(fullPath)
		if err != nil {
			if log.V(2) {
				log.Infof("matched path but failed to parse: %v, err: %v", fullPath, err)
			}
			continue
		}

		// Update to a member's port-channel assignment.
		if leafName == "aggregate-id" {
			newPCName := update.GetVal().GetStringVal()
			i.handleAggregateIDUpdate(interfaceName, newPCName, impactedPortChannels)
			continue
		}

		// This ensures the singleton port QOS counters are preserved.
		passthroughUpdates = append(passthroughUpdates, &gnmipb.Update{
			Path: relativeQoSPath(prefix, fullPath),
			Val:  update.GetVal(),
		})
		val := update.GetVal().GetUintVal()
		i.handleQoSUpdate(interfaceName, queueIDStr, leafName, val, impactedPortChannels)
	}

	// Aggregate and generate the aggregate updates.
	var aggregateUpdates []*gnmipb.Update
	for pcName := range impactedPortChannels {
		if log.V(2) {
			log.Infof("recalculating aggregates for impacted Port-Channel: %s", pcName)
		}
		newUpdates := i.aggregateAndBuildUpdates(pcName)
		aggregateUpdates = append(aggregateUpdates, newUpdates...)
	}

	// Combine the original passthrough updates with the new aggregated updates.
	finalUpdates := append(passthroughUpdates, aggregateUpdates...)

	if len(finalUpdates) == 0 {
		return nil, nil
	}

	// Factor the longest common OpenConfig prefix (/qos/interfaces) shared across
	// all member passthrough and Port-Channel aggregate updates into Notification.Prefix.
	outgoingNotification := &gnmipb.Notification{
		Timestamp: timestamp,
		Prefix: &gnmipb.Path{
			Origin: "openconfig",
			Target: target,
			Elem:   qosInterfacesPrefix,
		},
		Update: finalUpdates,
	}

	return &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{Update: outgoingNotification},
	}, nil
}

// relativeQoSPath returns fullPath relative to the outgoing Notification
// Prefix (/openconfig/qos/interfaces), stripping Origin, Target, and the
// leading [qos, interfaces] elements.
//
// When prefix != nil, ftutilities.Join already allocated a fresh Path struct,
// so updating it in place adds zero heap allocations on the collector path.
// When prefix == nil, ftutilities.Join returns the incoming Path pointer
// directly; allocating a shallow Path header avoids mutating shared input
// messages fanned out across concurrent subscribers.
func relativeQoSPath(prefix, fullPath *gnmipb.Path) *gnmipb.Path {
	if prefix == nil {
		return &gnmipb.Path{Elem: fullPath.GetElem()[2:]}
	}
	fullPath.Origin = ""
	fullPath.Target = ""
	fullPath.Elem = fullPath.Elem[2:]
	return fullPath
}
