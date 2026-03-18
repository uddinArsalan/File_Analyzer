package subscriber

import "file-analyzer/internals/domain"

type Subscriber interface {
	Notify(msg domain.DocEvent)
}
