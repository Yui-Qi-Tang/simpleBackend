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
	filePath := "PUT_YOUR_FILE_HERE"
	file := readFile(filePath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	sText := scanner.Text()
	sText = strings.Replace(sText, "\n", "", -1)

	resultList := strings.Split(sText, "+")

	for i, v := range resultList {
		if i == 0 {
			fmt.Println(strings.Split(v, "\t"))
		}
		//fmt.Println(i, v)
	}

}
