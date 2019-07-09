## 微信公众号golang-sdk

## project start

1. 参考config/example_config.yml, 本地新建 config/local.yml
2. export GO_WEIXIN_WORKDIR = <project_path>
3. go mod tidy
4. go run main.go


## 功能说明

1. cli
```
cd cli/  && go build cli  
./cli migreate  # 数据库初始化 需要创建正式数据库以及一个测试数据库
./cli createsuperuser username password  # 创建超级用户
```




## 开发计划 

- [ ] 系统用户管理
    - [ ] 登录
    - [ ] 注册
    - [ ] 修改密码
    - [ ] 权限管理
 