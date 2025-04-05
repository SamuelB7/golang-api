# DevBook API

A simple social media REST API built with Go (Golang) that allows users to create accounts, make posts, and interact with other users.

## Prerequisites

- Docker
- Docker Compose

## Getting Started

1. **Clone the repository:**

   ```sh
   git clone git@github.com:SamuelB7/golang-api.git
   cd golang-api
   ```

2. **Set up environment variables:**

   - Using the `.env.example` file, create a new `.env` file and add the necessary environment variables.

   ```sh
   cp .env.example .env
   ```

3. **Run the project:**

   ```sh
   docker compose up
   ```

   The API will be available at:

   `http://localhost:8080/api`

## API Documentation

Swagger UI

Access the interactive API documentation at:

`http://localhost:8080/swagger/index.html`

Postman/Insomnia

To import the API specification into Postman or other API tools:

`http://localhost:8080/swagger/doc.json`

## Development

Generate API Documentation

After making changes to the API endpoints, regenerate the Swagger documentation:

```sh
 swag init
```
