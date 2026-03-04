package pdfparser

import (
	"bytes"
	"file-analyzer/internals/domain"
	"io"
	"os"

	"github.com/ledongthuc/pdf"
)

type PdfParser struct {
}

func NewPdfParser() *PdfParser {
	return &PdfParser{}
}

func (pd *PdfParser) Parse(stream io.Reader) (domain.DocumentParseResult, error) {
	// creating temp file from stream
	out, err := os.Create("file.pdf")
	if err != nil {
		return domain.DocumentParseResult{
			Content: "",
		}, err
	}
	defer out.Close()
	_, err = io.Copy(out, stream)
	if err != nil {
		return domain.DocumentParseResult{
			Content: "",
		}, err
	}

	// Reading content from file
	pdf.DebugOn = true
	f, r, err := pdf.Open("./file.pdf")
	if err != nil {
		return domain.DocumentParseResult{
			Content: "",
		}, err
	}

	defer f.Close()
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		panic(err)
	}
	buf.ReadFrom(b)
	content := buf.String()

	return domain.DocumentParseResult{
		Content: content,
	}, nil
}
