package textparser

import (
	"file-analyzer/internals/domain"
	"io"
)

type TextParser struct {
}

func NewTextParser() *TextParser {
	return &TextParser{}
}

func (txt *TextParser) Parse(stream io.Reader, size int64) (domain.DocumentParseResult, error) {
	data, err := io.ReadAll(stream)
	if err != nil {
		return domain.DocumentParseResult{}, err
	}

	return domain.DocumentParseResult{
		Content: string(data),
	}, nil
}
