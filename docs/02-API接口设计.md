# 校园协作平台 · API 接口设计文档

> 版本：v0.1.0（草案）  
> 基础路径：`/api/v1`  
> 协议：HTTPS  
> 数据格式：JSON

---

## 目录

1. 接口总则
2. 认证模块
3. 班级模块
4. 课表模块
5. 空闲时段计算模块
6. 投票模块
7. 用户模块
8. 错误码清单

---

## 1. 接口总则

### 1.1 RESTful 设计规范

- 资源名使用**复数形式**：`/classes`、`/polls`、`/timetables`
- 使用标准 HTTP 方法：GET（查）、POST（增）、PUT（改）、DELETE（删）
- 状态码语义：200 成功、201 创建成功、400 参数错误、401 未认证、403 无权限、404 不存在、500 服务器错误

### 1.2 统一响应格式

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": { ... },
  "timestamp": 1700000000
}
```

**分页响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [ ... ],
    "total": 100,
    "page": 1,
    "page_size": 20
  },
  "timestamp": 1700000000
}
```

**错误响应：**
```json
{
  "code": 1101,
  "message": "学号格式不正确",
  "data": null,
  "timestamp": 1700000000
}
```

### 1.3 认证方式

所有需认证的接口，在请求头携带 JWT Token：
```
Authorization: Bearer <token>
```

### 1.4 分页参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| page | int | 1 | 页码，从 1 开始 |
| page_size | int | 20 | 每页条数，最大 100 |

### 1.5 通用请求头

| Header | 值 | 说明 |
|--------|------|------|
| Content-Type | application/json | 请求体格式 |
| Authorization | Bearer {token} | 认证令牌 |
| X-Request-ID | uuid | 请求追踪ID（可选，客户端生成） |

---

## 2. 认证模块（✅ 已实现）

### 2.1 注册

```
POST /api/v1/auth/register
```

**请求体：**
```json
{
  "student_id": "2024001234",
  "password": "StrongP@ss123",
  "nickname": "小明",
  "school_id": 1
}
```

**成功响应 (201)：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 1,
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_at": 1780882632
  }
}
```

> **实现说明**：学号经 SHA-256 哈希后存储，JWT 内含 `user_id` 和原始学号。

### 2.2 登录

```
POST /api/v1/auth/login
```

**请求体：**
```json
{
  "student_id": "2024001234",
  "password": "StrongP@ss123"
}
```

**成功响应 (200)：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 1,
    "nickname": "小明",
    "avatar_url": null,
    "school_id": 1,
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_at": 1780882632
  }
}
```

### 2.3 刷新令牌

```
POST /api/v1/auth/refresh-token
认证：需要
```

**成功响应 (200)：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_at": 1780882632
  }
}
```

### 2.4 获取当前用户信息

```
GET /api/v1/auth/me
认证：需要
```

**成功响应 (200)：**
```json
{
  "code": 0,
  "data": {
    "user_id": 1,
    "student_id": "2024****34",
    "nickname": "小明",
    "avatar_url": "https://...",
    "school": { "id": 1, "name": "XX大学" },
    "auth_status": 1,
    "privacy_level": 1,
    "created_at": "2025-09-01T10:00:00Z"
  }
}
```

---

## 3. 班级模块（🔲 待实现）

### 3.1 创建班级

```
POST /api/v1/classes
认证：需要
```

**请求体：**
```json
{
  "school_id": 1,
  "grade": "2024",
  "department": "计算机学院",
  "name": "计科2401班",
  "code": "CS2401"
}
```

**成功响应 (201)：**
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "name": "计科2401班",
    "code": "CS2401",
    "invite_code": "A3X9K2",
    "creator_user_id": 1,
    "timetable_status": 0,
    "member_count": 1,
    "created_at": "2025-09-01T10:00:00Z"
  }
}
```

### 3.2 获取班级详情

```
GET /api/v1/classes/:id
认证：需要
```

