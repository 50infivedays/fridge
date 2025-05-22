#!/bin/bash

# 显示标题
echo "====================================="
echo "  Fridge 智能冰箱管理系统停止脚本"
echo "====================================="
echo ""

# 尝试从PID文件获取进程ID
if [ -f "backend.pid" ]; then
  BACKEND_PID=$(cat backend.pid)
  echo "从PID文件找到后端进程ID: $BACKEND_PID"
else
  # 如果PID文件不存在，尝试查找进程
  BACKEND_PID=$(pgrep -f "go run.*backend")
  echo "通过进程查找找到后端进程ID: $BACKEND_PID"
fi

if [ -f "frontend.pid" ]; then
  FRONTEND_PID=$(cat frontend.pid)
  echo "从PID文件找到前端进程ID: $FRONTEND_PID"
else
  # 如果PID文件不存在，尝试查找进程
  FRONTEND_PID=$(pgrep -f "npm run dev")
  echo "通过进程查找找到前端进程ID: $FRONTEND_PID"
fi

# 停止后端服务
if [ ! -z "$BACKEND_PID" ]; then
  echo "正在停止后端服务 (PID: $BACKEND_PID)..."
  kill $BACKEND_PID 2>/dev/null || kill -9 $BACKEND_PID 2>/dev/null
  if [ $? -eq 0 ]; then
    echo "后端服务已停止"
    rm -f backend.pid
  else
    echo "无法停止后端服务，可能已经不在运行"
  fi
else
  echo "未找到运行中的后端服务"
fi

# 停止前端服务
if [ ! -z "$FRONTEND_PID" ]; then
  echo "正在停止前端服务 (PID: $FRONTEND_PID)..."
  kill $FRONTEND_PID 2>/dev/null || kill -9 $FRONTEND_PID 2>/dev/null
  if [ $? -eq 0 ]; then
    echo "前端服务已停止"
    rm -f frontend.pid
  else
    echo "无法停止前端服务，可能已经不在运行"
  fi
else
  echo "未找到运行中的前端服务"
fi

# 清理日志文件（可选，取消注释以启用）
# echo "清理日志文件..."
# rm -f backend.log frontend.log

# 查找可能遗留的Node.js进程（Next.js可能会启动多个进程）
NODE_PIDS=$(pgrep -f "node.*next")
if [ ! -z "$NODE_PIDS" ]; then
  echo "发现可能的Next.js相关进程，正在停止..."
  for PID in $NODE_PIDS; do
    echo "停止进程 $PID..."
    kill $PID 2>/dev/null || kill -9 $PID 2>/dev/null
  done
fi

echo ""
echo "====================================="
echo "  所有服务已停止"
echo "====================================="
