# OrdinaryBank (Transactly)

OrdinaryBank is a comprehensive banking application that allows users to create accounts, manage transfers, and track their financial transactions. This project demonstrates a robust full-stack implementation using Go for the backend and React for the frontend.

## 🚀 Technology Stack

### Backend
- **Language**: Go (Golang) 1.24
- **Framework**: [Gin](https://gin-gonic.com/) (HTTP web framework)
- **Database**: PostgreSQL 18
- **ORM/Query Builder**: [SQLC](https://sqlc.dev/) (Type-safe SQL compiler)
- **Authentication**: PASETO (Platform-Agnostic Security Tokens)
- **Migrations**: Golang Migrate

### Frontend
- **Framework**: React 18
- **Build Tool**: Vite
- **Styling**: CSS (Modular)
- **HTTP Client**: Axios

### DevOps & Tools
- **Containerization**: Docker
- **Orchestration**: Docker Compose
- **Testing**: Testify, MockGen

## 🛠️ Prerequisites

Ensure you have the following installed:
- [Docker](https://docs.docker.com/get-docker/) & Docker Compose
- [Go](https://go.dev/dl/) 1.24+ (for local development)
- [Node.js](https://nodejs.org/) & NPM (for frontend development)

## 🏁 Getting Started

### Method 1: Docker Compose (Recommended)

The easiest way to run the entire application (Database, Backend API, Frontend) is using Docker Compose.

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/nilesh0729/Transactly.git
    cd Transactly
    ```

2.  **Set up Environment Variables:**
    Copy the example environment file:
    ```bash
    cp .env.example .env
    ```
    *Note: The defaults in `docker-compose.yml` work out-of-the-box, but you can customize `.env` if needed.*

3.  **Start the Application:**
    ```bash
    docker-compose up --build
    ```

4.  **Access the App:**
    - Frontend: `http://localhost:80`
    - Backend API: `http://localhost:8080`

### Method 2: Manual Setup (Local Development)

If you prefer to run services individually for development.

#### 1. Database Setup
Start a Postgres container and run migrations.

```bash
# Start Postgres container (uses default 'root' user and 'secret' password from simple Makefile)
make Container

# Create Database
make Createdb

# Run Migrations
make MigrateUp
```

#### 2. Backend Setup
Run the Go API server.

```bash
# Install dependencies
go mod tidy

# Run the server
make Server
# OR
go run cmd/api/main.go
```
The server will start on `http://localhost:8080`.

#### 3. Frontend Setup
Run the React application.

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```
The frontend will start on the URL provided by Vite (usually `http://localhost:5173`).

## ⚙️ Configuration

The application uses environment variables for configuration. See `.env.example` for reference.

| Variable | Description | Default |
| :--- | :--- | :--- |
| `DB_DRIVER` | Database driver | `postgres` |
| `DB_SOURCE` | Postgres connection string | `postgresql://...` |
| `SERVER_ADDRESS` | Address for the API to listen on | `0.0.0.0:8080` |
| `TOKEN_SYMMETRIC_KEY`| Secret key for token signing | *Change this!* |
| `ACCESS_TOKEN_DURATION`| Token validity duration | `15m` |

## 🧪 Development Commands

The `Makefile` provides several helpful commands for development:

- `make Test`: Run Go tests.
- `make Sqlc`: Generate Go code from SQL queries using SQLC.
- `make Mock`: Generate mocks for testing.
- `make MigrateUp` / `make MigrateDown`: Manage database migrations.
