# Forum Application - Phase 1 MVP

基于文档 `guidelines/guidv1.0.md` 的第一阶段实现，这是一个现代化的论坛应用程序，采用高性能技术栈构建。

## 🎉 现已支持 Podman！

项目已更新为**优先使用 Podman**（更安全，无需 root），同时完全兼容 Docker。

### 快速开始（推荐 Podman）

```bash
# 1. 安装 Podman（一行命令，自动检测是否需要 sudo）
./INSTALL_PODMAN.sh

# 2. 重新加载配置（如果无 sudo 安装）
source ~/.bashrc

# 3. 启动应用（自动检测容器引擎）
./start.sh

# 4. 访问应用
# 前端: http://localhost:3000
# 后端: http://localhost:8080
```

**Podman 的优势**：🔒 更安全 | ⚡ 更快 | 💾 更轻 | 🔑 无需 root | ✨ 无需 sudo

详细说明请查看：
- **使用说明.md** - 中文快速指南
- **INSTALL_WITHOUT_SUDO.md** - 无 sudo 安装指南
- **QUICKSTART.md** - 快速入门
- **PODMAN_SETUP.md** - Podman 详细配置

---

## 技术栈

### 后端
- **Go 1.21** - 高性能后端语言
- **Echo Framework** - 轻量级 Web 框架
- **MariaDB** - 关系型数据库
- **JWT** - 无状态身份认证

### 前端
- **React 18** - UI 框架
- **TypeScript** - 类型安全
- **Material-UI (MUI)** - UI 组件库
- **React Router** - 路由管理
- **Axios** - HTTP 客户端

### 部署
- **Podman & Podman Compose** - 容器化部署（推荐，更安全）
- **Docker & Docker Compose** - 容器化部署（也支持）

## 核心功能

第一阶段 MVP 版本包含以下功能：

### 用户系统
- ✅ 用户注册
- ✅ 用户登录
- ✅ JWT 会话管理
- ✅ 个人资料查看和编辑

### 内容发布
- ✅ 板块（Forum）浏览
- ✅ 主题（Thread）创建和查看
- ✅ 帖子（Post）创建、编辑和删除
- ✅ 回复功能

### 数据库设计
- ✅ 规范化的表结构
- ✅ 外键约束和索引优化
- ✅ 帖子内容分离存储（性能优化）

## 快速开始

### 前置要求

- **Podman 和 Podman Compose**（推荐）或 Docker 和 Docker Compose
- （可选）Go 1.21+ 和 Node.js 18+ 用于本地开发

### 方法 1: 使用容器（推荐）

项目同时支持 Podman 和 Docker，**强烈推荐使用 Podman**，因为：
- ✅ 更安全（无守护进程）
- ✅ 无需 root 权限
- ✅ 更好的资源隔离

#### 安装 Podman（推荐）

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install podman podman-compose

# Fedora/RHEL (通常已预装)
sudo dnf install podman podman-compose

# 验证安装
podman --version
podman-compose --version
```

详细说明请查看 `PODMAN_SETUP.md`

#### 或安装 Docker（备选）

如果需要使用 Docker，请参考 `DOCKER_SETUP.md`

#### 启动服务

```bash
# 方式1：使用便捷脚本（推荐，自动检测 Podman/Docker）
./start.sh

# 方式2：直接使用 compose 命令
podman-compose up -d    # 使用 Podman
# 或
docker compose up -d    # 使用 Docker

# 查看日志
podman-compose logs -f  # 或 docker compose logs -f

# 停止所有服务
./stop.sh
# 或
podman-compose down     # 或 docker compose down
```

服务将在以下端口启动：
- 前端: http://localhost:3000
- 后端 API: http://localhost:8080
- 数据库: localhost:3306

### 方法 2: 本地开发模式

适合开发调试，只启动数据库，后端和前端在本地运行。

#### 1. 启动数据库

```bash
# 使用开发配置启动数据库
podman-compose -f docker-compose.dev.yml up -d
# 或
docker compose -f docker-compose.dev.yml up -d

# 或使用 make
make dev-db
```

#### 2. 启动后端

```bash
cd backend

# 安装依赖
go mod download

# 复制环境变量文件（如果需要修改配置）
# cp .env.example .env

# 运行后端
go run main.go
```

后端将在 http://localhost:8080 启动

#### 3. 启动前端

```bash
cd frontend

# 安装依赖
npm install

