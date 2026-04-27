package domain

type TextContent struct {
	Text string `json:"text"`
}

type DocumentSource struct {
	ID       *string                `json:"id,omitempty"`
	Document map[string]interface{} `json:"document,omitempty"`
}

type Source struct {
	Type     string `json:"type,omitempty"`
	Document *DocumentSource
}

type Citation struct {
	Start *int `json:"start,omitempty"`
	// End index of the cited snippet in the original source text.
	End *int `json:"end,omitempty"`
	// Text snippet that is being cited.
	Text    *string   `json:"text,omitempty"`
	Sources []*Source `json:"sources,omitempty"`
	// Index of the content block in which this citation appears.
	ContentIndex *int `json:"content_index,omitempty"`
}

type AskResponse struct {
	Content   []*TextContent `json:"content,omitempty"`
	Citations []*Citation    `json:"citations,omitempty"`
}
