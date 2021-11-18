#### 红包雨

## 本地开发准备

- **第三方软件**
  - golang
  - mysql
- **生成本地数据库**
  - 查看generate.sh脚本，根据本地设置修改root密码
  - sh generate.sh

## 文件说明

- **配置文件config.yml**

  项目中的所有配置将在此处进行设置，如数据库账号密码、服务监听端口等

  **在新的环境进行部署时，需要对routers.go中的令牌桶进行注释，或者将本文件下的令牌桶参数尽可能地调大，避免影响测试。在经过压力测试后，再将令牌桶参数调至合适的值**；

- **database文件夹**

  - database.go 连接数据库
  
- **wrktest.lua**

  - 本文件用于压力测试，文件中实现了每个请求自动迭代uid。使用时需要根据数据库已有uid进行设置

  - 命令行语句 

    ``````
    wrk -t4 -c40 -d10s -swrksnatch.lua http://localhost:8080/redenvelope/snatch
    ``````

    `-t`后面接线程数  `-c`后面接连接数，连接数需要根据情况调测

## 测试

- **性能测试**

  - main.go中在测试环境添加"github.com/gin-contrib/pprof"包用于性能分析，生产环境中删去该包

  - 使用方式：运行过程中，可以在 "http://localhost:8080/debug/pprof/" 中查看程序性能；
              压测场景下，可以先运行 "go tool pprof -http=:1234 http://localhost:8080/debug/pprof/profile?second=10" 再马上运行压力测试脚本，之后可以在 "http://localhost:1234" 中查看图形化堆栈调用（Graph）和火焰图（Flame Graph）等。

