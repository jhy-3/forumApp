# 无 Sudo 权限安装 Podman

如果你没有 sudo 权限，也可以安装和使用 Podman！

## 方式 1：使用安装脚本（推荐）

我们的安装脚本会自动检测是否有 sudo 权限，并选择合适的安装方式：

```bash
./INSTALL_PODMAN.sh
```

脚本会：
1. 检查是否已安装 Podman
2. 检测是否有 sudo 权限
3. 如果没有 sudo，自动使用用户空间安装方式
4. 安装到 `~/.local/bin/`
5. 自动配置 PATH

## 方式 2：手动安装（静态二进制）

### 步骤 1：下载 Podman 静态二进制

```bash
# 创建本地 bin 目录
mkdir -p ~/.local/bin
cd /tmp

# 下载静态二进制（x86_64）
curl -L -o podman.tar.gz \
  https://github.com/containers/podman/releases/download/v4.9.3/podman-remote-static-linux_amd64.tar.gz

# 解压
tar -xzf podman.tar.gz

# 安装
mv podman-remote-static ~/.local/bin/podman
chmod +x ~/.local/bin/podman

# 清理
rm -f podman.tar.gz
```

### 步骤 2：配置 PATH

```bash
# 添加到 .bashrc
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc

# 立即生效
source ~/.bashrc

# 验证
podman --version
```

### 步骤 3：安装 podman-compose

```bash
# 使用 pip3（无需 sudo）
pip3 install --user podman-compose

# 或下载脚本
curl -o ~/.local/bin/podman-compose \
  https://raw.githubusercontent.com/containers/podman-compose/main/podman_compose.py
chmod +x ~/.local/bin/podman-compose
```

## 方式 3：请求系统管理员安装

如果你的系统由管理员管理，可以请求他们安装：

```bash
# 系统管理员运行（需要 root 权限）
sudo apt-get update
sudo apt-get install podman podman-compose
```

安装后，所有用户都可以使用 Podman（无需 sudo）。

## 验证安装

```bash
# 检查 Podman
podman --version

# 检查 podman-compose
podman-compose --version

# 测试运行
podman run --rm alpine echo "Hello from Podman!"
```

## 启动应用

安装完成后，直接启动：

```bash
./start.sh
```

脚本会自动检测并使用 Podman。

## 常见问题

### Q: 为什么 Podman 可以无 sudo 安装？

A: Podman 支持 rootless 模式，可以作为普通用户运行。静态二进制版本不需要系统级安装。

### Q: 静态二进制和包管理器安装有什么区别？

A: 
- **静态二进制**: 不需要 root，安装在用户目录，只对当前用户可用
- **包管理器**: 需要 root，安装在系统目录，所有用户可用

功能上完全相同！

### Q: podman-compose 安装失败怎么办？

A: 可以使用 Podman 内置的 compose 功能：

```bash
# 不需要 podman-compose，直接使用
podman compose up -d
```

我们的启动脚本会自动处理这种情况。

### Q: PATH 配置不生效？

A: 确保重新加载配置：

```bash
# 重新加载
source ~/.bashrc

# 或打开新终端
# 或登出后重新登录

# 验证
echo $PATH | grep ".local/bin"
```

### Q: 下载速度很慢？

A: 可以使用国内镜像：

```bash
# 使用 GitHub 代理
export GITHUB_PROXY="https://ghproxy.com/"
curl -L -o podman.tar.gz \
  ${GITHUB_PROXY}https://github.com/containers/podman/releases/download/v4.9.3/podman-remote-static-linux_amd64.tar.gz
```

### Q: 我的架构是 ARM64 怎么办？

A: 下载 ARM64 版本：

```bash
curl -L -o podman.tar.gz \
  https://github.com/containers/podman/releases/download/v4.9.3/podman-remote-static-linux_arm64.tar.gz
```

脚本会自动检测架构。

## 故障排查

### Podman 命令找不到

```bash
# 检查是否安装
ls -lh ~/.local/bin/podman

# 检查 PATH
echo $PATH | grep ".local/bin"

# 手动添加到当前会话
export PATH="$HOME/.local/bin:$PATH"

# 测试
podman --version
```

### 权限问题

Podman rootless 模式需要配置 subuid/subgid：

```bash
# 检查配置
grep $USER /etc/subuid
grep $USER /etc/subgid

# 如果没有，请联系系统管理员添加
```

通常系统会自动配置，不需要手动操作。

### 容器运行失败

```bash
# 查看详细错误
podman run --rm alpine echo "test"

# 检查 Podman 信息
podman info

# 重置（如果需要）
podman system reset
```

## 完整安装流程示例

```bash
# 1. 运行安装脚本
./INSTALL_PODMAN.sh

# 2. 重新加载配置
source ~/.bashrc

# 3. 验证安装
podman --version
podman-compose --version

# 4. 启动应用
./start.sh

# 5. 访问应用
# 前端: http://localhost:3000
# 后端: http://localhost:8080
```

## 卸载

如果需要卸载用户空间安装的 Podman：

```bash
# 删除二进制文件
rm ~/.local/bin/podman
rm ~/.local/bin/podman-compose

# 删除数据（可选）
rm -rf ~/.local/share/containers
rm -rf ~/.config/containers

# 从 .bashrc 中删除 PATH 配置（可选）
# 手动编辑 ~/.bashrc，删除相关行
```

## 总结

- ✅ Podman 可以无 sudo 安装和使用
- ✅ 使用静态二进制或 pip 安装
- ✅ 安装到 `~/.local/bin/`
- ✅ 配置 PATH 后即可使用
- ✅ 功能与系统安装完全相同

**即使没有 sudo 权限，也能享受 Podman 的所有优势！** 🚀

## 相关资源

- [Podman 官方文档](https://docs.podman.io/)
- [Rootless Podman 教程](https://github.com/containers/podman/blob/main/docs/tutorials/rootless_tutorial.md)
- [静态二进制下载](https://github.com/containers/podman/releases)

