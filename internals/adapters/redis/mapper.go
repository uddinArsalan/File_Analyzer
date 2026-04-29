package redis

import (
	"encoding/json"
	"file-analyzer/queue"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type XPending struct {
	ID         string
	Consumer   string
	Idle       time.Duration
	RetryCount int64
}

func ToXPendingList(res []redis.XPendingExt) []XPending {
	pendingEntryList := make([]XPending, len(res))
	for i, val := range res {
		pendingEntryList[i] = XPending{
			ID:         val.ID,
			Consumer:   val.Consumer,
			Idle:       val.Idle,
			RetryCount: val.RetryCount,
		}
	}
	return pendingEntryList
}

func ToJobsList(msgsList []redis.XMessage) []queue.Job {
	jobs := make([]queue.Job, len(msgsList))
	for i, msg := range msgsList {
		raw, ok := msg.Values["data"].(string)
		if !ok {
			continue
		}
		job, err := ParseJob(raw)
		if err != nil {
			fmt.Printf("Error Parsing Job %v", err)
			continue
		}
		jobs[i] = *job
	}
	return jobs
}

func ParseJob(item string) (*queue.Job, error) {
	var job queue.Job
	err := json.Unmarshal([]byte(item), &job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}
