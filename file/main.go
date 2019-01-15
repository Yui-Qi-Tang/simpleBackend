package main

import (
	"bufio"
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
	filePath := "PUT_YOUR_DATA_PATH"
	file := readFile(filePath)
	defer file.Close()
	scanner := bufio.NewScanner(file)

	totalItems := 0
	correctHits := 0
	for scanner.Scan() {
		totalItems++
		sText := scanner.Text()
		sText = strings.Replace(sText, "\n", "", -1)
		resultList := strings.Split(sText, "+")
		for i, v := range resultList {
			if i == 0 {
				// top: [queryName, result, scope]
				top := strings.Split(v, "\t")
				if strings.Contains(top[0], top[1]) {
					// fmt.Println("+1")
					correctHits++
				}
			} //fi
			// otherwise: [result, scope]
			// TO-DO check other result?
		} // for
	} // for

	fmt.Println(totalItems, correctHits, float64(correctHits)/float64(totalItems))
} // main
