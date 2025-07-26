# 红包雨系统

一个基于Go语言开发的高性能红包雨系统，支持大规模并发抢红包场景，具备完整的反作弊机制和消息队列处理能力。

## 项目特性

- 🚀 **高性能并发**: 基于Gin框架，支持高并发抢红包请求
- 🛡️ **反作弊机制**: 内置多层反作弊策略，确保公平性
- 🔥 **消息队列**: 基于RocketMQ的异步消息处理
- 💾 **缓存优化**: Redis缓存提升系统性能
- 📊 **性能监控**: 集成pprof性能分析工具
- 🐳 **容器化部署**: 支持Docker容器化部署
- ⚡ **限流控制**: 令牌桶算法实现流量控制

## 项目架构

### 部署架构图
![部署架构图](/img/architecture.png)

### 技术架构图
![技术架构图](/img/technical-architecture.svg)

### 业务流程图
![业务流程图](/img/business-flow-diagram.svg)

### 架构说明

**分层架构设计**：
- **客户端层**: Web/移动端客户端，支持压力测试工具
- **负载均衡层**: 分发请求到多个应用实例
- **应用层**: Gin路由器 + 中间件（CORS、限流、反作弊、配置加载）
- **业务逻辑层**: 红包相关业务服务（抢红包、开红包、钱包、配置）
- **数据访问层**: MySQL、Redis、RocketMQ的客户端封装
- **存储层**: 数据持久化（MySQL）、缓存（Redis）、消息队列（RocketMQ）

**核心流程**：
1. **请求处理**: 客户端请求 → 负载均衡 → 中间件链 → 业务逻辑
2. **限流控制**: 令牌桶算法限制QPS，防止系统过载
3. **反作弊**: 检查用户请求间隔，防止恶意刷红包
4. **业务逻辑**: 概率性分配红包，检查用户限制
5. **数据处理**: Redis缓存用户状态，MySQL持久化红包记录
6. **异步处理**: RocketMQ处理异步任务，Consumer服务处理后续业务

## 技术栈

- **后端框架**: Gin (Go Web框架)
- **数据库**: MySQL + Redis
- **消息队列**: Apache RocketMQ
- **容器化**: Docker
- **性能分析**: pprof
- **限流**: 令牌桶算法

## API 接口

### 抢红包
```
POST /redenvelope/snatch
```
- 功能：用户抢红包
- 参数：uid (用户ID)
- 返回：红包ID、最大次数、当前次数等

### 开红包
```
POST /redenvelope/open
```
- 功能：打开已抢到的红包
- 参数：envelope_id (红包ID)
- 返回：红包金额

### 查看钱包
```
GET /redenvelope/wallet
```
- 功能：查看用户红包记录
- 参数：uid (用户ID)
- 返回：红包列表（包含金额、时间等）

### 配置管理
```
GET /redenvelope/config
POST /redenvelope/config
```
- 功能：获取/设置红包系统配置
- 配置项：概率、预算、数量、金额范围等

## 快速开始

### 环境准备

**必需软件**
- Go 1.16+
- MySQL 5.7+
- Redis 6.0+
- Docker (可选)

**初始化数据库**
```bash
# 修改generate.sh中的MySQL密码配置
vim generate.sh
# 执行数据库初始化脚本
sh generate.sh
```

### 配置文件

修改 `config.yaml` 文件，配置数据库连接、端口等参数：

```yaml
server:
  port: "8080"
database:
  host: "localhost"
  port: "3306"
  username: "root"
  password: "your_password"
  database: "red_envelope"
redis:
  host: "localhost"
  port: "6379"
```

### 启动服务

```bash
# 安装依赖
go mod download

# 启动主服务
go run main.go

# 启动消费者服务
cd consumer
go run consumer1.go
go run consumer2.go
```

### Docker 部署

```bash
# 构建镜像
docker build -t red-envelope:latest .

# 启动服务
docker-compose up -d
```

## 项目结构

