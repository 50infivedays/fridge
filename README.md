# Fridge Item Management API

这是一个使用Go实现的API服务，用于管理冰箱中的物品。该服务利用Google的Gemini API将自然语言描述转换为结构化的JSON数据。

## 功能

- `/record` 端点：接收关于要放入冰箱的物品的自然语言描述，并返回结构化的JSON数据

## 前提条件

- Go 1.21或更高版本
- Gemini API密钥

## 配置

应用程序支持从`config.yaml`文件读取配置。您可以基于提供的`config.yaml.example`模板创建此文件。

```yaml
# Gemini API配置
gemini_api_key: "your_api_key_here"
gemini_model: "gemini-pro"  # 默认模型，可以更改为其他Gemini模型
```

配置在服务启动时加载一次并缓存在内存中。

### 配置查找顺序

1. 应用程序首先在当前目录中查找`config.yaml`
2. 如果未找到，它会在可执行文件目录中查找
3. 如果仍未找到，它会回退到环境变量

### 环境变量

您也可以使用环境变量而不是配置文件：

- `GEMINI_API_KEY`：您的Gemini API密钥
- `GEMINI_MODEL`：（可选）要使用的Gemini模型（默认为"gemini-pro"）

## 安装和运行

1. 克隆此仓库
2. 创建配置文件或设置环境变量：
   ```
   cp config.yaml.example config.yaml
   # 编辑config.yaml设置您的API密钥
   ```
   或
   ```
   export GEMINI_API_KEY=your_api_key_here
   ```
3. 运行服务：
   ```
   go run *.go
   ```
   服务将在`http://localhost:8602`上启动

## API使用说明

### 记录物品 (POST /record)

**请求格式**:

```json
{
  "description": "三个苹果，过期时间是2025年6月1日"
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

## 示例

使用curl发送请求：

```bash
curl -X POST http://localhost:8080/record \
  -H "Content-Type: application/json" \
  -d '{"description": "两斤猪肉，过期时间是下周五"}'
```

## 注意事项

- 确保设置了有效的Gemini API密钥
- 自然语言描述应尽可能包含物品名称、数量、单位和过期时间信息
