package test

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"posts/src/create_posts"
	"testing"
)

func TestInitializeLogger(t *testing.T) {
	originalStdout := os.Stdout
	var outputBuffer bytes.Buffer
	log.SetOutput(&outputBuffer)
	create_posts.InitializeLogger(&outputBuffer)
	defer func() { os.Stdout = originalStdout }()

	if _, err := os.Stat("server.log"); os.IsNotExist(err) {
		t.Error("The server.log file was not created")
	}
	log.Println("Test message")

	expectedOutput := "Test message\n"
	actualOutput := outputBuffer.String()
	if actualOutput != expectedOutput {
		t.Errorf("Invalid logger output..\nОжидалось: %s\nПолучено: %s", expectedOutput, actualOutput)
	}

	fileContent, err := ioutil.ReadFile("server.log")
	if err != nil {
		t.Errorf("Error reading file server.log: %v", err)
	}

	if string(fileContent) != expectedOutput {
		t.Errorf("Invalid logger output..\nОжидалось: %s\nПолучено: %s", expectedOutput, string(fileContent))
	}
	os.Remove("server.log")
}
