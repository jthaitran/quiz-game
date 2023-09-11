package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type problem struct {
	q string
	a string
}

type quizGame struct {
}
type quiz interface {
	loadProblems(filename string) ([]problem, error)
	startQuiz([]problem)
}

func (qg quizGame) loadProblems(csvFileName *string) ([]problem, error) {
	flag.Parse()
	file, err := os.Open(*csvFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	return parseLines(records), nil
}

func main() {
	fmt.Println("################")
	fmt.Println("Start Quiz Game")
	fmt.Println("################")
	quiz := quizGame{}

	csvFileName := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 3, "the time limit for the quiz in seconds")

	problems, err := quiz.loadProblems(csvFileName)
	if err != nil {
		exit(err.Error())
	}

	correct := quiz.startQuiz(problems, timeLimit)

	fmt.Printf("Your score %d out of %d.\n", correct, len(problems))

}
func (qg quizGame) startQuiz(problems []problem, timeLimit *int) int {
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("time out!!!\n")
			return correct
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}
	return correct
}
func parseLines(records [][]string) []problem {
	var problems []problem

	for _, line := range records {
		problem := problem{q: line[0], a: line[1]}
		problems = append(problems, problem)
	}

	return problems
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
