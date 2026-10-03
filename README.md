# Lab 7: Chrono-Solver

Full assignment: `Lab_07_Chrono_Solver.md`.

## GitHub workflow
Use **Use this template** to create your own repository, clone it, complete the
starter TODOs, and commit and push your Go work. Submit the files and evidence
listed in the assignment through Blackboard.

The test suite runs 100 deterministic causal-order trials and prints the token
only if the whole package passes.

## Run
```bash
go test -v ./clocks
go test -race ./clocks
```
`MessageQueue.Enqueue` is already implemented in the starter. Complete
`LamportClock.Tick`/`Observe` and `MessageQueue.Drain` in
`clocks/lamport.go`, and wire up the TODO-marked call sites in
`clocks/simulation.go`. Success Token: `CAUSAL-ORDER-VERIFIED`.

## Submit
1. `Lab7_Theory.pdf` (or `.md`)
2. `clocks/lamport.go`, `clocks/simulation.go`, and any other modified files in `clocks/`
3. A transcript of both `go test -v ./clocks` and `go test -race ./clocks`
4. The Success Token
