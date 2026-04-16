package dto

import (
	cohere "github.com/cohere-ai/cohere-go/v2"
)

type UserQuestion struct {
	Question string `json:"question"`
}

type AskResponse struct {
	Content   []*cohere.AssistantMessageResponseContentItem `json:"content,omitempty" url:"content,omitempty"`
	Citations []*cohere.Citation                            `json:"citations,omitempty" url:"citations,omitempty"`
}
