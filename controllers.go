package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func controllerDeleteApplication(c *cli.Context) {
	if c.String("force") == "yes" {
		var applicationDirectory = fmt.Sprintf("%s", c.Args().First())
		if c.String("database") == "yes" {
			config := loadConfig(c.Args().First() + "/")
			dbConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
				config.Database.Username,
				config.Database.Password,
				config.Database.Server,
				config.Database.Port,
				config.Database.Database)

			db, err := gorm.Open("mysql", dbConn)
			if err == nil {

				err := db.Exec("SET foreign_key_checks = 0").Error

				if err != nil {
					log.Println("Cannot clear out databases: %s", err)
					return
				}

				rows, err := db.DB().Query("SHOW TABLES")

				if err != nil {
					log.Println("cannot get list of tables", err)
					return
				}

				for rows.Next() {
					var name string
					rows.Scan(&name)
					err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", name)).Error
					if err != nil {
						log.Println("Cannot drop table: %s", err)
					}
				}

				err = db.Exec("SET foreign_key_checks = 1").Error

				if err != nil {
					log.Println("Cannot clear out databases: %s", err)
					return
				}
			} else {
				log.Println("Cannot load mysql: %s", err)
			}
		}
		os.RemoveAll(applicationDirectory)
		log.Println("Application was deleted")
	} else {
		log.Println("Application not deleted, force was not used.  Use the force, Luke.")
	}

}

func controllerCreateApplication(c *cli.Context) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to get working directory: %s", pwd)
	}
	repo := strings.Split(pwd, "/src/")
	log.Printf("working directory: %s", pwd)
	if len(repo) != 2 {
		log.Fatalf("Unable to figure out the repo, please make sure you are working under your gopath src directory")
	}
	log.Printf("working repo: %s", repo[1])
	log.Printf("Creating application: %s", c.Args().First())
	log.Println("Creating application directory")

	// Create directories
	var applicationDirectory = fmt.Sprintf("%s", c.Args().First())
	var controllersDirectory = fmt.Sprintf("%s/controllers", c.Args().First())
	var modelsDirectory = fmt.Sprintf("%s/models", c.Args().First())
	var templatesDirectory = fmt.Sprintf("%s/templates", c.Args().First())
	os.MkdirAll(applicationDirectory, 0777)
	log.Println("Creating controllers directory")
	os.MkdirAll(controllersDirectory, 0777)
	log.Println("Creating models directory")
	os.MkdirAll(modelsDirectory, 0777)

	log.Println("Creating templates directory")
	os.MkdirAll(templatesDirectory, 0777)

	// Creating config file
	var config GgmmConfig
	config.CurrentRepo = repo[1] + "/" + c.Args().First()

	config.ApplicationName = c.Args().First()
	var homeController GgmmController
	homeController.Function = "controllerHome"
	homeController.Route = "/"
	config.GetControllers = append(config.GetControllers, homeController)
	config.Database.Username = c.String("dbusername")
	config.Database.Password = c.String("dbpassword")
	config.Database.Server = c.String("dbserver")
	config.Database.Port = c.String("dbport")
	config.Database.Database = c.String("dbname")
	config.SessionKey = randString(32)

	saveConfig(c, &config, fmt.Sprintf("%s/", c.Args().First()))
	createMainFile(c.Args().First()+"/", &config)
	createHomeController(c, &config)

	createTemplate("layout", templateTemplatesLayoutTmpl, c, &config, GgmmCrudController{}, true)
	createTemplate("generated_left_nav", templateTemplatesGeneratedLeftNavTmpl, c, &config, GgmmCrudController{}, true)
	createTemplate("siteModals", templateTemplatesSiteModalsTmpl, c, &config, GgmmCrudController{}, true)
}

func controllerCreateCrud(c *cli.Context) {
	config := loadConfig("")

	if c.Args().First() == "" {
		log.Println("model type required")
		return
	}

	for _, x := range config.CrudControllers {
		if "user" == strings.ToLower(c.String("template")) {
			log.Fatal("Users CRUD template has already been used")
		}
		if x.Name == c.Args().First() {
			log.Fatal("Cannot create CRUD, one of the same name already exists")
		}
	}

	switch {
	case "user" <= c.String("template"):
		log.Println("creating crud based off users")

		config.AllowRegistration = true
		createMainFile("", config)
		saveConfig(c, config, "")

		createUserCrud(c, config)
	case "" <= c.String("template"):
		log.Println("creating default crud")
	}

	updateCruds(c, config)
	saveConfig(c, config, "")
}
