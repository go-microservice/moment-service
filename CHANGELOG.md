
## 2023-10-19

- 增加清理函数 cleanup(main.go)
- cache中的redis客户端使用注入的方式代替使用全局变量
- 优化router, 文档路由和pprof增加开关控制
- 增加更多注释(post)
- 生成结构体说明文档，使用 make doc

## 2023-09-13

- 升级eagle到v1.8.0
- 修改post分页参数page_token和page_size，参考自: https://google.aip.dev/158