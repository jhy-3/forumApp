# 快速开始指南

## 推荐：使用 Podman（更安全，无需 root）

### 安装 Podman

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install podman podman-compose

# Fedora/RHEL
sudo dnf install podman podman-compose

# 验证安装
podman --version
podman-compose --version
```

### 启动应用

使用提供的启动脚本，一条命令启动所有服务：

```bash
./start.sh
```

脚本会自动检测并使用 Podman（优先）或 Docker。

**Podman 的优势：**
- ✅ 无需守护进程，更安全
- ✅ 无需 root 权限
- ✅ 完全兼容 Docker 命令
- ✅ 更好的资源隔离

启动脚本会：
1. 检查 Docker 和 Docker Compose 是否安装
2. 启动所有服务（数据库、后端、前端）
3. 等待服务就绪
4. 显示访问地址

服务启动后访问：
- **前端界面**: http://localhost:3000
- **后端 API**: http://localhost:8080

## 停止服务

```bash
./stop.sh
```

或者：

```bash
podman-compose down  # 使用 Podman
# 或
docker compose down  # 使用 Docker
```

## 完全清理（包括数据库数据）

```bash
podman-compose down -v  # 使用 Podman
# 或
docker compose down -v  # 使用 Docker
```

## 使用 Makefile（可选）

项目提供了 Makefile 简化操作，会自动检测使用 Podman 或 Docker：

```bash
# 查看所有可用命令
make help

# 检查使用的容器引擎
make check-engine

# 启动所有服务
make start

# 停止服务
make stop

# 查看日志
make logs

# 测试后端健康检查
make test-api

# 测试用户注册
make test-register

# 打开数据库命令行
make db-shell
```

## 第一次使用

1. **启动服务**
   ```bash
   ./start.sh
   ```

2. **打开浏览器访问** http://localhost:3000

3. **注册一个新账号**
   - 点击右上角 "Register"
   - 填写用户名、邮箱和密码（至少6位）
   - 点击 "Register" 按钮

4. **浏览论坛**
   - 登录后可以看到默认的4个板块
   - 点击任意板块进入

5. **创建主题**
   - 进入板块后，点击 "New Thread" 按钮
   - 填写标题和内容
   - 点击 "Create Thread" 提交

6. **发表回复**
   - 点击进入任意主题
   - 在底部的回复框输入内容
   - 点击 "Post Reply" 发送

## 备选：使用 Docker

如果你更熟悉 Docker，项目也完全支持：

```bash
# 安装 Docker（如果未安装）
# 参考 DOCKER_SETUP.md

# 启动应用（脚本会自动检测）
./start.sh
```

## 本地开发模式

如果你想修改代码并实时看到效果：

### 1. 只启动数据库

```bash
make dev-db
# 或
podman-compose -f docker-compose.dev.yml up -d
# 或
docker compose -f docker-compose.dev.yml up -d
```

### 2. 本地运行后端（终端1）

```bash
cd backend
go mod download  # 第一次需要
go run main.go
```

### 3. 本地运行前端（终端2）

```bash
cd frontend
npm install  # 第一次需要
npm start
```

前端会自动打开浏览器并启用热重载。

## 故障排查

### 端口被占用

如果提示端口被占用，可以：

1. 停止占用端口的程序
2. 或修改 `docker-compose.yml` 中的端口映射

### Podman 相关问题

```bash
# 检查 Podman 状态
podman info

# 重置 Podman（谨慎使用）
podman system reset

# 查看详细日志
podman-compose logs -f
```

详细的 Podman 配置请参考 `PODMAN_SETUP.md`

### Docker 权限问题

如果使用 Docker 遇到权限问题，请参考 `DOCKER_SETUP.md`

### 数据库初始化失败

重新创建数据库：

```bash
podman-compose down -v  # 或 docker compose down -v
./start.sh
```

### 前端无法连接后端

检查后端是否正常运行：

```bash
curl http://localhost:8080/health
```

应该返回：
```json
{"status":"ok"}
```

## 系统要求

**推荐配置：**
- Podman 4.0+
- Podman Compose 1.0+

**或备选配置：**
- Docker 20.10+
- Docker Compose 2.0+

**本地开发（可选）：**
- Go 1.21+
- Node.js 18+

## 默认账号信息

数据库初始不包含任何用户，需要自己注册。

默认创建了4个板块：
- General Discussion (general)
- Announcements (announcements)
- Technical Support (support)
- Off Topic (off-topic)

## 下一步

- 查看 `README.md` 了解完整功能
- 查看 `PODMAN_SETUP.md` 了解 Podman 详细配置
- 查看 `guidelines/guidv1.0.md` 了解架构设计
- 查看 API 文档部分了解如何使用 API

## 技术支持

遇到问题？

1. 查看日志: `podman-compose logs -f` 或 `docker compose logs -f`
2. 查看后端日志: `make logs` 或 `podman-compose logs backend`
3. 查看数据库日志: `podman-compose logs database`
4. 查看 `PODMAN_SETUP.md` 或 `DOCKER_SETUP.md` 的故障排查章节
5. 查看 README.md 的故障排查章节

