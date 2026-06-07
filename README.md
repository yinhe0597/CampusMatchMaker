<div align="center">

# 🏫 Campus Collab

> ✨ Timetable Sharing · Free-Slot Engine · Time Polling · One-Click Class Collaboration

![Version](https://img.shields.io/badge/version-open0.0.4-blue) ![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go) ![Vue](https://img.shields.io/badge/Vue-3.5-4FC08D?logo=vue.js) ![MySQL](https://img.shields.io/badge/MySQL-8.0-4479A1?logo=mysql) ![Redis](https://img.shields.io/badge/Redis-7.0-DC382D?logo=redis) ![Docker](https://img.shields.io/badge/Docker-ready-2496ED?logo=docker) ![License](https://img.shields.io/badge/license-MIT-green) ![Tests](https://img.shields.io/badge/tests-55%2B%20passing-brightgreen)

</div>

---

<div align="center">

**English** &nbsp;|&nbsp; [**简体中文**](README.zh-CN.md) &nbsp;|&nbsp; [**繁體中文**](README.zh-TW.md)

</div>

---

### 🎯 Overview

Campus Collab is a **timetable sharing & time polling tool** for university classes. It helps class members quickly align their free time — no more messy group chats and manual vote counting.

- 🗓️ **Class Timetable** — Admin enters the master timetable; members inherit with one click
- ⏰ **Free Slot Engine** — Computes common free slots based on all members' timetables
- 📊 **Time Polling** — Create polls, vote, and get real-time stats
- 🔒 **Privacy-First** — Student IDs hashed with SHA-256; timetable data scoped per class

### 🚀 Quick Start

| Tool | Version |
|------|---------|
| Go | ≥ 1.25 |
| Node.js | ≥ 18 |
| pnpm | ≥ 8 |
| MySQL | 8.0 |
| Redis | 7.0 |

```bash
# 1. Start dependencies
docker-compose up -d

# 2. Run database migrations
cd server && go run cmd/migrate/main.go

# 3. Start backend (localhost:8080)
go run cmd/main.go

# 4. Start frontend (localhost:5173)
cd web && pnpm install && pnpm dev
```

### 📁 Project Structure

```
wingo/time_tag_peidui/
├── server/                     # Go backend (Gin + GORM)
│   ├── api/v1/                 #   Route registration + DI
│   ├── cmd/                    #   Entry points
│   ├── internal/
│   │   ├── domain/             #   Domain models + repos (user/class/timetable/poll)
│   │   ├── engine/timeslot/    #   Free-slot computation engine (pure functions)
│   │   ├── handler/            #   HTTP handlers + middleware
│   │   ├── infra/              #   Infrastructure (config/db/logger/cache/redis)
│   │   └── service/            #   Business service layer + DTOs
│   ├── migrations/             #   DB migration SQL
│   └── pkg/                    #   Utilities (errors/response/utils)
├── web/                        # Vue3 frontend (Vite + Element Plus)
│   └── src/
│       ├── api/                #   Axios API wrappers
│       ├── stores/             #   Pinia state management
│       ├── views/              #   Page components (auth/class/timetable/poll)
│       └── router/             #   Vue Router + guards
└── docs/                       # Project documentation
```

### 🧩 Core Features

#### 🔐 Authentication (4 APIs)
| Feature | Endpoint |
|---------|----------|
| Register | `POST /api/v1/auth/register` |
| Login | `POST /api/v1/auth/login` |
| Refresh Token | `POST /api/v1/auth/refresh-token` |
| My Profile | `GET /api/v1/auth/me` |

#### 🏫 Class Management (7 APIs)
| Feature | Endpoint |
|---------|----------|
| List Classes | `GET /api/v1/classes` |
| Create Class | `POST /api/v1/classes` |
| Lookup by Code | `GET /api/v1/classes/by-code/:code` |
| Class Details | `GET /api/v1/classes/:id` |
| Join Class | `POST /api/v1/classes/:id/join` |
| Member List | `GET /api/v1/classes/:id/members` |
| Remove Member | `DELETE /api/v1/classes/:id/members/:userId` |

#### 📅 Timetable (10 APIs)
- Class timetable: create / view / update
- Personal timetable: CRUD with auto-inheritance
- Correction: submit & review

#### 🗳️ Polling (10 APIs)
- Create polls with auto-recommended free slots
- Vote (yes / no / maybe) with real-time stats
- Confirm final time slot

> **32 APIs total** — see [`docs/02-API接口设计.md`](docs/02-API接口设计.md)

### ⚙️ Architecture

```
┌─────────────────────────────────────────┐
│            Vue3 Frontend                │
│  Element Plus · Pinia · Vue Router      │
├─────────────────────────────────────────┤
│        Gin HTTP Server (:8080)          │
│  Handler → Service → Repository         │
├──────────────┬──────────────────────────┤
│  GORM/MySQL  │  Timeslot Engine (pure)  │
│  13 tables   │  Free-slot computation   │
└──────────────┴──────────────────────────┘
```

### 🔒 Privacy Design

| Data | Storage |
|------|---------|
| Student ID | SHA-256 hash (queryable, irreversible) |
| Password | bcrypt hash |
| JWT | HS256 signed, short-lived |
| Timetable | Class-scoped, members-only |

### 📊 Version History

| Version | Date | Milestone |
|---------|------|-----------|
| **open0.0.4** | 2026-06-07 | 🎨 Timetable delete+edit · Poll tabs · Full-stack Docker · CI/CD · Bug fixes (5) |
| **open0.0.3** | 2026-06-07 | ⚡ 17 E2E tests passing · Redis cache · Multi-class aggregation · Perf tuning · DB fix |
| **open0.0.2** | 2026-06-07 | 🧪 38 unit tests · Engine sort bug fix · auto_recommend root cause |
| **open0.0.1** | 2026-06-07 | 🚀 Phase 1: Auth + Class + Timetable + Poll 32 APIs · Code review fixes (9) |

### 📈 Development Phases

| Phase | Scope | Progress |
|-------|-------|:--------:|
| 🏗️ Infra | Scaffold + Env + Docs + 13-table migration | ✅ 100% |
| 🚀 Phase 1 | Auth + Class + Timetable + Poll (32 APIs + 23 UI) | ✅ 100% |
| 🔧 Phase 2 | Unit tests + Bug fixes + Integration | ✅ 100% |
| 🚢 Phase 3 | E2E tests + Redis cache + Perf + UI aggregation | ✅ 100% |
| 🎨 Phase 4 | UI polish + Docker deploy + CI/CD | ✅ 100% |

> Detailed progress: [`docs/06-开发进度.md`](docs/06-开发进度.md) | Dev log: [`docs/2026-06-07-开发日志.md`](docs/2026-06-07-开发日志.md)

### 🤝 Contributors

| Name | Email |
|------|-------|
| Vincentluo | yhkjsj@foxmail.com |



---

<p align="center">
  <sub>Made with ❤️ for campus collaboration</sub>
</p>
