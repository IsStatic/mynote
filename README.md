# cloudNoteBygrpc 
基于grpc的云笔记云笔记

# 项目结构介绍

### 

| idl        | proto接口定义文件      |                                                              | 文档/子目录介绍                                              |
| ---------- | ---------------------- | ------------------------------------------------------------ | ------------------------------------------------------------ |
| pkg        | constants              | 常量                                                         |                                                              |
| errno      | 错误码                 | 关于错误码的讨论                                             |                                                              |
| tracer     | Jarger 初始化          |                                                              |                                                              |
| cmd        | api                    | demoapi服务的业务代码                                        | handlers : 封装了 api 的业务逻辑rpc : 封装了调用其它 rpc 服务的逻辑 |
| note       | demonote服务的业务代码 | dal : 封装了数据库的访问逻辑service: 封装了业务逻辑rpc : 封装了调用其它 rpc 服务的逻辑pack : 数据打包/处理 |                                                              |
| user       | demouser服务的业务代码 |                                                              |                                                              |
| etcdserver |                        | 封装了服务注册和服务发现                                     |                                                              |

### 

# 使用

## 部署docker

在项目目录下运行

```
docker-compose up
```

运行User服务

```
cd cmd/user
sh build.sh
sh output/bootstrap.sh
```

运行note服务

```
cd cmd/note
sh build.sh
sh output/bootstrap.sh
```

运行api服务

```
cd cmd/api
chmod +x run.sh
./run.sh
```
