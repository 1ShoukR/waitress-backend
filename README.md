# waitress-backend

## Introduction

This is the project repository for Waotress.

## Getting Started

The primary source of information about this project is its documentation. This README is intended to only get you as far as building that documentation locally. From there you can switch to the docs to further learn how to run and develop locally.

1. Clone this repository, and navigate into it. 

       $ git clone https://github.com/1ShoukR/waitress-backend.git
       $ cd waitress-backend

1. Create and activate a virtual environment. Make sure you do this with a compatible version of Python (preferably 3.8). 

    Unix Example:

        $ python3 -m venv venv
        $ source ./venv/Scripts/activate

    Windows Example:

        $ python -m venv venv
        $ .\venv\Scripts\activate

    Git Bash Example
        $ python -m venv venv
        $ venv\Scripts\activate

1. Make sure your virtual environment has the latest version of Pip. 

        $ python -m pip install --upgrade pip


1. Install this project into your virtual environment.

        $ pip install -r requirements.txt

    Be sure to install `requirements.txt`, otherwise you won't get development dependencies such as Sphinx.

1. Navigate to the docs folder, and build the docs using Sphinx.

        $ cd docs
        $ sphinx-build -M html docs/source/ docs/build/

1. The `sphinx-build -M html docs/source/ docs/build/` script should generate a `build` directory. Open `build/html/index.html` to access the docs.
