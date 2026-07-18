# BlogX 后端代码

这是Fengfeng博客10代项目



最开始做七代博客的时候，是刚开始做go开发不到半年

表结构的搭建并不算是从0到1，没有原型，做前端页面的时候，一会儿一个想法

但是好在技术栈全面，和文档记录详细，还是帮助了不少人成功找到了自己心仪的工作



经过两年的技术沉淀，和gvb系列，gvd，fim，fai项目的授课经验

决定再从0到1开发一个博客系统

为此，我给它取了一个响当当的名字  BlogX，意为“博客十代”





## 项目需求

1. 管理员可以切换运行模式，社区、博客模式
2. 网站配置要尽可能全，直接运行不需要改代码就能变成自己的
3. 文章管理需要有人工审核机制，可以有效保障内容安全，同时也留了一个免审核功能
4. 私信功能全面升级，用户可给喜欢的博主发私信，好友关系分为 陌生人 已关注  粉丝  好友
5. 消息通知更加全面，包括评论，赞和收藏，系统通知
6. 收藏夹、文章分类功能，可以通过此功能整理同类型文章
7. 个人主页，每个用户都可以选择不同的样式
8. 用户发布文章，可以开启是否可评论
9. 登录注册功能，可定制化是否使用qq登录，邮箱注册，是否显示图形验证码
10. 文章最少只需要标题和正文，可大大简化发布流程
11. 要接入AI，要有差异化

    可以做AI知识库，根据用户输入的内容，在全站文章中AI匹配

    还可以在发布文章的时候，分析文章适合的标签，分类，标题
12. 要能实现前端发版之后，线上能够自动刷新
13. 日志的记录要尽可能全面，日志方法的使用要简单
14. 数据库优化方面，索引优化，主从同步、读写分离，水平分表
15. 可以不需要es依赖，对应的搜索功能则需要对应的服务降级
16. 一键部署，部署流程要尽可能简单



## 项目相关技术栈

后端：

gin，gorm，redis，mysql，websocket， SSE，elasticsearch，docker-compose

前端：

vue3/ts，arco-desgin

## 后端默认用户名密码
```text
admin/1234
```

## 运行项目注意事项

1. 先确认 `settings.yaml` 配置正确，重点检查 MySQL、Redis、ES 和系统端口。

   ```yaml
   system:
     port: 8080
   db:
   - user: root
     password: "root"
     host: 127.0.0.1
     port: 3306
     db: blogx
   redis:
     addr: 127.0.0.1:6379
   es:
     enable: false
   river:
     enable: false
   ```

2. MySQL 需要先创建数据库。

   ```sql
   create database blogx charset utf8mb4;
   ```

3. 第一次运行前先执行数据库迁移。

   ```bash
   go run main.go -db
   ```

4. 创建管理员用户可以使用命令行工具。

   ```bash
   go run main.go -t user -s create
   ```

   角色选择 `1` 为管理员。

5. 启动后端服务。

   ```bash
   go run main.go
   ```

6. ES 不是必需依赖。没有 ES 时，将 `es.enable` 和 `river.enable` 都设置为 `false`，搜索相关功能会走降级逻辑。

7. Redis 是必需依赖，用于验证码、JWT 黑名单、缓存等功能。运行项目时需要保证 Redis 服务可连接。

8. 如果端口被占用，修改 `settings.yaml` 中的 `system.port`，或先停止占用端口的进程。

9. 前端默认请求后端地址由前端项目 `.env` 中的 `VITE_SERVER_URL` 控制，本地联调时需要和后端端口保持一致。

## 如何学习这个项目

1. 先从 `main.go` 看启动流程，理解配置读取、日志、IP 库、数据库、Redis、ES、命令行参数、定时任务和路由注册的初始化顺序。

2. 再看 `settings.yaml` 和 `conf/` 目录，掌握项目有哪些可配置能力，例如站点信息、登录方式、上传、七牛、邮箱、QQ 登录、AI、ES 和 river 同步。

3. 接着看 `router/` 目录，按业务模块梳理 API 分组，例如用户、文章、评论、消息、搜索、站点配置、后台数据统计等。

4. 学习接口实现时，可以按照 `router -> api -> service -> models` 的顺序阅读。比如文章模块可以从 `router/article_router.go` 开始，再看 `api/article_api/`，最后看 `models/article_model.go`。

5. 数据库部分重点看 `models/` 和 `flags/flag_db.go`，理解 GORM 模型、关联关系、钩子函数和自动迁移。

6. 登录鉴权建议按下面顺序学习：

   ```text
   api/user_api/pwd_login.go
   utils/jwts/
   middleware/auth_middleware.go
   middleware/captcha_middleware.go
   service/redis_service/redis_jwt/
   ```

7. 缓存和 Redis 可以从 `middleware/cache_middleware.go`、`core/init_redis.go`、`service/redis_service/` 入手，理解接口缓存、用户状态和文章计数缓存。

8. 评论、点赞、收藏、关注、私信这些功能适合用来学习业务建模。建议结合 `models/` 里的表结构和对应 API 一起看。

9. ES 和 MySQL 同步属于进阶内容。先理解普通 MySQL 查询和降级搜索，再看 `service/river_service/`、`service/es_service/` 和 `models/mappings/`。

10. 最后再看部署目录 `init/`，理解 docker-compose、MySQL 配置、Nginx 配置、前后端部署文件的关系。

建议每学一个模块，就用接口工具或前端页面实际跑一遍请求，同时观察数据库、Redis 和日志变化，这样比只看代码更容易理解完整链路。



## 你将学到

1. 使用ws实现在线用户的即时通讯
2. mysql和es的数据同步
3. qq登录的相关知识点
4. 编写docker-compose，实现便捷部署



## 表结构

![](https://image.fengfengzhidao.com/pic/20240926232427.png)



## 原型

![](https://image.fengfengzhidao.com/rj_0912/20240926144115.png)



![](https://image.fengfengzhidao.com/rj_0912/20240926144221.png)



![](https://image.fengfengzhidao.com/rj_0912/20240926144301.png)



![](https://image.fengfengzhidao.com/rj_0912/20240926144339.png)



![](https://image.fengfengzhidao.com/rj_0912/20240926144456.png)



![](https://image.fengfengzhidao.com/rj_0912/20240926144718.png)



![](https://image.fengfengzhidao.com/rj_0912/20240926144752.png)



![](https://image.fengfengzhidao.com/rj_0912/20240926144826.png)



![](https://image.fengfengzhidao.com/rj_0912/20240926144856.png)



![](https://image.fengfengzhidao.com/rj_0912/20240926145011.png)



![](https://image.fengfengzhidao.com/rj_0912/20240926144954.png)



![](https://image.fengfengzhidao.com/rj_0912/20240926145040.png)



![](https://image.fengfengzhidao.com/rj_0912/20240926145116.png)



![](https://image.fengfengzhidao.com/rj_0912/20240930110631.png)



![](https://image.fengfengzhidao.com/rj_0912/20240926145312.png)



![](https://image.fengfengzhidao.com/rj_0912/20240926145345.png)
