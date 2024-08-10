# moment-service

## Features

- post
- comment
- like

## Layer

The application with the following layers:

- 1. gRPC Layer: Responsible for handling incoming gRPC requests and returning responses.
- 2. Service Layer: Responsible for handling business logic and communicating with the data layer.
- 3. Data Layer: Responsible for handling communication with the database or nosql or other external service.

## Run

```bash
# 运行
make run
```

## Develop

如果中间有修改 proto 文件，需要重新生成 `*.pb.go` 和 `*_grpc.pb.go`

```bash
# 生成 *.pb.go
make proto

# 生成 *_grpc.pb.go
make grpc

# 生成文档 在 api/index.html 可以直接浏览器查看
make doc

make run
```

## Reference

- https://medium.com/@leoantony102/how-i-made-twitter-back-end-57addbaa14f5
- https://mp.weixin.qq.com/s/ov1UPkhjIti0QuHdxm2t9Q
- https://dsysd-dev.medium.com/stop-using-dtos-in-go-its-not-java-96ef4794481a
