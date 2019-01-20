package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkIfFileExist(filePath string) {
	_, err := os.Open(filePath)
	check(err)
}

func readFile(filePath string) *os.File {
	f, err := os.Open(filePath)
	check(err)
	return f
}

func main() {

	filePath := flag.String("file", "./result", "file path")
	flag.Parse() // pare variables from commnad line
	file := readFile(*filePath)
	defer file.Close()

	var fileStreamErr error

	// no matched songs record
	var noMatch *os.File
	noMatch, fileStreamErr = os.Create("./noMatched.log")
	check(fileStreamErr)
	defer noMatch.Close() // close file stream

	scanner := bufio.NewScanner(file)

	totalItems := 0
	matchItems := 0
	correctHits := 0
	for scanner.Scan() {
		totalItems++
		sText := scanner.Text()
		sText = strings.Replace(sText, "\n", "", -1)
		resultList := strings.Split(sText, "+")
		correctName := ""
		hitFlag := false

		// check if query answer is in resultList
		for i, v := range resultList {
			if i == 0 {
				// top: [queryName, result, scope]
				top := strings.Split(v, "\t")
				correctName = top[0]
				if strings.Contains(correctName, top[1]) {
					correctHits++
					hitFlag = true
				}
			} else {
				// otherwise: [result, scope]
				rankResult := strings.Split(v, "\t")
				if strings.Contains(correctName, rankResult[0]) {
					matchItems++
					hitFlag = true
				}
			} //fi
		} // for
		if !hitFlag {
			data := []byte(correctName)
			_, err := noMatch.Write(append(data, byte('\n')))
			check(err)
			noMatch.Sync()
		}
	} // for

	fmt.Println(
		totalItems, // song total number
		correctHits,
		float64(correctHits)/float64(totalItems), // correct hit rate
		matchItems, // number of the  correct answer of the query not in top1
		((correctHits + matchItems) == totalItems), // if true -> queries are correct
		(correctHits + matchItems), float64(correctHits+matchItems)/float64(totalItems),
	)
} // main
