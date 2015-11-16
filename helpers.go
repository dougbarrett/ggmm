package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/codegangsta/cli"
)

func randString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func loadConfig(dir string) (config *GgmmConfig) {
	b, err := ioutil.ReadFile(fmt.Sprintf("%sconfig.json", dir))

	if err != nil {
		log.Fatalf("Cannot load config: %s", err)
	}

	err = json.Unmarshal(b, &config)

	if err != nil {
		log.Fatalf("config.json not in correct format: %s", err)
	}
	return
}

func createHomeController(c *cli.Context, config *GgmmConfig) {
	var homeControllerFile string = fmt.Sprintf("%s/controllers/home.go", c.Args().First())

	var output bytes.Buffer

	t, err := template.New("ControllerHomeGo").Parse(controllersHomeGo)

	if err != nil {
		log.Fatalf("Cannot load controller home template: %s", err)
	}

	err = t.Execute(&output, config)

	if err != nil {
		log.Fatalf("Cannot execute controller home template: %s", err)
	}

	createFile(homeControllerFile, output.Bytes())

	createTemplate("home", templateTemplatesHomeTmpl, c, config, GgmmCrudController{}, true)

	var controllerExtraFile string = fmt.Sprintf("%s/controllers/extra.go", c.Args().First())

	createFile(controllerExtraFile, []byte(controllerExtraGo))

	// createUserBootstraps(c, config)
	updateBootstraps(fmt.Sprintf("%s/", c.Args().First()), config)
}

func createUserBootstraps(c *cli.Context, config *GgmmConfig) {
	var controllerUserFile string = fmt.Sprintf("%s/controllers/user.go", c.Args().First())

	var output bytes.Buffer

	t, err := template.New("ControllerUser").Parse(controllerUserGo)

	if err != nil {
		log.Fatalf("Cannot load user controller bootstrap template: %s", err)
	}

	err = t.Execute(&output, config)

	if err != nil {
		log.Fatalf("Cannot execute user controller bootstrap template: %s", err)
	}

	createFile(controllerUserFile, output.Bytes())
}

func updateBootstraps(appname string, config *GgmmConfig) {
	var controllerBootstrapFile string = fmt.Sprintf("%scontrollers/bootstrap.go", appname)
	var output bytes.Buffer

	t, err := template.New("ControllerBootstrap").Parse(controllersBootstrapGo)

	if err != nil {
		log.Fatalf("Cannot load controller bootstrap template: %s", err)
	}

	err = t.Execute(&output, config)

	if err != nil {
		log.Fatalf("Cannot execute controller bootstrap template: %s", err)
	}

	createFile(controllerBootstrapFile, output.Bytes())

	var controllerModelBootstrapFile string = fmt.Sprintf("%smodels/bootstrap.go", appname)
	output = bytes.Buffer{}

	t, err = template.New("ModelBootstrap").Parse(modelsBootstrapGo)

	if err != nil {
		log.Fatalf("Cannot load model bootstrap template: %s", err)
	}

	err = t.Execute(&output, config)

	if err != nil {
		log.Fatalf("Cannot execute model bootstrap template: %s", err)
	}

	createFile(controllerModelBootstrapFile, output.Bytes())
}

func saveConfig(c *cli.Context, config *GgmmConfig, dir string) {
	b, err := json.MarshalIndent(config, "", "\t")

	if err != nil {
		log.Fatalf("Cannot read config correctly: %s", err)
	}
	createFile(fmt.Sprintf("%sconfig.json", dir), b)
}

func createFile(filename string, data []byte) {
	f, err := os.Create(filename)

	if err != nil {
		log.Fatalf("Cannot create file: %s", err)
	}

	_, err = f.Write(data)

	if err != nil {
		log.Fatalf("Cannot save file: %s", err)
	}
}

func createMainFile(dirName string, config *GgmmConfig) {
	var output bytes.Buffer

	t, err := template.New("MainFile").Parse(mainGo)

	if err != nil {
		log.Fatalf("Cannot load main.go template: %s", err)
	}

	err = t.Execute(&output, config)

	if err != nil {
		log.Fatalf("Cannot execute main.go template: %s", err)
	}

	createFile(fmt.Sprintf("%smain.go", dirName), output.Bytes())
}

func applyDefaultFields(crud *GgmmCrudController) {
	var ID GgmmCrudControllerField
	ID.Name = "ID"
	ID.Type = "uint"
	ID.Config = `form:"id"`
	ID.Form.Hidden = true

	var CreatedAt GgmmCrudControllerField
	CreatedAt.Name = "CreatedAt"
	CreatedAt.Type = "time.Time"
	CreatedAt.Config = `form:"created_at"`
	CreatedAt.Form.Hidden = true
	CreatedAt.Form.Type = "string"

	var UpdatedAt GgmmCrudControllerField
	UpdatedAt.Name = "UpdatedAt"
	UpdatedAt.Type = "time.Time"
	UpdatedAt.Config = `form:"updated_at"`
	UpdatedAt.Form.Hidden = true
	UpdatedAt.Form.Type = "string"

	var DeletedAt GgmmCrudControllerField
	DeletedAt.Name = "DeletedAt"
	DeletedAt.Type = "time.Time"
	DeletedAt.Config = `form:"deleted_at"`
	DeletedAt.Form.Hidden = true
	DeletedAt.Form.Type = "string"

	crud.Model = append(crud.Model, ID)
	crud.Model = append(crud.Model, CreatedAt)
	crud.Model = append(crud.Model, UpdatedAt)
	crud.Model = append(crud.Model, DeletedAt)
}

