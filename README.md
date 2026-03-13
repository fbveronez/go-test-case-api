go-test-case-api
================

RESTful API developed in **Go**, with **Docker**, **hot reload**, **migrations**, and **Swagger** documentation.

Technologies
------------

- `Go <https://golang.org/>`_ – main programming language
- `Docker <https://www.docker.com/>`_ – containerization
- `Docker Compose <https://docs.docker.com/compose/>`_ – orchestration
- `Swagger <https://swagger.io/>`_ – API documentation
- `Testify <https://github.com/stretchr/testify>`_ – unit testing
- `Makefile <https://www.gnu.org/software/make/>`_ – automated commands

Features
--------

- RESTful endpoints (GET and POST)
- Account creation and retrieval in the database
- Automated migrations
- Hot reload for fast development
- Swagger documentation at `/swagger/index.html`
- Unit tests using `testify`

Running with Docker
------------------

### Build and run the application



    make run

This command will:

- Build the Docker image of the API
- Start PostgreSQL and API containers
- Enable hot reload for development

### Stop containers



    docker compose down

Hot Reload
----------

During development, code changes automatically reload the application without rebuilding the container:



    make run

> Hot reload is done using a tool like `Air <https://github.com/cosmtrek/air>`_ integrated into the dev container.

Swagger Documentation
---------------------

After running the container, the API documentation is available at:
http://localhost:8080/swagger/api/ui/


> You will find all endpoints, parameters, and request/response examples.

Testing
-------

To run unit tests using `testify`:


    make test

To run functional tests:


    make test-functional


Makefile
--------

Example of main commands:

    run:
        docker compose up --build

    test:
        go test ./internal/... -v

Project Structure
-----------------
```
GO-TEST-CASE-API
│
├── cmd
│   └── api
│       └── main.go
│
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── internal
│   ├── db
│   ├── functional_tests
│   ├── handlers
│   ├── model
│   ├── repository
│   └── service
│
├── migrations
├── .air.toml
├── coverage.out
├── docker-compose.test.yml
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── LICENSE
├── Makefile
└── README.md
```


Tips
----

- Use `make run` during development for hot reload.  
- Use `make new-migration` whenever creating new migrations.  
- Access Swagger documentation to test endpoints without Postman.