**成功响应 (200)：**
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "name": "计科2401班",
    "grade": "2024",
    "department": "计算机学院",
    "school": { "id": 1, "name": "XX大学" },
    "timetable_status": 1,
    "member_count": 30,
    "created_at": "2025-09-01T10:00:00Z"
  }
}
```

### 3.3 加入班级

```
POST /api/v1/classes/:id/join
认证：需要
```

**请求体：**
```json
{
  "invite_code": "A3X9K2"
}
```

**成功响应 (200)：**
```json
{
  "code": 0,
  "message": "加入成功",
  "data": {
    "class_id": 1,
    "user_id": 2,
    "role": "member",
    "timetable_inherited": true,
    "inherited_count": 15
  }
}
```

> `timetable_inherited` 表示是否自动继承了班级课表，`inherited_count` 为继承的课表条目数。

### 3.4 获取班级成员列表

```
GET /api/v1/classes/:id/members?page=1&page_size=20
认证：需要
```

**成功响应 (200)：**
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "user_id": 1,
        "nickname": "小明",
        "avatar_url": "https://...",
        "role": "owner",
        "joined_at": "2025-09-01T10:00:00Z"
      }
    ],
    "total": 30,
    "page": 1,
    "page_size": 20
  }
}
```

### 3.5 移除班级成员

```
DELETE /api/v1/classes/:id/members/:userId
认证：需要（仅 owner/admin）
```

**成功响应 (200)：**
```json
{
  "code": 0,
  "message": "已移除"
}
```

---

## 4. 课表模块（🔲 待实现）

### 4.1 录入班级公共课表

```
POST /api/v1/timetables/class/:classId
认证：需要（仅班级成员）
```

**请求体：**
```json
{
  "entries": [
    {
      "day_of_week": 1,
      "period_start": 1,
      "period_end": 2,
      "course_name": "高等数学",
      "teacher": "张教授",
      "room": "教学楼A-301"
    },
    {
      "day_of_week": 3,
      "period_start": 3,
      "period_end": 4,
      "course_name": "数据结构",
      "teacher": "李教授",
      "room": "教学楼B-205"
    }
  ]
}
```

**成功响应 (201)：**
```json
{
  "code": 0,
  "data": {
    "created_count": 2,
    "timetable_status": 1
  }
}
```

### 4.2 获取班级公共课表

```
GET /api/v1/timetables/class/:classId
认证：需要
```

**成功响应 (200)：**
```json
{
  "code": 0,
  "data": {
    "class_id": 1,
    "entries": [
      {
        "id": 1,
        "day_of_week": 1,
        "period_start": 1,
        "period_end": 2,
        "course_name": "高等数学",
        "teacher": "张教授",
        "room": "教学楼A-301",
        "version": 1
      }
    ],
    "contributor": { "user_id": 1, "nickname": "小明" },
    "total_entries": 15
  }
}
```

### 4.3 更新班级公共课表

```
PUT /api/v1/timetables/class/:classId
认证：需要（仅 owner/admin）
```

**请求体（同 4.1）**

**成功响应 (200)：**
```json
{
  "code": 0,
  "message": "课表已更新",
  "data": {
    "updated_count": 3,
    "affected_members": 28
  }
}
```

### 4.4 添加个人课表条目

```
POST /api/v1/timetables/personal
认证：需要
```

**请求体：**
```json
{
  "class_id": 1,
  "day_of_week": 2,
  "period_start": 5,
  "period_end": 6,
  "course_name": "英语选修",
  "source": "personal"
}
```

**成功响应 (201)：**
```json
{
  "code": 0,
  "data": {
    "id": 100,
    "source": "personal",
    "is_overridden": false
  }
}
```

### 4.5 获取个人完整课表

```
GET /api/v1/timetables/personal?class_id=1
认证：需要
```

**成功响应 (200)：**
```json
{
  "code": 0,
  "data": {
    "user_id": 1,
    "class_id": 1,
    "entries": [
      {
        "id": 1,
        "day_of_week": 1,
        "period_start": 1,
        "period_end": 2,
        "course_name": "高等数学",
        "source": "inherited",
        "is_overridden": false
      },
      {
        "id": 100,
        "day_of_week": 2,
        "period_start": 5,
        "period_end": 6,
        "course_name": "英语选修",
        "source": "personal",
        "is_overridden": false
      }
    ],
    "inherited_count": 15,
    "personal_count": 3
  }
}
```

### 4.6 修改个人课表条目

```
PUT /api/v1/timetables/personal/:id
认证：需要
```

**请求体：**
```json
{
  "course_name": "英语选修（修改）",
  "is_overridden": true
}
```

