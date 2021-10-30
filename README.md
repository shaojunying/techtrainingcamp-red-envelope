#### 红包雨

## 本地开发准备

- **第三方软件**
  - golang
  - mysql
- **生成本地数据库**
  - 查看generate.sh脚本，根据本地设置修改root密码
  - sh generate.sh

## 文件结构说明

- **配置文件app.yml**

  项目中的所有配置将在此处进行设置，如数据库账号密码、服务监听端口等

- **database文件夹**

  - database.go 连接数据库

