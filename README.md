
## 简介

[在线访问](https://www.flywithme.top)

[前端源码地址](https://github.com/flywithbug-vue-admin/vue-eladmin)

#### 启动须知
  1. 参考 ``` config_simple.json ``` 添加配置文件（config.json）
  2. mongo数据库有两个：业务服务（doc_manager）和监控(monitor)
  3. 配置好启动文件之后，运行 ``` go run main.go ``` 即可启动服务



- **go主要框架** 
``` 
  github.com/gin-gonic/gin
  github.com/dgrijalva/jwt-go
  github.com/flywithbug/log4go
  gopkg.in/mgo.v2
```

-------------------------------
## 功能说明 (vue & go)
- [x] 登录 / 注销

- [x] 权限验证
  - [x] 页面权限  

- [x] 全局功能
  - [x] 国际化多语言
  - [x] 动态侧边栏（支持多级路由嵌套）
  - [x] 快捷导航(标签页)
  - [x] ScreenFull全屏
  - [x] 自适应收缩侧边栏




- [x] 系统管理
  - [x] 用户管理
  - [x] 角色管理
  - [x] 权限管理
  - [x] 菜单管理
  
- [x] 元数据 
  - [x] 应用管理
    - [x] 管理员权限
    - [x] 修改权限
  - [x] 版本管理
  
- [ ] 系统监控
  - [x] 操作日志
  - [x] 验证码
  - [ ] 服务器监控
  
  
- [ ] 开发工具
  - [ ] 数据模型 
  - [ ] API管理
  
  
### 用户权限设计 
![user_permission](/static/user_permission.png)
  
### 页面示例
![frontend](/static/dashboard.png)   




### Server项目布局

```
.
├── common
│   ├── com_definition.go
│   ├── common.go
│   ├── compare.go
│   └── compare_test.go
├── config
│   └── config.go
├── core
│   ├── errors
│   │   ├── errors.go
│   │   ├── errors_test.go
│   │   ├── reporter.go
│   │   └── reporter_test.go
│   ├── jwt
│   │   └── jwt.go
│   └── mongo
│       ├── Increment.go
│       ├── db.go
│       └── db_test.go
├── key_sources
│   ├── private_key
│   └── public_key.pub
├── model
│   ├── app_version.go
│   ├── application.go
│   ├── login.go
│   ├── model_func.go
│   ├── model_test.go
│   ├── response.go
│   ├── role.go
│   ├── user.go
│   └── user_role.go
└── server
    ├── handler
    │   ├── app_handler.go
    │   ├── app_version_handler.go
    │   ├── file_handler.go
    │   ├── html_handler.go
    │   ├── index.go
    │   ├── para_model.go
    │   ├── router.go
    │   └── user_handler.go
    ├── middleware
    │   ├── authentication.go
    │   ├── cookie.go
    │   └── logger.go
    ├── server.go
    └── web_server.go

11 directories, 37 files

``` 
