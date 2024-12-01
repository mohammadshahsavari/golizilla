# Golizilla API

This repository contains the Golizilla API, a simple REST API built using Go and Fiber. This guide will help you get started running the API either with Docker or locally. Additionally, you'll find instructions for using Swagger for API documentation and testing.

## Prerequisites

- **Go** (version 1.20 or higher)
- **Docker** (if you want to run with Docker)

## Running the App

### Option 1: Running with Docker

This option will spin up the API and a PostgreSQL database using Docker Compose. Follow these steps:

1. **Ensure Docker is installed and running**.

2. **Copy the example `.env` file**:
   ```sh
   cp .env.example .env
   ```

3. **Edit `.env` File** (if needed): Make sure that the `.env` configuration is suitable for your Docker environment. Leave `DB_HOST=golizilla-db` for Docker Compose to work properly.

4. **Run Docker Compose**:
   ```sh
   docker-compose up -d --build
   ```

   This will build and start both the `golizilla-api` service and the `golizilla-db` service.

5. **Access the API**:
   The API should be running on `http://localhost:8080`.

6. **Access Swagger Documentation**:
   Swagger documentation is available at `http://localhost:8080/swagger/index.html`.
   
   **Note**: Swagger documentation is accessible only when `ENV` is not set to production in the `.env` file.

### Option 2: Running Locally with Docker for Database Only

In this option, you will run the PostgreSQL database using Docker and run the API locally.

1. **Ensure Docker is installed and running**.

2. **Copy the example `.env` file**:
   ```sh
   cp .env.example .env
   ```

3. **Edit `.env` File**:
   - Change `DB_HOST` to `localhost`:
     ```
     DB_HOST=localhost
     ```

4. **Run Docker for the Database**:
   ```sh
   docker-compose up -d golizilla-db
   ```
   This will start only the PostgreSQL container.

5. **Install Go dependencies**:
   ```sh
   go mod tidy
   ```

6. **Run the API**:
   ```sh
   go run cmd/http/main.go
   ```

7. **Access the API**:
   The API should be running on `http://localhost:8080`.

8. **Access Swagger Documentation**:
   Swagger documentation is available at `http://localhost:8080/swagger/index.html`.
   
   **Note**: Swagger documentation is accessible only when `ENV` is not set to production in the `.env` file.

## Swagger API Documentation

Swagger is used to document and interact with the API.

- **URL**: Once the server is running, Swagger documentation is available at:
  - `http://localhost:8080/swagger/index.html`
- You can use Swagger to interact with the API endpoints, view request and response formats, and understand the behavior of different endpoints.

**Note**: Swagger documentation is accessible only when `ENV` is not set to production in the `.env` file.

