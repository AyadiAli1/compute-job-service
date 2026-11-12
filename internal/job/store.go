package job

import (
	"crypto/rand"
	"errors"
	"fmt"
	"sync"
	"time"

	computejobv1 "github.com/AyadiAli1/compute-job-service/gen/go/computejob/v1"
)

// Sentinel errors returned by Store; check them with errors.Is.
var (
	ErrNotFound          = errors.New("job not found")
	ErrInvalidTransition = errors.New("invalid state transition")
)

// Job is a snapshot of a compute job. Spec and Result are treated as
// immutable once set and may be shared between snapshots.
type Job struct {
	ID           string
	State        State
	Spec         *computejobv1.JobSpec
	Result       *computejobv1.JobResult
	ErrorMessage string
	CreatedAt    time.Time
	StartedAt    time.Time
	FinishedAt   time.Time
}

// maxUpdates is the maximum number of state changes a job goes through
// (pending -> running -> terminal). Watcher channels are buffered with
// this size, so notifying a watcher never blocks while the lock is held.
const maxUpdates = 2

type entry struct {
	job      Job
	watchers map[chan Job]struct{}
}

// Store is a concurrency-safe in-memory job store.
type Store struct {
	mu   sync.Mutex
	jobs map[string]*entry
	now  func() time.Time // replaceable in tests
}

// NewStore returns an empty Store.
func NewStore() *Store {
	return &Store{
		jobs: make(map[string]*entry),
		now:  time.Now,
	}
}

// Create stores a new pending job and returns its snapshot.
func (s *Store) Create(spec *computejobv1.JobSpec) Job {
	s.mu.Lock()
	defer s.mu.Unlock()

	j := Job{
		ID:        rand.Text(),
		State:     StatePending,
		Spec:      spec,
		CreatedAt: s.now(),
	}
	s.jobs[j.ID] = &entry{job: j, watchers: make(map[chan Job]struct{})}
	return j
}

// Get returns a snapshot of the job with the given ID.
func (s *Store) Get(id string) (Job, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.jobs[id]
	if !ok {
		return Job{}, fmt.Errorf("get %q: %w", id, ErrNotFound)
	}
	return e.job, nil
}

// Transition moves a job to the next state and records timestamps.
// The optional update function may set fields such as Result or
// ErrorMessage; it must not change State. All watchers receive the new
// snapshot; once the job is terminal, their channels are closed.
func (s *Store) Transition(id string, next State, update func(*Job)) (Job, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.jobs[id]
	if !ok {
		return Job{}, fmt.Errorf("transition %q: %w", id, ErrNotFound)
	}
	if !e.job.State.CanTransitionTo(next) {
		return e.job, fmt.Errorf("transition %q from %s to %s: %w",
			id, e.job.State, next, ErrInvalidTransition)
	}

	e.job.State = next
	if next == StateRunning {
		e.job.StartedAt = s.now()
	}
	if next.IsTerminal() {
		e.job.FinishedAt = s.now()
	}
	if update != nil {
		update(&e.job)
	}

	for ch := range e.watchers {
		ch <- e.job // never blocks: buffered with maxUpdates
		if next.IsTerminal() {
			close(ch)
			delete(e.watchers, ch)
		}
	}
	return e.job, nil
}

// Watch returns the current snapshot of a job and a channel that receives
// every subsequent state change. The channel is closed once the job is
// terminal. Call cancel when no longer interested; calling it multiple
// times or after the channel was closed is safe.
func (s *Store) Watch(id string) (snapshot Job, updates <-chan Job, cancel func(), err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.jobs[id]
	if !ok {
		return Job{}, nil, func() {}, fmt.Errorf("watch %q: %w", id, ErrNotFound)
	}

	ch := make(chan Job, maxUpdates)
	if e.job.State.IsTerminal() {
		close(ch)
		return e.job, ch, func() {}, nil
	}
	e.watchers[ch] = struct{}{}

	cancel = func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		if _, ok := e.watchers[ch]; ok {
			delete(e.watchers, ch)
			close(ch)
		}
	}
	return e.job, ch, cancel, nil
}
