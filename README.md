# go-login

本项目提供一个通用的go版本的注册登录服务框架，该框架特点：
- 用session来管理登录态
- 基于redis来管理session
- 支持黑名单、分布式限流中间件
- xorm操作mysql
- viper管理配置文件
- gin为web server框架

## 效果图

![登录页面](./images/3.png)
![注册页面](./images/5.png)
![注册页面跳转](./images/4.png)
![主页](./images/1.png)

![修改页码页面](./images/2.png)

## 运行
项目根目录下
```
go run .
```