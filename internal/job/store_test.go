package job

import (
	"errors"
	"sync"
	"testing"
	"time"
)

// newTestStore returns a store whose clock advances one second per call.
// s.now is only called while s.mu is held, so the counter needs no lock.
func newTestStore() *Store {
	s := NewStore()
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	var ticks int
	s.now = func() time.Time {
		ticks++
		return base.Add(time.Duration(ticks) * time.Second)
	}
	return s
}

// receive reads one value from ch or fails the test after a timeout.
func receive(t *testing.T, ch <-chan Job) (Job, bool) {
	t.Helper()
	select {
	case j, ok := <-ch:
		return j, ok
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for job update")
		return Job{}, false
	}
}

func TestCreateAndGet(t *testing.T) {
	t.Parallel()
	s := newTestStore()

	created := s.Create(nil)
	if created.ID == "" || created.State != StatePending || created.CreatedAt.IsZero() {
		t.Fatalf("unexpected created job: %+v", created)
	}

	got, err := s.Get(created.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.ID != created.ID || got.State != StatePending {
		t.Errorf("Get() = %+v, want %+v", got, created)
	}
}

func TestGetUnknownJob(t *testing.T) {
	t.Parallel()
	s := newTestStore()

	if _, err := s.Get("does-not-exist"); !errors.Is(err, ErrNotFound) {
		t.Errorf("Get() error = %v, want ErrNotFound", err)
	}
}

func TestTransitionLifecycle(t *testing.T) {
	t.Parallel()
	s := newTestStore()
	j := s.Create(nil)

	running, err := s.Transition(j.ID, StateRunning, nil)
	if err != nil {
		t.Fatalf("Transition(Running) error = %v", err)
	}
	if running.StartedAt.IsZero() {
		t.Error("StartedAt not set")
	}

	failed, err := s.Transition(j.ID, StateFailed, func(j *Job) {
		j.ErrorMessage = "singular matrix"
	})
	if err != nil {
		t.Fatalf("Transition(Failed) error = %v", err)
	}
	if failed.ErrorMessage != "singular matrix" {
		t.Errorf("ErrorMessage = %q", failed.ErrorMessage)
	}
	if !failed.FinishedAt.After(failed.StartedAt) {
		t.Errorf("FinishedAt %v not after StartedAt %v", failed.FinishedAt, failed.StartedAt)
	}
}

func TestTransitionInvalid(t *testing.T) {
	t.Parallel()
	s := newTestStore()
	j := s.Create(nil)

	if _, err := s.Transition(j.ID, StateSucceeded, nil); !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("error = %v, want ErrInvalidTransition", err)
	}
	if got, _ := s.Get(j.ID); got.State != StatePending {
		t.Errorf("state changed to %s after invalid transition", got.State)
	}
	if _, err := s.Transition("does-not-exist", StateRunning, nil); !errors.Is(err, ErrNotFound) {
		t.Errorf("error = %v, want ErrNotFound", err)
	}
}

func TestWatchReceivesUpdatesUntilTerminal(t *testing.T) {
	t.Parallel()
	s := newTestStore()
	j := s.Create(nil)

	snapshot, updates, cancel, err := s.Watch(j.ID)
	if err != nil {
		t.Fatalf("Watch() error = %v", err)
	}
	defer cancel()
	if snapshot.State != StatePending {
		t.Fatalf("snapshot state = %s, want PENDING", snapshot.State)
	}

	mustTransition(t, s, j.ID, StateRunning)
	mustTransition(t, s, j.ID, StateSucceeded)

	for _, want := range []State{StateRunning, StateSucceeded} {
		got, ok := receive(t, updates)
		if !ok || got.State != want {
			t.Fatalf("update = (%s, open=%v), want %s", got.State, ok, want)
		}
	}
	if _, ok := receive(t, updates); ok {
		t.Error("channel still open after terminal state")
	}
}

func TestWatchTerminalJob(t *testing.T) {
	t.Parallel()
	s := newTestStore()
	j := s.Create(nil)
	mustTransition(t, s, j.ID, StateCancelled)

	snapshot, updates, cancel, err := s.Watch(j.ID)
	if err != nil {
		t.Fatalf("Watch() error = %v", err)
	}
	defer cancel()
	if snapshot.State != StateCancelled {
		t.Errorf("snapshot state = %s, want CANCELLED", snapshot.State)
	}
	if _, ok := receive(t, updates); ok {
		t.Error("channel of terminal job should be closed")
	}
}

func TestWatchCancel(t *testing.T) {
	t.Parallel()
	s := newTestStore()
	j := s.Create(nil)

	_, updates, cancel, err := s.Watch(j.ID)
	if err != nil {
		t.Fatalf("Watch() error = %v", err)
	}
	cancel()
	cancel() // second call must be a no-op

	mustTransition(t, s, j.ID, StateRunning) // must not panic (send on closed channel)
	if _, ok := receive(t, updates); ok {
		t.Error("channel should be closed after cancel")
	}
}

func TestWatchUnknownJob(t *testing.T) {
	t.Parallel()
	s := newTestStore()

	_, _, cancel, err := s.Watch("does-not-exist")
	defer cancel()
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("error = %v, want ErrNotFound", err)
	}
}

// TestStoreConcurrentAccess is mainly useful together with -race.
func TestStoreConcurrentAccess(t *testing.T) {
	t.Parallel()
	s := NewStore()

	var wg sync.WaitGroup
	for range 50 {
		wg.Go(func() {
			j := s.Create(nil)
			_, updates, cancel, err := s.Watch(j.ID)
			if err != nil {
				t.Error(err)
				return
			}
			defer cancel()

			if _, err := s.Transition(j.ID, StateRunning, nil); err != nil {
				t.Error(err)
			}
			if _, err := s.Transition(j.ID, StateSucceeded, nil); err != nil {
				t.Error(err)
			}

			n := 0
			for range updates {
				n++
			}
			if n != maxUpdates {
				t.Errorf("received %d updates, want %d", n, maxUpdates)
			}
		})
	}
	wg.Wait()
}

func mustTransition(t *testing.T, s *Store, id string, next State) {
	t.Helper()
	if _, err := s.Transition(id, next, nil); err != nil {
		t.Fatalf("Transition(%s) error = %v", next, err)
	}
}