```
├── api/redenvelope/          # 红包业务逻辑
│   ├── controller.go         # 控制器层
│   ├── service.go           # 业务逻辑层
│   ├── model.go             # 数据模型
│   ├── mapper.go            # 数据访问层
│   ├── router.go            # 路由定义
│   └── utils.go             # 工具函数
├── config/                  # 配置管理
├── database/                # 数据库连接
│   ├── database.go          # MySQL连接
│   ├── redis_db.go          # Redis连接
│   └── mq.go               # 消息队列连接
├── middleware/              # 中间件
│   ├── CheatPreventing.go   # 反作弊中间件
│   ├── Limiter.go          # 限流中间件
│   └── ConfigLoading.go     # 配置加载中间件
├── routers/                 # 路由配置
├── consumer/                # 消息队列消费者
└── main.go                 # 程序入口
```

## 性能测试

### 压力测试

使用wrk工具进行压力测试：

```bash
# 基础压力测试
wrk -t4 -c40 -d10s -s wrksnatch.lua http://localhost:8080/redenvelope/snatch
```

参数说明：
- `-t4`: 4个线程
- `-c40`: 40个并发连接
- `-d10s`: 持续10秒
- `-s wrksnatch.lua`: 使用自定义脚本（自动迭代uid）

### 性能分析

**实时性能监控**
```bash
# 访问pprof面板
http://localhost:8080/debug/pprof/
```

**生成火焰图**
```bash
# 启动性能采集并运行压测
go tool pprof -http=:1234 http://localhost:8080/debug/pprof/profile?seconds=10

# 访问可视化界面
http://localhost:1234
```

## 消息队列消费者

消费者服务位于 `consumer/` 目录，负责处理异步消息：

### 本地运行
```bash
cd consumer
go run consumer1.go  # 消费者1
go run consumer2.go  # 消费者2
```

### 容器化部署
```bash
# 构建消费者镜像
docker build -t consumer1:v0.3 -f Dockerfile1 .
docker build -t consumer2:v0.3 -f Dockerfile2 .

# 推送到镜像仓库
docker push cr-cn-beijing.volces.com/group1/consumer1:v0.3
docker push cr-cn-beijing.volces.com/group1/consumer2:v0.3
```

## 重要配置说明

### 限流配置
⚠️ **重要提醒**: 在新环境部署时，需要调整限流参数：

1. 测试阶段：注释掉 `routers.go` 中的令牌桶限制，或调大限流参数
2. 生产环境：根据压力测试结果，设置合适的限流值

### 反作弊配置
系统内置多层反作弊机制，可在配置文件中调整相关参数。

## 开发指南

### 添加新功能
1. 在 `api/redenvelope/` 下添加相应的model、service、controller
2. 在 `router.go` 中注册新路由
3. 更新数据库表结构（如需要）

### 调试技巧
- 使用pprof分析性能瓶颈
- 查看Redis缓存命中率
- 监控消息队列消费速度

## 贡献指南

1. Fork 本项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 数据库设计

### 表结构

**red_envelope 表**
```sql
CREATE TABLE red_envelope (
    envelope_id INT NOT NULL,           -- 红包ID (主键)
    uid INT NOT NULL,                   -- 用户ID
    opened BOOL DEFAULT FALSE,          -- 是否已开启
    value INT DEFAULT NULL,             -- 红包金额(分)
    snatch_time INT NOT NULL,           -- 抢红包时间戳
    INDEX index_id (uid),
    PRIMARY KEY(envelope_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

## 配置参数详解

### 服务配置
```yaml
server:
  port: 8080                    # 服务端口

database:
  driverName: mysql             # 数据库驱动
  userName: developer           # 数据库用户名
  password: Group1234.          # 数据库密码
  host: 111.62.122.65          # 数据库地址
  database: red_envelope_rain   # 数据库名
  MaxOpenConns: 100            # 最大连接数
  MaxIdleConns: 100            # 最大空闲连接数
  ConnMaxLifeTime: 300         # 连接生命周期(秒)

redis:
  addr: redis:6379             # Redis地址
  username: developer          # Redis用户名
  password: Group1234.         # Redis密码
  dbNumber: 0                  # Redis数据库编号
