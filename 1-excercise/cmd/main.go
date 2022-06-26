package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	Wrong = iota
	TimeIsUp
	Correct
)

func Questionare(question, answer string, timeLimit int) int {
	var (
		input   string
		readVal = make(chan bool)
	)
	log.Print(question)
	go func() {
		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Panic(err)
		}
		readVal <- true
	}()
	select {
	case <-time.After(time.Duration(timeLimit) * time.Second):
		return TimeIsUp
	case <-readVal:
		if strings.Trim(input, " ") == strings.Trim(answer, " ") {
			return Correct
		} else {
			return Wrong
		}
	}
}

func StartQuiz(records [][]string, timeLimit *int) int {
	var correctAnswerCount = 0
	for _, col := range records {
		if len(col) < 2 {
			fmt.Printf("invalid input\n")
			return correctAnswerCount
		} else {
			if found := Questionare(col[0], col[1], *timeLimit); found == Correct {
				fmt.Printf("Correct\n")
				correctAnswerCount++
			} else if found == Wrong {
				fmt.Printf("Wrong answer, correct answer is %s\n", col[1])
				return correctAnswerCount
			} else {
				fmt.Printf("Time is up\n")
				return correctAnswerCount
			}
		}
	}
	return correctAnswerCount
}

func ParseProblems(path *string) [][]string {
	if path == nil {
		return [][]string{}
	}
	file, err := os.OpenFile(*path, os.O_RDWR, 0666)
	if err != nil {
		log.Panic(err)
	}
	scvReader := csv.NewReader(file)
	records, err := scvReader.ReadAll()
	if err != nil {
		log.Panic(err)
	}
	return records
}

func main() {
	var (
		pathToCSV = flag.String("csv", "problems.csv", "path to csv file")
		timeLimit = flag.Int("limit", 30, "time limit in seconds")
	)
	flag.Parse()
	fmt.Printf("Total correct answers count: %d", StartQuiz(ParseProblems(pathToCSV), timeLimit))
}