**成功响应 (200)：**
```json
{
  "code": 0,
  "message": "已更新"
}
```

### 4.7 删除个人课表条目

```
DELETE /api/v1/timetables/personal/:id
认证：需要
```

**成功响应 (200)：**
```json
{
  "code": 0,
  "message": "已删除"
}
```

### 4.8 提交纠错

```
POST /api/v1/timetables/corrections
认证：需要
```

**请求体：**
```json
{
  "class_timetable_id": 5,
  "correction_type": "error",
  "description": "周三第3-4节应该是数据结构，不是算法设计",
  "suggested_course_name": "数据结构",
  "suggested_period_start": 3,
  "suggested_period_end": 4
}
```

**成功响应 (201)：**
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "status": 0,
    "message": "纠错已提交，等待审核"
  }
}
```

### 4.9 获取纠错列表

```
GET /api/v1/timetables/corrections?class_id=1&status=0
认证：需要（仅 owner/admin）
```

**成功响应 (200)：**
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "class_timetable_id": 5,
        "reporter": { "user_id": 3, "nickname": "小红" },
        "correction_type": "error",
        "description": "周三第3-4节应该是数据结构",
        "suggested_course_name": "数据结构",
        "status": 0,
        "created_at": "2025-09-10T14:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 20
  }
}
```

### 4.10 处理纠错

```
PUT /api/v1/timetables/corrections/:id
认证：需要（仅 owner/admin）
```

**请求体：**
```json
{
  "action": "approve"
}
```

> `action` 可选值：`approve`（采纳）、`reject`（驳回）

**成功响应 (200)：**
```json
{
  "code": 0,
  "message": "纠错已采纳，课表已更新"
}
```

---

## 5. 空闲时段计算模块（✅ 引擎已实现，API待实现）

### 5.1 计算共同空闲时段

> **核心接口**：原子引擎对外唯一出口，只输出聚合结果，不暴露个人日程。

```
POST /api/v1/free-slots/calculate
认证：需要
```

**请求体：**
```json
{
  "user_ids": [1, 2, 3, 4, 5],
  "date_range": {
    "start": "2025-09-15",
    "end": "2025-09-21"
  },
  "time_preference": {
    "day_start_hour": 8,
    "day_end_hour": 22,
    "min_duration_minutes": 60
  }
}
```

**成功响应 (200)：**
```json
{
  "code": 0,
  "data": {
    "total_users": 5,
    "slots": [
      {
        "date": "2025-09-18",
        "day_of_week": 4,
        "start_time": "14:00",
        "end_time": "16:00",
        "available_count": 5,
        "total_count": 5,
        "rate": 1.0
      },
      {
        "date": "2025-09-19",
        "day_of_week": 5,
        "start_time": "10:00",
        "end_time": "12:00",
        "available_count": 4,
        "total_count": 5,
        "rate": 0.8
      },
      {
        "date": "2025-09-16",
        "day_of_week": 2,
        "start_time": "14:00",
        "end_time": "16:00",
        "available_count": 3,
        "total_count": 5,
        "rate": 0.6
      }
    ],
    "empty_hint": null
  }
}
```

> 当无人有空时，`slots` 为空数组，`empty_hint` 返回提示文案。

---

## 6. 投票模块（🔲 待实现）

### 6.1 创建投票

```
POST /api/v1/polls
认证：需要
```

**请求体：**
```json
{
  "title": "班级会议时间投票",
  "description": "请大家选择本周班会时间",
  "scope_type": "class",
  "scope_id": 1,
  "deadline": "2025-09-20T23:59:59Z",
  "auto_recommend": true,
  "time_preference": {
    "day_start_hour": 8,
    "day_end_hour": 22,
    "min_duration_minutes": 60,
    "max_recommendations": 5
  }
}
```

> `auto_recommend = true` 时，系统自动调用空闲时段引擎生成推荐选项。