```

### 业务配置
```yaml
snatchProbability: 70          # 抢红包成功概率(%)
snatchMaxCount: 5              # 用户最大抢红包次数
totalNum: 100                  # 红包总数量
totalAmount: 100000            # 红包总金额(分)
minAmount: 100                 # 单个红包最小金额(分)
maxAmount: 3000                # 单个红包最大金额(分)
```

### 限流配置
```yaml
limitRate: 5000                # 令牌桶每秒放入令牌数
limitCapacity: 5000            # 令牌桶容量
```

### 反作弊配置
```yaml
cheat-preventing:
  milliseconds: 1000           # 同一用户两次请求最小间隔(毫秒)
```

## API 错误码说明

| 错误码 | 含义 | 说明 |
|--------|------|------|
| 0 | 成功 | 操作成功 |
| 1 | 没有抢到红包 | 概率性失败，非错误 |
| 2 | 达到抢红包上限 | 用户已达到最大抢红包次数 |
| 3 | 参数错误 | 请求参数格式错误 |
| 4 | 红包已开启 | 重复开启红包 |
| 5 | 红包不存在 | 红包ID不存在或不属于该用户 |
| 6 | 系统错误 | 服务器内部错误 |
| 7 | 请求过于频繁 | 触发反作弊限制 |

## 测试指南

### 运行单元测试
```bash
# 运行所有测试
go test -v

# 运行特定测试
go test -v -run TestSnatchRedEnvelope
go test -v -run TestOpenRedEnvelope
go test -v -run TestGetWalletList
go test -v -run TestWorkflow
```

### 测试用例说明
- **TestSnatchRedEnvelope**: 测试抢红包功能
- **TestOpenRedEnvelope**: 测试开红包功能  
- **TestGetWalletList**: 测试获取钱包列表
- **TestWorkflow**: 测试完整业务流程

### 压力测试配置
修改 `wrksnatch.lua` 中的用户ID范围：
```lua
-- 根据数据库中的用户数量调整
local uid = math.random(1, 1000)  -- 1-1000为用户ID范围
```

## 监控和日志

### 性能监控
- **CPU和内存**: 通过pprof监控 `http://localhost:8080/debug/pprof/`
- **QPS监控**: 查看令牌桶使用情况
- **数据库连接**: 监控连接池状态

### 日志配置
```bash
# 查看实时日志
tail -f /var/log/red-envelope/app.log

# 查看错误日志
grep "ERROR" /var/log/red-envelope/app.log
```

## 安全机制

### 反作弊策略
1. **时间间隔限制**: 同一用户两次请求间隔不得少于1秒
2. **次数限制**: 每个用户最多抢5次红包
3. **令牌桶限流**: 全局请求频率控制
4. **参数验证**: 严格的输入参数校验

### 数据安全
- **密码保护**: 数据库密码加密存储
- **SQL注入防护**: 使用ORM防止SQL注入
- **Redis安全**: 配置密码认证

## 故障排查

### 常见问题

**1. 服务启动失败**
```bash
# 检查端口占用
lsof -i :8080

# 检查配置文件
cat config.yaml
```

**2. 数据库连接失败**
```bash
# 测试数据库连接
mysql -h111.62.122.65 -udeveloper -pGroup1234. -Dred_envelope_rain

# 检查网络连通性
ping 111.62.122.65
```

**3. Redis连接失败**
```bash
# 测试Redis连接
redis-cli -h redis -p 6379 -a Group1234.

# 检查Redis状态
redis-cli info
```

**4. 消息队列问题**
```bash
# 检查RocketMQ状态
docker logs rmqnamesrv
docker logs rmqbroker
```

### 性能优化建议

1. **数据库优化**
   - 添加适当索引
   - 定期分析慢查询
   - 调整连接池参数

2. **缓存优化**
   - 提高Redis命中率
   - 合理设置过期时间
   - 监控缓存使用情况

3. **限流调优**
   - 根据实际QPS调整令牌桶参数
   - 分布式限流考虑

## 生产环境部署

### 环境要求
- **服务器**: 4核8G以上
- **MySQL**: 5.7+ 主从配置
- **Redis**: 6.0+ 集群模式
- **RocketMQ**: 4.8+ 集群部署

### 部署检查清单
- [ ] 修改生产环境配置
- [ ] 调整限流参数
- [ ] 配置监控告警
- [ ] 备份数据库
- [ ] 配置日志轮转
- [ ] 测试故障恢复

## License

本项目采用 MIT 许可证，详情请查看 [LICENSE](LICENSE) 文件。