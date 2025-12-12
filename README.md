<a name="readme-top"></a>

<div align="center">
  <img src="https://img.icons8.com/fluency/96/null/bank-building.png" alt="Transactly Logo" width="80" height="80">

  <h3 align="center">Transactly</h3>

  <p align="center">
    A Secure, High-Performance Banking Service API
    <br />
    <a href="#api-reference"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="#getting-started">View Demo</a>
    ·
    <a href="https://github.com/your_username/transactly/issues">Report Bug</a>
    ·
    <a href="https://github.com/your_username/transactly/issues">Request Feature</a>
  </p>
</div>

## 🏦 About The Project

**Transactly** is a robust backend banking system designed to simulate core financial operations. It focuses on data integrity, high concurrency, and security.

The system handles the creation of user accounts, records balance changes, and performs safe money transfers between accounts using database transactions (ACID).

### Key Features

* **Multi-Method Authentication:** Supports both **JWT** and **PASETO** tokens for secure, stateless user sessions.
* **Safe Money Transfers:** Implements database transactions to ensure money is never lost during transfers, even in high-concurrency scenarios.
* **Account Management:** Create accounts, list balances, and track transaction history.
* **RESTful API:** Clean and structured API endpoints served via the Gin framework.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### 🛠️ Built With

* ![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) **Golang** - Core language.
* ![Gin](https://img.shields.io/badge/gin-%23008FC7.svg?style=for-the-badge&logo=go&logoColor=white) **Gin Gonic** - HTTP Web Framework.
* ![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white) **PostgreSQL** - Relational Database.
* ![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white) **Docker** - Containerization.
* **PASETO** - Platform-Agnostic Security Tokens.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## 🚀 Getting Started

To run Transactly locally, you will need **Docker** and **Go** installed.

### Prerequisites

* Go 1.21+
* Docker Desktop
* Make (Optional, for running makefile commands)

### Installation

1.  **Clone the repo**
    ```sh
    git clone [https://github.com/your_username/transactly.git](https://github.com/your_username/transactly.git)
    cd transactly
    ```

2.  **Setup Environment Variables**
    Create a `.env` file in the root directory:
    ```env
    DB_DRIVER=postgres
    DB_SOURCE=postgresql://root:secret@localhost:5432/transactly?sslmode=disable
    SERVER_ADDRESS=0.0.0.0:8080
    TOKEN_SYMMETRIC_KEY=12345678901234567890123456789012
    ACCESS_TOKEN_DURATION=15m
    ```

3.  **Start the Database**
    ```sh
    docker-compose up -d postgres
    ```

4.  **Run Migrations**
    ```sh
    make migrateup
    ```

5.  **Start the Server**
    ```sh
    go run main.go
    ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## 📖 API Reference

Here is a quick overview of the main endpoints.

### Users & Auth
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :---: |
| `POST` | `/users` | Create a new user | ❌ |
| `POST` | `/users/login` | Login and receive Access Token | ❌ |

### Accounts
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :---: |
| `POST` | `/accounts` | Create a new bank account | ✅ |
| `GET` | `/accounts/:id` | Get specific account details | ✅ |
| `GET` | `/accounts` | List accounts (with pagination) | ✅ |

### Transfers
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :---: |
| `POST` | `/transfers` | Transfer money between accounts | ✅ |

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## 🗄️ Database Schema

The project uses a normalized PostgreSQL schema with the following core entities:
* **Users:** Stores secure password hashes (bcrypt) and user info.
* **Accounts:** Holds balance and currency information, linked to Users.
* **Entries:** Records all account balance changes (Audit trail).
* **Transfers:** Records money movement between two accounts.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## 🤝 Contributing

Contributions are welcome!
1.  Fork the Project
2.  Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3.  Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4.  Push to the Branch (`git push origin feature/AmazingFeature`)
5.  Open a Pull Request

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>