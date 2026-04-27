package redis

import (
	"file-analyzer/queue"
	"strconv"
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
		userIDStr, ok := msg.Values["user_id"].(string)
		if !ok {
			continue
		}
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			continue
		}
		newJob := queue.Job{
			ID:        msg.Values["id"].(string),
			UserID:    userID,
			ObjectKey: msg.Values["object_key"].(string),
			DocID:     msg.Values["doc_id"].(string),
			Mime_Type: msg.Values["mime_type"].(string),
		}
		jobs[i] = newJob
	}
	return jobs
}
