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

func getinput(input chan string) {
	var answer string
	fmt.Scanf("%s\n", &answer)
	input <- answer
}
func parseQuiz(records [][]string) []quiz {
	result := make([]quiz, len(records))
	for i, record := range records {
		result[i] = quiz{
			question: record[0],
			answer:   strings.TrimSpace(record[1]),
		}
	}
	return result
}

type quiz struct {
	question string
	answer   string
}

func main() {
	quizfile := flag.String("csv", "problems.csv", "csv file containing the quiz")
	limit := flag.Int("limit", 10, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*quizfile)
	if err != nil {
		log.Fatalln(err)
	}
	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}
	problems := parseQuiz(records)

	correct := 0
	var flag int
	for _, p := range problems {
		fmt.Printf(" %s = ", p.question)
		input := make(chan string)
		go getinput(input)
		select {
		case result := <-input:
			if p.answer == result {
				correct++
			}
		case <-time.After((time.Duration)(*limit) * time.Second):
			fmt.Println("\n timed out")
			flag = 1
		}

		if flag == 1 {
			break
		}

	}

	fmt.Printf("%d qustions out of which %d  answered right", len(problems), correct)
}
