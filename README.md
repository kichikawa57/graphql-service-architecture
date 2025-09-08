# GraphQL Service Architecture

A high-performance GraphQL API service built with **Go**, **gqlgen**, **MySQL**, and **Docker** with DataLoader optimization for efficient data fetching.

## 🏗️ Architecture

This service provides a GraphQL API with optimized data loading patterns using DataLoader to prevent N+1 queries.

**Tech Stack:**
- **Backend:** Go with gqlgen (GraphQL server)
- **Database:** MySQL 8.0 
- **Optimization:** DataLoader for batched queries
- **Infrastructure:** Docker & Docker Compose
- **Migration:** Custom migration system

```
graphql-service-architecture/
├── cmd/               # Command line tools
│   └── migrate/       # Database migration tool
├── database/          # Database connection & config
├── graph/             # Generated GraphQL resolvers & schema
├── migrations/        # Database migration files
├── models/            # Data models & repository layer
│   └── loaders/       # DataLoader implementations
├── schema/            # GraphQL schema definitions
├── main.go            # Application entrypoint
├── docker-compose.yml # Multi-service orchestration
└── Makefile          # Development commands
```

## ✨ Features

### GraphQL API Endpoints
| Operation | Type | Description |
|-----------|------|-------------|
| `users` | Query | Fetch all users with their posts |
| `user(id: ID!)` | Query | Fetch user by ID with posts |
| `posts` | Query | Fetch all posts |
| `post(id: ID!)` | Query | Fetch post by ID |
| `createUser(input: CreateUserInput!)` | Mutation | Create new user |
| `updateUser(id: ID!, input: UpdateUserInput!)` | Mutation | Update existing user |
| `deleteUser(id: ID!)` | Mutation | Delete user |

### Key Features
- **DataLoader Integration**: Prevents N+1 queries when fetching related data
- **User-Post Relations**: Users can have multiple posts with optimized loading
- **Database Migrations**: Version-controlled schema management
- **CORS Support**: Cross-origin requests enabled
- **Docker Development**: Containerized development environment

## 🚀 Quick Start

### Prerequisites
- Docker 20.0+
- Docker Compose 2.0+
- Go 1.23+ (for local development)

### Launch Services

```bash
# Clone and navigate to project
cd graphql-service-architecture

# Setup and start all services
make setup

# Or manually:
docker-compose up --build -d
make migrate-up
```

### 🌐 Access Points

| Service | URL | Description |
|---------|-----|-------------|
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

### Posts Table
```sql
CREATE TABLE posts (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    content     TEXT,
    user_id     INT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

## 📋 GraphQL Examples

### Fetch Users with Posts (DataLoader Optimized)
```graphql
query GetUsersWithPosts {
  users {
    id
    name
    email
    posts {
      id
      title
      content
      createdAt
    }
    createdAt
    updatedAt
  }
}
```

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
    posts {
      id
      title
    }
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
    posts {
      id
      title
    }
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

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `mysql` | MySQL host (container name in Docker) |
| `DB_PORT` | `3306` | MySQL port |
| `DB_USER` | `root` | MySQL username |
| `DB_PASSWORD` | `password` | MySQL password |
| `DB_NAME` | `graphql_db` | Database name |
| `PORT` | `8080` | GraphQL server port |

## 🛠️ Development

### Local Development
```bash
# Run GraphQL server locally
go run main.go

# Build the application
go build -o main .

# Run tests
go test ./...
```

### GraphQL Schema Generation
```bash
# Regenerate GraphQL resolvers and types
make gen
# or directly:
gqlgen generate
```

### Database Management
```bash
# Create new migration
make migrate-create name=add_posts_table

# Run migrations
make migrate-up

# Rollback migrations
make migrate-down

# Check migration status
make migrate-status
```

### Docker Development
```bash
# Start all services
make dev

# Setup everything (build + migrate)
make setup

# Stop services
make down

# View logs
make logs

# Access MySQL directly
make db

# Clean up everything
make clean
```

## 🔧 Troubleshooting

### Common Issues

**MySQL Connection Error**
```bash
# Check if services are running
make ps

# Check MySQL logs
make log-mysql

# Verify database connectivity
make db
```

**GraphQL API Connection Error**
```bash
# Check backend service status
make log-backend

# Test GraphQL endpoint
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query":"query{users{id name}}"}'
```

**DataLoader Issues**
```bash
# Check for N+1 query problems in logs
make log-query

# Verify DataLoader middleware is working
make log-backend | grep -i "loader"
```

### Useful Commands
```bash
# Restart all services
make reset

# Rebuild containers
docker-compose up --build --force-recreate

# View specific service logs
make log-mysql    # MySQL logs
make log-backend  # Backend logs

# Clean up everything
make clean
```

## 📄 License

MIT License