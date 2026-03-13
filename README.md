go-test-case-api
================

RESTful API developed in **Go**, with **Docker**, **hot reload**, **migrations**, and **Swagger** documentation. Designed to be modular, professional, and easy to maintain.

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

Migrations
----------

The project uses migrations to manage the PostgreSQL database.

- Apply all migrations:



    make migrate-up

- Rollback all migrations:



    make migrate-down

> Migrations are located in the `migrations/` folder.

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
go-test-case-api/
├─ cmd/ # Main application entrypoint
├─ internal/
│ ├─ service/ # Business logic
│ ├─ repository/ # Database access
│ └─ handler/ # HTTP handlers
├─ migrations/ # Migration files
├─ swagger/ # Swagger documentation files
├─ Dockerfile
├─ docker-compose.yml
├─ Makefile
└─ go.mod
```


Tips
----

- Use `make run` during development for hot reload.  
- Use `make new-migration` whenever creating new migrations.  
- Access Swagger documentation to test endpoints without Postman.