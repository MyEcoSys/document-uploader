package queue

import (
	"github.com/MyEcoSys/document-uploader/jobs"
	"github.com/MyEcoSys/document-uploader/models"
)

// taskQueue is the channel-based queue
var taskQueue = make(chan models.Task, 100) // buffered channel

func Enqueue(task models.Task) {
    taskQueue <- task
}

func StartWorkerPool(workerCount int) {
    for i := range workerCount {
        go func(id int) {
            for task := range taskQueue {
                jobs.ProcessUploadTask(task)
            }
        }(i)
    }
}