# 运行开发服务器
npm start
```

前端将在 http://localhost:3000 启动，并自动打开浏览器。

## 项目结构

```
forumApp/
├── backend/                    # Go 后端
│   ├── internal/
│   │   ├── api/               # HTTP 处理器和路由
│   │   ├── config/            # 配置管理
│   │   ├── models/            # 数据模型
│   │   ├── repository/        # 数据访问层
│   │   └── service/           # 业务逻辑层
│   ├── main.go                # 入口文件
│   ├── go.mod                 # Go 依赖
│   └── Dockerfile             # 后端镜像
├── frontend/                   # React 前端
│   ├── src/
│   │   ├── api/               # API 调用
│   │   ├── components/        # React 组件
│   │   ├── contexts/          # Context 状态管理
│   │   ├── pages/             # 页面组件
│   │   └── types/             # TypeScript 类型
│   ├── package.json           # 前端依赖
│   └── Dockerfile             # 前端镜像
├── database/
│   └── init.sql               # 数据库初始化脚本
├── guidelines/
│   └── guidv1.0.md           # 项目架构蓝图
├── docker-compose.yml         # 生产环境配置
├── docker-compose.dev.yml     # 开发环境配置
└── README.md                  # 本文件
```

## API 端点

### 公共端点

- `POST /api/users/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `GET /api/forums` - 获取所有板块
- `GET /api/forums/:slug/threads` - 获取板块的主题列表
- `GET /api/threads/:id` - 获取主题详情
- `GET /api/threads/:id/posts` - 获取主题的帖子列表

### 需要认证的端点

需要在请求头中包含 `Authorization: Bearer <token>`

- `GET /api/users/me` - 获取当前用户信息
- `PATCH /api/users/me` - 更新当前用户信息
- `POST /api/threads` - 创建新主题
- `POST /api/threads/:id/posts` - 在主题中发表回复
- `PATCH /api/posts/:id` - 编辑帖子
- `DELETE /api/posts/:id` - 删除帖子

## 数据库架构

核心表结构：

- `users` - 用户信息
- `forums` - 板块信息
- `threads` - 主题
- `posts` - 帖子元数据
- `post_contents` - 帖子内容（分离存储以优化性能）

详细的数据库设计请参考 `database/init.sql`。

## 配置说明

### 后端配置

后端配置通过环境变量设置。主要配置项：

```bash
# 服务器配置
PORT=8080
ENVIRONMENT=development

# 数据库配置
DB_HOST=database
DB_PORT=3306
DB_USER=forumuser
DB_PASSWORD=forumpass123
DB_NAME=forumdb

# JWT 配置
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_ACCESS_TOKEN_EXPIRE_MINUTES=15
JWT_REFRESH_TOKEN_EXPIRE_DAYS=7

# CORS 配置
CORS_ALLOWED_ORIGINS=http://localhost:3000
```

### 前端配置

前端配置在 `.env` 文件中：

```bash
REACT_APP_API_URL=http://localhost:8080/api
```

## 开发指南

### 后端开发

1. 遵循 Go 的代码规范
2. 使用分层架构（Handler -> Service -> Repository）
3. 所有错误使用统一的错误响应格式
4. 使用 JWT 进行认证授权

### 前端开发

1. 使用 TypeScript 编写所有代码
2. 遵循 React Hooks 最佳实践
3. 使用 MUI 组件保持 UI 一致性
4. 通过 Context API 管理全局状态

### 添加新功能

1. 后端：在 `internal/models/` 定义模型
2. 后端：在 `internal/repository/` 实现数据访问
3. 后端：在 `internal/service/` 实现业务逻辑
4. 后端：在 `internal/api/` 添加 HTTP 处理器
5. 前端：在 `src/types/` 定义 TypeScript 类型
6. 前端：在 `src/api/` 添加 API 调用
7. 前端：在 `src/pages/` 或 `src/components/` 实现 UI

## 测试

```bash
# 测试后端 API
curl http://localhost:8080/health

# 注册用户
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'

# 登录
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

## 故障排查

### 数据库连接失败

```bash
# 检查数据库是否运行
docker-compose ps database

# 查看数据库日志
docker-compose logs database

# 重启数据库
docker-compose restart database
```

### 后端无法启动

```bash
# 查看后端日志
docker-compose logs backend

# 确保数据库已就绪
docker-compose logs database | grep "ready for connections"
```

### 前端无法连接后端

1. 确认后端正在运行: `curl http://localhost:8080/health`
2. 检查前端的 `.env` 文件中的 `REACT_APP_API_URL`
3. 检查浏览器控制台的 CORS 错误

## 下一步计划

参考 `guidelines/guidv1.0.md` 文档，后续阶段将实现：

### 第二阶段：增强内容与交互功能
- 集成 Vditor Markdown 编辑器
- 用户私信功能
- 实时通知系统
- Elasticsearch 全文搜索

### 第三阶段：社交化与管理功能
- @提及功能
- 点赞/反应系统
- 用户声望和积分
- 管理员后台

### 第四阶段：生产环境加固
- 性能优化和缓存
- 备份和恢复策略
- 监控和告警
- CI/CD 流水线

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

## 联系方式

如有问题，请查阅 `guidelines/guidv1.0.md` 文档或提交 Issue。
