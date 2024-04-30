# waitress-backend

## Introduction

This is the project repository for Waotress.

## Getting Started

The primary source of information about this project is its documentation. This README is intended to only get you as far as building that documentation locally. From there you can switch to the docs to further learn how to run and develop locally.

1. Clone this repository, and navigate into it.

       git clone https://github.com/1ShoukR/waitress-backend.git
       cd waitress-backend

1. Create and activate a virtual environment. Make sure you do this with a compatible version of Python (preferably 3.8).

    Unix Example:

        python3 -m venv venv
        source ./venv/Scripts/activate

    Windows Example:

        python -m venv venv
        .\venv\Scripts\activate

    Git Bash Example
        $ python -m venv venv
        $ venv\Scripts\activate

1. Make sure your virtual environment has the latest version of Pip.

        python -m pip install --upgrade pip

1. Install this project into your virtual environment.

        pip install -r requirements.txt

    Be sure to install `requirements.txt`, otherwise you won't get development dependencies such as Sphinx.

1. Navigate to the docs folder, and build the docs using Sphinx.

        cd docs
        sphinx-build -M html docs/source/ docs/build/

1. The `sphinx-build -M html docs/source/ docs/build/` script should generate a `build` directory. Open `build/html/index.html` to access the docs.

## Golang Rewrite

- This application is being rewritten in Goalng.
- This is the MAIN branch that will hold the code that will be eventually merged into the `main` branch 

## Getting Started With Golang

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.
    - Install the latest version of Golang.
    ( Add more instructions here )
## MakeFile

run all make commands with clean tests

```bash
make all build
```

build the application

```bash
make build
```

run the application

```bash
make run
```

Create DB container

```bash
make docker-run
```

Shutdown DB container

```bash
make docker-down
```

live reload the application

```bash
make watch
```

run the test suite

```bash
make test
```

clean up binary from the last build

```bash
make clean
```
