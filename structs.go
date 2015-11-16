package main

type GgmmConfig struct {
	ApplicationName   string
	CurrentRepo       string
	CurrentRepository string
	GetControllers    []GgmmController
	CrudControllers   []GgmmCrudController
	Database          struct {
		Username string
		Password string
		Server   string
		Port     string
		Database string
	}
	LoginCrud         string
	LoginCrudLower    string
	AllowRegistration bool
	SessionKey        string
}

type GgmmController struct {
	Route    string
	Function string
}

type GgmmCrudController struct {
	Name     string
	Routed   string
	Template string
	Model    []GgmmCrudControllerField
}

type GgmmCrudControllerField struct {
	Name   string
	Type   string
	Config string
	Form   struct {
		Hidden        bool
		Type          string
		SelectOptions []string
	}
}
