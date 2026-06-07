# 🏫 校园协作平台 — Campus Collab

> ✨ 让班级协作更简单 · 课表共享 · 时间投票 · 一键对齐

![Version](https://img.shields.io/badge/version-open0.0.4-blue) ![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go) ![Vue](https://img.shields.io/badge/Vue-3.5-4FC08D?logo=vue.js) ![MySQL](https://img.shields.io/badge/MySQL-8.0-4479A1?logo=mysql) ![Redis](https://img.shields.io/badge/Redis-7.0-DC382D?logo=redis) ![Docker](https://img.shields.io/badge/Docker-ready-2496ED?logo=docker) ![License](https://img.shields.io/badge/license-MIT-green) ![Tests](https://img.shields.io/badge/tests-55%2B%20passing-brightgreen)

---

## 🎯 项目简介

校园协作平台是一个面向大学生班级的 **课表共享与时间投票工具**，帮助班级成员快速对齐空闲时间，告别"拉群 — 接龙 — 数票"的低效流程。

- 🗓️ **班级公共课表** — 管理员录入课表，成员一键继承
- ⏰ **空闲时段计算** — 基于全员课表自动计算共同空闲时间
- 📊 **时间投票** — 发起投票、成员表态、自动统计结果
- 🔒 **隐私安全** — 学号 SHA-256 哈希存储，课表数据按班级隔离

---

## 🚀 快速开始

### 环境要求

| 工具 | 版本 |
|------|------|
| Go | ≥ 1.25 |
| Node.js | ≥ 18 |
| pnpm | ≥ 8 |
| MySQL | 8.0 |
| Redis | 7.0 |

### 本地启动

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

---

## 📁 项目结构

```
wingo/time_tag_peidui/
├── server/                     # 🔧 Go 后端 (Gin + GORM)
│   ├── api/v1/                 #    路由注册 + 依赖注入
│   ├── cmd/                    #    入口程序
│   ├── internal/
│   │   ├── domain/             #    领域模型 + 仓储 (user/class/timetable/poll)
│   │   ├── engine/timeslot/    #    ⚙️ 空闲时段计算引擎 (纯函数)
│   │   ├── handler/            #    HTTP 处理器 + 中间件
│   │   ├── infra/              #    基础设施 (config/db/logger/cache/redis)
│   │   └── service/            #    业务服务层 + DTO
│   ├── migrations/             #    数据库迁移 SQL
│   └── pkg/                    #    工具包 (errors/response/utils)
├── web/                        # 🎨 Vue3 前端 (Vite + Element Plus)
│   └── src/
│       ├── api/                #    Axios API 封装
│       ├── stores/             #    Pinia 状态管理
│       ├── views/              #    页面组件 (auth/class/timetable/poll)
│       └── router/             #    Vue Router + 路由守卫
└── docs/                       # 📚 项目文档
    ├── 01-数据库设计.md
    ├── 02-API接口设计.md
    ├── 03-系统架构设计.md
    ├── 04-项目工程指南.md
    ├── 05-开发规范补充.md
    └── 06-开发进度.md
```

---

## 🧩 核心功能

### 🔐 认证模块 (4 API)
| 功能 | 接口 |
|------|------|
| 📝 注册 | `POST /api/v1/auth/register` |
| 🔑 登录 | `POST /api/v1/auth/login` |
| 🔄 刷新令牌 | `POST /api/v1/auth/refresh-token` |
| 👤 个人信息 | `GET /api/v1/auth/me` |

### 🏫 班级模块 (7 API)
| 功能 | 接口 |
|------|------|
| 📋 班级列表 | `GET /api/v1/classes` |
| ✨ 创建班级 | `POST /api/v1/classes` |
| 🔍 查邀请码 | `GET /api/v1/classes/by-code/:code` |
| 📄 班级详情 | `GET /api/v1/classes/:id` |
| 🚪 加入班级 | `POST /api/v1/classes/:id/join` |
| 👥 成员列表 | `GET /api/v1/classes/:id/members` |
| ❌ 移除成员 | `DELETE /api/v1/classes/:id/members/:userId` |

### 📅 课表模块 (10 API)
- 📖 班级课表录入/查看/更新
- 👤 个人课表增删改查（含自动继承）
- 🛠️ 纠错提交与审核

### 🗳️ 投票模块 (10 API)
- 📝 创建投票 + ⚡ 自动推荐空闲时段
- ✅ 投票 (yes/no/maybe) + 📊 实时结果统计
- 🎯 确认最终时段

> **共 32 条 API**，详见 [`docs/02-API接口设计.md`](docs/02-API接口设计.md)

---

## ⚙️ 技术架构

```
┌─────────────────────────────────────────┐
│              Vue3 前端                    │
│  Element Plus · Pinia · Vue Router       │
├─────────────────────────────────────────┤
│           Gin HTTP 服务 (8080)            │
│  Handler → Service → Repository          │
├──────────────┬──────────────────────────┤
│   GORM/MySQL │  timeslot Engine (纯函数)  │
│   13 张表     │  空闲时段计算               │
└──────────────┴──────────────────────────┘
```

---

## 🔒 隐私设计

| 数据 | 存储方式 |
|------|---------|
| 学号 | SHA-256 哈希 (可查询，不可逆) |
| 密码 | bcrypt 哈希 |
| JWT | HS256 签名，短期有效 |
| 课表 | 班级范围隔离，仅成员可见 |

---

## 📊 版本历史

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

---

## 🤝 贡献者

| 姓名 | 邮箱 |
|------|------|
| Vincentluo | yhkjsj@foxmail.com |

---

<p align="center">
  <sub>Made with ❤️ for campus collaboration</sub>
</p>
