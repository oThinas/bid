# Bid - Real-time Auction Platform

A real-time auction platform built with Go, featuring WebSocket support for live bidding, user authentication, and PostgreSQL database integration.

## Introduction

Bid is a modern auction platform that enables real-time bidding through WebSocket connections. The application provides a complete auction experience with user registration, product creation, and live bidding functionality. Built with Go and following clean architecture principles, it offers a robust and scalable solution for online auctions.

## Features

- **User Authentication**: Secure user registration and login with session management
- **Product Management**: Create and manage auction products with base prices and end times
- **Real-time Bidding**: Live WebSocket-based bidding system
- **Auction Rooms**: Dynamic auction rooms that automatically close when time expires
- **Session Management**: Persistent sessions using PostgreSQL
- **RESTful API**: Clean REST API design with proper HTTP status codes
- **Database Migrations**: Automated database schema management
- **Docker Support**: Easy deployment with Docker Compose

## Tech Stack

- **Backend**: Go 1.24.5
- **Containerization**: Docker & Docker Compose
- **Web Framework**: Chi (HTTP router)
- **Database**: PostgreSQL
- **ORM**: SQLC (type-safe SQL)
- **Session Management**: SCS (Session Cookie Store)
- **WebSocket**: Gorilla WebSocket
- **Password Hashing**: bcrypt
- **Environment**: godotenv

## API Reference

### Authentication Endpoints

#### POST `/api/v1/users/signup`

Register a new user account.

**Request Body:**

```json
{
  "username": "string",
  "email": "string",
  "password": "string",
  "bio": "string"
}
```

**Response:**

```json
{
  "data": "uuid"
}
```

#### POST `/api/v1/users/login`

Authenticate and login user.

**Request Body:**

```json
{
  "email": "string",
  "password": "string"
}
```

**Response:**

```json
{
  "data": "logged in successfully"
}
```

#### POST `/api/v1/users/logout`

Logout current user (requires authentication).

**Response:**

```json
{
  "data": "logged out successfully"
}
```

### Product Endpoints

#### POST `/api/v1/products`

Create a new auction product (requires authentication).

**Request Body:**

```json
{
  "name": "string",
  "description": "string",
  "base_price": "decimal",
  "auction_end": "datetime"
}
```

**Response:**

```json
{
  "data": "uuid",
  "message": "auction room created"
}
```

### WebSocket Endpoints

#### GET `/api/v1/products/subscribe/{productID}`

Subscribe to real-time auction updates via WebSocket (requires authentication).

**Parameters:**

- `productID`: UUID of the product to subscribe to

**WebSocket Messages:**

- Bid updates
- Auction status changes
- Time remaining notifications

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
# Database Configuration
DATABASE_USER=your_db_user
DATABASE_PASSWORD=your_db_password
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=bid_db
```

You can use the `.env.example` file as a template.

## Run Locally

### Prerequisites

- Go 1.24.5 or higher
- PostgreSQL
- Docker

### Using Docker Compose (Recommended)

1. **Clone the repository:**

   ```bash
   git clone https://github.com/oThinas/bid.git
   cd bid
   ```

2. **Create environment file:**

   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start the database:**

   ```bash
   docker-compose up -d bid-db
   ```

4. **Run database migrations:**

   ```bash
   go run cmd/terndotenv/main.go
   ```

5. **Start the API server:**

   ```bash
   go run cmd/api/main.go
   ```

The server will be available at `http://localhost:8080`

### Manual Setup

1. **Install dependencies:**

   ```bash
   go mod download
   ```

2. **Set up PostgreSQL database**

3. **Run migrations:**

   ```bash
   go run cmd/terndotenv/main.go
   ```

4. **Start the server:**

   ```bash
   go run cmd/api/main.go
   ```

## Project Structure

```text
bid/
├── bin/                          # Binary files
├── cmd/                          # Application entry points
│   ├── api/                      # Main API server
│   │   └── main.go
│   └── terndotenv/               # Database migration tool
│       └── main.go
├── internal/                     # Internal application code
│   ├── api/                      # HTTP handlers and routing
│   │   ├── api.go                # API structure definition
│   │   ├── auction_handlers.go   # WebSocket auction handlers
│   │   ├── auth.go               # Authentication middleware
│   │   ├── constants.go          # API constants
│   │   ├── product_handlers.go   # Product CRUD handlers
│   │   ├── routes.go             # Route definitions
│   │   └── user_handlers.go      # User authentication handlers
│   ├── services/                 # Business logic layer
│   │   ├── auctions_service.go   # Auction room management
│   │   ├── bids_service.go       # Bidding logic
│   │   ├── constants.go          # Service constants
│   │   ├── products_service.go   # Product management
│   │   └── users_service.go      # User management
│   ├── store/                    # Data access layer
│   │   └── pg/                   # PostgreSQL implementation
│   │       ├── db.go             # Database connection
│   │       ├── migrations/       # Database migrations
│   │       ├── models.go         # Data models
│   │       ├── queries/          # SQL queries
│   │       └── *.sql.go          # Generated SQLC code
│   ├── usecase/                  # Application use cases
│   │   ├── products/             # Product use cases
│   │   └── users/                # User use cases
│   ├── utils/                    # Utility functions
│   │   └── json.go               # JSON encoding/decoding
│   └── validator/                # Input validation
│       └── validator.go
├── docker-compose.yml            # Docker services configuration
├── go.mod                        # Go module definition
├── go.sum                        # Go module checksums
└── README.md                     # This file
```

## Development

### Database Migrations

The project uses SQLC for type-safe database operations. To generate code after schema changes:

```bash
sqlc generate -f ./path-to-your-sqlc-file
```

Example:

```bash
sqlc generate -f ./internal/store/pg/sqlc.yaml
```

### Code Generation

```bash
# Generate SQLC code
sqlc generate

# Run migrations
go run cmd/terndotenv/main.go
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.
