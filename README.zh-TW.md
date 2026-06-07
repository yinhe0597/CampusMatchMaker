<div align="center">

# 🏫 Campus Collab

> ✨ 課表共享 · 空閒時段引擎 · 時間投票 · 一鍵班級協作

![Version](https://img.shields.io/badge/version-open0.0.4-blue) ![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go) ![Vue](https://img.shields.io/badge/Vue-3.5-4FC08D?logo=vue.js) ![MySQL](https://img.shields.io/badge/MySQL-8.0-4479A1?logo=mysql) ![Redis](https://img.shields.io/badge/Redis-7.0-DC382D?logo=redis) ![Docker](https://img.shields.io/badge/Docker-ready-2496ED?logo=docker) ![License](https://img.shields.io/badge/license-MIT-green) ![Tests](https://img.shields.io/badge/tests-55%2B%20passing-brightgreen)

</div>

---

<div align="center">

[**English**](README.md) &nbsp;|&nbsp; [**简体中文**](README.zh-CN.md) &nbsp;|&nbsp; **繁體中文**

</div>

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

---

<p align="center">
  <sub>Made with ❤️ for campus collaboration</sub>
</p>
