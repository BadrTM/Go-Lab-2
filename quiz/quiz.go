package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

// question struct stores a single question and its corresponding answer.
type question struct {
	q, a string
}

type score int

// check handles a potential error.
// It stops execution of the program ("panics") if an error has happened.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// questions reads in questions and corresponding answers from a CSV file into a slice of question structs.
func questions() []question {
	f, err := os.Open("quiz-questions.csv")
	check(err)
	reader := csv.NewReader(f)
	table, err := reader.ReadAll()
	check(err)
	var questions []question
	for _, row := range table {
		questions = append(questions, question{q: row[0], a: row[1]})
	}
	return questions
}

// ask asks a question and returns an updated score depending on the answer.
func ask(s score, question question, askChannel chan score) {
	fmt.Println(question.q)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter answer: ")
	scanner.Scan()
	text := scanner.Text()
	if strings.Compare(text, question.a) == 0 {
		fmt.Println("Correct!")
		s++
	} else {
		fmt.Println("Incorrect :-(")
	}
	askChannel <- s
}

func main() {
	s := score(0)
	qs := questions()
	askChannel := make(chan score)
	timer := time.After(5 * time.Second)
	timerBool := false

	for _, q := range qs {
		qAnswered := false
		if timerBool == false {
			go ask(s, q, askChannel)
		} else {
			break
		}
		for {
			if timerBool || qAnswered == true {
				break
			}
			select {
			case score := <-askChannel:
				s = score
				qAnswered = true
				break
			case <-timer:
				timerBool = true
				break
			}
		}
	}

	fmt.Println("\n\nFinal score", s)
}
