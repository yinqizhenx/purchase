## 写数据库同时发mq消息事务一致性的一种解决方案
1、begin tx 开启本地事务

2、do work 执行业务操作

3、insert message 向同实例消息库插入消息

4、end tx 事务提交

5、send message 网络向 server 发送消息 

6、reponse server 回应消息

7、delete message 如果 server 回复成功则删除消息

8、scan messages 补偿任务扫描未发送消息（使用定时任务）

9、send message 补偿任务补偿消息

10、delete messages 补偿任务删除补偿成功的消息 

11、消息发送加上重试

12、消息消费需要做好幂等