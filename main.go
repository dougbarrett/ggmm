package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "ggmm"
	app.Usage = "Easy way to bootstrap golang web applications"

	app.Commands = []cli.Command{
		{
			Name:  "app",
			Usage: "Application-wide commands",
			Subcommands: []cli.Command{
				{
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "dbusername, dbu",
							Value: "username",
							Usage: "set the db username for the config",
						},
						cli.StringFlag{
							Name:  "dbpassword, dbp",
							Value: "password",
							Usage: "set the db password for the config",
						},
						cli.StringFlag{
							Name:  "dbserver, dbs",
							Value: "localhost",
							Usage: "set the db server for the config",
						},
						cli.StringFlag{
							Name:  "dbport",
							Value: "3306",
							Usage: "set the db port for the config",
						},
						cli.StringFlag{
							Name:  "dbname, dbn",
							Value: "database",
							Usage: "set the db name for the config",
						},
					},
					Name:   "create",
					Usage:  "Create a ggmm application",
					Action: controllerCreateApplication,
				},
				{
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "force",
							Value: "n",
							Usage: "required to delete an application",
						},
						cli.StringFlag{
							Name:  "database",
							Value: "n",
							Usage: "required to delete the database",
						},
					},
					Name:   "delete",
					Action: controllerDeleteApplication,
				},
			},
		},
		{
			Name:  "crud",
			Usage: "Create a create/read/update/delete starter controller, model and view",
			Subcommands: []cli.Command{
				{
					Name:  "create",
					Usage: "create new crud handlers",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "template",
							Value: "",
							Usage: "set a template for the crud to set some nice defaults",
						},
					},
					Action: controllerCreateCrud,
				},
			},
		},
	}

	app.Run(os.Args)
}
