# Fridge 智能冰箱管理系统

这是一个使用Go和Next.js实现的智能冰箱管理系统，用于管理冰箱中的物品。该系统利用Google的Gemini API将自然语言描述转换为结构化的JSON数据，并提供直观的Web界面进行操作。

## 系统架构

- **后端**：基于Go语言实现的API服务
- **前端**：使用Next.js 15、React 19和TailwindCSS构建的响应式Web应用
- **AI处理**：利用Google Gemini API进行自然语言处理

## 功能特点

- 使用自然语言描述添加冰箱物品
- 自动识别物品名称、数量、单位和过期时间
- 直观的Web界面展示和管理物品

## 安装前提

- Go 1.21或更高版本
- Node.js 18或更高版本
- npm或yarn包管理工具
- Google Gemini API密钥

## 安装步骤

### 1. 克隆仓库

```bash
git clone https://github.com/yourusername/fridge.git
cd fridge
```

### 2. 配置后端

创建或编辑`config.yaml`文件，设置您的Gemini API密钥：

```yaml
# Gemini API配置
gemini_api_key: "your_api_key_here"
gemini_model: "gemini-1.5-flash"  # 默认模型，可以更改为其他Gemini模型
```

或者，您也可以使用环境变量进行配置：

```bash
export GEMINI_API_KEY=your_api_key_here
export GEMINI_MODEL=gemini-1.5-flash  # 可选
```

### 3. 编译后端

使用提供的构建脚本编译后端：

```bash
chmod +x build.sh
./build.sh
```

或者手动编译：

```bash
cd backend
go build -o ../fridge *.go
```

### 4. 安装前端依赖

```bash
cd frontend
npm install
```

## 运行系统

### 方式一：使用一键启动脚本

项目提供了一键启动和停止脚本，可以方便地管理整个应用的生命周期：

#### 启动应用

在项目根目录下运行：

```bash
./start.sh
```

该脚本会同时启动后端和前端服务。

#### 停止应用

在项目根目录下运行：

```bash
./stop.sh
```

该脚本会停止所有相关的服务进程。

### 方式二：手动启动各组件

如果您希望手动控制各组件的启动，可以按照以下步骤操作：

#### 1. 启动后端服务

在项目根目录下运行：

```bash
./fridge
```

或者进入backend目录运行：

```bash
cd backend
go run .
```

后端服务将在`http://localhost:8602`上启动。

#### 2. 启动前端开发服务器

在另一个终端窗口中，进入frontend目录并启动开发服务器：

```bash
cd frontend
npm run dev
```

前端应用将在`http://localhost:3000`上启动。

### 3. 访问应用

无论使用哪种方式启动，应用都可以通过浏览器访问：

- 打开浏览器，访问`http://localhost:3000`即可使用应用
- 后端API服务地址为`http://localhost:8602`

## 使用指南

### 1. 首页

首页提供了系统的主要功能概览：
- SmartExpiry：智能物品管理
- 即将过期：过期提醒（计划中）
- 购物清单：智能推荐（计划中）

### 2. 添加物品（SmartExpiry）

1. 点击首页上的"SmartExpiry"卡片中的"立即使用"按钮
2. 在"物品描述"输入框中，使用自然语言描述要添加的物品，例如：
   - "三个苹果，过期时间是下周一"
   - "两盒牛奶，保质期5天；一袋面包，下周五过期"
   - "500克猪肉，3天后过期"
3. 点击"添加到清单"按钮
4. 系统将解析您的描述，并在下方的表格中显示解析结果

### 系统处理逻辑

系统会自动识别：
- 物品名称（如"苹果"、"牛奶"）
- 数量（如"三个"、"两盒"）
- 单位（如"个"、"盒"）
- 过期时间（如"下周一"、"5天后"）

如果未指定某些信息，系统会使用合理的默认值：
- 未指定数量时，默认为1
- 未指定单位时，使用适合该物品的常见单位
- 未指定过期时间时，默认为7天后

## API接口说明

### 记录物品 (POST /record)

**请求格式**:

```json
{
  "description": "三个苹果，过期时间是2025年6月1日",
  "currentTime": "2023-05-20 10:30:00"  // 可选
}
```

**响应格式**:

```json
{
  "items": [
    {
      "item": "苹果",
      "quantity": 3,
      "unit": "个",
      "expireDate": "2025-06-01 00:00:00"
    }
  ]
}
```

## 示例请求

使用curl发送请求：

```bash
curl -X POST http://localhost:8602/record \
  -H "Content-Type: application/json" \
  -d '{"description": "两斤猪肉，过期时间是下周五"}'
```

## 注意事项

- 确保设置了有效的Google Gemini API密钥
- 自然语言描述应尽可能包含物品名称、数量、单位和过期时间信息
- 前端开发使用Next.js的Turbopack功能，确保您的Node.js版本兼容
- 在生产环境中，您可能需要设置跨域资源共享(CORS)策略

## 故障排除

如果遇到问题，请检查：

1. Gemini API密钥是否正确配置
2. 前端应用的API端点是否正确指向后端服务
3. 后端服务是否正常运行在端口8602上

## 许可证

[LICENSE文件中的相关许可证]
