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

package aristaqosaggregatecounters

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"github.com/openconfig/functional-translators/ftutilities"
)

// setupStateForTwoMembers pre-populates the cache with one member in a port-channel.
func setupStateForTwoMembers(state *qosState) {
	pcInfo := state.createOrRetrievePortChannel("Port-Channel10")
	state.setPortChannelForMember("Ethernet19/1", "Port-Channel10")
	member1 := pcInfo.createOrRetrieveMember("Ethernet19/1")
	q1 := member1.createOrGetQueue("0")
	q1.txBytes = 1000
	q1.txPackets = 100
	q1.droppedBytes = 20
	q1.droppedPackets = 2
}

// setupStateForMemberRemoval pre-populates the cache with two members.
// This is used to test that re-aggregation is correct when one member is removed.
func setupStateForMemberRemoval(state *qosState) {
	pcInfo := state.createOrRetrievePortChannel("Port-Channel10")

	// Member 1 (this one will be removed)
	state.setPortChannelForMember("Ethernet19/1", "Port-Channel10")
	member1 := pcInfo.createOrRetrieveMember("Ethernet19/1")
	q1 := member1.createOrGetQueue("0")
	q1.txBytes = 1000
	q1.txPackets = 100
	q1.droppedBytes = 20
	q1.droppedPackets = 2

	// Member 2 (this one will remain)
	state.setPortChannelForMember("Ethernet18/1", "Port-Channel10")
	member2 := pcInfo.createOrRetrieveMember("Ethernet18/1")
	q2 := member2.createOrGetQueue("0")
	q2.txBytes = 500
	q2.txPackets = 50
	q2.droppedBytes = 10
	q2.droppedPackets = 1
}

// setupStateForMultipleTargets pre-populates the cache for 'cx12.sql12'.
// This is used to verify that state is correctly isolated when an update for a new target arrives.
func setupStateForMultipleTargets(state *qosState) {
	pcInfo := state.createOrRetrievePortChannel("Port-Channel10")

	// Member on cx12.sql12
	state.setPortChannelForMember("Ethernet19/1", "Port-Channel10")
	member1 := pcInfo.createOrRetrieveMember("Ethernet19/1")
	q1 := member1.createOrGetQueue("0")
	q1.txBytes = 1000
	q1.txPackets = 100
	q1.droppedBytes = 20
	q1.droppedPackets = 2
}

// setupStateForCounterChange pre-populates the cache with two members.
func setupStateForCounterChange(state *qosState) {
	pcInfo := state.createOrRetrievePortChannel("Port-Channel10")

	state.setPortChannelForMember("Ethernet11/2", "Port-Channel10")
	member1 := pcInfo.createOrRetrieveMember("Ethernet11/2")
	q1 := member1.createOrGetQueue("0")
	q1.txBytes = 1000
	q1.txPackets = 100
	q1.droppedBytes = 10
	q1.droppedPackets = 1

	state.setPortChannelForMember("Ethernet22/3", "Port-Channel10")
	member2 := pcInfo.createOrRetrieveMember("Ethernet22/3")
	q2 := member2.createOrGetQueue("0")
	q2.txBytes = 500
	q2.txPackets = 50
	q2.droppedBytes = 5
	q2.droppedPackets = 0
}

