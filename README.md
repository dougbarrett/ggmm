# GGMM

# This project is not done and under heavy development, do not use on production sites yet

### What is GGMM?

GGMM is a tool used to easily create the bones of a Go web application.  It uses the following libraries:

* Gorm (for database handling)
* Go (for the language)
* Martini (for the web framework)
* MySQL (for the database)

### Requirements

- GGMM requires Go 1.5.1 or newer
- MySQL is required

### Getting Started

Run the following:

````
go get github.com/dougbarrett/ggmm
````

Now in your repository, run the following:

````
ggmm app create testapp
````

*testapp* should be replace with the name of your applcation.

After that, go into the directory and you'll see some files have been created.  In here, configure your config.json file to your correct database settings.

When you're all set, just run `go run *.go` or use your favorite live reload tool!

### Documentation

#### Create CRUD 'login' user

To create the logic that can handle user registration and login, run the following in your application directory:

````
ggmm crud create user --template user
````

This will create a 'user' model and by telling it you want to use the 'user' template, it will set that model up to be used for login and user management.

### Example commands to create first app

````
ggmm app create userCrud --dbu username --dbp password --dbname userCrud
cd userCrud
ggmm crud create user --template user
````

Replace 'username' and 'password' with your database login credentials, this will override the default settings

### TODO

For v1

- Create model handler generators
- Create controller handler generators
- Create normal-crud handler generators
- Add tests
- Improve documentation on using ggmm and how to use the generated code

For v1.2

- Add adons: stripe payments, oauth.io support, mailgun support

For v2

- Open up ability to use more than just martini (gin, gin-gonic, gorilla mux, standard net/http lib)