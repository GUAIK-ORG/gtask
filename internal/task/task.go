package task

import (
	"github.com/golang/glog"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/vmihailenco/msgpack"

	//"github.com/vmihailenco/msgpack"
	"sync"
)

type Task struct {
	mutex     sync.Mutex
	waitGroup *sync.WaitGroup
	jobs      map[string]*Job
	jobsCount int64
	state     int64
	db        *leveldb.DB
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

func (t *Task) Init(dbPath string) (err error) {
	t.state = 1
	t.db, err = leveldb.OpenFile(dbPath, nil)
	if err != nil {
		glog.Error(err)
	}
	t.load()
	return
}

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
		err := t.db.Delete([]byte(key), nil)
		if err != nil {
			glog.Error(err)
		}
		t.jobsCount--
	}
}

// 给Jos调用用来删除任务管理器中的key
func (t *Task) DeleteJob(key string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.state == 0 {
		return
	}
	if _, ok := t.jobs[key]; ok {
		delete(t.jobs, key)
		err := t.db.Delete([]byte(key), nil)
		if err != nil {
			glog.Error(err)
		}
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

func (t *Task) Release() {
	t.db.Close()
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

// 将任务存储到文件系统中
func (t *Task) Save(key string, processors []*Processor) {
	p := make([]map[string]interface{}, 0)
	for _, v := range processors {
		p = append(p, map[string]interface{}{
			"trigger": v.trigger,
			"bReset":  v.bReset,
			"bLoop":   v.bLoop,
			"bExit":   v.bExit,
			"code":    v.code,
		})
	}
	go func() {
		b, err := msgpack.Marshal(p)
		if err != nil {
			glog.Error(err)
		} else {
			err = t.db.Put([]byte(key), b, nil)
			if err != nil {
				glog.Error(err)
			}
		}
	}()
}

// 从文件系统中加载任务
func (t *Task) load() {
	iter := t.db.NewIterator(nil, nil)
	for iter.Next() {
		processors := make([]map[string]interface{}, 0)
		job := t.AddJob(string(iter.Key()))
		if job == nil {
			glog.Errorf("add job failed")
			continue
		}
		err := msgpack.Unmarshal(iter.Value(), &processors)
		if err != nil {
			glog.Error(err)
			continue
		}
		for _, p := range processors {
			job.AddProcessor(
				p["code"].(string),
				p["trigger"].(int64),
				p["bReset"].(bool),
				p["bLoop"].(bool),
				p["bExit"].(bool))
		}
		job.Run()
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		glog.Error(err)
	}
}
