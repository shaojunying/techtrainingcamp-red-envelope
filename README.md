# 红包雨系统

一个基于Go语言开发的高性能红包雨系统，采用标准Go项目布局，支持大规模并发抢红包场景。

## 📁 项目结构

```
red-envelope/
├── cmd/                          # 应用程序入口
│   └── server/                   # 主服务入口
│       └── main.go
├── internal/                     # 私有应用代码
│   ├── app/                      # 应用程序层
│   │   ├── config.go            # 配置管理
│   │   └── middleware/          # 中间件
│   ├── domain/                   # 领域层
│   │   └── envelope/            # 红包领域
│   ├── infrastructure/           # 基础设施层
│   │   └── database/            # 数据库、Redis、MQ
│   └── interface/               # 接口层
│       ├── http/                # HTTP接口
│       │   └── router/          # 路由定义
│       └── consumer/            # 消息队列消费者
├── configs/                     # 配置文件
│   └── config.yaml
├── scripts/                     # 脚本文件
│   └── generate.sh
├── deployments/                 # 部署配置
│   ├── docker-compose.yml
│   └── Dockerfile
├── docs/                        # 文档
│   ├── README.md               # 详细文档
│   └── img/                    # 架构图
├── test/                        # 测试文件
│   ├── main_test.go
│   └── load/
│       └── wrksnatch.lua
├── go.mod
└── go.sum
```

## 🚀 快速开始

### 环境准备
- Go 1.16+
- MySQL 5.7+
- Redis 6.0+
- RocketMQ 4.8+

### 初始化数据库
```bash
# 修改数据库配置
vim configs/config.yaml

# 执行数据库初始化
sh scripts/generate.sh
```

### 运行服务
```bash
# 安装依赖
go mod download

# 启动主服务
go run cmd/server/main.go

# 启动消费者服务
cd internal/interface/consumer
go run consumer1.go
go run consumer2.go
```

### Docker部署
```bash
cd deployments
docker-compose up -d
```

## 📚 详细文档

完整的项目文档请查看：[docs/README.md](docs/README.md)

## 🏗️ 架构特点

- ✅ **标准Go项目布局** - 符合Go社区最佳实践
- ✅ **清晰的分层架构** - 领域驱动设计
- ✅ **高性能并发** - 支持大规模并发请求
- ✅ **完整的监控** - 性能监控和日志记录
- ✅ **容器化部署** - Docker容器化支持

## 🔧 项目重构

本项目已完成标准化重构，主要改进：

1. **目录结构标准化** - 采用Go标准项目布局
2. **模块化设计** - 清晰的分层和职责分离
3. **依赖管理优化** - 合理的import路径组织
4. **部署配置分离** - 独立的部署和配置目录

## 📄 License

MIT License - 详见 [LICENSE](LICENSE) 文件