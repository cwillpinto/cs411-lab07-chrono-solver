// Lab 7: Chrono-Solver -- 3-node simulation.
//
// The node loop and message-send plumbing are given; you must call
// Tick()/Observe() at the TODO-marked points so the emitted
// timestamps are actually correct.
package clocks

import "sync"

// wireMessage is what travels over a simulated channel between nodes.
type wireMessage struct {
	fromNodeID string
	timestamp  int
	payload    string
}

// RunSimulation starts simulated nodes with individual Lamport clocks.
// Buffered messages are received after all sends complete, avoiding sleeps
// and allowing the harness to pair each send with its receive by payload.
func RunSimulation(nodeIDs []string, eventsPerNode int) *MessageQueue {
	queue := &MessageQueue{}
	clocks := make(map[string]*LamportClock, len(nodeIDs))
	inboxes := make(map[string]chan wireMessage, len(nodeIDs))
	for _, id := range nodeIDs {
		clocks[id] = &LamportClock{}
		inboxes[id] = make(chan wireMessage, 16)
	}

	var sends sync.WaitGroup

	for index, id := range nodeIDs {
		sends.Add(1)
		go func(id string, index int) {
			defer sends.Done()
			runNode(id, index, nodeIDs, eventsPerNode, clocks[id], inboxes, queue)
		}(id, index)
	}
	sends.Wait()
	for _, id := range nodeIDs { close(inboxes[id]) }
	var receives sync.WaitGroup
	for _, id := range nodeIDs {
		receives.Add(1)
		go func(id string) {
			defer receives.Done()
			for msg := range inboxes[id] {
				// TODO: call clocks[id].Observe(msg.timestamp), then enqueue
				// LogEvent{NodeID:id, Timestamp:receiveTS,
				// Payload:"receive:"+msg.payload}.
				_ = msg
			}
		}(id)
	}
	receives.Wait()
	return queue
}

func runNode(selfID string, index int, allNodeIDs []string, eventsPerNode int, clock *LamportClock, inboxes map[string]chan wireMessage, queue *MessageQueue) {
	for i := 0; i < eventsPerNode; i++ {
		// TODO: this is a LOCAL event. Call clock.Tick() and Enqueue a
		// LogEvent with the resulting timestamp, NodeID=selfID, and a
		// payload like "local-work".

		if len(allNodeIDs) > 1 {
			target := allNodeIDs[(index+i+1)%len(allNodeIDs)]
			if target == selfID { target = allNodeIDs[(index+i+2)%len(allNodeIDs)] }
			messageID := selfID + "-" + target + "-" + string(rune('A'+i))
			// TODO: this is a SEND event. Call clock.Tick() to get the
			// timestamp to attach to the outgoing message, Enqueue a
			// LogEvent{NodeID:selfID, Timestamp:sendTS,
			// Payload:"send:"+messageID}, then deliver the same sendTS.
			sendTS := 0 // TODO: replace with clock.Tick()
			inboxes[target] <- wireMessage{fromNodeID: selfID, timestamp: sendTS, payload: messageID}
		}
	}
}
