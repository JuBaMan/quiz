package main

import "encoding/csv"
import "os"
import "fmt"
import "log"
import "bufio"
import "io"
import "strings"
import "flag"
import "time"

func main(){

	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for each question in seconds")
	flag.Parse()

	csvFile, err := os.Open("problems.csv")
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
	}

	readerCSV := csv.NewReader(bufio.NewReader(csvFile))
	readerConsole := bufio.NewReader(os.Stdin)
	
	var correct int
	var count int

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	for {
        line, error := readerCSV.Read()
        count++

        if error == io.EOF {
        	fmt.Println(count - 1, correct)
            break
        } else if error != nil {
            log.Fatal(error)
        }

        fmt.Println(line[0])

        answerCh := make(chan string)
        
        go func() {
        	text, _ := readerConsole.ReadString('\n')
        	answerCh <- text
        }()

        select {
        	case <-timer.C:
        		fmt.Println("You failed.")
        		fmt.Println(count, correct)
        		return
        	case text := <-answerCh:
        		if strings.TrimSpace(strings.TrimRight(text, "\n")) == strings.TrimSpace(line[1]) {
        			correct++
        		}
        }
	}
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}