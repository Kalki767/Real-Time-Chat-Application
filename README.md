# Real-Time Chat Application

A modern, scalable real-time chat application built with Go, MongoDB, and WebSocket technology. This application implements clean architecture principles to ensure maintainability, testability, and scalability.

## Features

- User Management
  - User registration with email and username validation
  - Secure authentication using JWT tokens
  - Complete profile management including:
    - Update email, username and password
    - Delete account with associated data cleanup
    - View and modify profile information
  - Advanced user lookup functionality:
    - Search by username
    - Search by email
    - Get user by ID
    
- Real-time Chat
  - Instant one-on-one messaging with WebSocket
  - Persistent message storage with MongoDB
  - Message features include:
    - Timestamp tracking
    - Read receipts
    - Message history pagination
    - Real-time delivery status
  - Chat session management
    - Active/inactive status
    - Typing indicators
    - Online presence detection

## Technical Stack

- **Backend**: 
  - Go (Golang) for high-performance server-side operations
  - Gorilla WebSocket for real-time communication
  - Clean architecture implementation with domain-driven design
  
- **Database**: 
  - MongoDB for flexible document storage
  - Optimized indexes for quick lookups
  - Efficient data modeling for chat and user collections
  
- **Real-time Communication**: 
  - WebSocket protocol for bi-directional communication
  - Efficient connection pooling
  - Heartbeat mechanism for connection health
  
- **Testing**: 
  - Comprehensive unit tests using Go testing package
  - Mock interfaces with testify/mock
  - Integration tests for critical paths
  - Test coverage reporting

## Project Structure

The project follows clean architecture principles with clear separation of concerns across multiple layers:

- **Domain Layer**
  - Contains enterprise business rules and entities
  - Defines core interfaces and models
  - Independent of external frameworks and tools
  - Located in `/domain` directory

- **Use Case Layer**
  - Implements application-specific business rules
  - Orchestrates data flow between entities
  - Contains business logic implementations
  - Located in `/usecase` directory

- **Interface Adapters Layer**
  - Converts data between use cases and external agencies
  - Contains controllers, presenters, and gateways
  - Handles framework-specific implementations
  - Located in `/delivery` and `/repository` directories

- **Frameworks & Drivers Layer**
  - Contains frameworks and tools like databases
  - Implements interfaces defined by inner layers
  - Handles external communications
  - Located in `/infrastructure` directory
  