**成功响应 (201)：**
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "title": "班级会议时间投票",
    "status": "draft",
    "options": [
      {
        "id": 1,
        "slot_date": "2025-09-18",
        "slot_start_time": "14:00",
        "slot_end_time": "16:00",
        "is_recommended": true,
        "recommendation_rate": 1.0,
        "sort_order": 1
      },
      {
        "id": 2,
        "slot_date": "2025-09-19",
        "slot_start_time": "10:00",
        "slot_end_time": "12:00",
        "is_recommended": true,
        "recommendation_rate": 0.8,
        "sort_order": 2
      }
    ],
    "recommended_count": 2,
    "created_at": "2025-09-15T10:00:00Z"
  }
}
```

### 6.2 获取投票详情

```
GET /api/v1/polls/:id
认证：需要
```

**成功响应 (200)：**
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "title": "班级会议时间投票",
    "description": "请大家选择本周班会时间",
    "creator": { "user_id": 1, "nickname": "小明" },
    "scope_type": "class",
    "scope_id": 1,
    "status": "open",
    "deadline": "2025-09-20T23:59:59Z",
    "options": [
      {
        "id": 1,
        "slot_date": "2025-09-18",
        "slot_start_time": "14:00",
        "slot_end_time": "16:00",
        "is_recommended": true,
        "recommendation_rate": 1.0,
        "sort_order": 1
      }
    ],
    "total_voters": 30,
    "voted_count": 15,
    "user_vote_status": "not_voted",
    "created_at": "2025-09-15T10:00:00Z"
  }
}
```

### 6.3 编辑投票

```
PUT /api/v1/polls/:id
认证：需要（仅创建者）
限制：status = draft 时才可编辑
```

**请求体：**
```json
{
  "title": "修改后的标题",
  "description": "修改后的说明"
}
```

### 6.4 开启投票

```
POST /api/v1/polls/:id/open
认证：需要（仅创建者）
```

**成功响应 (200)：**
```json
{
  "code": 0,
  "message": "投票已开启",
  "data": {
    "status": "open",
    "deadline": "2025-09-20T23:59:59Z"
  }
}
```

### 6.5 关闭投票

```
POST /api/v1/polls/:id/close
认证：需要（仅创建者）
```

**成功响应 (200)：**
```json
{
  "code": 0,
  "message": "投票已关闭",
  "data": {
    "status": "closed",
    "closed_at": "2025-09-20T18:00:00Z"
  }
}
```

### 6.6 获取我参与的投票列表

```
GET /api/v1/polls?scope_type=class&scope_id=1&status=open&page=1&page_size=20
认证：需要
```

