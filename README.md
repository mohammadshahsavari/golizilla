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


## Project File Structure

This project is designed following **Hexagonal Architecture** (Ports and Adapters) to ensure scalability, maintainability, and separation of concerns.

### Root Directory
- **`.env` / `.env.example`**: Environment variable configuration files.
- **`docker-compose.yml` / `Dockerfile`**: Configuration for containerization and running the application.
- **`go.mod` / `go.sum`**: Go module and dependency management.
- **`README.md`**: Documentation for the project.

---

### Directories Overview

#### `.vscode`
Contains IDE-specific settings for development, including extensions and formatting rules.

---

#### `adapters`
Holds external-facing components, connecting the **core logic** with external systems and frameworks.

- **`http`**: Handles HTTP interactions.
  - **`handler`**: Processes incoming HTTP requests and invokes core services.
    - **`context`**: Manages request-specific context.
    - **`middleware`**: Authentication, authorization, and logging middleware.
    - **`presenter`**: Validates requests and formats responses.
  - **`route`**: Defines and initializes API routes.

- **`persistence`**: Manages database and logging integrations.
  - **`gorm`**: Contains database initialization and repository implementations.
  - **`logger`**: Structured logging system.
    - **`archive`**: Stores archived logs.
    - **`server`**: Contains logging service logic.


---

#### `cmd`
Entry point for the application.

- **`http`**: HTTP server setup and initialization.

---

#### `config`
Centralized configuration management, including environment variable parsing.

---

#### `core`
The **core business logic** and domain-driven components of the application.

- **`domain/model`**: Domain models such as `User`, `Questionnaire`, and `Answer`.
- **`port/repository`**: Repository interfaces for external system interactions (e.g., database, cache).
- **`service`**: Implements business logic and interacts with repositories.
  - **`utils`**: Contains helper functions used in services.

---

#### `internal`
Contains internal utilities, constants, and helpers.

- **`apperrors`**: Custom application errors and error handling mechanisms.
- **`email/template`**: Email templates for notifications (e.g., 2FA, verification emails).
- **`logmessages`**: Predefined log message templates.
- **`privilege`**: Privilege constants and role-based access definitions.

---

This structure ensures:
1. **Separation of Concerns**: Clear distinction between core logic, external integrations, and infrastructure.
2. **Maintainability**: Each module is isolated and independently testable.
3. **Scalability**: Easily extendable to support new features or integrations.

