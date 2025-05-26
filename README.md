# 项目介绍

## 如何启动
打开终端运行 make run，服务端起起来后可以运行客户端测试userservice，目前测试数据是已经注册过的，由幂等性可以得到信息：用户已注册过。systemservice用service目录下单元测试测 <br>
打开另一个终端运行 make pprof，下载Graphviz后在跳转到的浏览器8080端口页面可以看见可视化性能分析

## 有关说明
client为客户端，测试userservice <br>
cmd为服务端启动入口 <br>
middelware设置两个拦截器，为了方便测试并没有添加到grpc中 <br>
model初始化豆包向量化模型 <br>
pb存放grpc生成文件 <br>
service为业务具体代码 <br>
utils为工具包 <br>
使用user.sql初始化数据库 <br>



