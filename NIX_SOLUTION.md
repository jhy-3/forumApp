# 使用 Nix 安装 Podman（无需 Sudo）

## 🎯 完美的解决方案！

**Nix 包管理器**允许你在**完全无需 root 权限**的情况下安装**完整版 Podman**！

## 什么是 Nix？

Nix 是一个强大的包管理器，特点：
- ✅ **无需 root** - 完全在用户空间运行
- ✅ **完整软件包** - 安装完整版软件，不是阉割版
- ✅ **可重复构建** - 确保一致性
- ✅ **隔离环境** - 不污染系统
- ✅ **可回滚** - 随时恢复到之前版本

## 🚀 快速安装（一条命令）

```bash
./INSTALL_PODMAN_WITH_NIX.sh
```

脚本会自动：
1. 安装 Nix 包管理器（无需 sudo）
2. 通过 Nix 安装完整版 Podman
3. 安装 podman-compose
4. 测试是否工作
5. 提供下一步指引

## 📋 详细步骤

### 步骤 1：安装 Nix

```bash
# 单用户模式安装 Nix（无需 sudo）
curl -L https://nixos.org/nix/install | sh -s -- --no-daemon

# 加载 Nix 环境
source ~/.nix-profile/etc/profile.d/nix.sh
```

### 步骤 2：安装 Podman

```bash
# 通过 Nix 安装完整版 Podman
nix-env -iA nixpkgs.podman

# 验证
podman --version
```

### 步骤 3：配置环境

```bash
# 添加 Nix 到 .bashrc（可选但推荐）
echo 'if [ -e ~/.nix-profile/etc/profile.d/nix.sh ]; then . ~/.nix-profile/etc/profile.d/nix.sh; fi' >> ~/.bashrc

# 立即加载
source ~/.bashrc
```

### 步骤 4：安装 podman-compose

```bash
# 方式 1：使用 pip（推荐）
pip3 install --user podman-compose

# 方式 2：使用 Nix
nix-env -iA nixpkgs.podman-compose
```

### 步骤 5：启动应用

```bash
# 正常启动脚本现在可以工作了！
./start.sh
```

## 为什么 Nix 方案可行？

### 对比其他方案

| 方案 | Podman 版本 | 需要 sudo | 可行性 |
|------|-------------|-----------|--------|
| 静态二进制 | podman-remote | ❌ | ❌ 不可用 |
| 包管理器 | 完整版 | ✅ | ❌ 无权限 |
| **Nix** | **完整版** | **❌** | **✅ 可用** |

### Nix 的工作原理

```
传统包管理器：
sudo apt install → /usr/bin/ (需要 root)

Nix 包管理器：
nix-env -i → ~/.nix-profile/bin/ (用户空间)
             ↓
           /nix/store/hash-package/
```

- 所有软件安装在 `/nix/store/`
- 用户 profile 链接到需要的版本
- 完全隔离，互不影响

## 🔍 技术细节

### Nix 安装位置

```bash
/nix/store/          # 所有软件包（immutable）
~/.nix-profile/      # 用户 profile
~/.nix-profile/bin/  # 可执行文件链接
```

### Podman 配置

Nix 安装的 Podman 会自动配置：
- ✅ 完整的二进制文件
- ✅ 所有必需的工具
- ✅ Rootless 配置
- ✅ 所有依赖项

## 💡 优势

### 相比静态二进制

| 特性 | 静态二进制 | Nix 安装 |
|------|-----------|----------|
| 版本 | podman-remote | 完整版 ✅ |
| system service | ❌ | ✅ |
| 独立运行 | ❌ | ✅ |
| 所有命令 | 部分 | 全部 ✅ |
| 依赖处理 | 手动 | 自动 ✅ |

### 相比请求管理员

- ✅ 立即可用 - 不需要等待审批
- ✅ 完全控制 - 自己管理版本
- ✅ 隔离安装 - 不影响系统
- ❌ 占用空间 - Nix 需要 ~500MB-1GB

## 📖 使用流程

### 完整流程

```bash
# 1. 安装 Nix 和 Podman（一次性）
./INSTALL_PODMAN_WITH_NIX.sh

# 2. 加载环境
source ~/.nix-profile/etc/profile.d/nix.sh

# 3. 配置 Podman rootless
./FIX_PODMAN_ROOTLESS.sh

# 4. 启动应用
./start.sh

# 5. 访问应用
# http://localhost:3000
```

### 日常使用

以后每次使用，只需：

```bash
# 在新终端中
source ~/.nix-profile/etc/profile.d/nix.sh

# 启动应用
./start.sh
```

或者将 Nix 加载命令添加到 ~/.bashrc：

```bash
echo 'source ~/.nix-profile/etc/profile.d/nix.sh' >> ~/.bashrc
```

