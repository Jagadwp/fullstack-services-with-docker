# Full-Stack Multi-Service Application

A microservices application with PHP API, Go Scheduler, Python API, and React Frontend. Demonstrates clean service boundaries and proper Docker networking.

---

## Architecture

```
┌─────────────────────┐
│  React Frontend     │  (Browser → localhost:4135)
│  (Port 4135)        │
└──────────┬──────────┘
           │ HTTP
           ▼
┌─────────────────────┐      ┌──────────────────┐
│     PHP API         │◄────▶│  MySQL Database  │
│   (Port 8080)       │      │   (Port 3307)    │
└──────────┬──────────┘      └──────────────────┘
           ▲
           │ Docker Network (fullstack_network)
           │
┌──────────┴──────────┐
│   Go Scheduler      │
│  (Internal)         │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│   Python Service    │
│   (Port 5000)       │
└─────────────────────┘
```

**Key Design:**
- Browser traffic uses `localhost` (host network)
- Container-to-container uses service names (Docker network)
- Each service has independent `docker-compose.yml`
- Shared external network: `fullstack_network`

---

## Project Structure

```
fullstack-assignment/
├── README.md
├── setup-network.sh          # Creates shared Docker network
├── start-all.sh              # Starts all services
├── stop-all.sh               # Stops all services
│
├── php-api/
│   ├── docker-compose.yml
│   ├── Dockerfile
│   ├── public/index.php
│   └── src/
│       ├── Database.php
│       └── UserController.php
│
├── go-scheduler/
│   ├── docker-compose.yml
│   ├── Dockerfile
│   ├── main.go
│   └── data/                 # Stores user JSON files
│       ├── incoming/
│       └── processed/
│
├── python-service/
│   ├── docker-compose.yml
│   ├── Dockerfile
│   ├── app.py
│   └── data/
│       └── received/
│
└── react-frontend/
    ├── docker-compose.yml
    ├── Dockerfile
    └── src/App.jsx
```

---

## Services

### 1. PHP API (User Management)
- **Stack:** PHP 8.2 + Apache + MySQL 8.0
- **Port:** 8080
- **Endpoints:**
  - `POST /users` – Create user (validates email, detects duplicates)
  - `GET /users` – List all users
  - `GET /health` – Health check
- **Features:** Input validation, Xdebug enabled, auto-creates schema

### 2. Go Scheduler
- **Stack:** Go 1.24 (Alpine)
- **Purpose:** Background worker for orchestrating data flow

**Workflow:**
1. Periodically creates users via PHP API (`POST /users`)
2. Saves response to `/data/incoming/user_{id}.json`
3. Scans incoming files for names starting with `"David"`
4. Forwards matching users to Python service (`POST /process`)
5. Moves processed files to `/data/processed/`
6. Repeats on configurable interval

**File-Based State:**
```
/data
 ├── incoming/       # New users from PHP API
 └── processed/      # Already evaluated users
```

**Config (Environment Variables):**
- `PHP_API_BASE_URL` – PHP API URL (default: `http://php_api:80`)
- `PYTHON_API_BASE_URL` – Python service URL (default: `http://python_service:5000`)
- `SCHEDULER_INTERVAL_SECONDS` – Polling interval (default: `30`)

### 3. Python Service
- **Stack:** Flask + Python 3.11 (Alpine)
- **Port:** 5000 (internal)
- **Purpose:** Receives and persists filtered users from Go scheduler

**Endpoints:**
- `POST /process` – Receives user data, saves to `/data/received/user_{id}.json`
- `GET /health` – Health check

**File Storage:**
```
/data
 └── received/       # Users forwarded by Go scheduler
```

### 4. React Frontend
- **Stack:** React 18 + Vite + Tailwind CSS
- **Port:** 4135
- **Features:**
  - User creation form with validation
  - Display created user confirmation
  - Toggle-able user list (show/hide)
  - Responsive design

---

## How to Run

### Prerequisites
- Docker Desktop (running)
- Ports available: 3307, 4135, 5000, 8080

### Quick Start

```bash
# 1. Make scripts executable
chmod +x setup-network.sh start-all.sh stop-all.sh

# 2. Start all services
./start-all.sh

# 3. Access the application
# Frontend: http://localhost:4135
# PHP API:  http://localhost:8080/users
```

### Stop Services

```bash
./stop-all.sh
```

### Individual Service Management

```bash
# Start specific service
cd <service-folder>
docker-compose up -d

# View logs
docker-compose logs -f

# Restart after code changes
docker-compose restart

# Rebuild
docker-compose up --build
```

---

## Testing the Flow

### Verify Files

```bash
# Check incoming files
cd go-scheduler
docker compose exec go-scheduler sh -c "ls /data/incoming"
## Result
user_1.json  user_2.json  user_3.json

# Check processed files
cd go-scheduler
docker compose exec go-scheduler sh -c "ls /data/processed"
## Result
user_1.json  user_2.json  user_3.json

# Check received files
cd python-api
docker compose exec python-api sh -c "ls /data/received"
## Result
user_1.json  user_2.json  user_3.json
```

---

## Design Decisions

**Why separate docker-compose files?**
- Mimics real microservices (independent deployment)
- Each service is self-contained and testable

**Why external Docker network?**
- Allows services in separate compose files to communicate
- Simulates production service mesh behavior

**Why file-based state in Go?**
- Stateless in memory
- Survives container restarts
- Simple and debuggable

---

## Troubleshooting

**Services can't communicate:**
```bash
docker network inspect fullstack_network
./setup-network.sh  # Recreate if needed
```

**Port already in use:**
```bash
lsof -i :8080
kill -9 <PID>
```

**Fresh start:**
```bash
./stop-all.sh
docker system prune -a --volumes
./start-all.sh
```

---

## Author

**Jagad Wijaya Purnomo**