**成功响应 (200)：**
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "id": 1,
        "title": "班级会议时间投票",
        "status": "open",
        "creator": { "user_id": 1, "nickname": "小明" },
        "deadline": "2025-09-20T23:59:59Z",
        "total_voters": 30,
        "voted_count": 15,
        "created_at": "2025-09-15T10:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 20
  }
}
```

### 6.7 获取投票选项（含推荐排序）

```
GET /api/v1/polls/:id/options
认证：需要
```

**成功响应 (200)：**
```json
{
  "code": 0,
  "data": {
    "poll_id": 1,
    "options": [
      {
        "id": 1,
        "slot_date": "2025-09-18",
        "slot_start_time": "14:00",
        "slot_end_time": "16:00",
        "is_recommended": true,
        "recommendation_rate": 1.0,
        "sort_order": 1
      }
    ]
  }
}
```

### 6.8 提交投票

```
POST /api/v1/polls/:id/vote
认证：需要
```

**请求体：**
```json
{
  "votes": [
    { "option_id": 1, "choice": "yes" },
    { "option_id": 2, "choice": "maybe" },
    { "option_id": 3, "choice": "no" }
  ]
}
```

> `choice` 可选值：`yes`（可以）、`no`（不行）、`maybe`（也许可以）

**成功响应 (200)：**
```json
{
  "code": 0,
  "message": "投票成功",
  "data": {
    "voted_options": 3
  }
}
```

### 6.9 获取投票结果

```
GET /api/v1/polls/:id/results
认证：需要
```

**成功响应 (200)：**
```json
{
  "code": 0,
  "data": {
    "poll_id": 1,
    "poll_status": "closed",
    "total_voters": 30,
    "voted_count": 25,
    "participation_overview": 0.8333,
    "results": [
      {
        "option_id": 1,
        "slot_date": "2025-09-18",
        "slot_start_time": "14:00",
        "slot_end_time": "16:00",
        "yes_count": 20,
        "no_count": 2,
        "maybe_count": 3,
        "total_votes": 25,
        "participation_rate": 0.8
      },
      {
        "option_id": 2,
        "slot_date": "2025-09-19",
        "slot_start_time": "10:00",
        "slot_end_time": "12:00",
        "yes_count": 15,
        "no_count": 5,
        "maybe_count": 5,
        "total_votes": 25,
        "participation_rate": 0.6
      }
    ],
    "best_option": {
      "option_id": 1,
      "slot_date": "2025-09-18",
      "slot_start_time": "14:00",
      "slot_end_time": "16:00",
      "yes_count": 20,
      "participation_rate": 0.8
    },
    "final_option": null
  }
}
```

> `final_option` 在创建者确认最终时间后才有值。投票结果**只展示聚合数据**，不展示任何个人投票明细。

### 6.10 确认最终时间

```
POST /api/v1/polls/:id/finalize
认证：需要（仅创建者）
```

**请求体：**
```json
{
  "option_id": 1
}
```

**成功响应 (200)：**
```json
{
  "code": 0,
  "message": "已确认最终时间",
  "data": {
    "poll_status": "finalized",
    "final_option": {
      "option_id": 1,
      "slot_date": "2025-09-18",
      "slot_start_time": "14:00",
      "slot_end_time": "16:00"
    }
  }
}
```

---

## 7. 用户模块（🔲 待实现）

### 7.1 获取用户公开信息

```
GET /api/v1/users/:id/profile
认证：需要
```

**成功响应 (200)：**
```json
{
  "code": 0,
  "data": {
    "user_id": 1,
    "nickname": "小明",
    "avatar_url": "https://...",
    "school": { "id": 1, "name": "XX大学" }
  }
}
```

> 隐私保护：不返回学号、手机号、课表等敏感信息。

### 7.2 修改个人信息

```
PUT /api/v1/users/me
认证：需要
```

**请求体：**
```json
{
  "nickname": "新昵称",
  "avatar_url": "https://..."
}
```

**成功响应 (200)：**
```json
{
  "code": 0,
  "message": "已更新"
}
```

### 7.3 修改隐私设置

```
PUT /api/v1/users/me/privacy
认证：需要
```

**请求体：**
```json
{
  "privacy_level": 2
}
```

> `privacy_level`：1=默认最小暴露（推荐），2=允许同班成员查看部分信息

**成功响应 (200)：**
```json
{
  "code": 0,
  "message": "隐私设置已更新"
}
```

---

## 8. 错误码清单

### 8.1 通用错误码（1000-1099）

| 错误码 | 说明 |
|--------|------|
| 1000 | 未知服务器错误 |
| 1001 | 请求参数错误 |
| 1002 | 资源不存在 |
| 1003 | 请求频率超限 |
| 1004 | 功能暂未开放 |

### 8.2 认证错误码（1100-1199）

| 错误码 | 说明 |
|--------|------|
| 1100 | 未登录或 Token 已过期 |
| 1101 | 学号格式不正确 |
| 1102 | 密码错误 |
| 1103 | 学号已被注册 |
| 1104 | 学号认证失败 |
| 1105 | Token 无效 |
| 1106 | 无操作权限 |

### 8.3 班级错误码（1200-1299）

| 错误码 | 说明 |
|--------|------|
| 1200 | 班级不存在 |
| 1201 | 邀请码错误 |
| 1202 | 已是班级成员 |
| 1203 | 无管理权限 |
| 1204 | 班级代码已存在 |

### 8.4 课表错误码（1300-1399）

| 错误码 | 说明 |
|--------|------|
| 1300 | 课表条目不存在 |
| 1301 | 班级课表已存在（不可重复录入，应使用更新接口） |
| 1302 | 节次范围无效（如 start > end） |
| 1303 | 时间冲突（个人选修课与已有课表重叠） |
| 1304 | 纠错记录不存在 |
| 1305 | 纠错已处理 |

### 8.5 投票错误码（1400-1499）

| 错误码 | 说明 |
|--------|------|
| 1400 | 投票不存在 |
| 1401 | 投票已关闭 |
| 1402 | 已投过票（同一选项） |
| 1403 | 选项不属于该投票 |
| 1404 | 仅创建者可操作 |
| 1405 | 投票尚未开启 |
| 1406 | 已达截止时间 |
| 1407 | 无权限查看 |
