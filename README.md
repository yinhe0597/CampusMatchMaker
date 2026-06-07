<div align="center">

# 🏫 Campus Collab

> ✨ Smarter Class Collaboration · Timetable Sharing · Time Polling · One-Click Alignment

![Version](https://img.shields.io/badge/version-open0.0.4-blue) ![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go) ![Vue](https://img.shields.io/badge/Vue-3.5-4FC08D?logo=vue.js) ![MySQL](https://img.shields.io/badge/MySQL-8.0-4479A1?logo=mysql) ![Redis](https://img.shields.io/badge/Redis-7.0-DC382D?logo=redis) ![Docker](https://img.shields.io/badge/Docker-ready-2496ED?logo=docker) ![License](https://img.shields.io/badge/license-MIT-green) ![Tests](https://img.shields.io/badge/tests-55%2B%20passing-brightgreen)

</div>

---

<div align="center">

## 🌐 Language · 語言

</div>

<details open>
<summary><b>🇬🇧 English</b> (click to expand)</summary>

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

</details>

<details>
<summary><b>🇨🇳 简体中文</b> (点击展开)</summary>

---

### 🎯 项目简介

校园协作平台是一个面向大学生班级的 **课表共享与时间投票工具**，帮助班级成员快速对齐空闲时间，告别"拉群 — 接龙 — 数票"的低效流程。

- 🗓️ **班级公共课表** — 管理员录入课表，成员一键继承
- ⏰ **空闲时段计算** — 基于全员课表自动计算共同空闲时间
- 📊 **时间投票** — 发起投票、成员表态、自动统计结果
- 🔒 **隐私安全** — 学号 SHA-256 哈希存储，课表数据按班级隔离

### 🚀 快速开始

| 工具 | 版本 |
|------|------|
| Go | ≥ 1.25 |
| Node.js | ≥ 18 |
| pnpm | ≥ 8 |
| MySQL | 8.0 |
| Redis | 7.0 |

```bash
# 1. 启动依赖服务
docker-compose up -d

# 2. 数据库迁移
cd server && go run cmd/migrate/main.go

# 3. 启动后端 (localhost:8080)
go run cmd/main.go

# 4. 启动前端 (localhost:5173)
cd web && pnpm install && pnpm dev
```

### 📁 项目结构

```
wingo/time_tag_peidui/
├── server/                     # Go 后端 (Gin + GORM)
│   ├── api/v1/                 #   路由注册 + 依赖注入
│   ├── cmd/                    #   入口程序
│   ├── internal/
│   │   ├── domain/             #   领域模型 + 仓储 (user/class/timetable/poll)
│   │   ├── engine/timeslot/    #   空闲时段计算引擎 (纯函数)
│   │   ├── handler/            #   HTTP 处理器 + 中间件
│   │   ├── infra/              #   基础设施 (config/db/logger/cache/redis)
│   │   └── service/            #   业务服务层 + DTO
│   ├── migrations/             #   数据库迁移 SQL
│   └── pkg/                    #   工具包 (errors/response/utils)
├── web/                        # Vue3 前端 (Vite + Element Plus)
│   └── src/
│       ├── api/                #   Axios API 封装
│       ├── stores/             #   Pinia 状态管理
│       ├── views/              #   页面组件 (auth/class/timetable/poll)
│       └── router/             #   Vue Router + 路由守卫
└── docs/                       # 项目文档
```

### 🧩 核心功能

#### 🔐 认证模块 (4 API)
| 功能 | 接口 |
|------|------|
| 注册 | `POST /api/v1/auth/register` |
| 登录 | `POST /api/v1/auth/login` |
| 刷新令牌 | `POST /api/v1/auth/refresh-token` |
| 个人信息 | `GET /api/v1/auth/me` |

#### 🏫 班级模块 (7 API)
| 功能 | 接口 |
|------|------|
| 班级列表 | `GET /api/v1/classes` |
| 创建班级 | `POST /api/v1/classes` |
| 查邀请码 | `GET /api/v1/classes/by-code/:code` |
| 班级详情 | `GET /api/v1/classes/:id` |
| 加入班级 | `POST /api/v1/classes/:id/join` |
| 成员列表 | `GET /api/v1/classes/:id/members` |
| 移除成员 | `DELETE /api/v1/classes/:id/members/:userId` |

#### 📅 课表模块 (10 API)
- 班级课表录入 / 查看 / 更新
- 个人课表增删改查（含自动继承）
- 纠错提交与审核

#### 🗳️ 投票模块 (10 API)
- 创建投票 + 自动推荐空闲时段
- 投票 (yes/no/maybe) + 实时结果统计
- 确认最终时段

> **共 32 条 API**，详见 [`docs/02-API接口设计.md`](docs/02-API接口设计.md)

### ⚙️ 技术架构

```
┌─────────────────────────────────────────┐
│              Vue3 前端                    │
│  Element Plus · Pinia · Vue Router       │
├─────────────────────────────────────────┤
│           Gin HTTP 服务 (8080)            │
│  Handler → Service → Repository          │
├──────────────┬──────────────────────────┤
│   GORM/MySQL │  Timeslot Engine (纯函数)  │
│   13 张表     │  空闲时段计算               │
└──────────────┴──────────────────────────┘
```

### 🔒 隐私设计

| 数据 | 存储方式 |
|------|---------|
| 学号 | SHA-256 哈希 (可查询，不可逆) |
| 密码 | bcrypt 哈希 |
| JWT | HS256 签名，短期有效 |
| 课表 | 班级范围隔离，仅成员可见 |

### 📊 版本历史

| 版本 | 日期 | 里程碑 |
|------|------|--------|
| **open0.0.4** | 2026-06-07 | 🎨 课表删除+编辑 · 投票Tab完善 · 全栈Docker化 · CI/CD · 审查修复(5项) |
| **open0.0.3** | 2026-06-07 | ⚡ E2E 17 测试全通过 · Redis 缓存层 · 前端多班级聚合 · 性能压测 · DB 迁移修复 |
| **open0.0.2** | 2026-06-07 | 🧪 单元测试 38 个全覆盖 · 引擎排序 Bug 修复 · auto_recommend 根因定位 |
| **open0.0.1** | 2026-06-07 | 🚀 阶段一完成：认证+班级+课表+投票 32 API · 代码审查修复 9 项 |

### 📈 开发阶段

| 阶段 | 内容 | 进度 |
|------|------|:--:|
| 🏗️ 基础设施 | 脚手架 + 环境 + 文档 + 13 表迁移 | ✅ 100% |
| 🚀 阶段一 | 认证 + 班级 + 课表 + 投票 (32 API + 23 前端) | ✅ 100% |
| 🔧 阶段二 | 单元测试 + Bug 修复 + 集成验证 | ✅ 100% |
| 🚢 阶段三 | E2E 测试 + Redis 缓存 + 性能压测 + 前端聚合 | ✅ 100% |
| 🎨 阶段四 | 前端完善 + Docker 部署 + CI/CD | ✅ 100% |

> 详细进度见 [`docs/06-开发进度.md`](docs/06-开发进度.md) | 开发日志见 [`docs/2026-06-07-开发日志.md`](docs/2026-06-07-开发日志.md)

### 🤝 贡献者

| 姓名 | 邮箱 |
|------|------|
| Vincentluo | yhkjsj@foxmail.com |

</details>

<details>
<summary><b>🇭🇰 繁體中文</b> (點擊展開)</summary>

---

### 🎯 項目簡介

校園協作平台是一個面向大學生班級的 **課表共享與時間投票工具**，幫助班級成員快速對齊空閒時間，告別"拉群 — 接龍 — 數票"的低效流程。

- 🗓️ **班級公共課表** — 管理員錄入課表，成員一鍵繼承
- ⏰ **空閒時段計算** — 基於全員課表自動計算共同空閒時間
- 📊 **時間投票** — 發起投票、成員表態、自動統計結果
- 🔒 **隱私安全** — 學號 SHA-256 哈希存儲，課表數據按班級隔離

### 🚀 快速開始

| 工具 | 版本 |
|------|------|
| Go | ≥ 1.25 |
| Node.js | ≥ 18 |
| pnpm | ≥ 8 |
| MySQL | 8.0 |
| Redis | 7.0 |

```bash
# 1. 啟動依賴服務
docker-compose up -d

# 2. 數據庫遷移
cd server && go run cmd/migrate/main.go

# 3. 啟動後端 (localhost:8080)
go run cmd/main.go

# 4. 啟動前端 (localhost:5173)
cd web && pnpm install && pnpm dev
```

### 📁 項目結構

```
wingo/time_tag_peidui/
├── server/                     # Go 後端 (Gin + GORM)
│   ├── api/v1/                 #   路由註冊 + 依賴注入
│   ├── cmd/                    #   入口程式
│   ├── internal/
│   │   ├── domain/             #   領域模型 + 倉儲 (user/class/timetable/poll)
│   │   ├── engine/timeslot/    #   空閒時段計算引擎 (純函數)
│   │   ├── handler/            #   HTTP 處理器 + 中間件
│   │   ├── infra/              #   基礎設施 (config/db/logger/cache/redis)
│   │   └── service/            #   業務服務層 + DTO
│   ├── migrations/             #   數據庫遷移 SQL
│   └── pkg/                    #   工具包 (errors/response/utils)
├── web/                        # Vue3 前端 (Vite + Element Plus)
│   └── src/
│       ├── api/                #   Axios API 封裝
│       ├── stores/             #   Pinia 狀態管理
│       ├── views/              #   頁面組件 (auth/class/timetable/poll)
│       └── router/             #   Vue Router + 路由守衛
└── docs/                       # 項目文檔
```

### 🧩 核心功能

#### 🔐 認證模組 (4 API)
| 功能 | 接口 |
|------|------|
| 註冊 | `POST /api/v1/auth/register` |
| 登錄 | `POST /api/v1/auth/login` |
| 刷新令牌 | `POST /api/v1/auth/refresh-token` |
| 個人資訊 | `GET /api/v1/auth/me` |

#### 🏫 班級模組 (7 API)
| 功能 | 接口 |
|------|------|
| 班級列表 | `GET /api/v1/classes` |
| 創建班級 | `POST /api/v1/classes` |
| 查邀請碼 | `GET /api/v1/classes/by-code/:code` |
| 班級詳情 | `GET /api/v1/classes/:id` |
| 加入班級 | `POST /api/v1/classes/:id/join` |
| 成員列表 | `GET /api/v1/classes/:id/members` |
| 移除成員 | `DELETE /api/v1/classes/:id/members/:userId` |

#### 📅 課表模組 (10 API)
- 班級課表錄入 / 查看 / 更新
- 個人課表增刪改查（含自動繼承）
- 糾錯提交與審核

#### 🗳️ 投票模組 (10 API)
- 創建投票 + 自動推薦空閒時段
- 投票 (yes/no/maybe) + 即時結果統計
- 確認最終時段

> **共 32 條 API**，詳見 [`docs/02-API接口設計.md`](docs/02-API接口設計.md)

### ⚙️ 技術架構

```
┌─────────────────────────────────────────┐
│              Vue3 前端                    │
│  Element Plus · Pinia · Vue Router       │
├─────────────────────────────────────────┤
│           Gin HTTP 服務 (8080)            │
│  Handler → Service → Repository          │
├──────────────┬──────────────────────────┤
│   GORM/MySQL │  Timeslot Engine (純函數)  │
│   13 張表     │  空閒時段計算               │
└──────────────┴──────────────────────────┘
```

### 🔒 隱私設計

| 數據 | 存儲方式 |
|------|---------|
| 學號 | SHA-256 哈希 (可查詢，不可逆) |
| 密碼 | bcrypt 哈希 |
| JWT | HS256 簽名，短期有效 |
| 課表 | 班級範圍隔離，僅成員可見 |

### 📊 版本歷史

| 版本 | 日期 | 里程碑 |
|------|------|--------|
| **open0.0.4** | 2026-06-07 | 🎨 課表刪除+編輯 · 投票Tab完善 · 全棧Docker化 · CI/CD · 審查修復(5項) |
| **open0.0.3** | 2026-06-07 | ⚡ E2E 17 測試全通過 · Redis 緩存層 · 前端多班級聚合 · 效能壓測 · DB 遷移修復 |
| **open0.0.2** | 2026-06-07 | 🧪 單元測試 38 個全覆蓋 · 引擎排序 Bug 修復 · auto_recommend 根因定位 |
| **open0.0.1** | 2026-06-07 | 🚀 階段一完成：認證+班級+課表+投票 32 API · 代碼審查修復 9 項 |

### 📈 開發階段

| 階段 | 內容 | 進度 |
|------|------|:--:|
| 🏗️ 基礎設施 | 腳手架 + 環境 + 文檔 + 13 表遷移 | ✅ 100% |
| 🚀 階段一 | 認證 + 班級 + 課表 + 投票 (32 API + 23 前端) | ✅ 100% |
| 🔧 階段二 | 單元測試 + Bug 修復 + 集成驗證 | ✅ 100% |
| 🚢 階段三 | E2E 測試 + Redis 緩存 + 效能壓測 + 前端聚合 | ✅ 100% |
| 🎨 階段四 | 前端完善 + Docker 部署 + CI/CD | ✅ 100% |

> 詳細進度見 [`docs/06-開發進度.md`](docs/06-開發進度.md) | 開發日誌見 [`docs/2026-06-07-開發日誌.md`](docs/2026-06-07-開發日誌.md)

### 🤝 貢獻者

| 姓名 | 郵箱 |
|------|------|
| Vincentluo | yhkjsj@foxmail.com |

</details>

---

<p align="center">
  <sub>Made with ❤️ for campus collaboration</sub>
</p>
