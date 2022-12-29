# Gin Frame

Gin框架脚手架 by liwei

### 框架说明：

主要目录说明:

    conf：用于存储配置文件
    middleware：应用中间件
    models：数据操作层
    pkg：第三方包
    routers：路由逻辑处理
            -- api 存放逻辑处理层代码(controller)
    services:服务层：存放一些可复用的处理逻辑，作为一个服务单元

### 目录结构

```
├── conf                            // 项目配置                       
│   ├── dev
│       └── app.toml                // 测试环境
│   └── prod
│       └── app.toml                // 正式环境
├── datasource 
│   ├── dbhelper.go                 // 连接数据库
|   ├── hive.go                     // 连接hive
|   └── redishelper.go              // 连接redis
├── grpc_pb                         // grpc pb库
├── middleware                      // 中间件
│   ├── access                 
│   │   ├── jwt.go
│   │   ├── csrf.go
│   │   ├── recover.go
├── models                          // 数据操作层
│   ├── demo.go                     // 业务model demo
│   └── models.go                   // 初始化一些连接资源
├── pkg                             // 第三方包
│   ├── curl                        // http请求方法
│   ├── e                           // 错误码
│   ├── httpx                       // http请求方法
│   ├── mailx                       // 发邮件
│   └── setting                     // 加载配置文件
├── routers                         // 路由处理层
│   ├── api
|   │   └── v1                      // controller
|   │       └── demo.go  
│   └── router.go                   // 路由文件
├── go.mod
├── go.sum
├── main.go                         // 主程序入口
└── README.MD

```

用到的库:

    Gin:github.com/gin-gonic/gin
    日志：go.uber.org/zap
    配置文件管理：github.com/spf13/viper
    Redis库：https://github.com/go-redis/redis/v8
    orm: gorm.io/gorm
    JWT: github.com/dgrijalva/jwt-go


使用流程：
    
    1.routers/api/v1 新建文件，开始自己的controller
    2.routers/router.go 添加路由
    3.models 新建文件，添加数据存取方法
    
启动脚本：
    
    trunk: ./main 
    online: ./main --env=pro --config_dir=xxx

demo数据表
    
    CREATE TABLE `blog_article` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `tag_id` int unsigned DEFAULT '0' COMMENT '标签ID',
    `title` varchar(100) DEFAULT '' COMMENT '文章标题',
    `desc` varchar(255) DEFAULT '' COMMENT '简述',
    `content` text,
    `created_on` int DEFAULT NULL,
    `created_by` varchar(100) DEFAULT '' COMMENT '创建人',
    `modified_on` int unsigned DEFAULT '0' COMMENT '修改时间',
    `modified_by` varchar(255) DEFAULT '' COMMENT '修改人',
    `deleted_on` int unsigned DEFAULT '0',
    `state` tinyint unsigned DEFAULT '1' COMMENT '状态 0为禁用1为启用',
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb3 COMMENT='文章管理';

    
