# moment-service


## Features

- post
- comment
- like

## Layer

The application with the following layers:

- 1. HTTP Layer: Responsible for handling incoming HTTP requests and returning responses.
- 2. Service Layer: Responsible for handling business logic and communicating with the data layer.
- 3. Data Layer: Responsible for handling communication with the database.

## Run

```bash
go build
./moment-service -c=config
```

## Reference

- https://medium.com/@leoantony102/how-i-made-twitter-back-end-57addbaa14f5
- https://mp.weixin.qq.com/s/ov1UPkhjIti0QuHdxm2t9Q
- https://dsysd-dev.medium.com/stop-using-dtos-in-go-its-not-java-96ef4794481a
