package quiz

import "sync"

type JobManager struct {
	jobs map[string]chan ProgressEvent
	mu   sync.RWMutex
}

func NewJobManager() *JobManager {
	return &JobManager{
		jobs: make(map[string]chan ProgressEvent),
	}
}

func (jm *JobManager) CreateJob(jobID string) chan ProgressEvent {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	ch := make(chan ProgressEvent, 50)
	jm.jobs[jobID] = ch
	return ch
}

func (jm *JobManager) GetJob(jobID string) (chan ProgressEvent, bool) {
	jm.mu.RLock()
	defer jm.mu.RUnlock()

	ch, ok := jm.jobs[jobID]
	return ch, ok
}

func (jm *JobManager) RemoveJob(jobID string) {
	jm.mu.Lock()
	defer jm.mu.Unlock()
	delete(jm.jobs, jobID)
}