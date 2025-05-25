# taxin
同学你好，通过一面的筛选，恭喜您进入塔心科技的go实习生二轮面试。以下是二轮面试的题目及要求，如果对题意有不清晰的地方，可以随时联系我。
项目实践（二面）
- 安装中间件
    - Redis
    - PostgreSql (需要安装向量插件)
    - jaeger
      所用到的中间件可以本地安装，也可以使用docker或者docker-compose
- 业务逻辑
    - 使用uber的fx框架搭建grpc server
    - 创建proto文件user.proto，写一个 user service,user service中要写3个方法。
        - 用户注册    ( 需要实现幂等 )
        - 用户登录  （ 需要用到jwt，返回access_token ）
        - 获取用户信息 (需要登录用户才可以操作，只能获取自己的用户信息)
        - user表字段说明:
            - Id: 主键id
            - user_id: 用户分布式id
            - password: 用户密码
            - username：用户名
            - like: 喜好
            - like_embedding: 喜好的词嵌入向量值 (可以调用任意平台接口去获取embedding)
            - create_at: 创建时间
            - update_at: 更新时间
        - 用户表 user.sql 表结构文件，需要写在代码目录中。

    - 创建proto文件system.proto,写一个system service.包含一个方法，SendFile。
        -  读取一个本地文件（可以是音频，视频，文本）以流的形式返回。
- 单元测试
    - 为SendFile的service写测试用例，可以只写一个断言为发送出的字节流内容和发送文件一致。
-  链路追踪
- 对每一个请求使用jeager实现链路追踪，至少要在tag中记录用户id，方便通过user_id查询。
-  需要有方法查看你的go执行性能（比如使用 pprof）
-  写一个部署脚本，可以是shell，python，makefile。能通过脚本直接部署你的程序。（选作）
- 关于代码提交&时间限制
    - 可以直接给git或者gitee地址 （代码中注意设置 git ignore 文件，不要把重要的key泄漏）
    - 时间限制为48小时，如果时间不充分，可以联系增加答题时间。
- 幂等用unique索引实现