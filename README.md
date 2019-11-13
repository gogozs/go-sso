## golang sso

## project start

1. 参考config/example_config.yml, 本地新建 config/local.yml
2. export GO_WEIXIN_WORKDIR = <project_path>
3. go mod tidy
4. go run main.go


## 功能说明

1. cli 命令行工具
```

go run main.go migreate  # 数据库初始化 需要创建正式数据库以及一个测试数据库
go run main.go createsuperuser username password  # 创建超级用户
go run main.go refresh_permission  # 更新用户权限

```
2. swagger文档
```
# install swag
go get -u github.com/swaggo/swag/cmd/swag
# 自动生成文档
swag init
```

3. start server
```
go run -tags=doc main.go  // 带doc模式启动
go build main.go  // 无 doc
```



## 开发计划 
- [x] 通用组件
    - [x] swagger文档
    - [x] 配置管理
    - [x] zap日志服务
    - [x] 集中式err处理
    - [x] 命令行工具
    - [x] docker启动脚本
    - [ ] error邮件通知

- [ ] 系统用户管理
    - [ ] 登录
        - [x] 账号密码登录
        - [ ] 微信扫码登录
        - [ ] 手机号登录
        - [ ] 邮箱登录
    - [ ] 注册
        - [ ] 手机验证码认证
    - [ ] 修改密码
        - [ ] 手机验证码
    - [x] 用户认证 支持 `jwt`和`session`
    - [x] 权限管理 
    通过`casbin`实现`restful`风格权限管理

 