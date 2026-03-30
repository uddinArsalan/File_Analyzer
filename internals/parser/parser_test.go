package parser

import (
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	var (
		fileType = "pdf"
		filePath = "file.pdf"
	)
	out, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Error opening file %v", err.Error())
	}
	fileInfo, err := out.Stat()
	if err != nil {
		t.Fatalf("File path incorrect %v", err.Error())
	}
	pm := NewParserManager(out, fileInfo.Size())
	output, err := pm.ParseFile(fileType)
	if err != nil {
		t.Fatalf("Error parsing file %v", err.Error())
	}
	t.Logf("---CONTENT--- \n%v", output.Content)
}
