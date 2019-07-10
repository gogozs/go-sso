## 微信平台 golang-sdk

## project start

1. 参考config/example_config.yml, 本地新建 config/local.yml
2. export GO_WEIXIN_WORKDIR = <project_path>
3. go mod tidy
4. go run main.go


## 功能说明

1. cli 命令行工具
```
cd cli/  && go build cli  
./cli migreate  # 数据库初始化 需要创建正式数据库以及一个测试数据库
./cli createsuperuser username password  # 创建超级用户
```




## 开发计划 
**参考** [公众号官方文档](https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140842)

- [ ] 系统用户管理
    - [ ] 登录
        - [x] 账号密码登录
        - [ ] 微信授权登录
    - [ ] 注册
    - [ ] 修改密码
    - [ ] 权限管理

- [ ] SDK功能
    - [ ] 消息管理
    - [ ] 图文消息留言管理
    - [ ] 菜单管理
    - [ ] 素材管理
    - [ ] 数据统计
    - [ ] 帐号管理
 