package main

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v2"
)

func main() {
	var data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
c:
  a: 100
  b: 200
  test: [1, 2]
`
	// yaml tags: https://godoc.org/gopkg.in/yaml.v2
	type T struct {
		A string
		B struct {
			RenamedC int   `yaml:"c"`
			D        []int `yaml:"d,flow"`
		}
		C struct {
			RA int    `yaml:"a"`
			RB string `yaml:"b"`
			RC []int  `yaml:"test,flow"`
		}
	}

	t := T{}

	fmt.Println("Hello world, my first go mod, I use yaml pkg here!")

	err := yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//fmt.Printf("--- t:\n%v\n\n", t)
	fmt.Println(t)
	/*
		d, err := yaml.Marshal(&t)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		fmt.Printf("--- t dump:\n%s\n\n", string(d))

		m := make(map[interface{}]interface{})

		err = yaml.Unmarshal([]byte(data), &m)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		fmt.Printf("--- m:\n%v\n\n", m)

		d, err = yaml.Marshal(&m)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		fmt.Printf("--- m dump:\n%s\n\n", string(d))
	*/
}