## ⚙️ 配置 Podman

安装完成后，配置 Podman rootless 模式：

```bash
# 创建配置目录
mkdir -p ~/.config/containers

# 配置 storage
cat > ~/.config/containers/storage.conf << 'EOF'
[storage]
driver = "overlay"
graphroot = "$HOME/.local/share/containers/storage"
runroot = "/run/user/$UID/containers"
EOF

# 配置 registries
cat > ~/.config/containers/registries.conf << 'EOF'
unqualified-search-registries = ["docker.io"]

[[registry]]
location = "docker.io"
EOF
```

## 🧪 测试

```bash
# 测试 Podman
podman run --rm alpine echo "Hello from Podman!"

# 测试网络
podman run --rm alpine ping -c 3 google.com

# 查看信息
podman info

# 查看版本
podman version
```

## 🆘 故障排查

### Nix 安装失败

```bash
# 检查安装日志
cat ~/.nix-profile/etc/profile.d/nix.sh

# 手动重试
curl -L https://nixos.org/nix/install | sh -s -- --no-daemon
```

### Podman 找不到

```bash
# 加载 Nix 环境
source ~/.nix-profile/etc/profile.d/nix.sh

# 检查
which podman
podman --version
```

### 容器运行失败

```bash
# 检查 rootless 配置
podman info | grep -A 10 "graphRoot"

# 重置存储
rm -rf ~/.local/share/containers/storage
podman system reset
```

## 📊 资源占用

### 磁盘空间

```bash
# Nix 本身: ~500MB
# Podman: ~100MB
# 容器镜像: 根据使用情况

# 总计建议: 至少 2GB 可用空间
```

### 检查空间

```bash
df -h ~
du -sh ~/.nix-profile
du -sh ~/.local/share/containers
```

## 🎯 完整解决方案

### 方案 A：使用 Nix（推荐，无需 sudo）

```bash
# 1. 安装 Nix 和 Podman
./INSTALL_PODMAN_WITH_NIX.sh

# 2. 加载环境
source ~/.nix-profile/etc/profile.d/nix.sh

# 3. 测试
podman run hello-world

# 4. 启动应用
./start.sh
```

**优点：**
- ✅ 完整的 Podman
- ✅ 无需 sudo
- ✅ 立即可用

**缺点：**
- ⚠️ 需要下载 ~500MB
- ⚠️ 需要 1-2GB 磁盘空间

### 方案 B：请求管理员安装（最简单）

提交安装请求：
- 查看：`给管理员的安装请求.md`

**优点：**
- ✅ 系统级安装，所有用户受益
- ✅ 占用空间小
- ✅ 标准配置

### 方案 C：使用其他环境

- 个人电脑/虚拟机
- GitHub Codespaces
- 云服务器

## 🔄 Nix vs 包管理器

| 特性 | Nix | apt/dnf |
|------|-----|---------|
| 需要 root | ❌ | ✅ |
| 安装位置 | 用户空间 | 系统目录 |
| 软件完整性 | ✅ 完整 | ✅ 完整 |
| 磁盘占用 | 较大 | 较小 |
| 版本管理 | 优秀 | 一般 |
| 隔离性 | 优秀 | 一般 |

## 📚 Nix 简介

### 什么是 Nix？

Nix 是一个跨平台的包管理器：
- 🌟 **声明式配置** - 可重复的环境
- 🌟 **原子升级** - 要么全部成功，要么全部回滚
- 🌟 **多版本共存** - 可以同时安装多个版本
- 🌟 **纯函数式** - 保证构建一致性

### Nix 的用途

```bash
# 安装软件
nix-env -iA nixpkgs.package-name

# 搜索软件
nix search nixpkgs package-name

# 列出已安装
nix-env -q

# 删除软件
nix-env -e package-name

# 回滚到上一代
nix-env --rollback
```

### 学习资源

- 官网：https://nixos.org/
- 文档：https://nixos.org/manual/nix/stable/
- 搜索包：https://search.nixos.org/

## 🎉 总结

### 最佳方案：Nix + Podman

**这是在没有 sudo 权限时的最佳方案！**

**安装命令：**
```bash
./INSTALL_PODMAN_WITH_NIX.sh
```

**为什么选择这个方案：**
1. ✅ 完全无需 sudo
2. ✅ 获得完整版 Podman
3. ✅ 可以独立运行容器
4. ✅ 所有功能都可用
5. ✅ 安装后立即可用

**权衡：**
- ⚠️ 需要下载较大的安装包
- ⚠️ 需要 1-2GB 磁盘空间
- ⚠️ 首次安装需要10-15分钟

**但是：一次安装，永久使用！**

---

**准备好了吗？运行：**

```bash
./INSTALL_PODMAN_WITH_NIX.sh
```

然后跟随屏幕提示完成安装！🚀

