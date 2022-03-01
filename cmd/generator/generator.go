package main

import (
	"flag"
	"html/template"
	"log"
	"os"
	"path/filepath"
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

var EntityPath = "./../../pkg/domain/entity/"
var DomainRepositoryPath = "./../../pkg/domain/repository/"
var DomainServicePath = "./../../pkg/domain/service/"
var DatabaseRepositoryPath = "./../../pkg/driver/mysql/database/"
var InteractorPath = "./../../pkg/usecase/"
var ControllerPath = "./../../pkg/adapter/controllers/"

//  go run generator.go -model="User"
func main() {
	f := flag.String("model", "None", "This is to generate files based on the model name")
	flag.Parse()
	words := strings.Fields(*f)

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
	createEntityFile(model)
	createDomainRepositoryFile(model)
	createServiceFile(model, modelFCharLowerCase)
	createDatabaseRepositoryFile(model)
	createInteractorFile(model, modelFCharLowerCase)
	createControllerFile(model, modelFCharLowerCase)
}

func createEntityFile(model string) {
	if _, err := os.Stat(EntityPath); os.IsNotExist(err) {
		log.Println("such file or directory does not exist. err = ", err)
		return
	}
	data := Entity{
		Package: "entity",
		Model:   model,
		Models:  model + "s",
	}
	tmpl, err := template.New("entity.tmpl").ParseFiles("entity.tmpl")
	if err != nil {
		log.Println(err)
		return
	}

	fileName := strings.ToLower(model) + FileExtension

	// ファイルの生成
	fp, err := os.Create(filepath.Join(EntityPath, fileName))
	if err != nil {
		log.Println("error creating file", err)
	}
	defer fp.Close()

	// ファイルへの書き込み
	if err = tmpl.Execute(fp, data); err != nil {
		log.Println(err)
		return
	}
}

func createDomainRepositoryFile(model string) {
	if _, err := os.Stat(DomainRepositoryPath); os.IsNotExist(err) {
		log.Println("such file or directory does not exist. err = ", err)
		return
	}
	data := DomainRepository{
		Package: "repository",
		Model:   model,
	}
	tmpl, err := template.New("domain_repo.tmpl").ParseGlob("domain_repo.tmpl")
	if err != nil {
		log.Println(err)
		return
	}

	fileName := strings.ToLower(model) + "_repository" + FileExtension

	// ファイルの生成
	fp, err := os.Create(filepath.Join(DomainRepositoryPath, fileName))
	if err != nil {
		log.Println("error creating file", err)
	}
	defer fp.Close()

	// ファイルへの書き込み
	if err = tmpl.Execute(fp, data); err != nil {
		log.Println(err)
		return
	}
}

func createServiceFile(model, modelLower string) {
	if _, err := os.Stat(DomainServicePath); os.IsNotExist(err) {
		log.Println("such file or directory does not exist. err = ", err)
		return
	}
	data := DomainService{
		Package:        "service",
		Model:          model,
		ModelLowerCase: modelLower,
	}
	tmpl, err := template.New("domain_service.tmpl").ParseGlob("domain_service.tmpl")
	if err != nil {
		log.Println(err)
		return
	}

	fileName := data.ModelLowerCase + "_service" + FileExtension

	fp, err := os.Create(filepath.Join(DomainServicePath, fileName))
	if err != nil {
		log.Println("error creating file", err)
	}
	defer fp.Close()

	// ファイルへの書き込み
	if err = tmpl.Execute(fp, data); err != nil {
		log.Println(err)
		return
	}
}

func createDatabaseRepositoryFile(model string) {
	if _, err := os.Stat(DatabaseRepositoryPath); os.IsNotExist(err) {
		log.Println("such file or directory does not exist. err = ", err)
		return
	}
	data := DatabaseRepository{
		Package: "database",
		Model:   model,
	}
	tmpl, err := template.New("database_repository.tmpl").ParseGlob("database_repository.tmpl")
	if err != nil {
		log.Println(err)
		return
	}

	fileName := strings.ToLower(model) + "_repository" + FileExtension

	fp, err := os.Create(filepath.Join(DatabaseRepositoryPath, fileName))
	if err != nil {
		log.Println("error creating file", err)
	}
	defer fp.Close()

	// ファイルへの書き込み
	if err = tmpl.Execute(fp, data); err != nil {
		log.Println(err)
		return
	}
}

func createInteractorFile(model, modelLower string) {
	if _, err := os.Stat(InteractorPath); os.IsNotExist(err) {
		log.Println("such file or directory does not exist. err = ", err)
		return
	}
	data := Interactor{
		Package:        "interactor",
		Model:          model,
		ModelLowerCase: modelLower,
	}
	tmpl, err := template.New("interactor.tmpl").ParseGlob("interactor.tmpl")
	if err != nil {
		log.Println(err)
		return
	}

	fileName := data.ModelLowerCase + "_interactor" + FileExtension

	fp, err := os.Create(filepath.Join(InteractorPath, fileName))
	if err != nil {
		log.Println("error creating file", err)
	}
	defer fp.Close()

	// ファイルへの書き込み
	if err = tmpl.Execute(fp, data); err != nil {
		log.Println(err)
		return
	}
}

func createControllerFile(model, modelLower string) {
	if _, err := os.Stat(ControllerPath); os.IsNotExist(err) {
		log.Println("such file or directory does not exist. err = ", err)
		return
	}
	data := Controller{
		Package:        "controllers",
		Model:          model,
		ModelLowerCase: modelLower,
	}
	tmpl, err := template.New("controller.tmpl").ParseGlob("controller.tmpl")
	if err != nil {
		log.Println(err)
		return
	}

	fileName := data.ModelLowerCase + "_controller" + FileExtension

	fp, err := os.Create(filepath.Join(ControllerPath, fileName))
	if err != nil {
		log.Println("error creating file", err)
	}
	defer fp.Close()

	// ファイルへの書き込み
	if err = tmpl.Execute(fp, data); err != nil {
		log.Println(err)
		return
	}
}

// DI for databse repository

/*
	// load template file
	// template, err := template.New("index.html.tmpl").ParseFiles("index.html.tmpl")
	// ParseFiles -> which template file
	// New
*/

// error とfile name をマップの中に詰め込む
func FirstCharToLowerCase(model string) string {
	s := model
	a := []rune(s)
	a[0] = unicode.ToLower(a[0])
	s = string(a)
	return s
}
