package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
	"unicode"
)

type Entity struct {
	Package string
	Model   string
	Models  string
}
type DomainRepository struct {
	Package string
	Model   string
}
type DomainService struct {
	Package        string
	Model          string
	ModelLowerCase string
}
type DatabaseRepository struct {
	Package string
	Model   string
}
type Interactor struct {
	Package        string
	Model          string
	ModelLowerCase string
}
type Controller struct {
	Package        string
	Model          string
	ModelLowerCase string
}

var FileExtension = ".go"

//  go run generator.go -model="User"
func main() {
	// go run generator.go model="User"
	f := flag.String("model", "None", "This is to generate files based on the model name")
	flag.Parse()
	// fmt.Println(*f)
	words := strings.Fields(*f)

	// fmt.Println(words)

	var model string
	var modelFCharLowerCase string
	if len(words) > 1 {
		for _, v := range words {
			model += strings.Title(strings.ToLower(v))
		}
		modelFCharLowerCase = FirstCharToLowerCase(model)
	} else if len(words) == 1 {
		model = strings.Title(strings.ToLower(*f))
		modelFCharLowerCase = FirstCharToLowerCase(model)
	} else {
		log.Println("Does not accept empty word")
	}
	// createEntityFile(model)
	// createDomainRepositoryFile(model)
	createServiceFile(model, modelFCharLowerCase)
	// createDatabaseRepositoryFile(model)
	// createInteractorFile(model, modelFCharLowerCase)
	// createControllerFile(model, modelFCharLowerCase)
}

func createEntityFile(model string) {
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

	// ファイルへの書き込み
	if err = template.Execute(fp, data); err != nil {
		log.Fatal(err)
	}
}

func createDomainRepositoryFile(model string) {
	data := DomainRepository{
		Package: "entity",
		Model:   model,
	}
	template, err := template.New("domain_repo.tmpl").ParseGlob("domain_repo.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	fileName := strings.ToLower(model) + "_repository" + FileExtension

	// ファイルの生成
	fp, err := os.Create(fileName)
	if err != nil {
		log.Println("error creating file", err)
	}
	defer fp.Close()

	// ファイルへの書き込み
	if err = template.Execute(fp, data); err != nil {
		log.Fatal(err)
	}
}

func createServiceFile(model string, modelLower string) {
	fmt.Println(modelLower)
	data := DomainService{
		Package:        "service",
		Model:          model,
		ModelLowerCase: modelLower,
	}
	template, err := template.New("domain_service.tmpl").ParseGlob("domain_service.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	fileName := data.ModelLowerCase + "_service" + FileExtension

	fp, err := os.Create(fileName)
	if err != nil {
		log.Println("error creating file", err)
	}
	defer fp.Close()

	// ファイルへの書き込み
	if err = template.Execute(fp, data); err != nil {
		log.Fatal(err)
	}
}

func createDatabaseRepositoryFile(model string) {
	data := DatabaseRepository{
		Package: "database",
		Model:   model,
	}
	template, err := template.New("database_repository.tmpl").ParseGlob("database_repository.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	fileName := strings.ToLower(model) + "_repository" + FileExtension

	fp, err := os.Create(fileName)
	if err != nil {
		log.Println("error creating file", err)
	}
	defer fp.Close()

	// ファイルへの書き込み
	if err = template.Execute(fp, data); err != nil {
		log.Fatal(err)
	}
}

func createInteractorFile(model string, modelLower string) {
	data := Interactor{
		Package:        "interactor",
		Model:          model,
		ModelLowerCase: modelLower,
	}
	template, err := template.New("interactor.tmpl").ParseGlob("interactor.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	fileName := data.ModelLowerCase + "_interactor" + FileExtension

	fp, err := os.Create(fileName)
	if err != nil {
		log.Println("error creating file", err)
	}
	defer fp.Close()

	// ファイルへの書き込み
	if err = template.Execute(fp, data); err != nil {
		log.Fatal(err)
	}
}

func createControllerFile(model string, modelLower string) {
	data := Controller{
		Package:        "controller",
		Model:          model,
		ModelLowerCase: modelLower,
	}
	template, err := template.New("controller.tmpl").ParseGlob("controller.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	fileName := data.ModelLowerCase + "_controller" + FileExtension

	fp, err := os.Create(fileName)
	if err != nil {
		log.Println("error creating file", err)
	}
	defer fp.Close()

	// ファイルへの書き込み
	if err = template.Execute(fp, data); err != nil {
		log.Fatal(err)
	}
}

// DI for databse repository

/*
	// load template file
	// template, err := template.New("index.html.tmpl").ParseFiles("index.html.tmpl")
	// ParseFiles -> which template file
	// New
*/

func FirstCharToLowerCase(model string) string {
	s := model
	a := []rune(s)
	a[0] = unicode.ToLower(a[0])
	s = string(a)
	return s
}
