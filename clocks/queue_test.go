// Lab 7: Chrono-Solver -- verification suite (given, do not edit).
//
// Run: go test -v ./clocks
// Prints "SUCCESS TOKEN: CAUSAL-ORDER-VERIFIED" only if every check
// below passes across 100 deterministic trial sizes.
package clocks

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// TestMain gates the Success Token on the WHOLE package's result, not on
// one test's own t.Failed(). Without this, a broken Tick()/Observe() can
// still make TestCausalOrderRandomized pass on its own (Drain's sort makes
// its check tautological) and print the token even while
// TestLamportClockBasics fails in the same run -- confirmed by temporarily
// stubbing Tick() to `return 1`, which failed TestLamportClockBasics but
// still printed the token before this fix.
func TestMain(m *testing.M) {
	code := m.Run()
	if code == 0 {
		fmt.Println("SUCCESS TOKEN: CAUSAL-ORDER-VERIFIED")
	}
	os.Exit(code)
}

func TestLamportClockBasics(t *testing.T) {
	c := &LamportClock{}
	if got := c.Tick(); got != 1 {
		t.Errorf("first Tick() = %d, want 1", got)
	}
	if got := c.Tick(); got != 2 {
		t.Errorf("second Tick() = %d, want 2", got)
	}
	// Observe with a LOWER received value: max(2, 1) + 1 = 3
	if got := c.Observe(1); got != 3 {
		t.Errorf("Observe(1) after clock=2 = %d, want 3", got)
	}
	// Observe with a HIGHER received value: max(3, 10) + 1 = 11
	if got := c.Observe(10); got != 11 {
		t.Errorf("Observe(10) after clock=3 = %d, want 11", got)
	}
}

func TestDrainOrdersByTimestampThenNodeID(t *testing.T) {
	q := &MessageQueue{}
	// Enqueue deliberately out of order.
	q.Enqueue(LogEvent{NodeID: "P3", Timestamp: 5, Payload: "e"})
	q.Enqueue(LogEvent{NodeID: "P1", Timestamp: 2, Payload: "b"})
	q.Enqueue(LogEvent{NodeID: "P2", Timestamp: 2, Payload: "a"}) // ties with the P1 event at Timestamp 2
	q.Enqueue(LogEvent{NodeID: "P1", Timestamp: 1, Payload: "d"})

	drained := q.Drain()
	expected := []string{"d", "b", "a", "e"} // ts1(P1) -> ts2(P1 before P2, NodeID tiebreak) -> ts2(P2) -> ts5(P3)
	if len(drained) != len(expected) {
		t.Fatalf("Drain() returned %d events, want %d", len(drained), len(expected))
	}
	for i, ev := range drained {
		if ev.Payload != expected[i] {
			t.Errorf("drained[%d].Payload = %q, want %q (full result: %+v)", i, ev.Payload, expected[i], drained)
			break
		}
	}
}

// TestCausalOrderRandomized runs 100 deterministic trial sizes and verifies
// each paired send/receive timestamp and drained position. Sorting alone
// cannot make a missing Observe pass this test.
func TestCausalOrderRandomized(t *testing.T) {
	const trials = 100
	nodeIDs := []string{"P1", "P2", "P3"}

	for trial := 0; trial < trials; trial++ {
		queue := RunSimulation(nodeIDs, 1+trial%5)
		drained := queue.Drain()
		sends := map[string]int{}
		receives := map[string]int{}
		positions := map[string]int{}
		for i, ev := range drained {
			positions[ev.Payload] = i
			if strings.HasPrefix(ev.Payload, "send:") { sends[strings.TrimPrefix(ev.Payload, "send:")] = ev.Timestamp }
			if strings.HasPrefix(ev.Payload, "receive:") { receives[strings.TrimPrefix(ev.Payload, "receive:")] = ev.Timestamp }
		}
		want := len(nodeIDs)*(1+trial%5)
		if len(sends) != want || len(receives) != want {
			t.Fatalf("trial %d: want %d paired sends and receives, got %d and %d", trial, want, len(sends), len(receives))
		}
		for id, sent := range sends {
			received, ok := receives[id]
			if !ok || received <= sent || positions["send:"+id] >= positions["receive:"+id] {
				t.Fatalf("trial %d message %s: send=%d receive=%d present=%v", trial, id, sent, received, ok)
			}
		}

		// Sanity: timestamps must be non-decreasing after Drain (that's
		// the whole point of Drain's sort).
		for i := 1; i < len(drained); i++ {
			if drained[i].Timestamp < drained[i-1].Timestamp {
				t.Fatalf("trial %d: Drain() result not sorted by timestamp at index %d: %+v", trial, i, drained)
			}
		}
	}
}
