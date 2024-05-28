# waitress-backend

## Introduction

Welcome to Waitress. This is the backend repository for the Waitress project. This project is an application that allows users to order food from a restaurant, reserve their favorite seat, and have their food to their table in no time. The backend is built using the Golang, while utilizing the Gin Framework. The mobile application is built using React Native. The project is currently in development.

NOTE: This repository contains a Python Flask implementation of the backend. The Golang implementation is currently in development, and will be the main application that we shall be using. The Python implementation is only for reference purposes, as it was the origial implementation of the backend.

Please create a .env file in the root directory of the project, and ask the project owner for the environment variables that you need to add to the .env file.

## Getting Started With Golang

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.
    - Install the latest version of Golang [here](https://golang.org/dl/)
    - Verify that Golang is installed by running the following command in your terminal:
        ```
        go version
        ```
    - Clone the repository to your local machine:
        ```
        git clone https://github.com/1ShoukR/waitress-backend.git
        ```
    - Navigate to the project directory:
        ```
        cd waitress-backend
        ```
    - Install the project dependencies:
        ```
        go mod download
        ```

## MySQL Database

This project uses a MySQL database to store and manage data. You are required to install MySQL on your local machine, and create a database for the project. You can install MySQL by visiting this [link](https://dev.mysql.com/downloads/mysql/). It is recommended to install MySQL Workbench as well, to manage the database.


## MakeFile

This project uses a Makefile to automate the build process, and to make it easier to run the application.

If you cannot install or use Make commands, you can run the application normally by running the following commands:

To build the project, run: 
```
go build cmd/api/main.go
```

To run the project, run:
```
go run cmd/api/main.go
```


NOTE: If you are using Git Bash on Windows, you may need to install Make. You can install Make by visiting this stackoverflow [link](https://stackoverflow.com/questions/36770716/mingw64-make-build-error-bash-make-command-not-found).

Ensure that you have the latest version of Make installed on your machine. You can check if Make is installed by running the following command in your terminal:

```
make --version
```

The Makefile contains the following commands:

```
make all build
```

build the application

```
make build
```

run the application

```
make run
```

Create DB container

```
make docker-run
```

Shutdown DB container

```
make docker-down
```

live reload the application

```
make watch
```

run the test suite

```
make test
```

clean up binary from the last build

```
make clean
```
