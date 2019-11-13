# 微信SDK


## 结构
- component 开放平台


## 依赖

### 工具

- redis（用于缓存获取的access_token等）


### 依赖库

- [go-redis](https://github.com/go-redis/redis)，redis操作


## feature

- debug模式和release模式
- debug模式控制台直接输出请求日志，结果
    - 方法调用记录
    - http请求记录，url，参数
    - http返回结果记录
    - 方法调用结果
- 缓存interface
- 提供几种缓存方式
    - go-cache
    - redis
    - merchant
    - 文件