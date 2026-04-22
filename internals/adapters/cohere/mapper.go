package cohere

import (
	"file-analyzer/internals/domain"
	cohere "github.com/cohere-ai/cohere-go/v2"
)

func AskResultToResponse(r *cohere.AssistantMessageResponse) *domain.AskResponse {
	return &domain.AskResponse{
		Content:   mapContent(r.Content),
		Citations: mapCitations(r.Citations),
	}
}

func mapContent(cohereContent []*cohere.AssistantMessageResponseContentItem) []*domain.TextContent {
	content := make([]*domain.TextContent, len(cohereContent))
	for i, c := range cohereContent {
		content[i] = &domain.TextContent{
			Text: c.Text.Text,
		}
	}
	return content
}

func mapCitations(cohereCitations []*cohere.Citation) []*domain.Citation {
	citations := make([]*domain.Citation, len(cohereCitations))
	for i, c := range cohereCitations {
		citations[i] = &domain.Citation{
			Start:        c.Start,
			End:          c.End,
			Text:         c.Text,
			ContentIndex: c.ContentIndex,
			Sources:      mapSources(c.Sources),
		}
	}
	return citations
}

func mapSources(cohereSources []*cohere.Source) []*domain.Source {
	sources := make([]*domain.Source, len(cohereSources))
	for i, src := range cohereSources {
		sources[i] = &domain.Source{
			Type: src.Type,
			Document: &domain.DocumentSource{
				ID:       src.Document.Id,
				Document: src.Document.Document,
			},
		}
	}
	return sources
}
