# Podman 安装和配置指南

## 为什么选择 Podman？

Podman 是 Docker 的现代化替代方案，具有以下优势：

- ✅ **无守护进程** - 更安全，不需要后台服务
- ✅ **无需 root 权限** - 支持 rootless 模式
- ✅ **兼容 Docker** - 命令完全兼容
- ✅ **更安全** - 用户隔离更好
- ✅ **系统集成** - 更好的 systemd 集成

## 安装 Podman

### Ubuntu/Debian

```bash
# 更新包列表
sudo apt-get update

# 安装 Podman
sudo apt-get -y install podman

# 安装 podman-compose
sudo apt-get -y install podman-compose
# 或使用 pip
pip3 install podman-compose
```

### Fedora/RHEL/CentOS

```bash
# Podman 通常已预装
sudo dnf install podman podman-compose
```

### Arch Linux

```bash
sudo pacman -S podman podman-compose
```

### 验证安装

```bash
# 检查 Podman 版本
podman --version

# 检查 podman-compose 版本
podman-compose --version

# 测试运行
podman run hello-world
```

## Rootless 配置（推荐）

Podman 的最大优势是支持 rootless 模式，无需 root 权限：

```bash
# 配置 subuid 和 subgid（通常已自动配置）
grep $USER /etc/subuid
grep $USER /etc/subgid

# 如果没有配置，手动添加
echo "$USER:100000:65536" | sudo tee -a /etc/subuid
echo "$USER:100000:65536" | sudo tee -a /etc/subgid

# 启用 systemd 用户服务
systemctl --user enable --now podman.socket
```

## 端口映射配置

对于低于 1024 的端口，需要配置：

```bash
# 允许非特权用户绑定低端口（可选）
echo "net.ipv4.ip_unprivileged_port_start=80" | sudo tee -a /etc/sysctl.conf
sudo sysctl -p
```

## Docker 兼容性

如果你习惯使用 Docker 命令，可以创建别名：

```bash
# 添加到 ~/.bashrc 或 ~/.zshrc
echo 'alias docker=podman' >> ~/.bashrc
echo 'alias docker-compose=podman-compose' >> ~/.bashrc
source ~/.bashrc
```

## 网络配置

Podman 使用 CNI 网络插件：

```bash
# 检查网络
podman network ls

# 创建网络（通常自动创建）
podman network create forum-network
```

## 故障排查

### 权限问题

```bash
# 检查当前用户的 subuid/subgid
podman system info | grep -A2 IDMappings

# 重置 Podman
podman system reset
```

### 端口访问问题

```bash
# 检查 SELinux（如果启用）
sudo setenforce 0  # 临时禁用测试
getenforce

# 检查防火墙
sudo firewall-cmd --list-all
```

### 网络问题

```bash
# 重新创建网络
podman network rm forum-network
podman network create forum-network

# 检查 DNS
podman run --rm alpine ping -c 3 google.com
```

## 从 Docker 迁移

### 1. 导出 Docker 镜像

```bash
# 如果有 Docker 镜像
docker save -o images.tar image1 image2
podman load -i images.tar
```

### 2. 迁移卷

```bash
# Docker 卷路径通常在
sudo ls /var/lib/docker/volumes/

# Podman 卷路径在
podman volume ls
```

### 3. 清理旧的 Docker（可选）

```bash
# 停止 Docker 服务
sudo systemctl stop docker
sudo systemctl disable docker

# 卸载 Docker（谨慎操作）
# sudo apt-get remove docker docker-engine docker.io containerd runc
```

## 与 Docker 的主要区别

| 特性 | Docker | Podman |
|------|--------|--------|
| 守护进程 | 需要 | 不需要 |
| Root 权限 | 通常需要 | 可选 |
| 架构 | Client-Server | Fork-Exec |
| Systemd 集成 | 差 | 好 |
| Kubernetes 集成 | 需要额外工具 | 内置支持 |
| 命令兼容性 | - | 100% 兼容 |

## 性能优化

```bash
# 启用 cgroup v2（推荐）
# 检查当前版本
stat -fc %T /sys/fs/cgroup/

# 配置存储驱动
podman info --format '{{.Store.GraphDriverName}}'

# 推荐使用 overlay
```

## 常用命令对比

```bash
# Docker -> Podman
docker run         -> podman run
docker ps          -> podman ps
docker images      -> podman images
docker build       -> podman build
docker-compose up  -> podman-compose up
docker system prune -> podman system prune
```

## 验证配置

运行测试以确保 Podman 正常工作：

```bash
# 1. 测试基本功能
podman run --rm alpine echo "Podman works!"

# 2. 测试网络
podman run --rm --net=host alpine ping -c 3 google.com

# 3. 测试卷挂载
podman run --rm -v test-vol:/data alpine touch /data/test
podman volume ls

# 4. 测试端口映射
podman run --rm -d -p 8080:80 nginx
curl http://localhost:8080
podman stop $(podman ps -q)

# 5. 清理测试
podman volume rm test-vol
```

## 推荐配置

创建 `~/.config/containers/containers.conf`：

```toml
[containers]
# 网络后端
network_backend = "cni"

# 默认网络
default_network = "forum-network"

[engine]
# 并发下载层数
max_parallel_downloads = 5

# 事件日志
events_logger = "journald"
```

## 下一步

配置完成后，使用项目提供的脚本启动应用：

```bash
cd /home/hy/develop/forumApp
./start.sh
```

脚本会自动检测并使用 Podman！

## 参考资源

- [Podman 官方文档](https://docs.podman.io/)
- [从 Docker 迁移到 Podman](https://podman.io/getting-started/installation)
- [Rootless Podman 教程](https://github.com/containers/podman/blob/main/docs/tutorials/rootless_tutorial.md)

