#!/bin/bash

# 显示标题
echo "====================================="
echo "  Fridge 智能冰箱管理系统启动脚本"
echo "====================================="
echo ""

# 检查是否有旧进程在运行
echo "检查是否有旧进程在运行..."
BACKEND_PID=$(pgrep -f "go run.*backend")
FRONTEND_PID=$(pgrep -f "npm run dev")

if [ ! -z "$BACKEND_PID" ] || [ ! -z "$FRONTEND_PID" ]; then
  echo "检测到旧进程正在运行，请先运行 ./stop.sh 停止它们"
  exit 1
fi

# 启动后端服务
echo "启动后端服务..."
cd backend
go run . > ../backend.log 2>&1 &
BACKEND_PID=$!
cd ..

# 等待后端服务启动
echo "等待后端服务启动..."
sleep 2

# 检查后端服务是否成功启动
if ! ps -p $BACKEND_PID > /dev/null; then
  echo "后端服务启动失败，请检查 backend.log 文件获取详细信息"
  exit 1
fi

echo "后端服务已启动，PID: $BACKEND_PID"
echo "后端日志保存在: $(pwd)/backend.log"

# 启动前端服务
echo "启动前端服务..."
cd frontend
npm run dev > ../frontend.log 2>&1 &
FRONTEND_PID=$!
cd ..

# 等待前端服务启动
echo "等待前端服务启动..."
sleep 5

# 检查前端服务是否成功启动
if ! ps -p $FRONTEND_PID > /dev/null; then
  echo "前端服务启动失败，请检查 frontend.log 文件获取详细信息"
  kill $BACKEND_PID
  exit 1
fi

echo "前端服务已启动，PID: $FRONTEND_PID"
echo "前端日志保存在: $(pwd)/frontend.log"

# 保存进程ID到文件中，以便停止脚本使用
echo "$BACKEND_PID" > backend.pid
echo "$FRONTEND_PID" > frontend.pid

echo ""
echo "====================================="
echo "  所有服务已成功启动"
echo "  前端访问地址: http://localhost:3000"
echo "  后端API地址: http://localhost:8602"
echo "  使用 ./stop.sh 可以停止所有服务"
echo "====================================="
