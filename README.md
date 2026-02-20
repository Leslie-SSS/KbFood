# KbFood

> 美食价格监控平台 - 追踪多平台食品价格，智能降价提醒

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![React](https://img.shields.io/badge/React-19-61DAFB?style=flat&logo=react)](https://reactjs.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![SQLite](https://img.shields.io/badge/SQLite-3-003B57?style=flat&logo=sqlite)](https://sqlite.org)

## 特性

- **多平台数据聚合** - 支持探探糖、多堂、小蚕等平台
- **价格趋势图表** - 可视化价格变化历史，掌握价格走势
- **Bark 推送通知** - 降价到目标价格时推送提醒到 iPhone
- **自定义监控** - 设置目标价格，精准追踪心仪商品
- **响应式设计** - 完美适配移动端和桌面端

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- SQLite 3

### 安装

```bash
# 克隆项目
git clone https://github.com/Leslie-SSS/KbFood.git
cd KbFood

# 安装后端依赖
go mod download

# 安装前端依赖
cd frontend && npm install && cd ..
```

### 运行

```bash
# 启动后端服务 (默认端口 9000)
go run cmd/server/main.go

# 启动前端开发服务器 (默认端口 5173)
cd frontend && npm run dev
```

访问 http://localhost:5173 即可使用。

### Docker 部署

```bash
cd deployments
docker-compose up -d
```

访问 http://localhost:8200

## 配置

配置文件: `config.yaml`

```yaml
server:
  port: 9000

tantantang:
  # 已预配置，开箱即用
  token: "b69755a8-cff3-4b2f-9c4e-c7c6d19efd82"
  secret_key: "~~JSTtT*(lvlv!#^&%~~"
  base_url: "https://ttt.bjlxkjyxgs.cn/api/shop/activity"

database:
  path: "./data/food.db"
```

## 技术栈

| 后端 | 前端 |
|------|------|
| Go 1.24 | React 19 |
| Echo | TypeScript |
| SQLite | Vite 7 |
| Wire (DI) | Tailwind CSS 4 |
| sqlc | React Query |
| Cron | Chart.js |

## 项目结构

```
KbFood/
├── cmd/server/              # 入口程序
│   ├── main.go              # 主程序
│   └── wire.go              # 依赖注入
├── internal/
│   ├── config/              # 配置管理
│   ├── domain/              # 领域层
│   │   ├── entity/          # 实体
│   │   ├── repository/      # 仓库接口
│   │   └── service/         # 业务服务
│   ├── infra/               # 基础设施
│   │   ├── db/              # 数据库
│   │   ├── external/        # 外部服务 (Bark)
│   │   ├── platform/        # 平台适配器
│   │   ├── repository/      # 仓库实现
│   │   └── scheduler/       # 定时任务
│   └── interface/           # 接口层
│       └── http/            # HTTP 处理
├── frontend/                # React 前端
│   ├── src/
│   │   ├── components/      # UI 组件
│   │   ├── hooks/           # 自定义 Hooks
│   │   ├── services/        # API 服务
│   │   └── types/           # 类型定义
│   └── e2e/                 # E2E 测试
├── deployments/             # Docker 配置
│   ├── Dockerfile
│   ├── docker-compose.yml
│   └── nginx.conf
├── config.yaml              # 配置文件
├── Makefile                 # 构建脚本
└── LICENSE                  # MIT 许可证
```

## API 端点

| 方法 | 端点 | 描述 |
|------|------|------|
| GET | `/api/products` | 获取商品列表 |
| GET | `/api/products/:id/trend` | 获取价格趋势 |
| POST | `/api/notifications` | 设置价格提醒 |
| PUT | `/api/notifications/:id` | 更新价格提醒 |
| DELETE | `/api/notifications/:id` | 删除价格提醒 |
| POST | `/admin/test-notification` | 测试推送通知 |
| GET | `/health` | 健康检查 |

## 开发

### 构建

```bash
# 构建后端
make build

# 构建前端
cd frontend && npm run build
```

### 测试

```bash
# 后端单元测试
make test

# 前端单元测试
cd frontend && npm run test

# 前端 E2E 测试
cd frontend && npm run test:e2e
```

## 致谢

- 数据来源: [探探糖](https://ttt.bjlxkjyxgs.cn)
- 推送服务: [Bark](https://github.com/Finb/Bark)

## 许可证

[MIT License](LICENSE)
