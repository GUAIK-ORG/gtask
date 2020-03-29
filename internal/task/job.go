package task

import (
	"time"
)

type Job struct {
	Key            string
	task           *Task
	resetEvent     chan int64 // 检测存活
	doneEvent      chan int64 // 完成任务
	ProcessorList  []*Processor
	processorCount int64
	state          int64
	saveDelay      int64
	lastSaveTime   int64
}

func NewJob(task *Task, Key string) *Job {
	return &Job{
		Key:            Key,
		task:           task,
		resetEvent:     make(chan int64, 1),
		doneEvent:      make(chan int64, 1),
		ProcessorList:  make([]*Processor, 0),
		processorCount: 0,
		state:          0,
		saveDelay:      5, // (秒)
		lastSaveTime:   time.Now().Unix(),
	}
}

// 获取工作者数量
func (job *Job) GetProcessorCount() int64 {
	return job.processorCount
}

// 获取状态
func (job *Job) GetState() int64 {
	return job.state
}

// 添加超时处理者
func (job *Job) AddProcessor(code string, trigger int64, bReset, bLoop, bExit bool) {
	processor := NewProcessor(code, trigger, bReset, bLoop, bExit)
	if processor != nil {
		job.processorCount++
		job.ProcessorList = append(job.ProcessorList, processor)
	}

}

func (job *Job) Release() {
	for _, t := range job.ProcessorList {
		t.Release()
	}
}

func (job *Job) Run() {
	job.task.waitGroup.Add(1)
	go job.worker()
}

func (job *Job) Stop() {
	job.doneEvent <- 1
}

func (job *Job) Reset() {
	job.resetEvent <- 1
}

func (job *Job) worker() {
	defer job.task.waitGroup.Done()
	for {
		select {
		case <-job.resetEvent:
			for _, proc := range job.ProcessorList {
				if proc.bReset {
					proc.count = 0
				}
			}
		case <-job.doneEvent:
			job.Release()
			return
		case <-time.After(time.Second):
			for index, proc := range job.ProcessorList {
				proc.count++
				if proc.count >= proc.trigger {
					// 执行任务
					if !proc.Run(job.Key, proc.count) {
						job.Release()
						job.task.DeleteJob(job.Key)
						return
					}
					if proc.bExit {
						// 超时退出Page
						job.Release()
						job.task.DeleteJob(job.Key)
						return
					}
					if proc.bLoop {
						proc.count = 0
					} else {
						// 超时移除处理器
						proc.Release()
						job.ProcessorList = append(job.ProcessorList[:index], job.ProcessorList[index+1:]...)
					}
				}
			}
			// 更新数据
			if time.Now().Unix()-job.lastSaveTime > job.saveDelay {
				job.task.Save(job.Key, job.ProcessorList)
				job.lastSaveTime = time.Now().Unix()
			}
		}
	}
}
