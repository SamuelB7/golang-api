# Golang REST API

This is a simple REST API built with Golang.

## Prerequisites

- Docker
- Docker Compose

## How to Run the Project

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

   The API will start at port `8080`.

## API Endpoints

The following endpoints are available:

- **Create User:**
  - **POST** `/users`
- **Get All Users:**
  - **GET** `/users`
- **Get One User:**
  - **GET** `/users/{id}`
- **Update User:**
  - **PUT** `/users/{id}`
- **Delete User:**
  - **DELETE** `/users/{id}`
