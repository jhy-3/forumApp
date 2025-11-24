#!/bin/bash

echo "========================================"
echo "  论坛应用 - 完整安装和启动"
echo "========================================"
echo ""
echo "这个脚本会："
echo "  1. 使用 Nix 安装完整版 Podman"
echo "  2. 配置环境"
echo "  3. 启动论坛应用"
echo ""
echo "需要时间：约 15-20 分钟"
echo "需要空间：约 1-2 GB"
echo ""
read -p "按 Enter 继续，或 Ctrl+C 取消..." 

echo ""
echo "===== 步骤 1/4: 安装 Nix 和 Podman ====="
./INSTALL_PODMAN_WITH_NIX.sh

echo ""
echo "===== 步骤 2/4: 加载环境 ====="
if [ -e ~/.nix-profile/etc/profile.d/nix.sh ]; then
    source ~/.nix-profile/etc/profile.d/nix.sh
    echo "✓ Nix 环境已加载"
else
    echo "✗ 找不到 Nix 环境"
    exit 1
fi

echo ""
echo "===== 步骤 3/4: 配置 .bashrc ====="
if ! grep -q "nix-profile/etc/profile.d/nix.sh" ~/.bashrc; then
    echo 'if [ -e ~/.nix-profile/etc/profile.d/nix.sh ]; then . ~/.nix-profile/etc/profile.d/nix.sh; fi' >> ~/.bashrc
    echo "✓ 已添加到 .bashrc"
else
    echo "✓ 已经在 .bashrc 中"
fi

echo ""
echo "===== 步骤 4/4: 启动应用 ====="
if command -v podman &> /dev/null; then
    echo "✓ Podman 可用: $(podman --version)"
    echo ""
    echo "启动论坛应用..."
    ./start.sh
else
    echo "✗ Podman 未找到"
    echo "请手动运行："
    echo "  source ~/.nix-profile/etc/profile.d/nix.sh"
    echo "  ./start.sh"
fi
