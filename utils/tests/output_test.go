package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaptureOutput(t *testing.T) {
	reader, writer, originalStdout := CaptureOutput(t)
	defer func() {
		os.Stdout = originalStdout
	}()

	assert.NotNil(t, reader)
	assert.NotNil(t, writer)

	writer.Close()
}

func mockFunction() {
	fmt.Printf("Test Message!")
}

func TestReadOutput(t *testing.T) {
	reader, writer, originalStdout := CaptureOutput(t)
	defer func() {
		os.Stdout = originalStdout
	}()

	fmt.Fprint(os.Stdout, "Hello World!")

	output := ReadOutput(reader, writer)

	assert.Contains(t, output, "Hello World")
}

func TestCaptureAndReadOutput(t *testing.T) {
	reader, writer, originalStdout := CaptureOutput(t)
	defer func() {
		os.Stdout = originalStdout
	}()

	mockFunction()

	output := ReadOutput(reader, writer)

	assert.Contains(t, output, "Test Message")
}
