#### 红包雨

## 本地开发准备

- **第三方软件**
  - golang
  - mysql
- **生成本地数据库**
  - 查看generate.sh脚本，根据本地设置修改root密码
  - sh generate.sh

## 文件说明

- **配置文件app.yml**

  项目中的所有配置将在此处进行设置，如数据库账号密码、服务监听端口等

- **database文件夹**

  - database.go 连接数据库
  
- **wrktest.lua**

  - 本文件用于压力测试，文件中实现了每个请求自动迭代uid。使用时需要根据数据库已有uid进行设置

  - 命令行语句 

    ``````
    wrk -t4 -c40 -d10s -swrktest.lua http://localhost:8080/redenvelope/snatch
    ``````

    `-t`后面接线程数  `-c`后面接连接数，连接数需要根据情况调测

