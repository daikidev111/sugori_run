package main

import (
	"flag"
	"html/template"
	"log"
	"os"
	"strings"
)

type Entity struct {
	Package string
	Model   string
	Models  string
}

var FileExtension = ".go"

//  go run generator.go -model="User"
func main() {

	// go run generator.go model="User"
	f := flag.String("model", "None", "This is to generate files based on the model name")
	flag.Parse()

	createEntityFile(*f)
}

func createEntityFile(model string) {
	model = strings.Title(strings.ToLower(model))

	data := Entity{
		Package: "entity",
		Model:   model,
		Models:  model + "s",
	}
	template, err := template.New("entity.tmpl").ParseFiles("entity.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	fileName := strings.ToLower(model) + FileExtension

	// ファイルの生成
	fp, err := os.Create(fileName)
	if err != nil {
		log.Println("error creating file", err)
	}
	defer fp.Close()

	// write file
	if err = template.Execute(fp, data); err != nil {
		log.Fatal(err)
	}
}

//　関数化する
// entity の構造体を作る
// make the first character capitalized

/*
	// load template file
	// template, err := template.New("index.html.tmpl").ParseFiles("index.html.tmpl")
	// ParseFiles -> which template file
	// New
*/
