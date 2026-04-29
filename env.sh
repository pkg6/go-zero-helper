#!/bin/bash
set -euo pipefail

# ====================== 【文件映射表】======================
# 格式：goctl模板内文件路径=本地文件 或 远程URL
TEMPLATE_MAP=(
  "api/handler.tpl=https://raw.githubusercontent.com/pkg6/go-zero-helper/refs/heads/main/goctl/template/api/handler.tpl"
)
# =================================================================

# 从 go.mod 自动获取 go-zero 版本
if [ ! -f "go.mod" ]; then
  echo "错误：当前目录未找到 go.mod 文件"
  exit 1
fi
GO_ZERO_VERSION=$(grep 'github.com/zeromicro/go-zero' go.mod | awk '{print $2}' | sed 's/^v//')
if [ -z "$GO_ZERO_VERSION" ]; then
  echo "错误：未在 go.mod 中找到 go-zero 依赖"
  exit 1
fi

GOCTL_TEMPLATE_DIR="$HOME/.goctl/$GO_ZERO_VERSION"

echo "==================================================="
echo " 正在初始化 go-zero 开发环境 v$GO_ZERO_VERSION"
echo "==================================================="

# 检查 Go
if ! command -v go &> /dev/null; then
  echo "错误：未安装 Go 环境"
  exit 1
fi

# 安装 goctl
if ! command -v goctl &> /dev/null; then
  echo "正在安装 goctl..."
  go install github.com/zeromicro/go-zero/tools/goctl@v$GO_ZERO_VERSION
else
  echo "goctl 已安装"
fi

echo "检查 goctl 环境并安装依赖..."
if ! goctl env check --install --verbose --force; then
  echo "❌ goctl 环境检查失败"
  exit 1
fi

echo "初始化官方模板..."
goctl template init

# ====================== 自动处理模板文件 ======================
echo -e "\n==================================================="
echo " 开始覆盖自定义模板文件（支持本地/URL）"
echo "==================================================="

for item in "${TEMPLATE_MAP[@]}"; do
  IFS='=' read -r target_file source <<< "$item"
  dest="$GOCTL_TEMPLATE_DIR/$target_file"
  mkdir -p "$(dirname "$dest")"

  if [[ "$source" =~ ^https?:// ]]; then
    echo "🌐 下载远程文件: $source"
    echo "└─→ 覆盖到: $dest"
    if ! curl -fsSL "$source" -o "$dest"; then
      echo "❌ 下载失败: $source"
      exit 1
    fi
  else
    if [ -f "$source" ]; then
      echo "📁 复制本地文件: $source"
      echo "└─→ 覆盖到: $dest"
      cp -f "$source" "$dest"
    else
      echo "❌ 错误：本地文件不存在 $source"
      exit 1
    fi
  fi

  echo "✅ 成功处理: $target_file"
  echo
done

echo "==================================================="
echo " ✅ 环境初始化全部完成！"
echo "==================================================="