package pdfparser

import (
	"bytes"
	"file-analyzer/internals/domain"
	"io"
	"strings"

	"github.com/ledongthuc/pdf"
)

type PdfParser struct {
}

func NewPdfParser() *PdfParser {
	return &PdfParser{}
}

func (pd *PdfParser) Parse(stream io.Reader, size int64) (domain.DocumentParseResult, error) {
	data, err := io.ReadAll(stream)
	if err != nil {
		return domain.DocumentParseResult{}, err
	}

	readerAt := bytes.NewReader(data)

	pdfReader, err := pdf.NewReader(readerAt, int64(len(data)))

	pdf.DebugOn = true
	if err != nil {
		return domain.DocumentParseResult{
			Content: "",
		}, err
	}

	var str strings.Builder

	totalPage := pdfReader.NumPage()
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := pdfReader.Page(pageIndex)
		if p.V.IsNull() || p.V.Key("Contents").Kind() == pdf.Null {
			continue
		}
		content, err := p.GetPlainText(nil)
		if err != nil {
			continue
		}

		str.WriteString(content)
	}

	return domain.DocumentParseResult{
		Content: str.String(),
	}, nil
}
