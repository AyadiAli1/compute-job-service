package job

import "testing"

func TestStateCanTransitionTo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		from, to State
		want     bool
	}{
		{StatePending, StateRunning, true},
		{StatePending, StateCancelled, true},
		{StatePending, StateSucceeded, false},
		{StateRunning, StateSucceeded, true},
		{StateRunning, StateFailed, true},
		{StateRunning, StateCancelled, true},
		{StateRunning, StatePending, false},
		{StateSucceeded, StateRunning, false},
		{StateFailed, StateCancelled, false},
		{StateCancelled, StateRunning, false},
	}

	for _, tt := range tests {
		t.Run(tt.from.String()+"->"+tt.to.String(), func(t *testing.T) {
			t.Parallel()
			if got := tt.from.CanTransitionTo(tt.to); got != tt.want {
				t.Errorf("CanTransitionTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateIsTerminal(t *testing.T) {
	t.Parallel()

	for s, want := range map[State]bool{
		StatePending:   false,
		StateRunning:   false,
		StateSucceeded: true,
		StateFailed:    true,
		StateCancelled: true,
	} {
		if got := s.IsTerminal(); got != want {
			t.Errorf("%s.IsTerminal() = %v, want %v", s, got, want)
		}
	}
}
