# GraphQL Service

A modern full-stack GraphQL service built with **Go**, **Hono**, **MySQL**, and **Docker**.

## 🏗️ Architecture

This service provides both GraphQL and REST API endpoints with a complete CRUD interface for user management.

**Tech Stack:**
- **Backend:** Go with gqlgen (GraphQL server)
- **Frontend:** TypeScript with Hono framework
- **Database:** MySQL 8.0
- **Infrastructure:** Docker & Docker Compose

```
graphql-service/
├── backend/            # Go GraphQL API Server
│   ├── graph/         # Generated GraphQL resolvers
│   ├── models/        # Data models & repository layer
│   ├── database/      # Database connection & config
│   ├── schema/        # GraphQL schema definitions
│   ├── main.go        # Application entrypoint
│   └── Dockerfile     # Backend container config
├── frontend/          # Hono Web Application
│   ├── src/           # TypeScript source code
│   ├── public/        # Static assets
│   └── Dockerfile     # Frontend container config
├── database/          # Database setup
│   └── init.sql       # Schema & seed data
└── docker-compose.yml # Multi-service orchestration
```

## ✨ Features

### GraphQL API Endpoints
| Operation | Query/Mutation | Description |
|-----------|---------------|-------------|
| `users` | Query | Fetch all users |
| `user(id: ID!)` | Query | Fetch user by ID |
| `createUser(input: CreateUserInput!)` | Mutation | Create new user |
| `updateUser(id: ID!, input: UpdateUserInput!)` | Mutation | Update existing user |
| `deleteUser(id: ID!)` | Mutation | Delete user |

### REST API Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/users` | Get all users |
| `GET` | `/api/users/:id` | Get user by ID |
| `POST` | `/api/users` | Create new user |
| `PUT` | `/api/users/:id` | Update user |
| `DELETE` | `/api/users/:id` | Delete user |

## 🚀 Quick Start

### Prerequisites
- Docker 20.0+
- Docker Compose 2.0+

### Launch Services

```bash
# Clone and navigate to project
cd graphql-service

# Start all services with Docker Compose
docker-compose up --build
```

### 🌐 Access Points

| Service | URL | Description |
|---------|-----|-------------|
| **Frontend UI** | http://localhost:3000 | Web interface |
| **GraphQL Playground** | http://localhost:8080 | Interactive GraphQL IDE |
| **GraphQL API** | http://localhost:8080/query | GraphQL endpoint |
| **MySQL Database** | localhost:3306 | Database connection |

## 🗄️ Database Schema

### Users Table
```sql
CREATE TABLE users (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    email       VARCHAR(255) UNIQUE NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

**Sample Data:**
| ID | Name | Email |
|----|------|-------|
| 1 | John Doe | john@example.com |
| 2 | Jane Smith | jane@example.com |
| 3 | Bob Johnson | bob@example.com |

## 📋 GraphQL Examples

### Create User
```graphql
mutation CreateUser {
  createUser(input: {
    name: "Alice Johnson"
    email: "alice@example.com"
  }) {
    id
    name
    email
    createdAt
    updatedAt
  }
}
```

### Fetch All Users
```graphql
query GetUsers {
  users {
    id
    name
    email
    createdAt
    updatedAt
  }
}
```

### Update User
```graphql
mutation UpdateUser {
  updateUser(id: "1", input: {
    name: "Updated Name"
    email: "updated@example.com"
  }) {
    id
    name
    email
    updatedAt
  }
}
```

### Delete User
```graphql
mutation DeleteUser {
  deleteUser(id: "1")
}
```

## ⚙️ Environment Variables

### Backend (Go)
| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `localhost` | MySQL host |
| `DB_PORT` | `3306` | MySQL port |
| `DB_USER` | `root` | MySQL username |
| `DB_PASSWORD` | `password` | MySQL password |
| `DB_NAME` | `graphql_db` | Database name |
| `PORT` | `8080` | Server port |

### Frontend (Hono)
| Variable | Default | Description |
|----------|---------|-------------|
| `GRAPHQL_ENDPOINT` | `http://localhost:8080/query` | GraphQL API endpoint |
| `PORT` | `3000` | Server port |

## 🛠️ Development

### Backend Development
```bash
cd backend

# Run backend locally
go run main.go

# Build backend
go build -o main .

# Run tests
go test ./...
```

### Frontend Development
```bash
cd frontend
npm install
npm run dev
```

### Regenerate GraphQL Schema
```bash
cd backend
make gen
```

### Docker Development
```bash
cd backend

# Start all services
make dev

# Stop services
make down

# View logs
make logs

# Clean up
make clean
```

## 🔧 Troubleshooting

### Common Issues

**MySQL Connection Error**
```bash
# Check if services are running
docker-compose ps

# Check MySQL logs
docker-compose logs mysql
```

**GraphQL API Connection Error**
```bash
# Check backend service status
docker-compose logs backend

# Verify GraphQL endpoint
curl http://localhost:8080/query
```

**Frontend Display Issues**
```bash
# Check frontend service logs
docker-compose logs frontend

# Verify frontend is accessible
curl http://localhost:3000
```

### Useful Commands
```bash
# Restart all services
docker-compose restart

# Rebuild and restart
docker-compose up --build --force-recreate

# View logs for specific service
docker-compose logs -f [service-name]

# Clean up containers and volumes
docker-compose down -v
```

## 📄 License

MIT License