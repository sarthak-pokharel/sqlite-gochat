# SQLite D1 Go Chat

Multi-channel chat aggregator service that listens for incoming messages from external platforms (WhatsApp Business, Telegram, Instagram, etc.) and provides a unified API for managing conversations.

## Architecture

The service acts as a headless backend microservice designed to integrate with NestJS applications. It automatically creates users when they first message any connected platform, manages conversations, and emits events via Redis Pub/Sub for real-time updates.

## Tech Stack

- Go 1.25.5
- Echo v4 (HTTP framework)
- SQLite (modernc.org/sqlite - pure Go driver)
- Redis (event emission)
- JWT authentication

## Features

- Organization and channel management
- Auto-create external users on first message
- Conversation tracking with assignment support
- Message storage and retrieval
- Webhook receivers for external platforms
- Event emission to NestJS via Redis Pub/Sub
- CRM-agnostic design (no internal user management)


## Setup

1. Copy environment configuration:
```bash
cp .env.example .env
```

2. Configure environment variables in `.env`

3. Build and run:
```bash
make run
```

The server starts on port 8080. Database migrations run automatically on startup.

## API Endpoints

### Organizations
- `GET /api/v1/organizations` - List organizations
- `POST /api/v1/organizations` - Create organization
- `GET /api/v1/organizations/:id` - Get organization
- `PUT /api/v1/organizations/:id` - Update organization
- `DELETE /api/v1/organizations/:id` - Delete organization

### Conversations
- `GET /api/v1/conversations` - List conversations
- `GET /api/v1/conversations/:id` - Get conversation
- `PATCH /api/v1/conversations/:id/assign` - Assign conversation
- `PATCH /api/v1/conversations/:id/status` - Update status

### Messages
- `GET /api/v1/conversations/:id/messages` - List messages
- `POST /api/v1/conversations/:id/messages` - Send message

### Webhooks
- `POST /webhooks/:platform` - Receive webhook from external platform

## Events Emitted to NestJS

The service publishes events to Redis channels for NestJS consumption:

- `chat.conversation.new` - New conversation created
- `chat.message.new` - New message received
- `chat.conversation.assigned` - Conversation assigned to agent
- `chat.conversation.status_changed` - Conversation status updated

## Database Schema

- `organizations` - Business accounts
- `chat_channels` - Connected platform accounts
- `external_users` - Customers from external platforms
- `conversations` - Chat sessions
- `messages` - Message content
- `webhook_events` - Event log for debugging

## Development Principles

- MVC/Layered architecture
- SOLID principles
- Dependency injection
- Interface-based design
- KISS (Keep It Simple, Stupid)
