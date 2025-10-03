# SJAFoot API

SJAFoot is a backend RESTful API built with Go for managing football (soccer) championship data, user authentication, and fan notifications. It interfaces with an external API to provide real-time match and competition data.

## Features

-   **JWT Authentication**: Secure endpoints for user registration and login.
-   **External API Integration**: Fetches championship and match data from `api-football-data.org`.
-   **Fan Registration**: Allows users to register their favorite team to receive notifications.
-   **Broadcast System**: A protected endpoint to simulate sending notifications to registered fans of a specific team.
-   **Database Migrations**: Uses `golang-migrate` for version-controlled schema management.
-   **Dockerized Environment**: Fully containerized with Docker Compose for easy setup and consistent deployment.

## Technology Stack

-   **Backend**: Go
-   **Database**: PostgreSQL
-   **Routing**: `julienschmidt/httprouter`
-   **Migrations**: `golang-migrate/migrate`
-   **Authentication**: JWT (`golang-jwt/jwt`)
-   **Containerization**: Docker & Docker Compose

## Prerequisites

-   Docker
-   Docker Compose

For local development without Docker, you will also need:
-   Go (version 1.22+)
-   PostgreSQL
-   [migrate-cli](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

---

## Running the Application

There are two ways to run the application: with Docker (recommended) or locally.

### Running with Docker (Recommended)

This is the simplest and most reliable method. It sets up the Go application and the PostgreSQL database automatically.

1.  **Create an Environment File**

    Create a file named `.env` in the root of the project and copy the contents from the `.env.example` below.

    `.env.example`
    ```env
    # PostgreSQL Credentials
    DB_USER=sjafoot_user
    DB_PASSWORD=yourpassword
    DB_NAME=sjafoot

    # JWT Secret
    JWT_SECRET=a-very-strong-and-secret-key-that-is-long-and-secure
    ```

2.  **Build and Run**

    Open your terminal in the project root and run:
    ```sh
    docker-compose up --build
    ```
    The API will be available at `http://localhost:4000`.

### Running Locally (Without Docker)

1.  **Start PostgreSQL**

    Make sure you have a PostgreSQL server running and create the user and database.
    ```sql
    -- In psql
    CREATE ROLE sjafoot_user WITH LOGIN PASSWORD 'yourpassword';
    CREATE DATABASE sjafoot WITH OWNER = sjafoot_user;
    \c sjafoot
    CREATE EXTENSION IF NOT EXISTS citext;
    ```

2.  **Run Database Migrations**

    Use the `migrate-cli` to apply the database migrations.
    ```sh
    migrate -database "postgres://sjafoot_user:yourpassword@localhost/sjafoot?sslmode=disable" -path migrations up
    ```

3.  **Run the Application**

    Execute the `main.go` file with the required configuration flags.
    ```sh
    go run ./cmd/api \
        -port=4000 \
        -db-dsn="postgres://sjafoot_user:yourpassword@localhost/sjafoot?sslmode=disable" \
        -jwt-secret="a-very-strong-and-secret-key-that-is-long-and-secure"
    ```

---

## API Endpoints

**Base URL**: `http://localhost:4000`

### Authentication & Users

#### 1. Register a New User
-   **Endpoint**: `POST /users`
-   **Description**: Creates a new user in the database for accessing protected endpoints. The first user to register becomes an 'admin'.
-   **Protection**: Public
-   **Example Request**:
    ```sh
    curl -i -X POST -H "Content-Type: application/json" \
    -d '{"name": "Admin User", "email": "admin@example.com", "password": "password123"}' \
    http://localhost:4000/users
    ```
-   **Example Response (`201 Created`)**:
    ```json
    {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
    ```

#### 2. User Login
-   **Endpoint**: `POST /auth/login`
-   **Description**: Authenticates a user and returns a JWT.
-   **Protection**: Public
-   **Example Request**:
    ```sh
    curl -i -X POST -H "Content-Type: application/json" \
    -d '{"email": "admin@example.com", "password": "password123"}' \
    http://localhost:4000/auth/login
    ```
-   **Example Response (`200 OK`)**:
    ```json
    {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
    ```
    *(For the protected endpoints below, we will store this token in a shell variable)*
    ```sh
    TOKEN=$(curl -s -X POST -H "Content-Type: application/json" -d '{"email": "admin@example.com", "password": "password123"}' http://localhost:4000/auth/login | jq -r .token)
    ```

### Torcedores (Fans)

#### 3. Register a Fan
-   **Endpoint**: `POST /torcedores`
-   **Description**: Registers a fan to receive notifications for their favorite team.
-   **Protection**: Public
-   **Example Request**:
    ```sh
    curl -i -X POST -H "Content-Type: application/json" \
    -d '{"nome": "João Silva", "email": "joao.silva@example.com", "time": "Flamengo"}' \
    http://localhost:4000/torcedores
    ```
-   **Example Response (`201 Created`)**:
    ```json
    {
        "id": 1,
        "nome": "João Silva",
        "email": "joao.silva@example.com",
        "time": "Flamengo",
        "mensagem": "Cadastro realizado com sucesso"
    }
    ```

### Championships & Matches

#### 4. List Championships
-   **Endpoint**: `GET /v1/campeonatos`
-   **Description**: Returns a list of available football championships.
-   **Protection**: Protected (JWT Required)
-   **Example Request**:
    ```sh
    curl -i -H "Authorization: Bearer $TOKEN" http://localhost:4000/v1/campeonatos
    ```
-   **Example Response (`200 OK`)**:
    ```json
    [
        {
            "id": 2013,
            "nome": "Campeonato Brasileiro Série A",
            "temporada": "2025"
        },
        {
            "id": 2001,
            "nome": "UEFA Champions League",
            "temporada": "2025"
        }
    ]
    ```

#### 5. List Matches for a Championship
-   **Endpoint**: `GET /v1/campeonatos/{id}/partidas`
-   **Description**: Returns a list of matches for a given championship ID.
-   **Protection**: Public (as currently implemented)
-   **Example Request**:
    ```sh
    # Get matches for Campeonato Brasileiro (ID 2013)
    curl -i http://localhost:4000/v1/campeonatos/2013/partidas
    ```
-   **Example Response (`200 OK`)**:
    ```json
    {
        "rodada": 0,
        "partidas": [
            {
                "time_casa": "Flamengo",
                "time_fora": "Palmeiras",
                "placar": "2-1"
            },
            ...
        ]
    }
    ```

### Broadcast

#### 6. Send Broadcast
-   **Endpoint**: `POST /broadcast`
-   **Description**: Sends a notification to all fans registered for a specific team. Requires an 'admin' role. The notification is simulated by logging to the server console.
-   **Protection**: Protected (Admin JWT Required)
-   **Example Request**:
    ```sh
    curl -i -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" \
    -d '{"tipo": "inicio", "time": "Flamengo", "mensagem": "O jogo vai começar!"}' \
    http://localhost:4000/broadcast
    ```
-   **Example Response (`200 OK`)**:
    ```json
    {
        "status": "broadcast initiated",
        "team": "Flamengo",
        "event_type": "inicio",
        "notified_fans": 2
    }
    ```

### System

#### 7. Health Check
-   **Endpoint**: `GET /v1/healthcheck`
-   **Description**: Provides the status of the API.
-   **Protection**: Public
-   **Example Request**:
    ```sh
    curl -i http://localhost:4000/v1/healthcheck
    ```
-   **Example Response (`200 OK`)**:
    ```json
    {
        "status": "available",
        "environment": "development",
        "version": "1.0.0"
    }
    ```
