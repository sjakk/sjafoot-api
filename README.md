# ⚽ SJAFoot API

**SJAFoot** is a backend **RESTful API** built with Go for managing football (soccer) championship data, user authentication, and fan notifications.  
It integrates with the [Football Data API](https://www.football-data.org/) to provide real-time match and competition data.  

---

## ✨ Features

- 🔐 **JWT Authentication** — Secure endpoints for user registration and login  
- 🌍 **External API Integration** — Fetches championships and matches from `api.football-data.org`  
- 🙌 **Fan Registration** — Users can register their favorite team to receive notifications  
- 📢 **Broadcast System** — Simulate notifications to fans of a specific team (protected endpoint)  
- 📂 **Database Migrations** — Managed with `golang-migrate`  
- 🐳 **Dockerized Environment** — Easy setup & consistent deployment with Docker Compose  

---

## 🛠 Technology Stack

- **Backend:** Go  
- **Database:** PostgreSQL  
- **Routing:** [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter)  
- **Migrations:** [golang-migrate/migrate](https://github.com/golang-migrate/migrate)  
- **Authentication:** [golang-jwt/jwt](https://github.com/golang-jwt/jwt)  
- **Containerization:** Docker & Docker Compose  

---

## ⚡ Prerequisites

- 🐳 Docker  
- 📦 Docker Compose  

For **local development without Docker**, you will also need:  

- 🐹 Go `1.22+`  
- 🐘 PostgreSQL  
- 🔄 migrate-cli  

---

## 🚀 Running the Application

You can run the API in two ways: **with Docker (recommended)** or **locally**.

### ▶️ Running with Docker (Recommended)

1. **Create an Environment File**  
   Create a `.env` file in the project root and copy the contents from `.env.example`:  

   ```env
   # PostgreSQL Credentials
   DB_USER=sjafoot_user
   DB_PASSWORD=yourpassword
   DB_NAME=sjafoot

   # JWT Secret
   JWT_SECRET=a-very-strong-and-secret-key-that-is-long-and-secure
   ```

2. **Build & Run**  

   ```bash
   docker-compose up --build
   ```

   The API will be available at 👉 **http://localhost:4000**

---

### ▶️ Running Locally (Without Docker)

1. **Start PostgreSQL**  

   ```sql
   CREATE ROLE sjafoot_user WITH LOGIN PASSWORD 'yourpassword';
   CREATE DATABASE sjafoot WITH OWNER = sjafoot_user;
   \c sjafoot
   CREATE EXTENSION IF NOT EXISTS citext;
   ```

2. **Run Database Migrations**  

   ```bash
   migrate -database "postgres://sjafoot_user:yourpassword@localhost/sjafoot?sslmode=disable" -path migrations up
   ```

3. **Run the Application**  

   ```bash
   go run ./cmd/api      -port=4000      -db-dsn="postgres://sjafoot_user:yourpassword@localhost/sjafoot?sslmode=disable"      -jwt-secret="a-very-strong-and-secret-key-that-is-long-and-secure"
   ```

---

## 📡 API Endpoints

Base URL: **`http://localhost:4000`**

### 🔑 Authentication & Users

#### 1. Register a User
```http
POST /users
```
Creates a new user.  
✅ Public

**Example Request**
```bash
curl -X POST http://localhost:4000/users   -H "Content-Type: application/json"   -d '{"name":"Admin User","email":"admin@example.com","password":"password123"}'
```

---

#### 2. User Login
```http
POST /auth/login
```
Authenticates a user and returns a JWT.  
✅ Public

**Example Request**
```bash
curl -X POST http://localhost:4000/auth/login   -H "Content-Type: application/json"   -d '{"email":"admin@example.com","password":"password123"}'
```

---

### 👥 Fans (Torcedores)

#### 3. Register a Fan
```http
POST /torcedores
```
Registers a fan to receive notifications.  
✅ Public

**Example Request**
```bash
curl -X POST http://localhost:4000/torcedores   -H "Content-Type: application/json"   -d '{"nome":"João Silva","email":"joao.silva@example.com","time":"Flamengo"}'
```

---

### 🏆 Championships & Matches

#### 4. List Championships
```http
GET /v1/campeonatos
```
Returns available championships.  
🔒 Protected (JWT Required)

```bash
curl -H "Authorization: Bearer $TOKEN" http://localhost:4000/v1/campeonatos
```

---

#### 5. List Matches
```http
GET /v1/campeonatos/{id}/partidas
```
Returns matches for a given championship.  
✅ Public

```bash
curl http://localhost:4000/v1/campeonatos/2013/partidas
```

---

### 📢 Broadcasts

#### 6. Send Broadcast
```http
POST /broadcast
```
Sends a notification to fans of a team.  
🔒 Protected (JWT Required)

```bash
curl -X POST http://localhost:4000/broadcast   -H "Authorization: Bearer $TOKEN"   -H "Content-Type: application/json"   -d '{"tipo":"inicio","time":"Flamengo","mensagem":"O jogo vai começar!"}'
```

---

### 🩺 System

#### 7. Health Check
```http
GET /v1/healthcheck
```
Returns API status.  
✅ Public

```bash
curl http://localhost:4000/v1/healthcheck
```

---

## 📌 Project Status

✔️ MVP completed  
🚧 Future improvements: WebSocket support, email notifications, admin panel  
