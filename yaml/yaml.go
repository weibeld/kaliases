package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Token struct {
	Short, Long string
}

type Segment struct {
	Single   []Token
	Multiple []Token
}

func main() {
	//var s string
	var a []Segment
	//m := make(map[string][]string)
	//var s MyNewStruct
	readYAML("test.yaml", &a)
	fmt.Println(a)
	spew.Dump(a)
}

func readYAML(filename string, out interface{}) {
	data := readFile(filename)
	err := yaml.Unmarshal(data, out)
	if err != nil {
		log.Fatal(err)
	}
}

func readFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
