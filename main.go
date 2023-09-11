package main

import (
	"encoding/csv"
	"fmt"
	"os"
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

func (qg quizGame) loadProblems(filename string) ([]problem, error) {
	file, err := os.Open(filename)
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

	problems, err := quiz.loadProblems("problems.csv")
	if err != nil {
		exit(err.Error())
	}

	correct := quiz.startQuiz(problems)

	fmt.Printf("Your score %d out of %d.\n", correct, len(problems))

}
func (qg quizGame) startQuiz(problems []problem) int {
	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.q)

		var answer string
		fmt.Scanf("%s\n", &answer)
		if answer == p.a {
			correct++
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
