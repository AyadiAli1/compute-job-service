// Package job models compute jobs, their lifecycle and an in-memory store.
package job

import "slices"

// State is the lifecycle state of a job.
type State int

// Job lifecycle states. The zero value is intentionally invalid.
const (
	StatePending State = iota + 1
	StateRunning
	StateSucceeded
	StateFailed
	StateCancelled
)

var stateNames = map[State]string{
	StatePending:   "PENDING",
	StateRunning:   "RUNNING",
	StateSucceeded: "SUCCEEDED",
	StateFailed:    "FAILED",
	StateCancelled: "CANCELLED",
}

// String implements fmt.Stringer.
func (s State) String() string {
	if name, ok := stateNames[s]; ok {
		return name
	}
	return "UNKNOWN"
}

// IsTerminal reports whether no further transitions are possible.
func (s State) IsTerminal() bool {
	return s == StateSucceeded || s == StateFailed || s == StateCancelled
}

// transitions lists the allowed successor states of each state.
// Terminal states have no entry and therefore no successors.
var transitions = map[State][]State{
	StatePending: {StateRunning, StateCancelled},
	StateRunning: {StateSucceeded, StateFailed, StateCancelled},
}

// CanTransitionTo reports whether a job may move from s to next.
func (s State) CanTransitionTo(next State) bool {
	return slices.Contains(transitions[s], next)
}
