# TodoTist - A Modern Todo List Application

A full-stack todo list application built with Go + ConnectRPC backend and Next.js frontend.

## 🌟 Features

- ✅ **Add Tasks**: Create new todo items with text validation
- 📋 **View Tasks**: Display all tasks sorted by creation time (newest first)  
- 🗑️ **Delete Tasks**: Remove completed or unwanted tasks
- 🔍 **Input Validation**: Client and server-side validation for task text
- 📱 **Responsive UI**: Clean, modern interface built with Tailwind CSS
- ⚡ **Real-time Updates**: Instant UI updates after operations
- 🛡️ **Error Handling**: Comprehensive error handling and user feedback

## 🏗️ Architecture

### Backend (Go + ConnectRPC)
- **Framework**: Go with ConnectRPC for type-safe API
- **Storage**: In-memory storage with thread-safe operations
- **Validation**: Input validation with length limits (1-500 characters)
- **Error Handling**: Proper HTTP status codes and error messages
- **CORS**: Configured for frontend communication

### Frontend (Next.js + TypeScript)
- **Framework**: Next.js 15 with TypeScript
- **Styling**: Tailwind CSS for responsive design
- **State Management**: React hooks for component state
- **API Client**: Official ConnectRPC client with full protocol support
- **Validation**: Client-side validation with user feedback
- **Package Management**: pnpm for efficient dependency management

## 🚀 Quick Start

### Prerequisites
- Go 1.24+
- Node.js 18+
- pnpm (recommended) or npm/yarn

### Development Startup

**Open two terminal windows**

#### Terminal 1 - Start Backend
```bash
cd backend
go mod tidy          # Install dependencies (first time only)
go run server.go     # Start server on port 8080
```
Expected output: `Server running on http://localhost:8080`

#### Terminal 2 - Start Frontend
```bash
cd frontend
pnpm install        # Install dependencies (first time only)
pnpm dev            # Start dev server on port 3000 with Turbopack
```
Expected output: `ready - started server on 0.0.0.0:3000`

*Alternatively, you can use npm:*
```bash
npm install         # Install dependencies
npm run dev         # Start dev server
```

#### Access the Application
Open your browser and navigate to `http://localhost:3000`

### Testing API Connection

To verify the backend API is working:
```bash
node test-api.js    # Run from root directory (ensure backend is running first)
```

### Stopping Services
- **Backend**: Press `Ctrl+C` in Terminal 1
- **Frontend**: Press `Ctrl+C` in Terminal 2

### Troubleshooting Common Issues

#### Port Already in Use
- **Backend (8080)**: Change port in `backend/server.go`
- **Frontend (3000)**: Use `npm run dev -- -p 3001` for different port

#### CORS Errors
Ensure backend CORS configuration includes `http://localhost:3000` in `backend/server.go`

#### Dependency Issues
```bash
cd backend && go mod tidy          # Refresh Go dependencies
cd frontend && pnpm install        # Refresh Node dependencies (recommended)
# or
cd frontend && npm install         # Using npm
```

## 🧪 Testing

### Backend Tests
```bash
cd backend
go test -v
```

### Regenerate Protocol Buffer Types
```bash
cd backend
go generate ./...    # Regenerate Go types from todo.proto
```

### Frontend Build
```bash
cd frontend
pnpm build          # Using pnpm (recommended)
# or
npm run build        # Using npm
```

## 🔌 ConnectRPC Integration

This project uses ConnectRPC for type-safe API communication between the Go backend and Next.js frontend.

### Key Features

- **Type Safety**: Protocol Buffers define the API contract, generating both Go and TypeScript types
- **Protocol Support**: Full ConnectRPC v1 protocol implementation
- **Transport**: HTTP/1.1 and HTTP/2 support with graceful fallback
- **Serialization**: JSON format for development, with optional binary format for production
- **Error Handling**: Structured error responses with proper HTTP status codes
- **CORS**: Pre-configured for development environment

### Protocol Buffer Definition

The API contract is defined in `backend/todo.proto`:

```protobuf
service TodoService {
  rpc AddTask(AddTaskRequest) returns (AddTaskResponse) {}
  rpc GetTasks(GetTasksRequest) returns (GetTasksResponse) {}
  rpc DeleteTask(DeleteTaskRequest) returns (DeleteTaskResponse) {}
}
```

### Client Implementation

The frontend uses the official ConnectRPC Web Client (`@connectrpc/connect-web`) with:

- JSON transport for compatibility
- Proper ConnectRPC headers (`Connect-Protocol-Version: 1`)
- Timeout configuration (10s)
- Type-safe request/response handling

### Development Workflow

1. **Modify Protocol Buffer**: Edit `backend/todo.proto`
2. **Generate Go Types**: Run `go generate ./...` in backend directory
3. **Generate TypeScript Types**: Run `frontend/generate-types.sh`
4. **Update Client**: Implement new methods in `frontend/src/lib/todo_connect.ts`
5. **Update UI**: Use new client methods in React components

## 📡 API Endpoints

### Add Task
- **Endpoint**: `POST /todo.v1.TodoService/AddTask`
- **Request**: `{"text": "Task description"}`
- **Response**: `{"task": {"id": "...", "text": "...", "createdAt": 1234567890}}`

