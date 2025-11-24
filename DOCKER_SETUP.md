# Docker 权限配置指南

## 问题说明

如果你遇到以下错误：
```
permission denied while trying to connect to the Docker daemon socket
```

这是因为当前用户没有访问 Docker 的权限。

## 解决方案

### 方案 1：将用户添加到 docker 组（推荐）

这是最推荐的方式，配置后就不需要每次使用 sudo。

```bash
# 1. 将当前用户添加到 docker 组
sudo usermod -aG docker $USER

# 2. 重新加载用户组（二选一）
# 方式A：重新登录系统
# 方式B：运行以下命令
newgrp docker

# 3. 验证配置
docker run hello-world
```

配置完成后，就可以直接运行：
```bash
./start.sh
```

### 方案 2：使用 sudo 运行（临时方案）

如果你不想修改用户组，可以使用 sudo：

```bash
# 使用 sudo 启动
sudo docker compose up -d

# 或使用 make
sudo make start

# 查看日志
sudo docker compose logs -f

# 停止服务
sudo docker compose down
```

## 验证 Docker 权限

运行以下命令检查 Docker 是否可用：

```bash
# 检查是否在 docker 组中
groups | grep docker

# 尝试运行 Docker 命令
docker ps

# 检查 Docker 版本
docker --version
docker compose version
```

## 常见问题

### Q: 为什么需要将用户添加到 docker 组？

A: Docker 守护进程运行时需要 root 权限。将用户添加到 docker 组后，该用户就可以通过 Docker socket 与守护进程通信，而无需每次都使用 sudo。

### Q: 添加到 docker 组后还是提示权限错误？

A: 需要完全退出并重新登录系统，或者运行 `newgrp docker` 命令使组更改生效。

### Q: 使用 sudo 有什么缺点？

A: 
- 每次都需要输入密码
- 创建的容器和卷的所有者是 root
- 某些开发工具可能无法正常访问 Docker

### Q: 安全性如何？

A: docker 组的成员实际上拥有了 root 等效的权限。因此：
- 只在受信任的环境中使用
- 只将受信任的用户添加到 docker 组
- 对于生产服务器，考虑使用 rootless Docker

## 推荐配置步骤

完整的配置步骤：

```bash
# 1. 添加用户到 docker 组
sudo usermod -aG docker $USER

# 2. 重新登录或运行
newgrp docker

# 3. 测试 Docker
docker run hello-world

# 4. 清理测试容器
docker rm $(docker ps -aq --filter "ancestor=hello-world")

# 5. 启动论坛应用
cd /home/hy/develop/forumApp
./start.sh
```

## 其他资源

- [Docker 官方文档 - 以非 root 用户身份管理 Docker](https://docs.docker.com/engine/install/linux-postinstall/)
- [Rootless Docker](https://docs.docker.com/engine/security/rootless/)

