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

	createEntityFile(*f)
	createDomainRepositoryFile(*f)
	createServiceFile(*f)
	createDatabaseRepositoryFile(*f)
	createInteractorFile(*f)
	createControllerFile(*f)
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

	// ファイルへの書き込み
	if err = template.Execute(fp, data); err != nil {
		log.Fatal(err)
	}
}

func createDomainRepositoryFile(model string) {
	model = strings.Title(strings.ToLower(model))

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

func createServiceFile(model string) {
	model = strings.Title(strings.ToLower(model))
	modelLowerCase := strings.ToLower(model)
	data := DomainService{
		Package:        "service",
		Model:          model,
		ModelLowerCase: modelLowerCase,
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
	model = strings.Title(strings.ToLower(model))
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

func createInteractorFile(model string) {
	model = strings.Title(strings.ToLower(model))
	modelLowerCase := strings.ToLower(model)
	data := Interactor{
		Package:        "interactor",
		Model:          model,
		ModelLowerCase: modelLowerCase,
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

func createControllerFile(model string) {
	model = strings.Title(strings.ToLower(model))
	modelLowerCase := strings.ToLower(model)
	data := Controller{
		Package:        "controller",
		Model:          model,
		ModelLowerCase: modelLowerCase,
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
