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

const (
	fileExtension          = ".go"
	entityPath             = "./../../pkg/domain/entity/"
	domainRepositoryPath   = "./../../pkg/domain/repository/"
	domainServicePath      = "./../../pkg/domain/service/"
	databaseRepositoryPath = "./../../pkg/driver/mysql/database/"
	interactorPath         = "./../../pkg/usecase/"
	controllerPath         = "./../../pkg/adapter/controllers/"
)

//  go run generator.go -model="User"
func main() {
	f := flag.String("model", "", "This is to generate files based on the model name")
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
		return
	}
	createEntityFile(model)
	createDomainRepositoryFile(model)
	createServiceFile(model, modelFCharLowerCase)
	createDatabaseRepositoryFile(model)
	createInteractorFile(model, modelFCharLowerCase)
	createControllerFile(model, modelFCharLowerCase)
}

func createEntityFile(model string) {
	if _, err := os.Stat(entityPath); os.IsNotExist(err) {
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

	fileName := strings.ToLower(model) + fileExtension

	// ファイルの生成
	fp, err := os.Create(filepath.Join(entityPath, fileName))
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
	if _, err := os.Stat(domainRepositoryPath); os.IsNotExist(err) {
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

	fileName := strings.ToLower(model) + "_repository" + fileExtension

	// ファイルの生成
	fp, err := os.Create(filepath.Join(domainRepositoryPath, fileName))
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
	if _, err := os.Stat(domainServicePath); os.IsNotExist(err) {
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

	fileName := data.ModelLowerCase + "_service" + fileExtension

	fp, err := os.Create(filepath.Join(domainServicePath, fileName))
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
	if _, err := os.Stat(databaseRepositoryPath); os.IsNotExist(err) {
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

	fileName := strings.ToLower(model) + "_repository" + fileExtension

	fp, err := os.Create(filepath.Join(databaseRepositoryPath, fileName))
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
	if _, err := os.Stat(interactorPath); os.IsNotExist(err) {
		log.Println("such file or directory does not exist. err = ", err)
		return
	}
	data := Interactor{
		Package:        "usecase",
		Model:          model,
		ModelLowerCase: modelLower,
	}
	tmpl, err := template.New("interactor.tmpl").ParseGlob("interactor.tmpl")
	if err != nil {
		log.Println(err)
		return
	}

	fileName := data.ModelLowerCase + "_interactor" + fileExtension

	fp, err := os.Create(filepath.Join(interactorPath, fileName))
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
	if _, err := os.Stat(controllerPath); os.IsNotExist(err) {
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

	fileName := data.ModelLowerCase + "_controller" + fileExtension

	fp, err := os.Create(filepath.Join(controllerPath, fileName))
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

// error とfile name をマップの中に詰め込む
func FirstCharToLowerCase(model string) string {
	s := model
	a := []rune(s)
	a[0] = unicode.ToLower(a[0])
	s = string(a)
	return s
}
