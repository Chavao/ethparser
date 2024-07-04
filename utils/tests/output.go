package tests

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func CaptureOutput(t *testing.T) (r *os.File, w *os.File, s *os.File) {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	originalStdout := os.Stdout

	os.Stdout = writer

	return reader, writer, originalStdout
}

func ReadOutput(reader *os.File, writer *os.File) string {
	writer.Close()

	var buf bytes.Buffer
	io.Copy(&buf, reader)

	return buf.String()
}
