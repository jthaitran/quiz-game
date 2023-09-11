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

func main() {
	fmt.Println("################")
	fmt.Println("Start Quiz Game")
	fmt.Println("################")

	file, err := os.Open("problems.csv")
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s.\n", "problems"))
	}

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		exit("Failed to read record the provided CSV file.")
	}

	problems := parseLines(records)

	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.q)

		var answer string
		fmt.Scanf("%s\n", &answer)
		if answer == p.a {
			correct++
		}
	}
	fmt.Printf("Your score %d out of %d.\n", correct, len(problems))
}

func parseLines(records [][]string) []problem {
	var problems []problem

	for _, line := range records {
		var problem problem
		problem.q = line[0]
		problem.a = line[1]
		problems = append(problems, problem)
	}

	return problems
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