func TestTranslate(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(state *qosState)
		inputPath      string
		wantOutputPath string
		wantNil        bool
		wantErr        bool
	}{
		{
			name:           "first_member_joins_and_updates_counter",
			inputPath:      "testdata/join_pc_and_counter_update_input.txt",
			wantOutputPath: "testdata/join_pc_and_counter_update_output.txt",
		},
		{
			name:           "multiple_updates_in_a_single_notification",
			inputPath:      "testdata/counter_change_input_single_notification.txt",
			wantOutputPath: "testdata/counter_change_output_single_notification.txt",
		},
		{
			name:           "second_member_joins_triggering_re_aggregation",
			setup:          setupStateForTwoMembers, // Pre-populates cache with the first member for cx12.sql12
			inputPath:      "testdata/two_members_aggregation_input.txt",
			wantOutputPath: "testdata/two_members_aggregation_output.txt",
		},
		{
			name:           "removing_a_member_from_a_bundle_triggers_re_aggregation",
			setup:          setupStateForMemberRemoval, // Pre-populates with both members for cx12.sql12
			inputPath:      "testdata/remove_member_input.txt",
			wantOutputPath: "testdata/remove_member_output.txt",
		},
		{
			name:           "multiple_targets_state_isolation",
			setup:          setupStateForMultipleTargets, // Pre-populates cache for cx12.sql12
			inputPath:      "testdata/multiple_targets_input.txt",
			wantOutputPath: "testdata/multiple_targets_output.txt",
		},
		{
			name:           "counter_change_triggers_re_aggregation",
			setup:          setupStateForCounterChange, // Pre-populates cache for cx12.sql12
			inputPath:      "testdata/counter_change_input.txt",
			wantOutputPath: "testdata/counter_change_output.txt",
		},
		{
			name:           "member_removed_while_still_in_the_waiting_room",
			inputPath:      "testdata/remove_from_waiting_room_input.txt",
			wantOutputPath: "testdata/remove_from_waiting_room_output.txt",
		},
		{
			name:      "malformed_aggregate_id_path_that_matches_pattern_but_fails_parsing",
			inputPath: "testdata/malformed_aggregate_id_input.txt",
			wantNil:   true,
		},
		{
			name:      "delete_handler_ignores_non_aggregate_id_paths",
			inputPath: "testdata/delete_non_aggregate_id_input.txt",
			wantNil:   true,
		},
		{
			name:      "qos_path_with_unexpected_queue_name_format",
			inputPath: "testdata/unexpected_queue_format_input.txt",
			wantNil:   true,
		},
		{
			name:      "invalid_update_path_that_does_not_match_patterns",
			inputPath: "testdata/invalid_update_path_input.txt",
			wantNil:   true, // The translator should ignore this update completely.
		},
		{
			name:      "invalid_delete_path_that_does_not_match_patterns",
			inputPath: "testdata/invalid_delete_path_input.txt",
			wantNil:   true, // The translator should ignore this delete.
		},
		{
			name:      "malformed_qos_path_that_matches_pattern_but_fails_parsing",
			inputPath: "testdata/malformed_qos_path_input.txt",
			wantNil:   true, // The update is matched but then rejected by the parser, resulting in no output.
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			state := newQoSState()
			testImpl := &impl{state: state}
			if tc.setup != nil {
				tc.setup(state)
			}

			inputSR, err := ftutilities.LoadSubscribeResponse(tc.inputPath)
			if err != nil {
				t.Fatalf("failed to load input message: %v", err)
			}

			gotSR, err := testImpl.translate(inputSR)

			if gotNil, gotErr := gotSR == nil, err != nil; gotNil || gotErr {
				switch {
				case gotErr != tc.wantErr:
					t.Fatalf("unexpected error result returned from translate() = %v, want error %t", err, tc.wantErr)
				case err != nil:
					return
				case gotNil != tc.wantNil:
					t.Fatalf("unexpected nil result returned from translate() = %t, want nil %t", gotNil, tc.wantNil)
				default:
					return
				}
			}

			if tc.wantNil || tc.wantErr {
				return
			}

			wantSR, err := ftutilities.LoadSubscribeResponse(tc.wantOutputPath)
			if err != nil {
				t.Fatalf("failed to load want message: %v", err)
			}
			if diff := cmp.Diff(wantSR, gotSR, protocmp.Transform()); diff != "" {
				t.Errorf("translate() returned unexpected diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestMemberInterfaceInfoCreateOrGetQueueNilMap(t *testing.T) {
	member := &memberInterfaceInfo{interfaceName: "Ethernet1"}
	q := member.createOrGetQueue("0")
	if q == nil {
		t.Fatal("createOrGetQueue() returned nil")
	}
	if got := member.createOrGetQueue("0"); got != q {
		t.Errorf("createOrGetQueue() returned %v, want %v", got, q)
	}
}

func TestPortChannelInfoNilMembers(t *testing.T) {
	pc := &portChannelInfo{portChannelName: "Port-Channel1"}
	m := pc.createOrRetrieveMember("Ethernet1")
	if m == nil {
		t.Fatal("createOrRetrieveMember() returned nil")
	}

	pc2 := &portChannelInfo{portChannelName: "Port-Channel2"}
	m2 := newMemberInterfaceInfo("Ethernet2")
	pc2.addMemberInfo(m2)
	if pc2.members["Ethernet2"] != m2 {
		t.Errorf("addMemberInfo() member not added to nil map")
	}
}
