# FlashTix

Platform pemesanan tiket event high-concurrency menggunakan Go dan React.

## Arsitektur

### Backend (Go)
- **Repository Pattern**: Memisahkan akses database dari logic bisnis
- **Middleware Pattern**: Auth, Logging (Sentry), CORS
- **Database**: Prisma Client Go + Neon (PostgreSQL) + Redis untuk optimistic locking

### Frontend (React)
- **Atomic Design**: Komponen UI (Button, Input, Card)
- **Store Pattern (Zustand)**: Global state management
- **Global Types**: Interface/types terpusat

## Setup

### Prerequisites
- Go 1.23+
- Node.js 18+
- Neon PostgreSQL account
- Upstash Redis account (or similar)

1. **Clone repository**
   ```bash
   git clone <repo-url>
   cd FlashTix
   ```

2. **Setup environment**
   ```bash
   cp .env.example .env
   # Edit .env dengan credentials Neon dan Redis
   ```

3. **Setup database schema**
   ```bash
   cd server
   # Install Prisma CLI globally if not already installed
   npm install -g prisma

   # Generate Prisma Client
   npx prisma generate

   # Push schema to database (for development)
   npx prisma db push
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

## Database Configuration

### Neon PostgreSQL
1. Create account di [neon.tech](https://neon.tech)
2. Create new project
3. Copy connection string ke `.env`:
   ```
   DATABASE_URL=postgresql://username:password@hostname/database?sslmode=require
   ```

### Redis (Upstash - RECOMMENDED)
1. Create account di [upstash.com](https://upstash.com)
2. Create Redis database
3. Copy REST API credentials ke `.env`:
   ```
   UPSTASH_REDIS_REST_URL=https://your-database.upstash.io
   UPSTASH_REDIS_REST_TOKEN=your-rest-token
   ```

   **Catatan**: FlashTix menggunakan Upstash REST API untuk koneksi Redis yang lebih reliable dan tidak memerlukan TCP connection.

### Redis (Local Development)
Untuk development lokal, Anda bisa menggunakan:
1. **Docker**: `docker run -d -p 6379:6379 redis:7-alpine`
2. **Local install**: Install Redis server dan set `REDIS_URL=localhost:6379`

**Catatan**: Redis digunakan untuk optimistic locking pada sistem booking tiket untuk mencegah race conditions saat high-concurrency booking.

## Prisma Commands

```bash
# Generate Prisma Client
npx prisma generate

# Push schema changes to database
npx prisma db push

# Create migration
npx prisma migrate dev --name <migration-name>

# View database
npx prisma studio
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
- Type-safe database queries dengan Prisma Client Go