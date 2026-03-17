package parser

import (
	"errors"
	"file-analyzer/internals/adapters/pdfparser"
	"file-analyzer/internals/adapters/textparser"
	"file-analyzer/internals/domain"
	"io"
)

type Parser interface {
	// each parser has parser method
	Parse(stream io.Reader,size int64) (domain.DocumentParseResult, error)
}

type ParserManager struct {
	stream io.Reader
	size   int64
	// map of file type -> Parser for that file Type
	parsers map[string]Parser
}

func NewParserManager(stream io.Reader,size   int64) *ParserManager {
	pm := &ParserManager{
		stream: stream,
		size:   size,
		parsers: make(map[string]Parser),
	}
	pm.parsers["pdf"] = pdfparser.NewPdfParser()
	pm.parsers["text"] = textparser.NewTextParser()
	return pm
}

func (pm *ParserManager) AddNewParser(fileType string, parser Parser) {
	pm.parsers[fileType] = parser
}

func (pm *ParserManager) ParseFile(fileType string) (domain.DocumentParseResult, error) {
	parser, ok := pm.parsers[fileType]
	if !ok {
		return domain.DocumentParseResult{}, errors.New("unsupported file type")
	}
	return parser.Parse(pm.stream,pm.size)
}