### Get Tasks
- **Endpoint**: `GET /todo.v1.TodoService/GetTasks`
- **Response**: `{"tasks": [{"id": "...", "text": "...", "createdAt": 1234567890}]}`

### Delete Task
- **Endpoint**: `POST /todo.v1.TodoService/DeleteTask`
- **Request**: `{"id": "task-id"}`
- **Response**: `{"success": true}`

## 🔧 Configuration

### Backend Configuration
- **Port**: 8080 (configurable in `server.go`)
- **CORS Origins**: `http://localhost:3000`
- **Max Task Length**: 500 characters

### Frontend Configuration
- **API Base URL**: `http://localhost:8080`
- **Development Port**: 3000
- **Package Manager**: pnpm (recommended)
- **ConnectRPC**: Full protocol support with JSON transport

## 📁 Project Structure

```text
todoTist/
├── backend/
│   ├── server.go           # Main server implementation
│   ├── server_test.go      # Comprehensive test suite
│   ├── go.mod             # Go dependencies
│   ├── .gitignore         # Excludes generated *.pb.go files
│   ├── todo.proto         # Protocol Buffer definition
│   └── todo/
│       └── v1/
│           ├── todo.pb.go     # Generated Protocol Buffer types (from go generate)
│           └── todo.connect.go # ConnectRPC handlers
├── frontend/
│   ├── src/
│   │   ├── app/
│   │   │   ├── page.tsx       # Main page
│   │   │   └── layout.tsx     # App layout
│   │   ├── components/
│   │   │   └── TodoList.tsx   # Main todo component
│   │   └── lib/
│   │       ├── todo_pb.ts         # TypeScript types from Protocol Buffers
│   │       └── todo_connect.ts    # Official ConnectRPC client
│   ├── package.json           # Frontend dependencies
│   ├── pnpm-lock.yaml        # pnpm lock file
│   └── generate-types.sh     # Protocol Buffer type generation script
│   └── next.config.ts     # Next.js configuration
└── README.md              # This file
```

## 🎯 Key Improvements Made

### Code Quality & Architecture
- ✅ Added comprehensive input validation (empty text, max length)
- ✅ Improved error handling with proper HTTP status codes
- ✅ Added thread-safe operations with RWMutex
- ✅ Implemented task sorting (newest first)
- ✅ Fixed TypeScript linting issues
- ✅ **NEW**: Integrated official ConnectRPC client in frontend
- ✅ **NEW**: Migrated to pnpm for efficient dependency management
- ✅ **NEW**: Full ConnectRPC protocol support with JSON transport
- ✅ **NEW**: Type-safe API communication with Protocol Buffers
- ✅ **NEW**: Optimized Protocol Buffer workflow with go generate

### User Experience
- ✅ Enhanced UI with better error messages
- ✅ Added character counter for task input
- ✅ Improved visual feedback and loading states
- ✅ Better responsive design
- ✅ Added keyboard support (Enter to add task)

### Testing & Validation
- ✅ Extended test coverage with validation scenarios
- ✅ Added benchmark tests for performance
- ✅ Client-side and server-side validation
- ✅ Comprehensive error handling tests

### Documentation
- ✅ Added detailed README with setup instructions
- ✅ Documented API endpoints
- ✅ Improved code comments and structure

## 🔮 Future Enhancements

### Features
- [ ] Task completion/status toggle
- [ ] Task editing capability
- [ ] Task filtering (completed/pending)
- [ ] Task search functionality
- [ ] Task categories/tags
- [ ] Due dates and reminders

### Technical
- [ ] Database persistence (PostgreSQL/SQLite)
- [ ] User authentication and authorization
- [ ] Rate limiting and request throttling
- [ ] Docker containerization
- [ ] Environment-based configuration
- [ ] Logging middleware
- [ ] API versioning
- [ ] Integration tests
- [ ] Performance monitoring

### DevOps
- [ ] CI/CD pipeline
- [ ] Automated testing
- [ ] Production deployment guide
- [ ] Health check endpoints
- [ ] Metrics and monitoring

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## 📄 License

This project is open source and available under the [MIT License](LICENSE).

## 🛠️ Technology Stack

### Backend
- **Language**: Go 1.24+
- **RPC Framework**: ConnectRPC
- **HTTP Server**: net/http with h2c support
- **Storage**: In-memory with RWMutex for thread safety
- **Testing**: Go testing package with benchmarks
- **Build**: Go modules

### Frontend
- **Framework**: Next.js 15 with React 19
- **Language**: TypeScript 5
- **Styling**: Tailwind CSS
- **RPC Client**: ConnectRPC Web Client
- **Package Manager**: pnpm (recommended)
- **Build**: Turbopack for development, standard webpack for production
- **Development**: Hot reloading, TypeScript compilation

### Protocol & Communication
- **API Contract**: Protocol Buffers (todo.proto)
- **Transport**: HTTP/1.1 and HTTP/2 support
- **Serialization**: JSON format (binary optional)
- **Protocol**: ConnectRPC v1 with full compatibility