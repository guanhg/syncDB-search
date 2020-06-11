### 同步Mysql表数据到Elastic
使用ali的Canal数据库同步工具，把Mysql数据的增删改同步到RabbitMQ，然后从MQ中拿取消息同步到Elastic中
 - Mysql --> Canal --> RabbitMq --> Elastic

1. 依赖安装 [Canal](https://github.com/alibaba/canal) 、RabbitMq、MySql、ElasticSearch 7.5
2. 编辑config/config.json配置文件
3. main目录中有3个文件 
    - syncDb2Mq.go 利用Canal工具同步表数据到Mq中
    - syncElastic4Mq.go 消费Mq中消息，更新Elastic
    - webApplication.go 一些自定义搜索应用接口







