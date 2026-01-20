# E-Commerce Checkout Processor

A distributed e-commerce checkout system demonstrating the **Saga pattern** for distributed transactions.

## Tech Stack

- **Go** - Application language
- **Temporal** - Workflow orchestration engine
- **NATS** - High-performance messaging
- **PostgreSQL** - Temporal state storage

## Architecture

```
┌─────────────┐     ┌──────────────┐     ┌─────────────────┐
│   Trigger   │────▶│   Temporal   │────▶│  OrderWorkflow  │
└─────────────┘     └──────────────┘     └────────┬────────┘
                                                  │
                    ┌─────────────────────────────┼─────────────────────────────┐
                    │                             │                             │
                    ▼                             ▼                             ▼
            ┌──────────────┐             ┌──────────────┐             ┌──────────────┐
            │   Payment    │             │  Inventory   │             │    Refund    │
            │   Service    │             │   Service    │             │   Service    │
            └──────────────┘             └──────────────┘             └──────────────┘
                    │                             │                             │
                    └─────────────────────────────┴─────────────────────────────┘
                                                  │
                                                  ▼
                                            ┌──────────┐
                                            │   NATS   │
                                            └──────────┘
```

## Project Structure

```
e-comm_processor/
├── cmd/
│   └── processor/
│       └── main.go              # Application entry point
├── configs/
│   └── config.yaml              # Configuration file
├── internal/
│   ├── config/                  # Configuration loading
│   ├── domain/                  # Domain models
│   ├── mock/                    # Mock NATS services
│   ├── nats/                    # NATS client wrapper
│   └── workflow/                # Temporal workflows & activities
├── docker-compose.yml
└── README.md
```

## Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose

### 1. Start Infrastructure

```bash
docker-compose up -d
```

### 2. Run Worker

```bash
go run ./cmd/processor -mode=worker
```

### 3. Trigger Order (in another terminal)

```bash
go run ./cmd/processor -mode=trigger
```

## Configuration

Configuration via `configs/config.yaml` or environment variables:

| Config | Env Variable | Default |
|--------|--------------|---------|
| temporal.host | TEMPORAL_HOST | localhost:7233 |
| nats.url | NATS_URL | nats://localhost:4222 |

## Testing

```bash
# Run all tests
go test ./...

# With coverage
go test ./... -cover

# Verbose output
go test ./... -v
```

## Why This Architecture?

- **Temporal** - Provides durable execution, automatic retries, and workflow state persistence
- **NATS** - Enables fast pub/sub messaging with backpressure support
- **Saga Pattern** - Handles distributed transactions with compensation logic (refund + restock on cancellation)

---

# E-Commerce 结账处理器

一个展示分布式事务 **Saga 模式** 的电商结账系统。

## 技术栈

- **Go** - 应用程序语言
- **Temporal** - 工作流编排引擎
- **NATS** - 高性能消息中间件
- **PostgreSQL** - Temporal 状态存储

## 项目结构

```
e-comm_processor/
├── cmd/
│   └── processor/
│       └── main.go              # 程序入口
├── configs/
│   └── config.yaml              # 配置文件
├── internal/
│   ├── config/                  # 配置加载
│   ├── domain/                  # 领域模型
│   ├── mock/                    # 模拟 NATS 服务
│   ├── nats/                    # NATS 客户端封装
│   └── workflow/                # Temporal 工作流和 Activity
├── docker-compose.yml
└── README.md
```

## 快速开始

### 环境要求

- Go 1.21+
- Docker & Docker Compose

### 1. 启动基础设施

```bash
docker-compose up -d
```

### 2. 启动 Worker

```bash
go run ./cmd/processor -mode=worker
```

### 3. 触发订单（新终端）

```bash
go run ./cmd/processor -mode=trigger
```

## 配置

通过 `configs/config.yaml` 或环境变量配置：

| 配置项 | 环境变量 | 默认值 |
|--------|----------|--------|
| temporal.host | TEMPORAL_HOST | localhost:7233 |
| nats.url | NATS_URL | nats://localhost:4222 |

## 测试

```bash
# 运行所有测试
go test ./...

# 带覆盖率
go test ./... -cover
```

## 为什么选择这个架构？

- **Temporal** - 提供持久化执行、自动重试和工作流状态持久化
- **NATS** - 支持快速发布/订阅消息和背压机制
- **Saga 模式** - 通过补偿逻辑处理分布式事务（取消时退款+恢复库存）