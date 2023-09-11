package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

// Helper function to create a temporary CSV file for testing
func createTempCSV(csvFileName string, content string) error {
	file, err := os.Create(csvFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

func TestLoadProblems(t *testing.T) {
	// Create a temporary CSV file for testing
	csvFileName := "test_problems.csv"
	csvContent := "question1,answer1\nquestion2,answer2\n"

	err := createTempCSV(csvFileName, csvContent)
	if err != nil {
		t.Fatalf("Error creating temporary CSV file: %v", err)
	}
	defer os.Remove(csvFileName)

	quiz := quizGame{}
	problems, err := quiz.loadProblems(&csvFileName)
	if err != nil {
		t.Fatalf("Error loading problems from CSV: %v", err)
	}

	expectedProblems := []problem{
		{q: "question1", a: "answer1"},
		{q: "question2", a: "answer2"},
	}

	if len(problems) != len(expectedProblems) {
		t.Fatalf("Expected %d problems, but got %d", len(expectedProblems), len(problems))
	}

	for i, p := range problems {
		if p != expectedProblems[i] {
			t.Errorf("Problem #%d: Expected %v, but got %v", i+1, expectedProblems[i], p)
		}
	}
}

func TestStartQuiz(t *testing.T) {
	problems := []problem{
		{q: "1+1", a: "2"},
		{q: "2+2", a: "4"},
	}

	// Simulate user input by creating a string containing the user's answers separated by newline characters (\n).
	userInput := "2\n4\n"
	timeLimit := 5
	timeLimitP := &timeLimit
	reader := strings.NewReader(userInput)

	// Create a custom os.File using a pipe
	r, w, _ := os.Pipe()
	os.Stdin = r

	// Close the write end of the pipe to avoid blocking
	defer func() {
		os.Stdin = os.NewFile(uintptr(3), "/dev/stdin") // Restore os.Stdin
		w.Close()
	}()

	/**
		Copying data from reader to w is a blocking operation. By using a goroutine, we ensure that this operation
		doesn't block the main testing function. If the copy operation were done synchronously in the main function,
		it could potentially block indefinitely if the data isn't available yet.
	**/
	go func() {
		io.Copy(w, reader)
	}()

	quiz := quizGame{}
	correct := quiz.startQuiz(problems, timeLimitP)

	expectedCorrect := 2
	if correct != expectedCorrect {
		t.Errorf("Expected %d correct answers, but got %d", expectedCorrect, correct)
	}
}
