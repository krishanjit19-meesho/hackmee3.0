# Meesho Clone Backend

A Golang backend server that integrates with Meesho's staging API to provide homescreen data and user management.

## 🚀 Features

- **User Management**: Phone number-based login with MySQL storage
- **Meesho API Integration**: Real-time data from Meesho's staging API
- **RESTful APIs**: Clean API design with proper error handling
- **CORS Support**: Frontend integration ready
- **Database Management**: MySQL with GORM ORM

## 📁 Project Structure

```
backend/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── handlers/               # HTTP request handlers
│   │   ├── auth_handler.go     # Authentication endpoints
│   │   └── homescreen_handler.go # Homescreen data endpoints
│   ├── services/               # Business logic
│   │   ├── user_service.go     # User management
│   │   └── meesho_service.go   # Meesho API integration
│   ├── models/                 # Data models
│   │   └── user.go             # User and API models
│   └── middleware/             # HTTP middleware
│       ├── cors.go             # CORS handling
│       ├── logger.go           # Request logging
│       └── error_handler.go    # Error handling
├── configs/
│   └── database.go             # Database configuration
├── setup_database.sql          # SQL setup script
├── go.mod                      # Go module file
└── README.md                   # This file
```

## 🛠️ Setup Instructions

### 1. Prerequisites

- Go 1.21+ installed
- MySQL 8.0+ installed and running
- Git

### 2. Database Setup

First, set up the MySQL database:

```bash
# Connect to MySQL (enter your password when prompted)
mysql -u root -p

# Run the setup commands
CREATE DATABASE IF NOT EXISTS meesho_clone;
USE meesho_clone;

CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL UNIQUE,
    phone_number VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    
    INDEX idx_user_id (user_id),
    INDEX idx_phone_number (phone_number),
    INDEX idx_deleted_at (deleted_at)
);

# Verify setup
SHOW TABLES;
exit;
```

### 3. Environment Configuration

Create a `.env` file in the backend directory (optional, defaults work for local development):

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_mysql_password
DB_NAME=meesho_clone

# Server Configuration
SERVER_PORT=8080
GIN_MODE=debug
```

### 4. Install Dependencies

```bash
cd backend
go mod tidy
```

### 5. Run the Server

```bash
go run cmd/main.go
```

The server will start on `http://localhost:8080`

## 📋 API Endpoints

### Authentication

#### POST `/api/v1/auth/login`
Login with phone number (creates user if doesn't exist)

**Request:**
```json
{
    "phone_number": "9876543210"
}
```

**Response:**
```json
{
    "success": true,
    "data": {
        "user_id": "user_abc123456789",
        "phone_number": "9876543210",
        "name": "User 3210",
        "message": "Login successful"
    }
}
```

#### GET `/api/v1/auth/profile/:user_id`
Get user profile information

**Response:**
```json
{
    "success": true,
    "data": {
        "id": 1,
        "user_id": "user_abc123456789",
        "phone_number": "9876543210",
        "name": "User 3210",
        "created_at": "2024-01-01T10:00:00Z",
        "updated_at": "2024-01-01T10:00:00Z"
    }
}
```

### Homescreen Data

#### GET `/api/v1/home/?user_id=USER_ID`
Get complete homescreen data from Meesho API

**Response:**
```json
{
    "success": true,
    "user_id": "user_abc123456789",
    "top_nav_bar": {...},
    "widget_groups": [...],
    "categories": [...],
    "products": [...],
    "user_info": {
        "user_id": "user_abc123456789",
        "name": "User 3210",
        "phone": "9876543210"
    },
    "timestamp": 1704110400
}
```

#### GET `/api/v1/home/categories?user_id=USER_ID`
Get only categories data

#### GET `/api/v1/home/products?user_id=USER_ID`
Get only products data

### Product Search

#### GET `/api/v1/products/search?user_id=USER_ID&q=QUERY`
Search products (placeholder implementation)

#### GET `/api/v1/products/:product_id?user_id=USER_ID`
Get product details (placeholder implementation)

### Health Check

#### GET `/health`
Server health check

## 🧪 Testing the APIs

### 1. Test Health Check
```bash
curl http://localhost:8080/health
```

### 2. Test User Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone_number": "9876543210"}'
```

### 3. Test Homescreen Data
```bash
# Replace USER_ID with the user_id from login response
curl "http://localhost:8080/api/v1/home/?user_id=USER_ID"
```

### 4. Test Categories
```bash
curl "http://localhost:8080/api/v1/home/categories?user_id=USER_ID"
```

## 🔄 Integration Flow

1. **Frontend Login** → `POST /api/v1/auth/login` with phone number
2. **Get User ID** → Store returned `user_id` in frontend localStorage
3. **Fetch Homescreen** → `GET /api/v1/home/?user_id=USER_ID`
4. **Display Data** → Use returned categories and products to power frontend

## 🐛 Troubleshooting

### Database Connection Issues
- Ensure MySQL is running: `brew services start mysql`
- Check credentials in `.env` file
- Verify database exists: `mysql -u root -p -e "SHOW DATABASES;"`

### Port Already in Use
- Change port in `cmd/main.go` or kill existing process:
```bash
lsof -ti:8080 | xargs kill -9
```

### CORS Issues
- Frontend CORS is handled automatically
- For testing, use browser extension or Postman

## 📈 Performance Features

- **Connection Pooling**: GORM handles database connections efficiently
- **Error Handling**: Comprehensive error responses
- **Logging**: Request/response logging for debugging
- **Recovery**: Panic recovery middleware

## 🔒 Security Features

- **Input Validation**: Phone number format validation
- **SQL Injection Protection**: GORM parameterized queries
- **Error Sanitization**: No sensitive data in error responses

## 🚀 Production Deployment

For production deployment:

1. Set `GIN_MODE=release` in environment
2. Use proper MySQL credentials
3. Implement rate limiting
4. Add authentication tokens
5. Set up monitoring and logging

## 📊 Database Schema

### Users Table
```sql
id (INT, PRIMARY KEY, AUTO_INCREMENT)
user_id (VARCHAR(255), UNIQUE) 
phone_number (VARCHAR(20), UNIQUE)
name (VARCHAR(255))
created_at (TIMESTAMP)
updated_at (TIMESTAMP)
deleted_at (TIMESTAMP, nullable)
```

## 🤝 Contributing

1. Fork the repository
2. Create feature branch
3. Make changes
4. Test thoroughly
5. Submit pull request

---

**Built with ❤️ using Go, Gin, GORM, and MySQL** 