# User Service

The User Service is a core microservice within the GoalCircle platform responsible for user lifecycle management, profile management, onboarding, and administrative user operations.

It exposes gRPC APIs consumed by internal services such as the API Gateway and Admin Service.

---

## Overview

The service manages:

* User onboarding and profile management
* User role assignment
* Password management
* User retrieval and search
* Administrative user management
* User account status control (block/unblock)

---

## Features

### User Features

* Complete user onboarding
* Role selection and management
* Retrieve user profile
* Update user profile
* Change password
* Retrieve user details

### Administrative Features

* Retrieve users with pagination
* Search users
* View user details
* Block user accounts
* Unblock user accounts
* User management operations

---

## Technology Stack

* Go
* gRPC
* PostgreSQL
* GORM
* Docker
* Zap Logger
* Protocol Buffers

---

## Architecture

The service follows Clean Architecture principles to maintain separation of concerns, testability, and scalability.

```text
┌──────────────────┐
│   gRPC Handler   │
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│     Use Case     │
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│   Repository     │
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│   PostgreSQL     │
└──────────────────┘
```

---

## Project Structure

```text
user-service/
├── cmd/
│   └── server/
├── internal/
│   ├── domain/
│   ├── entity/
│   ├── repository/
│   ├── usecase/
│   ├── handler/
│   ├── infrastructure/
│   └── middleware/
├── proto/
├── configs/
├── migrations/
├── pkg/
├── Dockerfile
├── .air.toml
├── go.mod
└── README.md
```

---

## gRPC Services

### User Service

| Method         | Description                         |
| -------------- | ----------------------------------- |
| OnboardUser    | Complete user onboarding process    |
| GetProfile     | Retrieve authenticated user profile |
| UpdateProfile  | Update user profile information     |
| ChangePassword | Change user password                |
| GetUserByID    | Retrieve user details               |

### Admin User Management Service

| Method         | Description                    |
| -------------- | ------------------------------ |
| GetUsers       | Retrieve users with pagination |
| GetUserDetails | Retrieve specific user details |
| SearchUsers    | Search users                   |
| BlockUser      | Block a user account           |
| UnblockUser    | Unblock a user account         |

---

## Security

The service implements:

* JWT-based authentication
* Role-based authorization
* Protected administrative operations
* Password hashing and verification
* Request validation

---

## Logging

Structured logging is implemented using Zap Logger.

Logged information includes:

* Incoming requests
* Business errors
* Database failures
* gRPC errors
* Panic recovery events

---

## Error Handling

The service follows a centralized error handling strategy:

* Domain-specific errors
* gRPC status codes
* Repository-level error wrapping
* Recovery interceptor for panic handling

---

## Development

### Install Dependencies

```bash
go mod tidy
```

### Run Service

```bash
go run cmd/server/main.go
```

### Run with Air

```bash
air
```

---

## Docker

### Build Image

```bash
docker build -t user-service .
```

### Run Container

```bash
docker run -p 50052:50052 user-service
```

---

## Environment Variables

```env
PORT=50052

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=goalcircle

JWT_SECRET=your-secret-key
```

---

## Health & Monitoring

The service supports:

* Structured application logs
* gRPC health checks
* Containerized deployment
* Production-ready monitoring integration

---

## License

Internal Project – GoalCircle Platform.
