package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"
)

type Language struct {
	Name string
	URL  string
}

//  go run generator.go -model="User"
func main() {

	// load template file
	// template, err := template.New("index.html.tmpl").ParseFiles("index.html.tmpl")
	// ParseFiles -> which template file
	// New

	// go run generator.go model="User"
	f := flag.String("model", "None", "This is to generate files based on the model name")
	flag.Parse()
	fmt.Println(*f)

	data := Language{
		Name: `Go`,
		URL:  `https://golang.org/`,
	}
	template, err := template.New("test.txt").ParseFiles("test.txt")
	if err != nil {
		log.Fatal(err)
	}

	// ファイルの生成
	fp, err := os.Create("index.txt")
	if err != nil {
		log.Println("error creating file", err)
	}
	defer fp.Close()

	// write file
	if err = template.Execute(fp, data); err != nil {
		log.Fatal(err)
	}
}
