# FlashTix

Platform pemesanan tiket event high-concurrency menggunakan Go dan React.

## Arsitektur

### Backend (Go)
- **Repository Pattern**: Memisahkan akses database dari logic bisnis
- **Middleware Pattern**: Auth, Logging (Sentry), CORS
- **Database**: Neon (PostgreSQL) + Redis untuk optimistic locking

### Frontend (React)
- **Atomic Design**: Komponen UI (Button, Input, Card)
- **Store Pattern (Zustand)**: Global state management
- **Global Types**: Interface/types terpusat

## Setup

1. **Clone repository**
   ```bash
   git clone <repo-url>
   cd FlashTix
   ```

2. **Setup environment**
   ```bash
   cp .env.example .env
   # Edit .env dengan credentials yang sesuai
   ```

3. **Start databases**
   ```bash
   docker-compose up -d
   ```

4. **Setup backend**
   ```bash
   cd server
   go mod tidy
   go run cmd/main.go
   ```

5. **Setup frontend**
   ```bash
   cd client
   npm install
   npm run dev
   ```

## API Endpoints

- `GET /api/events` - Get all events
- `POST /api/events` - Create event (auth required)
- `POST /api/tickets/reserve` - Reserve seat (auth required)
- `POST /api/tickets/confirm` - Confirm purchase (auth required)

## Features

- Optimistic locking untuk seat reservation menggunakan Redis
- JWT authentication
- Atomic UI components
- Centralized state management dengan Zustand