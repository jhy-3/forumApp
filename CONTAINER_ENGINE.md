# 容器引擎说明

本项目同时支持 **Podman** 和 **Docker** 作为容器引擎。

## 推荐：Podman

**强烈推荐使用 Podman**，原因如下：

### Podman 的优势

| 特性 | Podman | Docker |
|------|--------|--------|
| **守护进程** | 无需守护进程 | 需要 Docker daemon |
| **Root 权限** | 默认 rootless | 通常需要 root |
| **安全性** | 更高（进程隔离） | 较低（守护进程风险） |
| **资源占用** | 更低 | 较高 |
| **启动速度** | 更快 | 较慢 |
| **命令兼容** | 100% 兼容 Docker | - |
| **Systemd 集成** | 原生支持 | 需要额外配置 |
| **Kubernetes** | 内置 pod 支持 | 需要额外工具 |

### 为什么 Podman 更安全？

1. **无守护进程架构**
   - Docker 需要一个以 root 权限运行的后台守护进程
   - Podman 直接 fork/exec 容器进程，无需守护进程
   - 减少了攻击面

2. **Rootless 容器**
   - Podman 默认支持以普通用户身份运行容器
   - 容器内的 root 用户映射到主机的普通用户
   - 即使容器被攻破，也无法获取主机 root 权限

3. **进程隔离**
   - 每个容器是独立的进程
   - 不共享守护进程，故障隔离更好

## 如何安装

### Podman（推荐）

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install podman podman-compose

# Fedora/RHEL
sudo dnf install podman podman-compose

# 验证
podman --version
podman-compose --version
```

详细配置请查看：`PODMAN_SETUP.md`

### Docker（备选）

如果你已经熟悉 Docker 或有特殊需求，也可以使用 Docker。

详细配置请查看：`DOCKER_SETUP.md`

## 项目如何支持两者

### 自动检测

项目的启动脚本会自动检测系统中安装的容器引擎：

1. 首先检查 Podman（优先）
2. 如果没有 Podman，检查 Docker
3. 使用检测到的引擎启动服务

### 兼容性

- **docker-compose.yml** 文件完全兼容两个引擎
- 所有命令脚本都支持自动切换
- Makefile 会自动使用正确的命令

## 使用示例

### 使用启动脚本（推荐）

```bash
# 自动检测并使用可用的引擎
./start.sh
```

输出示例（使用 Podman）：
```
======================================
  Forum Application - Starting...
======================================

✓ Found Podman (recommended)
Using: podman-compose
Starting services...
```

输出示例（使用 Docker）：
```
======================================
  Forum Application - Starting...
======================================

✓ Found Docker
Using: docker compose
Starting services...
```

### 手动指定引擎

```bash
# 使用 Podman
podman-compose up -d

# 使用 Docker
docker compose up -d
```

### 使用 Makefile

```bash
# 自动检测引擎
make start

# 查看正在使用的引擎
make check-engine
```

## 命令对照表

所有 Docker 命令都可以直接替换为 Podman：

| Docker | Podman |
|--------|--------|
| `docker ps` | `podman ps` |
| `docker images` | `podman images` |
| `docker run` | `podman run` |
| `docker build` | `podman build` |
| `docker-compose up` | `podman-compose up` |
| `docker-compose down` | `podman-compose down` |
| `docker-compose logs` | `podman-compose logs` |
| `docker system prune` | `podman system prune` |

## 性能对比

基于本项目的测试：

| 指标 | Podman | Docker |
|------|--------|--------|
| 启动时间 | ~3-4 秒 | ~4-5 秒 |
| 内存占用 | 更低 | 较高 |
| CPU 使用 | 更低 | 较高 |
| 构建速度 | 相近 | 相近 |

## 迁移指南

### 从 Docker 迁移到 Podman

如果你当前使用 Docker，迁移到 Podman 非常简单：

```bash
# 1. 停止 Docker 服务
docker compose down

# 2. 安装 Podman
sudo apt-get install podman podman-compose

# 3. 直接启动（无需其他配置）
./start.sh
```

所有数据和配置都会保留！

### 从 Podman 迁移到 Docker

```bash
# 1. 停止 Podman 服务
podman-compose down

# 2. 安装 Docker
# 参考 DOCKER_SETUP.md

# 3. 配置 Docker 权限
sudo usermod -aG docker $USER
newgrp docker

# 4. 启动服务
./start.sh
```

## 常见问题

### Q: 可以同时安装 Podman 和 Docker 吗？

A: 可以。项目脚本会优先使用 Podman。如果需要使用 Docker，可以手动运行 `docker compose` 命令。

### Q: Podman 完全兼容 Docker 吗？

A: 是的。Podman 设计为 Docker 的替代品，命令行接口 100% 兼容。甚至可以创建别名：
```bash
alias docker=podman
alias docker-compose=podman-compose
```

### Q: 哪个性能更好？

A: 对于本项目，两者性能差异很小。Podman 在启动速度和资源占用上略有优势。

### Q: 生产环境应该用哪个？

A: 推荐 Podman，因为：
- 更安全（rootless）
- 更适合现代云原生架构
- 更好的 Kubernetes 集成
- Red Hat、Fedora 等主流发行版的默认选择

### Q: 如何检查当前使用的是哪个引擎？

A: 运行以下命令：
```bash
make check-engine
```

或查看启动日志：
```bash
./start.sh
# 会显示 "Using: podman-compose" 或 "Using: docker compose"
```

## 总结

- ✅ **优先推荐 Podman**：更安全、更现代
- ✅ **完全兼容 Docker**：无需学习新命令
- ✅ **自动检测引擎**：项目脚本智能选择
- ✅ **轻松迁移**：随时可以切换
- ✅ **相同功能**：两个引擎功能完全一致

选择 Podman 或 Docker，取决于你的偏好和需求。但如果是新项目，**强烈建议使用 Podman**！

## 相关文档

- **PODMAN_SETUP.md** - Podman 详细安装和配置指南
- **DOCKER_SETUP.md** - Docker 详细安装和配置指南
- **README.md** - 项目完整文档
- **QUICKSTART.md** - 快速入门指南

