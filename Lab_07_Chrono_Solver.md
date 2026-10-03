# Lab 7: Chrono-Solver
## CS 411: Computing Paradigms (Fall 2026)

**Starter repository (GitHub template):** [https://github.com/cwillpinto/cs411-lab07-chrono-solver](https://github.com/cwillpinto/cs411-lab07-chrono-solver)

**Schedule:** Tuesday, October 6 is Fall Break. Thursday, October 8 is the single 75-minute Week 8 session. This 100-point Weekly Assessment is due Monday, October 12, 2026 at 11:59 PM Central.

### ⚡ The Mission
Three servers produce logs with wall-clock timestamps that disagree with the observed message path. A receive appears to occur before its send. Clock drift and network delay make wall time unreliable evidence of causality. Your task is to reconstruct what each server could know from its local events and received messages.

Lamport clocks preserve the direction of a known happened-before relationship: if A caused B, A has a smaller timestamp. The reverse is not guaranteed. You will implement a Go queue that sorts buffered events by **Lamport logical timestamp**, then use vector clocks and a causal trace to distinguish a verified dependency from an arbitrary order between concurrent events.

---

### Part A: The Systems Analyst (50 Points)
**Instructions:** Answer the following questions in a separate document (PDF or Markdown).

1. Three nodes `P1`, `P2`, `P3` each start with a Lamport clock of `0`. The following events occur in this real-time order (real time is given only for reference — your job is to compute logical time, not report real time):
   - `e1`: `P1` performs local work.
   - `e2`: `P1` sends message `m1` to `P2`.
   - `e3`: `P3` performs local work.
   - `e4`: `P2` receives `m1`.
   - `e5`: `P2` sends message `m2` to `P3`.
   - `e6`: `P3` receives `m2`.
   - `e7`: `P3` sends message `m3` to `P1`.
   - `e8`: `P1` receives `m3`.

   Using the Lamport clock rules (increment on every local event; on send, attach current clock; on receive, set clock to `max(local, received) + 1`), compute the Lamport timestamp of every event `e1` through `e8`. Show the clock value at each node immediately after each event.
2. Using your timestamps from Problem 1, list all pairs of events that are **definitely causally ordered** (one happened-before the other via a chain of local-order or send/receive relationships) and identify the causal chain connecting them for at least two such pairs (e.g., `e1 → e2 → e4 → e5 → e6` is a valid causal chain — state which relationship, program order or message send/receive, justifies each link).
3. Identify one pair of events from Problem 1 that receive **different** Lamport timestamps but are **not** actually causally related (i.e., neither happened-before the other in any real causal sense — they're concurrent, just coincidentally ordered by the total order Lamport timestamps impose). Explain why Lamport timestamps alone cannot reveal that this pair is truly concurrent.
4. Now assign **vector clocks** (a 3-element vector `[c1, c2, c3]`, one counter per node, incremented in the owning node's slot on every local event and merged element-wise via `max` on receive, exactly like Lamport but keeping a full vector instead of one integer) to the same 8 events from Problem 1, showing the full vector at each node after each event. Then, using the standard comparison rule (`VC(A) < VC(B)` iff every component of `VC(A)` is `<=` the corresponding component of `VC(B)` and at least one is strictly less; otherwise, if neither `VC(A) <= VC(B)` nor `VC(B) <= VC(A)`, the events are **concurrent**), determine the relationship between `e3` and `e4`, and separately between `e3` and `e6`. State for each pair whether it is happened-before, happened-after, or concurrent, and confirm whether this matches (or corrects) the Lamport-only conclusion you could have drawn from Problem 1's timestamps for the same two pairs.

---

### Part B: The Implementation (50 Points)
**Instructions:** You are provided a starter Go module (`chrono_solver/`) with a `clocks` package containing a `LamportClock` struct with `Tick()` and `Observe(received int)` methods **stubbed but unimplemented**, and a `MessageQueue` type that accepts incoming `LogEvent{NodeID string, Timestamp int, Payload string}` structs from multiple simulated nodes arriving out of order. It must emit a deterministic timestamp order. A sorted output respects happened-before relationships represented by valid Lamport timestamps, but it cannot establish that any two differently stamped events are causally related.

**Requirements:**
1. Implement `LamportClock.Tick() int` in `clocks/lamport.go`: increments the clock and returns the new value. Must be safe for concurrent use (protect internal state with a `sync.Mutex`).
2. Implement `LamportClock.Observe(received int) int`: on receiving a message stamped with `received`, sets the internal clock to `max(current, received) + 1` and returns the new value. Also safe for concurrent use.
3. Review the supplied `MessageQueue.Enqueue(event LogEvent)`, which already adds events to a mutex-protected buffer. Implement `MessageQueue.Drain() []LogEvent` to return all buffered events sorted primarily by `Timestamp` ascending and, for events with **equal** Lamport timestamps (a legitimate possibility across different nodes), by `NodeID` ascending. Document that the NodeID tiebreaker gives a deterministic total order; it does *not* imply real causal precedence.
4. Wire up the simulated 3-node producer scenario in `clocks/simulation.go` (the node loop and message-send logic are provided; you must call your `Tick()`/`Observe()` methods correctly at each local event, send, and receive point so the emitted timestamps are correct).
5. **Verification:** Run `go test -v ./clocks` and `go test -race ./clocks`. The suite checks send and receive timestamp relationships across 100 deterministic trials, as well as queue order and concurrent safety. If every test passes, it prints `SUCCESS TOKEN: CAUSAL-ORDER-VERIFIED`.

---

### 🚀 Submission
1. `Lab7_Theory.pdf` (or `.md`)
2. `clocks/lamport.go`, `clocks/simulation.go`, and any other modified files in `clocks/`
3. A terminal transcript of `go test -v ./clocks` and `go test -race ./clocks`
4. The **Success Token** generated by the test suite

### Grading Rubric (100 Points)

| Criterion | Points | Evidence |
| :--- | ---: | :--- |
| Lamport event trace | 20 | Correct `e1`–`e8` timestamps and per-node state after each event. |
| Causal and concurrent reasoning | 15 | Justified happened-before chains and a distinct-timestamp concurrent pair. |
| Vector-clock trace | 15 | Correct 3-element vectors and componentwise comparisons for `e3/e4` and `e3/e6`. |
| Race-safe Lamport clock | 20 | Correct `Tick` and `Observe` updates under a mutex. |
| Queue ordering | 10 | Timestamp order, deterministic NodeID ties, and comment that ties do not prove causality. |
| Three-node simulation | 10 | Tick/send/receive calls preserve send-before-receive relationships. |
| Reproducible verification | 10 | Verbose and race-test transcripts, 100 deterministic trials, and Success Token. |
