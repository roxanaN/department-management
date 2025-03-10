
# Department Management Application

## Description
The Department Management application is a web solution designed to manage department hierarchies and their associated users. It allows obtaining the complete department hierarchy, as well as listing users for a specific department. The application is managed by the admin/

## Prerequisites

Before running this application, ensure you have the following installed on your system:
 - Docker
 - Docker Compose


## Running the Application Using Docker Compose

**Start the application**: ```docker-compose up --build```. This command builds and starts the application in a container.

**Stop the application**: ```docker-compose down -v```

## Project Structure

.
├── Dockerfile
├── README.md
├── api
│   ├── api.go
│   ├── handlers
│   │   ├── auth.go
│   │   ├── departments.go
│   │   └── users.go
│   └── middleware
│       └── auth.go
├── cmd
│   └── server
│       └── main.go
├── config
│   └── config.go
├── db
│   └── db.go
├── docker-compose.yml
├── go.mod
├── go.sum
├── init_db.sql
├── models
│   └── models.go
├── stored_procedures.sql
└── utils
    └── utils.go

## Exposed Ports

8080: The application runs on this port inside the container and is mapped to the host.

## Notes

Ensure that go.mod and go.sum are updated if dependencies change.