func createUserCrud(c *cli.Context, config *GgmmConfig) {
	log.Println("Creating user crud")

	var user GgmmCrudController
	user.Name = strings.Title(c.Args().First())
	config.LoginCrud = user.Name
	user.Routed = strings.ToLower(c.Args().First())
	config.LoginCrudLower = user.Routed
	user.Template = "user"

	applyDefaultFields(&user)

	var Name GgmmCrudControllerField
	Name.Name = "Name"
	Name.Type = "string"
	Name.Config = `form:"name"`
	Name.Form.Type = "text"

	var Email GgmmCrudControllerField
	Email.Name = "Email"
	Email.Type = "string"
	Email.Config = `form:"email" sql:"unique_key"`
	Email.Form.Type = "email"

	var Password GgmmCrudControllerField
	Password.Name = "Password"
	Password.Type = "string"
	Password.Config = `form:"password"`
	Password.Form.Type = "password"

	var PasswordRepeat GgmmCrudControllerField
	PasswordRepeat.Name = "PasswordRepeat"
	PasswordRepeat.Type = "string"
	PasswordRepeat.Config = `form:"password_repeat" sql:"-"`
	PasswordRepeat.Form.Type = "password"

	user.Model = append(user.Model, Name)
	user.Model = append(user.Model, Email)
	user.Model = append(user.Model, Password)
	user.Model = append(user.Model, PasswordRepeat)

	config.CrudControllers = append(config.CrudControllers, user)
}

func updateCruds(c *cli.Context, config *GgmmConfig) {
	createCrudController(c, config)
	createCrudModel(c, config)
	updateBootstraps("", config)
	createCrudTemplates(c, config)
}

func createCrudTemplates(c *cli.Context, config *GgmmConfig) {
	for _, ctrl := range config.CrudControllers {
		if ctrl.Template == "user" {
			createTemplate("siteModals", templateTemplatesSiteModalsTmpl, c, config, ctrl, false)
			createTemplate("generated_left_nav", templateTemplatesGeneratedLeftNavTmpl, c, config, ctrl, false)
			createTemplate("Login", templateTemplatesLoginTmpl, c, config, ctrl, false)
			createTemplate("Register", templateTemplatesRegisterTmpl, c, config, ctrl, false)
			createTemplate("layout", templateTemplatesLayoutTmpl, c, config, ctrl, false)
		}
	}
}

func createTemplate(templateName string, templateSource string, c *cli.Context, config *GgmmConfig, ctrl GgmmCrudController, isRoot bool) {
	t, err := template.New(fmt.Sprintf("TemplateFile%s", templateName)).Parse(templateSource)

	if err != nil {
		log.Fatalf("Error parsing template: %s", err)
	}

	var retData struct {
		Config     *GgmmConfig
		Controller GgmmCrudController
	}

	retData.Config = config
	retData.Controller = ctrl

	var output bytes.Buffer

	err = t.Execute(&output, retData)

	if err != nil {
		log.Fatalf("Error executing template: %s", err)
	}

	bOutput := output.Bytes()

	bOutput = bytes.Replace(bOutput, []byte("[["), []byte("{{"), -1)
	bOutput = bytes.Replace(bOutput, []byte("]]"), []byte("}}"), -1)

	var fileName string

	if isRoot {
		fileName = c.Args().First() + "/"
	}

	fileName += fmt.Sprintf("templates/%s.tmpl", templateName)

	createFile(fileName, bOutput)
}

func createCrudController(c *cli.Context, config *GgmmConfig) {
	for _, c := range config.CrudControllers {
		_, err := ioutil.ReadFile(fmt.Sprintf("controllers/%s.go", strings.ToLower(c.Name)))
		if err == nil {
			log.Printf("%s already created, not creating again", c.Name)
			continue
		}

		var output bytes.Buffer

		t, err := template.New(fmt.Sprintf("ControllerTemplate%s", c.Name)).Parse(controllersCrudGo)

		if err != nil {
			log.Fatalf("Error parsing crud controller template: %s", err)
		}

		var retData struct {
			Config *GgmmConfig
			Crud   GgmmCrudController
		}

		retData.Config = config
		retData.Crud = c

		err = t.Execute(&output, retData)

		if err != nil {
			log.Fatalf("Error executing crud controller template: %s", err)
		}

		createFile(fmt.Sprintf("controllers/%s.go", strings.ToLower(c.Name)), output.Bytes())
	}

}

func createCrudModel(c *cli.Context, config *GgmmConfig) {
	for _, c := range config.CrudControllers {
		_, err := ioutil.ReadFile(fmt.Sprintf("models/%s.go", strings.ToLower(c.Name)))
		if err == nil {
			log.Printf("%s already created, not creating again", c.Name)
			continue
		}

		var output bytes.Buffer

		t, err := template.New(fmt.Sprintf("ModelTemplate%s", c.Name)).Parse(modelsCrudGo)

		if err != nil {
			log.Fatalf("Error parsing crud model template: %s", err)
		}

		err = t.Execute(&output, c)

		if err != nil {
			log.Fatalf("Error executing crud model template: %s", err)
		}

		createFile(fmt.Sprintf("models/%s.go", strings.ToLower(c.Name)), output.Bytes())
	}
}

func createCrudTemplate(c *cli.Context, config *GgmmConfig) {
	for _, c := range config.CrudControllers {
		_, err := ioutil.ReadFile(fmt.Sprintf("templates/%s/%s/all.tmpl", config.ApplicationName, strings.ToLower(c.Name)))

		if err == nil {
			log.Printf("Templates already created for %s", c.Name)
			continue
		}

	}
}
