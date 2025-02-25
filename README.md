
      ├── app # 应用层
      │   ├── assembler # 服务领域服务与 presentation DTO 的转化
      │   ├── cqe       # 定义应用服务的 Command、Query、Event DTO
      │   ├── cron       # 定时任务调用，轻量，主体逻辑写在domain层
      │   ├── domain_event   # 事件订阅，轻量，主体逻辑写在domain层
      │   └── order_app_service.go
      ├── domain # 领域层
      │   ├── service     # 放置领域服务
      │   ├── repo        # 放置仓储
      │   ├── vo          # 放置值对象
      │   ├── event       # 领域事件的定义放在这里
      │   └── entity      # 放置聚合根、实体、值对象，枚举、事件的定义放在idl里
      ├── infra # 基础设施层层
      │   ├── cache          # 本地缓存、redis等
      │   ├── mq             # 主要是处理消息的发送
      │   ├── persistence    # 数据库相关
      │   │   ├── convertor # 领域对象与 PO 之间的转化
      │   │   ├── dal       # 对 PO 的定义和 db 的 crud 操作
      │   │   └── order_repo.go
      │   └── sal            # 对rpc、http服务的访问，通常是对领域服务接口的实现
      │       └── ...        # 一个文件对应领域服务里一个接口的实现
      ├── conf # 配置dir
      │   └── config.yml
      ├── idl  # 当前服务的idl，放在其他地方也行
      │   ├── base.proto
      │   └── exam.proto
      ├── script  # kitex生成
      │   ├── bootstrap.sh
      │   └── settings.py
      ├── README.md
      ├── handler.go
      ├── main.go
      ├── build.sh
      └── wire.go # 解决应用的依赖注入问题

## 特性
* 基于DDD的项目结构
* 实现“零延迟”的异步任务队列
* 实现事务传播机制
* 实现了基于kafka的延迟队列，死信队列，幂等消费，基于本地消息表的事务消息
* 实现基于saga的分布式事务
* 实现基于redis的分布式锁
