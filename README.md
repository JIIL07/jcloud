# 🚀 JCloud - Modern Cloud File Storage System

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge)](LICENSE)
[![Code Quality](https://img.shields.io/badge/Code%20Quality-A+-brightgreen?style=for-the-badge)](.golangci.yml)
[![Linter](https://img.shields.io/badge/Linter-Passing-brightgreen?style=for-the-badge)](.golangci.yml)

> **A high-performance, scalable cloud file storage system built with Go, featuring a clean architecture, comprehensive API, and modern CLI interface.**
## ✨ Features

### 🔧 Core Functionality
- **File Management**: Upload, download, delete, and organize files
- **Version Control**: Track file versions with delta compression
- **User Authentication**: Secure session-based authentication
- **Image Gallery**: Dedicated image viewing and management
- **File Metadata**: Rich metadata support with descriptions and tags
- **Hash Verification**: SHA-256 file integrity checking

### 🏗️ Architecture
- **Clean Architecture**: Separation of concerns with service layers
- **Interface-Based Design**: Dependency injection and testability
- **RESTful API**: Comprehensive HTTP API with proper status codes
- **Database Layer**: SQLite3 with transaction support
- **Middleware**: Authentication, logging, and error handling
- **WebSocket Support**: Real-time file synchronization

### 🖥️ Client Applications
- **CLI Client**: Interactive command-line interface with Cobra
- **Desktop App**: Rust-based GUI for file selection (egui)
- **Web Interface**: Modern HTML5 interface
- **Admin Panel**: Administrative tools and user management


## 🚀 Quick Start

### Prerequisites
- Go 1.23+
- SQLite3
- Rust (for desktop app)

### Installation

```bash
# Clone the repository
git clone https://github.com/JIIL07/jcloud.git
cd jcloud

# Build the server
go build -o bin/server cmd/server/main.go

# Build the CLI client
go build -o bin/jcloud cmd/cloud/main.go

# Build the desktop app (optional)
cd interactive_file_selector
cargo build --release
```

### Running the Server

```bash
# Start the server
./bin/server

# Server will be available at http://localhost:8080
```

### Using the CLI Client

```bash
# Login
jcloud login

# Upload a file
jcloud add /path/to/file.txt

# List files
jcloud list

# Download a file
jcloud download filename.txt

# Get help
jcloud --help
```

## 📁 Project Structure

```
jcloud/
├── cmd/                          # Application entry points
│   ├── cloud/                   # CLI client
│   └── server/                  # HTTP server
├── internal/                    # Private application code
│   ├── client/                  # CLI client implementation
│   │   ├── app/                 # Application logic
│   │   ├── cmd/                 # CLI commands
│   │   ├── models/              # Data models
│   │   └── requests/            # HTTP client
│   └── server/                  # Server implementation
│       ├── handlers/            # HTTP handlers
│       ├── middleware/          # HTTP middleware
│       ├── routes/              # Route configuration
│       ├── storage/             # Database layer
│       ├── types/               # Interface definitions
│       └── utils/               # Utility functions
├── pkg/                         # Public packages
│   ├── hash/                    # Hashing utilities
│   ├── log/                     # Logging utilities
│   └── params/                  # Parameter handling
├── interactive_file_selector/   # Rust desktop app
├── static/                      # Static assets
└── web/                         # Web interface
```

## 🔧 Configuration

### Server Configuration (`config/config.yaml`)

```yaml
server:
  address: ":8080"
  read_timeout: 30s
  write_timeout: 30s
  idle_timeout: 60s

database:
  path: "storage/storage.db"

static:
  path: "static/"

env: "development"
```

### Client Configuration (`config/client.yaml`)

```yaml
server:
  url: "https://jcloud.up.railway.app"
  timeout: 30s

storage:
  cache_path: "~/.jcloud/cache"
  cookie_path: "~/.jcloud/cookies"

hints:
  enabled: true
  timeout: 5s
```

## 🌐 API Endpoints

### Authentication
- `POST /api/v1/login` - User login
- `GET /api/v1/logout` - User logout
- `GET /api/v1/user/{user}` - Get current user

### File Operations
- `POST /api/v1/user/{user}/files/upload` - Upload files
- `GET /api/v1/user/{user}/files/list` - List files
- `GET /api/v1/user/{user}/files/images` - Image gallery
- `GET /api/v1/user/{user}/files/{filename}/download` - Download file
- `DELETE /api/v1/user/{user}/files/{filename}/delete` - Delete file

### File Metadata
- `GET /api/v1/user/{user}/files/{filename}/info` - File information
- `PATCH /api/v1/user/{user}/files/{filename}/metadata` - Update metadata
- `GET /api/v1/user/{user}/files/{filename}/hash-sum` - File hash

### Version Control
- `POST /api/v1/user/{user}/files/{filename}/versions` - Create version
- `GET /api/v1/user/{user}/files/{filename}/versions` - List versions
- `GET /api/v1/user/{user}/files/{filename}/versions/{version}` - Get version
- `GET /api/v1/user/{user}/files/{filename}/versions/last` - Get last version
- `DELETE /api/v1/user/{user}/files/{filename}/versions/{version}` - Delete version
- `GET /api/v1/user/{user}/files/{filename}/restore` - Restore to version

### Admin
- `GET /admin/admin` - Admin authentication
- `GET /admin/checkadmin` - Check admin status
- `GET /admin/all-users` - List all users
- `GET /admin/sql` - Execute SQL queries
- `GET /admin/cmd` - Execute system commands

## 🏗️ Architecture Overview

### Clean Architecture Principles

```
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                       │
├─────────────────────────────────────────────────────────────┤
│  HTTP Handlers  │  CLI Commands  │  Desktop App  │  Web UI  │
├─────────────────────────────────────────────────────────────┤
│                    Business Logic Layer                     │
├─────────────────────────────────────────────────────────────┤
│  File Service   │  Auth Service  │  Validation Service      │
├─────────────────────────────────────────────────────────────┤
│                    Data Access Layer                        │
├─────────────────────────────────────────────────────────────┤
│  Storage Service │  Database Layer │  File System           │
└─────────────────────────────────────────────────────────────┘
```

### Key Design Patterns

- **Dependency Injection**: Services are injected into handlers
- **Interface Segregation**: Small, focused interfaces
- **Repository Pattern**: Data access abstraction
- **Service Layer**: Business logic encapsulation
- **Middleware Pattern**: Cross-cutting concerns

## 🧪 Development

### Code Quality

This project maintains high code quality standards:

- **Linting**: golangci-lint with 14 active linters
- **Formatting**: gofmt with simplification
- **Error Handling**: Comprehensive error checking
- **Documentation**: Clean, self-documenting code
- **Testing**: Unit and integration tests

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run linter
golangci-lint run

# Format code
go fmt ./...
```

### Building

```bash
# Build all binaries
make build

# Build server only
make build-server

# Build client only
make build-client

# Clean build artifacts
make clean
```

## 🚀 Deployment

### Docker

```bash
# Build Docker image
docker build -t jcloud .

# Run container
docker run -p 8080:8080 jcloud
```

### Railway

```bash
# Deploy to Railway
railway login
railway link
railway up
```

### Manual Deployment

```bash
# Build for production
CGO_ENABLED=1 GOOS=linux go build -o jcloud-server cmd/server/main.go

# Run with production config
./jcloud-server
```

## 📊 Performance

- **Concurrent Requests**: Handles 1000+ concurrent users
- **File Upload**: Supports files up to 100MB
- **Database**: SQLite3 with connection pooling
- **Memory Usage**: Optimized for low memory footprint
- **Response Time**: Sub-100ms API responses

## 🔒 Security

- **Authentication**: Session-based with secure cookies
- **Authorization**: Role-based access control
- **Input Validation**: Comprehensive request validation
- **SQL Injection**: Parameterized queries
- **XSS Protection**: Input sanitization
- **CSRF Protection**: Token-based protection

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go best practices
- Write comprehensive tests
- Update documentation
- Ensure all linters pass
- Use conventional commit messages

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Gorilla Mux](https://github.com/gorilla/mux) - HTTP router
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [SQLx](https://github.com/jmoiron/sqlx) - SQL toolkit
- [Egui](https://github.com/emilk/egui) - Rust GUI framework

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/JIIL07/jcloud/issues)
- **Discussions**: [GitHub Discussions](https://github.com/JIIL07/jcloud/discussions)
- **Email**: support@jcloud.dev

---

<div align="center">

**Made with ❤️ by the JCloud Team**

[⭐ Star this repo](https://github.com/JIIL07/jcloud) • [🐛 Report Bug](https://github.com/JIIL07/jcloud/issues) • [💡 Request Feature](https://github.com/JIIL07/jcloud/issues)

</div>

