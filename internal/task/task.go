package task

import "sync"

type Task struct {
	mutex     sync.Mutex
	waitGroup *sync.WaitGroup
	jobs      map[string]*Job
	jobsCount int64
	state     int64
}

var ins *Task
var once sync.Once

func Ins() *Task {
	once.Do(func() {
		ins = NewTask()
	})
	return ins
}

func NewTask() *Task {
	return &Task{
		waitGroup: &sync.WaitGroup{},
		jobs:      make(map[string]*Job),
		jobsCount: 0,
		state:     0,
	}
}

func (t *Task) Init() {}

func (t *Task) AddJob(key string) *Job {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.state == 0 {
		return nil
	}
	if _, ok := t.jobs[key]; ok {
		return nil
	}
	t.jobs[key] = NewJob(t, key)
	t.jobsCount++
	return t.jobs[key]
}

func (t *Task) GetJob(key string) *Job {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.state == 0 {
		return nil
	}
	if job, ok := t.jobs[key]; ok {
		return job
	}
	return nil
}

func (t *Task) RemoveJob(key string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.state == 0 {
		return
	}
	if job, ok := t.jobs[key]; ok {
		job.Stop()
		delete(t.jobs, key)
		t.jobsCount--
	}
}

func (t *Task) Loop() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.state == 1 {
		return
	}
	t.state = 1
	for _, job := range t.jobs {
		job.Run()
	}
	t.waitGroup.Wait()
}

func (t *Task) Start() {
	t.state = 1
}

func (t *Task) Wait() {
	t.waitGroup.Wait()
}

func (t *Task) Stop() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.state == 0 {
		return
	}
	t.state = 0
	for _, job := range t.jobs {
		job.Stop()
	}
}

func (t *Task) Reset() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.state == 0 {
		return
	}
	for _, job := range t.jobs {
		job.Reset()
	}
}
