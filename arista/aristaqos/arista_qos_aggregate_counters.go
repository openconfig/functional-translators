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
)

const (
	leafAggregateID    = "aggregate-id"
	leafTransmitOctets = "transmit-octets"
	leafTransmitPkts   = "transmit-pkts"
	leafDroppedOctets  = "dropped-octets"
	leafDroppedPkts    = "dropped-pkts"
)

// New creates a functional translator.
func New() *translator.FunctionalTranslator {
	ft, err := translator.NewFunctionalTranslator(
		translator.FunctionalTranslatorOptions{
			ID:               ftconsts.AristaQoSAggregateCountersTranslator,
			Translate:        translate,
			OutputToInputMap: ftutilities.MustStringMapPaths(translateMap),
			Metadata: []*translator.FTMetadata{
				{
					Vendor: ftconsts.VendorArista,
					SoftwareVersionRange: &translator.SWRange{
						InclusiveMin: "4.33.0F",
						ExclusiveMax: "4.35",
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

// returnPathForQOSAggregatedCounter creates a gNMI path for an aggregated QoS counter on a given port-channel.
func returnPathForQOSAggregatedCounter(interfaceName, queueID, leafName string) *gnmipb.Path {
	return &gnmipb.Path{
		Elem: []*gnmipb.PathElem{
			{Name: "qos"}, {Name: "interfaces"},
			{
				Name: "interface",
				Key: map[string]string{
					"interface-id": interfaceName,
				},
			},
			{Name: "output"}, {Name: "queues"},
			{
				Name: "queue",
				Key: map[string]string{
					"name": queueID,
				},
			},
			{Name: "state"},
			{Name: leafName},
		},
	}
}

// newCounterUpdate creates a gNMI update for a given port-channel, queueID, and counter leaf.
func newCounterUpdate(pcName, newCompositeQueueID, leafName string, value uint64) *gnmipb.Update {
	path := returnPathForQOSAggregatedCounter(pcName, newCompositeQueueID, leafName)
	return &gnmipb.Update{
		Path: path,
		Val:  &gnmipb.TypedValue{Value: &gnmipb.TypedValue_UintVal{UintVal: value}},
	}
}

// aggregateAndBuildUpdates calculates the sum of counters for a port-channel and creates gNMI updates.
func aggregateAndBuildUpdates(target, pcName string) []*gnmipb.Update {
	targetInfo, ok := ftutilities.QoSAggMap.RetrieveTargetQoSInfo(target)
	if !ok {
		return nil
	}
	pcInfo, ok := targetInfo.PortChannelInfo(pcName)
	if !ok {
		return nil
	}

	// aggregatedCounters maps a simple queue name to its summed counters.
	aggregatedCounters := make(map[string]*ftutilities.QueueCounters)

	// Sum counters from all members for each queue.
	for _, memberInfo := range pcInfo.CloneMembers() {
		// CloneQueues provides a thread-safe copy of the member's counters.
		for queueID, counters := range memberInfo.CloneQueues() {
			agg, ok := aggregatedCounters[queueID]
			if !ok {
				agg = new(ftutilities.QueueCounters)
				aggregatedCounters[queueID] = agg
			}

			agg.TxBytes += counters.TxBytes
			agg.TxPackets += counters.TxPackets
			agg.DroppedBytes += counters.DroppedBytes
			agg.DroppedPackets += counters.DroppedPackets
		}
	}
	var outgoingUpdates []*gnmipb.Update
	// Build gNMI Updates from the aggregated results.
	for simpleQueueName, counters := range aggregatedCounters {
		// Create the new composite queue ID for the aggregated path.
		// e.g., if pcName is "Port-Channel10" and simpleQueueName is "0", this becomes "Port-Channel10-0".
		newCompositeQueueID := fmt.Sprintf("%s-%s", pcName, simpleQueueName)

		outgoingUpdates = append(outgoingUpdates,
			newCounterUpdate(pcName, newCompositeQueueID, leafTransmitOctets, counters.TxBytes),
			newCounterUpdate(pcName, newCompositeQueueID, leafTransmitPkts, counters.TxPackets),
			newCounterUpdate(pcName, newCompositeQueueID, leafDroppedOctets, counters.DroppedBytes),
			newCounterUpdate(pcName, newCompositeQueueID, leafDroppedPkts, counters.DroppedPackets),
		)
	}
	return outgoingUpdates
}

// deleteHandler processes delete notifications to handle port-channel member removals.
// It updates the cache and returns a map of port-channel names that were affected.
func deleteHandler(n *gnmipb.Notification) (impactedPortChannels map[string]bool) {
	impactedPortChannels = make(map[string]bool)
	prefix := n.GetPrefix()
	target := prefix.GetTarget()

	for _, delPath := range n.GetDelete() {
		fullPath := ftutilities.Join(prefix, delPath)
		interfaceName, _, _, err := parsePath(fullPath)
		if err != nil {
			log.Warningf("failed to parse delete path %v: %v", fullPath, err)
			continue
		}

		if oldPCName, removed := processMemberRemoval(target, interfaceName); removed {
			impactedPortChannels[oldPCName] = true
		}
	}
	return impactedPortChannels
}

// processMemberRemoval updates the cache and returns the name of the impacted port-channel and whether a removal occurred.
func processMemberRemoval(target, interfaceName string) (string, bool) {
	targetInfo, ok := ftutilities.QoSAggMap.RetrieveTargetQoSInfo(target)
	if !ok {
		return "", false
	}
	// Also clean up the member if it's in the "waiting room"
	if _, found := targetInfo.UnassociatedMembers[interfaceName]; found {
		delete(targetInfo.UnassociatedMembers, interfaceName)
		log.V(2).Infof("removed unassociated member %s from the waiting room on target %s.", interfaceName, target)
		return "", true
	}

	// FindAndRemoveMember cleans up the cache and returns the old port-channel.
	if oldPCName, found := targetInfo.FindAndRemoveMember(interfaceName); found {
		log.V(2).Infof("processed removal of member %s from Port-Channel %s on target %s.", interfaceName, oldPCName, target)
		return oldPCName, true
	}
	return "", false
}

// handleAggregateIDUpdate processes membership changes based on an aggregate-id update.
// It updates the cache and the set of impacted port-channels.
func handleAggregateIDUpdate(targetInfo *ftutilities.TargetQoSInfo, interfaceName, newPCName string, impactedPortChannels map[string]bool) {
	// Process the implicit removal from any old port-channel.
	if oldPCName, removed := targetInfo.FindAndRemoveMember(interfaceName); removed {
		impactedPortChannels[oldPCName] = true
	}

	// If newPCName is empty, it signifies a removal.
	// We have no new port-channel to add to, so we can return early.
	if newPCName == "" {
		return
	}

	// Handle the addition to the new port-channel if one is specified.
	pcInfo := targetInfo.CreateOrRetrievePortChannel(newPCName)
	targetInfo.SetPortChannelForMember(interfaceName, newPCName)

	// Check the "waiting room" for pending counters for this interface.
	if unassociatedMember, found := targetInfo.UnassociatedMembers[interfaceName]; found {
		// Move the cached member info into the port-channel's member list.
		pcInfo.AddMemberInfo(unassociatedMember)
		// Clean up the waiting room.
		delete(targetInfo.UnassociatedMembers, interfaceName)
		log.V(2).Infof("moved pending counters for %s to Port-Channel %s", interfaceName, newPCName)
	} else {
		// If no pending data, just ensure the member struct exists.
		pcInfo.CreateOrRetrieveMember(interfaceName)
	}

	// Mark the port-channel as impacted to trigger an immediate aggregation.
	impactedPortChannels[newPCName] = true
	log.V(1).Infof("member %s assigned to new Port-Channel %s", interfaceName, newPCName)
}

// handleQoSUpdate processes a QoS counter update by updating the appropriate cache location.
func handleQoSUpdate(targetInfo *ftutilities.TargetQoSInfo, interfaceName, simpleQueueName, leafName string, value uint64, impactedPortChannels map[string]bool) {
	var memberInfo *ftutilities.MemberInterfaceInfo

	// The member is already part of a known port-channel. Update its state directly.
	pcName, ok := targetInfo.RetrievePortChannelForMember(interfaceName)
	if ok {
		impactedPortChannels[pcName] = true

		// Get the PortChannelInfo and then the MemberInterfaceInfo.
		pcInfo, pcOk := targetInfo.PortChannelInfo(pcName)
		if !pcOk {
			log.Errorf("cache inconsistency: port-channel %s not found for member %s", pcName, interfaceName)
			return
		}
		memberInfo = pcInfo.CreateOrRetrieveMember(interfaceName)
	} else {
		// We don't know its LAG yet. Put its data in the "waiting room".
		memberInfo, ok = targetInfo.UnassociatedMembers[interfaceName]
		if !ok {
			// Member is not in the waiting room yet, so create it.
			memberInfo = ftutilities.NewMemberInterfaceInfo(interfaceName)
			targetInfo.UnassociatedMembers[interfaceName] = memberInfo
		}
	}

	switch leafName {
	case leafTransmitOctets:
		memberInfo.SetTxBytes(simpleQueueName, value)
	case leafTransmitPkts:
		memberInfo.SetTxPackets(simpleQueueName, value)
	case leafDroppedOctets:
		memberInfo.SetDroppedBytes(simpleQueueName, value)
	case leafDroppedPkts:
		memberInfo.SetDroppedPackets(simpleQueueName, value)
	}
}

func translate(sr *gnmipb.SubscribeResponse) (*gnmipb.SubscribeResponse, error) {
	notification := sr.GetUpdate()
	if notification == nil {
		return nil, nil
	}

	prefix := notification.GetPrefix()
	target := prefix.GetTarget()
	timestamp := notification.GetTimestamp()

	// Handle deletes first and get the initial set of impacted port-channels.
	impactedPortChannels := deleteHandler(notification)

	// A slice to store all original QOS counter updates for singleton ports that should be passed through.
	var passthroughUpdates []*gnmipb.Update

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
			log.V(2).Infof("path not matched by patterns: %v", fullPath)
			continue
		}

		interfaceName, queueIDStr, leafName, err := parsePath(fullPath)
		if err != nil {
			log.V(2).Infof("matched path but failed to parse: %v, err: %v", fullPath, err)
			continue
		}

		targetInfo := ftutilities.QoSAggMap.CreateOrUpdateTargetQoSInfo(target)

		// Update to a member's port-channel assignment.
		if leafName == "aggregate-id" {
			newPCName := update.GetVal().GetStringVal()
			handleAggregateIDUpdate(targetInfo, interfaceName, newPCName, impactedPortChannels)
			continue
		}

		// This ensures the singleton port QOS counters are preserved.
		passthroughUpdates = append(passthroughUpdates, update)
		val := update.GetVal().GetUintVal()
		handleQoSUpdate(targetInfo, interfaceName, queueIDStr, leafName, val, impactedPortChannels)
	}

	// Aggregate and generate the aggregate updates.
	var aggregateUpdates []*gnmipb.Update
	for pcName := range impactedPortChannels {
		log.V(2).Infof("recalculating aggregates for impacted Port-Channel: %s", pcName)
		newUpdates := aggregateAndBuildUpdates(target, pcName)
		aggregateUpdates = append(aggregateUpdates, newUpdates...)
	}

	// Combine the original passthrough updates with the new aggregated updates.
	finalUpdates := append(passthroughUpdates, aggregateUpdates...)

	if len(finalUpdates) == 0 {
		return nil, nil
	}

	outgoingNotification := &gnmipb.Notification{
		Timestamp: timestamp,
		Prefix:    &gnmipb.Path{Origin: "openconfig", Target: target},
		Update:    finalUpdates,
	}

	return &gnmipb.SubscribeResponse{
		Response: &gnmipb.SubscribeResponse_Update{Update: outgoingNotification},
	}, nil
}
