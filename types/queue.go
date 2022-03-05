package types

import "time"

type QueueTask struct {
	Height int64
	retry  int
}

func (task *QueueTask) DoRetry() {
	task.retry += 1
}

// GetRetryTimeout - returns duration before retry block parsing
func (task *QueueTask) GetRetryTimeout() time.Duration {
	secondsAwait := 1*2 ^ task.retry
	return time.Duration(secondsAwait) * time.Second
}

func NewQueueTask(height int64) QueueTask {
	return QueueTask{
		Height: height,
		retry:  0,
	}
}

// HeightQueue is a simple type alias for a (buffered) channel of block heights.
type HeightQueue chan QueueTask

func NewQueue(size int) HeightQueue {
	return make(chan QueueTask, size)
}
